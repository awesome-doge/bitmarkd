// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package zmqutil

import (
	"github.com/bitmark-inc/bitmarkd/fault"
	"github.com/bitmark-inc/bitmarkd/util"
	zmq "github.com/pebbe/zmq4"
	"time"
)

// structure to hold a client connection
type Client struct {
	address string
	socket  *zmq.Socket
}

// create a cliet socket ususlly of type zmq.REQ or zmq.SUB
func NewClient(socketType zmq.Type, privateKey []byte, publicKey []byte, timeout time.Duration) (*Client, error) {
	socket, err := zmq.NewSocket(socketType)
	if nil != err {
		return nil, err
	}

	// set up as client
	socket.SetCurveServer(0)
	socket.SetCurvePublickey(string(publicKey))
	socket.SetCurveSecretkey(string(privateKey))

	socket.SetIdentity(string(publicKey)) // just use public key for identity

	// // basic socket options
	// socket.SetIpv6(true) // do not set here defer to connect
	// socket.SetRouterMandatory(0)   // discard unroutable packets
	// socket.SetRouterHandover(true) // allow quick reconnect for a given public key
	// socket.SetImmediate(false)     // queue messages sent to disconnected peer

	socket.SetReqCorrelate(1)
	socket.SetReqRelaxed(1)
	socket.SetSndtimeo(timeout)
	socket.SetRcvtimeo(timeout)
	socket.SetLinger(0)

	client := &Client{
		address: "",
		socket:  socket,
	}
	return client, nil
}

// disconnect old address and connect to new
func (client *Client) Connect(conn *util.Connection, serverPublicKey []byte) error {

	client.socket.SetCurveServerkey(string(serverPublicKey))

	connectTo, v6 := conn.CanonicalIPandPort("tcp://")

	// if already connected, disconnect first
	if "" != client.address {
		err := client.socket.Disconnect(client.address)
		if nil != err {
			return err
		}
	}
	client.address = ""

	// set IPv6 state before connect
	err := client.socket.SetIpv6(v6)
	if nil != err {
		return err
	}

	// new connection
	err = client.socket.Connect(connectTo)
	if nil != err {
		return err
	}
	client.address = connectTo

	return nil
}

func (client *Client) Reconnect() error {
	if "" == client.address {
		return nil
	}
	err := client.socket.Disconnect(client.address)
	if nil != err {
		return err
	}
	err = client.socket.Connect(client.address)
	if nil != err {
		return err
	}
	return nil
}

// disconnect old address and close
func (client *Client) Close() error {
	// if already connected, disconnect first
	if "" != client.address {
		client.socket.Disconnect(client.address)
	}
	client.address = ""

	// close socket
	err := client.socket.Close()
	client.socket = nil

	return err
}

// disconnect old addresses and close all
func CloseClients(clients []*Client) {
	for _, client := range clients {
		if nil != client {
			client.Close()
		}
	}
}

// send a message
func (client *Client) Send(items ...interface{}) error {
	if "" == client.address {
		return fault.ErrNotConnected
	}

	last := len(items) - 1
	for i, item := range items {

		flag := zmq.SNDMORE
		if i == last {
			flag = 0
		}
		switch item.(type) {
		case string:
			_, err := client.socket.Send(item.(string), flag)
			if nil != err {
				return err
			}
		case []byte:
			_, err := client.socket.SendBytes(item.([]byte), flag)
			if nil != err {
				return err
			}
		}
	}
	return nil
}

// receive a reply
func (client *Client) Receive(flags zmq.Flag) ([][]byte, error) {
	if "" == client.address {
		return nil, fault.ErrNotConnected
	}
	data, err := client.socket.RecvMessageBytes(flags)
	return data, err
}

// add to a poller
func (client *Client) Add(poller *zmq.Poller, state zmq.State) int {
	return poller.Add(client.socket, state)
}
