package reply

import (
	"bytes"
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"strconv"
)

/*自定义回复*/
var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

/* 单行字符串 */
type BulkReply struct {
	Arg []byte // 字符串 $len\r\n string \r\n
}

func (b *BulkReply) ToBytes() []byte {
	if len(b.Arg) == 0 {
		return nullBulkBytes
	}
	return []byte("$" + strconv.Itoa(len(b.Arg)) + CRLF + string(b.Arg) + CRLF)
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{Arg: arg}
}

/* 多行字符串 */
type MultiBulkReply struct {
	Args [][]byte
}

func (m *MultiBulkReply) ToBytes() []byte {
	length := len(m.Args)
	var buf bytes.Buffer

	// [SET, key, value]：*3\r\n $3\r\nSET\r\n $3\r\nkey\r\n $5\r\nvalue\r\n
	buf.WriteString("*" + strconv.Itoa(length) + CRLF)
	for _, arg := range m.Args {
		if arg == nil {
			buf.WriteString(string(nullBulkReplyBytes) + CRLF)
		} else {
			buf.WriteString("$" + strconv.Itoa(len(arg)) + CRLF + string(arg) + CRLF)
		}
	}
	return buf.Bytes()
}

func MakeMultiBulkReply(args [][]byte) *MultiBulkReply {
	return &MultiBulkReply{Args: args}
}

/* 状态回复 */
type StatusReply struct {
	Status string
}

func (s *StatusReply) ToBytes() []byte {
	return []byte("+" + s.Status + CRLF)
}
func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{Status: status}
}

/* 数字回复 */
type IntReply struct {
	Code int64
}

func (i *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.FormatInt(i.Code, 10) + CRLF)
}

func MakeIntReply(code int64) *IntReply {
	return &IntReply{Code: code}
}

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

/* 标准错误回复 */
type StandardErrReply struct {
	Status string
}

func (s *StandardErrReply) Error() string {
	return s.Status
}

func (s *StandardErrReply) ToBytes() []byte {
	return []byte("-" + s.Status + CRLF)
}

func MakeStandardErrReply(status string) *StandardErrReply {
	return &StandardErrReply{Status: status}
}

// 判断是正常回复还是异常？
// 判断第一个字符是否是 "-"
func IsErrReply(reply resp.Reply) bool {
	return reply.ToBytes()[0] == '-'
}
