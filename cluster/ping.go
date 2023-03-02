package cluster

import "github.com/A-walker-ninght/miniRedis/interface/resp"

func ping(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.db.Exec(c, cmdArgs)
}
