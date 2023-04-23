package gouring

import (
	"syscall"
	"unsafe"
)

func (ring *IoUring) do_register(opcode uint32, arg unsafe.Pointer, nrArgs uint32) (ret int, err error) {
	fd := ring.RingFd
	if ring.IntFlags&INT_FLAG_REG_REG_RING != 0 {
		fd = ring.EnterRingFd
		opcode |= IORING_REGISTER_USE_REGISTERED_RING
	}

	return io_uring_register(fd, uint32(opcode), arg, nrArgs)
}

func (ring *IoUring) io_uring_register_buffers_update_tag(off uint32,
	iov *syscall.Iovec,
	tags []uint64,
	nr uint32) error {
	up := &IoUringRsrcUpdate2{
		Offset: off,
		Data:   uint64(uintptr(unsafe.Pointer(iov))),
		Tags:   uint64(uintptr(unsafe.Pointer(&tags[0]))),
		Nr:     nr,
	}

	ret, err := ring.do_register(IORING_REGISTER_BUFFERS_UPDATE, unsafe.Pointer(up), uint32(unsafe.Sizeof(*up)))
	if err != nil {
		return err
	}
	_ = ret
	return nil
}

func (ring *IoUring) io_uring_register_buffers_tags(
	iov *syscall.Iovec,
	tags []uint64,
	nr uint32) error {
	reg := &IoUringRsrcRegister{
		Nr:   nr,
		Data: uint64(uintptr(unsafe.Pointer(iov))),
		Tags: uint64(uintptr(unsafe.Pointer(&tags[0]))),
	}
	ret, err := ring.do_register(IORING_REGISTER_BUFFERS2, unsafe.Pointer(reg), uint32(unsafe.Sizeof(*reg)))
	if err != nil {
		return err
	}
	_ = ret
	return nil
}

func (ring *IoUring) io_uring_register_buffers_sparse(nr uint32) error {
	reg := &IoUringRsrcRegister{
		Flags: IORING_RSRC_REGISTER_SPARSE,
		Nr:    nr,
	}
	ret, err := ring.do_register(IORING_RSRC_REGISTER_SPARSE, unsafe.Pointer(reg), uint32(unsafe.Sizeof(*reg)))
	if err != nil {
		return err
	}
	_ = ret
	return nil
}

func (ring *IoUring) io_uring_register_buffers(iov *syscall.Iovec, nrIov uint32) int {
	ret, err := ring.do_register(IORING_REGISTER_BUFFERS, unsafe.Pointer(iov), nrIov)
	if err != nil {
		return 0
	}
	return ret
}

func (ring *IoUring) io_uring_unregister_buffers() int {
	ret, err := ring.do_register(IORING_UNREGISTER_BUFFERS, nil, 0)
	if err != nil {
		return 0
	}
	return ret
}

func (ring *IoUring) io_uring_register_files_update_tag(off uint32,
	files []int, tags []uint64,
	nrFiles uint32) (int, error) {
	up := &IoUringRsrcUpdate2{
		Offset: off,
		Data:   uint64(uintptr(unsafe.Pointer(&files[0]))),
		Tags:   uint64(uintptr(unsafe.Pointer(&tags[0]))),
		Nr:     nrFiles,
	}
	return ring.do_register(IORING_REGISTER_FILES_UPDATE2, unsafe.Pointer(up), uint32(unsafe.Sizeof(*up)))
}

func (ring *IoUring) io_uring_register_files_update(off uint32,
	files []int, nrFiles uint32) (int, error) {
	up := &IoUringFilesUpdate{
		Offset: off,
		Fds:    uint64(uintptr(unsafe.Pointer(&files[0]))),
	}
	return ring.do_register(IORING_REGISTER_FILES_UPDATE,
		unsafe.Pointer(up), nrFiles)
}

func (ring *IoUring) io_uring_register_files_sparse(nr uint32) (ret int, err error) {
	reg := &IoUringRsrcRegister{
		Flags: IORING_RSRC_REGISTER_SPARSE,
		Nr:    nr,
	}
	var didIncrease bool
	for {
		ret, err = ring.do_register(IORING_REGISTER_FILES2, unsafe.Pointer(reg), uint32(unsafe.Sizeof(*reg)))
		if err == nil {
			break
		}
		if err == syscall.EMFILE && !didIncrease {
			increase_rlimit_nofile(uint64(nr))
			didIncrease = true
			continue
		}
		break
	}
	return
}

func (ring *IoUring) io_uring_register_files_tags(
	files []int,
	tags []uint64, nr uint32) (ret int, err error) {
	reg := &IoUringRsrcRegister{
		Nr:   nr,
		Data: uint64(uintptr(unsafe.Pointer(&files[0]))),
		Tags: uint64(uintptr(unsafe.Pointer(&tags[0]))),
	}
	var didIncrease bool
	for {
		ret, err = ring.do_register(IORING_REGISTER_FILES2, unsafe.Pointer(reg), uint32(unsafe.Sizeof(*reg)))
		if err == nil {
			break
		}
		if err == syscall.EMFILE && !didIncrease {
			increase_rlimit_nofile(uint64(nr))
			didIncrease = true
			continue
		}
		break
	}
	return
}

func (ring *IoUring) io_uring_register_files(files []int, nrFiles uint32) (ret int, err error) {
	var didIncrease bool
	for {
		ret, err = ring.do_register(IORING_REGISTER_FILES, unsafe.Pointer(&files[0]), nrFiles)
		if err == nil {
			break
		}
		if err == syscall.EMFILE && !didIncrease {
			increase_rlimit_nofile(uint64(nrFiles))
			didIncrease = true
			continue
		}
		break
	}
	return
}

func (ring *IoUring) io_uring_unregister_files() error {
	_, err := ring.do_register(IORING_UNREGISTER_FILES, nil, 0)
	return err
}

func (ring *IoUring) io_uring_register_eventfd(eventFd int) error {
	_, err := ring.do_register(IORING_REGISTER_EVENTFD, unsafe.Pointer(&eventFd), 1)
	return err
}

func (ring *IoUring) io_uring_unregister_eventfd() error {
	_, err := ring.do_register(IORING_UNREGISTER_EVENTFD, nil, 0)
	return err
}

func (ring *IoUring) io_uring_register_eventfd_async(eventFd int) error {
	_, err := ring.do_register(IORING_REGISTER_EVENTFD_ASYNC, unsafe.Pointer(&eventFd), 1)
	return err
}

func (ring *IoUring) io_uring_register_probe(p *IoUringProbe, nrOps uint32) error {
	_, err := ring.do_register(IORING_REGISTER_PROBE, unsafe.Pointer(p), nrOps)
	return err
}

func (ring *IoUring) io_uring_register_personality() (int, error) {
	return ring.do_register(IORING_REGISTER_PERSONALITY, nil, 0)
}

func (ring *IoUring) io_uring_unregister_personality(id uint32) (int, error) {
	return ring.do_register(IORING_UNREGISTER_PERSONALITY, nil, id)
}

func (ring *IoUring) io_uring_register_restrictions(res *IoUringRestriction, nrRes uint32) error {
	_, err := ring.do_register(IORING_REGISTER_RESTRICTIONS, unsafe.Pointer(res), nrRes)
	return err
}

func (ring *IoUring) io_uring_enable_rings() error {
	_, err := ring.do_register(IORING_REGISTER_ENABLE_RINGS, nil, 0)
	return err
}

// sched.h
// func io_uring_register_iowq_aff(ring *IoUring, cpuSz int, mask *CpuSet) {
// }
func (ring *IoUring) io_uring_unregister_iowq_aff() error {
	_, err := ring.do_register(IORING_UNREGISTER_IOWQ_AFF, nil, 0)
	return err
}

func (ring *IoUring) io_uring_register_iowq_max_workers(val *uint32) (int, error) {
	return ring.do_register(IORING_REGISTER_IOWQ_MAX_WORKERS, unsafe.Pointer(val), 2)
}

func (ring *IoUring) io_uring_register_ring_fd() (int, error) {
	up := &IoUringRsrcUpdate{
		Data:   uint64(ring.RingFd),
		Offset: ^uint32(0),
	}

	ret, err := ring.do_register(IORING_REGISTER_RING_FDS, unsafe.Pointer(up), 1)
	if err != nil {
		return 0, err
	}

	ring.EnterRingFd = int(up.Offset)
	ring.IntFlags |= INT_FLAG_REG_RING
	if ring.Features&IORING_FEAT_REG_REG_RING != 0 {
		ring.IntFlags |= INT_FLAG_REG_REG_RING
	}

	return ret, nil
}

func (ring *IoUring) io_uring_unregister_ring_fd() error {
	up := &IoUringRsrcUpdate{
		Offset: uint32(ring.EnterRingFd),
	}
	ret, err := ring.do_register(IORING_UNREGISTER_RING_FDS, unsafe.Pointer(up), 1)
	if err != nil {
		return err
	}
	if ret == 1 {
		ring.EnterRingFd = ring.RingFd
		ring.IntFlags &= ^(INT_FLAG_REG_RING | INT_FLAG_REG_REG_RING)
	}
	return nil
}

func (ring *IoUring) io_uring_register_buf_ring(reg *IoUringBufReg, flags uint32) (int, error) {
	return ring.do_register(IORING_REGISTER_PBUF_RING, unsafe.Pointer(reg), 1)
}

func (ring *IoUring) io_uring_unregister_buf_ring(bgId int32) (int, error) {
	reg := &IoUringBufReg{
		Bgid: uint16(bgId),
	}
	return ring.do_register(IORING_UNREGISTER_PBUF_RING, unsafe.Pointer(reg), 1)
}
