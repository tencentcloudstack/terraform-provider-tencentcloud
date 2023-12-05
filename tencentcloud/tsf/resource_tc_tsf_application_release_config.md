Provides a resource to create a tsf application_release_config

Example Usage

```hcl
resource "tencentcloud_tsf_application_release_config" "application_release_config" {
  config_id = "dcfg-nalqbqwv"
  group_id = "group-yxmz72gv"
  release_desc = "terraform-test"
}
```

Import

tsf application_release_config can be imported using the configId#groupId#configReleaseId, e.g.

```
terraform import tencentcloud_tsf_application_release_config.application_release_config dcfg-nalqbqwv#group-yxmz72gv#dcfgr-maeeq2ea
```