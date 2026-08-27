## ADDED Requirements

### Requirement: Manage Redis instance password policy
The system SHALL provide a Terraform resource `tencentcloud_redis_instance_password_policy` that allows users to declaratively manage the password complexity policy for a Redis instance.

#### Scenario: Read existing password policy
- **WHEN** a user applies a `tencentcloud_redis_instance_password_policy` resource with an existing `instance_id`
- **THEN** the system calls `DescribeInstancePasswordPolicy` and SHALL populate all fields including `enabled`, `min_letter_count`, `min_digit_count`, `min_special_count`, and `min_length` from the API response

#### Scenario: Update password policy
- **WHEN** a user changes any field in the `password_policy` block
- **THEN** the system SHALL call `ModifyInstancePasswordPolicy` with the full password policy configuration and SHALL refresh the state by calling `DescribeInstancePasswordPolicy`

#### Scenario: Enable password complexity
- **WHEN** a user sets `enabled` to `true` and provides values for `min_letter_count`, `min_digit_count`, `min_special_count`, and `min_length`
- **THEN** the system SHALL call `ModifyInstancePasswordPolicy` and the password complexity rules SHALL be enforced for all subsequent password changes on the Redis instance

#### Scenario: Disable password complexity
- **WHEN** a user sets `enabled` to `false`
- **THEN** the system SHALL call `ModifyInstancePasswordPolicy` with `Enabled: false` and the password complexity rules SHALL no longer be enforced

#### Scenario: Import existing policy
- **WHEN** a user imports a `tencentcloud_redis_instance_password_policy` resource using the Redis instance ID
- **THEN** the system SHALL read the existing password policy and populate the Terraform state

#### Scenario: Instance not found
- **WHEN** the `DescribeInstancePasswordPolicy` API returns `ResourceNotFound.InstanceNotExists`
- **THEN** the system SHALL clear the resource ID from state, allowing Terraform to detect the drift