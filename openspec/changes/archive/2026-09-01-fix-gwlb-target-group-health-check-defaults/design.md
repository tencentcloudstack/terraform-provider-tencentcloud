## Context

The `tencentcloud_gwlb_target_group` resource has a `health_check` nested block with several integer fields (`timeout`, `interval_time`, `health_num`, `un_health_num`) that have documented API defaults but no corresponding Terraform `Default` values in the schema.

When a user specifies a `health_check` block without setting these fields, Terraform stores the Go `int` zero value (`0`) for each omitted field. The CRUD code then passes these zero values to the GWLB API, which rejects them because they fall outside the valid ranges (e.g., `timeout` must be in [2, 30]).

The GWLB API documentation (in `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gwlb/v20240906/models.go`) explicitly states:
- `Timeout`: default 2 seconds, range [2, 30]
- `IntervalTime`: default 5 seconds, range [2, 300]
- `HealthNum`: default 3 times, range [2, 10]
- `UnHealthNum`: default 3 times, range [2, 10]

This is a pure Terraform Schema fix — no API changes needed.

## Goals / Non-Goals

**Goals:**
- Ensure that when a user specifies `health_check` but omits `timeout`/`interval_time`/`health_num`/`un_health_num`, the API receives the correct default values instead of `0`
- Maintain backward compatibility — existing configurations that explicitly set these fields continue to work unchanged

**Non-Goals:**
- Changing the CRUD logic (Create/Read/Update/Delete functions) — the schema fix alone resolves the issue
- Adding new fields or modifying the API interface
- Changing the `protocol` or `port` fields in `health_check` — they have different semantics (PING protocol doesn't need a port)

## Decisions

### Decision 1: Use Terraform Schema `Default` instead of code-level defaulting

**Rationale**: Using Terraform's built-in `Default` field in the schema is the idiomatic approach. It ensures:
- The default is applied at plan time, so users see the correct value in `terraform plan`
- The `d.GetOk()` pattern in CRUD code continues to work naturally — when the user doesn't set the field, `GetOk` returns the default value, not the zero value
- No changes to CRUD function logic are needed

**Alternative considered**: Adding `if !ok { targetGroupHealthCheck.Timeout = helper.IntInt64(2) }` in the Create/Update functions. This was rejected because it would be inconsistent with Terraform best practices and would not show the default in `terraform plan`.

### Decision 2: Apply defaults to all four fields with API-defined defaults

**Rationale**: All four fields (`timeout`, `interval_time`, `health_num`, `un_health_num`) have the same root cause — if the user sets `health_check` but omits them, `0` gets sent to the API. Fixing all four in one change prevents users from hitting the same class of error for a different field after the `timeout` fix is deployed.

## Risks / Trade-offs

- **[Risk] Users who previously relied on the buggy behavior** (i.e., they set `health_check` and expected the zero-value behavior) → **Mitigation**: The zero-value behavior was always broken and caused API errors, so no user could have been relying on it successfully.
- **[Risk] State drift after upgrade**: Users who applied the resource with the buggy code may see a diff in `terraform plan` after upgrading to the fixed provider → **Mitigation**: This is expected and desirable — the plan will show the defaults being applied, and `terraform apply` will align the state with the API defaults.