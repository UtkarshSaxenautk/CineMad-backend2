package transport

import (
	"authentication-ms/pkg/svc"
	"log"
)

type Worker struct {
	s svc.SVC
}

// NewWorker Instance
func NewWorker(svc svc.SVC) (*Worker, error) {
	return &Worker{svc}, nil
}

func (w *Worker) Initialize() {
}

func (w *Worker) Run() {
	log.Println("Worker: starting")

}

// Shutdown invoked when service termination
func (w *Worker) Shutdown() {
}

