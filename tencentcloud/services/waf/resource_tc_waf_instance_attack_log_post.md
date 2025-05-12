Provides a resource to create a WAF instance attack log post

~> **NOTE:** Only enterprise version and above are supported for activation

Example Usage

```hcl
resource "tencentcloud_waf_instance_attack_log_post" "example" {
  instance_id     = "waf_2kxtlbky11b4wcrb"
  attack_log_post = 1
}
```

Import

WAF instance attack log post can be imported using the id, e.g.

```
terraform import tencentcloud_waf_instance_attack_log_post.example waf_2kxtlbky11b4wcrb
```
