## Context

The `tencentcloud_mariadb_hour_db_instance` resource manages MariaDB hourly-billed instances. Currently, when creating an instance, the resource calls `InitDBInstances` with hardcoded initialization parameters (character_set_server=utf8mb4, lower_case_table_names=1, sync_mode=2, innodb_page_size=16384). The `CreateHourDBInstance` API also supports an `InitParams` field of type `[]*DBParamValue` that allows passing initialization parameters at creation time.

## Goals / Non-Goals

**Goals:**
- Allow users to configure database initialization parameters via a new `init_params` schema field
- Maintain backward compatibility: existing configurations without `init_params` continue to work with the same hardcoded defaults
- Pass user-provided `init_params` through the `CreateHourDBInstance` API's `InitParams` field

**Non-Goals:**
- Supporting update of init_params after creation (these are one-time initialization settings)
- Reading init_params back from the API (DescribeDBInstanceDetail does not return them)
- Removing the existing `InitDBInstances` call fallback for when `init_params` is not specified

## Decisions

1. **Schema type: TypeList of objects with `param` and `value` string fields**
   - Rationale: Maps directly to the `[]*DBParamValue` structure in the cloud API. Each element has a `Param` (parameter name) and `Value` (parameter value).
   - Alternative considered: TypeMap — rejected because it loses ordering and doesn't match the API's list-of-objects structure.

2. **ForceNew behavior**
   - Rationale: Init params can only be set at instance creation time. Changing them requires recreating the instance. This is consistent with the API behavior.

3. **Pass init_params via CreateHourDBInstance request.InitParams field**
   - Rationale: The `CreateHourDBInstance` API natively supports `InitParams`. When user provides `init_params`, set them on the create request directly. The separate `InitDBInstances` call with hardcoded values is only used when user does NOT provide `init_params`, preserving backward compatibility.

4. **No Read-back in Read method**
   - Rationale: The `DescribeDBInstanceDetail` API does not return initialization parameters. The field will be stored in Terraform state from the config only.

## Risks / Trade-offs

- [Risk] Users may provide invalid parameter names or values → The cloud API will return an error, which will be surfaced to the user via Terraform's error handling.
- [Risk] Backward compatibility concern if existing users import the resource → Since `init_params` is Optional and not read back from API, imported resources will simply not have this field set in state, which is acceptable.
