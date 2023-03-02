package cluster

import (
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
)

// del k1 k2 k3 k4 k5
func Del(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	repies := cluster.broadcast(c, cmdArgs)
	var errReply reply.ErrorReply
	var deleted int64 = 0

	for _, r := range repies {
		if reply.IsErrReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
		intReply, ok := r.(*reply.IntReply)
		if !ok {
			return reply.MakeStandardErrReply("error")
		}
		deleted += intReply.Code
	}

	if errReply == nil {
		return reply.MakeOKReply()
	}
	return reply.MakeStandardErrReply("error: " + errReply.Error())
}
