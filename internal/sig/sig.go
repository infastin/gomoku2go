package sig

import (
	"reflect"

	"github.com/imkira/go-observer"
)

type Signal interface {
	Connect(interface{})
	Emit(...interface{})
	Close()
}

type signal struct {
	observer.Property
}

type close struct{}
type update struct{}

func New() Signal {
	return &signal{observer.NewProperty(nil)}
}

func (s *signal) Emit(values ...interface{}) {
	s.Update(values)
}

func (s *signal) Close() {
	s.Update(close{})
}

func (s *signal) Connect(handler interface{}) {
	go func(prop observer.Property) {
		stream := prop.Observe()

		for {
			switch value := stream.WaitNext().(type) {
			case []interface{}:
				fptr := reflect.ValueOf(handler)

				var in []reflect.Value
				for _, arg := range value {
					in = append(in, reflect.ValueOf(arg))
				}

				fptr.Call(in)
			case close:
				return
			}
		}
	}(s.Property)
}
