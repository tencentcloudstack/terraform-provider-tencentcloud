Use this data source to query detailed information of CBS storages in parallel.

Example Usage

```hcl
data "tencentcloud_cbs_storages_set" "storages" {
  availability_zone = "ap-guangzhou-3"
}
```