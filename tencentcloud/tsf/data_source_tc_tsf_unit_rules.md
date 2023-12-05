Use this data source to query detailed information of tsf unit_rules

Example Usage

```hcl
data "tencentcloud_tsf_unit_rules" "unit_rules" {
  gateway_instance_id = "gw-ins-lvdypq5k"
  status = "disabled"
}
```