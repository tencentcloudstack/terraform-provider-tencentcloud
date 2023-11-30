Use this data source to query detailed information of CLB attachments

Example Usage

```hcl
data "tencentcloud_clb_attachments" "clblab" {
  listener_id = "lbl-hh141sn9"
  clb_id      = "lb-k2zjp9lv"
  rule_id     = "loc-4xxr2cy7"
}
```