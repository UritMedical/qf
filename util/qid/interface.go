package qid

//
// IIdAllocator
//  @Description: Id分配器接口
//
type IIdAllocator interface {
	Next(name string) uint64
}
