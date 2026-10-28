## ADDED Requirements

### Requirement: Resource registration
The provider SHALL expose a resource type named `tencentcloud_ga2_forwarding_rule` that manages a single Tencent Cloud Global Accelerator V2 layer-7 forwarding rule per resource block. The resource MUST be registered in `tencentcloud/provider.go` under the `ga2` namespace, adjacent to the other GA2 entries. The resource MUST also appear under the `Global Accelerator(GA2)` Resources section of `tencentcloud/provider.md` so that `make doc` generates a website page for it.

#### Scenario: Resource type is discoverable
- **WHEN** an operator runs `terraform plan` against a configuration that references `resource "tencentcloud_ga2_forwarding_rule" "<name>"`
- **THEN** Terraform resolves the type without an "unknown resource" error and shows the planned create.

#### Scenario: Provider compiles
- **WHEN** the codebase is built with `go build ./tencentcloud/...`
- **THEN** the build succeeds with no compilation or vet errors related to the new resource.

#### Scenario: Website doc generated
- **WHEN** an operator runs `make doc`
- **THEN** `website/docs/r/ga2_forwarding_rule.html.markdown` is generated from the resource Schema/Description and the example markdown file.

### Requirement: Schema mirrors `CreateForwardingRule`
The resource schema SHALL expose every input parameter accepted by the `CreateForwardingRule` API as a top-level attribute, with no renaming or merging:
- `global_accelerator_id` (string, required, **ForceNew**)
- `listener_id` (string, required, **ForceNew**)
- `forwarding_policy_id` (string, required, **ForceNew**)
- `rule_conditions` (set, required) — nested fields:
  - `rule_condition_type` (string, required)
  - `rule_condition_value` (set of string, required)
- `rule_actions` (set, required) — nested fields:
  - `rule_action_type` (string, required)
  - `rule_action_value` (string, required)
- `origin_headers` (set, optional+computed) — nested fields:
  - `key` (string, required)
  - `value` (string, required)
- `enable_origin_sni` (bool, optional+computed)
- `origin_sni` (string, optional+computed)
- `origin_host` (string, optional+computed)

The resource SHALL additionally expose the following read-only attribute hydrated from `DescribeForwardingRule` response:
- `forwarding_rule_id` (string, computed) — also stored as the 4th segment of `d.Id()`.

#### Scenario: All required SDK input fields are present
- **WHEN** a developer inspects the resource schema
- **THEN** every field declared in `ga2v20250115.CreateForwardingRuleRequestParams` (GlobalAcceleratorId, ListenerId, ForwardingPolicyId, RuleConditions, RuleActions, OriginHeaders, EnableOriginSni, OriginSni, OriginHost) appears in the schema with semantically equivalent typing.

#### Scenario: No undocumented schema fields
- **WHEN** a developer inspects the resource schema
- **THEN** there are no fields beyond those listed above; no derived flags or synthetic toggles are introduced.

### Requirement: Resource ID
The resource ID SHALL be the 4-segment composite `<GlobalAcceleratorId><FILED_SP><ListenerId><FILED_SP><ForwardingPolicyId><FILED_SP><ForwardingRuleId>`, using the project-standard separator `tccommon.FILED_SP`. The resource SHALL support `terraform import` using the composite ID.

#### Scenario: Create sets the composite ID
- **WHEN** `CreateForwardingRule` succeeds and the polled task transitions to `SUCCESS`
- **THEN** the resource calls `d.SetId(strings.Join([]string{gaId, listenerId, policyId, ruleId}, tccommon.FILED_SP))`.

#### Scenario: Import by composite ID
- **WHEN** an operator runs `terraform import tencentcloud_ga2_forwarding_rule.x ga-xxx#lsr-yyy#fpcy-zzz#frule-www`
- **THEN** the resource state is hydrated from `DescribeForwardingRule` using the parsed 4-tuple.

#### Scenario: Malformed import ID
- **WHEN** the import ID does not contain exactly three `tccommon.FILED_SP` separators, or has empty components
- **THEN** the resource returns a descriptive error before any SDK call.

### Requirement: Async create with task polling
On Create, the resource SHALL invoke `CreateForwardingRule`, capture the returned `TaskId` and `ForwardingRuleId`, and poll `DescribeTaskResult` via `Ga2Service.WaitForGa2TaskFinish(ctx, taskId, timeout)` until `Status == "SUCCESS"` or the user-supplied `Timeouts.Create` (default **5 minutes**) elapses.

#### Scenario: Successful async create
- **WHEN** `CreateForwardingRule` succeeds and the polled task transitions to `SUCCESS` within the timeout
- **THEN** the resource sets the composite ID, invokes Read, and returns no error.

#### Scenario: Async create timeout
- **WHEN** the task does not reach `SUCCESS` before the configured `Timeouts.Create`
- **THEN** the resource returns an error containing the task ID and last observed status.

#### Scenario: Empty TaskId or ForwardingRuleId
- **WHEN** `CreateForwardingRule` returns a nil `TaskId` or nil `ForwardingRuleId`
- **THEN** the resource returns an explicit error rather than dereferencing the nil pointer.

### Requirement: Read with retry and pagination
On Read, the resource SHALL call `Ga2Service.DescribeGa2ForwardingRuleById(ctx, gaId, listenerId, policyId, ruleId)`, which:
- Wraps the SDK call `DescribeForwardingRuleWithContext` in `resource.Retry(tccommon.ReadRetryTimeout, ...)`.
- Sets `request.GlobalAcceleratorId`, `request.ListenerId`, `request.ForwardingPolicyId` once outside the loop.
- Iterates pages with `Limit=100` (the documented maximum), constructing the request object once **outside** the loop and only mutating `Offset` / `Limit` per iteration.
- Strict-equals on `*item.ForwardingRuleId == ruleId` before returning. Items with mismatching parent IDs (`GlobalAcceleratorId`, `ListenerId`, `ForwardingPolicyId`) are skipped defensively.
- Returns `(nil, nil)` when no matching rule exists.

When the helper returns `(nil, nil)`, the resource SHALL clear `d.SetId("")` and log a `[WARN]` line indicating the rule may have been deleted out of band.

#### Scenario: Rule present
- **WHEN** the helper finds a matching `ForwardingRuleSet`
- **THEN** the resource populates all schema fields (input + computed) from the response, including the nested collections via the `RuleCondition` / `RuleAction` (singular) describe-side field names.

#### Scenario: Rule removed externally
- **WHEN** the helper returns `(nil, nil)` (no matching rule)
- **THEN** the resource calls `d.SetId("")` and returns no error.

#### Scenario: Pagination request reuse
- **WHEN** the helper paginates through more than one page
- **THEN** a single `DescribeForwardingRuleRequest` instance is reused across pages, with only `Offset` and `Limit` mutated.

### Requirement: Update path
The Update function SHALL:
- Short-circuit (skip the Modify call entirely) when no body field (`rule_conditions`, `rule_actions`, `origin_headers`, `enable_origin_sni`, `origin_sni`, `origin_host`) has changed.
- Always populate the 4 mandatory identifier fields on the `ModifyForwardingRuleRequest`: `GlobalAcceleratorId`, `ListenerId`, `ForwardingPolicyId`, `ForwardingRuleId`.
- Forward every body field whose schema getter returns a non-zero value.
- Wrap the SDK call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`.
- Wait for the returned `TaskId` via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate))`.

#### Scenario: Body field change
- **WHEN** `rule_conditions` or any other body field changes
- **THEN** `ModifyForwardingRule` is called with the updated body, awaited via the task helper.

#### Scenario: No-op update
- **WHEN** `terraform apply` runs but no body field has changed
- **THEN** the resource skips `ModifyForwardingRule` and immediately invokes Read.

### Requirement: Async delete
On Delete, the resource SHALL call `DeleteForwardingRule`, capture the returned `TaskId`, and poll the task to completion via `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete))` (default **5 minutes**). The request body MUST include all four identifier fields parsed from `d.Id()`.

#### Scenario: Successful async delete
- **WHEN** the delete task transitions to `SUCCESS`
- **THEN** the resource returns no error and Terraform marks the resource as destroyed.

### Requirement: Retry coverage
Every SDK call (`CreateForwardingRuleWithContext`, `DescribeForwardingRuleWithContext`, `ModifyForwardingRuleWithContext`, `DeleteForwardingRuleWithContext`) SHALL be invoked from inside a `resource.Retry` block. The retry budget is `tccommon.WriteRetryTimeout` for write operations and `tccommon.ReadRetryTimeout` for read operations.

#### Scenario: Transient SDK error
- **WHEN** any of the four SDK calls returns a transient TencentCloud SDK error
- **THEN** the call is retried via `tccommon.RetryError(e)` until it succeeds or the retry budget is exhausted.

#### Scenario: Nil response defense
- **WHEN** an SDK call returns a nil `Response` (or a nil critical sub-field such as `TaskId` or `ForwardingRuleId`)
- **THEN** the wrapper returns `resource.NonRetryableError` with a descriptive message rather than dereferencing the nil pointer.

### Requirement: Logging conventions
The resource SHALL emit:
- `defer tccommon.LogElapsed("resource.tencentcloud_ga2_forwarding_rule.<op>")()` at the top of every CRUD function.
- `defer tccommon.InconsistentCheck(d, meta)()` at the top of every CRUD function.
- A `[DEBUG]` line per SDK invocation containing the request action, request body, and response body (matching the existing GA2 log format).
- A `[CRITAL]%s ... failed, reason:%+v` line on every retry-block failure.
- A `[WARN]` line when the resource is detected as deleted out of band during Read.

#### Scenario: Standard log lines emitted
- **WHEN** any CRUD operation runs
- **THEN** the operation's elapsed time is logged via `tccommon.LogElapsed` and inconsistency is checked via `tccommon.InconsistentCheck`.

### Requirement: Documentation and tests
The change SHALL include:
- A markdown document at `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.md` containing a self-contained `terraform { ... } resource "tencentcloud_ga2_forwarding_rule" "..." { ... }` example and a `terraform import` example using the 4-segment composite ID. Filename pattern follows `resource_tc_config_compliance_pack.md`.
- An acceptance-test file at `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule_test.go` exercising at minimum: create, basic update (e.g. change `origin_host`), import, and destroy. Filename pattern follows `resource_tc_config_compliance_pack_test.go`.

#### Scenario: Documentation present
- **WHEN** the change is merged
- **THEN** the markdown documentation file exists and contains both an HCL example and an `import` example using the 4-segment composite ID.

#### Scenario: Test file present
- **WHEN** the change is merged
- **THEN** the `_test.go` file declares `TestAccTencentCloudGa2ForwardingRuleResource_basic` (or equivalent) using `resource.Test` with at least two `Steps` (create + update + import).

### Requirement: SDK constraint
The implementation SHALL NOT modify any file under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/`. If a required API is missing from the vendored SDK, the implementer MUST halt and request an SDK upgrade rather than authoring or patching SDK source.

#### Scenario: Vendored SDK is sufficient
- **WHEN** the implementation begins
- **THEN** the four ForwardingRule APIs plus `DescribeTaskResult` and the `RuleCondition` / `RuleAction` / `OriginHeader` / `ForwardingRuleSet` structs are confirmed present under `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115/` before any code is written.
