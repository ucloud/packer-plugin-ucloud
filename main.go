package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/hashicorp/packer-plugin-ucloud/builder/ucloud/uhost"
	ucloudimport "github.com/hashicorp/packer-plugin-ucloud/post-processor/ucloud-import"
	"github.com/hashicorp/packer-plugin-ucloud/version"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("uhost", new(uhost.Builder))
	pps.RegisterPostProcessor("import", new(ucloudimport.PostProcessor))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
