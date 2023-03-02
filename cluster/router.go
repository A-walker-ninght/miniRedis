package cluster

import "github.com/A-walker-ninght/miniRedis/interface/resp"

func makeRouter() map[string]CmdFunc {
	router := make(map[string]CmdFunc)
	router["exists"] = defaultFunc
	router["get"] = defaultFunc
	router["set"] = defaultFunc
	router["setnx"] = defaultFunc
	router["getset"] = defaultFunc
	router["type"] = defaultFunc
	router["ping"] = ping
	router["rename"] = Rename
	router["renamenx"] = Rename
	router["flushdb"] = flushdb
	router["del"] = Del
	router["select"] = execSelect
	return router
}

// GET key // set key val
func defaultFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	key := string(cmdArgs[1])
	peer := cluster.picker.PickNode(key)
	return cluster.relay(peer, c, cmdArgs)
}
