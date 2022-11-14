package proto

import (
	ml "github.com/hashicorp/memberlist"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type GossipNode struct {
	Id string
	Ip string
}

type GossipNodeList []GossipNode

func (nodes GossipNodeList) Len() int {
	return len(nodes)
}

type GossiperImpl struct {
	mlConf     *ml.Config
	memberList *ml.Memberlist

	nodes          GossipNodeList
	name           string
	gossipInterval time.Duration

	shutDown         bool
	selfNodeId       string
	hasJoinedCluster bool

	selfClusterDomain string
}

type GossipIntervals struct {
	GossipInterval   time.Duration
	PushPullInterval time.Duration
	ProbeInterval    time.Duration
	ProbeTimeout     time.Duration
	QuorumTimeout    time.Duration
	SuspicionMult    int
}

func (gossiper *GossiperImpl) Init(
	ipPort string,
	selfNodeId string,
	genNumber uint64,
	gossipIntervals GossipIntervals,
	gossipVersion string,
	clusterId string,
	selfClusterDomain string) {
	gossiper.name = ipPort
	gossiper.shutDown = false

	gossiper.nodes = make(GossipNodeList, 0)
	gossiper.gossipInterval = gossipIntervals.GossipInterval

	mlConf := ml.DefaultLANConfig()

	ip, port, _ := net.SplitHostPort(ipPort)

	port64, _ := strconv.ParseInt(port, 10, 64)

	nodeName := selfNodeId + gossipVersion
	mlConf.Name = nodeName
	mlConf.BindAddr = ip
	mlConf.BindPort = int(port64)

	//mlConf.Delegate = ml.Delegate(gossiper)
	//mlConf.Events = ml.EventDelegate(gossiper)
	//mlConf.Alive = ml.AliveDelegate(gossiper)
	//mlConf.Merge = ml.MergeDelegate(gossiper)

	gossiper.mlConf = mlConf
	gossiper.selfNodeId = selfNodeId
	gossiper.selfClusterDomain = selfClusterDomain
	rand.Seed(time.Now().UnixNano())
}
