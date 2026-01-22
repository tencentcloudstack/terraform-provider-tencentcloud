Provides a resource to create a VCube application and video

Example Usage

```hcl
resource "tencentcloud_vcube_application_and_video" "example" {
  app_name  = "tf-example"
  bundle_id = "com.example.appName"
}
```

Or

```hcl
resource "tencentcloud_vcube_application_and_video" "example" {
  app_name     = "tf-example"
  package_name = "com.example.appName"
}
```

Or

```hcl
resource "tencentcloud_vcube_application_and_video" "example" {
  app_name     = "tf-example"
  bundle_id    = "com.example.appName"
  package_name = "com.example.appName"
}
```

Import

VCube application and video can be imported using the id, e.g.

```
terraform import tencentcloud_vcube_application_and_video.example 1509
```
