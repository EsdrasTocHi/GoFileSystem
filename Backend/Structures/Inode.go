package structures

type Inode struct {
	I_uid   [8]byte
	I_gid   [8]byte
	I_size  [8]byte
	I_atime [19]byte
	I_ctime [19]byte
	I_mtime [19]byte
	I_block [16]int64
	I_type  byte
	I_perm  [8]byte
}
