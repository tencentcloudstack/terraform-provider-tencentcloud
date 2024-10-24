Use this data source to query CCN route table input policies.

Example Usage

```hcl
data "tencentcloud_ccn_route_table_input_policies" "example" {
  ccn_id         = "ccn-06jek8tf"
  route_table_id = "ccnrtb-4jv5ltb9"
}
```