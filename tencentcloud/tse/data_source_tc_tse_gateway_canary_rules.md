Use this data source to query detailed information of tse gateway_canary_rules

Example Usage

```hcl
data "tencentcloud_tse_gateway_canary_rules" "gateway_canary_rules" {
  gateway_id = "gateway-xxxxxx"
  service_id = "451a9920-e67a-4519-af41-fccac0e72005"
}
```