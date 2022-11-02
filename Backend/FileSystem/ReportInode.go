package filesystem

import (
	structs "Backend/Structures"
	"os"
	"strconv"
)

func ReportInodeTree(inode structs.Inode, link string, pointer int64, nodes *string, edges *string, istart int64, bstart int64, file *os.File) {
	*nodes += "i" + strconv.Itoa(int(pointer)) + "[label=\"Inode " + strconv.Itoa(int(pointer)) + "| {i_uid|" + strconv.Itoa(int(ToInt(inode.I_uid[:]))) + "}|{i_gid|" + strconv.Itoa(int(ToInt(inode.I_gid[:]))) + "}|{i_s|" + strconv.Itoa(int(ToInt(inode.I_size[:]))) + "}|\n"
	*nodes += "    {i_atime|" + ToString(inode.I_atime[:]) + "}|{i_ctime|" + ToString(inode.I_ctime[:]) + "}|{i_mtime|" + ToString(inode.I_mtime[:]) + "}|\n"
	for i := 1; i < 17; i++ {
		*nodes += "    {i_block " + strconv.Itoa(i) + "|<i" + strconv.Itoa(i) + "> " + strconv.Itoa(int(inode.I_block[i-1])) + "}|\n"
	}
	*nodes += "{i_type|"
	*nodes += string(inode.I_type)
	*nodes += "} | {i_perm|" + strconv.Itoa(int(ToInt(inode.I_perm[:]))) + "}\"];"

	if link != "" {
		*edges += link + "->i" + strconv.Itoa(int(pointer)) + ";\n"
	}
	//directorio ->dirblocks o pointerblocks
	if inode.I_type == '0' {
		for i := 0; i < 16; i++ {
			if inode.I_block[i] != -1 {
				db := structs.Dirblock{}
				file.Seek(bstart+(64*inode.I_block[i]), os.SEEK_SET)
				ReadDirBlock(&db, file)
				l := "i" + strconv.Itoa(int(pointer)) + ":i" + strconv.Itoa(i+1)
				ReportDirBlock(db, l, inode.I_block[i], nodes, edges, istart, bstart, file)
			}
		}
	} else {
		for i := 0; i < 16; i++ {
			if inode.I_block[i] != -1 {
				db := structs.FileBlock{}
				file.Seek(bstart+(64*inode.I_block[i]), os.SEEK_SET)
				ReadFileBlock(&db, file)
				l := "i" + strconv.Itoa(int(pointer)) + ":i" + strconv.Itoa(i+1)
				ReportFileBlock(db, l, inode.I_block[i], nodes, edges)
			}
		}
	}
}
