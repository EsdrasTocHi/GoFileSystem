package structures

var SizeOfEbr int64 = 42

type Ebr struct {
	Part_status byte
	Part_fit    byte
	Part_start  [8]byte
	Part_size   [8]byte
	Part_next   [8]byte
	Part_name   [16]byte
}
