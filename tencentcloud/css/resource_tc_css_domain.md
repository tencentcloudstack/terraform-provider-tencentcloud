Provides a resource to create a css domain

Example Usage

```hcl
resource "tencentcloud_css_domain" "domain" {
  domain_name = "iac-tf.cloud"
  domain_type = 0
  play_type = 1
  is_delay_live = 0
  is_mini_program_live = 0
  verify_owner_type = "dbCheck"
}
```

Import

css domain can be imported using the id, e.g.

```
terraform import tencentcloud_css_domain.domain domain_name
```