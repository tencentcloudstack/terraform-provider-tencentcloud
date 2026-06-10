## Context

The TencentCloud CFW (Cloud Firewall) product provides NAT CCN cluster mode firewall switches. Currently, the Terraform provider has no resource to manage these switches. The cloud APIs available are:
- `OpenClusterNatFwSwitch`: Opens (creates) a NAT CCN cluster mode firewall switch
- `DescribeNatCcnFwSwitch`: Queries the NAT CCN firewall switch configuration
- `ModifyClusterNatFwSwitch`: Modifies the NAT CCN cluster mode firewall switch configuration
- `CloseClusterNatFwSwitch`: Closes (deletes) the NAT CCN cluster mode firewall switch

The resource is identified by a composite ID of `nat_ins_id` and `ccn_id` (the NAT firewall instance ID and the CCN instance ID), since both are required to uniquely identify a switch.

The `NatCcnSwitchConfig` struct (used in Open and Modify) contains: `NatInsId`, `CcnId`, `SwitchMode`, `RoutingMode`, `AccessInstanceList`, and `LeadVpcCidr`.

## Goals / Non-Goals

**Goals:**
- Implement `tencentcloud_cfw_cluster_nat_fw_switch` as a RESOURCE_KIND_GENERAL resource
- Support full CRUD lifecycle: Create (OpenClusterNatFwSwitch), Read (DescribeNatCcnFwSwitch), Update (ModifyClusterNatFwSwitch), Delete (CloseClusterNatFwSwitch)
- Use composite ID `nat_ins_id#ccn_id` as the resource ID
- Register the resource in provider.go and provider.md
- Add unit tests using gomonkey mocks

**Non-Goals:**
- Managing NAT firewall instances themselves (handled by `tencentcloud_cfw_nat_instance`)
- Managing VPC firewall switches (separate resource)

## Decisions

### Decision 1: Composite ID using nat_ins_id and ccn_id
The `DescribeNatCcnFwSwitch` and `CloseClusterNatFwSwitch` APIs both require `NatInsId` and `CcnId` as input. Using `nat_ins_id#ccn_id` as the composite ID (with `tccommon.FILED_SP` separator) allows the resource to be uniquely identified and supports import.

### Decision 2: nat_ccn_switch as a nested object for Create/Update
The `OpenClusterNatFwSwitch` and `ModifyClusterNatFwSwitch` APIs accept a `NatCcnSwitch` parameter of type `NatCcnSwitchConfig`. This will be mapped to a `nat_ccn_switch` block in the Terraform schema containing `nat_ins_id`, `ccn_id`, `switch_mode`, `routing_mode`, `access_instance_list`, and `lead_vpc_cidr`.

### Decision 3: Read response fields as computed attributes
The `DescribeNatCcnFwSwitch` response returns `SwitchMode`, `RoutingMode`, `Bypass`, `CcnId`, and `AccessInstanceList`. These are mapped as computed attributes at the top level for easy access.

### Decision 4: ForceNew on nat_ins_id and ccn_id
Since the resource is identified by `nat_ins_id` and `ccn_id`, changing these requires destroying and recreating the resource.

## Risks / Trade-offs

- [Risk] The `DescribeNatCcnFwSwitch` API may return empty if the switch has not been fully created yet → Mitigation: Use retry logic in Read to handle eventual consistency.
- [Risk] The Open/Close/Modify APIs may be asynchronous in practice → Mitigation: After calling Open/Modify/Close, poll `DescribeNatCcnFwSwitch` to confirm the state change.

## Migration Plan

No migration needed. This is a new resource with no existing state.

## Open Questions

None.
