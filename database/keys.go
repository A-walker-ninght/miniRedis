package database

import (
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"github.com/A-walker-ninght/miniRedis/lib/utils"
	"github.com/A-walker-ninght/miniRedis/lib/wildcard"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
)

// DEL
func execDEL(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = string(arg)
	}
	deleted := db.Removes(keys...)
	if deleted > 0 {
		db.addAof(utils.ToCmdLine2("del", args...))
	}
	return reply.MakeIntReply(int64(deleted))
}

// EXISTS
func execExists(db *DB, args [][]byte) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, existed := db.GetEntry(key)
		if existed {
			result++
		}
	}
	return reply.MakeIntReply(result)
}

// KEYS *
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(args)
}

// FLUSHDB
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Clear()
	db.addAof(utils.ToCmdLine2("flushdb", args...))
	return reply.MakeOKReply()
}

// TYPE
func execType(db *DB, args [][]byte) resp.Reply {
	key1 := string(args[0])
	entry, exists := db.GetEntry(key1)
	if !exists {
		return reply.MakeStatusReply("none")
	}
	switch entry.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
	}
	return &reply.UnknowErrReply{}
}

// RENAME k1 k2
func execRename(db *DB, args [][]byte) resp.Reply {
	oldk := string(args[0])
	newk := string(args[1])

	entry, exists := db.GetEntry(oldk)
	if !exists {
		return reply.MakeStandardErrReply("no such key")
	}
	db.PutEntry(newk, entry)
	db.Remove(oldk)
	db.addAof(utils.ToCmdLine2("rename", args...))
	return reply.MakeOKReply()
}

// RENAMENX k1 k2
func execRenameNX(db *DB, args [][]byte) resp.Reply {
	oldk := string(args[0])
	newk := string(args[1])
	_, ok := db.GetEntry(newk)

	// 如果k2存在，则不修改
	if ok {
		return reply.MakeIntReply(0)
	}

	entry, exists := db.GetEntry(oldk)
	if !exists {
		return reply.MakeStandardErrReply("no such key")
	}
	db.PutEntry(newk, entry)
	db.Remove(oldk)
	db.addAof(utils.ToCmdLine2("renamenx", args...))
	return reply.MakeOKReply()
}

func init() {
	RegisterCommand("del", execDEL, -2)
	RegisterCommand("exists", execExists, -2)
	RegisterCommand("flushdb", execFlushDB, -1)
	RegisterCommand("type", execType, 2)
	RegisterCommand("rename", execRename, 3)
	RegisterCommand("type", execRenameNX, 3)
	RegisterCommand("keys", execKeys, 2)
}
