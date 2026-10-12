package app

import (
	"fmt"
	"reflect"
)

type Container struct {
	services map[reflect.Type]any
}

func NewContainer() *Container {
	return &Container{
		services: make(map[reflect.Type]any),
	}
}

type serviceProvider[T any] func() T

type containerService[T any] struct {
	instance T
	isMade   bool
	provider serviceProvider[T]
}

func Register[T any](c *Container, provider serviceProvider[T]) {
	tp := reflect.TypeFor[T]()
	c.services[tp] = containerService[T]{provider: provider}
}

func Resolve[T any](c *Container) T {
	tp := reflect.TypeFor[T]()

	service, ok := c.services[tp].(containerService[T])
	if !ok {
		panic(fmt.Sprintf("failed to get '%s' service provider", tp.String()))
	}

	if service.isMade {
		return service.instance
	}

	service.instance = service.provider()
	service.isMade = true
	return service.instance
}
