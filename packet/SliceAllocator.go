// =====================================================================
//
// SliceAllocator.go -
//
// Last Modified: 2024/08/02 17:39:34
//
// =====================================================================
package packet

var sizelist = []int{
	32, 64, 128, 256, 512, 1024, 2048,
	4096, 8192, 16384, 32768, 65536, 65536 + 100}

// SliceAllocator is a global memory pool for slices.
var SliceAllocator *MemoryPool = NewMemoryPool(sizelist)

func SliceAlloc(size int) []byte {
	return SliceAllocator.Alloc(size)
}

func SliceFree(b []byte) {
	SliceAllocator.Free(b)
}
