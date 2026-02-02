## ADDED Requirements

### Requirement: Network Configuration Fields
The data source SHALL expose network-related configuration fields for each instance type returned by the DescribeZoneInstanceConfigInfos API.

#### Scenario: Query instance network card type
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `network_card` field indicating network card type (e.g., 25 for 25G network card)
- **AND** `network_card` field SHALL be optional computed field of type integer
- **AND** `network_card` field SHALL be null if not provided by API

#### Scenario: Query instance bandwidth capacity
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `instance_bandwidth` field indicating internal network bandwidth in Gbps
- **AND** `instance_bandwidth` field SHALL be optional computed field of type float
- **AND** `instance_bandwidth` field SHALL be null if not provided by API

#### Scenario: Query instance packet forwarding capacity
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `instance_pps` field indicating network packet forwarding capacity in 10K PPS units
- **AND** `instance_pps` field SHALL be optional computed field of type integer
- **AND** `instance_pps` field SHALL be null if not provided by API

### Requirement: Instance Type Naming and Identification
The data source SHALL expose human-readable naming and identification fields for instance types.

#### Scenario: Query instance type display name
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `type_name` field with human-readable instance type name
- **AND** `type_name` field SHALL be optional computed field of type string
- **AND** `type_name` field SHALL be null if not provided by API

#### Scenario: Query instance remarks
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `remark` field with instance-specific remarks or notes
- **AND** `remark` field SHALL be optional computed field of type string
- **AND** `remark` field SHALL be null if not provided by API

### Requirement: CPU and Processor Information
The data source SHALL expose detailed CPU and processor specifications for instance types.

#### Scenario: Query CPU processor model
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `cpu_type` field indicating processor model
- **AND** `cpu_type` field SHALL be optional computed field of type string
- **AND** `cpu_type` field SHALL be null if not provided by API

#### Scenario: Query CPU frequency
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `frequency` field indicating CPU frequency information
- **AND** `frequency` field SHALL be optional computed field of type string
- **AND** `frequency` field SHALL be null if not provided by API

### Requirement: GPU and Accelerator Information
The data source SHALL expose GPU and accelerator hardware specifications for instance types.

#### Scenario: Query physical GPU card count
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `gpu_count` field indicating number of physical GPU cards mapped to instance
- **AND** `gpu_count` field SHALL be optional computed field of type float
- **AND** `gpu_count` SHALL distinguish between vGPU (value < 1) and direct-attach GPU (value >= 1)
- **AND** `gpu_count` field SHALL be null if not provided by API

#### Scenario: Query FPGA core count
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `fpga` field indicating number of FPGA cores
- **AND** `fpga` field SHALL be optional computed field of type integer
- **AND** `fpga` field SHALL be null if not provided by API

### Requirement: Storage Configuration
The data source SHALL expose local disk and storage specifications for instance types.

#### Scenario: Query local disk specifications
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `local_disk_type_list` field with list of local disk specifications
- **AND** `local_disk_type_list` SHALL be optional computed field of type list
- **AND** each local disk entry SHALL include fields: `type` (disk type), `partition_type` (partition type), `min_size` (minimum size in GB), `max_size` (maximum size in GB)
- **AND** `local_disk_type_list` SHALL be empty list if instance type does not support local disks

#### Scenario: Query storage block count
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `storage_block_amount` field indicating number of local storage blocks
- **AND** `storage_block_amount` field SHALL be optional computed field of type integer
- **AND** `storage_block_amount` field SHALL be null if not provided by API

### Requirement: Availability and Stock Information
The data source SHALL expose instance type availability and stock status information.

#### Scenario: Query sold out reason
- **WHEN** user queries instance types data source
- **AND** instance type status is "SOLD_OUT"
- **THEN** instance SHALL include `sold_out_reason` field explaining why instance is sold out
- **AND** `sold_out_reason` field SHALL be optional computed field of type string
- **AND** `sold_out_reason` field SHALL be null if instance is not sold out

#### Scenario: Query stock status category
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `status_category` field indicating inventory status
- **AND** `status_category` field SHALL be optional computed field of type string
- **AND** `status_category` SHALL have value of "EnoughStock", "NormalStock", "UnderStock", or "WithoutStock"
- **AND** `status_category` field SHALL be null if not provided by API

### Requirement: Pricing Information
The data source SHALL expose instance pricing information when available from the API.

#### Scenario: Query instance pricing details
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `price` field with nested pricing information structure
- **AND** `price` SHALL be optional computed field of type list with max 1 item
- **AND** `price` structure SHALL include available pricing fields from API response (unit_price, charge_unit, original_price, discount_price based on ItemPrice SDK structure)
- **AND** `price` field SHALL be empty list if pricing information not provided by API

### Requirement: Extended Attributes
The data source SHALL expose extended attributes and additional configuration details for instance types.

#### Scenario: Query instance extended attributes
- **WHEN** user queries instance types data source
- **THEN** each instance type SHALL include `externals` field with nested extended attributes structure
- **AND** `externals` SHALL be optional computed field of type list with max 1 item
- **AND** `externals` structure SHALL include all fields from Externals SDK structure
- **AND** `externals` field SHALL be empty list if extended attributes not provided by API

### Requirement: Backward Compatibility
The data source MUST maintain backward compatibility with existing Terraform configurations.

#### Scenario: Existing configurations remain functional
- **WHEN** user has existing Terraform configuration using instance_types data source
- **AND** configuration only references previously supported fields (availability_zone, instance_type, cpu_core_count, gpu_core_count, memory_size, family, instance_charge_type, status)
- **THEN** configuration SHALL continue to work without modification
- **AND** terraform plan SHALL not show any changes
- **AND** data source output SHALL contain same values as before for existing fields

#### Scenario: New fields do not break existing outputs
- **WHEN** user queries instance_types data source with existing configuration
- **THEN** new computed fields SHALL be populated automatically
- **AND** users SHALL be able to reference new fields in their configurations optionally
- **AND** omission of new fields in output references SHALL not cause errors
