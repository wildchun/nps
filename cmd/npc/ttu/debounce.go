// Package ttu
//
//	@Description
//	@Author  hd_0411_qxc  2025/1/10 15:52
//	@Update  hd_0411_qxc  2025/1/10 15:52
package ttu

import (
	"time"
)

func NewDebounce(after time.Duration) func(f func()) {
	var timer *time.Timer
	return func(f func()) {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(after, f)
	}
}
