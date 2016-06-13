// +build linux

package rpio

import (
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

// Open and memory map GPIO memory range from /dev/mem .
// Some reflection magic is used to convert it to a unsafe []uint32 pointer
func Open() (err error) {
	var file *os.File
	var base int64

	// Open fd for rw mem access; try gpiomem first
	if file, err = os.OpenFile(
		"/dev/gpiomem",
		os.O_RDWR|os.O_SYNC,
		0); os.IsNotExist(err) {
		file, err = os.OpenFile(
			"/dev/mem",
			os.O_RDWR|os.O_SYNC,
			0)
		base = getGPIOBase()
	}

	if err != nil {
		return
	}

	// FD can be closed after memory mapping
	defer file.Close()

	memlock.Lock()
	defer memlock.Unlock()

	// Memory map GPIO registers to byte array
	mem8, err = syscall.Mmap(
		int(file.Fd()),
		base,
		memLength,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED)

	if err != nil {
		return
	}

	// Convert mapped byte memory to unsafe []uint32 pointer, adjust length as needed
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&mem8))
	header.Len /= (32 / 8) // (32 bit = 4 bytes)
	header.Cap /= (32 / 8)

	mem = *(*[]uint32)(unsafe.Pointer(&header))

	return nil
}

// Close unmaps GPIO memory
func Close() error {
	memlock.Lock()
	defer memlock.Unlock()
	return syscall.Munmap(mem8)
}
