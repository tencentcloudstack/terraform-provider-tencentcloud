# AGS sandbox environment module

Pre-release module for the AGS Tool Resource in this checkout. Terraform 1.17.0 is the initial validated CLI target. There is no claimed published module address or minimum released AGS Provider version. Use the isolated local candidate installer supplied by the E2E harness, then pin the actual accepted release when upstream publishes it.

The module manages one Tool and references existing network, role, storage, log and image resources. It never discovers or imports dependencies by name, builds images, starts instances, or deletes storage data. Callers who want new dependencies must declare and own those resources separately and decide their retention policies.

Required inputs are `tool_name`, `tool_type`, and `network`. `custom_configuration` is required only for custom Tools. `description`, `default_timeout`, `persistent`, `role_arn`, `tags`, `storage_mounts`, `log_configuration`, and `timeouts` follow the Resource schema. Omitted Optional/Computed scalar settings accept remote defaults; an explicit empty collection requests reset. Environment values are stored in state and must not contain secrets.

Outputs are `tool_id`, `tool_name`, `status`, and `image_digest`. Applications consume `tool_id` with their runtime SDK. Active/paused instances must be explicitly released by their runtime owner before Tool destruction. Updating Tool configuration does not roll existing instances.

PUBLIC/SANDBOX/VPC are supported subject to the verified API matrix. VPC bindings, role, storage, log, default lifetime and persistence changes replace the Tool in V1; a same-name replacement may require downtime. Description clearing can fail to converge in the current service. Explicit read errors retain state; established managed state follows the Provider convention that a complete successful empty exact-ID query confirms absence. Do not treat successful schema validation as complete cloud acceptance.


Dependency outputs may be unknown during plan, including individual subnet IDs,
security-group IDs and log filenames in known collections. The current local
candidate validates known entries immediately and defers unknown entries until
apply; unresolved values cannot reach the SDK. See `../../e2e/README.md` for local validation and the independent SDK handoff.
