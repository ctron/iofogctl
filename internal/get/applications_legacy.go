/*
 *  *******************************************************************************
 *  * Copyright (c) 2020 Edgeworx, Inc.
 *  *
 *  * This program and the accompanying materials are made available under the
 *  * terms of the Eclipse Public License v. 2.0 which is available at
 *  * http://www.eclipse.org/legal/epl-2.0
 *  *
 *  * SPDX-License-Identifier: EPL-2.0
 *  *******************************************************************************
 *
 */

package get

import "github.com/eclipse-iofog/iofogctl/v2/pkg/util"

func (exe *applicationExecutor) initLegacy() (err error) {
	flows, err := exe.client.GetAllFlows()
	if err != nil {
		return
	}
	exe.flows = flows.Flows
	for _, flow := range exe.flows {
		listMsvcs, err := exe.client.GetMicroservicesPerFlow(flow.ID)
		if err != nil {
			return err
		}

		// Filter System microservices
		for _, ms := range listMsvcs.Microservices {
			if util.IsSystemMsvc(ms) {
				continue
			}
			exe.msvcsPerApplication[flow.ID] = append(exe.msvcsPerApplication[flow.ID], ms)
		}
	}
	return nil
}
