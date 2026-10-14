---
subcategory: "Secrets Manager(SSM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ssm_secret_versions"
sidebar_current: "docs-tencentcloud-datasource-ssm_secret_versions"
description: |-
  Use this data source to query detailed information of SSM secret version
---

# tencentcloud_ssm_secret_versions

Use this data source to query detailed information of SSM secret version

## Example Usage

```hcl
data "tencentcloud_ssm_secret_versions" "example" {
  secret_name = tencentcloud_ssm_secret_version.v1.secret_name
  version_id  = tencentcloud_ssm_secret_version.v1.version_id
}

resource "tencentcloud_ssm_secret" "example" {
  secret_name = "tf-example"
  description = "desc."

  tags = {
    createdBy = "terraform"
  }
}

resource "tencentcloud_ssm_secret_version" "v1" {
  secret_name   = tencentcloud_ssm_secret.example.secret_name
  version_id    = "v1"
  secret_binary = "MTIzMTIzMTIzMTIzMTIzQQ=="
}
```

## Argument Reference

The following arguments are supported:

* `secret_name` - (Required, String) Secret name used to filter result.
* `result_output_file` - (Optional, String) Used to save results.
* `version_id` - (Optional, String) VersionId used to filter result.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `secret_version_list` - A list of SSM secret versions. When secret status is `Disabled`, this field will not update anymore.
  * `secret_binary` - The base64-encoded binary secret.
  * `secret_string` - The string text of secret.
  * `version_id` - Version of secret.

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeSecret](https://cloud.tencent.com/document/api/1140/40526) - This API is used to obtain the detailed attribute information of a Secret.
* Read: [GetSecretValue](https://cloud.tencent.com/document/api/1140/40522) - For user-defined credentials, this API is used to get the plaintext information of a credential by specifying the credential name and version.
* Read: [ListSecretVersionIds](https://cloud.tencent.com/document/api/1140/40520) - This API is used to obtain list of versions of a Secret.

