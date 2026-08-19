## ADDED Requirements

### Requirement: cpu_topology schema definition
The `tencentcloud_instance` resource SHALL expose an Optional `cpu_topology` block parameter with `ForceNew: true` and `MaxItems: 1`. The block SHALL contain two Optional Int sub-fields: `core_count` (number of enabled CPU physical cores) and `thread_per_core` (threads per core, 1 for hyper-threading off, 2 for hyper-threading on).

#### Scenario: Schema includes cpu_topology block
- **WHEN** the Terraform schema for `tencentcloud_instance` is inspected
- **THEN** it SHALL contain a `cpu_topology` field of type List with MaxItems 1, ForceNew true, and Elem containing `core_count` (Optional, Int) and `thread_per_core` (Optional, Int)

### Requirement: cpu_topology passed to RunInstances on create
The resource Create function SHALL read the `cpu_topology` block from the Terraform configuration and, if specified, populate the `CpuTopology` field on the `RunInstancesRequest` with the corresponding `CoreCount` and `ThreadPerCore` values before calling the API.

#### Scenario: Instance created with cpu_topology specified
- **WHEN** a user specifies `cpu_topology { core_count = 4, thread_per_core = 1 }` in their Terraform config
- **THEN** the `RunInstances` API call SHALL include `CpuTopology` with `CoreCount = 4` and `ThreadPerCore = 1`

#### Scenario: Instance created without cpu_topology
- **WHEN** a user does not specify `cpu_topology` in their Terraform config
- **THEN** the `RunInstances` API call SHALL NOT include the `CpuTopology` field

### Requirement: cpu_topology read from DescribeInstances
The resource Read function SHALL read the `CpuTopology` field from the `Instance` struct returned by `DescribeInstances` and set it into the Terraform state. If the API returns nil for `CpuTopology`, the field SHALL NOT be set.

#### Scenario: Read instance with CpuTopology set
- **WHEN** the `DescribeInstances` API returns an Instance with `CpuTopology` containing `CoreCount = 4` and `ThreadPerCore = 2`
- **THEN** the Terraform state SHALL contain `cpu_topology.0.core_count = 4` and `cpu_topology.0.thread_per_core = 2`

#### Scenario: Read instance with CpuTopology nil
- **WHEN** the `DescribeInstances` API returns an Instance with `CpuTopology` as nil
- **THEN** the Terraform state SHALL NOT set the `cpu_topology` field

### Requirement: cpu_topology triggers ForceNew on change
The `cpu_topology` parameter SHALL be marked as ForceNew. Any change to `cpu_topology` (including adding, removing, or modifying sub-fields) SHALL trigger instance recreation.

#### Scenario: Changing cpu_topology triggers replacement
- **WHEN** a user changes `cpu_topology` from `{ core_count = 4, thread_per_core = 1 }` to `{ core_count = 8, thread_per_core = 2 }`
- **THEN** Terraform plan SHALL show the instance will be destroyed and recreated
