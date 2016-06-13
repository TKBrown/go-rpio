// +build windows

package rpio

// Open - mocked function for Windows
func Open() (err error) {

	memlock.Lock()

	//give up 1k to mock the memory mapping
	mem = make([]uint32, 256)

	defer memlock.Unlock()

	return nil
}

// Close - mocked function for Windows
func Close() error {
	memlock.Lock()
	mem = nil

	defer memlock.Unlock()
	return nil
}
