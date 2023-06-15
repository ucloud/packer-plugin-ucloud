# UCloud Plugins

The UCloud plugin is able to build
customized images based on an existing base image for use in UHost Instance.

## Components

The Scaffolding plugin is intended as a starting point for creating Packer plugins, containing:

### Builders

- [ucloud-uhost](/docs/builders/uhost.mdx) - The `ucloud-uhost` builder provides the capability to build
  customized images based on an existing base image for use in UHost Instance.

### Post-processors

- [ucloud-import](/docs/post-processors/import.mdx) - The UCloud Import post-processor takes the RAW, VHD, VMDK, or qcow2
  artifact from various builders and imports it to UCloud customized image list
  for UHost Instance.

## Installation

### Using pre-built releases

#### Using the `packer init` command

Starting from version 1.7, Packer supports a new `packer init` command allowing
automatic installation of Packer plugins. Read the
[Packer documentation](https://www.packer.io/docs/commands/init) for more information.

To install this plugin, copy and paste this code into your Packer configuration .
Then, run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    ucloud = {
      version = ">= 1.0.9"
      source  = "github.com/hashicorp/ucloud"
    }
  }
}
```

#### Manual installation

You can find pre-built binary releases of the plugin [here](https://github.com/hashicorp/packer-plugin-name/releases).
Once you have downloaded the latest archive corresponding to your target OS,
uncompress it to retrieve the plugin binary file corresponding to your platform.
To install the plugin, please follow the Packer documentation on
[installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


#### From Source

If you prefer to build the plugin from its source code, clone the GitHub
repository locally and run the command `go build` from the root
directory. Upon successful compilation, a `packer-plugin-ucloud` plugin
binary file can be found in the root directory.
To install the compiled plugin, please follow the official Packer documentation
on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).
