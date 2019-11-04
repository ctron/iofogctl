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

package cmd

import (
	delete "github.com/eclipse-iofog/iofogctl/internal/delete/all"
	"github.com/eclipse-iofog/iofogctl/pkg/util"
	"github.com/spf13/cobra"
)

func newDeleteAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all",
		Short: "Delete all resources within a namespace",
		Long: `Delete all resources within a namespace.

Tears down all components of an Edge Compute Network.

If you don't want to tear down the deployments but would like to free up the namespace, use the disconnect command instead.`,
		Example: `iofogctl delete all -n NAMESPACE`,
		Run: func(cmd *cobra.Command, args []string) {
			namespace, err := cmd.Flags().GetString("namespace")
			util.Check(err)
			// Execute command
			err = delete.Execute(namespace)
			util.Check(err)

			util.PrintSuccess("Successfully deleted all resources in namespace " + namespace)
		},
	}

	return cmd
}
