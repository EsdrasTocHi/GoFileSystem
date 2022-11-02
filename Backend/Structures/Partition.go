package structures

type Partition struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  [8]byte
	Part_size   [8]byte
	Part_name   [16]byte
}
