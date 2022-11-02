package structures

type Mbr struct {
	Mbr_tamano         [8]byte
	Mbr_fecha_creacion [19]byte
	Mbr_dsk_signature  [8]byte
	Mbr_dsk_fit        [8]byte
	Mbr_partition_1    Partition
	Mbr_partition_2    Partition
	Mbr_partition_3    Partition
	Mbr_partition_4    Partition
}
