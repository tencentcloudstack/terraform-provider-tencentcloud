Use this data source to query detailed information of cvm chc_hosts

Example Usage

```hcl
data "tencentcloud_cvm_chc_hosts" "chc_hosts" {
  chc_ids = ["chc-xxxxxx"]
  filters {
    name = "zone"
    values = ["ap-guangzhou-7"]
  }
}
```