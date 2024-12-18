Provides a resource to create a dasb bind device resource

Example Usage

```hcl
resource "tencentcloud_dasb_bind_device_resource" "example" {
  resource_id   = "bh-saas-weyosfym"
  device_id_set = [17, 18]
}
```

Or custom domain_id parameters

```hcl
resource "tencentcloud_dasb_bind_device_resource" "example" {
  resource_id   = "bh-saas-lx1pxhli"
  domain_id     = "net-31nssj3n"
  device_id_set = [115, 116]
}
```