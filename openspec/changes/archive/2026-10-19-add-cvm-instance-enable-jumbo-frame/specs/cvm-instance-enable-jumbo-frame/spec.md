## ADDED Requirements

### Requirement: Instance supports enable_jumbo_frame parameter

The `tencentcloud_instance` resource SHALL expose an `enable_jumbo_frame` parameter (Optional, Computed, boolean) that controls whether jumbo frames are enabled on the CVM instance. The parameter SHALL be configurable at create time, mutable at update time, and refreshed from the cloud API at read time.

#### Scenario: User enables jumbo frames at creation
- **WHEN** a user creates a `tencentcloud_instance` with `enable_jumbo_frame = true` on a supported instance type
- **THEN** the provider SHALL pass `EnableJumboFrame = true` to the `RunInstances` API, and after creation the Read SHALL set `enable_jumbo_frame` in state from the `DescribeInstances` response

#### Scenario: User omits enable_jumbo_frame at creation
- **WHEN** a user creates a `tencentcloud_instance` without specifying `enable_jumbo_frame`
- **THEN** the provider SHALL NOT send `EnableJumboFrame` in the `RunInstances` request, and the Read SHALL populate `enable_jumbo_frame` in state from the cloud-returned value (Computed)

#### Scenario: User updates enable_jumbo_frame on an existing instance
- **WHEN** a user changes `enable_jumbo_frame` on an existing `tencentcloud_instance` (e.g. from `false` to `true`)
- **THEN** the provider SHALL call `ModifyInstancesAttribute` with the new `EnableJumboFrame` value for that instance, and the subsequent Read SHALL reflect the updated value in state

#### Scenario: Imported instance refreshes enable_jumbo_frame
- **WHEN** an existing `tencentcloud_instance` is imported and read
- **THEN** the provider SHALL set `enable_jumbo_frame` in state from the `Instance.EnableJumboFrame` field returned by `DescribeInstances`
