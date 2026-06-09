Provides a resource to create a CLB redirection.

Example Usage

Manual Rewrite

```hcl
resource "tencentcloud_clb_redirection" "example" {
  clb_id             = "lb-ab09jtd2"
  source_listener_id = "lbl-qgtfowas"
  target_listener_id = "lbl-lpwdkukk"
  source_rule_id     = "loc-liz99mtg"
  target_rule_id     = "loc-4f53xn52"
  rewrite_code       = 307
  take_url           = true
  source_domian      = "www.demo.com"
}
```

Auto Rewrite

```hcl
resource "tencentcloud_clb_redirection" "example" {
  clb_id             = "lb-ab09jtd2"
  target_listener_id = "lbl-l7550kum"
  target_rule_id     = "loc-op7uz010"
  is_auto_rewrite    = true
}
```

Import

CLB redirection can be imported using the sourceLocId#targetLocId#sourceListenerId#targetListenerId#clbId, e.g.

```
terraform import tencentcloud_clb_redirection.example loc-ft8fmngv#loc-4xxr2cy7#lbl-jc1dx6ju#lbl-asj1hzuo#lb-p7olt9e5
```
