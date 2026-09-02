## ADDED Requirements

### Requirement: Health check fields omit zero values so the API applies its defaults
The `tencentcloud_gwlb_target_group` resource's `health_check` nested block SHALL NOT send zero values for omitted fields. When a user omits `timeout`, `interval_time`, `health_num`, or `un_health_num`, the provider SHALL omit the corresponding request field, allowing the GWLB API to apply its documented defaults (`timeout=2`, `interval_time=5`, `health_num=3`, `un_health_num=3`).

#### Scenario: User creates target group with health_check block but no timeout
- **WHEN** a user defines a `health_check` block without setting `timeout`
- **THEN** the Terraform provider SHALL NOT send `timeout=0`; it SHALL omit `Timeout` from the request so the GWLB API applies its default `timeout=2`

#### Scenario: User creates target group with health_check block but no interval_time
- **WHEN** a user defines a `health_check` block without setting `interval_time`
- **THEN** the Terraform provider SHALL NOT send `interval_time=0`; it SHALL omit `IntervalTime` from the request so the GWLB API applies its default `interval_time=5`

#### Scenario: User creates target group with health_check block but no health_num
- **WHEN** a user defines a `health_check` block without setting `health_num`
- **THEN** the Terraform provider SHALL NOT send `health_num=0`; it SHALL omit `HealthNum` from the request so the GWLB API applies its default `health_num=3`

#### Scenario: User creates target group with health_check block but no un_health_num
- **WHEN** a user defines a `health_check` block without setting `un_health_num`
- **THEN** the Terraform provider SHALL NOT send `un_health_num=0`; it SHALL omit `UnHealthNum` from the request so the GWLB API applies its default `un_health_num=3`

#### Scenario: User explicitly sets timeout to a custom value
- **WHEN** a user explicitly sets `timeout = 10` in the `health_check` block
- **THEN** the Terraform provider SHALL send `timeout=10` to the GWLB API (user value takes precedence over default)

#### Scenario: User does not specify health_check block at all
- **WHEN** a user creates a target group without a `health_check` block
- **THEN** the Terraform provider SHALL NOT send a `HealthCheck` parameter to the GWLB API (existing behavior preserved)
