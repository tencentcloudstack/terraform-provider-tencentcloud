## Requirements

### Requirement: Health check fields have correct API-aligned defaults
The `tencentcloud_gwlb_target_group` resource's `health_check` nested block SHALL use Terraform schema `Default` values that match the GWLB API's documented defaults for the following fields:
- `timeout` SHALL default to `2`
- `interval_time` SHALL default to `5`
- `health_num` SHALL default to `3`
- `un_health_num` SHALL default to `3`

#### Scenario: User creates target group with health_check block but no timeout
- **WHEN** a user defines a `health_check` block without setting `timeout`
- **THEN** the Terraform provider SHALL send `timeout=2` to the GWLB API

#### Scenario: User creates target group with health_check block but no interval_time
- **WHEN** a user defines a `health_check` block without setting `interval_time`
- **THEN** the Terraform provider SHALL send `interval_time=5` to the GWLB API

#### Scenario: User creates target group with health_check block but no health_num
- **WHEN** a user defines a `health_check` block without setting `health_num`
- **THEN** the Terraform provider SHALL send `health_num=3` to the GWLB API

#### Scenario: User creates target group with health_check block but no un_health_num
- **WHEN** a user defines a `health_check` block without setting `un_health_num`
- **THEN** the Terraform provider SHALL send `un_health_num=3` to the GWLB API

#### Scenario: User explicitly sets timeout to a custom value
- **WHEN** a user explicitly sets `timeout = 10` in the `health_check` block
- **THEN** the Terraform provider SHALL send `timeout=10` to the GWLB API (user value takes precedence over default)

#### Scenario: User does not specify health_check block at all
- **WHEN** a user creates a target group without a `health_check` block
- **THEN** the Terraform provider SHALL NOT send a `HealthCheck` parameter to the GWLB API (existing behavior preserved)