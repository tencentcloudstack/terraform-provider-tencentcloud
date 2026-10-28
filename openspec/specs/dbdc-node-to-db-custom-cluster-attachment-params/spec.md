# dbdc-node-to-db-custom-cluster-attachment-params Specification

## Purpose
Defines the schema and CRUD behavior of the `tencentcloud_dbdc_node_to_db_custom_cluster_attachment` resource, including the node-add parameters (labels, taints, host_name, host_name_type), computed node fields (network_mode, eni_ip), and the passing of login_settings on delete.

## Requirements
### Requirement: labels argument
The system SHALL define a `labels` argument on the `tencentcloud_dbdc_node_to_db_custom_cluster_attachment` resource schema as a TypeList (MaxItems 20) of schema.Resource, Optional and ForceNew. Each element SHALL contain:
- `key`: Required, TypeString — the label key
- `value`: Optional, TypeString — the label value

The system SHALL map each `labels` element to a `dbdcv20201029.Label{Key, Value}` and set `AddNodesToDBCustomClusterRequest.Labels` in the Create function.

#### Scenario: Create with labels provided
- **WHEN** the user provides `labels` with elements `[{key="env", value="prod"}]`
- **THEN** the system SHALL build a `[]*Label` containing `{Key:"env", Value:"prod"}` and set it on the `AddNodesToDBCustomClusterRequest` before calling the API

#### Scenario: Create without labels
- **WHEN** the user does not provide `labels`
- **THEN** the system SHALL NOT set `Labels` on the `AddNodesToDBCustomClusterRequest`

### Requirement: taints argument
The system SHALL define a `taints` argument on the `tencentcloud_dbdc_node_to_db_custom_cluster_attachment` resource schema as a TypeList (MaxItems 5) of schema.Resource, Optional and ForceNew. Each element SHALL contain:
- `key`: Required, TypeString — the taint key
- `effect`: Required, TypeString — the taint effect (NoSchedule, PreferNoSchedule, NoExecute)
- `value`: Optional, TypeString — the taint value

The system SHALL map each `taints` element to a `dbdcv20201029.Taint{Key, Effect, Value}` and set `AddNodesToDBCustomClusterRequest.Taints` in the Create function.

#### Scenario: Create with taints provided
- **WHEN** the user provides `taints` with elements `[{key="dedicated", effect="NoSchedule", value="true"}]`
- **THEN** the system SHALL build a `[]*Taint` containing `{Key:"dedicated", Effect:"NoSchedule", Value:"true"}` and set it on the `AddNodesToDBCustomClusterRequest`

#### Scenario: Create without taints
- **WHEN** the user does not provide `taints`
- **THEN** the system SHALL NOT set `Taints` on the `AddNodesToDBCustomClusterRequest`

### Requirement: host_name and host_name_type arguments
The system SHALL define `host_name` (Optional, TypeString, ForceNew) and `host_name_type` (Optional, TypeInt, ForceNew) arguments on the resource schema. In the Create function, the system SHALL set `AddNodesToDBCustomClusterRequest.HostName` when `host_name` is provided, and SHALL set `AddNodesToDBCustomClusterRequest.HostNameType` when `host_name_type` is provided.

#### Scenario: Create with host_name and host_name_type
- **WHEN** the user provides `host_name = "node-{R:1}"` and `host_name_type = 1`
- **THEN** the system SHALL set `request.HostName = helper.String("node-{R:1}")` and `request.HostNameType = helper.Int64(1)` on the `AddNodesToDBCustomClusterRequest`

#### Scenario: Create without host_name fields
- **WHEN** the user does not provide `host_name` or `host_name_type`
- **THEN** the system SHALL NOT set `HostName` or `HostNameType` on the request

### Requirement: network_mode computed field
The system SHALL define a `network_mode` Computed TypeString field on the resource schema. In the Read function, the system SHALL read `NetworkMode` from the node data returned by the read API and set it on the resource data when the field is not nil.

#### Scenario: Read with network_mode present
- **WHEN** the read API returns node data with `NetworkMode = "cross_tenant_eni"`
- **THEN** the system SHALL call `_ = d.Set("network_mode", "cross_tenant_eni")`

#### Scenario: Read with network_mode nil
- **WHEN** the read API returns node data with `NetworkMode = nil`
- **THEN** the system SHALL skip setting `network_mode`

### Requirement: eni_ip computed field
The system SHALL define an `eni_ip` Computed TypeString field on the resource schema. In the Read function, the system SHALL read `EniIP` from the node data returned by the read API and set it on the resource data when the field is not nil.

#### Scenario: Read with eni_ip present
- **WHEN** the read API returns node data with `EniIP = "10.0.0.1"`
- **THEN** the system SHALL call `_ = d.Set("eni_ip", "10.0.0.1")`

#### Scenario: Read with eni_ip nil
- **WHEN** the read API returns node data with `EniIP = nil`
- **THEN** the system SHALL skip setting `eni_ip`

### Requirement: login_settings passed on delete
The system SHALL read the `login_settings` block from the resource data in the Delete function and, when present, build a `dbdcv20201029.LoginSettings{Password, KeyIds, KeepImageLogin}` and set it on the `RemoveNodesFromDBCustomClusterRequest.LoginSettings` before calling the API.

#### Scenario: Delete with login_settings present
- **WHEN** the Delete function executes and the resource data has a `login_settings` block with `password = "xxx"`
- **THEN** the system SHALL set `request.LoginSettings = &LoginSettings{Password: helper.String("xxx")}` on the `RemoveNodesFromDBCustomClusterRequest`

#### Scenario: Delete without login_settings
- **WHEN** the Delete function executes and the resource data has no `login_settings` block
- **THEN** the system SHALL NOT set `LoginSettings` on the `RemoveNodesFromDBCustomClusterRequest`

### Requirement: Backward compatibility
The system SHALL NOT change the type, Required/Optional status, or ForceNew behavior of the existing schema fields (`cluster_id`, `node_id`, `image_id`, `login_settings` and its sub-fields, `node_name`, `lan_ip`, `ssh_endpoint`, `status`, `zone`, `node_type`). The system SHALL NOT change the composite ID format or the resource registration name. All new input arguments SHALL be Optional and ForceNew; all new output fields SHALL be Computed.

#### Scenario: Existing configuration still valid
- **WHEN** a user applies a Terraform configuration that only uses the pre-existing arguments (`cluster_id`, `node_id`, `image_id`, `login_settings`)
- **THEN** the system SHALL accept the configuration and behave identically to before this change

### Requirement: Documentation update
The system SHALL update `resource_tc_dbdc_node_to_db_custom_cluster_attachment.md` to include the new arguments (`labels`, `taints`, `host_name`, `host_name_type`) in the Example Usage section. The Import section SHALL remain unchanged (composite ID). No `Argument Reference` or `Attribute Reference` sections SHALL be added manually (auto-generated).

#### Scenario: Documentation contains new arguments
- **WHEN** the documentation file is generated
- **THEN** the Example Usage SHALL demonstrate usage of the new optional arguments

### Requirement: Unit tests with mocks
The system SHALL extend `resource_tc_dbdc_node_to_db_custom_cluster_attachment_test.go` with gomonkey-mocked unit tests covering the new arguments: verification that Create builds the request with `Labels`, `Taints`, `HostName`, `HostNameType`; Read populates `network_mode` and `eni_ip`; Delete passes `LoginSettings`. Tests SHALL be runnable with `go test -gcflags=all=-l` and SHALL NOT use the Terraform test suite.

#### Scenario: Mock test verifies labels/taints/host_name in Create request
- **WHEN** the Create unit test is executed with gomonkey mocks
- **THEN** the test SHALL assert that the constructed `AddNodesToDBCustomClusterRequest` contains the expected `Labels`, `Taints`, `HostName`, and `HostNameType` values

#### Scenario: Mock test verifies computed fields in Read
- **WHEN** the Read unit test is executed with gomonkey mocks returning node data with `NetworkMode` and `EniIP`
- **THEN** the test SHALL assert that `network_mode` and `eni_ip` are set in the resource data
