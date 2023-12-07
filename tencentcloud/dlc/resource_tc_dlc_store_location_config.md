Provides a resource to create a dlc store_location_config

Example Usage

```hcl
resource "tencentcloud_dlc_store_location_config" "store_location_config" {
  store_location = "cosn://bucketname/"
  enable = 1
}
```

Import

dlc store_location_config can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_store_location_config.store_location_config store_location_config_id
```