## Requirements

### Requirement: bot_management_lite schema definition in web_security_template
The `tencentcloud_teo_web_security_template` resource SHALL include an Optional `bot_management_lite` parameter (TypeList, MaxItems: 1) under the `security_policy` block, containing two sub-blocks:
- `captcha_page_challenge` (TypeList, MaxItems: 1, Optional, Computed): CAPTCHA page challenge configuration
- `ai_crawler_detection` (TypeList, MaxItems: 1, Optional, Computed): AI crawler detection configuration

#### Scenario: Resource accepts bot_management_lite configuration
- **WHEN** a user provides a `bot_management_lite` block in the `security_policy` of `tencentcloud_teo_web_security_template`
- **THEN** the resource SHALL accept and process the configuration without errors

#### Scenario: bot_management_lite is optional
- **WHEN** a user creates a `tencentcloud_teo_web_security_template` resource without specifying `bot_management_lite`
- **THEN** the resource SHALL be created successfully and the existing behavior SHALL be preserved

### Requirement: captcha_page_challenge configuration in web_security_template
The `captcha_page_challenge` block under `bot_management_lite` SHALL contain a single field:
- `enabled` (TypeString, Required): Whether CAPTCHA page challenge is enabled. Valid values: "on", "off".

#### Scenario: CAPTCHA page challenge enabled
- **WHEN** a user configures `captcha_page_challenge` with `enabled = "on"` in `tencentcloud_teo_web_security_template`
- **THEN** the resource SHALL set `BotManagementLite.CAPTCHAPageChallenge.Enabled` to "on" in the CreateWebSecurityTemplate/ModifyWebSecurityTemplate API request

#### Scenario: CAPTCHA page challenge disabled
- **WHEN** a user configures `captcha_page_challenge` with `enabled = "off"` in `tencentcloud_teo_web_security_template`
- **THEN** the resource SHALL set `BotManagementLite.CAPTCHAPageChallenge.Enabled` to "off" in the CreateWebSecurityTemplate/ModifyWebSecurityTemplate API request

### Requirement: ai_crawler_detection configuration in web_security_template
The `ai_crawler_detection` block under `bot_management_lite` SHALL contain:
- `enabled` (TypeString, Required): Whether AI crawler detection is enabled. Valid values: "on", "off".
- `action` (TypeList, MaxItems: 1, Optional): Execution action when Enabled is "on".

#### Scenario: AI crawler detection enabled with action
- **WHEN** a user configures `ai_crawler_detection` with `enabled = "on"` and an `action` block in `tencentcloud_teo_web_security_template`
- **THEN** the resource SHALL set `BotManagementLite.AICrawlerDetection.Enabled` to "on" and set the Action in the CreateWebSecurityTemplate/ModifyWebSecurityTemplate API request

#### Scenario: AI crawler detection disabled
- **WHEN** a user configures `ai_crawler_detection` with `enabled = "off"` in `tencentcloud_teo_web_security_template`
- **THEN** the resource SHALL set `BotManagementLite.AICrawlerDetection.Enabled` to "off" in the CreateWebSecurityTemplate/ModifyWebSecurityTemplate API request

### Requirement: ai_crawler_detection action schema in web_security_template
The `action` block under `ai_crawler_detection` SHALL contain:
- `name` (TypeString, Required): The security action name. Valid values: "Deny", "Monitor", "Allow", "Challenge", "Redirect", "BlockIP", "ReturnCustomPage", "JSChallenge", "ManagedChallenge".
- `deny_action_parameters` (TypeList, MaxItems: 1, Optional): Additional parameters when Name is "Deny".
- `allow_action_parameters` (TypeList, MaxItems: 1, Optional): Additional parameters when Name is "Allow".
- `challenge_action_parameters` (TypeList, MaxItems: 1, Optional): Additional parameters when Name is "Challenge".
- `redirect_action_parameters` (TypeList, MaxItems: 1, Optional): Additional parameters when Name is "Redirect".
- `block_ip_action_parameters` (TypeList, MaxItems: 1, Optional): Additional parameters when Name is "BlockIP".
- `return_custom_page_action_parameters` (TypeList, MaxItems: 1, Optional): Additional parameters when Name is "ReturnCustomPage".

#### Scenario: Action with Deny
- **WHEN** a user configures `action` with `name = "Deny"` and `deny_action_parameters`
- **THEN** the resource SHALL set `SecurityAction.Name` to "Deny" and populate `SecurityAction.DenyActionParameters` in the API request

#### Scenario: Action with Allow
- **WHEN** a user configures `action` with `name = "Allow"` and `allow_action_parameters`
- **THEN** the resource SHALL set `SecurityAction.Name` to "Allow" and populate `SecurityAction.AllowActionParameters` in the API request

#### Scenario: Action with Challenge
- **WHEN** a user configures `action` with `name = "Challenge"` and `challenge_action_parameters`
- **THEN** the resource SHALL set `SecurityAction.Name` to "Challenge" and populate `SecurityAction.ChallengeActionParameters` in the API request

### Requirement: deny_action_parameters schema for ai_crawler_detection in web_security_template
The `deny_action_parameters` block under `action` SHALL contain:
- `block_ip` (TypeString, Optional): Whether to extend IP blocking. Valid values: "on", "off".
- `block_ip_duration` (TypeString, Optional): IP blocking duration when block_ip is "on".
- `return_custom_page` (TypeString, Optional): Whether to use custom pages. Valid values: "on", "off".
- `response_code` (TypeString, Optional): Custom page status code.
- `error_page_id` (TypeString, Optional): Custom page ID.
- `stall` (TypeString, Optional): Whether to suspend request processing. Valid values: "on", "off".

#### Scenario: Deny action with IP blocking
- **WHEN** a user configures `deny_action_parameters` with `block_ip = "on"` and `block_ip_duration = "120s"`
- **THEN** the resource SHALL set `DenyActionParameters.BlockIp` to "on" and `DenyActionParameters.BlockIpDuration` to "120s" in the API request

### Requirement: allow_action_parameters schema for ai_crawler_detection in web_security_template
The `allow_action_parameters` block under `action` SHALL contain:
- `min_delay_time` (TypeString, Optional): Minimum delay response time. Supported unit: seconds, range 0-5.
- `max_delay_time` (TypeString, Optional): Maximum delay response time. Supported unit: seconds, range 5-10.

#### Scenario: Allow action with delay configuration
- **WHEN** a user configures `allow_action_parameters` with `min_delay_time = "0s"` and `max_delay_time = "5s"`
- **THEN** the resource SHALL set `AllowActionParameters.MinDelayTime` to "0s" and `AllowActionParameters.MaxDelayTime` to "5s" in the API request

### Requirement: challenge_action_parameters schema for ai_crawler_detection in web_security_template
The `challenge_action_parameters` block under `action` SHALL contain:
- `challenge_option` (TypeString, Required): Specific challenge action. Valid values: "JSChallenge", "ManagedChallenge".
- `interval` (TypeString, Optional): Time interval for repeating the challenge.
- `attester_id` (TypeString, Optional): Client authentication method ID.

#### Scenario: Challenge action with JSChallenge
- **WHEN** a user configures `challenge_action_parameters` with `challenge_option = "JSChallenge"`
- **THEN** the resource SHALL set `ChallengeActionParameters.ChallengeOption` to "JSChallenge" in the API request

### Requirement: Read operation for bot_management_lite in web_security_template
The resource Read operation SHALL populate `bot_management_lite` from the DescribeWebSecurityTemplate API response, extracting `SecurityPolicy.BotManagementLite` and flattening its nested structures (`CAPTCHAPageChallenge`, `AICrawlerDetection`) into the Terraform state.

#### Scenario: Read with BotManagementLite present
- **WHEN** the DescribeWebSecurityTemplate response contains a non-nil `BotManagementLite` field
- **THEN** the resource SHALL flatten the `BotManagementLite` data into the `bot_management_lite` state attribute

#### Scenario: Read with BotManagementLite absent
- **WHEN** the DescribeWebSecurityTemplate response has a nil `BotManagementLite` field
- **THEN** the resource SHALL not set the `bot_management_lite` attribute (it remains empty)

### Requirement: Nil checks in Read operation for web_security_template
The Read operation for `bot_management_lite` SHALL check for nil at each nesting level before accessing sub-fields:
- Check `respData.BotManagementLite != nil` before accessing its sub-fields
- Check `respData.BotManagementLite.CAPTCHAPageChallenge != nil` before accessing `Enabled`
- Check `respData.BotManagementLite.AICrawlerDetection != nil` before accessing its sub-fields
- Check `respData.BotManagementLite.AICrawlerDetection.Action != nil` before accessing action sub-fields
- Check each action parameter block for nil before accessing its fields

#### Scenario: Partial BotManagementLite response
- **WHEN** the DescribeWebSecurityTemplate response has `BotManagementLite` with only `CAPTCHAPageChallenge` set (and `AICrawlerDetection` is nil)
- **THEN** the resource SHALL only flatten the `captcha_page_challenge` data and not set `ai_crawler_detection`

### Requirement: Create and Update operations for bot_management_lite in web_security_template
The resource Create and Update operations SHALL expand the `bot_management_lite` Terraform configuration into the `SecurityPolicy.BotManagementLite` field of the CreateWebSecurityTemplate/ModifyWebSecurityTemplate API request, including all nested structures.

#### Scenario: Create with bot_management_lite
- **WHEN** a user creates a resource with `bot_management_lite` configuration
- **THEN** the Create operation SHALL expand the configuration and include `BotManagementLite` in the CreateWebSecurityTemplate request

#### Scenario: Update with bot_management_lite changes
- **WHEN** a user updates the resource's `bot_management_lite` configuration
- **THEN** the Update operation SHALL expand the new configuration and include `BotManagementLite` in the ModifyWebSecurityTemplate request