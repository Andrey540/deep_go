package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}
type UserRepository struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	constructors map[string]interface{}
	singletons   map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		constructors: make(map[string]interface{}),
		singletons:   make(map[string]interface{}),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	c.constructors[name] = constructor
}

func (c *Container) RegisterSingletonType(name string, constructor interface{}) {
	c.constructors[name] = constructor
	c.singletons[name] = nil
}

func (c *Container) Resolve(name string) (interface{}, error) {
	singleton, ok1 := c.singletons[name]
	if ok1 && singleton != nil {
		return singleton, nil
	}
	constructor, ok2 := c.constructors[name]
	if !ok2 {
		return nil, fmt.Errorf("type: %s is not registered", name)
	}
	switch constructor.(type) {
	case func() any:
		res := constructor.(func() any)()
		if ok1 {
			c.singletons[name] = res
		}
		return res, nil
	default:
		return nil, fmt.Errorf("constructor: %s is not a function", name)
	}
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})
	container.RegisterSingletonType("UserRepository", func() interface{} {
		return &UserRepository{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)

	userRepository1, err := container.Resolve("UserRepository")
	assert.NoError(t, err)
	userRepository2, err := container.Resolve("UserRepository")
	assert.NoError(t, err)

	r1 := userRepository1.(*UserRepository)
	r2 := userRepository2.(*UserRepository)
	assert.True(t, r1 == r2)
}
