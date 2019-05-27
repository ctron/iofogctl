package deploycontroller

import (
	"github.com/eclipse-iofog/cli/pkg/config"
	"fmt"
)

type controller struct {
	configManager *config.Manager
}

func new() *controller {
	c := &controller{}
	c.configManager = config.NewManager()
	return c
}

func (ctrl *controller) execute(namespace, name string) error {
	// TODO (Serge) Execute back-end logic

	// Update configuration
	configEntry := config.Controller{ Name: name, User: "none" }
	err := ctrl.configManager.AddController(namespace, configEntry)

	// TODO (Serge) Handle config file error, retry..?

	if err == nil {
		fmt.Printf("\nController %s/%s successfully deployed.\n", namespace, name)
	}
	return err
}