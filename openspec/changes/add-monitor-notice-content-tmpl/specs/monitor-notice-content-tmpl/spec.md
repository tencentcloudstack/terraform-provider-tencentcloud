# Spec: Monitor Notice Content Template Resource

## ADDED Requirements

### Requirement: Resource Registration
The system SHALL register a new Terraform resource named `tencentcloud_monitor_notice_content_tmpl` in the provider.

#### Scenario: Provider initialization
- **WHEN** the Terraform provider initializes
- **THEN** the resource `tencentcloud_monitor_notice_content_tmpl` SHALL be available for use

---

### Requirement: Create Template
The system SHALL create a notification content template using the Monitor API.

#### Scenario: Successful template creation
- **GIVEN** valid template parameters (name, monitor type, language, contents)
- **WHEN** user applies Terraform configuration with the resource
- **THEN** the system SHALL call `CreateNoticeContentTmpl` API
- **AND** store the returned `TmplID` in resource state
- **AND** set resource ID as `tmplID#tmplName` composite format

#### Scenario: Missing required fields
- **GIVEN** template configuration with missing required fields
- **WHEN** user applies Terraform configuration
- **THEN** the system SHALL return validation error before API call

---

### Requirement: Read Template
The system SHALL query existing template configuration using the Monitor API.

#### Scenario: Successful template read
- **GIVEN** a resource with valid composite ID `tmplID#tmplName`
- **WHEN** Terraform performs refresh or read operation
- **THEN** the system SHALL parse the composite ID
- **AND** call `DescribeNoticeContentTmpl` with `TmplIDs` array and `TmplName`
- **AND** populate resource state with returned data

#### Scenario: Template not found
- **GIVEN** a resource ID for a deleted template
- **WHEN** Terraform performs read operation
- **THEN** the system SHALL detect empty response
- **AND** clear resource ID from state
- **AND** log warning message

#### Scenario: Malformed ID
- **GIVEN** a resource with malformed ID (missing separator)
- **WHEN** Terraform performs read operation
- **THEN** the system SHALL return error "id is broken"

---

### Requirement: Update Template
The system SHALL update existing template configuration using the Monitor API.

#### Scenario: Update template contents
- **GIVEN** existing template resource
- **WHEN** user modifies `tmpl_contents` in configuration
- **THEN** the system SHALL detect changes using `d.HasChange()`
- **AND** call `ModifyNoticeContentTmpl` with updated contents
- **AND** refresh resource state

#### Scenario: Update immutable fields
- **GIVEN** existing template resource
- **WHEN** user modifies ForceNew fields (`tmpl_name`, `monitor_type`, `tmpl_language`)
- **THEN** Terraform SHALL destroy and recreate the resource

---

### Requirement: Delete Template
The system SHALL delete template using the Monitor API.

#### Scenario: Successful template deletion
- **GIVEN** existing template resource
- **WHEN** user removes resource from configuration
- **THEN** the system SHALL parse composite ID to extract `tmplID`
- **AND** call `DeleteNoticeContentTmpls` with `TmplIDs` array
- **AND** confirm deletion success

---

### Requirement: Import Template
The system SHALL support importing existing templates into Terraform state.

#### Scenario: Import by composite ID
- **GIVEN** existing template in Tencent Cloud with known `tmplID` and `tmplName`
- **WHEN** user runs `terraform import tencentcloud_monitor_notice_content_tmpl.example "ntpl-xxx#TemplateName"`
- **THEN** the system SHALL parse the composite ID
- **AND** query template details
- **AND** populate full resource state

---

### Requirement: Schema Definition
The system SHALL define comprehensive schema for all template fields.

#### Scenario: Required fields validation
- **GIVEN** resource configuration
- **WHEN** user omits required fields (`tmpl_name`, `monitor_type`, `tmpl_language`, `tmpl_contents`)
- **THEN** Terraform SHALL return validation error

#### Scenario: Complex nested structure support
- **GIVEN** template with multiple notification channels
- **WHEN** user defines nested `tmpl_contents` with channels like `we_work_robot`, `ding_ding_robot`
- **THEN** the system SHALL correctly marshal and unmarshal nested structures
- **AND** preserve all channel configurations

---

### Requirement: Service Layer Integration
The system SHALL implement service layer methods for template operations.

#### Scenario: Service method for query
- **GIVEN** template ID and name
- **WHEN** `DescribeMonitorNoticeContentTmplById` is called
- **THEN** the method SHALL construct proper API request with pagination
- **AND** return template details or nil if not found
- **AND** handle errors with proper retry logic

---

### Requirement: Error Handling
The system SHALL handle API errors gracefully with retry logic.

#### Scenario: Transient API failure
- **GIVEN** temporary API connectivity issue
- **WHEN** performing CRUD operation
- **THEN** the system SHALL retry with exponential backoff
- **AND** return error only after max retries exceeded

#### Scenario: API validation error
- **GIVEN** invalid template configuration
- **WHEN** API returns validation error
- **THEN** the system SHALL return error immediately without retry
- **AND** include error details in message

---

### Requirement: Logging
The system SHALL log all API operations for debugging.

#### Scenario: Successful API call
- **WHEN** any API operation succeeds
- **THEN** the system SHALL log request body and response body at DEBUG level

#### Scenario: API failure
- **WHEN** any API operation fails
- **THEN** the system SHALL log error details at CRITICAL level

---

### Requirement: Context Management
The system SHALL use context for lifecycle management.

#### Scenario: Context propagation
- **WHEN** any CRUD operation is invoked
- **THEN** the system SHALL create resource lifecycle context
- **AND** propagate context to service layer and API calls
- **AND** include log ID for request tracing

---

## Technical Constraints

### API Constraints
- **Monitor API Version**: 2023-06-16
- **Request Rate Limit**: 20 requests/second
- **TmplIDs Parameter**: Must be string array format
- **Pagination**: Required for DescribeNoticeContentTmpl (PageNumber, PageSize)

### Schema Constraints
- **Composite ID Format**: `tmplID#tmplName` using `tccommon.FILED_SP` separator
- **ForceNew Fields**: `tmpl_name`, `monitor_type`, `tmpl_language`
- **MaxItems**: Use `MaxItems: 1` for single-object nested structures
- **Type Safety**: All API pointer fields must be safely dereferenced

### Implementation Constraints
- **Reference Code**: Must follow patterns from `resource_tc_igtm_strategy.go`
- **Helper Usage**: Use `helper.String()`, `helper.IntInt64()` for conversions
- **Retry Logic**: Use `tccommon.WriteRetryTimeout` and `tccommon.ReadRetryTimeout`
- **Nil Checks**: All pointer access must check for nil first

---

## Non-Functional Requirements

### Performance
- **Read Operations**: Should complete within 5 seconds under normal conditions
- **Bulk Operations**: Support efficient handling of nested structures

### Reliability
- **Idempotency**: All operations must be idempotent
- **Retry Safety**: Retry logic must not cause duplicate resources

### Maintainability
- **Code Organization**: Follow existing service directory structure
- **Documentation**: Include inline comments for complex logic
- **Testing**: Provide acceptance tests for all CRUD operations
