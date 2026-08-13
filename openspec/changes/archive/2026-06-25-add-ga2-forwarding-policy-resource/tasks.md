## 1. Service Layer

- [x] 1.1 Add `DescribeGa2ForwardingPolicyById` helper to `tencentcloud/services/ga2/service_tencentcloud_ga2.go`: build `DescribeForwardingPolicyRequest` with `GlobalAcceleratorId` and `ListenerId`, paginate with `Limit=100` (documented maximum), strict-equal match on `*item.ForwardingPolicyId == policyId`, return `(nil, nil)` when not found.

## 2. Resource Implementation

- [x] 2.1 Create `tencentcloud/services/ga2/resource_tc_ga2_forwarding_policy.go` with:
  - Resource schema: `global_accelerator_id` (Required, ForceNew, TypeString), `listener_id` (Required, ForceNew, TypeString), `host` (Required, TypeString), `forwarding_policy_id` (Computed, TypeString), `default_host_flag` (Computed, TypeBool).
  - Timeouts block with Create/Update/Delete defaults of 5 minutes.
  - ID parser: parse 3-segment composite ID `<gaId>#<listenerId>#<policyId>` using `tccommon.FILED_SP`.
  - Create function: call `CreateForwardingPolicy`, extract `ForwardingPolicyId`, set composite ID via `d.SetId()`, await async task via `WaitForGa2TaskFinish`, then Read. Include nil-checks on response and `ForwardingPolicyId`.
  - Read function: parse composite ID, call `DescribeGa2ForwardingPolicyById`, set schema fields from `ForwardingPolicySet` response (nil-check all fields before setting), call `d.SetId("")` with log when not found.
  - Update function: parse composite ID, call `ModifyForwardingPolicy` with `Host`, await async task, then Read.
  - Delete function: parse composite ID, call `DeleteForwardingPolicy`, await async task.
  - Import support via `schema.Importer`.

- [x] 2.2 Register `tencentcloud_ga2_forwarding_policy` in `tencentcloud/provider.go` under the `ga2` namespace block, adjacent to existing GA2 resources.

- [x] 2.3 Append `tencentcloud_ga2_forwarding_policy` to the `Global Accelerator(GA2)` Resources section in `tencentcloud/provider.md`.

## 3. Documentation

- [x] 3.1 Create `tencentcloud/services/ga2/resource_tc_ga2_forwarding_policy.md`
  - One-line description: "Provides a resource to create a GA2 forwarding policy"
  - Example Usage HCL snippet showing resource creation with required fields.
  - Import section showing `terraform import tencentcloud_ga2_forwarding_policy.example <gaId>#<listenerId>#<policyId>`.

## 4. Testing

- [x] 4.1 Create `tencentcloud/services/ga2/resource_tc_ga2_forwarding_policy_test.go` with unit tests using gomonkey mocks.

## 5. Finalization

- [ ] 5.1 Run `gofmt` on all changed Go files.
- [ ] 5.2 Run `make doc` to generate website documentation.
- [ ] 5.3 Create `.changelog/<PR_NUMBER>.txt` with changelog entry.