package wiremelt

import (
	"context"
	"fmt"
	"log"
	"time"

	"wiremelt/shell"
	"wiremelt/worker"
)

var sess SessionConfiguration

// 	Initialize Client, populate workspaces for workers to complete jobs
func InitClient(session *SessionConfiguration) {
	fmt.Println()
	fmt.Println("CLIENT INIT")

	sess = *session
	ctx, cancel := context.WithCancel(context.TODO())

	factories := session.FactoryFocus // via SessionConfig
	constructFactories(ctx, cancel, factories, &sess)
}

// Start of a lifecycle, Initialization of Factories
func constructFactories(ctx context.Context, cancel context.CancelFunc, factories map[int]string, session *SessionConfiguration) {
	for fID, targetFactory := range factories { // Loop over each factory from Config, instanciate new factory
		targetFactory := targetFactory // mandatory

		fmt.Println()
		log.Println("] FACTORY INITIALIZED:", targetFactory)
		newFactory := Factory{ID: fID, Focus: targetFactory}
		SproutedFactories = append(SproutedFactories, newFactory)

		collector := StartDispatcher(ctx, newFactory, session.WorkerQuantity, session) // Start Worker Pool per instanciated factory

		for jID, job := range CreateJobs(session.JobsPerFactory, session) { // Create Jobs for workers per instanciated factory
			collector.JobQueue <- worker.Job{ID: jID, Macro: job.Macro, ParamArg: job.ParamArg} // Pass a new Job into the job queue for collector
		}
	}

	func() {
		defer cancel() // Cancel application context, killing all attached jobs
		fmt.Println("\nCLEANUP.")

		SproutedWorkers = nil
		SproutedJobs = nil
	}()

	defer func() {
		repeatCycle := session.RepeatCycle != 0
		shellCycle := session.ShellCycle != 0

		if repeatCycle {
			defer InitClient(&sess)
			time.Sleep(1 * time.Second) // Wait a short amount of time to give wiremelt time to process the canceled context and finish running
		} else {
			if shellCycle {
				defer shell.InitShell(sess.MacroLibrary)
			}

			fmt.Println("DONE.")
			fmt.Println()
		}
	}()
}
