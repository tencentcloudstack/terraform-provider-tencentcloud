Provides a resource to create a dcdb instance_config

Example Usage

```hcl
resource "tencentcloud_dcdb_instance_config" "instance_config" {
  instance_id = local.dcdb_id
  rs_access_strategy = 0
}
```

Import

dcdb instance_config can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_instance_config.instance_config instance_config_id
```