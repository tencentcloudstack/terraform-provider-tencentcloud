Provides a resource to create a TencentCloud EdgeOne (TEO) shared CNAME

Example Usage

```hcl
resource "tencentcloud_teo_shared_cname" "example" {
  zone_id             = "zone-39quuimqg8r6"
  shared_cname_prefix = "test-api"
  description         = "example shared cname"
}
```

Import

TEO shared CNAME can be imported using the composite id (zone_id#shared_cname), e.g.

```
terraform import tencentcloud_teo_shared_cname.example zone-39quuimqg8r6#test-api.sai2ig51kaa5.share.dnse2.com
```
