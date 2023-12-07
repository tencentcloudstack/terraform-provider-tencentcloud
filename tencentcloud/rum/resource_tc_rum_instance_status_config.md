Provides a resource to create a rum instance_status_config

Example Usage

```hcl
resource "tencentcloud_rum_instance_status_config" "instance_status_config" {
  instance_id = "rum-pasZKEI3RLgakj"
  operate     = "stop"
}
```

Import

rum instance_status_config can be imported using the id, e.g.

```
terraform import tencentcloud_rum_instance_status_config.instance_status_config instance_id
```