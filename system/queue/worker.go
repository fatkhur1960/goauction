package queue

import (
	"log"
	"reflect"

	"github.com/gin-gonic/gin"
)

//Worker … simple worker that handles queueable tasks
type Worker struct {
	Name       string
	WorkerPool chan chan Queuable
	JobChannel chan Queuable
	quit       chan bool
}

// NewWorker --
func NewWorker(workerPool chan chan Queuable) Worker {
	job := make(chan Queuable, 10)
	return Worker{WorkerPool: workerPool, JobChannel: job}
}

//Start ... initiate worker to start listening for upcoming queueable jobs
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				if gin.IsDebugging() {
					log.Println("Worker] Got Event:", reflect.TypeOf(job).String())
				}
				// we have received a work request.
				if err := job.Handle(); err != nil {
					log.Println("Worker] Got Error:", err.Error())
				}
			}
		}
	}()
}
