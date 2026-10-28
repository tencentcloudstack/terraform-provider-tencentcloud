## Context

The Terraform Provider for TencentCloud currently manages CDB (MySQL) instances through various resources, but there is no resource to manage CPU elastic expansion (弹性扩容) for CDB instances. The cloud API provides three interfaces for this capability: `StartCpuExpand` (enable expansion), `DescribeCPUExpandStrategyInfo` (query expansion strategy), and `StopCpuExpand` (disable expansion). Both `StartCpuExpand` and `StopCpuExpand` are asynchronous interfaces that return `AsyncRequestId`, requiring polling after invocation.

This is a RESOURCE_KIND_ATTACHMENT resource — it manages the binding/unbinding relationship between a CDB instance and its CPU elastic expansion strategy. The resource uses only `instance_id` as the composite ID component since the expansion strategy is uniquely tied to a single instance.

## Goals / Non-Goals

**Goals:**
- Add a new Terraform resource `tencentcloud_cdb_start_cpu_expand` to manage CPU elastic expansion binding for CDB instances
- Support all four expansion types: auto, manual, timeInterval, and period
- Support async operation polling after Create and Delete
- Support Import for existing expansion configurations
- Follow the attachment resource pattern (CRD only, no Update operation)
- Implement unit tests using gomonkey mock approach

**Non-Goals:**
- Supporting Update operation for the resource (attachment resources are immutable — changes require re-creation)
- Adding datasource for querying CPU expansion strategies (this is handled by the Read operation within the resource)
- Modifying existing CDB resources or their schemas

## Decisions

1. **Resource ID Strategy**: Use `instance_id` as the sole resource ID. Since CPU elastic expansion is uniquely bound to a single CDB instance, no composite ID is needed. The resource ID is simply the `instance_id` value.

2. **CRD-only Pattern (No Update)**: Since this is an RESOURCE_KIND_ATTACHMENT, the resource only supports Create, Read, and Delete. All top-level fields besides `instance_id` are immutable. The Update method will check `immutableArgs` and return an error if any immutable field has changed. This follows the standard pattern for attachment resources that only have CRD interfaces.

3. **Async Operation Polling**: Both `StartCpuExpand` and `StopCpuExpand` return `AsyncRequestId`. After calling these APIs, the resource MUST poll the `DescribeCPUExpandStrategyInfo` API until the operation completes:
   - After Create: Poll until the expansion strategy appears in the Read response
   - After Delete: Poll until the expansion strategy disappears from the Read response

4. **Schema Design for Strategy Parameters**: The three strategy types (auto_strategy, time_interval_strategy, period_strategy) are modeled as `TypeList` with `MaxItems: 1` nested blocks, following the standard Terraform pattern for single-instance complex objects. This allows users to configure strategy parameters based on the `type` field.

5. **Nested Structure Flattening**: The `PeriodStrategy` contains `TimeCycle` and `TimeInterval` sub-structures. These will be flattened into the `period_strategy` block as nested `time_cycle` and `time_interval` sub-blocks rather than creating separate top-level schema entries.

6. **Deprecated Fields Handling**: The `AutoStrategy` struct contains deprecated fields `ExpandPeriod` and `ShrinkPeriod` (replaced by `ExpandSecondPeriod` and `ShrinkSecondPeriod`). Only the non-deprecated fields will be included in the Terraform schema.

7. **Resource Naming**: The resource will be named `tencentcloud_cdb_start_cpu_expand` following the convention. The Go file will be `resource_tc_cdb_start_cpu_expand_attachment.go` following the attachment naming pattern. The provider registration name will use the `tencentcloud_cdb_start_cpu_expand` key.

8. **Error Handling in Read**: If `DescribeCPUExpandStrategyInfo` returns NULL for the `Type` field, it means the instance has no elastic expansion enabled. This indicates the resource has been deleted outside of Terraform, and `d.SetId("")` should be called. Before clearing the ID, a log message must be printed to preserve diagnostic information.

## Risks / Trade-offs

- [Risk] The `DescribeCPUExpandStrategyInfo` API returns NULL `Type` when no expansion is enabled, which could be confused with a transient API error → Mitigation: In the Read method, after retry exhaustion, log the empty response before clearing the ID, and treat NULL `Type` as a definitive signal that the resource no longer exists
- [Risk] Async operations may take significant time to complete → Mitigation: Implement proper polling with configurable timeouts using the Terraform SDK Timeouts feature
- [Risk] The `StopCpuExpand` async delete might fail, leaving the expansion still enabled → Mitigation: Use retry logic for the delete operation and poll until confirmed deletion
- [Trade-off] No Update operation means users must destroy and recreate the resource to change expansion strategy parameters → This is acceptable since it matches the cloud API behavior (changing strategy type requires stopping the current expansion first)
