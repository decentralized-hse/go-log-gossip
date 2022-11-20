package proto

import (
	"bufio"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type GossipNode struct {
	Id      string
	Address string
}

type Gossip interface {
	Send(address string, message string)
	Publish(message string)
}

type GossipNodeList []GossipNode

func (nodes GossipNodeList) Len() int {
	return len(nodes)
}

// GossiperImpl implements Gossiper
type GossiperImpl struct {
	ip             string
	port           int64
	nodes          GossipNodeList
	name           string
	gossipInterval time.Duration

	selfNodeId string

	selfClusterDomain string
}

func (gossiper *GossiperImpl) Init(
	ipPort string,
	selfNodeId string,
	selfClusterDomain string) {
	gossiper.name = ipPort

	gossiper.nodes = make(GossipNodeList, 0)

	ip, port, _ := net.SplitHostPort(ipPort)

	port64, _ := strconv.ParseInt(port, 10, 64)

	gossiper.ip = ip
	gossiper.port = port64

	gossiper.selfNodeId = selfNodeId
	gossiper.selfClusterDomain = selfClusterDomain
	rand.Seed(time.Now().UnixNano())
}

func (gossiper *GossiperImpl) Send(address string, message string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString(message)

	if err != nil {
		log.Fatal(err)
	}
}

func (gossiper *GossiperImpl) Publish(message string) {
	for _, node := range gossiper.nodes {
		gossiper.Send(node.Address, message)
	}
}
