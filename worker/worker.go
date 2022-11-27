package worker

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"

	"wiremelt/macro"
	"wiremelt/utils"
)

type MacroSpec map[string]string
type MacroMapping map[string]interface{}

var MacroFileDir = "custom/"      // Directory to store external macros (JavaScript)
var MacroSpecs = MacroSpec{}      // Map to track macros from SessionConfig
var MacroLibrary = MacroMapping{} // Map for default macros

var ProductChannel = make(chan macro.ProductSignal, math.MaxInt8)

type Job struct {
	ID       int
	Macro    string
	ParamArg string
}

type Worker struct {
	ID            int
	Factory       string
	Role          string
	WorkerChannel chan chan Job // Used to communicate between dispatcher and workers
	JobChannel    chan Job      // Used to communicate between job and workers
	EndShift      chan bool     // Used to communicate between workers and dispatcher
}

func checkProductQuality(ctx context.Context, job Job, productSignal macro.ProductSignal) {
	qualityCheck, err := productSignal.QualityCheck()
	if err != nil {
		log.Println("checkProductQuality error:", err)
	}

	if qualityCheck { // Product Validated for Neural Network
		ProductChannel <- productSignal // Send Product to ProductChannel
		neuron, accuracy, err := productSignal.InitNeuron()
		if err != nil {
			log.Println("initNeuron error:", err)
		}

		qc := fmt.Sprintf("\t[!✓] #%d NEURON_RESULT: (macro.%s) Accuracy = %f [Neuron \"%v\" - %v]", job.ID, job.Macro, accuracy, productSignal.Product, neuron)
		fmt.Println(qc)
	} else {
		fmt.Println("\t\t\t[xX] QUALITY CHECK: product is nil!")
	}
}

// Start Worker
func (w *Worker) StartWorker(ctx context.Context, useV8Isolates bool) {
	neuralEnabled := ctx.Value("neuralEnabled") != 0

	go func() {
		for {
			w.WorkerChannel <- w.JobChannel // When the worker is available place channel in queue

			select {
			case job := <-w.JobChannel: // Worker has received job
				job.Macro = strings.Replace(job.Macro, "*", "", -1) // Remove '*' signifing ignore of ParamArg

				// Check if Default Macro or Custom Macro
				if !strings.Contains(job.ParamArg, ".js") { // Built-in Macro
					executeMacro, err := macro.CallEmbedded(job.Macro, job.ParamArg)
					if err != nil {
						log.Fatalln("macroIdentity error:", err)
					}

					var execMacro interface{} = executeMacro
					product := macro.ExecuteMacro(w.ID, w.Factory, w.Role, job.ID, job.Macro, job.ParamArg, execMacro) // Macro execution, returns product or nil

					ctx = context.WithValue(ctx, job, execMacro) // parent context, key, value

					if neuralEnabled {
						checkProductQuality(ctx, job, product)
					} else {
						qcBypass := fmt.Sprintf("[✓][%d] %s .: %s @ %s :: (#%d) PRODUCT = \"%v\"", w.ID, job.Macro, w.Role, w.Factory, job.ID, execMacro)
						fmt.Println(qcBypass)
					}
				} else { // Custom / External Macro
					jsScript := fmt.Sprintf("%s%s", MacroFileDir, job.ParamArg) // custom/macro.js
					if _, err := os.Stat(jsScript); !os.IsNotExist(err) {
						// Parse file
						fileContent, err := ioutil.ReadFile(jsScript)
						if err != nil {
							log.Fatalln("customMacro error:", err)
						}
						convScript := string(fileContent) // Convert []byte to string

						var execJS interface{}

						// Depict when to use V8 Isolates > Node.js (V8Go memory limit)
						if useV8Isolates { // V8 Isolates
							execJS, err = utils.V8Isolates(convScript, false) // Execute script via V8 Isolates (speedy execution, yet JS heap / stack trace errors)
							if err != nil {
								log.Println("execJS error:", err)
							}
						} else { // Node.js
							execJS, err = utils.V8NodeJS(convScript, false) // Execute script via Node.js
							if err != nil {
								log.Println("execNodeJS error:", err)
							}
						}

						var execMacro interface{} = execJS
						product := macro.ExecuteMacro(w.ID, w.Factory, w.Role, job.ID, job.Macro, jsScript, execMacro) // Macro execution, returns product or nil

						ctx = context.WithValue(ctx, job, execMacro) // parent context, key, value

						if neuralEnabled {
							checkProductQuality(ctx, job, product)
						} else {
							qcBypass := fmt.Sprintf("[✓][%d] %s .: %s @ %s :: (#%d) PRODUCT = \"%v\"", w.ID, job.Macro, w.Role, w.Factory, job.ID, execMacro)
							fmt.Println(qcBypass)
						}
					}
				}
			case <-w.EndShift: // Worker has completed job, return to WorkerChannel to wait for the next available job
				<-ctx.Done()
				return
			}
		}
	}()
}

// End Worker
func (w *Worker) StopWorker(ctx context.Context) {
	log.Printf("Worker [%d @ %s] has halted!", w.ID, w.Factory)
	w.EndShift <- true
}
