Provides a resource to create a VCube application and web player license

Example Usage

```hcl
resource "tencentcloud_vcube_application_and_web_player_license" "example" {
  app_name = "tf-example"
  domain_list = [
    "www.example1.com",
    "www.example2.com",
    "www.example3.com",
  ]
}
```

Import

VCube application and web player license can be imported using the id, e.g.

```
terraform import tencentcloud_vcube_application_and_web_player_license.example 1513
```
