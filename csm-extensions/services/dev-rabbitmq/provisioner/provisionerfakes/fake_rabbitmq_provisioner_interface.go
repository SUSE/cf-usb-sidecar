// This file was generated by counterfeiter
package provisionerfakes

import (
	"sync"

	go_dockerclient "github.com/fsouza/go-dockerclient"
	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-rabbitmq/provisioner"
)

type FakeRabbitmqProvisionerInterface struct {
	CreateContainerStub        func(string) error
	createContainerMutex       sync.RWMutex
	createContainerArgsForCall []struct {
		arg1 string
	}
	createContainerReturns struct {
		result1 error
	}
	DeleteContainerStub        func(string) error
	deleteContainerMutex       sync.RWMutex
	deleteContainerArgsForCall []struct {
		arg1 string
	}
	deleteContainerReturns struct {
		result1 error
	}
	ContainerExistsStub        func(string) (bool, error)
	containerExistsMutex       sync.RWMutex
	containerExistsArgsForCall []struct {
		arg1 string
	}
	containerExistsReturns struct {
		result1 bool
		result2 error
	}
	CreateUserStub        func(string, string, string) (map[string]string, error)
	createUserMutex       sync.RWMutex
	createUserArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	createUserReturns struct {
		result1 map[string]string
		result2 error
	}
	DeleteUserStub        func(string, string) error
	deleteUserMutex       sync.RWMutex
	deleteUserArgsForCall []struct {
		arg1 string
		arg2 string
	}
	deleteUserReturns struct {
		result1 error
	}
	UserExistsStub        func(string, string) (bool, error)
	userExistsMutex       sync.RWMutex
	userExistsArgsForCall []struct {
		arg1 string
		arg2 string
	}
	userExistsReturns struct {
		result1 bool
		result2 error
	}
	FindImageStub        func(string) (*go_dockerclient.Image, error)
	findImageMutex       sync.RWMutex
	findImageArgsForCall []struct {
		arg1 string
	}
	findImageReturns struct {
		result1 *go_dockerclient.Image
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRabbitmqProvisionerInterface) CreateContainer(arg1 string) error {
	fake.createContainerMutex.Lock()
	fake.createContainerArgsForCall = append(fake.createContainerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("CreateContainer", []interface{}{arg1})
	fake.createContainerMutex.Unlock()
	if fake.CreateContainerStub != nil {
		return fake.CreateContainerStub(arg1)
	} else {
		return fake.createContainerReturns.result1
	}
}

func (fake *FakeRabbitmqProvisionerInterface) CreateContainerCallCount() int {
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	return len(fake.createContainerArgsForCall)
}

func (fake *FakeRabbitmqProvisionerInterface) CreateContainerArgsForCall(i int) string {
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	return fake.createContainerArgsForCall[i].arg1
}

func (fake *FakeRabbitmqProvisionerInterface) CreateContainerReturns(result1 error) {
	fake.CreateContainerStub = nil
	fake.createContainerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitmqProvisionerInterface) DeleteContainer(arg1 string) error {
	fake.deleteContainerMutex.Lock()
	fake.deleteContainerArgsForCall = append(fake.deleteContainerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("DeleteContainer", []interface{}{arg1})
	fake.deleteContainerMutex.Unlock()
	if fake.DeleteContainerStub != nil {
		return fake.DeleteContainerStub(arg1)
	} else {
		return fake.deleteContainerReturns.result1
	}
}

func (fake *FakeRabbitmqProvisionerInterface) DeleteContainerCallCount() int {
	fake.deleteContainerMutex.RLock()
	defer fake.deleteContainerMutex.RUnlock()
	return len(fake.deleteContainerArgsForCall)
}

func (fake *FakeRabbitmqProvisionerInterface) DeleteContainerArgsForCall(i int) string {
	fake.deleteContainerMutex.RLock()
	defer fake.deleteContainerMutex.RUnlock()
	return fake.deleteContainerArgsForCall[i].arg1
}

func (fake *FakeRabbitmqProvisionerInterface) DeleteContainerReturns(result1 error) {
	fake.DeleteContainerStub = nil
	fake.deleteContainerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitmqProvisionerInterface) ContainerExists(arg1 string) (bool, error) {
	fake.containerExistsMutex.Lock()
	fake.containerExistsArgsForCall = append(fake.containerExistsArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("ContainerExists", []interface{}{arg1})
	fake.containerExistsMutex.Unlock()
	if fake.ContainerExistsStub != nil {
		return fake.ContainerExistsStub(arg1)
	} else {
		return fake.containerExistsReturns.result1, fake.containerExistsReturns.result2
	}
}

func (fake *FakeRabbitmqProvisionerInterface) ContainerExistsCallCount() int {
	fake.containerExistsMutex.RLock()
	defer fake.containerExistsMutex.RUnlock()
	return len(fake.containerExistsArgsForCall)
}

func (fake *FakeRabbitmqProvisionerInterface) ContainerExistsArgsForCall(i int) string {
	fake.containerExistsMutex.RLock()
	defer fake.containerExistsMutex.RUnlock()
	return fake.containerExistsArgsForCall[i].arg1
}

func (fake *FakeRabbitmqProvisionerInterface) ContainerExistsReturns(result1 bool, result2 error) {
	fake.ContainerExistsStub = nil
	fake.containerExistsReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeRabbitmqProvisionerInterface) CreateUser(arg1 string, arg2 string, arg3 string) (map[string]string, error) {
	fake.createUserMutex.Lock()
	fake.createUserArgsForCall = append(fake.createUserArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("CreateUser", []interface{}{arg1, arg2, arg3})
	fake.createUserMutex.Unlock()
	if fake.CreateUserStub != nil {
		return fake.CreateUserStub(arg1, arg2, arg3)
	} else {
		return fake.createUserReturns.result1, fake.createUserReturns.result2
	}
}

func (fake *FakeRabbitmqProvisionerInterface) CreateUserCallCount() int {
	fake.createUserMutex.RLock()
	defer fake.createUserMutex.RUnlock()
	return len(fake.createUserArgsForCall)
}

func (fake *FakeRabbitmqProvisionerInterface) CreateUserArgsForCall(i int) (string, string, string) {
	fake.createUserMutex.RLock()
	defer fake.createUserMutex.RUnlock()
	return fake.createUserArgsForCall[i].arg1, fake.createUserArgsForCall[i].arg2, fake.createUserArgsForCall[i].arg3
}

func (fake *FakeRabbitmqProvisionerInterface) CreateUserReturns(result1 map[string]string, result2 error) {
	fake.CreateUserStub = nil
	fake.createUserReturns = struct {
		result1 map[string]string
		result2 error
	}{result1, result2}
}

func (fake *FakeRabbitmqProvisionerInterface) DeleteUser(arg1 string, arg2 string) error {
	fake.deleteUserMutex.Lock()
	fake.deleteUserArgsForCall = append(fake.deleteUserArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("DeleteUser", []interface{}{arg1, arg2})
	fake.deleteUserMutex.Unlock()
	if fake.DeleteUserStub != nil {
		return fake.DeleteUserStub(arg1, arg2)
	} else {
		return fake.deleteUserReturns.result1
	}
}

func (fake *FakeRabbitmqProvisionerInterface) DeleteUserCallCount() int {
	fake.deleteUserMutex.RLock()
	defer fake.deleteUserMutex.RUnlock()
	return len(fake.deleteUserArgsForCall)
}

func (fake *FakeRabbitmqProvisionerInterface) DeleteUserArgsForCall(i int) (string, string) {
	fake.deleteUserMutex.RLock()
	defer fake.deleteUserMutex.RUnlock()
	return fake.deleteUserArgsForCall[i].arg1, fake.deleteUserArgsForCall[i].arg2
}

func (fake *FakeRabbitmqProvisionerInterface) DeleteUserReturns(result1 error) {
	fake.DeleteUserStub = nil
	fake.deleteUserReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRabbitmqProvisionerInterface) UserExists(arg1 string, arg2 string) (bool, error) {
	fake.userExistsMutex.Lock()
	fake.userExistsArgsForCall = append(fake.userExistsArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("UserExists", []interface{}{arg1, arg2})
	fake.userExistsMutex.Unlock()
	if fake.UserExistsStub != nil {
		return fake.UserExistsStub(arg1, arg2)
	} else {
		return fake.userExistsReturns.result1, fake.userExistsReturns.result2
	}
}

func (fake *FakeRabbitmqProvisionerInterface) UserExistsCallCount() int {
	fake.userExistsMutex.RLock()
	defer fake.userExistsMutex.RUnlock()
	return len(fake.userExistsArgsForCall)
}

func (fake *FakeRabbitmqProvisionerInterface) UserExistsArgsForCall(i int) (string, string) {
	fake.userExistsMutex.RLock()
	defer fake.userExistsMutex.RUnlock()
	return fake.userExistsArgsForCall[i].arg1, fake.userExistsArgsForCall[i].arg2
}

func (fake *FakeRabbitmqProvisionerInterface) UserExistsReturns(result1 bool, result2 error) {
	fake.UserExistsStub = nil
	fake.userExistsReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeRabbitmqProvisionerInterface) FindImage(arg1 string) (*go_dockerclient.Image, error) {
	fake.findImageMutex.Lock()
	fake.findImageArgsForCall = append(fake.findImageArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("FindImage", []interface{}{arg1})
	fake.findImageMutex.Unlock()
	if fake.FindImageStub != nil {
		return fake.FindImageStub(arg1)
	} else {
		return fake.findImageReturns.result1, fake.findImageReturns.result2
	}
}

func (fake *FakeRabbitmqProvisionerInterface) FindImageCallCount() int {
	fake.findImageMutex.RLock()
	defer fake.findImageMutex.RUnlock()
	return len(fake.findImageArgsForCall)
}

func (fake *FakeRabbitmqProvisionerInterface) FindImageArgsForCall(i int) string {
	fake.findImageMutex.RLock()
	defer fake.findImageMutex.RUnlock()
	return fake.findImageArgsForCall[i].arg1
}

func (fake *FakeRabbitmqProvisionerInterface) FindImageReturns(result1 *go_dockerclient.Image, result2 error) {
	fake.FindImageStub = nil
	fake.findImageReturns = struct {
		result1 *go_dockerclient.Image
		result2 error
	}{result1, result2}
}

func (fake *FakeRabbitmqProvisionerInterface) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	fake.deleteContainerMutex.RLock()
	defer fake.deleteContainerMutex.RUnlock()
	fake.containerExistsMutex.RLock()
	defer fake.containerExistsMutex.RUnlock()
	fake.createUserMutex.RLock()
	defer fake.createUserMutex.RUnlock()
	fake.deleteUserMutex.RLock()
	defer fake.deleteUserMutex.RUnlock()
	fake.userExistsMutex.RLock()
	defer fake.userExistsMutex.RUnlock()
	fake.findImageMutex.RLock()
	defer fake.findImageMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeRabbitmqProvisionerInterface) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ provisioner.RabbitmqProvisionerInterface = new(FakeRabbitmqProvisionerInterface)
