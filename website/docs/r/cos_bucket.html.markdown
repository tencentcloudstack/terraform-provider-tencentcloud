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

~> **NOTE:** The following capabilities do not support cdc scenarios: `multi_az`, `website`, and bucket replication `replica_role`.

~> **NOTE:** If `chdfs_ofs` is `true`, cannot set `acl_body`, `acl`, `origin_pull_rules`, `origin_domain_rules`, `website`, `encryption_algorithm`, `kms_id`, `versioning_enable`, `acceleration_enable` at the same time. For more information, please refer to `https://www.tencentcloud.com/document/product/436/43305`.

## Example Usage

### Private Bucket

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "private_bucket" {
  bucket = "private-bucket-${local.app_id}"
  acl    = "private"
}
```

### Private Bucket with CDC cluster

```hcl
provider "tencentcloud" {
  cos_domain = "https://${local.cdc_id}.cos-cdc.${local.region}.myqcloud.com/"
  region     = local.region
}

locals {
  region = "ap-guangzhou"
  cdc_id = "cluster-262n63e8"
}

data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "private_bucket" {
  bucket            = "private-bucket-${local.app_id}"
  acl               = "private"
  versioning_enable = true
  force_clean       = true
}
```

### Enable SSE-KMS encryption

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_kms_key" "example" {
  alias                = "tf-example-kms-key"
  description          = "example of kms key"
  key_rotation_enabled = false
  is_enabled           = true

  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_cos_bucket" "bucket_basic" {
  bucket               = "tf-bucket-cdc-${local.app_id}"
  acl                  = "private"
  encryption_algorithm = "KMS"
  kms_id               = tencentcloud_kms_key.example.id
  versioning_enable    = true
  acceleration_enable  = false
  force_clean          = true
}
```

### Creation of multiple available zone bucket

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "multi_zone_bucket" {
  bucket            = "multi-zone-bucket-${local.app_id}"
  acl               = "private"
  multi_az          = true
  versioning_enable = true
  force_clean       = true
}
```

### Using verbose acl

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "bucket_with_acl" {
  bucket = "bucketwith-acl-${local.app_id}"
  # NOTE: Specify the acl_body by the priority sequence of permission and user type with the following sequence: `CanonicalUser with READ`, `CanonicalUser with WRITE`, `CanonicalUser with FULL_CONTROL`, `CanonicalUser with WRITE_ACP`, `CanonicalUser with READ_ACP`, then specify the `Group` of permissions same as `CanonicalUser`.
  acl_body = <<EOF
<AccessControlPolicy>
	<Owner>
		<ID>qcs::cam::uin/100022975249:uin/100022975249</ID>
		<DisplayName>qcs::cam::uin/100022975249:uin/100022975249</DisplayName>
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
				<ID>qcs::cam::uin/100022975249:uin/100022975249</ID>
				<DisplayName>qcs::cam::uin/100022975249:uin/100022975249</DisplayName>
			</Grantee>
			<Permission>FULL_CONTROL</Permission>
		</Grant>
		<Grant>
			<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="CanonicalUser">
				<ID>qcs::cam::uin/100022975249:uin/100022975249</ID>
				<DisplayName>qcs::cam::uin/100022975249:uin/100022975249</DisplayName>
			</Grantee>
			<Permission>WRITE_ACP</Permission>
		</Grant>
		<Grant>
			<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="Group">
				<URI>http://cam.qcloud.com/groups/global/AllUsers</URI>
			</Grantee>
			<Permission>READ_ACP</Permission>
		</Grant>
		<Grant>
			<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="Group">
				<URI>http://cam.qcloud.com/groups/global/AllUsers</URI>
			</Grantee>
			<Permission>WRITE_ACP</Permission>
		</Grant>
		<Grant>
			<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="CanonicalUser">
				<ID>qcs::cam::uin/100022975249:uin/100022975249</ID>
				<DisplayName>qcs::cam::uin/100022975249:uin/100022975249</DisplayName>
			</Grantee>
			<Permission>READ</Permission>
		</Grant>
		<Grant>
			<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="CanonicalUser">
				<ID>qcs::cam::uin/100022975249:uin/100022975249</ID>
				<DisplayName>qcs::cam::uin/100022975249:uin/100022975249</DisplayName>
			</Grantee>
			<Permission>WRITE</Permission>
		</Grant>
		<Grant>
			<Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="Group">
				<URI>http://cam.qcloud.com/groups/global/AllUsers</URI>
			</Grantee>
			<Permission>FULL_CONTROL</Permission>
		</Grant>
	</AccessControlList>
</AccessControlPolicy>
EOF
}
```

### Using verbose acl with CDC cluster

```hcl
provider "tencentcloud" {
  cos_domain = "https://${local.cdc_id}.cos-cdc.${local.region}.myqcloud.com/"
  region     = local.region
}

locals {
  region = "ap-guangzhou"
  cdc_id = "cluster-262n63e8"
}

data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "bucket_with_acl" {
  bucket   = "private-bucket-${local.app_id}"
  acl      = "private"
  acl_body = <<EOF
<AccessControlPolicy>
    <Owner>
        <ID>qcs::cam::uin/100023201586:uin/100023201586</ID>
        <DisplayName>qcs::cam::uin/100023201586:uin/100023201586</DisplayName>
    </Owner>
    <AccessControlList>
        <Grant>
            <Grantee type="CanonicalUser">
                <ID>qcs::cam::uin/100015006748:uin/100015006748</ID>
                <DisplayName>qcs::cam::uin/100015006748:uin/100015006748</DisplayName>
            </Grantee>
            <Permission>WRITE</Permission>
        </Grant>
        <Grant>
            <Grantee type="CanonicalUser">
                <ID>qcs::cam::uin/100023201586:uin/100023201586</ID>
                <DisplayName>qcs::cam::uin/100023201586:uin/100023201586</DisplayName>
            </Grantee>
            <Permission>FULL_CONTROL</Permission>
        </Grant>
    </AccessControlList>
</AccessControlPolicy>
EOF
}
```

### Static Website

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "bucket_with_static_website" {
  bucket = "bucket-with-static-website-${local.app_id}"

  website {
    index_document           = "index.html"
    error_document           = "error.html"
    redirect_all_requests_to = "https"
    routing_rules {
      rules {
        condition_error_code        = "404"
        redirect_protocol           = "https"
        redirect_replace_key_prefix = "/test"
      }

      rules {
        condition_prefix     = "/test"
        redirect_protocol    = "https"
        redirect_replace_key = "key"
      }
    }
  }
}

output "endpoint_test" {
  value = tencentcloud_cos_bucket.bucket_with_static_website.website.0.endpoint
}
```

### Using CORS

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "bucket_with_cors" {
  bucket = "bucket-with-cors-${local.app_id}"
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

### Using Origin pull

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id    = data.tencentcloud_user_info.info.app_id
  uin       = data.tencentcloud_user_info.info.uin
  owner_uin = data.tencentcloud_user_info.info.owner_uin
}

resource "tencentcloud_cos_bucket" "example" {
  bucket = "tf-bucket-basic10-${local.app_id}"
  acl    = "public-read"
  origin_pull_rules {
    priority            = 1
    back_to_source_mode = "Redirect"
    http_redirect_code  = "301"
    protocol            = "FOLLOW"
    host                = "1.1.1.1"
    follow_query_string = true
  }
}
```

### Using CORS with CDC

```hcl
provider "tencentcloud" {
  cos_domain = "https://${local.cdc_id}.cos-cdc.${local.region}.myqcloud.com/"
  region     = local.region
}

locals {
  region = "ap-guangzhou"
  cdc_id = "cluster-262n63e8"
}

data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "bucket_with_cors" {
  bucket = "bucket-with-cors-${local.app_id}"

  cors_rules {
    allowed_origins = ["http://*.abc.com"]
    allowed_methods = ["PUT", "POST"]
    allowed_headers = ["*"]
    max_age_seconds = 300
    expose_headers  = ["Etag"]
  }
}
```

### Using object lifecycle

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "bucket_with_lifecycle" {
  bucket = "bucket-with-lifecycle-${local.app_id}"
  acl    = "public-read-write"

  lifecycle_rules {
    filter_prefix = "path1/"

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    expiration {
      days = 90
    }
  }
}
```

### Using object lifecycle with CDC

```hcl
provider "tencentcloud" {
  cos_domain = "https://${local.cdc_id}.cos-cdc.${local.region}.myqcloud.com/"
  region     = local.region
}

locals {
  region = "ap-guangzhou"
  cdc_id = "cluster-262n63e8"
}

data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "bucket_with_lifecycle" {
  bucket = "bucket-with-lifecycle-${local.app_id}"
  acl    = "private"

  lifecycle_rules {
    filter_prefix = "path1/"

    expiration {
      days = 90
    }
  }
}
```

### Using replication

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id    = data.tencentcloud_user_info.info.app_id
  uin       = data.tencentcloud_user_info.info.uin
  owner_uin = data.tencentcloud_user_info.info.owner_uin
  region    = "ap-guangzhou"
}

resource "tencentcloud_cos_bucket" "bucket_replicate" {
  bucket            = "bucket-replicate-${local.app_id}"
  acl               = "private"
  versioning_enable = true
}

resource "tencentcloud_cos_bucket" "bucket_with_replication" {
  bucket            = "bucket-with-replication-${local.app_id}"
  acl               = "private"
  versioning_enable = true
  replica_role      = "qcs::cam::uin/${local.owner_uin}:uin/${local.uin}"
  replica_rules {
    id       = "rule1"
    status   = "Enabled"
    priority = 1
    prefix   = "/prefix"
    filter {
      and {
        tag {
          key   = "tagKey1"
          value = "tagValue1"
        }

        tag {
          key   = "tagKey2"
          value = "tagValue2"
        }
      }
    }
    destination_bucket                = "qcs::cos:${local.region}::${tencentcloud_cos_bucket.bucket_replicate.bucket}"
    destination_storage_class         = "Standard"
    destination_encryption_kms_key_id = "4f14a617-7c7d-11ef-9a62-525400d3a886"
    delete_marker_replication {
      status = "Disabled"
    }

    source_selection_criteria {
      sse_kms_encrypted_objects {
        status = "Enabled"
      }
    }
  }

  replica_rules {
    id                                = "rule2"
    status                            = "Enabled"
    priority                          = 2
    destination_bucket                = "qcs::cos:${local.region}::${tencentcloud_cos_bucket.bucket_replicate.bucket}"
    destination_storage_class         = "Standard"
    destination_encryption_kms_key_id = "4f14a617-7c7d-11ef-9a62-525400d3a886"
    delete_marker_replication {
      status = "Enabled"
    }

    source_selection_criteria {
      sse_kms_encrypted_objects {
        status = "Enabled"
      }
    }
  }
}
```

### Using intelligent tiering, Only enable intelligent tiering

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "example" {
  bucket                               = "bucket-intelligent-tiering-${local.app_id}"
  acl                                  = "private"
  enable_intelligent_tiering           = true
  intelligent_tiering_days             = 30
  intelligent_tiering_request_frequent = 1
}
```

### Using intelligent tiering and configure the intelligent tiered storage archiving and deep archiving rules list.

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "example" {
  bucket                               = "bucket-intelligent-tiering-${local.app_id}"
  acl                                  = "private"
  enable_intelligent_tiering           = true
  intelligent_tiering_days             = 30
  intelligent_tiering_request_frequent = 1
  intelligent_tiering_archiving_rule_list {
    rule_id = "rule1"
    status  = "Enabled"

    tiering {
      access_tier = "ARCHIVE_ACCESS"
      days        = 91
    }

    tiering {
      access_tier = "DEEP_ARCHIVE_ACCESS"
      days        = 180
    }
  }

  intelligent_tiering_archiving_rule_list {
    rule_id = "rule2"
    status  = "Enabled"

    filter {
      prefix = "/prefix"
      tag {
        key   = "tagKey"
        value = "tagValue"
      }
    }

    tiering {
      access_tier = "ARCHIVE_ACCESS"
      days        = 91
    }
  }
}
```

### Using object lock config

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "example" {
  bucket = "bucket-intelligent-tiering-${local.app_id}"
  acl    = "private"
  object_lock_configuration {
    enabled = true
    rule {
      days = 30
    }
  }
}
```

### Using OFS

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "example" {
  bucket    = "private-ofs-bucket-${local.app_id}"
  acl       = "private"
  chdfs_ofs = true
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String, ForceNew) The name of a bucket to be created. Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.
* `acceleration_enable` - (Optional, Bool) Enable bucket acceleration.
* `acl_body` - (Optional, String) ACL XML body for multiple grant info. NOTE: this argument will overwrite `acl`. Check https://intl.cloud.tencent.com/document/product/436/7737 for more detail.
* `acl` - (Optional, String) The canned ACL to apply. Valid values: private, public-read, and public-read-write. Defaults to private.
* `cdc_id` - (Optional, String, ForceNew) CDC cluster ID.
* `chdfs_ofs` - (Optional, Bool, ForceNew) Indicates whether to create a bucket of metadata acceleration. For more information, please refer to `https://www.tencentcloud.com/document/product/436/43305`.
* `cors_rules` - (Optional, List) A rule of Cross-Origin Resource Sharing (documented below).
* `enable_intelligent_tiering` - (Optional, Bool) Enable intelligent tiering. NOTE: When intelligent tiering configuration is enabled, it cannot be turned off or modified.
* `encryption_algorithm` - (Optional, String) The server-side encryption algorithm to use. Valid values are `AES256`, `KMS` and `SM4`.
* `force_clean` - (Optional, Bool) Force cleanup all objects before delete bucket.
* `intelligent_tiering_archiving_rule_list` - (Optional, List) List of intelligent tiered storage, archiving, and deep archiving rules. NOTE: only `enable_intelligent_tiering` is true can configure this argument.
* `intelligent_tiering_days` - (Optional, Int) Specifies the limit of days for standard-tier data to low-frequency data in an intelligent tiered storage configuration, with optional days of 30, 60, 90. Default value is 30.
* `intelligent_tiering_request_frequent` - (Optional, Int) Specify the access limit for converting standard layer data into low-frequency layer data in the configuration. The default value is once, which can be used in combination with the number of days to achieve the conversion effect. For example, if the parameter is set to 1 and the number of access days is 30, it means that objects with less than one visit in 30 consecutive days will be reduced from the standard layer to the low frequency layer.
* `kms_id` - (Optional, String) The KMS Master Key ID. This value is valid only when `encryption_algorithm` is set to KMS. Set kms id to the specified value. If not specified, the default kms id is used.
* `lifecycle_rules` - (Optional, List) A configuration of object lifecycle management (documented below).
* `log_enable` - (Optional, Bool) Indicate the access log of this bucket to be saved or not. Default is `false`. If set `true`, the access log will be saved with `log_target_bucket`. To enable log, the full access of log service must be granted. [Full Access Role Policy](https://intl.cloud.tencent.com/document/product/436/16920).
* `log_prefix` - (Optional, String) The prefix log name which saves the access log of this bucket per 5 minutes. Eg. `MyLogPrefix/`. The log access file format is `log_target_bucket`/`log_prefix`{YYYY}/{MM}/{DD}/{time}_{random}_{index}.gz. Only valid when `log_enable` is `true`.
* `log_target_bucket` - (Optional, String) The target bucket name which saves the access log of this bucket per 5 minutes. The log access file format is `log_target_bucket`/`log_prefix`{YYYY}/{MM}/{DD}/{time}_{random}_{index}.gz. Only valid when `log_enable` is `true`. User must have full access on this bucket.
* `multi_az` - (Optional, Bool, ForceNew) Indicates whether to create a bucket of multi available zone.
* `object_lock_configuration` - (Optional, List) Object locking configuration. Once enabled, this feature cannot be disabled.
* `origin_domain_rules` - (Optional, List) Bucket Origin Domain settings.
* `origin_pull_rules` - (Optional, List) Bucket Origin-Pull settings.
* `replica_role` - (Optional, String) Request initiator identifier, format: `qcs::cam::uin/<owneruin>:uin/<subuin>`. NOTE: only `versioning_enable` is true can configure this argument.
* `replica_rules` - (Optional, List) List of replica rule. NOTE: only `versioning_enable` is true and `replica_role` set can configure this argument.
* `tags` - (Optional, Map) The tags of a bucket.
* `versioning_enable` - (Optional, Bool) Enable bucket versioning. NOTE: The `multi_az` feature is true for the current bucket, cannot disable version control.
* `website` - (Optional, List) A website object(documented below).

The `abort_incomplete_multipart_upload` object of `lifecycle_rules` supports the following:

* `days_after_initiation` - (Required, Int) Specifies the number of days after the multipart upload starts that the upload must be completed. The maximum value is 3650.

The `and` object of `filter` supports the following:

* `prefix` - (Optional, String) Filter objects by prefix; you can specify at most one prefix.
* `tag` - (Optional, List) When filtering objects for analysis, you can use object tags (multiple tags are supported) as filtering criteria.

The `and` object of `filter` supports the following:

* `prefix` - (Optional, String) Filter objects by prefix; you can specify at most one prefix.
* `tag` - (Optional, List) When filtering objects to be copied, you can use object tags (multiple tags are supported) as filtering criteria, with a maximum of 10 tags allowed. After adding tags as filtering criteria, the `delete_marker_replication.status` option must be set to false.

The `cors_rules` object supports the following:

* `allowed_headers` - (Required, List) Specifies which headers are allowed.
* `allowed_methods` - (Required, List) Specifies which methods are allowed. Can be `GET`, `PUT`, `POST`, `DELETE` or `HEAD`.
* `allowed_origins` - (Required, List) Specifies which origins are allowed.
* `expose_headers` - (Optional, List) Specifies expose header in the response.
* `max_age_seconds` - (Optional, Int) Specifies time in seconds that browser can cache the response for a preflight request.

The `delete_marker_replication` object of `replica_rules` supports the following:

* `status` - (Optional, String) Whether to synchronously delete the tag, supports Disabled or Enabled. The default value is Enabled, meaning the tag will be deleted synchronously.

The `expiration` object of `lifecycle_rules` supports the following:

* `date` - (Optional, String) Specifies the date after which you want the corresponding action to take effect.
* `days` - (Optional, Int) Specifies the number of days after object creation when the specific rule action takes effect.
* `delete_marker` - (Optional, Bool) Indicates whether the delete marker of an expired object will be removed.

The `filter` object of `intelligent_tiering_archiving_rule_list` supports the following:

* `and` - (Optional, List) For filtering conditions, if both prefix and object tag conditions are required simultaneously, they need to be wrapped with an `And` operator.
* `prefix` - (Optional, String) Filter objects by prefix; you can specify at most one prefix.
* `tag` - (Optional, List) When filtering objects for analysis, you can use object tags (multiple tags are supported) as filtering criteria.

The `filter` object of `replica_rules` supports the following:

* `and` - (Optional, List) When filtering objects to be copied, if both prefix and object tag conditions are required simultaneously, or if multiple object tag conditions are needed, they must be enclosed in an `And` statement.
* `prefix` - (Optional, String) Filter objects by prefix; you can specify at most one prefix.

The `intelligent_tiering_archiving_rule_list` object supports the following:

* `rule_id` - (Required, String) The name of the intelligent tiering rule name list task, with the ID set to a non-default string, indicates that this rule is a conversion rule for archive and deep archive tiers.
* `status` - (Required, String) Indicates whether the intelligent tiering rule is enabled. Possible values: Enabled, Disabled. When the ID is `default`, only `Enabled` is supported.
* `tiering` - (Required, List) Specifies configuration information related to data transformation in the intelligent tiered storage configuration.
* `filter` - (Optional, List) Specifies configuration information related to data transformation in the intelligent tiered storage configuration.

The `lifecycle_rules` object supports the following:

* `abort_incomplete_multipart_upload` - (Optional, Set) Set the maximum time a multipart upload is allowed to remain running.
* `expiration` - (Optional, Set) Specifies a period in the object's expire (documented below).
* `filter_prefix` - (Optional, String) Object key prefix identifying one or more objects to which the rule applies.
* `id` - (Optional, String) A unique identifier for the rule. It can be up to 255 characters.
* `non_current_expiration` - (Optional, Set) Specifies when non current object versions shall expire.
* `non_current_transition` - (Optional, Set) Specifies a period in the non current object's transitions.
* `transition` - (Optional, Set) Specifies a period in the object's transitions (documented below).

The `non_current_expiration` object of `lifecycle_rules` supports the following:

* `non_current_days` - (Optional, Int) Number of days after non current object creation when the specific rule action takes effect. The maximum value is 3650.

The `non_current_transition` object of `lifecycle_rules` supports the following:

* `storage_class` - (Required, String) Specifies the storage class to which you want the non current object to transition. Available values include `STANDARD_IA`, `MAZ_STANDARD_IA`, `INTELLIGENT_TIERING`, `MAZ_INTELLIGENT_TIERING`, `ARCHIVE`, `DEEP_ARCHIVE`. For more information, please refer to: https://cloud.tencent.com/document/product/436/33417.
* `non_current_days` - (Optional, Int) Number of days after non current object creation when the specific rule action takes effect.

The `object_lock_configuration` object supports the following:

* `enabled` - (Required, Bool) Enable object lock configuration.
* `rule` - (Optional, List) Object locking configuration.

The `origin_domain_rules` object supports the following:

* `domain` - (Required, String) Specify domain host.
* `status` - (Optional, String) Domain status, default: `ENABLED`.
* `type` - (Optional, String) Specify origin domain type, available values: `REST`, `WEBSITE`, `ACCELERATE`, default: `REST`.

The `origin_pull_rules` object supports the following:

* `host` - (Required, String) Allows only a domain name or IP address. You can optionally append a port number to the address.
* `priority` - (Required, Int) Priority of origin-pull rules, do not set the same value for multiple rules.
* `back_to_source_mode` - (Optional, String) Back to source mode. Allow value: Proxy, Mirror, Redirect.
* `custom_http_headers` - (Optional, Map) Specifies the custom headers that you can add for COS to access your origin server.
* `follow_http_headers` - (Optional, Set) Specifies the pass through headers when accessing the origin server.
* `follow_query_string` - (Optional, Bool) Specifies whether to pass through COS request query string when accessing the origin server.
* `follow_redirection` - (Optional, Bool) Specifies whether to follow 3XX redirect to another origin server to pull data from.
* `http_redirect_code` - (Optional, String) Redirect code. Effective when `back_to_source_mode` is `Redirect`. ex: 301, 302, 307. Default is 302.
* `prefix` - (Optional, String) Triggers the origin-pull rule when the requested file name matches this prefix.
* `protocol` - (Optional, String) the protocol used for COS to access the specified origin server. The available value include `HTTP`, `HTTPS` and `FOLLOW`.
* `sync_back_to_source` - (Optional, Bool, **Deprecated**) It has been deprecated from version 1.81.196. Please use `back_to_source_mode` instead. If `true`, COS will not return 3XX status code when pulling data from an origin server. Current available zone: ap-beijing, ap-shanghai, ap-singapore, ap-mumbai.

The `replica_rules` object supports the following:

* `destination_bucket` - (Required, String) Destination bucket identifier, format: `qcs::cos:<region>::<bucketname-appid>`. NOTE: destination bucket must enable versioning.
* `delete_marker_replication` - (Optional, List) Synchronized deletion marker.
* `destination_encryption_kms_key_id` - (Optional, String) This field must be included when `source_selection_criteria.sse_kms_encrypted_objects.status` is set to Enabled. It is used to specify the KMS key used for KMS-encrypted objects copied to the destination bucket.
* `destination_storage_class` - (Optional, String) Storage class of destination, available values: `Standard`, `Intelligent_Tiering`, `Standard_IA`. default is following current class of destination.
* `filter` - (Optional, List) Filter the objects to be copied. The bucket feature will copy objects that match the prefixes and tags specified in the Filter settings.
* `id` - (Optional, String) Name of a specific rule.
* `prefix` - (Optional, String) Prefix matching policy. Policies cannot overlap; otherwise, an error will be returned. To match the root directory, leave this parameter empty.
* `priority` - (Optional, Int) Execution priority, used to handle scenarios where the target storage buckets are the same and multiple replication rules match the same object. Note: Supports setting positive integers in the range of 1-1000. The Priority values of different rules cannot be duplicated. Storage bucket replication rules must either all have Priority set or all not have Priority set. When all rules have Priority set, overlapping prefixes are allowed for different rules when the target storage buckets are the same. When different rules match the same object, the rule with the smallest Priority value will be triggered first. When none of the rules have Priority set, overlapping prefixes are not allowed for different rules.
* `source_selection_criteria` - (Optional, List) This is used to specify additional conditions for objects supported by bucket replication rules. Currently, only the option to replicate KMS-encrypted objects is supported.
* `status` - (Optional, String) Status identifier, available values: `Enabled`, `Disabled`.

The `routing_rules` object of `website` supports the following:

* `rules` - (Required, List) Routing rule list.

The `rule` object of `object_lock_configuration` supports the following:

* `days` - (Required, Int) Object lock default duration (range: 1-36500).

The `rules` object of `routing_rules` supports the following:

* `condition_error_code` - (Optional, String) Specifies the error code as the match condition for the routing rule. Valid values: only 4xx return codes, such as 403 or 404.
* `condition_prefix` - (Optional, String) Specifies the object key prefix as the match condition for the routing rule.
* `redirect_protocol` - (Optional, String) Specifies the target protocol for the routing rule. Only HTTPS is supported.
* `redirect_replace_key_prefix` - (Optional, String) Specifies the object key prefix to replace the original prefix in the request. You can set this parameter only if the condition is KeyPrefixEquals.
* `redirect_replace_key` - (Optional, String) Specifies the target object key to replace the original object key in the request.

The `source_selection_criteria` object of `replica_rules` supports the following:

* `sse_kms_encrypted_objects` - (Optional, List) Choose whether to copy the KMS-encrypted objects.

The `sse_kms_encrypted_objects` object of `source_selection_criteria` supports the following:

* `status` - (Optional, String) Choose whether to copy KMS encrypted objects; supported values are Enabled and Disabled.

The `tag` object of `and` supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

The `tag` object of `filter` supports the following:

* `key` - (Required, String) Tag key.
* `value` - (Required, String) Tag value.

The `tiering` object of `intelligent_tiering_archiving_rule_list` supports the following:

* `access_tier` - (Required, String) When `rule_id` is not `default`, this parameter is used to specify the archiving or deep archiving tier.  The possible value are: ARCHIVE_ACCESS, DEEP_ARCHIVE_ACCESS.
* `days` - (Required, Int) When the `rule_id` is not set to default, this specifies the number of days after which data is transitioned to the archive or deep archive tier in the intelligent tiering storage configuration. The archive tier (ARCHIVE_ACCESS) supports a range of 91 to 730 days. The deep archive tier (DEEP_ARCHIVE_ACCESS) supports a range of 180 to 730 days. Within the same rule, the number of days for the deep archive tier must be greater than the number of days for the archive tier.

The `transition` object of `lifecycle_rules` supports the following:

* `storage_class` - (Required, String) Specifies the storage class to which you want the object to transition. Available values include `STANDARD_IA`, `MAZ_STANDARD_IA`, `INTELLIGENT_TIERING`, `MAZ_INTELLIGENT_TIERING`, `ARCHIVE`, `DEEP_ARCHIVE`. For more information, please refer to: https://cloud.tencent.com/document/product/436/33417.
* `date` - (Optional, String) Specifies the date after which you want the corresponding action to take effect.
* `days` - (Optional, Int) Specifies the number of days after object creation when the specific rule action takes effect.

The `website` object supports the following:

* `error_document` - (Optional, String) An absolute path to the document to return in case of a 4XX error.
* `index_document` - (Optional, String) COS returns this index document when requests are made to the root domain or any of the subfolders.
* `redirect_all_requests_to` - (Optional, String) Redirects all request configurations. Valid values: http, https. Default is `http`.
* `routing_rules` - (Optional, List) Routing rule configuration. A RoutingRules container can contain up to 100 RoutingRule elements.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cos_bucket_url` - The URL of this cos bucket.


## Import

COS bucket can be imported, e.g.

```
$ terraform import tencentcloud_cos_bucket.bucket bucket-name
```

