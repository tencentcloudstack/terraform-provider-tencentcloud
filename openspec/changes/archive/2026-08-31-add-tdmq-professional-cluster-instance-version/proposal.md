## Why

The TDMQ professional cluster (`tencentcloud_tdmq_professional_cluster`) resource currently does not expose the `InstanceVersion` parameter. The cloud API `CreateProCluster` already supports specifying an instance version at creation time, and `DescribePulsarProInstances` returns the instance version in its response. Adding this parameter allows users to specify and read the cluster version, enabling better version management and visibility.

## What Changes

- Add a new optional + computed `instance_version` string field to the `tencentcloud_tdmq_professional_cluster` resource schema
- Wire the parameter in the Create method to pass through to `CreateProCluster` API
- Read the parameter back from `DescribePulsarProInstances` API response in the Read method
- Add corresponding documentation in the `.md` resource doc file
- Add unit test coverage for the new parameter

## Capabilities

### New Capabilities
- `tdmq-professional-cluster-instance-version`: Support for specifying and reading the instance version of TDMQ professional clusters

### Modified Capabilities
<!-- No existing capabilities are being modified at the spec level -->

## Impact

- **Affected code**: `tencentcloud/services/tpulsar/resource_tc_tdmq_professional_cluster.go`, `tencentcloud/services/tpulsar/resource_tc_tdmq_professional_cluster_test.go`, `tencentcloud/services/tpulsar/resource_tc_tdmq_professional_cluster.md`
- **Affected APIs**: `CreateProCluster` (tdmq v20200217), `DescribePulsarProInstances` (tdmq v20200217)
- **Backward compatibility**: Fully backward compatible — new field is optional and computed
- **Dependencies**: No new dependencies; uses existing `tencentcloud-sdk-go/tencentcloud/tdmq/v20200217` package