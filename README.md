# Packer Plugin Ucloud
The `Ucloud` multi-component plugin can be used with HashiCorp
[Packer](https://www.packer.io) to create custom images. For the full list of
available features for this plugin see [docs](docs).

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
      source  = "github.com/ucloud/ucloud"
    }
  }
}
```


#### Manual installation

You can find pre-built binary releases of the plugin
[here](https://github.com/ucloud/packer-plugin-ucloud/releases). Once you have
downloaded the latest archive corresponding to your target OS, uncompress it to
retrieve the plugin binary file corresponding to your platform. To install the
plugin, please follow the Packer documentation on [installing a
plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


### From Sources

If you prefer to build the plugin from sources, clone the GitHub repository
locally and run the command `go build` from the root
directory. Upon successful compilation, a `packer-plugin-ucloud` plugin
binary file can be found in the root directory.
To install the compiled plugin, please follow the official Packer documentation
on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


### Configuration

For more information on how to configure the plugin, please read the
documentation located in the [`docs/`](docs) directory.

## Usage Example

Before you try the example in this section, you should get public, private key
and project id from UCloud console and export the following environment varialbes:

- `UCLOUD_PRIVATE_KEY`
- `UCLOUD_PUBLIC_KEY`
- `UCLOUD_PROJECT_ID`

Then create a packer instruction file in HCL syntax as follows:


    packer {
      required_plugins {
        ucloud = {
          version = ">= 1.0.9"
          source  = "github.com/ucloud/ucloud"
        }
      }
    }

    variable "ucloud_private_key" {
      type    = string
      default = "${env("UCLOUD_PRIVATE_KEY")}"
    }

    variable "ucloud_project_id" {
      type    = string
      default = "${env("UCLOUD_PROJECT_ID")}"
    }

    variable "ucloud_public_key" {
      type    = string
      default = "${env("UCLOUD_PUBLIC_KEY")}"
    }

    source "ucloud-uhost" "centos7-image-seed" {
      availability_zone = "cn-sh2-02"
      image_name        = "uk8s-centos-7.9"
      instance_type     = "o-standard-2"
      private_key       = "${var.ucloud_private_key}"
      project_id        = "${var.ucloud_project_id}"
      public_key        = "${var.ucloud_public_key}"
      region            = "cn-sh2"
      source_image_id   = "uimage-vj5ui3"
      ssh_username      = "root"
    }

    build {
      name = "UK8S-Centos7-BaseImage"
      sources = ["source.ucloud-uhost.centos7-image-seed"]


      provisioner "shell" {
          inline = ["echo hello"]
      }

      provisioner "ansible" {
        playbook_file = "./playbook-centos.yml"
      }
    }

Name the file as `xxx.pkr.hcl` and run `packer init` command under the
directory where the above file located:

    packer init xxx.pkr.hcl

Packer will install the required plugins automatically. Then you run packer to
build the image for uhost. It is recommended that you enable the packer log to
aid diagnosis if something goes wrong. You run commands as follows:


    export UCLOUD_PRIVATE_KEY="****************"
    export UCLOUD_PUBLIC_KEY="******************"
    export UCLOUD_PROJECT_ID="xxx"
    export PACKER_LOG=1
    export PACKER_LOG_PATH="packer.log"
    packer build xxx.pkr.hcl

## Contributing

* If you think you've found a bug in the code or you have a question regarding
  the usage of this software, please reach out to us by opening an issue in
  this GitHub repository.
* Contributions to this project are welcome: if you want to add a feature or a
  fix a bug, please do so by opening a Pull Request in this GitHub repository.
  In case of feature contribution, we kindly ask you to open an issue to
  discuss it beforehand.

## GPG public key
You may verify the released binaries with this GPG public key:

    -----BEGIN PGP PUBLIC KEY BLOCK-----

    mDMEZCbtBxYJKwYBBAHaRw8BAQdAzED5JgKHHe17uuPWzoU8IJRK0bqBF9S+KdVn
    aCq1tg+0N0p1c3RpbiBaaGFuZyAoVUNsb3VkIHdvcmsgR1BHKSA8anVzdGluLnpo
    YW5nQHVjbG91ZC5jbj6ImgQTFgoAQhYhBPi+ba2/BwN794NuPH7pFs0ifsLFBQJk
    Ju0HAhsDBQkDwmcABQsJCAcCAyICAQYVCgkICwIEFgIDAQIeBwIXgAAKCRB+6RbN
    In7CxTwnAQD739fTMO0e4LSvBYYqMT0OEel4/MYJMVXdftLB9CnfQAEAkel163cV
    YCaW553KyQOJCI9aJSasopYju3lIUvkWfgGJATMEEAEIAB0WIQRFzhtgzKPwfHuQ
    /nYgH5GJOOn1+gUCZCwziwAKCRAgH5GJOOn1+nG4CACT066mVIGq2dBPmf/7oNQS
    JTZt1IdbYAkP3E7YhLolCDinDHjaw25JevR3iDIhPqqlaogrlIC2RBxsDg7zrv+q
    SSEIkebW8BYYgc7K+1hpse4/V/jXJ1b+aFPKC3uvtQX8wU5bgOIctvhxl9agXHqd
    MLtlbsGav0KGu65DElTLsoyEhaakWVatyGZ4tigc13o3CH6uUJCgaptYPc6J4p6N
    we6q0s1FnNRYe3+ov+WvJvz+k2Qjr5o4n9HFajlueZcs1rAf6Eehpa/Imcp5NXCT
    0q453b0atLykK567Xc9v+gRRsJtDXV+rHNV38eAHazfwrEyVapkQ/JGJtAPI1/5n
    uDgEZCbtBxIKKwYBBAGXVQEFAQEHQLZykGFXbDtodQLzuIclSAK4//NnT6F3kXC5
    M7flA55EAwEIB4h+BBgWCgAmFiEE+L5trb8HA3v3g248fukWzSJ+wsUFAmQm7QcC
    GwwFCQPCZwAACgkQfukWzSJ+wsU9HwD+NAsTt8Odq4arxJd1t+bflUFqplvqbBwB
    eXTsLTAKa2QBAI/ODojRoHtRU/rka6kiQJzFjN4nEINvoNyeem/PMaoE
    =6uAG
    -----END PGP PUBLIC KEY BLOCK-----

