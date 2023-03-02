package cluster

import "github.com/A-walker-ninght/miniRedis/interface/resp"

func execSelect(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	return cluster.Exec(c, cmdArgs)
}
