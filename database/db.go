package database

import (
	"github.com/A-walker-ninght/miniRedis/datastruct/dict"
	"github.com/A-walker-ninght/miniRedis/interface/database"
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
	"strings"
)

type DB struct {
	index  int
	data   dict.Dict
	addAof func(line CmdLine)
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply
type CmdLine = [][]byte

func makeDB() *DB {
	return &DB{
		data:   dict.NewDict(),
		addAof: func(line CmdLine) {},
	}
}

func (db *DB) Exec(c resp.Connection, line CmdLine) resp.Reply {
	// ping set setnx
	cmdName := strings.ToLower(string(line[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeStandardErrReply("ERR unknown command " + cmdName)
	}
	if !validArity(cmd.arity, line) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	fun := cmd.exector
	return fun(db, line[1:])
}

// arity < 0: 变长的参数，例如：exists k1, k2, k3 ...
// arity = -2表示变长
func validArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}

func (db *DB) GetEntry(key string) (*database.DataEntry, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entry, _ := raw.(*database.DataEntry)
	return entry, true
}

func (db *DB) PutEntry(key string, val *database.DataEntry) int {
	return db.data.Put(key, val)
}

func (db *DB) PutIfExists(key string, val *database.DataEntry) int {
	return db.data.PutIfExists(key, val)
}

func (db *DB) PutIfAbsent(key string, val *database.DataEntry) int {
	return db.data.PutIfAbsent(key, val)
}

func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) (deleted int) {
	for _, key := range keys {
		_, existed := db.data.Get(key)
		if existed {
			db.data.Remove(key)
			deleted++
		}
	}
	return
}

func (db *DB) Clear() {
	db.data.Clear()
}
