---
subcategory: "COS Vectors"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_vector_bucket_policy"
sidebar_current: "docs-tencentcloud-resource-cos_vector_bucket_policy"
description: |-
  Provides a resource to manage a COS Vectors bucket policy.
---

# tencentcloud_cos_vector_bucket_policy

Provides a resource to manage a COS Vectors bucket policy.

The resource manages the complete policy, not a single statement. Policy changes update in place; changing the target bucket replaces the policy resource. Removing it deletes the policy without restoring any prior unmanaged policy. Use only one policy resource per bucket.

Use the bucket and index resources from their respective examples. The principal must be a real, authorized business identity. This narrowly scoped query policy does not grant Terraform deployment permissions. Credentials must remain in the execution environment, never in policy JSON or state.

COS policy syntax uses `version = "2.0"` and lowercase `statement/effect/principal/action/resource`. Each statement requires effect `allow` or `deny`, non-empty `principal.qcs`, `action` and `resource` arrays of non-empty strings. CAM evaluates action names, principal identities, resource scope and conditions. AWS IAM JSON, `aws_iam_policy_document` output and AWS ARNs must be converted, not reused verbatim. JSON formatting and statement/action/resource/principal-set ordering do not cause a diff.

For validation and semantic comparison, top-level `version/statement` and statement field names `effect/principal/action/resource/condition/sid` accept case variants; effect values are normalized to lowercase. This does not make the whole policy case-insensitive: `principal.qcs`, identity/action/resource strings, condition contents, Sid values and unknown fields are not case-folded. Conflicting case variants of a known field are rejected. The original policy JSON is sent to the API, so write configuration using the documented COS syntax. Successful policy CRUD or round-trip does not establish CAM authorization; verify the intended access separately with the authorized target identity.

## Example Usage

```hcl
variable "query_principal_qcs" {
  type     = string
  nullable = false

  validation {
    condition     = can(regex("^qcs::cam::uin/[0-9]+:uin/[0-9]+$", var.query_principal_qcs))
    error_message = "Provide an explicitly authorized CAM user QCS; wildcards are not allowed in this allow-policy example."
  }
}

resource "tencentcloud_cos_vector_bucket_policy" "example" {
  vector_bucket_qcs = tencentcloud_cos_vector_bucket.example.vector_bucket_qcs
  policy = jsonencode({
    version = "2.0"
    statement = [{
      effect    = "allow"
      principal = { qcs = [var.query_principal_qcs] }
      action    = ["name/cos:QueryVectors"]
      resource  = [tencentcloud_cos_vector_index.example.index_qcs]
    }]
  })
}
```

## Argument Reference

The following arguments are supported:

* `policy` - (Required, String) COS policy JSON string, normally produced by jsonencode. Uses version 2.0, lowercase statement/effect/principal/action/resource keys and name/cos: actions. JSON formatting and COS policy set order are ignored. Policy content is stored in Terraform state; never include credentials.
* `vector_bucket_qcs` - (Required, String) Target bucket QCS. Reference the vector bucket resource to establish dependencies. Changing the bucket replaces the policy resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Bucket QCS identity.
* `region` - Creation/import region. A different provider region is rejected.
* `vector_bucket_name` - Full bucket name resolved from QCS.

## API Reference

This resource uses `PutVectorBucketPolicy`, `GetVectorBucketPolicy`, `DeleteVectorBucketPolicy` and `GetVectorBucket` to resolve imported bucket names.

## Import

Use the actual bucket QCS:

```hcl
import {
  to = tencentcloud_cos_vector_bucket_policy.example
  id = "qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-vector-example-1250000000"
}
```

The full bucket name is also accepted and resolved through GetVectorBucket:

```bash
terraform import tencentcloud_cos_vector_bucket_policy.example tf-vector-example-1250000000
```

Define configuration preserving the entire existing policy before import. Review the plan before changing it; importing a policy does not merge it with the example. For the optional instance in the complete example, use address `tencentcloud_cos_vector_bucket_policy.example[0]` and configure its principal variable.

