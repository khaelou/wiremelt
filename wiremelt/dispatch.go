package wiremelt

import (
	"context"
	"fmt"
	"log"
	"math"
	"runtime"

	"wiremelt/worker"
)

var SproutedFactories []Factory     // Array to track instanciated factories
var SproutedWorkers []worker.Worker // Array to track instanciated workers
var SproutedJobs []worker.Job       // Array to track instanciated jobs

var FactoryChannel = make(chan chan Factory, math.MaxInt8)   // Stores all of the channels of available workers
var WorkerChannel = make(chan chan worker.Job, math.MaxInt8) // Stores all of the channels of available workers

type Factory struct {
	ID    int
	Focus string
}

type Collector struct {
	JobQueue chan worker.Job // Receives jobs to send to workers
	EndCycle chan bool       // Receives signal to stop workers
}

// Dispatcher instantiates and connects all of the workers with the Factory pool
func StartDispatcher(ctx context.Context, targetFactory Factory, workerCount int, session *SessionConfiguration) Collector {
	var i int

	input := make(chan worker.Job) // Channel to receive jobs
	end := make(chan bool)         // Channel to halt workers
	collector := Collector{JobQueue: input, EndCycle: end}

	roles := session.WorkerRoles // via SessionConfig

	for i < workerCount {
		for _, role := range roles {
			i++

			log.Printf("~ Starting Worker #%d (%s @ %s)\n", i, role, targetFactory.Focus)

			// Add New Worker
			worker := worker.Worker{
				ID:            i,
				Factory:       targetFactory.Focus,
				Role:          role,
				WorkerChannel: WorkerChannel,
				JobChannel:    make(chan worker.Job),
				EndShift:      make(chan bool),
			}

			useV8Isolates := true
			if session.RepeatCycle == 1 {
				useV8Isolates = false
			}

			worker.StartWorker(ctx, useV8Isolates)            // Worker, grabs a waiting job and then does it's task
			SproutedWorkers = append(SproutedWorkers, worker) // Store Worker for reference
		}
	}

	// Collector; Receives jobs and pushes them to the job queue for available workers
	go func() {
		for {
			select {
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					_ = fmt.Sprintf("\nCTX DONE: %v", err)
				}

				// Close channels
				close(collector.JobQueue)
				close(collector.EndCycle)

				//fmt.Printf("COMPLETE.\n\n")
				return
			case <-end:
				for _, w := range SproutedWorkers {
					w.StopWorker(ctx) // Stop worker
				}
			case signal := <-worker.ProductChannel:
				macroID := fmt.Sprintf("macro.%s", signal.Macro)
				fmt.Printf("[✓][%d] %s .: %s @ %s :: (#%d) PRODUCT = \"%v\"\n", signal.WorkerID, macroID, signal.WorkerRole, signal.WorkerFactory, signal.JobID, signal.Product)
			case job := <-input:
				worker := <-WorkerChannel // Wait for available worker on channel
				worker <- job             // Dispatch job to worker waiting on channel
			}
		}
	}()

	return collector
}

// Create Workload of Jobs
func CreateJobs(amount int, session *SessionConfiguration) []worker.Job {
	fmt.Println("\n~ Active Threads:", runtime.NumGoroutine())
	fmt.Println("+ Macros:", session.MacroLibrary)
	fmt.Println()

	for i := 0; i < amount; i++ {
		for macro, paramArg := range session.MacroLibrary {
			newJob := worker.Job{ID: i, Macro: macro, ParamArg: paramArg}
			SproutedJobs = append(SproutedJobs, newJob)
		}
	}

	return SproutedJobs
}
