package gouring

import "syscall"

func USecToTimespec(uSec int64) (ts *syscall.Timespec) {
	ts = &syscall.Timespec{}
	uSecRoundedToSec := uSec / 1000_000 * 1000_000
	nSec := (uSec - uSecRoundedToSec) * 1000
	ts.Sec = uSecRoundedToSec / 1000_000
	ts.Nsec = nSec

	return
}
