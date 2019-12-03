/*
 *  *******************************************************************************
 *  * Copyright (c) 2019 Edgeworx, Inc.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package install

import (
	"fmt"

	"github.com/eclipse-iofog/iofogctl/internal/config"
	"github.com/eclipse-iofog/iofogctl/pkg/iofog"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
)

// LocalAgent uses Container exec commands
type LocalAgent struct {
	defaultAgent
	client           *LocalContainer
	localAgentConfig *LocalAgentConfig
}

func NewLocalAgent(agentConfig *LocalAgentConfig, client *LocalContainer) *LocalAgent {
	return &LocalAgent{
		defaultAgent:     defaultAgent{name: agentConfig.Name},
		localAgentConfig: agentConfig,
		client:           client,
	}
}

func (agent *LocalAgent) Bootstrap() error {
	// Nothing to do for local agent, bootstraping is done inside the image.
	return nil
}

func (agent *LocalAgent) Configure(ctrl *config.Controller, user IofogUser) (uuid string, err error) {
	key, uuid, err := agent.getProvisionKey(ctrl.Endpoint, user)
	if err != nil {
		return "", err
	}

	controllerEndpoint := ctrl.Endpoint
	if util.IsLocalHost(ctrl.Host) {
		IP, err := agent.client.GetContainerIP(GetLocalContainerName("controller"))
		if err != nil {
			return uuid, err
		}
		controllerEndpoint = fmt.Sprintf("%s:%s", IP, iofog.ControllerPortString)
	}

	// Instantiate provisioning commands
	controllerBaseURL := fmt.Sprintf("http://%s/api/v3", controllerEndpoint)
	cmds := [][]string{
		{"iofog-agent", "config", "-idc", "off"},
		{"iofog-agent", "config", "-a", controllerBaseURL},
		{"iofog-agent", "provision", key},
	}

	// TODO: Verify provisioning succeeded

	// Execute commands
	for _, cmd := range cmds {
		_, err = agent.client.ExecuteCmd(agent.localAgentConfig.ContainerName, cmd)
		if err != nil {
			return
		}
	}

	return
}
