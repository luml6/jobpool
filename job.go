package jobpool

// use exmple
/*type YourJob struct {
	Name string
	Age  int
}

func (c *YourJob) Do() error {
	fmt.Println("name is %v,age is %v", c.Name, c.Age)
	return nil
}

func main() {
	dispath := jobpool.NewDispatcher(conf.Conf.MaxWorker, conf.Conf.MaxWorker)
	dispath.Run()
	youjob := &YourJob{Name: "albert", Age: 12}
	dispath.Add(youjob)
}
*/
import (
	"fmt"
)

// 调用的工作方法
type Job interface {
	Do() error
}

const (
	MAXWORKERS  = 5
	MAXJOBQUEUE = 10
)

// job 工作者worker
type Worker struct {
	WorkerPool chan chan Job //woker pool
	JobChannel chan Job
	quit       chan bool
}

// create new worker
func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

func (w *Worker) Start() {
	w.WorkerPool <- w.JobChannel
	go func() {
		for {
			select {
			case job := <-w.JobChannel:
				if err := job.Do(); err != nil {
					fmt.Println("excut job failed with error:%v", err.Error())
				}
			case <-w.quit:
				return

			}

			w.WorkerPool <- w.JobChannel
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

// 总调度者安排工作者工作
type Dispatcher struct {
	WorkerPool chan chan Job
	JobQueue   chan Job
	maxWorkers int
	quit       chan bool
}

//create new dispatcher
func NewDispatcher(maxWorkers, maxJobQueue int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	queue := make(chan Job, maxJobQueue)
	return &Dispatcher{WorkerPool: pool, JobQueue: queue, maxWorkers: maxWorkers, quit: make(chan bool)}

}

// add job into jobqueue
func (d *Dispatcher) Add(job Job) {
	d.JobQueue <- job
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)

		worker.Start()
	}

	go func() {
		for {
			select {
			case job := <-d.JobQueue:
				go func(job Job) {
					jobChan := <-d.WorkerPool
					jobChan <- job
				}(job)
			// stop dispatcher
			case <-d.quit:
				return
			}
		}
	}()
}

func (d *Dispatcher) Stop() {
	go func() {
		d.quit <- true
	}()
}
