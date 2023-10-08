package di

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

var creatorMutex = sync.Mutex{}
var instanceMutex = sync.Mutex{}
var creators = make(map[string]func() any)
var instances = make(map[string]any)

// Injectable marks a constructor Function of a Struct for DI
func Injectable[T any](creator func() *T) {
	creatorMutex.Lock()
	defer creatorMutex.Unlock()
	creators[getSelector[T]()] = func() any {
		return creator()
	}
}

// Inject gets or create a Instance of the Struct used the Injectable constructor Function
func Inject[T any](identifier ...string) *T {
	selector := getSelector[T]()
	instanceSelector := getSelector[T](identifier...)
	_, instanceExists := instances[instanceSelector].(*T)
	if !instanceExists {
		creator, creatorExists := creators[selector]
		if !creatorExists {
			return nil
		}
		createdInstance, instanceCreated := creator().(*T)
		if instanceCreated {
			instanceMutex.Lock()
			defer instanceMutex.Unlock()
			instance, instanceExists := instances[instanceSelector].(*T)
			if instanceExists {
				return instance
			}
			instances[instanceSelector] = createdInstance
		}
	}
	return instances[instanceSelector].(*T)
}

func getSelector[T any](identifier ...string) string {
	var def T
	typeName := reflect.TypeOf(def).String()
	additionalKey := strings.Join(identifier, "_")
	return fmt.Sprintf("%s_%s", typeName, additionalKey)
}
