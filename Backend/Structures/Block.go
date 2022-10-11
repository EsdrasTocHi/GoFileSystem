package structures

type Content struct {
	B_name  [12]byte
	B_inodo [8]byte
}

type Dirblock struct {
	B_content [4]Content
}

type FileBlock struct {
	B_content [64]byte
}
