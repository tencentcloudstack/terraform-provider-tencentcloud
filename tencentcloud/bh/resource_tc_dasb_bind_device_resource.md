Provides a resource to create a dasb bind_device_resource

Example Usage

```hcl
resource "tencentcloud_dasb_bind_device_resource" "example" {
  resource_id   = "bh-saas-ocmzo6lgxiv"
  device_id_set = [17, 18]
}
```