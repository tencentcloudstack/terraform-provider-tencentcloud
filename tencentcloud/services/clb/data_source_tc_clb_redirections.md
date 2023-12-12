Use this data source to query detailed information of CLB redirections

Example Usage

```hcl
data "tencentcloud_clb_redirections" "foo" {
  clb_id             = "lb-p7olt9e5"
  source_listener_id = "lbl-jc1dx6ju"
  target_listener_id = "lbl-asj1hzuo"
  source_rule_id     = "loc-ft8fmngv"
  target_rule_id     = "loc-4xxr2cy7"
  result_output_file = "mytestpath"
}
```