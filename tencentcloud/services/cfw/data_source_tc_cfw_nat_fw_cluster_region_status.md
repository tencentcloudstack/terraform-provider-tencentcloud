Use this data source to query detailed information of CFW NAT firewall cluster region status

Example Usage

```hcl
data "tencentcloud_cfw_nat_fw_cluster_region_status" "example" {
  nat_cluster_region_status_query_list {
    ccn_id       = "ccn-p3mlp0tj"
    nat_ins_id   = "nat-h1i1mf4n"
    asset_type   = "nat_ccn"
    routing_mode = 0
  }
}
```
