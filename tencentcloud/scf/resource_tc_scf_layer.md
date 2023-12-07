Provide a resource to create a SCF layer.

Example Usage

```hcl
resource "tencentcloud_scf_layer" "foo" {
  layer_name = "foo"
  compatible_runtimes = ["Python3.6"]
  content {
    cos_bucket_name = "test-bucket"
    cos_object_name = "/foo.zip"
    cos_bucket_region = "ap-guangzhou"
  }
  description = "foo"
  license_info = "foo"
}
```
Import

Scf layer can be imported, e.g.

```
$ terraform import tencentcloud_scf_layer.layer layerId#layerVersion
```