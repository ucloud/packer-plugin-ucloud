# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "UCloud"
  description = "The UCloud plugin is able to build customized images for use in UHost Instance."
  identifier = "packer/hashicorp/ucloud"
  component {
    type = "builder"
    name = "UCloud Uhost"
    slug = "uhost"
  }
  component {
    type = "post-processor"
    name = "UCloud Import"
    slug = "import"
  }
}
