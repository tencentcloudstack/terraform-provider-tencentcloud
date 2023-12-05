Provides a resource to create a waf anti_fake

~> **NOTE:** Uri: Please configure static resources such as. html,. shtml,. txt,. js,. css,. jpg,. png, or access paths for static resources..

Example Usage

```hcl
resource "tencentcloud_waf_anti_fake" "example" {
  domain = "www.waf.com"
  name   = "tf_example"
  uri    = "/anti_fake_url.html"
  status = 1
}
```

Import

waf anti_fake can be imported using the id, e.g.

```
terraform import tencentcloud_waf_anti_fake.example 3200035516#www.waf.com
```