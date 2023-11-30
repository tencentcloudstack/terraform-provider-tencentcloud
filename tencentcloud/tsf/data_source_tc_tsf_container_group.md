Use this data source to query detailed information of tsf container_group

Example Usage

```hcl
data "tencentcloud_tsf_container_group" "container_group" {
  application_id = "application-a24x29xv"
  search_word = "keep"
  order_by = "createTime"
  order_type = 0
  cluster_id = "cluster-vwgj5e6y"
  namespace_id = "namespace-aemrg36v"
}
```