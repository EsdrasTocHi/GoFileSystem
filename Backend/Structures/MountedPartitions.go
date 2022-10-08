package structures

type MountedPartition struct {
	Par      Partition
	LogicPar Ebr
	IsLogic  bool
	Path     string
	Id       string
}
