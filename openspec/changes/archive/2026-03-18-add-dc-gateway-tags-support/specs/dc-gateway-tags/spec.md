# Specification: DC Gateway Tags Support

**Capability ID**: `dc-gateway-tags`  
**Resource**: `tencentcloud_dc_gateway`  
**Status**: Draft

---

## ADDED Requirements

### Requirement: DCG-TAGS-001 - Schema Support for Tags
**Priority**: High  
**Type**: Functional

The `tencentcloud_dc_gateway` resource MUST support a `tags` parameter in its schema.

**Acceptance Criteria**:
- `tags` field is of type `TypeMap`
- `tags` field is optional
- `tags` field accepts string key-value pairs
- Description clearly explains the field's purpose

#### Scenario: Define tags schema field
**Given** the DC gateway resource schema  
**When** a user defines the resource in Terraform configuration  
**Then** they can specify a `tags` map parameter  
**And** the parameter accepts multiple key-value pairs  
**And** the parameter is optional (resource works without tags)

```hcl
resource "tencentcloud_dc_gateway" "example" {
  name                = "test-dcg"
  network_type        = "VPC"
  network_instance_id = "vpc-123456"
  gateway_type        = "NORMAL"
  
  tags = {
    Environment = "production"
    Team        = "networking"
    CostCenter  = "IT-001"
  }
}
```

---

### Requirement: DCG-TAGS-002 - Create Gateway with Tags
**Priority**: High  
**Type**: Functional

When creating a DC gateway, any tags specified MUST be set on the resource using the `CreateDirectConnectGateway` API's `Tags` parameter.

**Acceptance Criteria**:
- Tags are extracted from Terraform schema
- Tags are converted to `[]*vpc.Tag` format
- Tags are included in `CreateDirectConnectGatewayRequest`
- Gateway is created with all specified tags
- Tags are immediately visible after creation

#### Scenario: Create gateway with tags
**Given** a Terraform configuration with tags specified  
**When** `terraform apply` is executed  
**Then** the CreateDirectConnectGateway API is called with Tags parameter  
**And** the gateway is created successfully  
**And** all specified tags are applied to the gateway  
**And** tags are retrievable via tag service API

```hcl
resource "tencentcloud_dc_gateway" "tagged" {
  name                = "tagged-gateway"
  network_type        = "VPC"
  network_instance_id = "vpc-abc123"
  
  tags = {
    Owner = "ops-team"
    Stage = "testing"
  }
}
```

#### Scenario: Create gateway without tags
**Given** a Terraform configuration without tags specified  
**When** `terraform apply` is executed  
**Then** the gateway is created successfully  
**And** no tags are applied  
**And** the resource operates normally

---

### Requirement: DCG-TAGS-003 - Read Gateway Tags
**Priority**: High  
**Type**: Functional

When reading a DC gateway's state, tags MUST be retrieved using the universal tag service and stored in Terraform state.

**Acceptance Criteria**:
- Tag service is initialized with correct client
- `DescribeResourceTags()` is called with correct parameters
  - ServiceType: `"vpc"`
  - ResourceType: `"dcg"`
  - Region: Current region
  - ResourceId: Gateway ID
- Retrieved tags are set in Terraform state
- Tag retrieval errors are properly handled

#### Scenario: Read tags from existing gateway
**Given** a DC gateway exists with tags  
**When** Terraform reads the resource state  
**Then** the tag service DescribeResourceTags API is called  
**And** all tags are retrieved successfully  
**And** tags are stored in Terraform state  
**And** `terraform show` displays all tags correctly

#### Scenario: Handle gateway with no tags
**Given** a DC gateway exists without any tags  
**When** Terraform reads the resource state  
**Then** the tag service returns an empty map  
**And** the `tags` field in state is set to empty map  
**And** no errors are raised

#### Scenario: Handle tag service API errors
**Given** the tag service API is unavailable  
**When** Terraform attempts to read tags  
**Then** the error is properly logged  
**And** the error is returned to the user  
**And** the Terraform operation fails gracefully

---

### Requirement: DCG-TAGS-004 - Update Gateway Tags
**Priority**: High  
**Type**: Functional

Tags MUST be updatable without recreating the gateway, using the universal tag service's `ModifyTags` API.

**Acceptance Criteria**:
- Tag changes are detected using `d.HasChange("tags")`
- Old and new tag values are compared
- Tag diff is calculated using `svctag.DiffTags()`
- `ModifyTags()` is called with correct resource name format
- Tags are updated without gateway recreation
- Update handles additions, modifications, and deletions

#### Scenario: Add new tags to gateway
**Given** a gateway exists with tags `{Env: "dev"}`  
**When** user adds tag `{Team: "backend"}`  
**And** `terraform apply` is executed  
**Then** ModifyTags is called with replaceTags containing both tags  
**And** the gateway has both tags after update  
**And** the gateway is not recreated

```hcl
# Before
tags = {
  Env = "dev"
}

# After
tags = {
  Env  = "dev"
  Team = "backend"
}
```

#### Scenario: Modify existing tag values
**Given** a gateway with tags `{Env: "dev", Team: "backend"}`  
**When** user changes `Env` to `"prod"`  
**And** `terraform apply` is executed  
**Then** ModifyTags is called with updated Env value  
**And** the tag value is changed  
**And** other tags remain unchanged  
**And** the gateway is not recreated

```hcl
# Before
tags = {
  Env  = "dev"
  Team = "backend"
}

# After
tags = {
  Env  = "prod"
  Team = "backend"
}
```

#### Scenario: Remove tags from gateway
**Given** a gateway with tags `{Env: "dev", Team: "backend"}`  
**When** user removes the `Team` tag  
**And** `terraform apply` is executed  
**Then** ModifyTags is called with deleteTags containing "Team"  
**And** the Team tag is removed from gateway  
**And** Env tag remains  
**And** the gateway is not recreated

```hcl
# Before
tags = {
  Env  = "dev"
  Team = "backend"
}

# After
tags = {
  Env = "dev"
}
```

#### Scenario: Replace all tags
**Given** a gateway with tags `{Old: "value"}`  
**When** user replaces all tags with `{New: "value"}`  
**And** `terraform apply` is executed  
**Then** ModifyTags is called with appropriate replaceTags and deleteTags  
**And** old tags are removed  
**And** new tags are applied  
**And** the gateway is not recreated

---

### Requirement: DCG-TAGS-005 - Import Gateway with Tags
**Priority**: Medium  
**Type**: Functional

When importing an existing DC gateway, tags MUST be imported and stored in Terraform state.

**Acceptance Criteria**:
- Import command works: `terraform import tencentcloud_dc_gateway.example dcg-xxxxx`
- After import, tags are read via tag service
- All tags are present in Terraform state
- Subsequent `terraform plan` shows no changes for tags

#### Scenario: Import gateway with tags
**Given** a DC gateway exists in Tencent Cloud with tags  
**When** user runs `terraform import tencentcloud_dc_gateway.test dcg-abc123`  
**Then** the gateway is imported successfully  
**And** all tags are read and stored in state  
**And** `terraform show` displays all tags  
**And** `terraform plan` shows no tag changes

```bash
$ terraform import tencentcloud_dc_gateway.test dcg-abc123
$ terraform show
resource "tencentcloud_dc_gateway" "test" {
  id                  = "dcg-abc123"
  name                = "imported-gateway"
  # ... other fields ...
  tags = {
    Environment = "production"
    Managed     = "terraform"
  }
}
```

---

### Requirement: DCG-TAGS-006 - Tag Validation and Constraints
**Priority**: Medium  
**Type**: Functional

Tag operations MUST respect Tencent Cloud tag service constraints and provide clear error messages.

**Acceptance Criteria**:
- Tag keys and values are passed as strings
- Empty tag maps are handled gracefully
- API errors are caught and reported clearly
- Retry logic is applied for transient failures

#### Scenario: Handle tag service rate limits
**Given** tag service API rate limit is exceeded  
**When** Terraform attempts tag operations  
**Then** the operation is retried with exponential backoff  
**And** eventually succeeds or fails with clear error message

#### Scenario: Handle invalid tag format
**Given** user provides invalid tag format (if possible via schema)  
**When** Terraform validates the configuration  
**Then** validation fails with clear error message  
**And** user is informed of correct format

#### Scenario: Empty tags map
**Given** user specifies `tags = {}`  
**When** Terraform applies the configuration  
**Then** no tags are set on the resource  
**And** no errors occur  
**And** resource operates normally

---

### Requirement: DCG-TAGS-007 - Resource Name Format
**Priority**: High  
**Type**: Technical

Tag resource names MUST follow the Tencent Cloud QCS format for proper tag service integration.

**Acceptance Criteria**:
- Resource name uses `BuildTagResourceName()` helper
- Format: `qcs::vpc:{region}:account:/dcg/{gateway-id}`
- ServiceType is `"vpc"`
- ResourceType is `"dcg"`
- Region is current region from client
- Resource ID is gateway ID from `d.Id()`

#### Scenario: Build correct resource name
**Given** a DC gateway with ID `dcg-12345678`  
**And** region is `ap-guangzhou`  
**When** building tag resource name  
**Then** the format MUST be `qcs::vpc:ap-guangzhou:account:/dcg/dcg-12345678`  
**And** the tag service accepts this format

---

### Requirement: DCG-TAGS-008 - Backward Compatibility
**Priority**: High  
**Type**: Non-Functional

Adding tags support MUST NOT break existing DC gateway resources or configurations.

**Acceptance Criteria**:
- Existing gateways without tags continue to work
- Existing Terraform configurations without tags remain valid
- No forced replacement of existing resources
- Tags field is purely additive

#### Scenario: Existing gateway without tags
**Given** a DC gateway created before tags support  
**When** Terraform refreshes state  
**Then** the gateway remains functional  
**And** tags field shows empty map  
**And** no forced replacement occurs  
**And** user can add tags without recreation

#### Scenario: Existing configuration without tags
**Given** a Terraform configuration without tags field  
**When** applying with new provider version  
**Then** the configuration remains valid  
**And** the gateway operates normally  
**And** no warnings or errors occur

---

### Requirement: DCG-TAGS-009 - Error Handling
**Priority**: High  
**Type**: Non-Functional

Tag operations MUST handle errors gracefully with clear, actionable error messages.

**Acceptance Criteria**:
- Tag service errors are logged with context
- Error messages include operation type (create/read/update)
- API errors are wrapped with additional context
- Transient errors trigger retry logic
- Fatal errors prevent resource creation/update

#### Scenario: Handle tag creation failure
**Given** CreateDirectConnectGateway succeeds  
**But** subsequent tag setting fails  
**When** applying configuration  
**Then** an error is logged with full context  
**And** the error message indicates tag operation failed  
**And** user can take corrective action

#### Scenario: Handle tag read failure
**Given** DescribeResourceTags API fails  
**When** reading gateway state  
**Then** the error is logged with gateway ID and operation  
**And** the read operation fails  
**And** user sees clear error message

#### Scenario: Handle tag update failure
**Given** ModifyTags API fails  
**When** updating tags  
**Then** the update operation fails  
**And** gateway state remains unchanged  
**And** user sees error message with details

---

## Success Metrics

- **Functional**: All 9 requirements implemented and tested
- **Testing**: 100% scenario coverage in acceptance tests
- **Compatibility**: No breaking changes to existing resources
- **Documentation**: Complete examples and argument reference
- **Code Quality**: Passes all linters and format checks

---

## Related Changes

- None (standalone capability)

---

## Notes

- This implementation follows the established pattern from `resource_tc_ccn.go` and other VPC resources
- The universal tag service is production-ready and widely used across the provider
- No changes to native DC gateway APIs; purely leveraging existing tag service infrastructure
