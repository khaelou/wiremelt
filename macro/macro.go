package macro

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/atedja/go-vector"
)

type ProductSignal struct {
	Product       interface{} // Return / Result
	WorkerID      int
	WorkerFactory string
	WorkerRole    string
	JobID         int
	Macro         string
	ParamArg      string
}

// Ensure return value of macro is not nil
func (ps *ProductSignal) QualityCheck() (bool, error) {
	var isValid bool
	var err error

	if ps.Product != nil {
		isValid = true
	} else {
		err = errors.New("macro returned a nil value")
		return !isValid, err
	}

	return isValid, err
}

// Converts strings to numeric type for neural network manipulations
func (ps *ProductSignal) InitNeuron() (float64, error) {
	// Convert Product to Numeric Value (string > bytes > vector > float64)
	parseProduct := fmt.Sprintf("%v", ps.Product) // interface{} to string
	parseBytes := []byte(parseProduct)            // string to bytes
	floatEncode := []float64{}
	for _, v := range parseBytes {
		floatEncode = append(floatEncode, float64(v)) // []byte to []float64
	}

	productVector := vector.NewWithValues(floatEncode) // []float64 to vector
	neuronInbound := productVector.Magnitude()

	fmt.Println("\t[✓✓] NEURON_INBOUND:", neuronInbound, "<<", ps.Product, "|", ps.WorkerRole, "@", ps.WorkerFactory)
	return neuronInbound, nil

	// Add conversion to CSV as training data
}

// Execute target built-in macro specified in MacroLibrary
func CallEmbedded(funcName string, params ...interface{}) (result interface{}, err error) {
	f := reflect.ValueOf(MacroLibrary[funcName])
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is out of index")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	var res []reflect.Value = f.Call(in)
	result = res[0].Interface()
	return
}

// Create ProductSignal instance from specified job, return new ProductSignal type from macro return value type interface{}
func ExecuteMacro(id int, factory, role string, jobID int, job string, paramArg string, execMacro interface{}) ProductSignal {
	product := fmt.Sprintf("%v", execMacro)                                                                                                                // Product represents return value of executed macro
	productSignal := ProductSignal{Product: product, WorkerID: id, WorkerFactory: factory, WorkerRole: role, JobID: jobID, Macro: job, ParamArg: paramArg} // Create new ProductSignal from Product

	return productSignal
}
