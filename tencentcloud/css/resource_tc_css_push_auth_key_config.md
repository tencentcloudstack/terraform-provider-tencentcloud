Provides a resource to create a css push_auth_key_config

Example Usage

```hcl
resource "tencentcloud_css_push_auth_key_config" "push_auth_key_config" {
  domain_name = "your_push_domain_name"
  enable = 1
  master_auth_key = "testmasterkey"
  backup_auth_key = "testbackkey"
  auth_delta = 1800
}
```

Import

css push_auth_key_config can be imported using the id, e.g.

```
terraform import tencentcloud_css_push_auth_key_config.push_auth_key_config push_auth_key_config_id
```