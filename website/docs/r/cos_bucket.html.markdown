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

Using verbose acl

```hcl
resource "tencentcloud_cos_bucket" "with_acl_body" {
  bucket   = "mycos-1258798060"
  acl_body = <<EOF
<AccessControlPolicy>
    <Owner>
        <ID>qcs::cam::uin/100000000001:uin/100000000001</ID>
    </Owner>
    <AccessControlList>
        <Grant>
            <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="Group">
                <URI>http://cam.qcloud.com/groups/global/AllUsers</URI>
            </Grantee>
            <Permission>READ</Permission>
        </Grant>
        <Grant>
            <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="CanonicalUser">
                <ID>qcs::cam::uin/100000000001:uin/100000000001</ID>
            </Grantee>
            <Permission>WRITE</Permission>
        </Grant>
        <Grant>
            <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="CanonicalUser">
                <ID>qcs::cam::uin/100000000001:uin/100000000001</ID>
            </Grantee>
            <Permission>READ_ACP</Permission>
        </Grant>
    </AccessControlList>
</AccessControlPolicy>
EOF
}
```

Static Website

```hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"

  website {
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

Using custom origin domain settings

```hcl
resource "tencentcloud_cos_bucket" "with_origin" {
  bucket = "mycos-1258798060"
  acl    = "private"
  origin_domain_rules {
    domain = "abc.example.com"
    type   = "REST"
    status = "ENABLE"
  }
}
```

Using origin-pull settings

```hcl
resource "tencentcloud_cos_bucket" "with_origin" {
  bucket = "mycos-1258798060"
  acl    = "private"
  origin_pull_rules {
    priority            = 1
    sync_back_to_source = false
    host                = "abc.example.com"
    prefix              = "/"
    protocol            = "FOLLOW" // "HTTP" "HTTPS"
    follow_query_string = true
    follow_redirection  = true
    follow_http_headers = ["origin", "host"]
    custom_http_headers = {
      "x-custom-header" = "custom_value"
    }
  }
}
```

Using replication

```hcl
resource "tencentcloud_cos_bucket" "replica1" {
  bucket            = "tf-replica-foo-1234567890"
  acl               = "private"
  versioning_enable = true
}

resource "tencentcloud_cos_bucket" "with_replication" {
  bucket            = "tf-bucket-replica-1234567890"
  acl               = "private"
  versioning_enable = true
  replica_role      = "qcs::cam::uin/100000000001:uin/100000000001"
  replica_rules {
    id                 = "test-rep1"
    status             = "Enabled"
    prefix             = "dist"
    destination_bucket = "qcs::cos:%s::${tencentcloud_cos_bucket.replica1.bucket}"
  }
}
```

Setting log status

```hcl
resource "tencentcloud_cam_role" "cosLogGrant" {
  name     = "CLS_QcsRole"
  document = <<EOF
{
  "version": "2.0",
  "statement": [
    {
      "action": [
        "name/sts:AssumeRole"
      ],
      "effect": "allow",
      "principal": {
        "service": [
          "cls.cloud.tencent.com"
        ]
      }
    }
  ]
}
EOF

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
* `acl_body` - (Optional) ACL XML body for multiple grant info.
* `acl` - (Optional) The canned ACL to apply. Valid values: private, public-read, and public-read-write. Defaults to private.
* `cors_rules` - (Optional) A rule of Cross-Origin Resource Sharing (documented below).
* `encryption_algorithm` - (Optional) The server-side encryption algorithm to use. Valid value is `AES256`.
* `lifecycle_rules` - (Optional) A configuration of object lifecycle management (documented below).
* `log_enable` - (Optional) Indicate the access log of this bucket to be saved or not. Default is `false`. If set `true`, the access log will be saved with `log_target_bucket`. To enable log, the full access of log service must be granted. [Full Access Role Policy](https://intl.cloud.tencent.com/document/product/436/16920).
* `log_prefix` - (Optional) The prefix log name which saves the access log of this bucket per 5 minutes. Eg. `MyLogPrefix/`. The log access file format is `log_target_bucket`/`log_prefix`{YYYY}/{MM}/{DD}/{time}_{random}_{index}.gz. Only valid when `log_enable` is `true`.
* `log_target_bucket` - (Optional) The target bucket name which saves the access log of this bucket per 5 minutes. The log access file format is `log_target_bucket`/`log_prefix`{YYYY}/{MM}/{DD}/{time}_{random}_{index}.gz. Only valid when `log_enable` is `true`. User must have full access on this bucket.
* `multi_az` - (Optional, ForceNew) Indicates whether to create a bucket of multi available zone.
* `origin_domain_rules` - (Optional) Bucket Origin Domain settings.
* `origin_pull_rules` - (Optional) Bucket Origin-Pull settings.
* `replica_role` - (Optional) Request initiator identifier, format: `qcs::cam::uin/<owneruin>:uin/<subuin>`. NOTE: only `versioning_enable` is true can configure this argument.
* `replica_rules` - (Optional) List of replica rule. NOTE: only `versioning_enable` is true and `replica_role` set can configure this argument.
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

The `origin_domain_rules` object supports the following:

* `domain` - (Required) Specify domain host.
* `status` - (Optional) Domain status, default: `ENABLED`.
* `type` - (Optional) Specify origin domain type, available values: `REST`, `WEBSITE`, `ACCELERATE`, default: `REST`.

The `origin_pull_rules` object supports the following:

* `host` - (Required) Allows only a domain name or IP address. You can optionally append a port number to the address.
* `priority` - (Required) Priority of origin-pull rules, do not set the same value for multiple rules.
* `custom_http_headers` - (Optional) Specifies the custom headers that you can add for COS to access your origin server.
* `follow_http_headers` - (Optional) Specifies the pass through headers when accessing the origin server.
* `follow_query_string` - (Optional) Specifies whether to pass through COS request query string when accessing the origin server.
* `follow_redirection` - (Optional) Specifies whether to follow 3XX redirect to another origin server to pull data from.
* `prefix` - (Optional) Triggers the origin-pull rule when the requested file name matches this prefix.
* `protocol` - (Optional) the protocol used for COS to access the specified origin server. The available value include `HTTP`, `HTTPS` and `FOLLOW`.
* `sync_back_to_source` - (Optional) If `true`, COS will not return 3XX status code when pulling data from an origin server. Current available zone: ap-beijing, ap-shanghai, ap-singapore, ap-mumbai.

The `replica_rules` object supports the following:

* `destination_bucket` - (Required) Destination bucket identifier, format: `qcs::cos:<region>::<bucketname-appid>`. NOTE: destination bucket must enable versioning.
* `status` - (Required) Status identifier, available values: `Enabled`, `Disabled`.
* `destination_storage_class` - (Optional) Storage class of destination, available values: `STANDARD`, `INTELLIGENT_TIERING`, `STANDARD_IA`. default is following current class of destination.
* `id` - (Optional) Name of a specific rule.
* `prefix` - (Optional) Prefix matching policy. Policies cannot overlap; otherwise, an error will be returned. To match the root directory, leave this parameter empty.

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

