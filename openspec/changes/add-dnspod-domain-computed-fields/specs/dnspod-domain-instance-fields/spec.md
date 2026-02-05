# Spec: DNSPod Domain Instance Resource - Additional Computed Fields

## Overview

This spec defines the schema and behavior changes for the `tencentcloud_dnspod_domain_instance` resource to expose additional computed fields from the DNSPod API's `DomainInfo` response.

## MODIFIED Requirements

### Requirement: Domain Instance Schema Definition

**ID**: `dnspod-domain-instance-schema-v2`

The `tencentcloud_dnspod_domain_instance` resource MUST expose computed fields that reflect the complete state of a DNSPod domain instance returned by the `DescribeDomain` API.

#### Scenario: Status field is read-only

**Given** a user configures a `tencentcloud_dnspod_domain_instance` resource  
**When** the configuration includes a `status` parameter  
**Then** Terraform MUST reject the configuration with an error indicating that `status` is a computed (read-only) field  
**And** the error message SHOULD suggest removing the `status` parameter from the configuration

**Rationale**: The `status` field represents the system-managed state of the domain (enable/pause/spam/lock) rather than a user-configurable property. Making it computed-only prevents confusion between user intent and actual system state.

---

#### Scenario: Record count is exposed

**Given** a domain exists with N DNS records  
**When** Terraform reads the domain state via `DescribeDomain` API  
**Then** the `record_count` computed field MUST be set to N  
**And** the value MUST match the `RecordCount` field from the API response  
**And** the value MUST be of type `int`

**Example**:
```hcl
resource "tencentcloud_dnspod_domain_instance" "example" {
  domain = "example.com"
}

output "total_records" {
  value = tencentcloud_dnspod_domain_instance.example.record_count
}
# Output: total_records = 5
```

---

#### Scenario: Domain grade/plan is exposed

**Given** a domain has a DNSPod plan/package grade  
**When** Terraform reads the domain state  
**Then** the `grade` computed field MUST be set to the domain's package grade  
**And** the value MUST match the `Grade` field from the API response  
**And** the value MUST be of type `string`  
**And** typical values include: "DP_Free", "DP_Plus", "DP_Extra", "DP_Expert", etc.

**Example**:
```hcl
resource "tencentcloud_dnspod_domain_instance" "example" {
  domain = "example.com"
}

output "dns_plan" {
  value = tencentcloud_dnspod_domain_instance.example.grade
}
# Output: dns_plan = "DP_Free"
```

---

#### Scenario: Domain status reflects actual state

**Given** a domain exists with a specific status (enable/pause/spam/lock)  
**When** Terraform reads the domain state  
**Then** the `status` computed field MUST reflect the actual domain status  
**And** the value MUST match the `Status` field from the API response exactly  
**And** NO status transformation or mapping MUST occur (e.g., "pause" should remain "pause", not converted to "disable")  
**And** possible values are: "enable", "pause", "spam", "lock"

**Example**:
```hcl
resource "tencentcloud_dnspod_domain_instance" "example" {
  domain = "example.com"
}

output "domain_state" {
  value = tencentcloud_dnspod_domain_instance.example.status
}
# Output: domain_state = "enable"
```

**Breaking Change Note**: Previously, `status` was an optional configurable field. Users who have `status = "enable"` or similar in their configurations MUST remove this parameter.

---

#### Scenario: Last update time is exposed

**Given** a domain has been modified at a specific timestamp  
**When** Terraform reads the domain state  
**Then** the `updated_on` computed field MUST be set to the last modification timestamp  
**And** the value MUST match the `UpdatedOn` field from the API response  
**And** the value MUST be of type `string`  
**And** the timestamp format MUST match the API's response format (typically ISO 8601 / RFC3339)

**Example**:
```hcl
resource "tencentcloud_dnspod_domain_instance" "example" {
  domain = "example.com"
}

output "last_modified" {
  value = tencentcloud_dnspod_domain_instance.example.updated_on
}
# Output: last_modified = "2024-01-15 10:30:00"
```

---

### Requirement: Domain Instance Create Behavior

**ID**: `dnspod-domain-instance-create-no-status`

The domain instance creation MUST NOT attempt to set the domain status through an explicit API call.

#### Scenario: Create without status configuration

**Given** a user creates a domain instance without specifying a `status` field  
**When** Terraform executes the create operation  
**Then** the domain MUST be created using the `CreateDomain` API  
**And** NO subsequent call to `ModifyDnsPodDomainStatus` MUST be made  
**And** the domain's initial status MUST be determined by the DNSPod service defaults  
**And** the `status` field MUST be populated from the API response during the subsequent read operation

**Code Impact**: Remove lines 115-123 from `resourceTencentCloudDnspodDomainInstanceCreate`:
```go
// REMOVE THIS BLOCK:
if v, ok := d.GetOk("status"); ok {
    domainId := response.Response.DomainInfo.Domain
    status := v.(string)
    err := service.ModifyDnsPodDomainStatus(ctx, *domainId, status)
    if err != nil {
        log.Printf("[CRITAL]%s set DnsPod Domain status failed, reason:%s\n", logId, err.Error())
        return err
    }
}
```

---

### Requirement: Domain Instance Update Behavior

**ID**: `dnspod-domain-instance-update-no-status`

The domain instance update MUST NOT handle changes to the `status` field.

#### Scenario: Update ignores status changes

**Given** a domain instance exists  
**When** Terraform detects any state drift or executes an update plan  
**Then** the update function MUST NOT check for `status` field changes  
**And** NO call to `ModifyDnsPodDomainStatus` MUST be made based on `status` field changes  
**And** the `status` field MUST only be refreshed through the read operation

**Code Impact**: Remove lines 199-206 from `resourceTencentCloudDnspodDomainInstanceUpdate`:
```go
// REMOVE THIS BLOCK:
if d.HasChange("status") {
    status := d.Get("status").(string)
    err := service.ModifyDnsPodDomainStatus(ctx, id, status)
    if err != nil {
        log.Printf("[CRITAL]%s modify DnsPod Domain status failed, reason:%s\n", logId, err.Error())
        return err
    }
}
```

---

### Requirement: Domain Instance Read Behavior

**ID**: `dnspod-domain-instance-read-all-fields`

The domain instance read operation MUST populate all computed fields from the `DomainInfo` API response.

#### Scenario: Read operation maps all new fields

**Given** the `DescribeDomain` API returns a `DomainInfo` object  
**When** Terraform executes the read operation  
**Then** all of the following fields MUST be mapped to the resource state:
  - `status` from `DomainInfo.Status` (without transformation)
  - `record_count` from `DomainInfo.RecordCount` (converted to int)
  - `grade` from `DomainInfo.Grade`
  - `updated_on` from `DomainInfo.UpdatedOn`  
**And** each field MUST have nil-safety checks  
**And** if a field is nil in the API response, it MUST NOT be set in the state (or set to an appropriate zero value)

**Code Example**:
```go
// In resourceTencentCloudDnspodDomainInstanceRead function:
info := response.Response.DomainInfo

// Existing fields
_ = d.Set("domain_id", info.DomainId)
_ = d.Set("domain", info.Domain)
_ = d.Set("create_time", info.CreatedOn)
_ = d.Set("is_mark", info.IsMark)
_ = d.Set("slave_dns", info.SlaveDNS)

// NEW: Direct status mapping (no transformation)
if info.Status != nil {
    _ = d.Set("status", info.Status)
}

// NEW: Record count
if info.RecordCount != nil {
    _ = d.Set("record_count", int(*info.RecordCount))
}

// NEW: Grade
if info.Grade != nil {
    _ = d.Set("grade", info.Grade)
}

// NEW: Updated timestamp
if info.UpdatedOn != nil {
    _ = d.Set("updated_on", info.UpdatedOn)
}

// Existing remark and group_id handling...
```

---

#### Scenario: Null API fields are handled safely

**Given** the `DescribeDomain` API returns a `DomainInfo` with some nil fields  
**When** Terraform executes the read operation  
**Then** the resource MUST NOT panic or error  
**And** nil fields MUST be skipped (not set in state) or set to appropriate zero values  
**And** non-nil fields MUST be set correctly

**Example**: If `RecordCount` is nil, the `record_count` field MUST remain unset or be set to 0 (depending on Terraform SDK behavior).

---

## REMOVED Requirements

### Requirement: Domain Status as Configurable Field

**ID**: `dnspod-domain-instance-status-configurable` (REMOVED)

**Reason for Removal**: The `status` field is now computed-only to accurately reflect system state rather than user configuration. Domain status control should be managed through separate resources or data sources if needed in the future.

**Migration Note**: Users who previously used `status = "enable"` or `status = "disable"` in their configurations MUST:
1. Remove the `status` parameter from resource blocks
2. Use the `status` computed attribute to read the domain's actual state
3. If status control is required, wait for a future resource like `tencentcloud_dnspod_domain_status` (not in this change)

---

## API Mapping

### DescribeDomain Response Mapping

| Terraform Field | API Field (DomainInfo) | Type Conversion      | Notes                                    |
|-----------------|------------------------|----------------------|------------------------------------------|
| `domain_id`     | `DomainId`             | uint64 → int         | Existing                                 |
| `domain`        | `Domain`               | string → string      | Existing                                 |
| `create_time`   | `CreatedOn`            | string → string      | Existing                                 |
| `is_mark`       | `IsMark`               | string → string      | Existing                                 |
| `slave_dns`     | `SlaveDNS`             | string → string      | Existing                                 |
| `status`        | `Status`               | string → string      | **MODIFIED**: Now direct mapping, no transformation |
| `record_count`  | `RecordCount`          | uint64 → int         | **NEW**                                  |
| `grade`         | `Grade`                | string → string      | **NEW**                                  |
| `updated_on`    | `UpdatedOn`            | string → string      | **NEW**                                  |

---

## Schema Definition

```go
Schema: map[string]*schema.Schema{
    // ... existing required fields ...
    
    "domain": {
        Type:        schema.TypeString,
        Required:    true,
        Description: "The Domain.",
    },
    
    // ... existing optional fields ...
    
    // MODIFIED: status is now Computed only
    "status": {
        Type:        schema.TypeString,
        Computed:    true,  // Changed from Optional: true
        Description: "The status of domain. Possible values: `enable`, `pause`, `spam`, `lock`.",
    },
    
    // ... existing computed fields ...
    
    "domain_id": {
        Type:        schema.TypeInt,
        Computed:    true,
        Description: "ID of the domain.",
    },
    
    // NEW computed fields
    "record_count": {
        Type:        schema.TypeInt,
        Computed:    true,
        Description: "Number of DNS records under this domain.",
    },
    
    "grade": {
        Type:        schema.TypeString,
        Computed:    true,
        Description: "The DNS plan/package grade of the domain (e.g., DP_Free, DP_Plus).",
    },
    
    "updated_on": {
        Type:        schema.TypeString,
        Computed:    true,
        Description: "Last modification time of the domain.",
    },
}
```

---

## Testing Requirements

### Unit Test Expectations

#### Test: Computed fields are populated

**Test Case**: `TestAccTencentCloudDnspodDoamin`

The test MUST verify that all new computed fields are set after resource creation:

```go
Check: resource.ComposeTestCheckFunc(
    testAccCheckDnspodDomainExists("tencentcloud_dnspod_domain_instance.domain"),
    resource.TestCheckResourceAttr("tencentcloud_dnspod_domain_instance.domain", "domain_id", "0"),
    resource.TestCheckResourceAttr("tencentcloud_dnspod_domain_instance.domain", "domain", "terraformer.com"),
    resource.TestCheckResourceAttr("tencentcloud_dnspod_domain_instance.domain", "is_mark", "no"),
    resource.TestCheckResourceAttr("tencentcloud_dnspod_domain_instance.domain", "slave_dns", "no"),
    
    // NEW: Verify new computed fields are set
    resource.TestCheckResourceAttrSet("tencentcloud_dnspod_domain_instance.domain", "record_count"),
    resource.TestCheckResourceAttrSet("tencentcloud_dnspod_domain_instance.domain", "grade"),
    resource.TestCheckResourceAttrSet("tencentcloud_dnspod_domain_instance.domain", "status"),
    resource.TestCheckResourceAttrSet("tencentcloud_dnspod_domain_instance.domain", "updated_on"),
),
```

#### Test: Status configuration is rejected

**Test Case**: `TestAccTencentCloudDnspodDoamin_StatusConfigError` (NEW)

```go
{
    Config:      testAccTencentCloudDnspodDomainWithStatus,
    ExpectError: regexp.MustCompile("status.*computed.*read-only"),
}
```

Test configuration:
```hcl
resource "tencentcloud_dnspod_domain_instance" "domain" {
  domain  = "terraformer.com"
  status  = "enable"  # This should cause an error
}
```

---

## Documentation Requirements

### Resource Documentation Updates

**File**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.md`

#### Argument Reference Section

REMOVE `status` from optional arguments (if present).

#### Attributes Reference Section

ADD the following:

```markdown
In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `domain_id` - ID of the domain.
* `create_time` - Create time of the domain.
* `slave_dns` - Is secondary DNS enabled.
* `record_count` - Number of DNS records under this domain.
* `grade` - The DNS plan/package grade of the domain (e.g., DP_Free, DP_Plus).
* `status` - The status of domain. Possible values: `enable`, `pause`, `spam`, `lock`.
* `updated_on` - Last modification time of the domain.
```

#### Add Migration Notice

```markdown
~> **NOTE**: Starting from version X.X.X, the `status` field is computed-only. 
   If you have `status` configured in your resource blocks, please remove it. 
   The actual domain status can be read from the `status` attribute.
```

---

## Backward Compatibility

### Breaking Changes

**Impact**: Users with `status = "enable"` or `status = "disable"` in their configurations will experience a validation error.

**Severity**: Medium - affects only users explicitly setting the `status` field.

**Migration Path**:
1. Identify resources using `status` parameter: `grep -r 'status.*=' *.tf`
2. Remove `status` from resource configuration
3. Access status through the computed attribute: `resource.status`

**Example Migration**:

Before:
```hcl
resource "tencentcloud_dnspod_domain_instance" "example" {
  domain = "example.com"
  status = "enable"  # Remove this line
}
```

After:
```hcl
resource "tencentcloud_dnspod_domain_instance" "example" {
  domain = "example.com"
  # status is now read-only, access it as an attribute
}

output "domain_status" {
  value = tencentcloud_dnspod_domain_instance.example.status
}
```

---

## Implementation Checklist

- [ ] Schema: Modify `status` field to be Computed-only
- [ ] Schema: Add `record_count` field (TypeInt, Computed)
- [ ] Schema: Add `grade` field (TypeString, Computed)
- [ ] Schema: Add `updated_on` field (TypeString, Computed)
- [ ] Read: Remove status transformation logic
- [ ] Read: Add `status` direct mapping with nil check
- [ ] Read: Add `record_count` mapping with type conversion and nil check
- [ ] Read: Add `grade` mapping with nil check
- [ ] Read: Add `updated_on` mapping with nil check
- [ ] Create: Remove status setting logic (lines 115-123)
- [ ] Update: Remove status update logic (lines 199-206)
- [ ] Test: Remove status from test configurations
- [ ] Test: Add checks for new computed fields
- [ ] Test: Add negative test for status configuration
- [ ] Docs: Update resource documentation
- [ ] Docs: Generate website docs
- [ ] Docs: Add CHANGELOG entry with BREAKING CHANGES tag
- [ ] Code: Run gofmt
- [ ] Code: Run golangci-lint
- [ ] Code: Compile successfully
- [ ] Manual: Test domain creation with new fields
- [ ] Manual: Verify status configuration is rejected

---

## Related Resources

**None** - This change is isolated to the `tencentcloud_dnspod_domain_instance` resource.

Future enhancements could include:
- A separate `tencentcloud_dnspod_domain_status` resource for explicit status control
- Additional computed fields from `DomainInfo` (e.g., `owner`, `grade_level`, `vip_end_at`)

---

## References

- **DNSPod API**: `DescribeDomain` / `CreateDomain` interfaces
- **SDK Source**: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323/models.go:5171-5270`
- **Existing Code**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`
