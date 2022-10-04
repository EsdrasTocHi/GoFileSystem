package structures

type Partition struct {
	Part_status [1]byte
	Part_type   [1]byte
	Part_fit    [1]byte
	Part_start  [8]byte
	Part_size   [8]byte
	Part_name   [16]byte
}
