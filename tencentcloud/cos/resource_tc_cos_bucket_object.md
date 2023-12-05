Provides a COS object resource to put an object(content or file) to the bucket.

Example Usage

Uploading a file to a bucket

```hcl
resource "tencentcloud_cos_bucket_object" "myobject" {
  bucket = "mycos-1258798060"
  key    = "new_object_key"
  source = "path/to/file"
}
```

Uploading a content to a bucket

```hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"
  acl    = "public-read"
}

resource "tencentcloud_cos_bucket_object" "myobject" {
  bucket  = tencentcloud_cos_bucket.mycos.bucket
  key     = "new_object_key"
  content = "the content that you want to upload."
}
```