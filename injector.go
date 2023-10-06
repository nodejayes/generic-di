package di

import (
	"reflect"
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
func Inject[T any]() *T {
	selector := getSelector[T]()
	_, instanceExists := instances[selector].(*T)
	if !instanceExists {
		creator, creatorExists := creators[selector]
		if !creatorExists {
			return nil
		}
		createdInstance, instanceCreated := creator().(*T)
		if instanceCreated {
			instanceMutex.Lock()
			defer instanceMutex.Unlock()
			instance, instanceExists := instances[selector].(*T)
			if instanceExists {
				return instance
			}
			instances[selector] = createdInstance
		}
	}
	return instances[selector].(*T)
}

func getSelector[T any]() string {
	var def T
	return reflect.TypeOf(def).String()
}
