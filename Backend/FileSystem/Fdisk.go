package filesystem

import (
	structs "Backend/Structures"
	"bytes"
	"encoding/binary"
	"net/http"
	"os"
	"sort"
	"time"
)

func FDisk(size int, unit rune, path string, tyype rune, fit rune, name string, w http.ResponseWriter) {
	if Exist(path) {
		if len(name) > 16 {
			WriteResponse(w, "$Error: partition name cannot be longer than 16 characters")
			return
		}

		if tyype == 'p' || tyype == 'e' {
			file, _ := os.OpenFile(path, os.O_RDWR, 0777)
			defer file.Close()
			var mbr structs.Mbr
			file.Seek(0, os.SEEK_SET)
			ReadMbr(&mbr, file)

			var newPar *structs.Partition
			if ToInt(mbr.Mbr_partition_1.Part_start[:]) == 0 {
				newPar = &mbr.Mbr_partition_1
			} else if ToInt(mbr.Mbr_partition_2.Part_start[:]) == 0 {
				newPar = &mbr.Mbr_partition_2
			} else if ToInt(mbr.Mbr_partition_3.Part_start[:]) == 0 {
				newPar = &mbr.Mbr_partition_3
			} else if ToInt(mbr.Mbr_partition_4.Part_start[:]) == 0 {
				newPar = &mbr.Mbr_partition_4
			} else {
				WriteResponse(w, "$Error: maximum partitions reached")
				return
			}

			if ToString(mbr.Mbr_partition_1.Part_name[:]) == name ||
				ToString(mbr.Mbr_partition_2.Part_name[:]) == name ||
				ToString(mbr.Mbr_partition_3.Part_name[:]) == name ||
				ToString(mbr.Mbr_partition_4.Part_name[:]) == name {
				WriteResponse(w, "$Error: a partitions with that name already exists")
				return
			}

			if tyype == 'e' {
				if mbr.Mbr_partition_1.Part_type == 'e' || mbr.Mbr_partition_2.Part_type == 'e' ||
					mbr.Mbr_partition_3.Part_type == 'e' || mbr.Mbr_partition_4.Part_type == 'e' {
					WriteResponse(w, "$Error: an extended partition already exist")
					return
				}
			}

			partitions := make([]structs.Partition, 0)[:]
			numOfPartitions := 0
			if ToInt(mbr.Mbr_partition_1.Part_start[:]) > 0 {
				partitions = append(partitions, mbr.Mbr_partition_1)
				numOfPartitions++
			}
			if ToInt(mbr.Mbr_partition_2.Part_start[:]) > 0 {
				partitions = append(partitions, mbr.Mbr_partition_2)
				numOfPartitions++
			}
			if ToInt(mbr.Mbr_partition_3.Part_start[:]) > 0 {
				partitions = append(partitions, mbr.Mbr_partition_3)
				numOfPartitions++
			}
			if ToInt(mbr.Mbr_partition_4.Part_start[:]) > 0 {
				partitions = append(partitions, mbr.Mbr_partition_4)
				numOfPartitions++
			}

			sort.Slice(partitions, func(i, j int) bool {
				return ToInt(partitions[i].Part_start[:]) < ToInt(partitions[j].Part_start[:])
			})

			availableSpace := make([]int64, 0)
			startOfSpace := make([]int64, 0)
			if len(partitions) > 0 {
				for i := 0; i < len(partitions); i++ {
					if i == len(partitions)-1 {
						availableSpace = append(availableSpace,
							ToInt(mbr.Mbr_tamano[:])-(ToInt(partitions[i].Part_start[:])-1+ToInt(partitions[i].Part_size[:])))

						startOfSpace = append(startOfSpace, ToInt(partitions[i].Part_start[:])-1+ToInt(partitions[i].Part_size[:]))
						continue
					}

					availableSpace = append(availableSpace,
						ToInt(partitions[i+1].Part_start[:])-(ToInt(partitions[i].Part_start[:])+ToInt(partitions[i].Part_size[:])))

					startOfSpace = append(startOfSpace, ToInt(partitions[i].Part_start[:])+ToInt(partitions[i].Part_size[:]))
				}
			} else {
				availableSpace = append(availableSpace, ToInt(mbr.Mbr_tamano[:])-173)
				startOfSpace = append(startOfSpace, 173)
			}

			var tam int64 = 0
			if unit == 'm' {
				tam = int64(size * 1024 * 1024)
			} else if unit == 'k' {
				tam = int64(size * 1024)
			} else {
				tam = int64(size)
			}

			index := -1
			if ToRune(mbr.Mbr_dsk_fit[:]) == 'b' {
				for i := 0; i < len(availableSpace)-1; i++ {
					for j := 0; j < len(availableSpace)-i-1; j++ {
						if availableSpace[j] < availableSpace[j+1] {
							temp := availableSpace[j]
							availableSpace[j] = availableSpace[j+1]
							availableSpace[j+1] = temp

							temp = startOfSpace[j]
							startOfSpace[j] = startOfSpace[j+1]
							startOfSpace[j+1] = temp
						}
					}
				}

				for i := 0; i < len(availableSpace); i++ {
					if availableSpace[i] > tam {
						index = 1
					} else if availableSpace[i] == tam {
						index = i
						break
					} else {
						break
					}
				}
			} else if ToRune(mbr.Mbr_dsk_fit[:]) == 'w' {
				for i := 0; i < len(availableSpace)-1; i++ {
					for j := 0; j < len(availableSpace)-i-1; j++ {
						if availableSpace[j] > availableSpace[j+1] {
							temp := availableSpace[j]
							availableSpace[j] = availableSpace[j+1]
							availableSpace[j+1] = temp

							temp = startOfSpace[j]
							startOfSpace[j] = startOfSpace[j+1]
							startOfSpace[j+1] = temp
						}
					}
				}

				for i := 0; i < len(availableSpace); i++ {
					if availableSpace[i] >= tam {
						index = i
					}
				}
			} else {
				for i := 0; i < len(availableSpace); i++ {
					if availableSpace[i] >= tam {
						index = i
						break
					}
				}
			}

			if index == -1 {
				WriteResponse(w, "$Error: no space available")
				return
			}

			newPar.Part_status = 'A'
			newPar.Part_type = byte(tyype)
			newPar.Part_fit = byte(fit)
			binary.BigEndian.PutUint64(newPar.Part_start[:], uint64(startOfSpace[index]))
			binary.BigEndian.PutUint64(newPar.Part_size[:], uint64(tam))
			for i := 0; i < 16; i++ {
				if i < len(name) {
					newPar.Part_name[i] = name[i]
					continue
				}
				newPar.Part_name[i] = 0
			}

			file.Seek(0, os.SEEK_SET)
			var buffer bytes.Buffer
			binary.Write(&buffer, binary.BigEndian, &mbr)
			writeBinary(file, buffer.Bytes())

			WriteResponse(w, "PARTITION CREATED SUCCESFULLY")
			return
		} else {
			file, _ := os.OpenFile(path, os.O_RDWR, 0777)
			defer file.Close()
			var mbr structs.Mbr
			file.Seek(0, os.SEEK_SET)
			ReadMbr(&mbr, file)

			var extendedPartition *structs.Partition

			if mbr.Mbr_partition_1.Part_type == 'e' {
				extendedPartition = &mbr.Mbr_partition_1
			} else if mbr.Mbr_partition_2.Part_type == 'e' {
				extendedPartition = &mbr.Mbr_partition_2
			} else if mbr.Mbr_partition_3.Part_type == 'e' {
				extendedPartition = &mbr.Mbr_partition_3
			} else if mbr.Mbr_partition_4.Part_type == 'e' {
				extendedPartition = &mbr.Mbr_partition_4
			} else {
				WriteResponse(w, "$Error: there is no extended partition")
				return
			}

			file.Seek(ToInt(extendedPartition.Part_start[:]), os.SEEK_SET)
			var ebr structs.Ebr
			ReadEbr(&ebr, file)

			var tam int64 = 0
			if unit == 'm' {
				tam = int64(size * 1024 * 1024)
			} else if unit == 'k' {
				tam = int64(size * 1024)
			} else {
				tam = int64(size)
			}

			spaceNeeded := tam + structs.SizeOfEbr

			if ToInt(ebr.Part_start[:]) == 0 && ToInt(ebr.Part_next[:]) == 0 {
				if ToInt(extendedPartition.Part_size[:]) > spaceNeeded {
					var aux int64 = 1
					binary.BigEndian.PutUint64(ebr.Part_start[:], uint64(ToInt(extendedPartition.Part_start[:])+structs.SizeOfEbr))
					binary.BigEndian.PutUint64(ebr.Part_next[:], uint64(aux))
					ebr.Part_fit = byte(fit)
					binary.BigEndian.PutUint64(ebr.Part_size[:], uint64(tam))
					ebr.Part_status = 'A'
					for i := 0; i < 16; i++ {
						if i < len(name) {
							ebr.Part_name[i] = name[i]
							continue
						}
						ebr.Part_name[i] = 0
					}

					file.Seek(ToInt(extendedPartition.Part_start[:]), os.SEEK_SET)
					var buffer bytes.Buffer
					binary.Write(&buffer, binary.BigEndian, &ebr)
					writeBinary(file, buffer.Bytes())
					WriteResponse(w, "PARTITION CREATED SUCCESFULLY")
					return
				} else {
					WriteResponse(w, "$Error: no space available")
					return
				}
			}

			pointer := ToInt(extendedPartition.Part_start[:])
			partitions := make([]structs.Ebr, 0)
			availableSpace := make([]int64, 0)

			for pointer < ToInt(extendedPartition.Part_start[:])+ToInt(extendedPartition.Part_size[:])-1 {
				var ebr structs.Ebr
				file.Seek(pointer, os.SEEK_SET)
				ReadEbr(&ebr, file)
				time.Sleep(3000 * time.Millisecond)
				if ToInt(ebr.Part_next[:]) == 1 {
					availableSpace = append(availableSpace, (ToInt(extendedPartition.Part_start[:])+ToInt(extendedPartition.Part_size[:]))-(ToInt(ebr.Part_start[:])+ToInt(ebr.Part_size[:])))
					partitions = append(partitions, ebr)
					break
				}

				partitions = append(partitions, ebr)
				availableSpace = append(availableSpace, ToInt(ebr.Part_next[:])-(ToInt(ebr.Part_start[:])+ToInt(ebr.Part_size[:])))

				pointer = ToInt(ebr.Part_next[:])
			}

			if len(partitions) == 23 {
				WriteResponse(w, "$Error: maximum of logic partitions reached")
				return
			}

			index := 0
			if extendedPartition.Part_fit == 'b' {
				actualSpace := -1
				for i := 0; i < len(availableSpace); i++ {
					if availableSpace[i] >= spaceNeeded {
						if actualSpace == -1 {
							actualSpace = int(availableSpace[i])
							index = i
							continue
						} else {
							if availableSpace[i] < int64(actualSpace) {
								actualSpace = int(availableSpace[i])
								index = i
							}
							continue
						}
					}
				}
			} else if extendedPartition.Part_fit == 'w' {
				actualSpace := -1
				for i := 0; i < len(availableSpace); i++ {
					if availableSpace[i] >= spaceNeeded {
						if actualSpace == -1 {
							actualSpace = int(availableSpace[i])
							index = i
							continue
						} else {
							if availableSpace[i] > int64(actualSpace) {
								actualSpace = int(availableSpace[i])
								index = i
							}
							continue
						}
					}
				}
			} else {
				for i := 0; i < len(availableSpace); i++ {
					if availableSpace[i] >= spaceNeeded {
						index = i
						break
					}
				}
			}

			if index == -1 {
				WriteResponse(w, "$Error: no space available")
				return
			}

			binary.BigEndian.PutUint64(ebr.Part_start[:],
				uint64(ToInt(partitions[index].Part_start[:])+ToInt(partitions[index].Part_size[:])+structs.SizeOfEbr))

			binary.BigEndian.PutUint64(ebr.Part_next[:], uint64(1))
			if index < len(availableSpace)-1 {
				binary.BigEndian.PutUint64(ebr.Part_next[:], uint64(ToInt(partitions[index].Part_next[:])))
			}
			ebr.Part_fit = byte(fit)
			binary.BigEndian.PutUint64(ebr.Part_size[:], uint64(tam))
			ebr.Part_status = 'A'
			for i := 0; i < len(name); i++ {
				ebr.Part_name[i] = name[i]
			}

			binary.BigEndian.PutUint64(partitions[index].Part_next[:], uint64(ToInt(ebr.Part_start[:])-structs.SizeOfEbr))

			file.Seek(ToInt(partitions[index].Part_start[:])-structs.SizeOfEbr, os.SEEK_SET)
			var buffer bytes.Buffer
			binary.Write(&buffer, binary.BigEndian, &(partitions[index]))
			writeBinary(file, buffer.Bytes())

			var buffer2 bytes.Buffer
			file.Seek(ToInt(ebr.Part_start[:])-structs.SizeOfEbr, os.SEEK_SET)
			binary.Write(&buffer2, binary.BigEndian, &ebr)
			writeBinary(file, buffer2.Bytes())

			WriteResponse(w, "PARTITION CREATED SUCCESFULLY")
		}
	} else {
		WriteResponse(w, "$Error: "+path+" does not exist")
	}
}
