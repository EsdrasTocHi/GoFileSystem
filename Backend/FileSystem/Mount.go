package filesystem

import (
	structs "Backend/Structures"
	"bytes"
	"encoding/binary"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func Mount(path string, nameString string, partitions *[]structs.MountedPartition, number *int, w http.ResponseWriter) {
	if Exist(path) {
		file, _ := os.OpenFile(path, os.O_RDWR, 0777)
		defer file.Close()
		file.Seek(0, os.SEEK_SET)
		var mbr structs.Mbr
		ReadMbr(&mbr, file)

		var par *structs.Partition
		found := false
		partitionNumber := 'z'

		if strings.ToLower(ToString(mbr.Mbr_partition_1.Part_name[:])) == strings.ToLower(nameString) {
			par = &mbr.Mbr_partition_1
			found = true
			partitionNumber = 'a'
		} else if strings.ToLower(ToString(mbr.Mbr_partition_2.Part_name[:])) == strings.ToLower(nameString) {
			par = &mbr.Mbr_partition_2
			found = true
			partitionNumber = 'b'
		} else if strings.ToLower(ToString(mbr.Mbr_partition_3.Part_name[:])) == strings.ToLower(nameString) {
			par = &mbr.Mbr_partition_3
			found = true
			partitionNumber = 'c'
		} else if strings.ToLower(ToString(mbr.Mbr_partition_4.Part_name[:])) == strings.ToLower(nameString) {
			par = &mbr.Mbr_partition_4
			found = true
			partitionNumber = 'd'
		}

		numberOfPartitions := 0
		if ToInt(mbr.Mbr_partition_1.Part_start[:]) != 0 {
			numberOfPartitions++
		}
		if ToInt(mbr.Mbr_partition_2.Part_start[:]) != 0 {
			numberOfPartitions++
		}
		if ToInt(mbr.Mbr_partition_3.Part_start[:]) != 0 {
			numberOfPartitions++
		}
		if ToInt(mbr.Mbr_partition_4.Part_start[:]) != 0 {
			numberOfPartitions++
		}

		if found {
			if par.Part_type == 'e' {
				WriteResponse(w, "$Error: extender partitions cannot be mounted")
				return
			}

			var newMount structs.MountedPartition
			newMount.Par = *par
			newMount.IsLogic = false
			newMount.Path = path
			(*number)++
			newMount.Id = "73" + strconv.Itoa(*number) + string(partitionNumber)

			*partitions = append(*partitions, newMount)

			WriteResponse(w, "PARTITION MOUNTED, ID -> "+newMount.Id)

			file.Seek(ToInt(par.Part_start[:]), os.SEEK_SET)

			sp := structs.SuperBlock{}
			ReadSuperBlock(&sp, file)

			if ToInt(sp.S_filesystem_type[:]) != 0 {
				binary.BigEndian.PutUint64(sp.S_mnt_count[:], uint64(ToInt(sp.S_mnt_count[:])+1))
				copy(sp.S_mtime[:], []byte(getDate()))

				file.Seek(ToInt(par.Part_start[:]), os.SEEK_SET)
				buffer := bytes.Buffer{}
				binary.Write(&buffer, binary.BigEndian, &sp)
				writeBinary(file, buffer.Bytes())
			}

			return
		}

		if mbr.Mbr_partition_1.Part_type == 'e' {
			par = &mbr.Mbr_partition_1
		} else if mbr.Mbr_partition_2.Part_type == 'e' {
			par = &mbr.Mbr_partition_2
		} else if mbr.Mbr_partition_3.Part_type == 'e' {
			par = &mbr.Mbr_partition_3
		} else if mbr.Mbr_partition_4.Part_type == 'e' {
			par = &mbr.Mbr_partition_4
		} else {
			WriteResponse(w, "$Error: the partition doesn't exist")
			return
		}

		var ebr structs.Ebr
		file.Seek(ToInt(par.Part_start[:]), os.SEEK_SET)
		ReadEbr(&ebr, file)

		if ToInt(ebr.Part_start[:]) == 0 && ebr.Part_next == 0 {
			WriteResponse(w, "$Error: the partition doesn't exist")
			return
		}

		pointer := ToInt(par.Part_start[:])
		logicPartitions := make([]structs.Ebr, 0)
		found = false
		index := 0

		for pointer < ToInt(par.Part_start[:])+ToInt(par.Part_size[:])-1 {
			var ebr structs.Ebr
			file.Seek(pointer, os.SEEK_SET)
			ReadEbr(&ebr, file)

			if strings.ToLower(ToString(ebr.Part_name[:])) == strings.ToLower(nameString) {
				found = true
				index = len(logicPartitions)
			}

			numberOfPartitions++
			if ebr.Part_next == -1 {
				logicPartitions = append(logicPartitions, ebr)
				break
			}

			logicPartitions = append(logicPartitions, ebr)
			pointer = int64(ebr.Part_next)
		}

		if !found {
			WriteResponse(w, "$Error: the partition doesn't exist")
			return
		}

		var newMount structs.MountedPartition
		newMount.LogicPar = logicPartitions[index]
		newMount.IsLogic = true
		newMount.Path = path
		(*number)++
		newMount.Id = "73" + strconv.Itoa(*number) + string(byte(96+*number))

		file.Seek(ToInt(logicPartitions[index].Part_start[:]), os.SEEK_SET)

		sp := structs.SuperBlock{}
		ReadSuperBlock(&sp, file)

		if ToInt(sp.S_filesystem_type[:]) != 0 {
			binary.BigEndian.PutUint64(sp.S_mnt_count[:], uint64(ToInt(sp.S_mnt_count[:])+1))
			copy(sp.S_mtime[:], []byte(getDate()))

			file.Seek(ToInt(logicPartitions[index].Part_start[:]), os.SEEK_SET)
			buffer := bytes.Buffer{}
			binary.Write(&buffer, binary.BigEndian, &sp)
			writeBinary(file, buffer.Bytes())
		}

		*partitions = append(*partitions, newMount)
		WriteResponse(w, "PATITION MOUNTED, ID -> "+newMount.Id)
	} else {
		WriteResponse(w, "$Error: "+path+" doesn't exist")
	}
}
