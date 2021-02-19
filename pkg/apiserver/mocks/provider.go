// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	api "github.com/scoir/canis/pkg/didcomm/doorman/api/protogen"
	apiprotogen "github.com/scoir/canis/pkg/didcomm/loadbalancer/api/protogen"

	datastore "github.com/scoir/canis/pkg/datastore"

	engine "github.com/scoir/canis/pkg/credential/engine"

	indy "github.com/scoir/canis/pkg/credential/engine/indy"

	kms "github.com/hyperledger/aries-framework-go/pkg/kms"

	mediatorapiprotogen "github.com/scoir/canis/pkg/didcomm/mediator/api/protogen"

	mock "github.com/stretchr/testify/mock"

	presentproofengine "github.com/scoir/canis/pkg/presentproof/engine"

	protogen "github.com/scoir/canis/pkg/didcomm/issuer/api/protogen"

	verifierapiprotogen "github.com/scoir/canis/pkg/didcomm/verifier/api/protogen"
)

// Provider is an autogenerated mock type for the provider type
type Provider struct {
	mock.Mock
}

// GetCredentialEngineRegistry provides a mock function with given fields:
func (_m *Provider) GetCredentialEngineRegistry() (engine.CredentialRegistry, error) {
	ret := _m.Called()

	var r0 engine.CredentialRegistry
	if rf, ok := ret.Get(0).(func() engine.CredentialRegistry); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(engine.CredentialRegistry)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDoormanClient provides a mock function with given fields:
func (_m *Provider) GetDoormanClient() (api.DoormanClient, error) {
	ret := _m.Called()

	var r0 api.DoormanClient
	if rf, ok := ret.Get(0).(func() api.DoormanClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.DoormanClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetIssuerClient provides a mock function with given fields:
func (_m *Provider) GetIssuerClient() (protogen.IssuerClient, error) {
	ret := _m.Called()

	var r0 protogen.IssuerClient
	if rf, ok := ret.Get(0).(func() protogen.IssuerClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(protogen.IssuerClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLoadbalancerClient provides a mock function with given fields:
func (_m *Provider) GetLoadbalancerClient() (apiprotogen.LoadbalancerClient, error) {
	ret := _m.Called()

	var r0 apiprotogen.LoadbalancerClient
	if rf, ok := ret.Get(0).(func() apiprotogen.LoadbalancerClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(apiprotogen.LoadbalancerClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMediatorClient provides a mock function with given fields:
func (_m *Provider) GetMediatorClient() (mediatorapiprotogen.MediatorClient, error) {
	ret := _m.Called()

	var r0 mediatorapiprotogen.MediatorClient
	if rf, ok := ret.Get(0).(func() mediatorapiprotogen.MediatorClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mediatorapiprotogen.MediatorClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPresentationEngineRegistry provides a mock function with given fields:
func (_m *Provider) GetPresentationEngineRegistry() (presentproofengine.PresentationRegistry, error) {
	ret := _m.Called()

	var r0 presentproofengine.PresentationRegistry
	if rf, ok := ret.Get(0).(func() presentproofengine.PresentationRegistry); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(presentproofengine.PresentationRegistry)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVerifierClient provides a mock function with given fields:
func (_m *Provider) GetVerifierClient() (verifierapiprotogen.VerifierClient, error) {
	ret := _m.Called()

	var r0 verifierapiprotogen.VerifierClient
	if rf, ok := ret.Get(0).(func() verifierapiprotogen.VerifierClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(verifierapiprotogen.VerifierClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IndyVDR provides a mock function with given fields:
func (_m *Provider) IndyVDR() (indy.VDRClient, error) {
	ret := _m.Called()

	var r0 indy.VDRClient
	if rf, ok := ret.Get(0).(func() indy.VDRClient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(indy.VDRClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// KMS provides a mock function with given fields:
func (_m *Provider) KMS() kms.KeyManager {
	ret := _m.Called()

	var r0 kms.KeyManager
	if rf, ok := ret.Get(0).(func() kms.KeyManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(kms.KeyManager)
		}
	}

	return r0
}

// MediatorKMS provides a mock function with given fields:
func (_m *Provider) MediatorKMS() kms.KeyManager {
	ret := _m.Called()

	var r0 kms.KeyManager
	if rf, ok := ret.Get(0).(func() kms.KeyManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(kms.KeyManager)
		}
	}

	return r0
}

// Store provides a mock function with given fields:
func (_m *Provider) Store() datastore.Store {
	ret := _m.Called()

	var r0 datastore.Store
	if rf, ok := ret.Get(0).(func() datastore.Store); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(datastore.Store)
		}
	}

	return r0
}
