package gouring

/*--- setup ---*/

func New(entries uint32, flags uint32) (*IoUring, error) {
	ring := &IoUring{}
	p := new(IoUringParams)
	p.Flags = flags
	err := io_uring_queue_init_params(entries, ring, p)
	if err != nil {
		return nil, err
	}
	return ring, nil
}

func NewWithParams(entries uint32, params *IoUringParams) (*IoUring, error) {
	ring := &IoUring{}
	if params == nil {
		params = new(IoUringParams)
	}
	err := io_uring_queue_init_params(entries, ring, params)
	if err != nil {
		return nil, err
	}
	return ring, nil
}

/*--- sq ---*/

func (h *IoUring) GetSqe() *IoUringSqe {
	return h.io_uring_get_sqe()
}

/*--- cq ---*/

func (h *IoUring) SeenCqe(cqe *IoUringCqe) {
	h.io_uring_cqe_seen(cqe)
}

func (h *IoUring) Advance(nr uint32) {
	h.io_uring_cq_advance(nr)
}

func (h *IoUring) PeekBatchCqe(cqes []*IoUringCqe, count uint32) uint32 {
	return h.io_uring_peek_batch_cqe(cqes, count)
}

/*--- enter ---*/

func (h *IoUring) Submit() (int, error) {
	return h.io_uring_submit()
}

func (h *IoUring) SubmitAndWait(waitNr uint32) (int, error) {
	return h.io_uring_submit_and_wait(waitNr)
}

func (h *IoUring) WaitCqe(cqePtr **IoUringCqe) error {
	return h.io_uring_wait_cqe(cqePtr)
}

func (h *IoUring) SubmitAndWaitTimeOut(cqePtr **IoUringCqe, waitNtr uint32, uSec int64, sigmask *Sigset_t) error {
	return h.io_uring_submit_and_wait_timeout(cqePtr, waitNtr, USecToTimespec(uSec), sigmask)
}

/*--- register ---*/

func (h *IoUring) RegisterRingFD() (int, error) {
	return h.io_uring_register_ring_fd()
}

func (h *IoUring) RegisterEventFd(efd int) error {
	return h.io_uring_register_eventfd(efd)
}

func (h *IoUring) UnRegisterEventFd() error {
	return h.io_uring_unregister_eventfd()
}

/*--- free sqe sq cq mmap and close ring fd ---*/
func (h *IoUring) Close() {
	h.io_uring_queue_exit()
}
