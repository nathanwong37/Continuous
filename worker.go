package temp

import (
	"errors"
	"fmt"
	"time"

	proto "github.com/temp/plugins"
)

//Worker struct is used to start and run the timer
type Worker struct {
	TimerInfo *proto.TimerInfo
	finished  chan bool
}

//NewWorker is a constructor for worker
func NewWorker(timerInfo *proto.TimerInfo) *Worker {
	return &Worker{
		TimerInfo: timerInfo,
	}
}

//RunTimer is used for the worker to run the timer
func (worker *Worker) RunTimer() error {
	var curr int = int(worker.TimerInfo.GetCount())
	//Set a channel up for delete, or if the timer is finished
	worker.finished = make(chan bool)
	seconds, err := parseInterval(worker.TimerInfo.GetInterval())
	if err != nil {
		//Chances of error should be zero... should be authenticated at api
		fmt.Println("Failed to parse interval")
	}
	//dealing with recent time
	if worker.TimerInfo.GetMostRecent() != "" {
		t, _ := time.Parse("2006-01-02 15:04:05", worker.TimerInfo.GetMostRecent())
		for t.Add(time.Duration(seconds)).Before(time.Now()) {
			t.Add(time.Duration(seconds))
			curr++
		}
		//timer is expired
		if int32(curr) > worker.TimerInfo.GetCount() && worker.TimerInfo.GetCount() != -1 {
			return errors.New("Timer is expired")

		}
		transporter := Transport{}
		_, err := transporter.Update(worker.TimerInfo.GetTimerID(), worker.TimerInfo.GetMostRecent(), worker.TimerInfo.GetNameSpace(), curr)
		if err != nil {
			fmt.Println("Failed to write to database")
		}
		time.Sleep(time.Now().Sub(t))
		curr++
	}
	//sleep then run go routine when it is time
	//start ticker go routine
	count := int(worker.TimerInfo.GetCount())
	go func(dur, count, curr int) {
		ticker := time.NewTicker(time.Second * time.Duration(dur))
		for {
			select {
			case <-worker.finished:
				//worker is finished or timer deleted
				return
			case fin := <-ticker.C:
				fmt.Println("tick at", fin)
				count++
				if curr >= count {
					//timer over
				}

			}

		}
	}(seconds, count, curr)
	return nil

}

//CallBack is used to start the run timer, when it is appropriate
func (worker *Worker) CallBack(start string, run func() error) error {
	t, err := time.Parse("2006-01-02 15:04:05", start)
	if err != nil {
		return run()
	}
	if t.Before(time.Now()) {
		return run()
	}
	diff := t.Sub(time.Now())
	time.Sleep(diff)
	return run()
}

//ParseInterval parses the string into seconds
func parseInterval(interval string) (int, error) {
	var hour, minute, second int
	num, err := fmt.Sscanf(interval, "%d:%d:%d", &hour, &minute, &second)
	if err != nil || num != 3 {
		//err scanning
		return -1, err
	}
	return hour*3600 + minute*60 + second, nil
}
