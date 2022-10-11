package structures

var SizeOfSuperBlock int64 = 139

type SuperBlock struct {
	S_filesystem_type   [8]byte
	S_inodes_count      [8]byte
	S_blocks_count      [8]byte
	S_free_blocks_count [8]byte
	S_free_inodes_count [8]byte
	S_mtime             [19]byte
	S_mnt_count         [8]byte
	S_magic             [8]byte
	S_inode_size        [8]byte
	S_block_size        [8]byte
	S_first_ino         [8]byte
	S_first_blo         [8]byte
	S_bm_inode_start    [8]byte
	S_bm_block_start    [8]byte
	S_inode_start       [8]byte
	S_block_start       [8]byte
}
