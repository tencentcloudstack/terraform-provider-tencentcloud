# Design: Public Network Access Configuration

## Context
Tencent Cloud RabbitMQ VIP instances support public network access configuration at creation time. The Terraform provider needs to expose these capabilities while maintaining consistency with the platform's immutability constraints.

## Goals
- Enable users to configure public network access during instance creation
- Maintain field immutability after creation (aligned with API behavior)
- Provide clear error messages when users attempt to modify immutable fields
- Ensure accurate state representation through proper data type conversions

## Non-Goals
- Enabling modification of public access settings after instance creation (API limitation)
- Adding bandwidth throttling or traffic management features (out of scope)
- Implementing client-side validation for bandwidth limits (delegated to API)

## Decisions

### 1. Field Naming Convention
**Decision**: Use `band_width` (snake_case) and `enable_public_access` to match existing Terraform resource conventions in this provider.

**Rationale**: 
- Consistent with provider's naming pattern (e.g., `auto_renew_flag`, `storage_size`)
- Clear and descriptive for Terraform users
- Maps cleanly to API fields `Bandwidth` and `EnablePublicAccess`

**Alternatives considered**:
- `public_bandwidth` - Less explicit about what "public" refers to
- `public_network_bandwidth` - Too verbose

### 2. Immutability Implementation
**Decision**: Add both fields to the `immutableArgs` list in the Update function, preventing changes after creation.

**Rationale**:
- Aligns with Tencent Cloud API behavior (no modify operation for these fields)
- Prevents users from expecting modification support
- Clear error messages guide users to recreate resources if needed
- Consistent with other immutable fields in the resource (e.g., `zone_ids`, `vpc_id`)

**Implementation**:
```go
immutableArgs := []string{
    "zone_ids", "vpc_id", "subnet_id", "node_spec", "node_num",
    "storage_size", "enable_create_default_ha_mirror_queue",
    "auto_renew_flag", "time_span", "pay_mode", "cluster_version",
    "band_width", "enable_public_access",  // New additions
}
```

### 3. String-to-Bool Conversion for Read Operation
**Decision**: Convert API response string values ("ON"/"OFF") to boolean for `enable_public_access` state.

**Rationale**:
- Terraform users expect boolean types for enable/disable flags
- Matches schema definition (`schema.TypeBool`)
- Simplifies configuration (users write `true`/`false`, not `"ON"`/`"OFF"`)

**Implementation**:
```go
if rabbitmqVipInstance.ClusterNetInfo != nil && rabbitmqVipInstance.ClusterNetInfo.PublicDataStreamStatus != nil {
    enablePublicAccess := *rabbitmqVipInstance.ClusterNetInfo.PublicDataStreamStatus == "ON"
    _ = d.Set("enable_public_access", enablePublicAccess)
}
```

**Alternatives considered**:
- Keep as string type - Less idiomatic for Terraform, requires users to remember string values
- Use custom validation - Unnecessary complexity, API handles validation

### 4. Optional Fields with No Default Values in Schema
**Decision**: Mark both fields as `Optional: true` without `Default` in schema definition.

**Rationale**:
- Allows API to apply its own defaults when fields are not specified
- Prevents Terraform from forcing values when users don't set them
- Maintains backward compatibility (existing configs without these fields work unchanged)

**API behavior when fields are omitted**:
- `EnablePublicAccess`: Defaults to `false` (no public access)
- `Bandwidth`: Ignored when public access is disabled

### 5. Bandwidth Field Dependency on Public Access
**Decision**: Do NOT enforce client-side validation that `band_width` requires `enable_public_access = true`.

**Rationale**:
- API accepts bandwidth values regardless of public access flag
- Client-side validation adds complexity and maintenance burden
- Users can learn from API errors if they misconfigure
- Bandwidth field is harmless when public access is disabled (simply ignored)

**Documentation approach**: Clearly state in docs that `band_width` "only takes effect when enable_public_access is true".

## Data Flow

### Create Flow
```
Terraform Config
  ├─ enable_public_access: true (bool)
  └─ band_width: 100 (int)
        ↓
Provider Schema Parsing
  ├─ d.GetOkExists("enable_public_access") → bool
  └─ d.GetOkExists("band_width") → int
        ↓
Request Construction
  ├─ request.EnablePublicAccess = helper.Bool(true)
  └─ request.Bandwidth = helper.IntInt64(100)
        ↓
API Call: CreateRabbitMQVipInstance
        ↓
Instance Created with Public Access
```

### Read Flow
```
API Call: DescribeRabbitMQVipInstances
        ↓
Response Parsing
  ├─ ClusterNetInfo.PublicDataStreamStatus: "ON" (string)
  └─ ClusterSpecInfo.PublicNetworkTps: 100 (*int64)
        ↓
Type Conversion
  ├─ "ON" → true (bool)
  └─ *int64 → int
        ↓
State Update
  ├─ d.Set("enable_public_access", true)
  └─ d.Set("band_width", 100)
        ↓
Terraform State
  ├─ enable_public_access: true
  └─ band_width: 100
```

## Error Handling

### Nil Value Handling in Read
- **ClusterNetInfo nil**: Skip setting `enable_public_access` (field absent in state)
- **PublicDataStreamStatus nil**: Default to `false` or skip setting
- **ClusterSpecInfo nil**: Skip setting `band_width`
- **PublicNetworkTps nil**: Skip setting `band_width`

**Implementation Pattern**:
```go
if rabbitmqVipInstance.ClusterNetInfo != nil && 
   rabbitmqVipInstance.ClusterNetInfo.PublicDataStreamStatus != nil {
    // Safe to access and convert
}
```

### API Error Propagation
- Creation errors with public access configuration → Propagate to user with full error context
- Quota/permission errors → Clear error messages from API
- No retry logic for quota errors (non-retryable)

## Risks and Mitigations

### Risk: Users attempt to modify immutable fields
**Mitigation**: 
- Clear error messages: "argument `enable_public_access` cannot be changed"
- Documentation explicitly states immutability
- Immutability enforced at Update function level (fails fast)

### Risk: String-to-bool conversion assumes only "ON"/"OFF" values
**Mitigation**:
- API contract guarantees these values (based on API documentation)
- If new values appear, default to `false` (safe fallback)
- Read operation won't crash (defensive nil checks)

### Risk: Bandwidth without public access causes confusion
**Mitigation**:
- Documentation clearly explains field relationship
- API ignores bandwidth when public access is disabled (no harm)
- No client-side validation to avoid false positives

## Migration Plan

### For Existing Resources
1. Users upgrade provider version with new fields
2. Existing resources without public access fields:
   - State refresh reads current values from API
   - No plan changes if fields are not added to config
   - Backward compatible (no forced recreation)

### For New Configurations
1. Users can immediately use new fields in new resources
2. Optional fields don't break existing configurations
3. API handles default values when fields are omitted

### Adding Public Access to Existing Resources
**User workflow**:
1. User adds `enable_public_access = true` to existing resource config
2. `terraform plan` shows immutability error
3. User must either:
   - Remove the change (keep current state)
   - Manually destroy and recreate resource with new settings
4. Documentation provides clear guidance on this limitation

## Open Questions
None - All design decisions are finalized based on API documentation and provider conventions.
