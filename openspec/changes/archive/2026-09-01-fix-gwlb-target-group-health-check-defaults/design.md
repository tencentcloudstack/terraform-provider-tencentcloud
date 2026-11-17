## Context

The `tencentcloud_gwlb_target_group` resource has a `health_check` nested block with several integer fields (`timeout`, `interval_time`, `health_num`, `un_health_num`) that have documented API defaults. These fields are declared as `Optional: true, Computed: true` in the schema.

When a user specifies a `health_check` block without setting these fields, the flattened map produced by `helper.InterfacesHeadMap` contains the Go `int` zero value (`0`) for each omitted field. The CRUD code then passes these zero values to the GWLB API, which rejects them because they fall outside the valid ranges (e.g., `timeout` must be in [2, 30]).

The GWLB API documentation (in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gwlb/v20240906/models.go`) explicitly states:
- `Timeout`: default 2 seconds, range [2, 30]
- `IntervalTime`: default 5 seconds, range [2, 300]
- `HealthNum`: default 3 times, range [2, 10]
- `UnHealthNum`: default 3 times, range [2, 10]

The SDK request fields are pointers with `omitempty`, so a `nil` field is omitted from the request and the API applies its own default.

## Goals / Non-Goals

**Goals:**
- Ensure that when a user specifies `health_check` but omits `timeout`/`interval_time`/`health_num`/`un_health_num`, the provider does not send `0`; the fields are omitted so the API applies its defaults
- Maintain backward compatibility — the schema is unchanged, and configurations that explicitly set these fields continue to work unchanged

**Non-Goals:**
- Changing the schema — the four fields remain `Optional: true, Computed: true`
- Adding new fields or modifying the API interface
- Changing the `protocol` or `port` fields in `health_check` — they have different semantics (PING protocol doesn't need a port)

## Decisions

### Decision 1: Omit zero values in CRUD code instead of adding schema `Default`

**Rationale**: The four fields are `Computed`, so the schema cannot carry a `Default` (terraform-plugin-sdk/v2 rejects `Default` on a `Computed` field with "Default must be nil if computed"). The correct fix is to stop sending the zero value. Using `d.GetOk("health_check.0.<field>")` returns `false` for an unset integer field, so the field pointer stays `nil` and is omitted from the request (`omitempty`), letting the API apply its own default.

**Alternative considered**: Adding `Default: 2/5/3/3` to the schema. This was rejected because the fields are `Computed`, and `Computed` + `Default` is invalid in terraform-plugin-sdk/v2 (schema validation fails with "Default must be nil if computed").

### Decision 2: Apply the fix to all four fields with API-defined defaults

**Rationale**: All four fields (`timeout`, `interval_time`, `health_num`, `un_health_num`) have the same root cause — if the user sets `health_check` but omits them, `0` gets sent to the API. Fixing all four in one change prevents users from hitting the same class of error for a different field after the `timeout` fix is deployed.

## Risks / Trade-offs

- **[Risk] Users who previously relied on the buggy behavior** (i.e., they set `health_check` and expected the zero-value behavior) → **Mitigation**: The zero-value behavior was always broken and caused API errors, so no user could have been relying on it successfully.
- **[Risk] API default mismatch**: If the API's actual defaults differ from the documented defaults, the value reflected in state (from `Read`) may differ from what the user expects → **Mitigation**: The API defaults are documented and stable; `Read` reflects the API's actual values, which is the source of truth.
