Use this data source to query detailed information of CBS storages in parallel.

Example Usage

Query CBS by storage set by zone

```hcl
data "tencentcloud_cbs_storages_set" "example" {
  availability_zone = "ap-guangzhou-3"
}
```