## ADDED Requirements

### Requirement: Resource tencentcloud_cfw_cluster_nat_fw_switch can be created
The system SHALL allow users to create a NAT CCN cluster mode firewall switch by calling `OpenClusterNatFwSwitch` with a `NatCcnSwitchConfig` object containing `nat_ins_id`, `ccn_id`, `switch_mode`, `routing_mode`, `access_instance_list`, and optionally `lead_vpc_cidr`.

#### Scenario: Successful creation of cluster nat fw switch
- **WHEN** user provides valid `nat_ccn_switch` block with `nat_ins_id`, `ccn_id`, `switch_mode`, and `routing_mode`
- **THEN** the resource is created by calling `OpenClusterNatFwSwitch` and the resource ID is set to `nat_ins_id#ccn_id`

#### Scenario: Creation fails when nat_ins_id is missing
- **WHEN** user omits `nat_ins_id` from the `nat_ccn_switch` block
- **THEN** Terraform returns a validation error before calling the API

### Requirement: Resource tencentcloud_cfw_cluster_nat_fw_switch can be read
The system SHALL allow users to read the current state of a NAT CCN cluster mode firewall switch by calling `DescribeNatCcnFwSwitch` with `nat_ins_id` and `ccn_id`.

#### Scenario: Successful read of existing switch
- **WHEN** the resource exists with a valid `nat_ins_id#ccn_id` composite ID
- **THEN** the system calls `DescribeNatCcnFwSwitch` and populates `switch_mode`, `routing_mode`, `bypass`, `ccn_id`, and `access_instance_list` from the response

#### Scenario: Read returns empty response
- **WHEN** `DescribeNatCcnFwSwitch` returns nil or empty response
- **THEN** the system logs the resource ID and sets the resource ID to empty string to indicate the resource no longer exists

### Requirement: Resource tencentcloud_cfw_cluster_nat_fw_switch can be updated
The system SHALL allow users to update the NAT CCN cluster mode firewall switch configuration by calling `ModifyClusterNatFwSwitch` with an updated `NatCcnSwitchConfig`.

#### Scenario: Successful update of switch configuration
- **WHEN** user changes `nat_ccn_switch` fields such as `switch_mode`, `routing_mode`, or `access_instance_list`
- **THEN** the system calls `ModifyClusterNatFwSwitch` with the updated configuration

### Requirement: Resource tencentcloud_cfw_cluster_nat_fw_switch can be deleted
The system SHALL allow users to delete (close) a NAT CCN cluster mode firewall switch by calling `CloseClusterNatFwSwitch` with `nat_ins_id` and `ccn_id`.

#### Scenario: Successful deletion of cluster nat fw switch
- **WHEN** user runs `terraform destroy` on the resource
- **THEN** the system calls `CloseClusterNatFwSwitch` with the `nat_ins_id` and `ccn_id` from the resource ID

### Requirement: Resource supports import
The system SHALL allow users to import an existing NAT CCN cluster mode firewall switch using the composite ID `nat_ins_id#ccn_id`.

#### Scenario: Successful import using composite ID
- **WHEN** user runs `terraform import tencentcloud_cfw_cluster_nat_fw_switch.example nat_ins_id#ccn_id`
- **THEN** the resource state is populated from `DescribeNatCcnFwSwitch`

### Requirement: Resource is registered in provider
The system SHALL register `tencentcloud_cfw_cluster_nat_fw_switch` in `tencentcloud/provider.go` and document it in `tencentcloud/provider.md`.

#### Scenario: Resource is available in provider
- **WHEN** user references `tencentcloud_cfw_cluster_nat_fw_switch` in a Terraform configuration
- **THEN** Terraform recognizes the resource type and can plan/apply changes
