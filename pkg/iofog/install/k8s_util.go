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

package install

import (
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newService(namespace string, ms *microservice) *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ms.name,
			Namespace: namespace,
			Labels: map[string]string{
				"name": ms.name,
			},
		},
		Spec: corev1.ServiceSpec{
			Type:                  corev1.ServiceTypeLoadBalancer,
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeLocal,
			LoadBalancerIP:        ms.IP,
			Selector: map[string]string{
				"name": ms.name,
			},
		},
	}
	// Add ports
	for idx, port := range ms.ports {
		svcPort := corev1.ServicePort{
			Name:       ms.name + strconv.Itoa(idx),
			Port:       int32(port),
			TargetPort: intstr.FromInt(int(port)),
			Protocol:   corev1.Protocol("TCP"),
		}
		svc.Spec.Ports = append(svc.Spec.Ports, svcPort)
	}
	return svc
}

func newDeployment(namespace string, ms *microservice) *appsv1.Deployment {
	maxUnavailable := intstr.FromInt(0)
	maxSurge := intstr.FromInt(1)
	strategy := appsv1.DeploymentStrategy{
		Type: appsv1.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &appsv1.RollingUpdateDeployment{
			MaxUnavailable: &maxUnavailable,
			MaxSurge:       &maxSurge,
		},
	}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ms.name,
			Namespace: namespace,
			Labels: map[string]string{
				"name": ms.name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &ms.replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": ms.name,
				},
			},
			Strategy: strategy,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": ms.name,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: ms.name,
				},
			},
		},
	}
	containers := &dep.Spec.Template.Spec.Containers
	for _, msCont := range ms.containers {
		cont := corev1.Container{
			Name:            msCont.name,
			Image:           msCont.image,
			ImagePullPolicy: corev1.PullPolicy(msCont.imagePullPolicy),
			Args:            msCont.args,
			ReadinessProbe:  msCont.readinessProbe,
			Ports:           msCont.ports,
			Env:             msCont.env,
			Command:         msCont.command,
		}
		*containers = append(*containers, cont)
	}
	return dep
}

func newServiceAccount(namespace string, ms *microservice) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ms.name,
			Namespace: namespace,
		},
	}
}

func getClusterRoleBindingName(namespace, resourceName string) string {
	return namespace + "-" + resourceName
}

func newClusterRoleBinding(namespace string, ms *microservice) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: getClusterRoleBindingName(namespace, ms.name),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      ms.name,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}

func newRoleBinding(namespace string, ms *microservice) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ms.name,
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "ServiceAccount",
				Name: ms.name,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     ms.name,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}

func newRole(namespace string, ms *microservice) *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ms.name,
			Namespace: namespace,
		},
		Rules: ms.rbacRules,
	}
}
