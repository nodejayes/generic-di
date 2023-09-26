package di

import "reflect"

var creators = make(map[string]func() any)
var instances = make(map[string]any)

// Injectable marks a constructor Function of a Struct for DI
func Injectable[T any](creator func() *T) {
	creators[getSelector[T]()] = func() any {
		return creator()
	}
}

// Inject gets or create a Instance of the Struct used the Injectable constructor Function
func Inject[T any]() *T {
	selector := getSelector[T]()
	instance, instanceExists := instances[selector].(*T)
	if !instanceExists {
		creator, creatorExists := creators[selector]
		if !creatorExists {
			return nil
		}
		instance, instanceExists = creator().(*T)
	}
	if !instanceExists {
		return nil
	}
	return instance
}

func getSelector[T any]() string {
	var def T
	return reflect.TypeOf(def).String()
}
