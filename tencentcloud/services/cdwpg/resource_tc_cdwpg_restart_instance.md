Provides a resource to create a cdwpg cdwpg_restart_instance

Example Usage

```hcl
resource "tencentcloud_cdwpg_restart_instance" "cdwpg_restart_instance" {
	instance_id = "cdwpg-zpiemnyd"
	node_types = ["gtm"]
}
```
