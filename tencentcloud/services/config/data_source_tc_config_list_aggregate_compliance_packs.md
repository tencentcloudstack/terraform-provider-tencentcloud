Use this data source to query detailed information of TEO config aggregate compliance packs by account group.

Example Usage

```hcl
data "tencentcloud_config_list_aggregate_compliance_packs" "example" {
  account_group_id = "aggr-xxxxxxxx"
}

output "first_pack_name" {
  value = data.tencentcloud_config_list_aggregate_compliance_packs.example.items[0].compliance_pack_name
}
```

Query with filters:

```hcl
data "tencentcloud_config_list_aggregate_compliance_packs" "filtered" {
  account_group_id = "aggr-xxxxxxxx"
  compliance_pack_name = "CIS-CWP-Best-Practices"
  risk_level = [1, 2]
  status = "ACTIVE"
  compliance_result = ["NON_COMPLIANT"]
  order_type = "desc"
  include_compliance_pack_rule_result = "NO"
}
```
