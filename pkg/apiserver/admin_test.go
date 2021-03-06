/*
Copyright Scoir Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package apiserver

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hyperledger/indy-vdr/wrappers/golang/identifiers"

	api "github.com/scoir/canis/pkg/apiserver/api/protogen"
	apimocks "github.com/scoir/canis/pkg/apiserver/mocks"
	emocks "github.com/scoir/canis/pkg/credential/engine/mocks"
	"github.com/scoir/canis/pkg/datastore"
	"github.com/scoir/canis/pkg/datastore/mocks"
	dmocks "github.com/scoir/canis/pkg/didexchange/mocks"
	"github.com/scoir/canis/pkg/protogen/common"
)

type AdminTestSuite struct {
	Store             *mocks.Store
	Bouncer           *dmocks.Bouncer
	CredRegistry      *emocks.CredentialRegistry
	LoadbalanceClient *apimocks.MockLoadbalancer
	IndyClient        *apimocks.MockVDRClient
	KMS               *apimocks.MockKMS
	MediatorKMS       *apimocks.MockKMS
	Doorman           *apimocks.MockDoorman
	Issuer            *apimocks.MockIssuer
	Verifier          *apimocks.MockVerifier
	Mediator          *apimocks.MockMediator
}

func SetupTest() (*APIServer, *AdminTestSuite) {
	suite := &AdminTestSuite{}
	suite.Store = &mocks.Store{}
	suite.Bouncer = &dmocks.Bouncer{}
	suite.CredRegistry = &emocks.CredentialRegistry{}
	suite.IndyClient = &apimocks.MockVDRClient{}
	suite.KMS = &apimocks.MockKMS{}
	suite.MediatorKMS = &apimocks.MockKMS{}
	suite.Doorman = &apimocks.MockDoorman{}
	suite.LoadbalanceClient = &apimocks.MockLoadbalancer{
		EndpointValue: "0.0.0.0:420",
	}
	suite.Issuer = &apimocks.MockIssuer{}
	suite.Verifier = &apimocks.MockVerifier{}
	suite.Mediator = &apimocks.MockMediator{}

	target := &APIServer{
		keyMgr:         suite.KMS,
		mediatorKeyMgr: suite.KMS,
		agentStore:     suite.Store,
		schemaStore:    suite.Store,
		store:          suite.Store,
		client:         suite.IndyClient,
		schemaRegistry: suite.CredRegistry,
		doorman:        suite.Doorman,
		issuer:         suite.Issuer,
		verifier:       suite.Verifier,
		loadbalancer:   suite.LoadbalanceClient,
		mediator:       suite.Mediator,
	}

	return target, suite
}

func TestCreateAgent(t *testing.T) {

	t.Run("no public did", func(t *testing.T) {
		target, suite := SetupTest()
		request := &api.CreateAgentRequest{
			Agent: &api.NewAgent{
				Name:                  "Test Agent",
				EndorsableSchemaNames: nil,
			},
		}

		match := func(m *datastore.Agent) bool {
			return m.Name == "Test Agent"
		}

		suite.Store.On("GetAgent", "Test Agent").Return(nil, errors.New("not found"))
		suite.Store.On("InsertAgent", mock.MatchedBy(match)).Return("123", nil)

		resp, err := target.CreateAgent(context.Background(), request)
		require.Nil(t, err)
		require.Equal(t, resp.Id, "123")
	})
	t.Run("with public did", func(t *testing.T) {
		target, suite := SetupTest()

		request := &api.CreateAgentRequest{
			Agent: &api.NewAgent{
				Name:                  "Test Agent",
				EndorsableSchemaNames: []string{"test-schema-id"},
				PublicDid:             true,
			},
		}

		d, err := identifiers.CreateDID(&identifiers.MyDIDInfo{
			PublicKey:  []byte("abcdefghijklmnopqrs"),
			Cid:        true,
			MethodName: "scr",
		})
		require.NoError(t, err)

		did := &datastore.DID{
			DID:     d,
			KeyPair: &datastore.KeyPair{},
		}
		s := &datastore.Schema{}
		match := func(m *datastore.Agent) bool {
			return m.Name == "Test Agent" && m.HasPublicDID
		}

		suite.Store.On("GetAgent", "Test Agent").Return(nil, errors.New("not found"))
		suite.Store.On("InsertAgent", mock.MatchedBy(match)).Return("123", nil)
		suite.Store.On("GetPublicDID").Return(did, nil)
		suite.Store.On("GetSchema", "test-schema-id").Return(s, nil)
		suite.CredRegistry.On("RegisterSchema", mock.AnythingOfType("*datastore.DID"), s).Return(nil)

		_, err = target.CreateAgent(context.Background(), request)
		require.Nil(t, err)
		//require.Equal(t, resp.Id, "123")
	})
}

func TestCreateAgentFails(t *testing.T) {
	target, suite := SetupTest()
	request := &api.CreateAgentRequest{
		Agent: &api.NewAgent{
			Name:                  "Test Agent",
			EndorsableSchemaNames: nil,
		},
	}
	match := func(m *datastore.Agent) bool {
		return m.Name == "Test Agent"
	}

	suite.Store.On("GetAgent", "Test Agent").Return(nil, errors.New("not found"))
	suite.Store.On("InsertAgent", mock.MatchedBy(match)).Return("", errors.New("Boom"))

	resp, err := target.CreateAgent(context.Background(), request)
	require.Nil(t, resp)
	require.NotNil(t, err)
	require.Equal(t, "rpc error: code = Internal desc = failed to create agent Test Agent: Boom", err.Error())
}

func TestCreateAgentMissingRequiredField(t *testing.T) {
	target, _ := SetupTest()
	request := &api.CreateAgentRequest{
		Agent: &api.NewAgent{
			EndorsableSchemaNames: nil,
		},
	}

	resp, err := target.CreateAgent(context.Background(), request)
	require.Nil(t, resp)
	require.NotNil(t, err)
	require.Equal(t, "rpc error: code = InvalidArgument desc = name is a required field", err.Error())
}

func TestCreateAgentAlreadyExists(t *testing.T) {
	target, suite := SetupTest()
	request := &api.CreateAgentRequest{
		Agent: &api.NewAgent{
			Name:                  "Test Agent",
			EndorsableSchemaNames: nil,
		},
	}

	suite.Store.On("GetAgent", "Test Agent").Return(nil, nil)

	resp, err := target.CreateAgent(context.Background(), request)
	require.Nil(t, resp)
	require.NotNil(t, err)
	require.Equal(t, "rpc error: code = AlreadyExists desc = agent with name Test Agent already exists", err.Error())
}

func TestGetAgent(t *testing.T) {
	target, suite := SetupTest()
	request := &api.GetAgentRequest{
		Id: "123",
	}

	suite.Store.On("GetAgent", "123").Return(&datastore.Agent{ID: "123", Name: "test Agent"}, nil)

	resp, err := target.GetAgent(context.Background(), request)
	require.Nil(t, err)
	require.Equal(t, "test Agent", resp.Agent.Name)
}

func TestGetAgentErr(t *testing.T) {
	target, suite := SetupTest()
	request := &api.GetAgentRequest{
		Id: "123",
	}

	suite.Store.On("GetAgent", "123").Return(nil, errors.New("BOOM"))

	resp, err := target.GetAgent(context.Background(), request)
	require.Nil(t, resp)
	require.NotNil(t, err)
	require.Equal(t, "rpc error: code = Internal desc = unable to get agent: BOOM", err.Error())
}

func TestListAgent(t *testing.T) {
	target, suite := SetupTest()
	request := &api.ListAgentRequest{}

	suite.Store.On("ListAgent", &datastore.AgentCriteria{}).Return(&datastore.AgentList{
		Count:  1,
		Agents: []*datastore.Agent{{ID: "123", Name: "test agent"}},
	}, nil)

	resp, err := target.ListAgent(context.Background(), request)
	require.Nil(t, err)
	require.Equal(t, "test agent", resp.Agents[0].Name)
}

func TestListAgentErr(t *testing.T) {
	target, suite := SetupTest()
	request := &api.ListAgentRequest{}

	suite.Store.On("ListAgent", &datastore.AgentCriteria{}).Return(nil, errors.New("BOOM"))

	resp, err := target.ListAgent(context.Background(), request)
	require.Nil(t, resp)
	require.NotNil(t, err)
	require.Equal(t, "rpc error: code = Internal desc = unable to list agent: BOOM", err.Error())
}

func TestDeleteAgent(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()
		request := &api.DeleteAgentRequest{
			Id: "123",
		}
		agent := &datastore.Agent{ID: "123"}
		suite.Store.On("GetAgent", "123").Return(agent, nil)
		suite.Store.On("DeleteAgent", "123").Return(nil)

		resp, err := target.DeleteAgent(context.Background(), request)
		require.Nil(t, err)
		require.NotNil(t, resp)

	})
	t.Run("store get error", func(t *testing.T) {
		target, suite := SetupTest()
		request := &api.DeleteAgentRequest{
			Id: "123",
		}

		suite.Store.On("GetAgent", "123").Return(nil, errors.New("BOOM"))

		resp, err := target.DeleteAgent(context.Background(), request)
		require.Error(t, err)
		require.Nil(t, resp)

	})

	t.Run("store save error", func(t *testing.T) {
		target, suite := SetupTest()
		request := &api.DeleteAgentRequest{
			Id: "123",
		}

		agent := &datastore.Agent{ID: "123"}
		suite.Store.On("GetAgent", "123").Return(agent, nil)
		suite.Store.On("DeleteAgent", "123").Return(errors.New("BOOM"))

		resp, err := target.DeleteAgent(context.Background(), request)
		require.Nil(t, resp)
		require.NotNil(t, err)
		require.Equal(t, "rpc error: code = Internal desc = failed to delete agent 123: BOOM", err.Error())
	})
}

func TestUpdateAgent(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()
		request := &api.UpdateAgentRequest{
			Agent: &api.Agent{
				Id:                    "123",
				Name:                  "Test Agent",
				EndorsableSchemaNames: []string{"test-schema-id"},
				PublicDid:             true,
			},
		}

		a := &datastore.Agent{
			ID:                    "123",
			Name:                  "Test Agent",
			EndorsableSchemaNames: []string{"test-schema-id"},
			HasPublicDID:          true,
		}
		d, err := identifiers.CreateDID(&identifiers.MyDIDInfo{
			PublicKey:  []byte("abcdefghijklmnopqrs"),
			Cid:        true,
			MethodName: "scr",
		})
		require.NoError(t, err)
		did := &datastore.DID{
			DID:     d,
			KeyPair: &datastore.KeyPair{},
		}
		s := &datastore.Schema{}

		match := func(m *datastore.Agent) bool {
			return m.Name == "Test Agent" && m.HasPublicDID
		}

		suite.Store.On("GetAgent", "123").Return(a, nil)
		suite.Store.On("GetPublicDID").Return(did, nil)
		suite.Store.On("GetSchema", "test-schema-id").Return(s, nil)
		suite.CredRegistry.On("RegisterSchema", mock.AnythingOfType("*datastore.DID"), s).Return(nil)
		suite.Store.On("UpdateAgent", mock.MatchedBy(match)).Return(nil)

		resp, err := target.UpdateAgent(context.Background(), request)
		require.Nil(t, err)
		require.NotNil(t, resp)
	})
	t.Run("update agent error", func(t *testing.T) {
		target, suite := SetupTest()
		request := &api.UpdateAgentRequest{
			Agent: &api.Agent{
				Id:                    "123",
				Name:                  "Test Agent",
				EndorsableSchemaNames: nil,
			},
		}

		a := &datastore.Agent{
			ID:                    "123",
			Name:                  "Test Agent",
			EndorsableSchemaNames: nil,
		}

		suite.Store.On("GetAgent", "123").Return(a, nil)
		suite.Store.On("UpdateAgent", a).Return(errors.New("BOOM"))

		resp, err := target.UpdateAgent(context.Background(), request)
		require.Error(t, err)
		require.Nil(t, resp)
	})
	t.Run("get agent error", func(t *testing.T) {
		target, suite := SetupTest()
		request := &api.UpdateAgentRequest{
			Agent: &api.Agent{
				Id:                    "123",
				Name:                  "Test Agent",
				EndorsableSchemaNames: nil,
			},
		}

		suite.Store.On("GetAgent", "123").Return(nil, errors.New("BOOM"))

		resp, err := target.UpdateAgent(context.Background(), request)
		require.Error(t, err)
		require.Nil(t, resp)
	})
	t.Run("empty agent ID", func(t *testing.T) {
		target, _ := SetupTest()
		request := &api.UpdateAgentRequest{
			Agent: &api.Agent{
				Id: "",
			},
		}

		resp, err := target.UpdateAgent(context.Background(), request)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestCreateSchema(t *testing.T) {
	target, suite := SetupTest()
	request := &api.CreateSchemaRequest{
		Schema: &api.NewSchema{
			Name:    "Test Schema",
			Version: "0.0.1",
			Attributes: []*api.Attribute{{
				Name: "City",
				Type: api.Attribute_STRING,
			}},
		},
	}

	match := func(m *datastore.Schema) bool {
		return m.Name == "Test Schema" &&
			m.Version == "0.0.1" &&
			1 == len(m.Attributes)
	}

	suite.Store.On("GetSchema", "Test Schema").Return(nil, errors.New("not found"))
	suite.CredRegistry.On("CreateSchema", mock.MatchedBy(match)).Return("abc", nil)
	suite.Store.On("InsertSchema", mock.MatchedBy(match)).Return("123", nil)

	resp, err := target.CreateSchema(context.Background(), request)
	require.Nil(t, err)
	require.Equal(t, resp.Id, "123")
}

func TestCreateSchemaFails(t *testing.T) {
	target, suite := SetupTest()
	request := &api.CreateSchemaRequest{
		Schema: &api.NewSchema{
			Name: "Test Schema",
		},
	}

	match := func(m *datastore.Schema) bool {
		return m.Name == "Test Schema"
	}

	suite.Store.On("GetSchema", "Test Schema").Return(nil, errors.New("not found"))
	suite.CredRegistry.On("CreateSchema", mock.MatchedBy(match)).Return("abc", nil)
	suite.Store.On("InsertSchema", mock.MatchedBy(match)).Return("", errors.New("Boom"))

	resp, err := target.CreateSchema(context.Background(), request)
	require.Nil(t, resp)
	require.NotNil(t, err)
	require.Equal(t, "rpc error: code = Internal desc = failed to create schema Test Schema: Boom", err.Error())
}

func TestCreateSchemaMissingRequiredField(t *testing.T) {
	target, _ := SetupTest()
	request := &api.CreateSchemaRequest{
		Schema: &api.NewSchema{
			Name: "",
		},
	}

	resp, err := target.CreateSchema(context.Background(), request)
	require.Nil(t, resp)
	require.NotNil(t, err)
	require.Equal(t, "rpc error: code = InvalidArgument desc = name is a required field", err.Error())
}

func TestCreateSchemaAlreadyExists(t *testing.T) {
	target, suite := SetupTest()
	request := &api.CreateSchemaRequest{
		Schema: &api.NewSchema{
			Name: "Test Schema",
		},
	}

	suite.Store.On("GetSchema", "Test Schema").Return(nil, nil)

	resp, err := target.CreateSchema(context.Background(), request)
	require.Nil(t, resp)
	require.NotNil(t, err)
	require.Equal(t, "rpc error: code = AlreadyExists desc = schema with id Test Schema already exists", err.Error())
}

func TestGetSchema(t *testing.T) {
	target, suite := SetupTest()
	request := &api.GetSchemaRequest{
		Id: "123",
	}

	suite.Store.On("GetSchema", "123").Return(&datastore.Schema{
		ID:   "123",
		Name: "test schema",
		Attributes: []*datastore.Attribute{
			{
				Name: "City",
				Type: 1,
			},
		},
	}, nil)

	resp, err := target.GetSchema(context.Background(), request)
	require.Nil(t, err)
	require.Equal(t, "test schema", resp.Schema.Name)
}

func TestGetSchemaErr(t *testing.T) {
	target, suite := SetupTest()
	request := &api.GetSchemaRequest{
		Id: "123",
	}

	suite.Store.On("GetSchema", "123").Return(nil, errors.New("BOOM"))

	resp, err := target.GetSchema(context.Background(), request)
	require.Nil(t, resp)
	require.NotNil(t, err)
	require.Equal(t, "rpc error: code = Internal desc = unable to get schema: BOOM", err.Error())
}

func TestListSchema(t *testing.T) {
	target, suite := SetupTest()
	request := &api.ListSchemaRequest{}

	suite.Store.On("ListSchema", &datastore.SchemaCriteria{}).Return(&datastore.SchemaList{
		Count: 1,
		Schema: []*datastore.Schema{{
			ID:   "123",
			Name: "test schema",
			Attributes: []*datastore.Attribute{
				{
					Name: "City",
					Type: 1,
				},
			},
		}},
	}, nil)

	resp, err := target.ListSchema(context.Background(), request)
	require.Nil(t, err)
	require.Equal(t, "test schema", resp.Schema[0].Name)
}

func TestListSchemaErr(t *testing.T) {
	target, suite := SetupTest()
	request := &api.ListSchemaRequest{}

	suite.Store.On("ListSchema", &datastore.SchemaCriteria{}).Return(nil, errors.New("BOOM"))

	resp, err := target.ListSchema(context.Background(), request)
	require.Nil(t, resp)
	require.NotNil(t, err)
	require.Equal(t, "rpc error: code = Internal desc = unable to list schema: BOOM", err.Error())
}

func TestDeleteSchema(t *testing.T) {
	target, suite := SetupTest()
	request := &api.DeleteSchemaRequest{
		Id: "123",
	}

	suite.Store.On("DeleteSchema", "123").Return(nil)

	resp, err := target.DeleteSchema(context.Background(), request)
	require.Nil(t, err)
	require.NotNil(t, resp)
}

func TestDeleteSchemaErr(t *testing.T) {
	target, suite := SetupTest()
	request := &api.DeleteSchemaRequest{
		Id: "123",
	}

	suite.Store.On("DeleteSchema", "123").Return(errors.New("BOOM"))

	resp, err := target.DeleteSchema(context.Background(), request)
	require.Nil(t, resp)
	require.NotNil(t, err)
	require.Equal(t, "rpc error: code = Internal desc = failed to delete schema 123: BOOM", err.Error())
}

func TestUpdateSchema(t *testing.T) {
	target, suite := SetupTest()
	request := &api.UpdateSchemaRequest{
		Schema: &api.Schema{
			Id:      "123",
			Name:    "Test Schema",
			Version: "0.0.1",
			Attributes: []*api.Attribute{{
				Name: "City",
				Type: api.Attribute_STRING,
			}},
		},
	}

	a := &datastore.Schema{
		ID:      "123",
		Name:    "Test Schema",
		Version: "0.0.1",
		Attributes: []*datastore.Attribute{{
			Name: "City",
			Type: int32(api.Attribute_STRING),
		}},
	}

	suite.Store.On("GetSchema", "123").Return(a, nil)
	suite.Store.On("UpdateSchema", a).Return(nil)

	resp, err := target.UpdateSchema(context.Background(), request)
	require.Nil(t, err)
	require.NotNil(t, resp)
}

func TestUpdateSchemaErr(t *testing.T) {
	target, suite := SetupTest()
	request := &api.UpdateSchemaRequest{
		Schema: &api.Schema{
			Id:      "123",
			Name:    "Test Schema",
			Version: "0.0.1",
			Attributes: []*api.Attribute{{
				Name: "City",
				Type: api.Attribute_STRING,
			}},
		},
	}

	suite.Store.On("GetSchema", "123").Return(nil, errors.New("BOOM"))

	resp, err := target.UpdateSchema(context.Background(), request)
	require.NotNil(t, err)
	require.Nil(t, resp)
}

func TestUpdateSchemaDataMissing(t *testing.T) {
	target, _ := SetupTest()
	request := &api.UpdateSchemaRequest{
		Schema: &api.Schema{
			Id:      "",
			Name:    "Test Schema",
			Version: "0.0.1",
			Attributes: []*api.Attribute{{
				Name: "City",
				Type: api.Attribute_STRING,
			}},
		},
	}

	resp, err := target.UpdateSchema(context.Background(), request)
	require.NotNil(t, err)
	require.Nil(t, resp)
}

func TestGetAgentInvitation(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()
		req := &common.InvitationRequest{}

		resp := &common.InvitationResponse{}
		suite.Doorman.InviteResponse = resp

		result, err := target.GetAgentInvitation(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, result)
	})
	t.Run("doorman error", func(t *testing.T) {
		target, suite := SetupTest()
		req := &common.InvitationRequest{}

		suite.Doorman.InviteErr = errors.New("BOOM")

		result, err := target.GetAgentInvitation(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, result)
	})

}

func TestGetAgentInvitationIamge(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()
		req := &common.InvitationRequest{}

		resp := &common.InvitationResponse{
			Invitation: `{"test": "abc"}`,
		}
		suite.Doorman.InviteResponse = resp

		result, err := target.GetAgentInvitationImage(context.Background(), req)
		require.NoError(t, err)
		require.Equal(t, result.ContentType, "image/png")
		require.Len(t, result.Data, 437)
	})
	t.Run("doorman error", func(t *testing.T) {
		target, suite := SetupTest()
		req := &common.InvitationRequest{}

		suite.Doorman.InviteErr = errors.New("BOOM")

		result, err := target.GetAgentInvitationImage(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, result)
	})

}

func TestSeedPublicDID(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()
		req := &api.SeedPublicDIDRequest{
			Seed: "b2352b32947e188eb72871093ac6217e",
		}

		suite.Store.On("GetPublicDID").Return(nil, errors.New("not found"))
		suite.Store.On("SetPublicDID", mock.AnythingOfType("*datastore.DID")).Return(nil)
		suite.Store.On("GetMediatorDID").Return(&datastore.DID{}, nil)

		resp, err := target.SeedPublicDID(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

	})
	t.Run("happy with mediator DID creation", func(t *testing.T) {
		target, suite := SetupTest()
		req := &api.SeedPublicDIDRequest{
			Seed: "b2352b32947e188eb72871093ac6217e",
		}

		suite.Mediator.EndpointResponse = &common.EndpointResponse{Endpoint: "ws://test:1000"}

		var newPubDID datastore.DID
		match := func(did *datastore.DID) bool {
			newPubDID = *did
			return true
		}

		medMatch := func(did *datastore.DID) bool {
			return did.Endpoint == "ws://test:1000"
		}

		suite.Store.On("GetPublicDID").Return(nil, errors.New("not found")).Once()
		suite.Store.On("SetPublicDID", mock.MatchedBy(match)).Return(nil)
		suite.Store.On("GetMediatorDID").Return(nil, errors.New("not found"))
		suite.Store.On("GetPublicDID").Return(&newPubDID, nil).Once()
		suite.Store.On("SetMediatorDID", mock.MatchedBy(medMatch)).Return(nil)

		resp, err := target.SeedPublicDID(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

	})
	t.Run("happy with random seed", func(t *testing.T) {
		target, suite := SetupTest()
		req := &api.SeedPublicDIDRequest{
			Seed: "",
		}

		suite.Store.On("GetPublicDID").Return(nil, errors.New("not found"))
		suite.Store.On("SetPublicDID", mock.AnythingOfType("*datastore.DID")).Return(nil)
		suite.Store.On("GetMediatorDID").Return(&datastore.DID{}, nil)

		resp, err := target.SeedPublicDID(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

	})
	t.Run("unregistered public DID", func(t *testing.T) {
		target, suite := SetupTest()
		req := &api.SeedPublicDIDRequest{
			Seed: "b2352b32947e188eb72871093ac6217e",
		}

		suite.Store.On("GetPublicDID").Return(nil, errors.New("not found"))
		suite.IndyClient.GetNymErr = errors.New("not found")

		resp, err := target.SeedPublicDID(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)

	})
	t.Run("public did already set", func(t *testing.T) {
		target, suite := SetupTest()
		req := &api.SeedPublicDIDRequest{
			Seed: "abc",
		}

		suite.Store.On("GetPublicDID").Return(nil, nil)

		resp, err := target.SeedPublicDID(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)

	})

	t.Run("unable to import private key", func(t *testing.T) {
		target, suite := SetupTest()
		req := &api.SeedPublicDIDRequest{
			Seed: "b2352b32947e188eb72871093ac6217e",
		}

		suite.Store.On("GetPublicDID").Return(nil, errors.New("not found"))
		suite.KMS.ImportPrivateKeyErr = errors.New("unexpected")

		resp, err := target.SeedPublicDID(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)

	})

	t.Run("save fails", func(t *testing.T) {
		target, suite := SetupTest()
		req := &api.SeedPublicDIDRequest{
			Seed: "b2352b32947e188eb72871093ac6217e",
		}

		suite.Store.On("GetPublicDID").Return(nil, errors.New("not found"))
		suite.Store.On("SetPublicDID", mock.AnythingOfType("*datastore.DID")).Return(errors.New("BOOM"))

		resp, err := target.SeedPublicDID(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)

	})

	t.Run("bad seed", func(t *testing.T) {
		target, suite := SetupTest()
		req := &api.SeedPublicDIDRequest{
			Seed: "abc",
		}

		suite.Store.On("GetPublicDID").Return(nil, errors.New("not found"))

		resp, err := target.SeedPublicDID(context.Background(), req)

		require.Error(t, err)
		require.Nil(t, resp)

	})
}

func TestIssueCredential(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()
		req := &common.IssueCredentialRequest{
			Credential: &common.Credential{
				SchemaId: "test-schema-id",
				Comment:  "test comment",
				Type:     "lds/json-ld-proof",
				Preview: []*common.CredentialAttribute{
					{
						Name:  "test-field",
						Value: "test-value",
					},
				},
			},
		}

		suite.Issuer.IssueCredResponse = &common.IssueCredentialResponse{
			CredentialId: "new-cred-id",
		}

		resp, err := target.IssueCredential(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "new-cred-id", resp.CredentialId)
	})
	t.Run("issuer error", func(t *testing.T) {
		target, suite := SetupTest()
		req := &common.IssueCredentialRequest{
			Credential: &common.Credential{
				SchemaId: "test-schema-id",
				Comment:  "test comment",
				Type:     "lds/json-ld-proof",
				Preview: []*common.CredentialAttribute{
					{
						Name:  "test-field",
						Value: "test-value",
					},
				},
			},
		}

		suite.Issuer.IssueCredErr = errors.New("BOOM")

		resp, err := target.IssueCredential(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestRequestPresentation(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()

		req := &common.RequestPresentationRequest{
			AgentName:  "test-agent-id",
			ExternalId: "test-external-id",
			Presentation: &common.RequestPresentation{
				RequestedAttributes: map[string]*common.AttrInfo{
					"test-attr": {
						Name: "test-attr",
					},
				},
				RequestedPredicates: map[string]*common.PredicateInfo{
					"test-predicate": {
						Name: "test-predicate",
					},
				},
			},
		}

		suite.Verifier.RequestPresResponse = &common.RequestPresentationResponse{
			RequestPresentationId: "test-presentation-id",
		}

		resp, err := target.RequestPresentation(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, "test-presentation-id", resp.RequestPresentationId)
	})
	t.Run("verifier fails", func(t *testing.T) {
		target, suite := SetupTest()

		req := &common.RequestPresentationRequest{
			AgentName:  "test-agent-id",
			ExternalId: "test-external-id",
			Presentation: &common.RequestPresentation{
				RequestedAttributes: map[string]*common.AttrInfo{
					"test-attr": {
						Name: "test-attr",
					},
				},
				RequestedPredicates: map[string]*common.PredicateInfo{
					"test-predicate": {
						Name: "test-predicate",
					},
				},
			},
		}

		suite.Verifier.RequestPresErr = errors.New("BOOM")

		resp, err := target.RequestPresentation(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestCreateWebhook(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()

		req := &api.CreateWebhookRequest{
			Id: "test-webhook",
			Webhook: []*api.Webhook{
				{
					Url: "test-url",
				},
				{
					Url: "test-url2",
				},
			},
		}

		hook := &datastore.Webhook{
			Type: "test-webhook",
			URL:  "test-url",
		}

		hook2 := &datastore.Webhook{
			Type: "test-webhook",
			URL:  "test-url2",
		}

		suite.Store.On("AddWebhook", hook).Return(nil).Once()
		suite.Store.On("AddWebhook", hook2).Return(nil).Once()

		resp, err := target.CreateWebhook(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

	})
	t.Run("store error", func(t *testing.T) {
		target, suite := SetupTest()

		req := &api.CreateWebhookRequest{
			Id: "test-webhook",
			Webhook: []*api.Webhook{
				{
					Url: "test-url",
				},
			},
		}

		hook := &datastore.Webhook{
			Type: "test-webhook",
			URL:  "test-url",
		}

		suite.Store.On("AddWebhook", hook).Return(errors.New("BOOM")).Once()

		resp, err := target.CreateWebhook(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)

	})
}

func TestListWebhook(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()

		req := &api.ListWebhookRequest{
			Id: "test-webhook-id",
		}

		hooks := []*datastore.Webhook{
			{
				Type: "test-webhook-id",
				URL:  "test-url",
			},
			{
				Type: "test-webhook-id2",
				URL:  "test-url2",
			},
		}
		suite.Store.On("ListWebhooks", "test-webhook-id").Return(hooks, nil)

		resp, err := target.ListWebhook(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

		require.Equal(t, []*api.Webhook{
			{
				Url: "test-url",
			},
			{
				Url: "test-url2",
			},
		}, resp.Hooks)

	})
	t.Run("store error", func(t *testing.T) {
		target, suite := SetupTest()

		req := &api.ListWebhookRequest{
			Id: "test-webhook-id",
		}

		suite.Store.On("ListWebhooks", "test-webhook-id").Return(nil, errors.New("BOOM"))

		resp, err := target.ListWebhook(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)

	})
}

func TestDeleteWebhook(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()

		req := &api.DeleteWebhookRequest{
			Id: "test-webhook-id",
		}

		suite.Store.On("DeleteWebhook", "test-webhook-id").Return(nil)

		resp, err := target.DeleteWebhook(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

	})
	t.Run("store error", func(t *testing.T) {
		target, suite := SetupTest()

		req := &api.DeleteWebhookRequest{
			Id: "test-webhook-id",
		}

		suite.Store.On("DeleteWebhook", "test-webhook-id").Return(errors.New("BOOM"))

		resp, err := target.DeleteWebhook(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)

	})
}

func TestListConnections(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()

		req := &api.ListConnectionRequest{AgentName: "agent-1"}
		agent := &datastore.Agent{}
		connections := []*datastore.AgentConnection{
			{
				TheirLabel: "their label",
			},
		}

		suite.Store.On("GetAgent", "agent-1").Return(agent, nil)
		suite.Store.On("ListAgentConnections", agent).Return(connections, nil)

		resp, err := target.ListConnections(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

	})
	t.Run("agent connection errors", func(t *testing.T) {
		target, suite := SetupTest()

		req := &api.ListConnectionRequest{AgentName: "agent-1"}
		agent := &datastore.Agent{}

		suite.Store.On("GetAgent", "agent-1").Return(agent, nil)
		suite.Store.On("ListAgentConnections", agent).Return(nil, errors.New("boom"))

		resp, err := target.ListConnections(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)

	})
	t.Run("agent not found", func(t *testing.T) {
		target, suite := SetupTest()

		req := &api.ListConnectionRequest{AgentName: "agent-1"}

		suite.Store.On("GetAgent", "agent-1").Return(nil, errors.New("boom"))

		resp, err := target.ListConnections(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)

	})
}

func TestDeleteConnections(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()

		req := &api.DeleteConnectionRequest{AgentName: "agent-1", ExternalId: "external-id"}
		agent := &datastore.Agent{}

		suite.Store.On("GetAgent", "agent-1").Return(agent, nil)
		suite.Store.On("DeleteAgentConnection", agent, "external-id").Return(nil)

		resp, err := target.DeleteConnection(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)

	})
	t.Run("agent connection errors", func(t *testing.T) {
		target, suite := SetupTest()

		req := &api.DeleteConnectionRequest{AgentName: "agent-1", ExternalId: "external-id"}
		suite.Store.On("GetAgent", "agent-1").Return(nil, errors.New("not found"))

		resp, err := target.DeleteConnection(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)
	})
	t.Run("agent not found", func(t *testing.T) {
		target, suite := SetupTest()

		req := &api.DeleteConnectionRequest{AgentName: "agent-1", ExternalId: "external-id"}
		agent := &datastore.Agent{}

		suite.Store.On("GetAgent", "agent-1").Return(agent, nil)
		suite.Store.On("DeleteAgentConnection", agent, "external-id").Return(errors.New("boom"))

		resp, err := target.DeleteConnection(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestAcceptInvitation(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()

		req := &common.AcceptInvitationRequest{AgentName: "agent-1", ExternalId: "external-id"}
		suite.Doorman.AcceptInviteResponse = &common.AcceptInvitationResponse{}

		resp, err := target.AcceptInvitation(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
	t.Run("doorman error", func(t *testing.T) {
		target, suite := SetupTest()

		req := &common.AcceptInvitationRequest{AgentName: "agent-1", ExternalId: "external-id"}
		suite.Doorman.AcceptInviteErr = errors.New("boom")

		resp, err := target.AcceptInvitation(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestRegisterEdgeAgent(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		target, suite := SetupTest()
		req := &common.RegisterEdgeAgentRequest{}

		resp := &common.RegisterEdgeAgentResponse{}
		suite.Mediator.RegisterResponse = resp

		result, err := target.RegisterEdgeAgent(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, result)
	})
	t.Run("mediator error", func(t *testing.T) {
		target, suite := SetupTest()
		req := &common.RegisterEdgeAgentRequest{}

		suite.Mediator.RegisterErr = errors.New("BOOM")

		result, err := target.RegisterEdgeAgent(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, result)
	})

}
