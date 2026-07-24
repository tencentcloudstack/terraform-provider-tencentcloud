## ADDED Requirements

### Requirement: Resource shall manage the full SCF trigger lifecycle

The `tencentcloud_scf_trigger` resource SHALL manage the complete lifecycle (create, read, update, delete) of a TencentCloud SCF (Serverless Cloud Function) function trigger by mapping to the native SCF APIs `CreateTrigger`, `ListTriggers`, `UpdateTrigger`, and `DeleteTrigger` from the `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416` package.

#### Scenario: Create a trigger
- **WHEN** a user applies a `tencentcloud_scf_trigger` configuration with required fields `function_name`, `trigger_name`, `type`, and `trigger_desc`
- **THEN** the provider SHALL call the SCF `CreateTrigger` API with the supplied parameters and set the Terraform state id to `function_name#namespace#trigger_name` (using `tccommon.FILED_SP` as the separator)

#### Scenario: Read a trigger
- **WHEN** the provider refreshes the resource state
- **THEN** the provider SHALL call the SCF `ListTriggers` API filtered by `FunctionName`, `Namespace`, and a `TriggerName` filter parsed from the composite id, and populate the schema fields from the returned `TriggerInfo`

#### Scenario: Update a trigger
- **WHEN** a user changes a mutable field (`enable`, `qualifier`, `trigger_desc`, `description`, or `custom_argument`) on an existing `tencentcloud_scf_trigger`
- **THEN** the provider SHALL call the SCF `UpdateTrigger` API with the function name, trigger name, type, namespace (parsed from the id) and the changed mutable fields

#### Scenario: Delete a trigger
- **WHEN** a user destroys a `tencentcloud_scf_trigger` resource
- **THEN** the provider SHALL call the SCF `DeleteTrigger` API with `FunctionName`, `TriggerName`, `Type`, `Namespace`, `Qualifier`, and `TriggerDesc` (when applicable) parsed from state

### Requirement: Schema SHALL define trigger parameters and computed attributes

The resource schema SHALL define the following user-settable parameters mapped from the SCF API: `function_name`, `trigger_name`, `type`, `trigger_desc`, `namespace`, `qualifier`, `enable`, `custom_argument`, and `description`. The schema SHALL also define computed read-back attributes from `TriggerInfo`: `available_status`, `add_time`, and `mod_time`.

#### Scenario: Required identity fields are ForceNew
- **WHEN** `function_name`, `trigger_name`, or `type` is changed after creation
- **THEN** the provider SHALL treat the change as a ForceNew (destroy and recreate), because these form the trigger identity and cannot be mutated by `UpdateTrigger`

#### Scenario: Namespace defaults to default and is ForceNew
- **WHEN** `namespace` is omitted in the configuration
- **THEN** the provider SHALL default `namespace` to `"default"` and treat any subsequent change as ForceNew

#### Scenario: enable is a string with OPEN/CLOSE values
- **WHEN** the user sets `enable` to `"OPEN"` or `"CLOSE"`
- **THEN** the provider SHALL pass the string to `CreateTrigger`/`UpdateTrigger` as-is, and on Read SHALL convert the int64 `TriggerInfo.Enable` (1/0) back to `"OPEN"`/`"CLOSE"`

#### Scenario: Computed metadata fields are populated on read
- **WHEN** the Read operation succeeds and the `TriggerInfo` contains `AvailableStatus`, `AddTime`, or `ModTime`
- **THEN** the provider SHALL set the computed fields `available_status`, `add_time`, and `mod_time` respectively, only when the API field is non-nil

### Requirement: Composite id SHALL be used and importable

The resource SHALL use a composite id composed of `function_name`, `namespace`, and `trigger_name` joined by `tccommon.FILED_SP` (`#`). The resource SHALL support Terraform import via `schema.ImportStatePassthrough`.

#### Scenario: Import by composite id
- **WHEN** a user runs `terraform import tencentcloud_scf_trigger.foo <function_name>#<namespace>#<trigger_name>`
- **THEN** the provider SHALL parse the three-part id and run the Read operation to populate state

#### Scenario: Broken id is rejected
- **WHEN** the composite id does not contain exactly three `#`-separated parts during Read/Update/Delete
- **THEN** the provider SHALL return an error indicating the id is broken

### Requirement: Create SHALL validate non-empty API response

The Create operation SHALL verify, after a successful `CreateTrigger` call, that the API response and its `TriggerInfo` field are non-nil and non-empty. If the response is empty, the provider SHALL return a `NonRetryableError` to avoid persisting a broken id.

#### Scenario: Empty create response fails
- **WHEN** `CreateTrigger` returns a nil response or a nil/empty `TriggerInfo`
- **THEN** the provider SHALL return a `NonRetryableError` and SHALL NOT set the Terraform state id

### Requirement: Read SHALL handle resource disappearance

When the `ListTriggers` response returns an empty list (trigger no longer exists), the Read operation SHALL log the disappearance with the current id and then call `d.SetId("")` so Terraform removes the resource from state.

#### Scenario: Trigger not found on read
- **WHEN** `ListTriggers` returns no triggers matching the composite id
- **THEN** the provider SHALL log `[CRUD] scf_trigger id=<id>` and set the resource id to empty string

### Requirement: API calls SHALL use retry and timeout helpers

All SCF API calls in Create, Read, Update, and Delete SHALL be wrapped with the provider's retry helpers (`resource.Retry` with `tccommon.WriteRetryTimeout` for write operations and `tccommon.ReadRetryTimeout` for read operations). API errors SHALL be wrapped with `tccommon.RetryError()`.

#### Scenario: Transient API failure is retried
- **WHEN** an SCF API call fails with a retryable error
- **THEN** the provider SHALL retry within the configured timeout before surfacing the error to the user

### Requirement: Resource SHALL be registered in the provider

The `tencentcloud_scf_trigger` resource SHALL be registered in `tencentcloud/provider.go` under the provider resources map and documented in `tencentcloud/provider.md`.

#### Scenario: Resource is available in the provider
- **WHEN** a user references `tencentcloud_scf_trigger` in a Terraform configuration
- **THEN** the provider SHALL recognize the resource type and expose its schema and CRUD operations
