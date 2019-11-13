package p2p

import (
	"fmt"
	"sync"
	"time"

	"github.com/bitmark-inc/bitmarkd/blockheader"

	"github.com/bitmark-inc/bitmarkd/blockdigest"
	"github.com/bitmark-inc/bitmarkd/util"
	"github.com/bitmark-inc/logger"
	peerlib "github.com/libp2p/go-libp2p-core/peer"
)

const (
	votingCycleInterval = 10 * time.Second
	votingQueryTimeout  = 5 * time.Second
)

//MetricsPeersVoting  is to get all metrics for voting
type MetricsPeersVoting struct {
	mutex      *sync.Mutex
	watchNode  *Node
	Candidates []*P2PCandidatesImpl
	Log        *logger.L
}

//NewMetricsPeersVoting return a MetricsPeersVoting for voting
func NewMetricsPeersVoting(thisNode *Node) MetricsPeersVoting {
	var mutex = &sync.Mutex{}
	metrics := MetricsPeersVoting{mutex: mutex, watchNode: thisNode, Log: logger.New("votingMetrics")}
	metrics.UpdateCandidates()
	return metrics
}

//UpdateCandidates update Candidate by registered peer
func (p *MetricsPeersVoting) UpdateCandidates() {
	var Candidates []*P2PCandidatesImpl
	p.mutex.Lock()
	for peerID, registered := range p.watchNode.Registers {
		util.LogInfo(p.Log, util.CoWhite, fmt.Sprintf("UpdateCandidates:peerID:%v!", peerID.ShortString()))
		if registered && !util.IDEqual(p.watchNode.Host.ID(), peerID) { // register and not self
			peerInfo := p.watchNode.Host.Peerstore().PeerInfo(peerID)
			if len(peerInfo.Addrs) > 0 {
				Candidates = append(Candidates, &P2PCandidatesImpl{ID: peerID, Addr: peerInfo.Addrs[0]})
			} else {
				Candidates = append(Candidates, &P2PCandidatesImpl{ID: peerID})
			}
		}
	}
	p.Candidates = Candidates
	p.mutex.Unlock()
	util.LogInfo(p.Log, util.CoWhite, fmt.Sprintf("UpdateCandidates:%d Candidates!", len(Candidates)))
}

//Run  is a Routine to get peer info
func (p *MetricsPeersVoting) Run(args interface{}, shutdown <-chan struct{}) {
	log := p.Log
	delay := time.After(nodeInitial)
	//nodeChain:= mode.ChainName()
loop:
	for {
		log.Debug("waiting…")
		select {
		case <-shutdown:
			continue loop
		case <-delay: //update voting metrics
			delay = time.After(votingCycleInterval)
			p.UpdateCandidates()
			if nil == p.Candidates {
				continue loop
			}
			for _, peer := range p.Candidates {
				go func(id peerlib.ID) {
					util.LogInfo(p.Log, util.CoGreen, fmt.Sprintf("UpdateMetrics ID:%v!", id.ShortString()))
					height, err := p.watchNode.QueryBlockHeight(id)
					if err != nil {
						util.LogWarn(p.Log, util.CoRed, fmt.Sprintf("QueryBlockHeight Error:%v", err))
						return
					}
					digest, err := p.watchNode.RemoteDigestOfHeight(id, height)
					if err != nil {
						util.LogWarn(p.Log, util.CoRed, fmt.Sprintf("RemoteDigestOfHeight Error:%v", err))
						return
					}
					util.LogInfo(p.Log, util.CoWhite, fmt.Sprintf("Query Return height: %d candidates ID:%v", height, id.ShortString()))
					p.setMetrics(id, height, digest)
				}(peer.ID)
			}
		}
	}
}

func (p *MetricsPeersVoting) setMetrics(id peerlib.ID, height uint64, digest blockdigest.Digest) {
	for _, candidate := range p.Candidates {
		if util.IDEqual(candidate.ID, id) {
			p.mutex.Lock()
			localheight := blockheader.Height()
			respTime := time.Now()
			candidate.UpdateMetrics(id.String(), height, localheight, digest, respTime)
			p.mutex.Unlock()
			p.Log.Debugf("\x1b[32m: ID:%s, height:%d, digest:%s respTime:%d\x1b[0m",
				candidate.ID, candidate.Metrics.localHeight, candidate.Metrics.remoteDigestOfLocalHeight, candidate.Metrics.lastResponseTime.Unix())
			break
		}
	}
}

func (p *MetricsPeersVoting) allCandidates(
	f func(candidate *P2PCandidatesImpl),
) {
	for _, candidate := range p.Candidates {
		f(candidate)
	}
}
