package temp

import (
	proto "github.com/temp/plugins"
)

//Worker struct is used to start and run the timer
type Worker struct {
	TimerInfo *proto.TimerInfo
}

//NewWorker is a constructor for worker
func NewWorker(timerInfo *proto.TimerInfo) *Worker {
	return &Worker{
		TimerInfo: timerInfo,
	}
}
