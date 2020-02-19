package p2p

import (
	"context"
	"fmt"
	"time"

	proto "github.com/golang/protobuf/proto"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
	peerlib "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	tls "github.com/libp2p/go-libp2p-tls"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/bitmark-inc/bitmarkd/fault"
	"github.com/bitmark-inc/bitmarkd/messagebus"
	"github.com/bitmark-inc/bitmarkd/util"
)

//Setup setup a node
func (n *Node) Setup(configuration *Configuration, version string, dnsPeerOnly dnsOnlyType) error {
	globalData.Version = version
	if nodeType(configuration.NodeType) == ClientNode {
		globalData.NodeType = ClientNode
		n.Log.Info("The NodeType is ClientNode")
	} else {
		globalData.NodeType = ServerNode
		n.Log.Info("The NodeType is ServerNode")
	}

	globalData.dnsPeerOnly = dnsPeerOnly
	listenIPPorts := makeDualStackAddrs(configuration.Listen)
	if len(listenIPPorts) == 0 {
		return fault.AddressIsNil
	}
	maAddrs := IPPortToMultiAddr(listenIPPorts)
	n.Registers = make(map[peerlib.ID]RegisterStatus)
	prvKey, err := util.DecodePrivKeyFromHex(configuration.PrivateKey) //Hex Decoded binaryString
	if err != nil {
		return err
	}

	n.PrivateKey = prvKey
	n.NewHost(globalData.NodeType, maAddrs, n.PrivateKey)

	if n.NodeType != ClientNode {
		n.setAnnounce(configuration.Announce)
	}

	go n.listen(configuration.Announce)
	n.MetricsNetwork = NewMetricsNetwork(n.Host, n.Log)

	//Start Broadcsting
	ps, err := pubsub.NewGossipSub(context.Background(), n.Host)
	if err != nil {
		return err
	}
	n.Multicast = ps
	sub, err := n.Multicast.Subscribe(TopicMulticasting)
	if err != nil {
		return err
	}
	go n.SubHandler(context.Background(), sub)

	globalData.initialised = true
	return nil
}

// NewHost create a NewHost according to nodetype
func (n *Node) NewHost(nodetype nodeType, listenAddrs []ma.Multiaddr, prvKey crypto.PrivKey) error {
	options := []libp2p.Option{libp2p.Identity(prvKey), libp2p.Security(tls.ID, tls.New)}
	if ClientNode != nodetype {
		options = append(options, libp2p.ListenAddrs(listenAddrs...))
	}
	newHost, err := libp2p.New(context.Background(), options...)
	if err != nil {
		return err
	}
	n.Host = newHost
	for _, a := range newHost.Addrs() {
		n.Log.Info(fmt.Sprintf("Host Address: %s/%v/%s\n", a, nodeProtocol, newHost.ID()))
	}
	return nil
}

//setAnnounce: Set Announce address in Routing
func (n *Node) setAnnounce(announceAddrs []string) {
	maAddrs := IPPortToMultiAddr(announceAddrs)
	fullAddr := announceMuxAddr(maAddrs, nodeProtocol, n.Host.ID())
	n.Announce = fullAddr
	util.LogInfo(n.Log, util.CoReset, fmt.Sprintf("setAnnounce:%v", util.PrintMaAddrs(fullAddr)))
	byteMessage, err := proto.Marshal(&Addrs{Address: util.GetBytesFromMultiaddr(fullAddr)})
	param0, idErr := n.Host.ID().Marshal()

	if nil == err && nil == idErr {
		messagebus.Bus.Announce.Send("self", param0, byteMessage)
	}
}

func (n *Node) listen(announceAddrs []string) {
	maAddrs := IPPortToMultiAddr(announceAddrs)
	shandler := NewListenHandler(n.Host.ID(), n, n.Log)
	n.Host.SetStreamHandler("p2pstream", shandler.handleStream)
	n.Log.Infof("A servant is listen to %s", util.PrintMaAddrs(maAddrs))
	// Hang forever
	<-make(chan struct{})
}

func (n *Node) addRegister(id peerlib.ID) {
	n.Lock()
	status, ok := n.Registers[id]
	if ok {
		status.Registered = true
		status.RegisterTime = time.Now()
		n.Unlock()
		return
	}
	n.Registers[id] = RegisterStatus{Registered: true, RegisterTime: time.Now()}
	util.LogInfo(n.Log, util.CoGreen, fmt.Sprintf("addRegister ID:%s Registered:%v time:%v", id.ShortString(), n.Registers[id].Registered, n.Registers[id].RegisterTime.String()))
	n.Unlock()
}

//unRegister unRegister change a peers's  Registered status  to false,  but it doe not not delete the register in the Registers
func (n *Node) unRegister(id peerlib.ID) {
	n.Lock()
	status, ok := n.Registers[id]
	if ok { // keep RegisterTime for last record purpose
		status.Registered = false
		n.Unlock()
	}
	util.LogInfo(n.Log, util.CoGreen, fmt.Sprintf("unRegister ID:%s Registered:%v time:%v", id.ShortString(), n.Registers[id].Registered, n.Registers[id].RegisterTime.String()))
}

//delRegister delete a Registerer  in the Registers map
func (n *Node) delRegister(id peerlib.ID) {
	n.Lock()
	_, ok := n.Registers[id]
	if ok { // keep RegisterTime for last record purpose
		delete(n.Registers, id)
		n.Unlock()
	}
}

//IsRegister if given id has a registered stream
func (n *Node) IsRegister(id peerlib.ID) (registered bool) {
	n.Lock()
	if status, ok := n.Registers[id]; ok && status.Registered {
		registered = true
	}
	n.Unlock()
	return
}

//IsExpire is the register expire
func (n *Node) IsExpire(id peerlib.ID) bool {
	if status, ok := n.Registers[id]; ok && status.Registered {
		expire := status.RegisterTime.Add(registerExpireTime)
		passInterval := time.Since(expire)
		if passInterval > 0 { // expire
			return true
		}
	}
	return false
}

//updateRegistersExpiry mark Registered false when time is expired
func (n *Node) updateRegistersExpiry() {
	for id, status := range n.Registers {
		if n.IsExpire(id) { //Keep time for record of last registered time
			n.Lock()
			status.Registered = false
			n.Unlock()
			util.LogDebug(n.Log, util.CoWhite, fmt.Sprintf("IsExpire ID:%v is expire", id.ShortString()))
		}
	}
}
