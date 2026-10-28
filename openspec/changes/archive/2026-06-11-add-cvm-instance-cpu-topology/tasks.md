## 1. Schema Definition

- [x] 1.1 Add `cpu_topology` block parameter to the `tencentcloud_instance` resource schema in `tencentcloud/services/cvm/resource_tc_instance.go` with type List, MaxItems 1, Optional, ForceNew, containing sub-fields `core_count` (Optional, Int, ForceNew) and `thread_per_core` (Optional, Int, ForceNew)

## 2. Create Logic

- [x] 2.1 In the Create function of `tencentcloud/services/cvm/resource_tc_instance.go`, read the `cpu_topology` block from config and populate `CpuTopology` on the `RunInstancesRequest` with `CoreCount` and `ThreadPerCore` values

## 3. Read Logic

- [x] 3.1 In the Read function of `tencentcloud/services/cvm/resource_tc_instance.go`, read `CpuTopology` from the `Instance` struct returned by `DescribeInstances` and set it into state (with nil check before setting)

## 4. Documentation

- [x] 4.1 Update `tencentcloud/services/cvm/resource_tc_instance.md` to include `cpu_topology` in the example usage

## 5. Testing

- [x] 5.1 Add unit tests in `tencentcloud/services/cvm/resource_tc_instance_cpu_topology_unit_test.go` for the `cpu_topology` parameter using gomonkey mock approach, and verify tests pass with `go test -gcflags=all=-l`
