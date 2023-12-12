Use this data source to query detailed information of tsf groups

Example Usage

```hcl
data "tencentcloud_tsf_groups" "groups" {
  search_word = "keep"
  application_id = "application-a24x29xv"
  order_by = "createTime"
  order_type = 0
  namespace_id = "namespace-aemrg36v"
  cluster_id = "cluster-vwgj5e6y"
  group_resource_type_list = ["DEF"]
  status = "Running"
  group_id_list = ["group-yrjkln9v"]
}
```