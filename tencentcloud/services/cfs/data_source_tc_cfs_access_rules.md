Use this data source to query the detail information of CFS access rule.

Example Usage

```hcl
data "tencentcloud_cfs_access_rules" "access_rules" {
  access_group_id = "pgroup-7nx89k7l"
  access_rule_id  = "rule-qcndbqzj"
}
```