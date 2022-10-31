package filesystem

import (
	structs "Backend/Structures"
	"encoding/binary"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetExtension(path string) string {
	res := ""
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			break
		}
		res = string(path[i]) + res
	}

	return res
}

func RemoveFileName(path string) string {
	p := ""
	add := false
	for i := len(path) - 1; i >= 0; i-- {
		if !add {
			if path[i] == '/' {
				add = true
			}
			continue
		}

		p = string(path[i]) + p
	}

	return p
}

func SaveImageGV(file_path string, content string, w http.ResponseWriter) {
	paux := RemoveFileName(file_path)

	if paux != "" {
		exec.Command("mkdir", "-p", paux).Run()
		exec.Command("chmod", "-R", "777", paux).Run()
	}

	dot, _ := os.OpenFile("temp.dot", os.O_CREATE, 0777)
	dot.Close()
	dot, _ = os.OpenFile("temp.dot", os.O_RDWR, 0777)
	dot.WriteString(content)
	exec.Command("dot", "-T", GetExtension(file_path), "temp.dot", "-o", file_path).Run()
	WriteResponse(w, "REPORT COMPLETED")
}

func ReportInode(inode structs.Inode, nodes *string, edges *string, lastInode int64, actualInode int64) {
	*nodes += "<tr>\n"
	*nodes += "                            <td>I_uid</td>\n"
	*nodes += "                            <td>" + strconv.Itoa(int(ToInt(inode.I_uid[:]))) + "</td>\n"
	*nodes += "                        </tr>\n"
	*nodes += "                        <tr>\n"
	*nodes += "                            <td>I_gid</td>\n"
	*nodes += "                            <td>" + strconv.Itoa(int(ToInt(inode.I_gid[:]))) + "</td>\n"
	*nodes += "                        </tr>\n"
	*nodes += "                        <tr>\n"
	*nodes += "                            <td>I_s</td>\n"
	*nodes += "                            <td>" + strconv.Itoa(int(ToInt(inode.I_size[:]))) + "</td>\n"
	*nodes += "                        </tr>\n"
	*nodes += "                        <tr>\n"
	*nodes += "                            <td>I_atime</td>\n"
	*nodes += "                            <td>" + string(inode.I_atime[:]) + "</td>\n"
	*nodes += "                        </tr>\n"
	*nodes += "                        <tr>\n"
	*nodes += "                            <td>I_ctime</td>\n"
	*nodes += "                            <td>" + string(inode.I_ctime[:]) + "</td>\n"
	*nodes += "                        </tr>\n"
	*nodes += "                        <tr>\n"
	*nodes += "                            <td>I_mtime</td>\n"
	*nodes += "                            <td>" + string(inode.I_mtime[:]) + "</td>\n"
	*nodes += "                        </tr>\n"
	for i := 0; i < 16; i++ {
		*nodes += "                        <tr>\n"
		*nodes += "                            <td>I_block" + strconv.Itoa(i+1) + "</td>\n"
		*nodes += "                            <td>" + strconv.Itoa(int(inode.I_block[i])) + "</td>\n"
		*nodes += "                        </tr>\n"
	}
	*nodes += "                        <tr>\n"
	*nodes += "                            <td>I_type</td>\n"
	*nodes += "                            <td>"
	*nodes += string(inode.I_type)
	*nodes += "</td>\n"
	*nodes += "                        </tr>\n"
	*nodes += "                        <tr>\n"
	*nodes += "                            <td>I_perm</td>\n"
	*nodes += "                            <td>" + strconv.Itoa(int(ToInt(inode.I_perm[:]))) + "</td>\n"
	*nodes += "                        </tr>\n"

	if lastInode != -1 {
		*edges += "i" + strconv.Itoa(int(lastInode)) + "->i" + strconv.Itoa(int(actualInode)) + ";\n"
	}
}

func ReportFileBlock(fb structs.FileBlock, link string, pointer int64, nodes *string, edges *string) {
	c := ""
	if len(fb.B_content) > 64 {
		for i := 0; i < 64; i++ {
			c += string(fb.B_content[i])
		}
	} else {
		c = ToString(fb.B_content[:])
	}
	*nodes += "b" + strconv.Itoa(int(pointer)) + "[label=<\n"
	*nodes += "        <table>\n"
	*nodes += "            <tr>\n"
	*nodes += "                <td>\n"
	*nodes += "                    <table>\n"
	*nodes += "                        <tr>\n"
	*nodes += "                            <td>FileBlock</td>\n"
	*nodes += "                            <td>" + strconv.Itoa(int(pointer)) + "</td>\n"
	*nodes += "                        </tr>\n"
	*nodes += "                    </table>\n"
	*nodes += "                </td>\n"
	*nodes += "            </tr>\n"
	*nodes += "            <tr>\n"
	*nodes += "                <td>" + c + "</td>\n"
	*nodes += "            </tr>\n"
	*nodes += "        </table>\n"
	*nodes += "    >];\n"

	*edges += link + "->" + "b" + strconv.Itoa(int(pointer)) + ";\n"
}

func ReportDirBlock(db structs.Dirblock, link string, pointer int64, nodes *string, edges *string, istart int64, bstart int64, file *os.File) {
	*nodes += "b" + strconv.Itoa(int(pointer)) + "[label=\"DirBlock " + strconv.Itoa(int(pointer))
	for i := 1; i < 5; i++ {
		*nodes += "|{" + ToString(db.B_content[i-1].B_name[:]) + "|<b" + strconv.Itoa(i) + ">" + strconv.Itoa(int(db.B_content[i-1].B_inodo)) + "}"
	}
	*nodes += "\"];\n"

	*edges += link + "->b" + strconv.Itoa(int(pointer)) + ";\n"
	for i := 0; i < 4; i++ {
		if ToString(db.B_content[i].B_name[:]) != "." && ToString(db.B_content[i].B_name[:]) != ".." {
			if db.B_content[i].B_inodo == -2 {
				continue
			}
			if db.B_content[i].B_inodo != -1 {
				aux := structs.Inode{}
				file.Seek(istart+(int64(binary.Size(aux))*int64(db.B_content[i].B_inodo)), os.SEEK_SET)
				ReadInode(&aux, file)

				link := "b" + strconv.Itoa(int(pointer)) + ":<b" + strconv.Itoa(i+1) + ">"
				ReportInodeTree(aux, link, int64(db.B_content[i].B_inodo), nodes, edges, istart, bstart, file)
			}
		}
	}
}

func ReportTree(partition structs.MountedPartition, path string, w http.ResponseWriter) {
	file, _ := os.OpenFile(partition.Path, os.O_RDWR, 0777)
	defer file.Close()
	start := int64(0)
	if partition.IsLogic {
		start = ToInt(partition.LogicPar.Part_start[:])
	} else {
		start = ToInt(partition.Par.Part_start[:])
	}

	file.Seek(start, os.SEEK_SET)
	sb := structs.SuperBlock{}
	ReadSuperBlock(&sb, file)
	if ToInt(sb.S_filesystem_type[:]) == 0 {
		WriteResponse(w, "$Error: the partition is not formatted")
		return
	}
	root := structs.Inode{}
	file.Seek(ToInt(sb.S_inode_start[:]), os.SEEK_SET)
	ReadInode(&root, file)
	nodes := ""
	edges := ""

	ReportInodeTree(root, "", 0, &nodes, &edges, ToInt(sb.S_inode_start[:]), ToInt(sb.S_block_start[:]), file)

	content := "digraph G {node[shape = record];rankdir = LR;\n" + nodes + edges + "}"

	SaveImageGV(path, content, w)
}

func Report(id string, name string, path string, partitions *[]structs.MountedPartition, ruta string, currentUser structs.Sesion, w http.ResponseWriter) {
	var mountedPartition *structs.MountedPartition
	i := 0
	for i = 0; i < len(*partitions); i++ {
		if id == (*partitions)[i].Id {
			mountedPartition = &(*partitions)[i]
			break
		}
	}

	if i == len(*partitions) {
		WriteResponse(w, "$Error: "+id+" is not mounted")
		return
	}

	if strings.ToLower(name) == "tree" {
		ReportTree(*mountedPartition, path, w)
	} else if strings.ToLower(name) == "file" {
		File(currentUser, ruta, path, w)
	} else if strings.ToLower(name) == "sb" {
		ReportSb(*mountedPartition, path, w)
	}
}

func File(currentUser structs.Sesion, reportPath string, filePath string, w http.ResponseWriter) {
	var mountedPartition *structs.MountedPartition
	mountedPartition = &(currentUser.Mounted)
	file, _ := os.OpenFile(mountedPartition.Path, os.O_RDWR, 0777)
	defer file.Close()

	var sp structs.SuperBlock
	start := int64(0)
	if mountedPartition.IsLogic {
		start = ToInt(mountedPartition.LogicPar.Part_start[:])
	} else {
		start = ToInt(mountedPartition.Par.Part_start[:])
	}

	file.Seek(start, os.SEEK_SET)
	ReadSuperBlock(&sp, file)

	root := structs.Inode{}
	aux := structs.Inode{}
	file.Seek(ToInt(sp.S_inode_start[:]), os.SEEK_SET)
	ReadInode(&root, file)
	pointerOfFile := int64(0)
	aux = SearchFile(file, root, SplithPath("users.txt"), ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), &pointerOfFile)
	c := ReadFile(file, aux, ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), w)

	finalContent := ""
	aux = SearchFile(file, root, SplithPath(filePath), ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), &pointerOfFile)

	if aux.I_type == 'n' {
		WriteResponse(w, "$Error: file does not exist")
		return
	}

	if GetPermission(aux, currentUser.Usr.Id, int64(GetGroupId(ToString(currentUser.Usr.Group[:]), c)), ToInt(aux.I_perm[:]), true, false, false) {
		finalContent += ReadFile(file, aux, ToInt(sp.S_inode_start[:]), ToInt(sp.S_block_start[:]), w) + "\n"
	} else {
		WriteResponse(w, "$Error: You do not have permission to read "+filePath)
		return
	}

	report := "digraph G{\nnode[shape = record];rankdir = LR;\nfile[label=\"" + filePath + "|" + finalContent + "\"];}"
	SaveImageGV(reportPath, report, w)
}

func ReportSb(partition structs.MountedPartition, path string, w http.ResponseWriter) {
	file, _ := os.OpenFile(partition.Path, os.O_RDWR, 0777)
	defer file.Close()
	start := int64(0)
	if partition.IsLogic {
		start = ToInt(partition.LogicPar.Part_start[:])
	} else {
		start = ToInt(partition.Par.Part_start[:])
	}

	sb := structs.SuperBlock{}
	file.Seek(start, os.SEEK_SET)
	ReadSuperBlock(&sb, file)

	report := ""
	report += "digraph G {\n"
	report += "    node[shape = record];rankdir = LR;\n"
	report += "    \n"
	report += "    superBlock[label=\"SuperBlock|{s_filesystem_type|" + strconv.Itoa(int(ToInt(sb.S_filesystem_type[:]))) + "}|{s_inodes_count|" + strconv.Itoa(int(ToInt(sb.S_inodes_count[:]))) + "}|\n"
	report += "    {s_blocks_count|" + strconv.Itoa(int(ToInt(sb.S_blocks_count[:]))) + "}|{s_free_blocks_count|" + strconv.Itoa(int(ToInt(sb.S_free_blocks_count[:]))) + "}|{s_free_inodes_count|" + strconv.Itoa(int(ToInt(sb.S_free_inodes_count[:]))) + "}|\n"
	report += "    {s_mtime|" + ToString(sb.S_mtime[:]) + "}|{s_mnt_count|" + strconv.Itoa(int(ToInt(sb.S_mnt_count[:]))) + "}|{s_magic|" + strconv.Itoa(int(ToInt(sb.S_magic[:]))) + "}|{s_inode_s|" + strconv.Itoa(int(ToInt(sb.S_inode_size[:]))) + "}|\n"
	report += "    {s_block_s|" + strconv.Itoa(int(ToInt(sb.S_block_size[:]))) + "}|{s_first_ino|" + strconv.Itoa(int(ToInt(sb.S_first_ino[:]))) + "}|{s_first_blo|" + strconv.Itoa(int(ToInt(sb.S_first_blo[:]))) + "}|{s_bm_inode_start|" + strconv.Itoa(int(ToInt(sb.S_bm_inode_start[:]))) + "}|\n"
	report += "    {s_bm_block_start|" + strconv.Itoa(int(ToInt(sb.S_bm_block_start[:]))) + "}|{s_inode_start|" + strconv.Itoa(int(ToInt(sb.S_inode_start[:]))) + "}|{s_block_start|" + strconv.Itoa(int(ToInt(sb.S_block_start[:]))) + "}\"];\n"
	report += "}"

	SaveImageGV(path, report, w)
}
