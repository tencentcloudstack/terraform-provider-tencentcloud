## Requirements

### Requirement: User can specify instance version at creation
The `tencentcloud_tdmq_professional_cluster` resource SHALL support an optional `instance_version` argument of type `TypeString` that allows users to specify the cluster version when creating a TDMQ professional cluster.

#### Scenario: Create cluster with instance version specified
- **WHEN** user creates a `tencentcloud_tdmq_professional_cluster` resource with `instance_version = "PULSAR_PRO_2.0.0"`
- **THEN** the Create request to `CreateProCluster` API SHALL include `InstanceVersion` set to `"PULSAR_PRO_2.0.0"`

#### Scenario: Create cluster without instance version
- **WHEN** user creates a `tencentcloud_tdmq_professional_cluster` resource without specifying `instance_version`
- **THEN** the Create request to `CreateProCluster` API SHALL NOT include `InstanceVersion`

### Requirement: Instance version is readable after creation
The `tencentcloud_tdmq_professional_cluster` resource SHALL expose the instance version as a computed attribute, readable from the `DescribePulsarProInstances` API response.

#### Scenario: Read instance version from existing cluster
- **WHEN** the provider reads the state of an existing `tencentcloud_tdmq_professional_cluster`
- **THEN** the `instance_version` attribute SHALL be set to the value returned by `DescribePulsarProInstances.Instances[*].InstanceVersion`

#### Scenario: Instance version is nil in API response
- **WHEN** the `DescribePulsarProInstances` API returns an instance where `InstanceVersion` is nil
- **THEN** the `instance_version` attribute SHALL NOT be set (preserving existing state or default)