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

package deploycontrolplane

import (
	"github.com/eclipse-iofog/iofogctl/internal/config"
	"github.com/eclipse-iofog/iofogctl/internal/deploy/connector"
	"github.com/eclipse-iofog/iofogctl/internal/deploy/controller"
	"github.com/eclipse-iofog/iofogctl/internal/execute"
	"github.com/eclipse-iofog/iofogctl/pkg/iofog/client"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
)

type Options struct {
	Namespace string
	InputFile string
}

func Execute(opt Options) error {
	// Check the namespace exists
	_, err := config.GetNamespace(opt.Namespace)
	if err != nil {
		return err
	}

	// Read the input file
	controlPlane, err := UnmarshallYAML(opt.InputFile)
	if err != nil {
		return err
	}

	// Require Controllers
	if len(controlPlane.Controllers) == 0 {
		return util.NewInputError("You must specify atleast one Controller in the Control Plane")
	}
	// Require IofogUser
	if controlPlane.IofogUser.Email == "" || controlPlane.IofogUser.Name == "" || controlPlane.IofogUser.Password == "" || controlPlane.IofogUser.Surname == "" {
		return util.NewInputError("You must specify an ioFog user with a name, surname, email, and password")
	}

	// Instantiate executors
	var executors []execute.Executor

	// Execute Controllers
	for idx := range controlPlane.Controllers {
		exe, err := deploycontroller.NewExecutor(opt.Namespace, controlPlane.Controllers[idx], controlPlane)
		if err != nil {
			return err
		}
		executors = append(executors, exe)
	}
	if err := runExecutors(executors); err != nil {
		return err
	}

	// Execute Connectors
	executors = executors[:0]
	for idx := range controlPlane.Connectors {
		exe, err := deployconnector.NewExecutor(opt.Namespace, controlPlane.Connectors[idx], controlPlane)
		if err != nil {
			return err
		}
		executors = append(executors, exe)
	}
	if err := runExecutors(executors); err != nil {
		return err
	}

	// Create new user
	// TODO: replace with controlplane variable for endpoint
	ctrlClient := client.New(controlPlane.Controllers[0].Endpoint)
	if err = ctrlClient.CreateUser(client.User(controlPlane.IofogUser)); err != nil {
		return err
	}

	// Update config
	if err = config.UpdateControlPlane(opt.Namespace, controlPlane); err != nil {
		return err
	}

	return config.Flush()
}

func runExecutors(executors []execute.Executor) error {
	if errs, failedExes := execute.ForParallel(executors); len(errs) > 0 {
		for idx := range errs {
			util.PrintNotify("Error from " + failedExes[idx].GetName() + ": " + errs[idx].Error())
		}
		return util.NewError("Failed to deploy")
	}
	return nil
}