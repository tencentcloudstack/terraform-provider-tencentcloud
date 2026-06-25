## 1. Service-layer plumbing

- [x] 1.1 Add `DescribeGa2ForwardingRuleById(ctx, gaId, listenerId, policyId, ruleId string) (*ga2v20250115.ForwardingRuleSet, error)` to `tencentcloud/services/ga2/service_tencentcloud_ga2.go`. Reuse the existing pattern from `DescribeGa2ListenerById`: build the `DescribeForwardingRuleRequest` and set `request.GlobalAcceleratorId`, `request.ListenerId`, `request.ForwardingPolicyId` **outside** the for-loop; iterate with `Limit=100` (literal, no new constant); strict-equal on `*item.ForwardingRuleId == ruleId` (and defensively skip mismatching parent IDs); wrap each SDK page in `resource.Retry(tccommon.ReadRetryTimeout, ...)`; return `(nil, nil)` when not found; emit `[CRITAL]` log on retry failure.
- [x] 1.2 Confirm `connectivity.UseGa2V20250115Client` is already wired (no change required) and that the SDK package import path remains `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115`.

## 2. Resource implementation

- [x] 2.1 Create `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.go` (single file). Top-level layout: package + imports → `ResourceTencentCloudGa2ForwardingRule()` schema → `resourceTencentCloudGa2ForwardingRuleCreate/Read/Update/Delete` → `parseGa2ForwardingRuleId` → `buildRuleConditions` / `buildRuleActions` / `buildOriginHeaders` build helpers → corresponding `flatten*` helpers.
- [x] 2.2 Schema fields per the spec, in this declaration order: `global_accelerator_id` (Required, ForceNew), `listener_id` (Required, ForceNew), `forwarding_policy_id` (Required, ForceNew), `rule_conditions` (Required, TypeSet, nested `rule_condition_type` Required + `rule_condition_value` Required TypeSet of String), `rule_actions` (Required, TypeSet, nested `rule_action_type` Required + `rule_action_value` Required), `origin_headers` (Optional+Computed, TypeSet, nested `key`/`value` Required), `enable_origin_sni` (Optional+Computed Bool), `origin_sni` (Optional+Computed String), `origin_host` (Optional+Computed String), and computed-only: `forwarding_rule_id`. Add a `Timeouts` block with `Create/Update/Delete: 5 * time.Minute`.
- [x] 2.3 Implement `Create`: build `CreateForwardingRuleRequest` from schema (translate sets via `(*schema.Set).List()` then via the build helpers); wrap SDK call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`; defend against nil `Response` / `ForwardingRuleId` / `TaskId`; assert `TaskId != ""` and call `Ga2Service.WaitForGa2TaskFinish(ctx, *resp.Response.TaskId, d.Timeout(schema.TimeoutCreate))`; finally `d.SetId(strings.Join([]string{gaId, listenerId, policyId, ruleId}, tccommon.FILED_SP))` and return `Read`.
- [x] 2.4 Implement `Read`: parse 4-segment composite ID into `(gaId, listenerId, policyId, ruleId)` via `parseGa2ForwardingRuleId`; call `DescribeGa2ForwardingRuleById`; on `(nil, nil)` log `[WARN]` and `d.SetId("")` then return; otherwise `_ = d.Set(...)` for every input field and the computed `forwarding_rule_id`. Use the "set only when non-nil" idiom (no `else` branches). Translate the SDK response field names `RuleCondition` / `RuleAction` (singular) into the schema's `rule_conditions` / `rule_actions` (plural).
- [x] 2.5 Implement `Update`: short-circuit if no body field changed (`rule_conditions`, `rule_actions`, `origin_headers`, `enable_origin_sni`, `origin_sni`, `origin_host`); otherwise build `ModifyForwardingRuleRequest` populated with the 4 mandatory ID fields plus every body field whose schema getter returns a non-zero value. Wrap SDK call in `resource.Retry(tccommon.WriteRetryTimeout, ...)`; defend against nil response/TaskId; call `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate))`.
- [x] 2.6 Implement `Delete`: call `DeleteForwardingRuleWithContext` inside `resource.Retry(tccommon.WriteRetryTimeout, ...)`; populate all 4 identifier fields from the parsed composite ID; defend against nil `Response` / `TaskId`; call `WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete))`.
- [x] 2.7 At the top of every CRUD function, add `defer tccommon.LogElapsed("resource.tencentcloud_ga2_forwarding_rule.<op>")()` and `defer tccommon.InconsistentCheck(d, meta)()`.
- [x] 2.8 Importer: use `schema.ImportStatePassthrough` (matches the other GA2 resources).
- [x] 2.9 ID parser: `parseGa2ForwardingRuleId(id string) (gaId, listenerId, policyId, ruleId string, err error)` — split on `tccommon.FILED_SP`, expect exactly 4 non-empty parts, otherwise return a descriptive error.

## 3. Provider registration

- [x] 3.1 In `tencentcloud/provider.go`, register `"tencentcloud_ga2_forwarding_rule": ga2.ResourceTencentCloudGa2ForwardingRule()` in `ResourcesMap`. Place it adjacent to the other GA2 entries to keep the namespace contiguous.
- [x] 3.2 In `tencentcloud/provider.md`, append `tencentcloud_ga2_forwarding_rule` to the existing `Global Accelerator(GA2)` Resources section so `gendoc` picks it up on the next `make doc`.

## 4. Documentation example

- [x] 4.1 Create `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule.md` (filename pattern matches `resource_tc_config_compliance_pack.md`). Content sections: short description; HCL example showing the 3 required ID fields plus a representative `rule_conditions` block (e.g. host-match), a `rule_actions` block, plus a small `origin_headers` map; an `Import` section showing `terraform import tencentcloud_ga2_forwarding_rule.example ga-xxx#lsr-yyy#fpcy-zzz#frule-www`. **Do not** hand-edit `website/docs/r/ga2_forwarding_rule.html.markdown` — that file is regenerated from the resource Schema/Description by `make doc`.

## 5. Acceptance test scaffolding

- [x] 5.1 Create `tencentcloud/services/ga2/resource_tc_ga2_forwarding_rule_test.go` (filename pattern matches `resource_tc_config_compliance_pack_test.go`). Include `TestAccTencentCloudGa2ForwardingRuleResource_basic` with `resource.Test` containing at minimum: a create step, an update step (e.g. change `origin_host`), and an `ImportState` verification step. Use `tcacctest.AccPreCheck` and the standard provider factories.

## 6. Build & verification

- [x] 6.1 Run `go build ./tencentcloud/...` and confirm clean compilation.
- [x] 6.2 Run `go vet ./tencentcloud/services/ga2/...` and confirm no new findings.
- [x] 6.3 Run `read_lints` on the changed files and confirm no new ERROR/WARN beyond the project-wide deprecated HINTs (`resource.Retry`, `schema.ImportStatePassthrough`, `d.GetOkExists`).
- [x] 6.4 Run `make doc` to regenerate `website/docs/r/ga2_forwarding_rule.html.markdown` from the resource Schema/Description and the `.md` example file. Verify the generated doc renders the example block, all input fields, the nested rule blocks, the computed `forwarding_rule_id`, and the import syntax. **Do not** hand-edit the generated file.
- [ ] 6.5 (Optional, requires `TENCENTCLOUD_SECRET_ID` / `TENCENTCLOUD_SECRET_KEY`) Run the acceptance test: `TF_ACC=1 go test -timeout 30m -run TestAccTencentCloudGa2ForwardingRuleResource_basic ./tencentcloud/services/ga2/...`.
