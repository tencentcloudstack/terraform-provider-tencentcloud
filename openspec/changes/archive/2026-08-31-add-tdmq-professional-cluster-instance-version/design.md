## Context

The `tencentcloud_tdmq_professional_cluster` resource manages TDMQ Pulsar professional clusters. The cloud API `CreateProCluster` already supports an `InstanceVersion` parameter to specify the cluster version at creation time, and `DescribePulsarProInstances` returns the `InstanceVersion` in the `PulsarProInstance` response. The Terraform resource currently does not expose this parameter. This is a straightforward addition of an optional+computed field.

The vendor SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217`) already has the field in both request and response structs, confirmed via vendor directory inspection.

## Goals / Non-Goals

**Goals:**
- Add `instance_version` as an optional + computed string field to the `tencentcloud_tdmq_professional_cluster` resource
- Pass the value through to `CreateProCluster` API when specified
- Read the value back from `DescribePulsarProInstances` API response
- Maintain full backward compatibility

**Non-Goals:**
- Modifying or exposing any other TDMQ parameters
- Adding validation logic for the version string (the API handles validation)
- Changing the resource's update/delete behavior

## Decisions

1. **Schema type: `TypeString`, Optional + Computed**: The parameter is optional (user may or may not specify it at creation) and computed (the API may return a default value). This is the standard pattern for parameters that are passed at creation and read back.

2. **Not adding to immutableArgs**: The `instance_version` parameter is only settable at creation time (the `ModifyCluster` API does not support changing it). However, since the field is not in the Update method's request (`ModifyClusterRequest`), we don't need to add it to `immutableArgs` — the user simply cannot change it via Update. Adding it to `immutableArgs` would produce a confusing error when the user hasn't actually modified it.

3. **Read from `professionalCluster.InstanceVersion`**: The Read method already calls `DescribePulsarProInstances` which returns a `PulsarProInstance` with `InstanceVersion`. We read from this field, consistent with how other fields like `SpecName`, `InstanceName`, `MaxStorage`, `AutoRenewFlag`, `VpcId`, `SubnetId` are read.

4. **Create: set `request.InstanceVersion`**: In the Create method, we add the same pattern as other fields — use `d.GetOk("instance_version")` and set `request.InstanceVersion`.

## Risks / Trade-offs

- **Risk**: The `InstanceVersion` field in `PulsarProInstance` may be `nil` for existing clusters created before this field was added to the API. → **Mitigation**: Check for `nil` before setting in Read, consistent with existing patterns (e.g., `if professionalCluster.InstanceVersion != nil`).
- **Risk**: Users may specify an invalid version string. → **Mitigation**: The API validates this and returns an error. No client-side validation needed.