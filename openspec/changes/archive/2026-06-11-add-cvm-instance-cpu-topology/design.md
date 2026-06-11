## Context

The `tencentcloud_instance` resource manages CVM instances. The TencentCloud CVM API supports a `CpuTopology` parameter in the `RunInstances` request that allows specifying CPU physical core count and hyper-threading configuration. The `DescribeInstances` API returns `CpuTopology` in the `Instance` struct. However, no Update API (e.g., `ResetInstancesType`, `ModifyInstancesAttribute`) supports modifying `CpuTopology` after creation.

The `CpuTopology` struct in the SDK (`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312/models.go`) contains:
- `CoreCount *int64` - Number of enabled CPU physical cores
- `ThreadPerCore *int64` - Threads per core (1 = hyper-threading off, 2 = hyper-threading on)

## Goals / Non-Goals

**Goals:**
- Add `cpu_topology` as an Optional, ForceNew block parameter to `tencentcloud_instance`
- Pass `CpuTopology` to `RunInstances` during instance creation
- Read `CpuTopology` from `DescribeInstances` response and populate Terraform state
- Maintain backward compatibility (existing configurations without `cpu_topology` continue to work)

**Non-Goals:**
- Supporting in-place update of CPU topology (not supported by any CVM Update API)
- Adding CPU topology to data sources or other CVM resources
- Validating core_count/thread_per_core values against instance type constraints (left to API-side validation)

## Decisions

### 1. Schema type: List with MaxItems=1

**Decision**: Use `schema.TypeList` with `MaxItems: 1` containing an `Elem: &schema.Resource` with two Optional Int fields.

**Rationale**: This is the standard pattern in this provider for nested object parameters (e.g., `system_disk`, `internet_max_bandwidth_out`). It provides clear structure and allows Terraform to diff individual sub-fields.

**Alternative considered**: Using two top-level fields `cpu_topology_core_count` and `cpu_topology_thread_per_core`. Rejected because the API models this as a nested struct, and grouping them in a block is more idiomatic Terraform.

### 2. ForceNew behavior

**Decision**: Mark the entire `cpu_topology` block as `ForceNew`.

**Rationale**: `CpuTopology` is only available in `RunInstances` (Create). No Update API supports modifying it. Changing CPU topology requires destroying and recreating the instance.

### 3. Read from DescribeInstances Instance struct

**Decision**: Read `CpuTopology` from the `Instance` struct returned by `DescribeInstances`.

**Rationale**: The `Instance` struct (line 6833 in models.go) includes `CpuTopology *CpuTopology`. This is the standard read path for instance attributes.

## Risks / Trade-offs

- [ForceNew causes instance recreation] → Users are warned via Terraform plan output. This is unavoidable given API limitations.
- [API may return nil CpuTopology for instances created without it] → Read logic must check for nil before setting state, following the existing pattern in the resource.
