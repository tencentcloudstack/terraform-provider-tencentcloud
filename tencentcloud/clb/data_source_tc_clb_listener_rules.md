Use this data source to query detailed information of CLB listener rule

Example Usage

```hcl
data "tencentcloud_clb_listener_rules" "foo" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-mwr6vbtv"
  rule_id     = "loc-inem40hz"
  domain      = "abc.com"
  url         = "/"
  scheduler   = "WRR"
}
```