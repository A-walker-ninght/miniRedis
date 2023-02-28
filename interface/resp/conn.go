package resp

type connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}
