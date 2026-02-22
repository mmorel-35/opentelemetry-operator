// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/open-telemetry/opamp-go/protobufs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAgent(t *testing.T) {
	logger := logr.Discard()
	instanceId := uuid.New()
	conn := &mockConnection{}

	agent := NewAgent(logger, instanceId, conn)
	require.NotNilf(t, agent, "agent should not be nil")
	assert.Equalf(t, instanceId, agent.InstanceId, "instance ID should match")
	assert.Equalf(t, conn, agent.conn, "connection should match")
}

func TestAgent_hasCapability(t *testing.T) {
	logger := logr.Discard()
	instanceId := uuid.New()
	conn := &mockConnection{}
	agent := NewAgent(logger, instanceId, conn)

	agent.Status = &protobufs.AgentToServer{
		Capabilities: uint64(protobufs.AgentCapabilities_AgentCapabilities_ReportsEffectiveConfig),
	}

	assert.Truef(t, agent.hasCapability(protobufs.AgentCapabilities_AgentCapabilities_ReportsEffectiveConfig), "should have capability")
	assert.Falsef(t, agent.hasCapability(protobufs.AgentCapabilities_AgentCapabilities_ReportsHealth), "should not have capability")
}

func TestAgent_UpdateStatus(t *testing.T) {
	logger := logr.Discard()
	instanceId := uuid.New()
	conn := &mockConnection{}
	agent := NewAgent(logger, instanceId, conn)

	newStatus := &protobufs.AgentToServer{
		SequenceNum: 1,
		AgentDescription: &protobufs.AgentDescription{
			IdentifyingAttributes: []*protobufs.KeyValue{
				{Key: "service.name", Value: &protobufs.AnyValue{Value: &protobufs.AnyValue_StringValue{StringValue: "test-service"}}},
			},
		},

		Capabilities: uint64(protobufs.AgentCapabilities_AgentCapabilities_ReportsEffectiveConfig),
	}

	response := &protobufs.ServerToAgent{}
	agent.UpdateStatus(newStatus, response)

	assert.Equalf(t, newStatus, agent.Status, "status should be updated")
	assert.NotEqualf(t, 0, response.Flags&uint64(protobufs.ServerToAgentFlags_ServerToAgentFlags_ReportFullState), "should request full state")
}

func TestAgent_GetHealth(t *testing.T) {
	logger := logr.Discard()
	instanceId := uuid.New()
	conn := &mockConnection{}
	agent := NewAgent(logger, instanceId, conn)

	health := &protobufs.ComponentHealth{
		Healthy: true,
	}
	agent.health = health

	assert.Equalf(t, health, agent.GetHealth(), "health should match")
}

func TestAgent_GetConfiguration(t *testing.T) {
	logger := logr.Discard()
	instanceId := uuid.New()
	conn := &mockConnection{}
	agent := NewAgent(logger, instanceId, conn)

	config := &protobufs.EffectiveConfig{
		ConfigMap: &protobufs.AgentConfigMap{
			ConfigMap: map[string]*protobufs.AgentConfigFile{
				"config.yaml": {Body: []byte("receivers:\n  otlp:\n")},
			},
		},
	}
	agent.effectiveConfig = config

	assert.Equalf(t, config, agent.GetConfiguration(), "configuration should match")
}
