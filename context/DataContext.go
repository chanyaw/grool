package context

import (
	"fmt"
	"github.com/newm4n/grool/pkg"
	"reflect"
	"strings"
)

func NewDataContext() *DataContext {
	return &DataContext{
		ObjectStore: make(map[string]interface{}),
	}
}

type DataContext struct {
	ObjectStore map[string]interface{}
}

func (ctx *DataContext) Add(key string, obj interface{}) {
	ctx.ObjectStore[key] = obj
}

func (ctx *DataContext) GetType(variable string) (reflect.Type, error) {
	varArray := strings.Split(variable, ".")
	if val, ok := ctx.ObjectStore[varArray[0]]; ok {
		return traceType(val, varArray[1:])
	} else {
		return nil, fmt.Errorf("data context not found '%s'", varArray[0])
	}
}

func (ctx *DataContext) GetValue(variable string) (reflect.Value, error) {
	varArray := strings.Split(variable, ".")
	if val, ok := ctx.ObjectStore[varArray[0]]; ok {
		return traceValue(val, varArray[1:])
	} else {
		return reflect.ValueOf(nil), fmt.Errorf("data context not found '%s'", varArray[0])
	}
}

func (ctx *DataContext) SetValue(variable string, newValue interface{}) error {
	varArray := strings.Split(variable, ".")
	if val, ok := ctx.ObjectStore[varArray[0]]; ok {
		return traceSetValue(val, varArray[1:], newValue)
	} else {
		return fmt.Errorf("data context not found '%s'", varArray[0])
	}
}

func traceType(obj interface{}, path []string) (reflect.Type, error) {
	if len(path) == 1 {
		return pkg.GetAttributeType(obj, path[0])
	} else if len(path) > 1 {
		objVal, err := pkg.GetAttributeValue(obj, path[0])
		if err != nil {
			return nil, err
		}
		return traceType(pkg.ValueToInterface(objVal), path[1:])
	} else {
		return reflect.TypeOf(obj), nil
	}
}

func traceValue(obj interface{}, path []string) (reflect.Value, error) {
	if len(path) == 1 {
		return pkg.GetAttributeValue(obj, path[0])
	} else if len(path) > 1 {
		objVal, err := pkg.GetAttributeValue(obj, path[0])
		if err != nil {
			return objVal, err
		}
		return traceValue(pkg.ValueToInterface(objVal), path[1:])
	} else {
		return reflect.ValueOf(obj), nil
	}
}

func traceSetValue(obj interface{}, path []string, newValue interface{}) error {
	if len(path) == 1 {
		return pkg.SetAttributeInterface(obj, path[0], newValue)
	} else if len(path) > 1 {
		objVal, err := pkg.GetAttributeValue(obj, path[0])
		if err != nil {
			return err
		}
		return traceSetValue(objVal, path[1:], newValue)
	} else {
		return fmt.Errorf("no attribute path specified")
	}
}
