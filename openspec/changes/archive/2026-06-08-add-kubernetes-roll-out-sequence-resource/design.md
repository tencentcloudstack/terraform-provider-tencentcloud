## Context

TKE provides cloud APIs for managing cluster roll-out sequences (`CreateRollOutSequence`, `DescribeRollOutSequences`, `ModifyRollOutSequence`, `DeleteRollOutSequence`). These APIs are already vendored in `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525`. The TKE service layer already exists at `tencentcloud/services/tke/`.

A roll-out sequence consists of:
- `Name` (string): The sequence name
- `SequenceFlows` (list of SequenceFlow): Each flow has `Tags` (list of SequenceTag with Key/Value) and `SoakTime` (int64, seconds)
- `Enabled` (bool): Whether the sequence is enabled
- `ID` (int64): The sequence identifier (returned on create, used for read/update/delete)

The `DescribeRollOutSequences` API returns a paginated list of all sequences. To read a specific sequence, we must iterate through the list and filter by ID.

## Goals / Non-Goals

**Goals:**
- Implement a fully functional RESOURCE_KIND_GENERAL Terraform resource `tencentcloud_kubernetes_roll_out_sequence` with CRUD operations
- Follow existing provider patterns (retry logic, error handling, resource registration)
- Provide unit tests using gomonkey mocks for cloud API calls
- Provide documentation with example usage

**Non-Goals:**
- Data source for listing roll-out sequences (separate resource if needed later)
- Managing `DescribeClusterRollOutSequenceTags` or `ModifyClusterRollOutSequenceTags` APIs (different resource concern)
- Import support for complex nested structures beyond basic ID import

## Decisions

### 1. Resource ID uses int64 converted to string
The cloud API uses `int64` for the sequence ID. Terraform resource IDs must be strings. We will convert using `strconv.FormatInt` / `strconv.ParseInt`.

**Rationale**: Standard pattern in the provider for numeric IDs. Simple and reversible.

### 2. Read implementation filters from list API
The `DescribeRollOutSequences` API only provides a paginated list endpoint with no single-resource query. The Read function will paginate through all results and filter by ID.

**Rationale**: This is the only available API. We use pagination with Limit=20 (the API default max) and iterate until we find the matching ID or exhaust results.

### 3. Schema structure for sequence_flows
```
sequence_flows (Required, List):
  tags (Required, List):
    key (Required, String)
    value (Required, List of String)
  soak_time (Required, Int)
```

**Rationale**: Direct mapping from the cloud API `SequenceFlow` and `SequenceTag` structures. All fields are required since a sequence flow without tags or soak time is meaningless.

### 4. All fields are updatable via ModifyRollOutSequence
The `ModifyRollOutSequence` API accepts the same fields as Create (Name, SequenceFlows, Enabled) plus the ID. No fields need to be ForceNew.

**Rationale**: The API supports in-place updates for all fields.

### 5. Unit tests use gomonkey mocks
Following the requirement for new terraform resources, unit tests will mock the TKE client methods using gomonkey rather than using the Terraform acceptance test framework.

**Rationale**: Avoids dependency on real cloud credentials and infrastructure for testing business logic.

## Risks / Trade-offs

- [Risk] The Describe API has no filter-by-ID parameter, requiring full pagination scan → Mitigation: For most accounts, the number of roll-out sequences is small. Pagination with Limit=20 is efficient enough.
- [Risk] The ID is int64 but Terraform IDs are strings → Mitigation: Use standard strconv conversion, well-established pattern in the provider.
- [Risk] No documented maximum for Limit in DescribeRollOutSequences → Mitigation: Use 20 as stated in the API comment ("默认为20").
