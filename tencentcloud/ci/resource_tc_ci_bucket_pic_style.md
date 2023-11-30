Provides a resource to create a ci bucket_pic_style

Example Usage

```hcl
resource "tencentcloud_ci_bucket_pic_style" "bucket_pic_style" {
  bucket     = "terraform-ci-xxxxxx"
  style_name = "rayscale_2"
  style_body = "imageMogr2/thumbnail/20x/crop/20x20/gravity/center/interlace/0/quality/100"
}
```

Import

ci bucket_pic_style can be imported using the bucket#styleName, e.g.

```
terraform import tencentcloud_ci_bucket_pic_style.bucket_pic_style terraform-ci-xxxxxx#rayscale_2
```