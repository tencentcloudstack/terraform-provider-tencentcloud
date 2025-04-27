Provides a resource to create a WAF log post cls flow

Example Usage

If log_type is 1

```hcl
resource "tencentcloud_waf_log_post_cls_flow" "example" {
  cls_region     = "ap-guangzhou"
  logset_name    = "waf_post_logset"
  log_type       = 1
  log_topic_name = "waf_post_logtopic"
}
```

If log_type is 2

```hcl
resource "tencentcloud_waf_log_post_cls_flow" "example" {
  cls_region     = "ap-guangzhou"
  logset_name    = "waf_post_logset"
  log_type       = 2
  log_topic_name = "waf_post_logtopic"
}
```

Import

WAF log post cls flow can be imported using the id, e.g.

```
# If log_type is 1
terraform import tencentcloud_waf_log_post_cls_flow.example 111462#1

# If log_type is 2
terraform import tencentcloud_waf_log_post_cls_flow.example 111467#2
```
