Use this data source to query detailed information of TEO l7 acc rule

~> **NOTE:** This feature only supports sites in plans of Standard Edition and Enterprise Edition.

Example Usage

Query l7 acc rules by zone Id

```hcl
data "tencentcloud_teo_l7_acc_rule" "example" {
  zone_id = "zone-36bjhygh1bxe"
}
```

Query l7 acc rules with offset for pagination

```hcl
data "tencentcloud_teo_l7_acc_rule" "example_with_offset" {
  zone_id = "zone-36bjhygh1bxe"
  offset  = 10
}
```

Query l7 acc rules by rule Id

```hcl
data "tencentcloud_teo_l7_acc_rule" "example_by_rule_id" {
  zone_id = "zone-36bjhygh1bxe"
  rule_id = "rule-12345678"
}
```
