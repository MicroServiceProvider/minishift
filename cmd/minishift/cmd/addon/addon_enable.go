/*
Copyright (C) 2017 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package addon

import (
	"fmt"
	"github.com/minishift/minishift/pkg/minishift/addon/manager"
	"github.com/minishift/minishift/pkg/util/os/atexit"
	"github.com/spf13/cobra"
)

const (
	priorityFlag   = "priority"
	addOnConfigKey = "addons"

	emptyEnableError       = "An addon name needs to be specified. Use `minishift addons list` to view installed addons."
	noAddOnToEnableMessage = "No addon with name %s installed"
)

var priority int

var addonsEnableCmd = &cobra.Command{
	Use:   "enable ADDON_NAME",
	Short: "Enables the specified addon.",
	Long:  "Enables the specified addon to be run after cluster creation.",
	Run:   runEnableAddon,
}

func init() {
	addonsEnableCmd.Flags().IntVar(&priority, priorityFlag, 0, "The priority of this addon during addon execution.")
	AddonsCmd.AddCommand(addonsEnableCmd)
}

func runEnableAddon(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println(emptyEnableError)
		atexit.Exit(1)
	}

	addonName := args[0]
	addOnManager := GetAddOnManager()

	if !addOnManager.IsInstalled(addonName) {
		fmt.Println(fmt.Sprintf(noAddOnToEnableMessage, addonName))
		return
	}

	enableAddon(addOnManager, addonName, priority)
}

func enableAddon(addOnManager *manager.AddOnManager, addonName string, priority int) {
	addOnConfig, err := addOnManager.Enable(addonName, priority)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to enable plugin %s: %s", addonName, err.Error()))
		atexit.Exit(1)
	}

	addOnConfigMap := getAddOnConfiguration()
	addOnConfigMap[addOnConfig.Name] = addOnConfig
	writeAddOnConfig(addOnConfigMap)
}
