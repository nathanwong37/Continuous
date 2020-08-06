package temp

import (
	"errors"
	"fmt"
	"math"
	"time"

	proto "github.com/temp/plugins"
)

//Worker struct is used to start and run the timer
type Worker struct {
	manager   *Manager
	TimerInfo *proto.TimerInfo
	finished  chan bool
}

//NewWorker is a constructor for worker
func NewWorker(timerInfo *proto.TimerInfo, mngr *Manager) *Worker {
	return &Worker{
		TimerInfo: timerInfo,
		manager:   mngr,
	}
}

//RunTimer is used for the worker to run the timer
func (worker *Worker) RunTimer(sleep time.Duration, curr int) error {
	//Set a channel up for delete, or if the timer is finished
	worker.finished = make(chan bool)
	seconds, err := parseInterval(worker.TimerInfo.GetInterval())
	if err != nil {
		//Chances of error should be zero... should be authenticated at api
		fmt.Println("Failed to parse interval")
	}
	count := int(worker.TimerInfo.GetCount())

	go func(dur, count, curr int, uuidstr, namespace string, done chan bool, sleep time.Duration, manager *Manager) {
		time.Sleep(sleep)
		ticker := time.NewTicker(time.Second * time.Duration(dur))
		var buffer int = 0
		var pass int = 0
		transporter := Transport{}
		fmt.Println("STARTING TIMER ON " + uuidstr)
		for {
			select {
			case val := <-done:
				//worker is finished or timer deleted
				if val {
					manager.DeleteEntry(uuidstr, namespace)
					return
				}
				fmt.Println("STOPPING TIMER ON " + uuidstr)
				return
			case fin := <-ticker.C:
				fmt.Println("tick at", fin)
				curr++
				if dur <= 5 {
					pass++
				} else {
					pass = buffer
				}
				if pass >= buffer {
					_, errs := transporter.Update(uuidstr, fin.Format("2006-01-02 15:04:05"), namespace, int(curr))
					if errs != nil {
						fmt.Println("Failed to update")
					}
					pass = 0
					buffer = int(math.Min(float64(buffer+1), float64(10%dur)))
				}
				//timer over
				if curr >= count && count != -1 {
					//final update
					_, errs := transporter.Update(uuidstr, fin.Format("2006-01-02 15:04:05"), namespace, int(curr))
					if errs != nil {
						//failed to write, try 3 times then give up, entry will be deleted anyways
						for i := 0; i < 3; i++ {
							_, errs := transporter.Update(uuidstr, fin.Format("2006-01-02 15:04:05"), namespace, int(curr))
							if errs == nil {
								break
							}
							time.Sleep(5 * time.Millisecond)
						}
					}
					done <- true
				}

			}

		}
	}(seconds, count, curr, worker.TimerInfo.GetTimerID(), worker.TimerInfo.GetNameSpace(), worker.finished, sleep, worker.manager)
	return nil

}

//CallBack is used to start the run timer, when it is appropriate
func (worker *Worker) CallBack(start string, run func(time.Duration, int) error) error {
	//timer starts now
	if start == "" {
		return run(time.Second*0, 0)
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", start, time.Local)
	if err != nil {
		return errors.New("Incorrect format")
	}
	//Timer got picked up, have to sync
	if t.Before(time.Now().Local()) {
		curr := worker.TimerInfo.GetAmountFired()
		seconds, err := parseInterval(worker.TimerInfo.GetInterval())
		if err != nil {
			return err
		}
		for t.Add(time.Duration(seconds)*time.Second).Before(time.Now().Local()) || (curr > worker.TimerInfo.GetCount() && worker.TimerInfo.GetCount() != -1) {
			t = t.Add(time.Second * time.Duration(seconds))
			curr++

			//call whatever happens when a timer
		}
		//timer is expired
		if int32(curr) > worker.TimerInfo.GetCount() && worker.TimerInfo.GetCount() != -1 {
			//delete from database then return an error
			worker.manager.remove(worker.TimerInfo.GetTimerID(), worker.TimerInfo.GetNameSpace())
			return errors.New("Timer has expired")
		}
		transporter := Transport{}
		_, errs := transporter.Update(worker.TimerInfo.GetTimerID(), t.Format("2006-01-02 15:04:05"), worker.TimerInfo.GetNameSpace(), int(curr))
		if errs != nil {
			fmt.Println(errs.Error())
			return errs
		}
		curr++
		return run(time.Now().Local().Sub(t), int(curr))
	}

	//timer starts later
	return run(t.Sub(time.Now().Local()), 0)
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

//SendSignal is used to send the signal to worker to stop timer
func (worker *Worker) SendSignal(done bool) {
	worker.finished <- done
}
