---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_bucket"
sidebar_current: "docs-tencentcloud-resource-cos_bucket"
description: |-
  Provides a COS resource to create a COS bucket and set its attributes.
---

# tencentcloud_cos_bucket

Provides a COS resource to create a COS bucket and set its attributes.

## Example Usage

Private Bucket

```hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"
  acl    = "private"
}
```

Static Website

```hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"

  website = {
    index_document = "index.html"
    error_document = "error.html"
  }
}
```

Using CORS

```hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"
  acl    = "public-read-write"

  cors_rules {
    allowed_origins = ["http://*.abc.com"]
    allowed_methods = ["PUT", "POST"]
    allowed_headers = ["*"]
    max_age_seconds = 300
    expose_headers  = ["Etag"]
  }
}
```

Using object lifecycle

```hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"
  acl    = "public-read-write"

  lifecycle_rules {
    filter_prefix = "path1/"

    transition {
      date          = "2019-06-01"
      storage_class = "STANDARD_IA"
    }

    expiration {
      days = 90
    }
  }
}
```

Setting log status

```hcl
resource "tencentcloud_cam_role" "cosLogGrant" {
  name     = "CLS_QcsRole"
  document = "{\"version\":\"2.0\",\"statement\":[{\"action\":[\"name/sts:AssumeRole\"],\"effect\":\"allow\",\"principal\":{\"service\":[\"cls.cloud.tencent.com\"]}}]}"

  description = "cos log enable grant"
}

data "tencentcloud_cam_policies" "cosAccess" {
  name = "QcloudCOSAccessForCLSRole"
}

resource "tencentcloud_cam_role_policy_attachment" "cosLogGrant" {
  role_id   = tencentcloud_cam_role.cosLogGrant.id
  policy_id = data.tencentcloud_cam_policies.cosAccess.policy_list.0.policy_id
}

resource "tencentcloud_cos_bucket" "mylog" {
  bucket = "mylog-1258798060"
  acl    = "private"
}

resource "tencentcloud_cos_bucket" "mycos" {
  bucket            = "mycos-1258798060"
  acl               = "private"
  log_enable        = true
  log_target_bucket = "mylog-1258798060"
  log_prefix        = "MyLogPrefix"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, ForceNew) The name of a bucket to be created. Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.
* `acl` - (Optional) The canned ACL to apply. Valid values: private, public-read, and public-read-write. Defaults to private.
* `cors_rules` - (Optional) A rule of Cross-Origin Resource Sharing (documented below).
* `encryption_algorithm` - (Optional) The server-side encryption algorithm to use. Valid value is `AES256`.
* `lifecycle_rules` - (Optional) A configuration of object lifecycle management (documented below).
* `log_enable` - (Optional) Indicate the access log of this bucket to be saved or not. Default is `false`. If set `true`, the access log will be saved with `log_target_bucket`. To enable log, the full access of log service must be granted. [Full Access Role Policy](https://cloud.tencent.com/document/product/436/16920#.E5.90.AF.E7.94.A8.E6.97.A5.E5.BF.97.E7.AE.A1.E7.90.86).
* `log_prefix` - (Optional) The prefix log name which saves the access log of this bucket per 5 minutes. Eg. `MyLogPrefix/`. The log access file format is `log_target_bucket`/`log_prefix`{YYYY}/{MM}/{DD}/{time}_{random}_{index}.gz. Only valid when `log_enable` is `true`.
* `log_target_bucket` - (Optional) The target bucket name which saves the access log of this bucket per 5 minutes. The log access file format is `log_target_bucket`/`log_prefix`{YYYY}/{MM}/{DD}/{time}_{random}_{index}.gz. Only valid when `log_enable` is `true`. User must have full access on this bucket.
* `tags` - (Optional) The tags of a bucket.
* `versioning_enable` - (Optional) Enable bucket versioning.
* `website` - (Optional) A website object(documented below).

The `cors_rules` object supports the following:

* `allowed_headers` - (Required) Specifies which headers are allowed.
* `allowed_methods` - (Required) Specifies which methods are allowed. Can be `GET`, `PUT`, `POST`, `DELETE` or `HEAD`.
* `allowed_origins` - (Required) Specifies which origins are allowed.
* `expose_headers` - (Optional) Specifies expose header in the response.
* `max_age_seconds` - (Optional) Specifies time in seconds that browser can cache the response for a preflight request.

The `expiration` object supports the following:

* `date` - (Optional) Specifies the date after which you want the corresponding action to take effect.
* `days` - (Optional) Specifies the number of days after object creation when the specific rule action takes effect.

The `lifecycle_rules` object supports the following:

* `filter_prefix` - (Required) Object key prefix identifying one or more objects to which the rule applies.
* `expiration` - (Optional) Specifies a period in the object's expire (documented below).
* `transition` - (Optional) Specifies a period in the object's transitions (documented below).

The `transition` object supports the following:

* `storage_class` - (Required) Specifies the storage class to which you want the object to transition. Available values include `STANDARD`, `STANDARD_IA` and `ARCHIVE`.
* `date` - (Optional) Specifies the date after which you want the corresponding action to take effect.
* `days` - (Optional) Specifies the number of days after object creation when the specific rule action takes effect.

The `website` object supports the following:

* `error_document` - (Optional) An absolute path to the document to return in case of a 4XX error.
* `index_document` - (Optional) COS returns this index document when requests are made to the root domain or any of the subfolders.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cos_bucket_url` - The URL of this cos bucket.


## Import

COS bucket can be imported, e.g.

```
$ terraform import tencentcloud_cos_bucket.bucket bucket-name
```

