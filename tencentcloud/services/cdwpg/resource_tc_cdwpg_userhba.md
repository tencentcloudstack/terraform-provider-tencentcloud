Provides a resource to create a cdwpg cdwpg_userhba

Example Usage

```hcl
resource "tencentcloud_cdwpg_userhba" "cdwpg_userhba" {
  instance_id = "cdwpg-zpiemnyd"
  hba_configs {
	type = "host"
	database = "all"
	user = "all"
	address = "0.0.0.0/0"
	method = "md5"
  }
}
```

Import

cdwpg cdwpg_userhba can be imported using the id, e.g.

```
terraform import tencentcloud_cdwpg_userhba.cdwpg_userhba cdwpg_userhba_id
```
