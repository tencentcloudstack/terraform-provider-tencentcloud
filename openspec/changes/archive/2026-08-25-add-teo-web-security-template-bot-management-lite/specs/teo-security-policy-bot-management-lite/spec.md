## MODIFIED Requirements

### Requirement: bot_management_lite schema definition
The `tencentcloud_teo_security_policy_config` resource SHALL include an Optional `bot_management_lite` parameter (TypeList, MaxItems: 1) under the `security_policy` block, containing two sub-fields:
- `captcha_page_challenge` (TypeList, MaxItems: 1, Optional): CAPTCHA page challenge configuration
- `ai_crawler_detection` (TypeList, MaxItems: 1, Optional): AI crawler detection configuration

The `tencentcloud_teo_web_security_template` resource SHALL also include the same `bot_management_lite` parameter (TypeList, MaxItems: 1) under its `security_policy` block, with identical sub-fields and behavior, but using the CreateWebSecurityTemplate/ModifyWebSecurityTemplate/DescribeWebSecurityTemplate APIs instead of ModifySecurityPolicy/DescribeSecurityPolicy.

#### Scenario: Resource accepts bot_management_lite configuration
- **WHEN** a user provides a `bot_management_lite` block in the `security_policy` of `tencentcloud_teo_security_policy_config` or `tencentcloud_teo_web_security_template`
- **THEN** the resource SHALL accept and process the configuration without errors

#### Scenario: bot_management_lite is optional
- **WHEN** a user creates a `tencentcloud_teo_security_policy_config` or `tencentcloud_teo_web_security_template` resource without specifying `bot_management_lite`
- **THEN** the resource SHALL be created successfully and the existing behavior SHALL be preserved

## ADDED Requirements

### Requirement: bot_management_lite in web_security_template APIs
The `tencentcloud_teo_web_security_template` resource SHALL use the following APIs for `bot_management_lite`:
- CreateWebSecurityTemplate: Set `request.SecurityPolicy.BotManagementLite` with CAPTCHAPageChallenge and AICrawlerDetection
- ModifyWebSecurityTemplate: Set `request.SecurityPolicy.BotManagementLite` with CAPTCHAPageChallenge and AICrawlerDetection
- DescribeWebSecurityTemplate: Read `response.SecurityPolicy.BotManagementLite` and flatten to state

#### Scenario: CreateWebSecurityTemplate sends BotManagementLite
- **WHEN** a user creates a `tencentcloud_teo_web_security_template` with `bot_management_lite` configuration
- **THEN** the CreateWebSecurityTemplate API request SHALL include `SecurityPolicy.BotManagementLite` with the configured values

#### Scenario: DescribeWebSecurityTemplate reads BotManagementLite
- **WHEN** the DescribeWebSecurityTemplate API returns `SecurityPolicy.BotManagementLite`
- **THEN** the Read function SHALL flatten the response into the Terraform state under `security_policy.0.bot_management_lite`