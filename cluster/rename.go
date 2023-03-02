package cluster

import (
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
)

// rename k1 k2
func Rename(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeStandardErrReply("ERR wrong number args")
	}

	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])
	srcPeer := cluster.picker.PickNode(src)
	destPeer := cluster.picker.PickNode(dest)
	if srcPeer != destPeer {
		return reply.MakeStandardErrReply("ERR rename must wighin on peer")
	}
	return cluster.relay(srcPeer, c, cmdArgs)
}
