package wiremelt

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"wiremelt/shell"
	"wiremelt/worker"

	"github.com/fatih/color"
)

var sess SessionConfiguration

// Initialize Client, populate workspaces for workers to complete jobs
func InitClient(session *SessionConfiguration) {
	fmt.Println()

	// Remove previous input / training data
	removeCSV := func(filePath string) {
		_, err := os.Stat(filePath)
		if err != nil {
			_ = err // Ignore
		} else {
			removeTrainCSV := os.Remove(filePath)
			if removeTrainCSV != nil {
				log.Fatalln(removeTrainCSV)
			}
		}
	}
	removeCSV("neural/data/train.csv")

	sess = *session
	ctx, cancel := context.WithCancel(context.TODO())

	factories := session.FactoryFocus // via SessionConfig
	constructFactories(ctx, cancel, factories, &sess)
}

// Start of a lifecycle, Initialization of Factories
func constructFactories(ctx context.Context, cancel context.CancelFunc, factories map[int]string, session *SessionConfiguration) {
	for fID, targetFactory := range factories { // Loop over each factory from Config, instanciate new factory
		targetFactory := targetFactory // mandatory

		log.Println("] FACTORY INITIALIZED:", targetFactory)
		newFactory := Factory{ID: fID, Focus: targetFactory}
		SproutedFactories = append(SproutedFactories, newFactory)

		collector := StartDispatcher(ctx, newFactory, session.WorkerQuantity, session) // Start Worker Pool per instanciated factory

		for jID, job := range CreateJobs(session.JobsPerMacro, session) { // Create Jobs for workers per instanciated factory
			collector.JobQueue <- worker.Job{ID: jID, Macro: job.Macro, ParamArg: job.ParamArg} // Pass a new Job into the job queue for collector
		}
	}

	func() {
		defer cancel() // Cancel application context, killing all attached jobs
		color.HiRed("\nCLEANUP.")

		SproutedWorkers = nil
		SproutedJobs = nil
	}()

	defer func() {
		trainFilePath := "neural/data/train.csv"
		testFilePath := "neural/data/test.csv"

		repeatCycle := session.RepeatCycle != 0
		shellCycle := session.ShellCycle != 0
		neuralEnabled := session.NeuralEnabled != 0

		// RepeatCycle, Neural Network secures priority
		if repeatCycle && !neuralEnabled {
			defer InitClient(&sess)
			time.Sleep(1 * time.Second) // Wait a short amount of time to give wiremelt time to process the canceled context and finish running
		} else {
			if shellCycle {
				defer shell.InitShell(sess.MacroLibrary)
			}

			// Neural Network
			if neuralEnabled {
				// Copy input / training data
				trainFile, err := ioutil.ReadFile(trainFilePath)
				if err != nil {
					_ = err // Ignore
				}
				trainLines := string(trainFile)
				trainFileLines := strings.Split(trainLines, "\n")

				// Retrieve test data for line count
				testFile, err := ioutil.ReadFile(testFilePath)
				if err != nil {

					// Create Test file
					f, envInitErr := os.Create(testFilePath)
					if envInitErr != nil {
						log.Fatalln(err, envInitErr)
					}
					defer f.Close()

					testFile, _ = ioutil.ReadFile(testFilePath)

					_ = err // Ignore
				}
				testLines := string(testFile)
				testFileLines := strings.Split(testLines, "\n")

				// Ensure test data is only updated if under trainLimit
				if len(testFileLines) < session.TrainLimit { // 50-14K lines
					// Update test data with copied input / training data
					output := strings.Join(trainFileLines, "\n")
					testData, updateErr := os.OpenFile(testFilePath, os.O_APPEND|os.O_WRONLY, 0644)
					if updateErr != nil {
						log.Fatalln(updateErr)
					}
					defer testData.Close()

					update := fmt.Sprintf("\n%v", strings.TrimSpace(output))
					_, writeErr := testData.WriteString(update)
					if writeErr != nil {
						log.Fatalln(writeErr)
					}
				}
			}

			color.HiMagenta("DONE.")
			WiremeltAscii()
		}
	}()
}
