Provides a resource to create a WAF instance attack log post config

~> **NOTE:** Only enterprise version and above are supported for activation

Example Usage

```hcl
resource "tencentcloud_waf_instance_attack_log_post_config" "example" {
  instance_id     = "waf_2kxtlbky11b4wcrb"
  attack_log_post = 1
}
```

Import

WAF instance attack log post config can be imported using the id, e.g.

```
terraform import tencentcloud_waf_instance_attack_log_post_config.example waf_2kxtlbky11b4wcrb
```
