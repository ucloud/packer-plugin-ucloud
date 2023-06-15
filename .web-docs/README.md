The UCloud plugin is able to build customized images based on an existing base image for use in UHost Instance.

### Installation
To install this plugin add this code into your Packer configuration and run [packer init](/packer/docs/commands/init)

```hcl
packer {
    required_plugins {
        ucloud = {
          version = "~> 1"
          source = "github.com/hashicorp/ucloud"
        }
    }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
packer plugins install github.com/hashicorp/ucloud
```
### Components

#### Builders

- [ucloud-uhost](/packer/integrations/hashicorp/ucloud/latest/components/builder/uhost) - The `ucloud-uhost` builder provides the capability to build
  customized images based on an existing base image for use in UHost Instance.

#### Post-processors

- [ucloud-import](/packer/integrations/hashicorp/uhost/latest/components/post-processor/import) - The UCloud Import post-processor takes the RAW, VHD, VMDK, or qcow2
  artifact from various builders and imports it to UCloud customized image list
  for UHost Instance.
