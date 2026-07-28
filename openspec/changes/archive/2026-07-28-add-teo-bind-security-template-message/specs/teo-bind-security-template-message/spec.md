## ADDED Requirements

### Requirement: Resource exposes Message from DescribeSecurityTemplateBindings
The `tencentcloud_teo_bind_security_template` resource SHALL expose a `message` computed attribute of type `TypeString` in its schema. The value SHALL be sourced from the `Message` field of `EntityStatus` returned by the `DescribeSecurityTemplateBindings` API.

#### Scenario: Message is set during Read
- **WHEN** the `DescribeSecurityTemplateBindings` API returns an `EntityStatus` with a non-nil `Message` field
- **THEN** the resource Read function SHALL set `message` on the resource data to the value of `EntityStatus.Message`

#### Scenario: Message is nil in API response
- **WHEN** the `DescribeSecurityTemplateBindings` API returns an `EntityStatus` with a nil `Message` field
- **THEN** the resource Read function SHALL NOT call `d.Set("message", ...)` (the field retains its prior value or zero value)

#### Scenario: Message is not set by user during Create
- **WHEN** a user creates a `tencentcloud_teo_bind_security_template` resource without specifying `message`
- **THEN** the resource SHALL still be created successfully because `message` is `Computed` only