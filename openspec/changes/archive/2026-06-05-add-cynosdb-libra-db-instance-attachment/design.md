## Context

TencentCloud CynosDB (TDSQL-C) provides LibraDB read-only analytics engine instances that can be attached to existing CynosDB clusters. The cloud API SDK (`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107`) already includes the necessary APIs:

- `AddLibraDBInstances` - Creates (attaches) a LibraDB analytics engine instance to a cluster
- `DescribeLibraDBInstanceDetail` - Queries the detail of a LibraDB instance
- `IsolateLibraDBCluster` - Isolates (detaches) the LibraDB analytics cluster

The resource follows the RESOURCE_KIND_ATTACHMENT pattern: Create (bind), Read, Delete (unbind). No Update operation is needed since this is an attachment resource with only CRD interfaces.

## Goals / Non-Goals

**Goals:**
- Implement a new Terraform resource `tencentcloud_cynosdb_libra_db_instance` that manages the lifecycle of LibraDB analytics engine instance attachments
- Support Create (AddLibraDBInstances), Read (DescribeLibraDBInstanceDetail), and Delete (IsolateLibraDBCluster) operations
- Use composite ID (`cluster_id` + `instance_id`) with `tccommon.FILED_SP` separator for resource identification
- Follow existing patterns in the cynosdb service directory
- Register the resource in provider.go and provider.md
- Provide documentation with example usage

**Non-Goals:**
- Update operation (not supported by the API for attachment resources)
- Managing the underlying CynosDB cluster itself
- Supporting import of existing LibraDB instances (attachment resources with composite IDs support import via composite ID format)

## Decisions

### 1. Resource ID Strategy
**Decision**: Use composite ID with `cluster_id` + FILED_SP + `instance_id`.

**Rationale**: The Read API (`DescribeLibraDBInstanceDetail`) requires both `ClusterId` and `InstanceId`. The `instance_id` is obtained from the Create response's `ResourceIds` field. Using composite ID allows the Read and Delete operations to extract both identifiers from `d.Id()`.

**Alternatives considered**: Using only `instance_id` as the resource ID would require storing `cluster_id` separately and reading it from state, which is less idiomatic for this provider.

### 2. Create Response Handling
**Decision**: After calling `AddLibraDBInstances`, extract the first element from `ResourceIds` as the `instance_id`, then poll with `DescribeLibraDBInstanceDetail` until the instance status indicates it's ready.

**Rationale**: The `AddLibraDBInstances` API returns `ResourceIds` which contains the created instance IDs. Since this is an async operation (creates infrastructure), we need to poll until the instance is available. The response also returns `TranId`, `BigDealIds`, and `DealNames` which are stored as computed attributes.

### 3. Delete Strategy
**Decision**: Use `IsolateLibraDBCluster` API for deletion, passing `cluster_id` from the composite ID. The `isolate_reason_types` and `isolate_reason` parameters are optional inputs for the delete operation.

**Rationale**: The `IsolateLibraDBCluster` API isolates the entire analytics cluster associated with the CynosDB cluster. This is the designated "unbind" operation for RESOURCE_KIND_ATTACHMENT.

### 4. Schema Design - ForceNew vs Immutable
**Decision**: Mark `cluster_id` as ForceNew (since it's part of the resource ID). All other create-only parameters will be checked in an immutableArgs list in the update function, returning an error if changed.

**Rationale**: Per the requirements for resources with only CRD interfaces, only the ID field should be ForceNew. Other fields should use the immutableArgs pattern in the update method.

### 5. Objects Parameter Handling
**Decision**: Model the `objects` parameter as a nested block with `database_tables` sub-block containing `migrate_db_mode` and `databases` list.

**Rationale**: The `Objects` struct in the SDK contains a `DatabaseTables` field of type `MigrateObject`, which has `MigrateDBMode` and `Databases` (list of `MigrateDBItem` with `DbName`, `MigrateTableMode`, and `Tables`). This nested structure maps naturally to Terraform schema blocks.

## Risks / Trade-offs

- **[Risk] Async Create Operation** → After calling `AddLibraDBInstances`, poll `DescribeLibraDBInstanceDetail` until the instance status is ready. Use `resource.Retry` with `tccommon.ReadRetryTimeout` for polling.
- **[Risk] IsolateLibraDBCluster isolates the entire analytics cluster** → This is the intended behavior for the attachment resource. Users should be aware that deleting this resource isolates the analytics cluster, not just a single instance.
- **[Risk] ResourceIds may be empty in Create response** → Check for nil/empty `ResourceIds` and return `NonRetryableError` if the response doesn't contain the expected instance ID.
- **[Trade-off] No Update support** → Since this is a CRD-only attachment resource, any parameter change (except ForceNew fields) will return an error via the immutableArgs pattern rather than silently recreating the resource.
