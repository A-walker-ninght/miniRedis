package cluster

import (
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
)

func flushdb(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	repies := cluster.broadcast(c, cmdArgs)
	var errReply reply.ErrorReply
	for _, r := range repies {
		if reply.IsErrReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
	}
	if errReply == nil {
		return reply.MakeOKReply()
	}
	return reply.MakeStandardErrReply("error: " + errReply.Error())
}
