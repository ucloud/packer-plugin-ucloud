Type: `ucloud-import`
Artifact BuilderId: `packer.post-processor.ucloud-import`

The Packer UCloud Import post-processor takes the RAW, VHD, VMDK, or qcow2 artifact from various builders and imports it to UCloud customized image list for UHost Instance.

~> **Note** Some regions don't support image import. You may refer to [ucloud console](https://console.ucloud.cn/uhost/uimage) for more detail. If you want to import to unsupported regions, please import the image in `cn-bj2` first, and then copy the image to the target region.

## How Does it Work?

The import process operates by making a temporary copy of the RAW, VHD, VMDK, or qcow2 to an UFile bucket, and calling an import task in UHost on the RAW, VHD, VMDK, or qcow2 file. Once completed, an UCloud UHost Image is returned. The temporary RAW, VHD, VMDK, or qcow2 copy in UFile can be discarded after the import is complete.

## Configuration

There are some configuration options available for the post-processor. There
are two categories: required and optional parameters.

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


<!-- Code generated from the comments of the Config struct in post-processor/ucloud-import/post-processor.go; DO NOT EDIT MANUALLY -->

- `ufile_bucket_name` (string) - The name of the UFile bucket where the RAW, VHD, VMDK, or qcow2 file will be copied to for import.
   This bucket must exist when the post-processor is run.

- `image_name` (string) - The name of the user-defined image, which contains 1-63 characters and only
  supports Chinese, English, numbers, '-\_,.:[]'.

- `image_os_type` (string) - Type of the OS. Possible values are: `CentOS`, `Ubuntu`, `Windows`, `RedHat`, `Debian`, `Other`.
  You may refer to [ucloud_api_docs](https://docs.ucloud.cn/api/uhost-api/import_custom_image) for detail.

- `image_os_name` (string) - The name of OS. Such as: `CentOS 7.2 64位`, set `Other` When `image_os_type` is `Other`.
  You may refer to [ucloud_api_docs](https://docs.ucloud.cn/api/uhost-api/import_custom_image) for detail.

- `format` (string) - The format of the import image , Possible values are: `raw`, `vhd`, `vmdk`, or `qcow2`.

<!-- End of code generated from the comments of the Config struct in post-processor/ucloud-import/post-processor.go; -->


### Optional:

<!-- Code generated from the comments of the AccessConfig struct in builder/ucloud/common/access_config.go; DO NOT EDIT MANUALLY -->

- `base_url` (string) - This is the base url. (Default: `https://api.ucloud.cn`).

- `profile` (string) - This is the UCloud profile name as set in the shared credentials file, it can
  also be sourced from the `UCLOUD_PROFILE` environment variables.

- `shared_credentials_file` (string) - This is the path to the shared credentials file, it can also be sourced from
  the `UCLOUD_SHARED_CREDENTIAL_FILE` environment variables. If this is not set
  and a profile is specified, `~/.ucloud/credential.json` will be used.

<!-- End of code generated from the comments of the AccessConfig struct in builder/ucloud/common/access_config.go; -->


<!-- Code generated from the comments of the Config struct in post-processor/ucloud-import/post-processor.go; DO NOT EDIT MANUALLY -->

- `ufile_key_name` (string) - The name of the object key in
   `ufile_bucket_name` where the RAW, VHD, VMDK, or qcow2 file will be copied
   to import. This is a [template engine](/packer/docs/templates/legacy_json_templates/engine).
   Therefore, you may use user variables and template functions in this field.

- `skip_clean` (bool) - Whether we should skip removing the RAW, VHD, VMDK, or qcow2 file uploaded to
  UFile after the import process has completed. Possible values are: `true` to
  leave it in the UFile bucket, `false` to remove it. (Default: `false`).

- `image_description` (string) - The description of the image.

- `wait_image_ready_timeout` (int) - Timeout of importing image. The default timeout is 3600 seconds if this option is not set or is set.

<!-- End of code generated from the comments of the Config struct in post-processor/ucloud-import/post-processor.go; -->


## Basic Example

Here is a basic example. This assumes that the builder has produced a RAW artifact for us to work with. This will take the RAW image generated by a builder and upload it to UFile. Once uploaded, the import process will start, creating an UCloud UHost image to the region `cn-bj2`.

```json
"post-processors":[
    {
      "type":"ucloud-import",
      "public_key": "{{user `ucloud_public_key`}}",
      "private_key": "{{user `ucloud_private_key`}}",
      "project_id": "{{user `ucloud_project_id`}}",
      "region":"cn-bj2",
      "ufile_bucket_name": "packer-import",
      "image_name": "packer_import",
      "image_os_type": "CentOS",
      "image_os_name": "CentOS 6.10 64位",
      "format": "raw"
    }
  ]
```
