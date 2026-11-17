Provides a resource to manage a COS Vectors bucket.

This is a vector bucket, not a normal COS object bucket or a Terraform state backend. The resource uses the existing TencentCloud provider credentials and the dedicated regional HTTPS endpoint `vectors.<region>.coslake.com`. Ordinary `cos_domain`, `domain` and `protocol` overrides do not apply.

New buckets default to explicit AES256 (SSE-COS) encryption, visible in the creation plan. This is a provider default, not a statement about the API default. Import preserves actual encryption, including an unencrypted bucket. Switching between explicitly configured AES256 and an omitted configuration is a no-op for an existing AES256 bucket; explicitly adding AES256 to an imported unencrypted bucket replaces it. Indexes and bucket policies are managed by `tencentcloud_cos_vector_index` and `tencentcloud_cos_vector_bucket_policy`. Tags and KMS are not supported.

-> **Note:** Changing the bucket name or explicitly configured encryption replaces the bucket. `force_destroy` defaults to false. If explicitly applied as true, destroy deletes all indexes, including unmanaged indexes and all their vectors, before deleting the bucket. This is irreversible. Index deletion is destructive even without `force_destroy`. Use separate provider aliases for different regions rather than changing the region of a managed bucket.

Replace the sample APPID with your actual account APPID. Authentication uses the full provider's shared client: `profile`/`shared_credentials_dir`, `assume_role`, `cam_role_name`, and `assume_role_with_web_identity` are available, as are environment `TENCENTCLOUD_SECRET_ID`, `TENCENTCLOUD_SECRET_KEY` and temporary `TENCENTCLOUD_SECURITY_TOKEN`. Keep secrets outside HCL. TencentCloud profiles use tccli JSON files, not AWS INI; the current web identity token file must be JSON containing `web_identity_token`, not a raw JWT. CVM role metadata access must be explicitly authorized. Temporary credentials are initialization snapshots, not a verified automatically refreshing AWS credential chain.

API

This resource uses `CreateVectorBucket`, `GetVectorBucket` and `DeleteVectorBucket`. Deletion may use `ListIndexes`; explicitly applied `force_destroy` also uses `DeleteIndex` and `GetIndex`.

Example Usage

```hcl
provider "tencentcloud" {
  region = "ap-guangzhou"
}

resource "tencentcloud_cos_vector_bucket" "example" {
  vector_bucket_name = "tf-vector-example-1250000000"
  encryption_configuration {
    sse_type = "AES256"
  }
}
```

Import

A vector bucket can be imported using the full bucket name including APPID:

```
terraform import tencentcloud_cos_vector_bucket.example tf-vector-example-1250000000
```

The bucket QCS is also accepted. Its region must match the provider region, and its account must match the bucket APPID. AWS ARNs are not accepted.

```
terraform import tencentcloud_cos_vector_bucket.example 'qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-vector-example-1250000000'
```

When importing an unencrypted bucket, omit `encryption_configuration` to preserve it. Adding AES256 afterwards plans a replacement.

## Index and Policy Resources

The example in `examples/tencentcloud-cos-vector-bucket/` contains a vector bucket and index, with an optional narrowly scoped bucket policy. Use all `.tf` files in that directory: `providers.tf`, `variables.tf`, `main.tf` and `outputs.tf`. Resource references establish dependency ordering: bucket first, index next, policy last; destroy reverses dependencies. No vector data is stored in Terraform state.

`tencentcloud_cos_vector_index` requires `vector_bucket_name`, `index_name`, `data_type` (`float32`), `dimension` (1-4096 per the current SDK contract), and `distance_metric` (`cosine` or `euclidean`). All are immutable and changing them replaces the index. It returns `id` (`bucket-name/index-name`), `index_qcs`, `region`, and raw `creation_time`. Optional `metadata_configuration.non_filterable_metadata_keys` is a set of 1-10 keys matching `^[a-zA-Z_][a-zA-Z0-9_]{0,62}$`; configured changes replace the index, while omitting it preserves imported metadata. Deletion removes all vectors in that index irreversibly.

```
terraform import tencentcloud_cos_vector_index.example tf-vector-example-1250000000/example-index
terraform import tencentcloud_cos_vector_index.example 'qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-vector-example-1250000000/index/example-index'
```

`tencentcloud_cos_vector_bucket_policy` requires `vector_bucket_qcs` and a `policy` JSON string. Target changes replace the resource, while policy changes update in place. The full policy is overwritten; removing the resource deletes the policy, not the bucket, and does not restore a previously unmanaged policy. Use only one policy resource per bucket. COS uses `version = "2.0"`, lowercase `statement/effect/principal/action/resource`, QCS principals/resources and `name/cos:` actions; do not copy AWS IAM policy JSON. Formatting and statement/action/resource/principal-set order are ignored, while real policy changes remain visible. Policy content is persisted in state, so never include credentials.

```
terraform import tencentcloud_cos_vector_bucket_policy.example 'qcs::cosvector:ap-guangzhou:uid/1250000000:bucket/tf-vector-example-1250000000'
```

Bucket-name import is also supported for policies; it resolves the QCS through GetVectorBucket. If using the optional example resource, its import address is `tencentcloud_cos_vector_bucket_policy.example[0]` and `policy_principal_qcs` must be configured. The principal must be an authorized test identity, never anonymous. Recorded authorized cloud verification covered bucket/index/policy lifecycle, name/QCS imports, no-op plans, metadata preservation, refusal to delete a bucket containing an unmanaged index with `force_destroy = false`, and cleanup after applying `force_destroy = true`. Policy CRUD and round-trip verification do not prove CAM enforcement; the observed Deny authorization effect remains unconfirmed.

## Local Verification

Run the new unit tests without cloud access:

```bash
go test ./tencentcloud/connectivity ./tencentcloud/ratelimit ./tencentcloud/services/cosvectors -run 'TestCosVector|TestProCheckContext|TestCheckContext' -count=1
```

Run Terraform lifecycle tests against an automatically started loopback mock, using an installed Terraform CLI:

```bash
CHECKPOINT_DISABLE=1 TENCENTCLOUD_COS_VECTOR_LOCAL_TEST=1 \
TF_ACC_TERRAFORM_PATH=/absolute/path/to/terraform \
go test ./tencentcloud/services/cosvectors -run 'TestCosVector.*Local' -count=1 -v -timeout=5m
```

These tests use fake credentials and do not require `TF_ACC=1`. They exercise the production resource and COS SDK through a test-only provider; full provider mux/schema compatibility is covered separately by `make check-mux`. They do not prove real CAM authorization, service availability or cloud-side lifecycle behavior.

For authorized cloud verification, first confirm that Terraform Registry has published a TencentCloud provider version containing these resources. PR merge alone does not make a Registry version available. Use all `.tf` files in `examples/tencentcloud-cos-vector-bucket/`; the example declares the official provider source without a version constraint. For production, constrain the provider to a verified release and retain `.terraform.lock.hcl`. Run normal `terraform init`, `validate`, `plan` and `apply`, supplying the region and a dedicated disposable full bucket name. For an existing configuration pinned to an older version, update its constraints as needed, run `terraform init -upgrade`, verify the selected version and review the plan. Local development packages still require a matching filesystem mirror and an explicit local version constraint.

Review plan, then apply, require a subsequent plan to have no differences, verify import, and finally destroy only the disposable bucket. A bucket creation whitelist may be required. Do not use benchmark or production buckets. Only the official public HTTPS endpoint is enabled; private endpoints require explicit authorization and a separate restricted configuration design.

Create is sent once. Delete normally runs once, but a valid DeleteVectorBucket response exactly matching HTTP 409, `ConflictException`, and `bucket not empty` triggers a ListIndexes check. If no visible indexes remain, deletion is retried with bounded backoff and jitter while asynchronous GC completes, within the five-minute delete context. Visible indexes or any list error stop deletion; ListIndexes permission is required for this path. Other conflicts and unknown network outcomes are not retried by this waiter. Delete success or a recognized DeleteBucket NotFound is independently confirmed through Get polling, up to 30 seconds within the remaining context; list emptiness or list NotFound alone never removes state. Timeout or confirmation failure keeps the bucket tracked. Cancellation interrupts rate-limit waits. The vector SDK client and connection pool are reused within the same provider client; changing credentials or region rebuilds that client without mutating old credentials. If a creation request fails before a confirmed success response, query the bucket and import it if present instead of blindly retrying apply. After a confirmed Create 200, identity is retained even if response parsing or refresh fails; Terraform may mark the resource tainted, so inspect state and the replacement plan before recovery. Malformed/unknown 404, 403, 409 and 500 responses do not remove state. `creation_time` preserves the SDK integer without inferring units; an unreported/zero value remains null.
