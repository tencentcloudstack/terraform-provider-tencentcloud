Provides a resource to manage TEO (EdgeOne) acceleration domain shared CNAME binding attachment

Example Usage

```hcl
resource "tencentcloud_teo_domain_shared_cname_attachment" "example" {
  zone_id = "zone-2qtuhspy7cr6"

  bind_shared_cname_maps {
    shared_cname = "shared.example.com"
    domain_names = ["domain1.example.com", "domain2.example.com"]
  }
}
```
Import

teo domain_shared_cname_attachment can be imported using the zone_id#shared_cname#domain_names (domain names joined by comma), e.g.
```
terraform import tencentcloud_teo_domain_shared_cname_attachment.example zone-2qtuhspy7cr6#shared.example.com#domain1.example.com,domain2.example.com
```