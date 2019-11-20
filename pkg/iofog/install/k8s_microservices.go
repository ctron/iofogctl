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
	"github.com/eclipse-iofog/iofogctl/pkg/iofog"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type microservice struct {
	name        string
	IP          string
	ports       []int32
	serviceType string
	replicas    int32
	containers  []container
}

type container struct {
	name            string
	image           string
	imagePullPolicy string
	args            []string
	readinessProbe  *v1.Probe
	env             []v1.EnvVar
	command         []string
	ports           []v1.ContainerPort
}

var controllerMicroservice = microservice{
	name:        "controller",
	ports:       []int32{iofog.ControllerPort, 80},
	replicas:    1,
	serviceType: "LoadBalancer",
	containers: []container{
		{
			name:            "controller",
			image:           "iofog/controller:" + util.GetControllerTag(),
			imagePullPolicy: "Always",
			readinessProbe: &v1.Probe{
				Handler: v1.Handler{
					HTTPGet: &v1.HTTPGetAction{
						Path: "/api/v3/status",
						Port: intstr.FromInt(iofog.ControllerPort),
					},
				},
				InitialDelaySeconds: 1,
				PeriodSeconds:       4,
				FailureThreshold:    3,
			},
		},
	},
}

var connectorMicroservice = microservice{
	name: "connector",
	ports: []int32{
		iofog.ConnectorPort,
		6000, 6001, 6002, 6003, 6004, 6005, 6006, 6007, 6008, 6009,
		6010, 6011, 6012, 6013, 6014, 6015, 6016, 6017, 6018, 6019,
		6020, 6021, 6022, 6023, 6024, 6025, 6026, 6027, 6028, 6029,
		6030, 6031, 6032, 6033, 6034, 6035, 6036, 6037, 6038, 6039,
		6040, 6041, 6042, 6043, 6044, 6045, 6046, 6047, 6048, 6049,
		6050,
	},
	replicas:    1,
	serviceType: "LoadBalancer",
	containers: []container{
		{
			name:            "connector",
			image:           "iofog/connector:" + util.GetConnectorTag(),
			imagePullPolicy: "Always",
		},
	},
}

var schedulerMicroservice = microservice{
	name:     "scheduler",
	replicas: 1,
	containers: []container{
		{
			name:            "scheduler",
			image:           "iofog/iofog-scheduler:" + util.GetSchedulerTag(),
			imagePullPolicy: "Always",
		},
	},
}

var operatorMicroservice = microservice{
	name:     "iofog-operator",
	ports:    []int32{60000},
	replicas: 1,
	containers: []container{
		{
			name:            "iofog-operator",
			image:           "iofog/iofog-operator:" + util.GetOperatorTag(),
			imagePullPolicy: "Always",
			readinessProbe: &v1.Probe{
				Handler: v1.Handler{
					Exec: &v1.ExecAction{
						Command: []string{
							"stat",
							"/tmp/operator-sdk-ready",
						},
					},
				},
				InitialDelaySeconds: 4,
				PeriodSeconds:       10,
				FailureThreshold:    1,
			},
			env: []v1.EnvVar{
				{
					Name: "WATCH_NAMESPACE",
					ValueFrom: &v1.EnvVarSource{
						FieldRef: &v1.ObjectFieldSelector{
							FieldPath: "metadata.namespace",
						},
					},
				},
				{
					Name: "POD_NAME",
					ValueFrom: &v1.EnvVarSource{
						FieldRef: &v1.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				},
				{
					Name:  "OPERATOR_NAME",
					Value: "iofog-operator",
				},
			},
			ports: []v1.ContainerPort{
				{
					ContainerPort: int32(60000),
					Name:          "metrics",
				},
			},
			command: []string{
				"iofog-operator",
			},
		},
	},
}

var kubeletMicroservice = microservice{
	name:     "kubelet",
	ports:    []int32{60000},
	replicas: 1,
	containers: []container{
		{
			name:            "kubelet",
			image:           "iofog/iofog-kubelet:" + util.GetKubeletTag(),
			imagePullPolicy: "Always",
			args: []string{
				"--namespace",
				"",
				"--iofog-token",
				"",
				"--iofog-url",
				"",
			},
		},
	},
}
