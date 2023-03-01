package parser

import (
	"bufio"
	"errors"
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"github.com/A-walker-ninght/miniRedis/lib/logger"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
)

// 客户端解析后的数据
type Payload struct {
	Data resp.Reply
	Err  error
}

// 解析器的状态
type readState struct {
	readingMultiLine  bool // 多行还是单行？
	expectedArgsCount int  // 字符串的个数
	msgType           byte
	args              [][]byte // [set key value]
	bulkLen           int64
}

func (r *readState) finished() bool {
	return r.expectedArgsCount > 0 && len(r.args) == r.expectedArgsCount
}
func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(debug.Stack())
		}
	}()

	bufReader := bufio.NewReader(reader)
	var state readState
	var err error
	var msg []byte
	for {
		var ioEOF bool
		// 每行读取
		msg, ioEOF, err = readLine(bufReader, &state)
		if err != nil {
			if ioEOF {
				ch <- &Payload{
					Err: err,
				}
				close(ch)
				return
			}
			ch <- &Payload{
				Err: err,
			}
			state = readState{}
			continue
		}

		// 判断是否解析，或者是否单行模式
		// *3\r\n 或者 $7\r\n
		if !state.readingMultiLine {
			// 多行
			if msg[0] == '*' {
				err = parseMutiBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: err,
					}
					state = readState{}
					continue
				}

				// *0
				if state.expectedArgsCount == 0 {
					ch <- &Payload{
						Data: reply.MakeEmptyMultiBulkReply(),
					}
					state = readState{}
					continue
				}
			} else if msg[0] == '$' {
				err = parseBulkHeader(msg, &state)
				if err != nil {
					ch <- &Payload{
						Err: err,
					}
					state = readState{}
					continue
				}

				if state.bulkLen == -1 {
					ch <- &Payload{
						Data: reply.MakeNullBulkReply(),
					}
					state = readState{}
					continue
				}
			} else {
				// singleLine
				result, err := parseSingleLine(msg, &state)
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{}
				continue
			}
		} else {
			// 多行模式
			err = readBody(msg, &state)
			if err != nil {
				ch <- &Payload{
					Err: err,
				}
				state = readState{}
				continue
			}

			// 判断是否读完
			if state.finished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.MakeMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.MakeBulkReply(state.args[0])
				}
				ch <- &Payload{
					Data: result,
					Err:  err,
				}
				state = readState{}
			}
		}
	}
}

// 每行读取，
// 内容，是否是io.EOF，error
func readLine(reader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var msg []byte
	var err error

	if state.bulkLen == 0 {
		msg, err = reader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}

		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
	} else {
		// 之前读到了长度，接着读即可
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(reader, msg)
		if err != nil {
			return nil, true, err
		}

		if len(msg) == 0 || msg[len(msg)-1] != '\n' || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error: " + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

// 解析器
// 单行模式

// 固定消息：+OK\r\n, -ERR\r\n  :12312\r\n
func parseSingleLine(msg []byte, state *readState) (resp.Reply, error) {
	str := strings.TrimSuffix(string(msg), "\r\n") // +Ok, -err
	var result resp.Reply
	switch msg[0] {
	case '+':
		result = reply.MakeStatusReply(str[1:])
	case '-':
		result = reply.MakeStandardErrReply(str[1:])
	case ':':
		n, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, errors.New("protocol error: " + string(msg))
		}
		result = reply.MakeIntReply(n)
	}
	return result, nil
}

// 非固定消息：$7\r\n | dfdssdf \r\n | $-1\r\n |
func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}

	// -1, 空消息
	if state.bulkLen == -1 {
		return nil
	} else if state.bulkLen >= 0 {
		state.args = make([][]byte, 0, 1)
		state.msgType = msg[0]
		state.expectedArgsCount = 1
		state.readingMultiLine = true
		return nil
	}
	return nil
}

// 多行解析：*24324\r\n
func parseMutiBulkHeader(msg []byte, state *readState) error {
	var err error
	var expectedLine int64
	expectedLine, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}

	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		state.expectedArgsCount = int(expectedLine)
		state.args = make([][]byte, 0, expectedLine)
		state.msgType = msg[0]
		state.readingMultiLine = true
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

// 读取数据
// 可能是多行 $xxx \r\n
// 可能是单行 dfasfda\r\n
func readBody(msg []byte, state *readState) error {
	line := msg[0 : len(msg)-2]
	var err error

	if msg[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
		if err != nil {
			return errors.New("protocol error: " + string(msg))
		}
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}
