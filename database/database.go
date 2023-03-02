package database

import (
	"github.com/A-walker-ninght/miniRedis/aof"
	"github.com/A-walker-ninght/miniRedis/config"
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"github.com/A-walker-ninght/miniRedis/lib/logger"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
	"strconv"
	"strings"
)

type Database struct {
	dbSet      []*DB
	aofHandler *aof.AofHandler
}

func NewDatabase() *Database {
	database := &Database{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := range database.dbSet {
		db := makeDB()
		db.index = i
		database.dbSet[i] = db
	}

	if config.Properties.AppendOnly {
		aofHandler, err := aof.NewAofHandler(database)
		if err != nil {
			panic(err)
		}
		database.aofHandler = aofHandler
		for _, db := range database.dbSet {
			sdb := db
			sdb.addAof = func(line CmdLine) {
				database.aofHandler.AddAof(sdb.index, line)
			}
		}
	}
	return database
}

func (database *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	cmd := strings.ToLower(string(args[0]))
	if cmd == "select" {
		if len(args) != 2 {
			return reply.MakeArgNumErrReply(cmd)
		}
		return execSelect(client, database, args[1:])
	}
	dbIndex := client.GetDBIndex()
	if dbIndex >= len(database.dbSet) {
		return reply.MakeStandardErrReply("ERR DB index is out of range")
	}
	db := database.dbSet[dbIndex]
	return db.Exec(client, args)
}

func (database *Database) Close() error {
	return nil
}

func (database *Database) AfterClientClose(c resp.Connection) {

}

// select 1
func execSelect(c resp.Connection, database *Database, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeStandardErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(database.dbSet) {
		return reply.MakeStandardErrReply("ERR DB index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.MakeOKReply()
}
