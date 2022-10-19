package filesystem

import (
	structs "Backend/Structures"
	"bytes"
	"encoding/binary"
	"math"
	"net/http"
	"os"
)

func NewInode() structs.Inode {
	var inode structs.Inode
	binary.BigEndian.PutUint64(inode.I_uid[:], 1)
	binary.BigEndian.PutUint64(inode.I_gid[:], 1)
	binary.BigEndian.PutUint64(inode.I_size[:], 0)
	copy(inode.I_atime[:], []byte(getDate()))
	copy(inode.I_ctime[:], []byte(getDate()))
	copy(inode.I_mtime[:], []byte(getDate()))
	for i := 0; i < 16; i++ {
		inode.I_block[i] = -1
	}
	inode.I_type = byte('0')
	binary.BigEndian.PutUint64(inode.I_perm[:], 0)

	return inode
}

func ext2(mountedPartition *structs.MountedPartition, w http.ResponseWriter) {
	var sizeOfPartition int64
	var start int64

	if mountedPartition.IsLogic {
		sizeOfPartition = ToInt(mountedPartition.LogicPar.Part_size[:])
		start = ToInt(mountedPartition.LogicPar.Part_start[:])
	} else {
		sizeOfPartition = ToInt(mountedPartition.Par.Part_size[:])
		start = ToInt(mountedPartition.Par.Part_start[:])
	}

	var sp structs.SuperBlock
	num_structures := math.Floor(float64((sizeOfPartition - structs.SizeOfSuperBlock) / (4 + structs.SizeOfInode + 3*64)))
	num_blocks := 3 * num_structures

	binary.BigEndian.PutUint64(sp.S_filesystem_type[:], uint64(2))
	binary.BigEndian.PutUint64(sp.S_inodes_count[:], uint64(num_structures))
	binary.BigEndian.PutUint64(sp.S_blocks_count[:], uint64(num_blocks))
	binary.BigEndian.PutUint64(sp.S_free_blocks_count[:], uint64(num_blocks-2))
	binary.BigEndian.PutUint64(sp.S_free_inodes_count[:], uint64(num_structures-2))
	binary.BigEndian.PutUint64(sp.S_mnt_count[:], 0)
	binary.BigEndian.PutUint64(sp.S_magic[:], 0xEF53)
	binary.BigEndian.PutUint64(sp.S_inode_size[:], uint64(binary.Size(structs.Inode{})))
	binary.BigEndian.PutUint64(sp.S_block_size[:], uint64(binary.Size(structs.FileBlock{})))
	binary.BigEndian.PutUint64(sp.S_first_ino[:], 2)
	binary.BigEndian.PutUint64(sp.S_first_blo[:], 2)
	binary.BigEndian.PutUint64(sp.S_bm_inode_start[:], uint64(start+int64(binary.Size(sp))))
	binary.BigEndian.PutUint64(sp.S_bm_block_start[:], uint64(ToInt(sp.S_bm_inode_start[:])+int64(num_structures)))
	binary.BigEndian.PutUint64(sp.S_inode_start[:], uint64(ToInt(sp.S_bm_block_start[:])+int64(num_blocks)))
	binary.BigEndian.PutUint64(sp.S_block_start[:], uint64(ToInt(sp.S_inode_start[:])+int64(num_structures)*int64(binary.Size(structs.Inode{}))))

	file, _ := os.OpenFile(mountedPartition.Path, os.O_RDWR, 0777)
	defer file.Close()
	file.Seek(start, os.SEEK_SET)
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, &sp)
	writeBinary(file, buffer.Bytes())
	buffer = bytes.Buffer{}

	binary.Write(&buffer, binary.BigEndian, byte('1'))
	writeBinary(file, buffer.Bytes())
	writeBinary(file, buffer.Bytes())
	buffer = bytes.Buffer{}

	binary.Write(&buffer, binary.BigEndian, byte('0'))
	for i := 2; i < int(num_structures); i++ {
		writeBinary(file, buffer.Bytes())
	}
	buffer = bytes.Buffer{}

	binary.Write(&buffer, binary.BigEndian, byte('1'))
	writeBinary(file, buffer.Bytes())
	writeBinary(file, buffer.Bytes())
	buffer = bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, byte('0'))
	for i := 2; i < int(num_blocks); i++ {
		writeBinary(file, buffer.Bytes())
	}
	buffer = bytes.Buffer{}

	inode := NewInode()
	inode.I_block[0] = 0
	binary.BigEndian.PutUint64(inode.I_perm[:], uint64(664))

	binary.Write(&buffer, binary.BigEndian, &inode)
	writeBinary(file, buffer.Bytes())
	buffer = bytes.Buffer{}

	inode = NewInode()
	binary.BigEndian.PutUint64(inode.I_size[:], uint64(27))
	inode.I_block[0] = 1
	inode.I_type = byte('1')
	binary.BigEndian.PutUint64(inode.I_perm[:], uint64(755))
	binary.Write(&buffer, binary.BigEndian, &inode)
	writeBinary(file, buffer.Bytes())
	buffer = bytes.Buffer{}

	var dirBlock structs.Dirblock
	copy(dirBlock.B_content[0].B_name[:], ".")
	dirBlock.B_content[0].B_inodo = 0
	copy(dirBlock.B_content[1].B_name[:], "..")
	dirBlock.B_content[1].B_inodo = 0
	copy(dirBlock.B_content[2].B_name[:], "users.txt")
	dirBlock.B_content[2].B_inodo = 1

	copy(dirBlock.B_content[3].B_name[:], "")
	dirBlock.B_content[3].B_inodo = -1

	file.Seek(ToInt(sp.S_block_start[:]), os.SEEK_SET)
	binary.Write(&buffer, binary.BigEndian, &dirBlock)
	writeBinary(file, buffer.Bytes())
	buffer = bytes.Buffer{}

	var fileBlock structs.FileBlock
	file.Seek(ToInt(sp.S_block_start[:])+64, os.SEEK_SET)
	copy(fileBlock.B_content[:], "1,G,root\n1,U,root,root,123\n")

	binary.Write(&buffer, binary.BigEndian, &fileBlock)
	writeBinary(file, buffer.Bytes())

	WriteResponse(w, "EXT2 FORMAT DONE SUCCESFULLY")
}

func Mkfs(id string, partitions *[]structs.MountedPartition, w http.ResponseWriter) {
	var mountedPartition *structs.MountedPartition
	i := 0
	for i = 0; i < len(*partitions); i++ {
		if id == (*partitions)[i].Id {
			mountedPartition = &((*partitions)[i])
			break
		}
	}

	if i == len(*partitions) {
		WriteResponse(w, "$Error: "+id+" is not mounted")
		return
	}

	ext2(mountedPartition, w)
}
