Type: `ucloud-uhost`
Artifact BuilderId: `ucloud.uhost`

The `ucloud-uhost` Packer builder plugin provides the capability to build
customized images based on an existing base image for use in UHost Instance.

This builder builds an UCloud image by launching an UHost instance from a source image,
provisioning that running machine, and then creating an image from that machine.

## Configuration Reference

The following configuration options are available for building UCloud images. They are
segmented below into two categories: required and optional parameters.

In addition to the options listed here, a
[communicator](/packer/docs/templates/legacy_json_templates/communicator) can be configured for this
builder.

~> **Note:** The builder doesn't support Windows images for now and only supports CentOS and Ubuntu images via SSH authentication with `ssh_username` (Required) and `ssh_password` (Optional). The `ssh_username` must be `root` for CentOS images and `ubuntu` for Ubuntu images. The `ssh_password` may contain 8-30 characters, and must consist of at least 2 items out of the capital letters, lower case letters, numbers and special characters. The special characters include `()~!@#\$%^&\*-+=\_|{}\[]:;'<>,.?/`.

### Required:

<!-- Code generated from the comments of the AccessConfig struct in builder/ucloud/common/access_config.go; DO NOT EDIT MANUALLY -->

- `public_key` (string) - This is the UCloud public key. It must be provided unless `profile` is set,
  but it can also be sourced from the `UCLOUD_PUBLIC_KEY` environment variable.

- `private_key` (string) - This is the UCloud private key. It must be provided unless `profile` is set,
  but it can also be sourced from the `UCLOUD_PRIVATE_KEY` environment variable.

- `region` (string) - This is the UCloud region. It must be provided, but it can also be sourced from
  the `UCLOUD_REGION` environment variables.

- `project_id` (string) - This is the UCloud project id. It must be provided, but it can also be sourced
  from the `UCLOUD_PROJECT_ID` environment variables.

<!-- End of code generated from the comments of the AccessConfig struct in builder/ucloud/common/access_config.go; -->


<!-- Code generated from the comments of the RunConfig struct in builder/ucloud/common/run_config.go; DO NOT EDIT MANUALLY -->

- `availability_zone` (string) - This is the UCloud availability zone where UHost instance is located. such as: `cn-bj2-02`.
  You may refer to [list of availability_zone](https://docs.ucloud.cn/api/summary/regionlist)

- `source_image_id` (string) - This is the ID of base image which you want to create your customized images with.

- `instance_type` (string) - The type of UHost instance.
  You may refer to [list of instance type](https://docs.ucloud.cn/compute/terraform/specification/instance)

<!-- End of code generated from the comments of the RunConfig struct in builder/ucloud/common/run_config.go; -->


<!-- Code generated from the comments of the ImageConfig struct in builder/ucloud/common/image_config.go; DO NOT EDIT MANUALLY -->

- `image_name` (string) - The name of the user-defined image, which contains 1-63 characters and only
  support Chinese, English, numbers, '-\_,.:[]'.

<!-- End of code generated from the comments of the ImageConfig struct in builder/ucloud/common/image_config.go; -->


### Optional:

<!-- Code generated from the comments of the AccessConfig struct in builder/ucloud/common/access_config.go; DO NOT EDIT MANUALLY -->

- `base_url` (string) - This is the base url. (Default: `https://api.ucloud.cn`).

- `profile` (string) - This is the UCloud profile name as set in the shared credentials file, it can
  also be sourced from the `UCLOUD_PROFILE` environment variables.

- `shared_credentials_file` (string) - This is the path to the shared credentials file, it can also be sourced from
  the `UCLOUD_SHARED_CREDENTIAL_FILE` environment variables. If this is not set
  and a profile is specified, `~/.ucloud/credential.json` will be used.

<!-- End of code generated from the comments of the AccessConfig struct in builder/ucloud/common/access_config.go; -->


<!-- Code generated from the comments of the RunConfig struct in builder/ucloud/common/run_config.go; DO NOT EDIT MANUALLY -->

- `instance_name` (string) - The name of instance, which contains 1-63 characters and only support Chinese,
  English, numbers, '-', '\_', '.'.

- `boot_disk_type` (string) - The type of boot disk associated to UHost instance.
  Possible values are: `cloud_ssd` and `cloud_rssd` for cloud boot disk, `local_normal` and `local_ssd`
  for local boot disk. (Default: `cloud_ssd`). The `cloud_ssd` and `local_ssd` are not fully supported
  by all regions as boot disk type, please proceed to UCloud console for more details.
  
  ~> **Note:** It takes around 10 mins for boot disk initialization when `boot_disk_type` is `local_normal` or `local_ssd`.

- `boot_disk_size` (int) - The size of boot disk associated to UHost instance, which cannot be smaller than the size of source image.
  The unit is `GB`. Default value is the size of source image.

- `vpc_id` (string) - The ID of VPC linked to the UHost instance. If not defined `vpc_id`, the instance will use the default VPC in the current region.

- `subnet_id` (string) - The ID of subnet under the VPC. If `vpc_id` is defined, the `subnet_id` is mandatory required.
  If `vpc_id` and `subnet_id` are not defined, the instance will use the default subnet in the current region.

- `security_group_id` (string) - The ID of the fire wall associated to UHost instance. If `security_group_id` is not defined,
  the instance will use the non-recommended web fire wall, and open port include 22, 3389 by default.
  It is supported by ICMP fire wall protocols.
  You may refer to [security group_id](https://docs.ucloud.cn/network/firewall/firewall).

- `eip_bandwidth` (int) - Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second). (Default: `10`).

- `eip_charge_mode` (string) - Elastic IP charge mode. Possible values are: `traffic` as pay by traffic, `bandwidth` as pay by bandwidth,
  `post_accurate_bandwidth` as post pay mode. (Default: `traffic`).
  Note currently default `traffic` eip charge mode not not fully support by all `availability_zone`
  in the `region`, please proceed to [UCloud console](https://console.ucloud.cn/unet/eip/create) for more details.
  You may refer to [eip introduction](https://docs.ucloud.cn/unet/eip/introduction).

- `user_data` (string) - User data to apply when launching the instance.
  Note that you need to be careful about escaping characters due to the templates
  being JSON. It is often more convenient to use user_data_file, instead.
  Packer will not automatically wait for a user script to finish before
  shutting down the instance this must be handled in a provisioner.
  You may refer to [user_data_document](https://docs.ucloud.cn/uhost/guide/metadata/userdata)

- `user_data_file` (string) - Path to a file that will be used for the user data when launching the instance.

- `min_cpu_platform` (string) - Specifies a minimum CPU platform for the the VM instance. (Default: `Intel/Auto`).
  You may refer to [min_cpu_platform](https://docs.ucloud.cn/uhost/introduction/uhost/type_new)
     - The Intel CPU platform:
         - `Intel/Auto` as the Intel CPU platform version will be selected randomly by system;
         - `Intel/IvyBridge` as Intel V2, the version of Intel CPU platform selected by system will be `Intel/IvyBridge` and above;
         - `Intel/Haswell` as Intel V3,  the version of Intel CPU platform selected by system will be `Intel/Haswell` and above;
         - `Intel/Broadwell` as Intel V4, the version of Intel CPU platform selected by system will be `Intel/Broadwell` and above;
         - `Intel/Skylake` as Intel V5, the version of Intel CPU platform selected by system will be `Intel/Skylake` and above;
         - `Intel/Cascadelake` as Intel V6, the version of Intel CPU platform selected by system will be `Intel/Cascadelake`;
     - The AMD CPU platform:
         - `Amd/Auto` as the Amd CPU platform version will be selected randomly by system;
         - `Amd/Epyc2` as the version of Amd CPU platform selected by system will be `Amd/Epyc2` and above;

- `use_ssh_private_ip` (bool) - If this value is true, packer will connect to the created UHost instance via a private ip
  instead of allocating an EIP (elastic public ip).(Default: `false`).

<!-- End of code generated from the comments of the RunConfig struct in builder/ucloud/common/run_config.go; -->


<!-- Code generated from the comments of the ImageConfig struct in builder/ucloud/common/image_config.go; DO NOT EDIT MANUALLY -->

- `image_description` (string) - The description of the image.

- `image_copy_to_mappings` ([]ImageDestination) - The array of mappings regarding the copied images to the destination regions and projects.
  
   - `project_id` (string) - The destination project id, where copying image in.
  
   - `region` (string) - The destination region, where copying image in.
  
   - `name` (string) - The copied image name. If not defined, builder will use `image_name` as default name.
  
   - `description` (string) - The copied image description.
  
  ```json
  {
    "image_copy_to_mappings": [
      {
        "project_id": "{{user `ucloud_project_id`}}",
        "region": "cn-sh2",
        "description": "test",
        "name": "packer-test-basic-sh"
      }
    ]
  }
  ```

- `wait_image_ready_timeout` (int) - Timeout of creating image or copying image. The default timeout is 3600 seconds if this option
  is not set or is set to 0.

<!-- End of code generated from the comments of the ImageConfig struct in builder/ucloud/common/image_config.go; -->


## Examples

Here is a basic example for build UCloud CentOS image:

**JSON**

```json
{
  "variables": {
    "ucloud_public_key": "{{env `UCLOUD_PUBLIC_KEY`}}",
    "ucloud_private_key": "{{env `UCLOUD_PRIVATE_KEY`}}",
    "ucloud_project_id": "{{env `UCLOUD_PROJECT_ID`}}"
  },

  "builders": [
    {
      "type": "ucloud-uhost",
      "public_key": "{{user `ucloud_public_key`}}",
      "private_key": "{{user `ucloud_private_key`}}",
      "project_id": "{{user `ucloud_project_id`}}",
      "region": "cn-bj2",
      "availability_zone": "cn-bj2-02",
      "instance_type": "n-basic-2",
      "source_image_id": "uimage-f1chxn",
      "ssh_username": "root",
      "image_name": "packer-test{{timestamp}}"
    }
  ]
}
```

**HCL2**

```hcl
// .pkr.hcl file
variable "ucloud_public_key" {
  type = string
  default = "xxx"
}

variable "ucloud_private_key" {
  type = string
  default = "xxx"
}

variable "ucloud_project_id" {
  type = string
  default = "xxx"
}

source "ucloud-uhost" "basic-example" {
  public_key        =  var.ucloud_public_key
  private_key       =  var.ucloud_private_key
  project_id        =  var.ucloud_project_id
  region            =  "cn-bj2"
  availability_zone =  "cn-bj2-02"
  instance_type     =  "n-basic-2"
  source_image_id   =  "uimage-f1chxn"
  ssh_username      =  "root"
}

build {
  source "sources.ucloud-uhost.basic-example" {
    image_name =  "packer-test-${timestamp()}"
  }
}
```


Here is a example for build UCloud Ubuntu image:

**JSON**

```json
{
  "variables": {
    "ucloud_public_key": "{{env `UCLOUD_PUBLIC_KEY`}}",
    "ucloud_private_key": "{{env `UCLOUD_PRIVATE_KEY`}}",
    "ucloud_project_id": "{{env `UCLOUD_PROJECT_ID`}}",
    "password": "ucloud_2020"
  },

  "builders": [
    {
      "type": "ucloud-uhost",
      "public_key": "{{user `ucloud_public_key`}}",
      "private_key": "{{user `ucloud_private_key`}}",
      "project_id": "{{user `ucloud_project_id`}}",
      "region": "cn-bj2",
      "availability_zone": "cn-bj2-02",
      "instance_type": "n-basic-2",
      "source_image_id": "uimage-irofn4",
      "ssh_password": "{{user `password`}}",
      "ssh_username": "ubuntu",
      "image_name": "packer-test-ubuntu{{timestamp}}"
    }
  ],

  "provisioners": [
    {
      "type": "shell",
      "execute_command": "echo '{{user `password`}}' | sudo -S '{{.Path}}'",
      "inline": ["sleep 30", "sudo apt update", "sudo apt install nginx -y"]
    }
  ]
}
```

**HCL2**

```hcl
// .pkr.hcl file
variable "ucloud_public_key" {
  type = string
  default = "xxx"
}

variable "ucloud_private_key" {
  type = string
  default = "xxx"
}

variable "ucloud_project_id" {
  type = string
  default = "xxx"
}

variable "password" {
  type  = string
  default = "ucloud_2020"
}

source "ucloud-uhost" "basic-example" {
  public_key        =  var.ucloud_public_key
  private_key       =  var.ucloud_private_key
  project_id        =  var.ucloud_project_id
  region            =  "cn-bj2"
  availability_zone =  "cn-bj2-02"
  instance_type     =  "n-basic-2"
  ssh_password      =  var.password
  source_image_id   =  "uimage-irofn4"
  ssh_username      =  "ubuntu"
}

build {
  source "sources.ucloud-uhost.basic-example" {
    image_name =  "packer-test-ubuntu-${timestamp()}"
  }
  provisioner "shell" {
      execute_command = "echo '${var.password}' | sudo -S '{{.Path}}'"
      inline          = ["sleep 30", "sudo apt update", "sudo apt install nginx -y"]
  }
}
```


-> **Note:** Packer can also read the public key and private key from
environmental variables. See the configuration reference in the section above
for more information on what environmental variables Packer will look for.

~> **Note:** Source image may be deprecated after a while, you can use the tools like [UCloud CLI](https://docs.ucloud.cn/cli/intro) to run `ucloud image list` to find one that exists.
