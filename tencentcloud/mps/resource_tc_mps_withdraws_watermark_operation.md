Provides a resource to create a mps withdraws_watermark_operation

Example Usage

Withdraw the watermark from COS

```hcl
resource "tencentcloud_cos_bucket" "example" {
  bucket = "tf-test-mps-wm-${local.app_id}"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket_object" "example" {
  bucket = tencentcloud_cos_bucket.example.bucket
  key    = "/test-file/test.mov"
  source = "/Users/luoyin/Downloads/file_example_MOV_480_700kB.mov"
}

resource "tencentcloud_mps_withdraws_watermark_operation" "operation" {
  input_info {
    type = "COS"
    cos_input_info {
      bucket = tencentcloud_cos_bucket_object.example.bucket
      region = "%s"
      object = tencentcloud_cos_bucket_object.example.key
    }
  }

  session_context = "this is a example session context"
}
```