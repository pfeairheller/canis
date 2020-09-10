/*
Copyright Scoir Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package apiserver

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/scoir/canis/pkg/apiserver/api"
	"github.com/scoir/canis/pkg/datastore"
)

func (suite *AdminTestSuite) TestCreateAgent() {
	request := &api.CreateAgentRequest{
		Agent: &api.Agent{
			Id:                  "123",
			Name:                "Test Agent",
			AssignedSchemaId:    "",
			EndorsableSchemaIds: nil,
		},
	}

	a := &datastore.Agent{
		ID:                  "123",
		Name:                "Test Agent",
		AssignedSchemaId:    "",
		EndorsableSchemaIds: nil,
	}

	suite.Store.On("GetAgent", "123").Return(nil, errors.New("not found"))
	suite.Store.On("InsertAgent", a).Return("123", nil)

	resp, err := target.CreateAgent(context.Background(), request)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), resp.Id, "123")
}

func (suite *AdminTestSuite) TestCreateAgentFails() {
	request := &api.CreateAgentRequest{
		Agent: &api.Agent{
			Id:                  "123",
			Name:                "Test Agent",
			AssignedSchemaId:    "",
			EndorsableSchemaIds: nil,
		},
	}

	a := &datastore.Agent{
		ID:                  "123",
		Name:                "Test Agent",
		AssignedSchemaId:    "",
		EndorsableSchemaIds: nil,
	}

	suite.Store.On("GetAgent", "123").Return(nil, errors.New("not found"))
	suite.Store.On("InsertAgent", a).Return("", errors.New("Boom"))

	resp, err := target.CreateAgent(context.Background(), request)
	assert.Nil(suite.T(), resp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "rpc error: code = Internal desc = failed to create agent 123: Boom", err.Error())
}

func (suite *AdminTestSuite) TestCreateAgentMissingRequiredField() {
	request := &api.CreateAgentRequest{
		Agent: &api.Agent{
			Id:                  "",
			Name:                "Test Agent",
			AssignedSchemaId:    "",
			EndorsableSchemaIds: nil,
		},
	}

	resp, err := target.CreateAgent(context.Background(), request)
	assert.Nil(suite.T(), resp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "rpc error: code = InvalidArgument desc = name and id are required fields", err.Error())
}

func (suite *AdminTestSuite) TestCreateAgentAlreadyExists() {
	request := &api.CreateAgentRequest{
		Agent: &api.Agent{
			Id:                  "123",
			Name:                "Test Agent",
			AssignedSchemaId:    "",
			EndorsableSchemaIds: nil,
		},
	}

	suite.Store.On("GetAgent", "123").Return(nil, nil)

	resp, err := target.CreateAgent(context.Background(), request)
	assert.Nil(suite.T(), resp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "rpc error: code = AlreadyExists desc = agent with id 123 already exists", err.Error())
}

func (suite *AdminTestSuite) TestGetAgent() {
	request := &api.GetAgentRequest{
		Id: "123",
	}

	suite.Store.On("GetAgent", "123").Return(&datastore.Agent{ID: "123", Name: "test Agent"}, nil)

	resp, err := target.GetAgent(context.Background(), request)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test Agent", resp.Agent.Name)
}

func (suite *AdminTestSuite) TestGetAgentErr() {
	request := &api.GetAgentRequest{
		Id: "123",
	}

	suite.Store.On("GetAgent", "123").Return(nil, errors.New("BOOM"))

	resp, err := target.GetAgent(context.Background(), request)
	assert.Nil(suite.T(), resp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "rpc error: code = Internal desc = unable to get agent: BOOM", err.Error())
}

func (suite *AdminTestSuite) TestListAgent() {
	request := &api.ListAgentRequest{}

	suite.Store.On("ListAgent", &datastore.AgentCriteria{}).Return(&datastore.AgentList{
		Count:  1,
		Agents: []*datastore.Agent{{ID: "123", Name: "test agent"}},
	}, nil)

	resp, err := target.ListAgent(context.Background(), request)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test agent", resp.Agents[0].Name)
}

func (suite *AdminTestSuite) TestListAgentErr() {
	request := &api.ListAgentRequest{}

	suite.Store.On("ListAgent", &datastore.AgentCriteria{}).Return(nil, errors.New("BOOM"))

	resp, err := target.ListAgent(context.Background(), request)
	assert.Nil(suite.T(), resp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "rpc error: code = Internal desc = unable to list agent: BOOM", err.Error())
}

func (suite *AdminTestSuite) TestDeleteAgent() {
	request := &api.DeleteAgentRequest{
		Id: "123",
	}
	agent := &datastore.Agent{ID: "123"}
	suite.Store.On("GetAgent", "123").Return(agent, nil)
	suite.Store.On("DeleteAgent", "123").Return(nil)

	resp, err := target.DeleteAgent(context.Background(), request)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), resp)
}

func (suite *AdminTestSuite) TestDeleteAgentErr() {
	request := &api.DeleteAgentRequest{
		Id: "123",
	}

	agent := &datastore.Agent{ID: "123"}
	suite.Store.On("GetAgent", "123").Return(agent, nil)
	suite.Store.On("DeleteAgent", "123").Return(errors.New("BOOM"))

	resp, err := target.DeleteAgent(context.Background(), request)
	assert.Nil(suite.T(), resp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "rpc error: code = Internal desc = failed to delete agent 123: BOOM", err.Error())
}

func (suite *AdminTestSuite) TestUpdateAgent() {
	request := &api.UpdateAgentRequest{
		Agent: &api.Agent{
			Id:                  "123",
			Name:                "Test Agent",
			AssignedSchemaId:    "",
			EndorsableSchemaIds: nil,
		},
	}

	a := &datastore.Agent{
		ID:                  "123",
		Name:                "Test Agent",
		AssignedSchemaId:    "",
		EndorsableSchemaIds: nil,
	}

	suite.Store.On("GetAgent", "123").Return(a, nil)
	suite.Store.On("UpdateAgent", a).Return(nil)

	resp, err := target.UpdateAgent(context.Background(), request)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), resp)
}

func (suite *AdminTestSuite) TestCreateSchema() {
	request := &api.CreateSchemaRequest{
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

	s := &datastore.Schema{
		ID:      "123",
		Name:    "Test Schema",
		Version: "0.0.1",
		Attributes: []*datastore.Attribute{{
			Name: "City",
			Type: int32(api.Attribute_STRING),
		}},
	}
	s2 := &datastore.Schema{
		ID:               "123",
		Name:             "Test Schema",
		Version:          "0.0.1",
		ExternalSchemaID: "abc",
		Attributes: []*datastore.Attribute{{
			Name: "City",
			Type: int32(api.Attribute_STRING),
		}},
	}

	suite.Store.On("GetSchema", "123").Return(nil, errors.New("not found"))
	suite.CredRegistry.On("CreateSchema", s).Return("abc", nil)
	suite.Store.On("InsertSchema", s2).Return("123", nil)

	resp, err := target.CreateSchema(context.Background(), request)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), resp.Id, "123")
}

func (suite *AdminTestSuite) TestCreateSchemaFails() {
	request := &api.CreateSchemaRequest{
		Schema: &api.Schema{
			Id:   "123",
			Name: "Test Schema",
		},
	}

	s := &datastore.Schema{
		ID:         "123",
		Name:       "Test Schema",
		Attributes: []*datastore.Attribute{},
	}
	s2 := &datastore.Schema{
		ID:               "123",
		Name:             "Test Schema",
		ExternalSchemaID: "abc",
		Attributes:       []*datastore.Attribute{},
	}

	suite.Store.On("GetSchema", "123").Return(nil, errors.New("not found"))
	suite.CredRegistry.On("CreateSchema", s).Return("abc", nil)
	suite.Store.On("InsertSchema", s2).Return("", errors.New("Boom"))

	resp, err := target.CreateSchema(context.Background(), request)
	assert.Nil(suite.T(), resp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "rpc error: code = Internal desc = failed to create schema 123: Boom", err.Error())
}

func (suite *AdminTestSuite) TestCreateSchemaMissingRequiredField() {
	request := &api.CreateSchemaRequest{
		Schema: &api.Schema{
			Id:   "",
			Name: "Test Schema",
		},
	}

	resp, err := target.CreateSchema(context.Background(), request)
	assert.Nil(suite.T(), resp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "rpc error: code = InvalidArgument desc = name and id are required fields", err.Error())
}

func (suite *AdminTestSuite) TestCreateSchemaAlreadyExists() {
	request := &api.CreateSchemaRequest{
		Schema: &api.Schema{
			Id:   "123",
			Name: "Test Schema",
		},
	}

	suite.Store.On("GetSchema", "123").Return(nil, nil)

	resp, err := target.CreateSchema(context.Background(), request)
	assert.Nil(suite.T(), resp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "rpc error: code = AlreadyExists desc = schema with id 123 already exists", err.Error())
}

func (suite *AdminTestSuite) TestGetSchema() {
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
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test schema", resp.Schema.Name)
}

func (suite *AdminTestSuite) TestGetSchemaErr() {
	request := &api.GetSchemaRequest{
		Id: "123",
	}

	suite.Store.On("GetSchema", "123").Return(nil, errors.New("BOOM"))

	resp, err := target.GetSchema(context.Background(), request)
	assert.Nil(suite.T(), resp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "rpc error: code = Internal desc = unable to get schema: BOOM", err.Error())
}

func (suite *AdminTestSuite) TestListSchema() {
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
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "test schema", resp.Schema[0].Name)
}

func (suite *AdminTestSuite) TestListSchemaErr() {
	request := &api.ListSchemaRequest{}

	suite.Store.On("ListSchema", &datastore.SchemaCriteria{}).Return(nil, errors.New("BOOM"))

	resp, err := target.ListSchema(context.Background(), request)
	assert.Nil(suite.T(), resp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "rpc error: code = Internal desc = unable to list schema: BOOM", err.Error())
}

func (suite *AdminTestSuite) TestDeleteSchema() {
	request := &api.DeleteSchemaRequest{
		Id: "123",
	}

	suite.Store.On("DeleteSchema", "123").Return(nil)

	resp, err := target.DeleteSchema(context.Background(), request)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), resp)
}

func (suite *AdminTestSuite) TestDeleteSchemaErr() {
	request := &api.DeleteSchemaRequest{
		Id: "123",
	}

	suite.Store.On("DeleteSchema", "123").Return(errors.New("BOOM"))

	resp, err := target.DeleteSchema(context.Background(), request)
	assert.Nil(suite.T(), resp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "rpc error: code = Internal desc = failed to delete schema 123: BOOM", err.Error())
}

func (suite *AdminTestSuite) TestUpdateSchema() {
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
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), resp)
}

func (suite *AdminTestSuite) TestUpdateSchemaErr() {
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
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), resp)
}

func (suite *AdminTestSuite) TestUpdateSchemaDataMissing() {
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
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), resp)
}