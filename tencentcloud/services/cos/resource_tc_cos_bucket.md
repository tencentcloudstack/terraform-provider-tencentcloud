Provides a COS resource to create a COS bucket and set its attributes.

~> **NOTE:** The following capabilities do not support cdc scenarios: `multi_az`, `website`, and bucket replication `replica_role`.

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
  encryption_algorithm = "KMS" #cos/kms for cdc cos
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
    index_document = "index.html"
    error_document = "error.html"
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
    id                 = "test-rep1"
    status             = "Enabled"
    prefix             = "dist"
    destination_bucket = "qcs::cos:${local.region}::${tencentcloud_cos_bucket.bucket_replicate.bucket}"
  }
}
```

Import

COS bucket can be imported, e.g.

```
$ terraform import tencentcloud_cos_bucket.bucket bucket-name
```