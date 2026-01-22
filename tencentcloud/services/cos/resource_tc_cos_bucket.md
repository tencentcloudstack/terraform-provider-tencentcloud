Provides a COS resource to create a COS bucket and set its attributes.

~> **NOTE:** The following capabilities do not support cdc scenarios: `multi_az`, `website`, and bucket replication `replica_role`.

~> **NOTE:** If `chdfs_ofs` is `true`, cannot set `acl_body`, `acl`, `origin_pull_rules`, `origin_domain_rules`, `website`, `encryption_algorithm`, `kms_id`, `versioning_enable`, `acceleration_enable` at the same time. For more information, please refer to `https://www.tencentcloud.com/document/product/436/43305`.

Example Usage

Private Bucket

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

Private Bucket with CDC cluster

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

Enable SSE-KMS encryption 

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

Creation of multiple available zone bucket

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

Using verbose acl

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

Using verbose acl with CDC cluster

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

Static Website

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
        condition_prefix            = "/test"
        redirect_protocol           = "https"
        redirect_replace_key        = "key"
      }
    }
  }
}

output "endpoint_test" {
    value = tencentcloud_cos_bucket.bucket_with_static_website.website.0.endpoint
}
```

Using CORS

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

Using Origin pull

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

Using CORS with CDC

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
    expose_headers = ["Etag"]
  }
}
```

Using object lifecycle

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

Using object lifecycle with CDC

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

Using replication

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

Using intelligent tiering, Only enable intelligent tiering

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

Using intelligent tiering and configure the intelligent tiered storage archiving and deep archiving rules list.

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

Using object lock config

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

Using OFS

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
Import

COS bucket can be imported, e.g.

```
$ terraform import tencentcloud_cos_bucket.bucket bucket-name
```