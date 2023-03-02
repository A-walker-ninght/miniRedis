package cluster

import (
	"context"
	"github.com/A-walker-ninght/miniRedis/config"
	database2 "github.com/A-walker-ninght/miniRedis/database"
	"github.com/A-walker-ninght/miniRedis/interface/database"
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"github.com/A-walker-ninght/miniRedis/lib/consistenthash"
	pool "github.com/jolestar/go-commons-pool/v2"
)

type ClusterDatabase struct {
	self            string
	nodes           []string
	picker          *consistenthash.NodeMap
	peerconncection map[string]*pool.ObjectPool //
	db              database.DataBase
}

func NewClusterDatabase() *ClusterDatabase {
	cluster := &ClusterDatabase{
		self:            config.Properties.Self,
		db:              database2.NewStandaloneDatabase(),
		picker:          consistenthash.NewNodeMap(nil, 6),
		peerconncection: make(map[string]*pool.ObjectPool),
	}
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, cluster.self)
	cluster.picker.AddNodes(nodes...)
	ctx := context.Background()

	for _, peer := range config.Properties.Peers {
		pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{
			Peer: peer,
		})
	}
	cluster.nodes = nodes
	return cluster
}

type CmdFunc func(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply

var router = makeRouter()

func (CD *ClusterDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	//TODO implement me
	panic("implement me")
}

func (CD *ClusterDatabase) Close() error {
	//TODO implement me
	panic("implement me")
}

func (CD *ClusterDatabase) AfterClientClose(c resp.Connection) {
	//TODO implement me
	panic("implement me")
}
