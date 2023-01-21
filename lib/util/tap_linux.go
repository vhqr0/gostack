package util

import (
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	devPath = "/dev/net/tun"
	ifrSize = unix.IFNAMSIZ + 64
)

func OpenTAP(ifaceName string) (*os.File, error) {
	// open device
	fd, err := unix.Open(devPath, unix.O_RDWR|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}

	// prepare ifreq
	var ifr [ifrSize]byte
	ifaceBytes := []byte(ifaceName)
	if len(ifaceBytes) >= unix.IFNAMSIZ {
		unix.Close(fd)
		return nil, &IfaceNameOverflow{Name: ifaceName}
	}
	copy(ifr[:], ifaceBytes)
	*(*uint16)(unsafe.Pointer(&ifr[unix.IFNAMSIZ])) = unix.IFF_TAP | unix.IFF_NO_PI

	// set iface
	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(fd),
		uintptr(unix.TUNSETIFF),
		uintptr(unsafe.Pointer(&ifr[0])))
	if errno != 0 {
		unix.Close(fd)
		return nil, errno
	}

	// set nonblock
	if err := unix.SetNonblock(fd, true); err != nil {
		unix.Close(fd)
		return nil, err
	}

	// create file
	return os.NewFile(uintptr(fd), devPath), nil
}
