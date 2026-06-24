## Context

The TencentCloud Terraform provider currently has no resource for STS (Security Token Service) AssumeRole operations. Users who need temporary credentials for cross-account access or service role assumption must use external scripts or the CLI. Adding a `tencentcloud_sts_assume_role_operation` resource enables this workflow natively within Terraform.

The STS SDK package (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813`) is already vendored. An existing `tencentcloud/services/sts/` directory may or may not exist; the service client accessor `UseStsClient()` needs to be verified or created.

## Goals / Non-Goals

**Goals:**
- Implement a RESOURCE_KIND_OPERATION resource that calls the STS AssumeRole API
- Expose all AssumeRole input parameters as ForceNew schema fields
- Expose response credentials (token, tmp_secret_id, tmp_secret_key), expired_time, and expiration as Computed fields
- Follow the established operation resource pattern (empty Read/Delete)
- Provide unit tests using gomonkey mocks (no real API calls)
- Register the resource in provider.go and provider.md

**Non-Goals:**
- Implementing AssumeRoleWithSAML or AssumeRoleWithWebIdentity (separate resources if needed)
- Creating a service layer file (the API call is simple enough to inline in the resource Create function)
- Supporting import (operation resources are one-time and have no persistent remote state)
- Storing credentials in a secure backend (standard Terraform state handling applies)

## Decisions

### 1. Resource Pattern: OPERATION type with empty Read/Delete

**Decision**: Follow the same pattern as `resource_tc_config_start_config_rule_evaluation_operation.go`.

**Rationale**: AssumeRole is a one-time API call that returns temporary credentials. There is no remote resource to read or delete. The operation pattern (Create calls API, Read/Delete are no-ops) is the established convention.

**Alternative considered**: Using a data source instead. Rejected because data sources don't support `Sensitive` fields well for credential output, and the operation semantics (side-effect of generating credentials) align better with a resource.

### 2. ID Generation: helper.BuildToken()

**Decision**: Use `helper.BuildToken()` to generate a random ID for the resource.

**Rationale**: There is no meaningful remote ID for a one-time operation. This matches the pattern used by other operation resources.

### 3. Tags Schema: TypeList of objects with key/value

**Decision**: Model `tags` as `TypeList` with `Elem` being a `schema.Resource` containing `key` (string) and `value` (string) fields.

**Rationale**: The STS API `Tags` field is `[]*Tag` where Tag has `Key` and `Value`. This maps naturally to a list of objects in Terraform schema. Note: these are session tags (not resource tags), so the standard `tags` helper pattern does not apply.

### 4. Credentials Output: TypeList with single element containing sensitive fields

**Decision**: Model `credentials` as a `TypeList` (MaxItems: 1) containing `token`, `tmp_secret_id`, and `tmp_secret_key` as Sensitive Computed strings.

**Rationale**: The API returns a single Credentials object with three fields. Using a nested list with MaxItems 1 is the standard Terraform pattern for single complex objects. Marking fields as Sensitive prevents them from appearing in logs/plan output.

### 5. Testing: gomonkey mock-based unit tests

**Decision**: Use gomonkey to mock the STS client's AssumeRoleWithContext method and test the resource Create logic.

**Rationale**: Per project requirements, new terraform resources use gomonkey mocks rather than the terraform acceptance test suite. This avoids real API calls and allows testing business logic in isolation.

## Risks / Trade-offs

- [Credentials in state] Temporary credentials will be stored in Terraform state in plaintext → Users should use encrypted state backends; document this in the resource description.
- [Token expiry] The credentials expire after `duration_seconds` but Terraform state retains them → This is acceptable for operation resources; users understand the one-time nature.
- [No refresh] Since Read is empty, `terraform refresh` won't update credential values → Expected behavior for operation resources; re-apply to get new credentials.
