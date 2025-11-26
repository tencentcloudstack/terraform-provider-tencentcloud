Provides a resource to create a IGTM instance

~> **NOTE:** Currently, executing the `terraform destroy` command to delete this resource is not supported. If you need to destroy it, please contact Tencent Cloud IGTM through a ticket.

Example Usage

```hcl
resource "tencentcloud_igtm_instance" "example" {
  domain            = "domain.com"
  access_type       = "CUSTOM"
  global_ttl        = 60
  package_type      = "STANDARD"
  instance_name     = "tf-example"
  access_domain     = "domain.com"
  access_sub_domain = "sub_domain.com"
  remark            = "remark."
  resource_id       = "ins-lnpnnwvwxgs"
}
```

Import

IGTM instance can be imported using the id, e.g.

```
terraform import tencentcloud_igtm_instance.example gtm-uukztqtoaru
```
