package cluster

import (
	"context"
	"errors"
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"github.com/A-walker-ninght/miniRedis/lib/utils"
	client2 "github.com/A-walker-ninght/miniRedis/resp/client"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
	"strconv"
)

func (cluster *ClusterDatabase) getPeerClient(peer string) (*client2.Client, error) {
	pool, ok := cluster.peerconncection[peer]
	if !ok {
		return nil, errors.New("connection pool not found")
	}
	object, err := pool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}
	client, ok := object.(*client2.Client)
	if !ok {
		return nil, errors.New("wrong client type")
	}
	return client, nil
}

func (cluster *ClusterDatabase) putbackPeerClient(peer string, client *client2.Client) error {
	pool, ok := cluster.peerconncection[peer]
	if !ok {
		return errors.New("connection pool not found")
	}
	return pool.ReturnObject(context.Background(), client)
}

// 转发
func (cluster *ClusterDatabase) relay(peer string, c resp.Connection, args [][]byte) resp.Reply {
	if peer == cluster.self {
		return cluster.Exec(c, args)
	}
	peerClient, err := cluster.getPeerClient(peer)
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}
	defer func() {
		_ = cluster.putbackPeerClient(peer, peerClient)
	}()
	for {
		repl := peerClient.Send(utils.ToCmdLine("select", strconv.Itoa(c.GetDBIndex())))
		if len(repl.ToBytes()) != 0 {
			break
		}
	}
	return peerClient.Send(args)
}

// 广播
func (cluster *ClusterDatabase) broadcast(c resp.Connection, args [][]byte) map[string]resp.Reply {
	results := make(map[string]resp.Reply)

	for _, node := range cluster.nodes {
		result := cluster.relay(node, c, args)
		results[node] = result
	}
	return results
}
