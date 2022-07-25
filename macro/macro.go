package macro

import (
	"errors"
	"fmt"
	"reflect"
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
