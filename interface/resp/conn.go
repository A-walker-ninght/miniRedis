package resp

<<<<<<< HEAD
type connection interface {
=======
type Connection interface {
>>>>>>> 70f3717 (resp 2023.3.1)
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}
