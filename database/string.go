package database

import (
	"github.com/A-walker-ninght/miniRedis/interface/database"
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
)

// GET
func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val, ok := db.GetEntry(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	bytes := val.Data.([]byte)
	return reply.MakeBulkReply(bytes)
}

// SET
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entry := &database.DataEntry{
		Data: val,
	}
	db.PutEntry(key, entry)
	return reply.MakeOKReply()
}

// SETNX
func execSetNX(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entry := &database.DataEntry{
		Data: val,
	}
	result := db.PutIfAbsent(key, entry)
	return reply.MakeIntReply(int64(result))
}

// GETSET
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	entry, ok := db.GetEntry(key)
	db.PutEntry(key, &database.DataEntry{Data: val})
	if !ok {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply(entry.Data.([]byte))
}

// STRLEN
func execStrLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entry, ok := db.GetEntry(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	bytes := entry.Data.([]byte)
	return reply.MakeIntReply(int64(len(bytes)))
}

func init() {
	RegisterCommand("get", execGet, 2)
	RegisterCommand("set", execSet, 3)
	RegisterCommand("setnx", execSetNX, 3)
	RegisterCommand("getset", execGetSet, 3)
	RegisterCommand("strlen", execStrLen, 2)
}
