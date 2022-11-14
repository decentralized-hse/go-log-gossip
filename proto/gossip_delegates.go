package proto

import (
	"sync"
	"time"
)

type GossipDelegate struct {
	// GossipStoreImpl implements the GossipStoreInterface
	GossipStoreImpl
	nodeId string
	// last gossip time
	lastGossipTsLock sync.Mutex
	lastGossipTs     time.Time
	// channel to receive state change events
	stateEvent chan types.StateEvent
	// current State object
	currentState     state.State
	currentStateLock sync.Mutex
	// quorum timeout to change the quorum status of a node
	quorumTimeout            time.Duration
	timeoutVersion           uint64
	timeoutVersionLock       sync.Mutex
	nodeDownProbationManager probation.Probation
	quorumProvider           state.Quorum
	// ping is a callback function from Gossiper that uses memberlist
	// apis to ping a peer node
	ping func(types.NodeId, string) (time.Duration, error)
}

func (gd *GossipDelegate) InitGossipDelegate(
	genNumber uint64,
	selfNodeId types.NodeId,
	gossipVersion string,
	quorumTimeout time.Duration,
	clusterId string,
	selfClusterDomain string,
	ping func(types.NodeId, string) (time.Duration, error),
) {
	gd.GenNumber = genNumber
	gd.nodeId = string(selfNodeId)
	gd.stateEvent = make(chan types.StateEvent)
	gd.ping = ping
	// We start with a NOT_IN_QUORUM status
	gd.InitStore(
		selfNodeId,
		gossipVersion,
		types.NODE_STATUS_NOT_IN_QUORUM,
		clusterId,
		selfClusterDomain,
	)
	gd.quorumTimeout = quorumTimeout
	gd.nodeDownProbationManager = probation.NewProbationManager(
		"node-suspected-down-probation-manager",
		suspectNodeDownTimeout,
		gd.probationExpiredOnSuspectedDownNode,
	)
}
