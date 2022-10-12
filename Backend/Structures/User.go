package structures

type User struct {
	Id       int64
	Group    [10]byte
	Name     [10]byte
	Password [10]byte
}
