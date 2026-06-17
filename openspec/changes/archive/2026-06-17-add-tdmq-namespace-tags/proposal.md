## Why

The `tencentcloud_tdmq_namespace` resource currently does not support tagging namespaces during creation. The TDMQ `CreateEnvironment` API already supports a `Tags` parameter, and the `DescribeEnvironments` API returns tags in the `Environment` response. Users need the ability to assign tags to TDMQ namespaces at creation time for resource organization and management purposes.

## What Changes

- Add a `Tags` parameter (TypeList of key/value pairs) to the `tencentcloud_tdmq_namespace` resource schema
- The `Tags` parameter will be Optional + ForceNew (since `ModifyEnvironmentAttributes` API does not support updating tags)
- Pass `Tags` to the `CreateEnvironment` API request during resource creation
- Read `Tags` from the `DescribeEnvironments` API response during resource read
- Update the service layer `CreateTdmqNamespace` function to accept and pass Tags
- Update the resource `.md` documentation example

## Capabilities

### New Capabilities
- `tdmq-namespace-tags`: Add Tags parameter support to tencentcloud_tdmq_namespace resource, enabling tag assignment during namespace creation

### Modified Capabilities

## Impact

- `tencentcloud/services/tpulsar/resource_tc_tdmq_namespace.go` - Add Tags schema, create/read logic
- `tencentcloud/services/tdmq/service_tencentcloud_tdmq.go` - Update `CreateTdmqNamespace` to accept Tags parameter
- `tencentcloud/services/tpulsar/resource_tc_tdmq_namespace.md` - Update documentation with Tags example
- `tencentcloud/services/tpulsar/resource_tc_tdmq_namespace_test.go` - Add unit tests for Tags parameter
