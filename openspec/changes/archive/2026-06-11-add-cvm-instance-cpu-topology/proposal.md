## Why

The `tencentcloud_instance` resource currently does not expose the `CpuTopology` parameter, which allows users to configure CPU physical core count and hyper-threading settings when creating CVM instances. Users need this capability to optimize CPU performance for specific workloads (e.g., disabling hyper-threading for latency-sensitive applications).

## What Changes

- Add a new `cpu_topology` block parameter to the `tencentcloud_instance` resource schema
  - Sub-field `core_count` (Optional, Int): Number of enabled CPU physical cores
  - Sub-field `thread_per_core` (Optional, Int): Threads per core (1 = hyper-threading off, 2 = hyper-threading on)
- Set `cpu_topology` as `ForceNew` since it is only supported in the `RunInstances` (Create) API and not in any Update API
- Read `CpuTopology` from the `DescribeInstances` response (Instance struct) to populate state

## Capabilities

### New Capabilities
- `cpu-topology-param`: Add cpu_topology parameter to tencentcloud_instance resource for configuring CPU topology during instance creation

### Modified Capabilities

## Impact

- `tencentcloud/services/cvm/resource_tc_instance.go`: Add schema definition, create logic, and read logic for `cpu_topology`
- `tencentcloud/services/cvm/resource_tc_instance.md`: Update example usage documentation
- `tencentcloud/services/cvm/resource_tc_instance_test.go`: Add unit tests for the new parameter
- Cloud API dependency: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312` (already vendored)
