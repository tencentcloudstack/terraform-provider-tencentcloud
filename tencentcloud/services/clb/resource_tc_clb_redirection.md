Provides a resource to create a CLB redirection.

Example Usage

Manual Rewrite

```hcl
resource "tencentcloud_clb_redirection" "foo" {
  clb_id             = "lb-p7olt9e5"
  source_listener_id = "lbl-jc1dx6ju"
  target_listener_id = "lbl-asj1hzuo"
  source_rule_id     = "loc-ft8fmngv"
  target_rule_id     = "loc-4xxr2cy7"
}
```

Auto Rewrite

```hcl
resource "tencentcloud_clb_redirection" "foo" {
  clb_id             = "lb-p7olt9e5"
  target_listener_id = "lbl-asj1hzuo"
  target_rule_id     = "loc-4xxr2cy7"
  is_auto_rewrite    = true
}
```

Import

CLB redirection can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_redirection.foo loc-ft8fmngv#loc-4xxr2cy7#lbl-jc1dx6ju#lbl-asj1hzuo#lb-p7olt9e5
```