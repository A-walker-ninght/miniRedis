package aof

import (
	"github.com/A-walker-ninght/miniRedis/config"
	dbinterface "github.com/A-walker-ninght/miniRedis/interface/database"
	"github.com/A-walker-ninght/miniRedis/lib/logger"
	"github.com/A-walker-ninght/miniRedis/lib/utils"
	"github.com/A-walker-ninght/miniRedis/resp/connetion"
	"github.com/A-walker-ninght/miniRedis/resp/parser"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
	"io"
	"os"
	"strconv"
)

type CmdLine = [][]byte
type payload struct {
	cmdLine CmdLine
	dbIndex int
}

type AofHandler struct {
	database    dbinterface.DataBase
	aofFile     *os.File
	aofFileName string
	aofChan     chan *payload
	currentDB   int // 当前插入的db
}

// 创建
func NewAofHandler(database dbinterface.DataBase) (*AofHandler, error) {
	handler := &AofHandler{}
	handler.aofFileName = config.Properties.AppendFilename
	handler.database = database

	handler.AofRecover()
	fd, err := os.OpenFile(handler.aofFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofChan = make(chan *payload, 1<<16)
	go func() {
		handler.handlerAof()
	}()

	handler.aofFile = fd
	return handler, nil
}

// 添加指令
func (handler *AofHandler) AddAof(dbIndex int, cmd CmdLine) {
	if config.Properties.AppendOnly && handler.aofChan != nil {
		handler.aofChan <- &payload{
			dbIndex: dbIndex,
			cmdLine: cmd,
		}
	}
}

// 落盘
func (handler *AofHandler) handlerAof() {
	handler.currentDB = 0

	for p := range handler.aofChan {
		// 记录db是否切换
		if p.dbIndex != handler.currentDB {
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Error(err)
				continue
			}
			handler.currentDB = p.dbIndex
		}
		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(data)
		if err != nil {
			logger.Error(err)
		}
	}
}

// 恢复
func (handler *AofHandler) AofRecover() {
	open, err := os.Open(handler.aofFileName)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		_ = open.Close()
	}()

	ch := parser.ParseStream(open)
	dummyConn := &connetion.Connection{}
	for p := range ch {
		if p.Err != nil {
			if p.Err == io.EOF {
				break
			}
			logger.Error(p.Err.Error())
			continue
		}

		if p.Data == nil {
			logger.Error("empty payload")
			continue
		}

		repl, ok := p.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("need multibulk")
			continue
		}
		rep := handler.database.Exec(dummyConn, repl.Args)
		if reply.IsErrReply(rep) {
			logger.Error("exec err", rep.ToBytes())
		}
	}
}
