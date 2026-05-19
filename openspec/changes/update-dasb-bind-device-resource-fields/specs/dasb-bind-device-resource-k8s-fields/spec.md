## ADDED Requirements

### Requirement: Schema exposes K8S managed account dimension fields
The resource `tencentcloud_dasb_bind_device_resource` SHALL define the following optional fields in its schema, corresponding to `BindDeviceResource` API parameters. These fields SHALL be wired to the API call in Create and Update once the Go SDK is upgraded to include them in `BindDeviceResourceRequest`.

- `manage_dimension` (TypeInt, Optional): K8S cluster managed account dimension. Valid values: `1`-Cluster, `2`-Namespace, `3`-Workload.
- `manage_account_id` (TypeInt, Optional): K8S cluster managed account ID.
- `manage_account` (TypeString, Optional): K8S cluster managed account name.
- `manage_kubeconfig` (TypeString, Optional, Sensitive): K8S cluster managed account kubeconfig credential.
- `namespace` (TypeString, Optional): K8S cluster managed namespace.
- `workload` (TypeString, Optional): K8S cluster managed workload.

#### Scenario: User configures K8S dimension fields
- **WHEN** a user sets `manage_dimension = 2`, `namespace = "default"` in their configuration
- **THEN** the provider SHALL include these values in the `BindDeviceResource` API request after SDK upgrade

#### Scenario: K8S fields are omitted
- **WHEN** a user does not specify any K8S dimension fields
- **THEN** the provider SHALL call `BindDeviceResource` without those parameters, preserving existing behaviour

### Requirement: Schema exposes domain_name as computed read field
The resource SHALL expose `domain_name` as a Computed `TypeString` field, populated from the `DomainName` field of the first `Device` returned by `DescribeDevices`.

#### Scenario: domain_name is returned by API
- **WHEN** `DescribeDevices` returns devices with a non-nil `DomainName`
- **THEN** the provider SHALL set `domain_name` in state from the first device's value

#### Scenario: domain_name is nil
- **WHEN** `DescribeDevices` returns devices with `DomainName = nil`
- **THEN** the provider SHALL leave `domain_name` unset in state without error

### Requirement: Read correctly populates device_id_set with integer values
The Read function SHALL dereference `*uint64` pointers from `Device.Id` and convert them to `int` before appending to the `device_id_set` state value.

#### Scenario: Multiple devices are bound to a resource
- **WHEN** `DescribeDevices` returns N devices for a given `ResourceIdSet`
- **THEN** `device_id_set` in state SHALL contain exactly N integer values equal to each device's `Id`

### Requirement: Read sets domain_id once from the first matching device
The Read function SHALL set `domain_id` from the first device with a non-nil `DomainId` and SHALL NOT overwrite it in subsequent loop iterations.

#### Scenario: All devices share the same domain_id
- **WHEN** `DescribeDevices` returns multiple devices all with the same `DomainId`
- **THEN** `domain_id` in state SHALL be set to that value exactly once

#### Scenario: No device has a domain_id
- **WHEN** all returned devices have `DomainId = nil`
- **THEN** `domain_id` SHALL remain as previously stored in state
