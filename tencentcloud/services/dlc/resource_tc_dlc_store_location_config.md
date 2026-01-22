Provides a resource to create a DLC store location config

Example Usage

Select user-defined COS path storage

```hcl
resource "tencentcloud_dlc_store_location_config" "example" {
  store_location = "cosn://tf-example-1308135196/demo"
  enable         = 1
}
```

Select DLC internal storage

```hcl
resource "tencentcloud_dlc_store_location_config" "example" {
  store_location = ""
  enable         = 0
}
```
