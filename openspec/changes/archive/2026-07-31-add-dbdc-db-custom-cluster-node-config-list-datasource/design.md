## Context

The `dbdc` (DB Dedicated Cloud) service in the TencentCloud Terraform provider
already ships several read-only data sources for DB Custom clusters, nodes,
images, and cluster-node membership. Operators can discover *which* nodes exist,
but cannot currently discover the Kubernetes scheduling configuration — labels
and taints — attached to each node. The cloud API
`DescribeDBCustomClusterNodeConfig` (`dbdc` SDK `v20201029`) returns exactly
this: given a `ClusterId` (and optionally a list of `NodeIds`), it returns a
`NodeSet` of `DBCustomClusterNodeConfig` items, each containing `NodeId`,
`Labels` (`[]*Label` with `Key`/`Value`) and `Taints` (`[]*Taint` with
`Key`/`Value`/`Effect`).

This design adds a new `RESOURCE_KIND_DATASOURCE`
`tencentcloud_dbdc_db_custom_cluster_node_config_list` that wraps that API. It
follows the established pattern of the sibling
`tencentcloud_dbdc_db_custom_cluster_nodes` data source (same package
`tencentcloud/services/dbdc`, same SDK alias `dbdcv20201029`, same
`UseDbdcV20201029Client()` accessor, same `result_output_file` convention).

## Goals / Non-Goals

**Goals:**
- Provide a read-only Terraform data source that returns the labels and taints
  configuration for nodes in a DB Custom cluster.
- Flatten the `NodeSet` list so each element exposes `node_id`, `labels`, and
  `taints` directly (no extra `xxx_set` nesting layer), per the provider
  convention for Describe-style data sources.
- Support filtering by an optional `node_ids` list.
- Reuse the shared retry/ratelimit conventions (`tccommon.ReadRetryTimeout`,
  `tccommon.RetryError`, `ratelimit.Check`).
- Guard all pointer/slice accesses against nil, and follow the strict empty
  response handling rule for `RESOURCE_KIND_DATASOURCE` (return
  `NonRetryableError` instead of clearing the local state id).

**Non-Goals:**
- No Create/Update/Delete — this is a data source (read only).
- No pagination loop: the `DescribeDBCustomClusterNodeConfig` API does not
  expose `Offset`/`Limit`; it returns the full `NodeSet` for the requested
  `ClusterId`/`NodeIds` in a single call. (The `node_ids` input itself is
  capped at 100 by the cloud API.)
- No website markdown is hand-written; the `data_source_*.md` example file is
  generated into `website/docs/d/` via `make doc`.

## Decisions

### Decision 1: Flattened `node_set` schema (no extra nesting)

**Choice**: Expose `node_set` as a `schema.TypeList` whose `Elem` is a
`schema.Resource` with top-level `node_id`, `labels`, `taints` — rather than a
single nested block that wraps everything.

**Rationale**: Per the project's data-source convention, the list returned by a
Describe endpoint must be **flattened**: each list element's fields are surfaced
directly, never hidden behind an extra `xxx_set`/`xxx_list` wrapper layer. This
keeps each field individually set/readable by Terraform and matches the sibling
`dbdc_db_custom_cluster_nodes` data source's `node_set` shape.

**Alternative considered**: A single nested block mirroring the API response
verbatim. Rejected because it violates the flatten convention and complicates
downstream `for_each` / `splat` expressions.

### Decision 2: `node_ids` is `TypeList` of `TypeString` (not `TypeSet`)

**Choice**: `node_ids` argument uses
`TypeList` / `Elem: TypeString`, Optional.

**Rationale**: The cloud API `NodeIds` is an ordered `[]*string` (max 100). Order
is not semantically significant for the result, but `TypeList` keeps the
mapping to `[]*string` straightforward and consistent with how other dbdc data
sources map list inputs. Conversion to `[]*string` via `helper.String` is done
in the Read function before calling the service layer.

### Decision 3: Empty-response handling returns `NonRetryableError`

**Choice**: Inside the retry block, if
`result == nil || result.Response == nil ||
len(result.Response.NodeSet) == 0`, return
`resource.NonRetryableError(...)` instead of clearing `d.SetId("")`.

**Rationale**: This is a hard rule for `RESOURCE_KIND_DATASOURCE` resources in
this provider. Clearing the id on a transient cloud API hiccup would lose local
state; surfacing `NonRetryableError` lets the outer retry keep trying and, on
exhaustion, fail loudly for human intervention. On the failure path we still
emit `log.Printf("[DATASOURCE] read empty, skip SetId")` for diagnostics.

### Decision 4: Service helper `DescribeDBCustomClusterNodeConfigByFilter`

**Choice**: Add a single service-layer helper
`DescribeDBCustomClusterNodeConfigByFilter(ctx, param map[string]interface{})
(ret []*dbdcv20201029.DBCustomClusterNodeConfig, errRet error)` to
`service_tencentcloud_dbdc.go`, mirroring the existing
`DescribeDBCustomClusterNodesByFilter` helper shape (same `param` map style,
same `defer` log on error, same `resource.Retry` +
`ratelimit.Check` wrapping).

**Rationale**: Keeps the cloud-API call, retry, and error logging centralized
in the service layer so the data-source Read function stays a thin
schema-mapping layer — exactly like the sibling cluster-nodes data source.

### Decision 5: Resource ID is a generated token

**Choice**: `d.SetId(helper.BuildToken())` after a successful read, identical to
the sibling data sources.

**Rationale**: Data sources have no persistent cloud identity; a generated token
satisfies the Terraform requirement that `Read` set an id and keeps behavior
consistent with the rest of the `dbdc` data sources.

## Risks / Trade-offs

- [API returns null for `Labels`/`Taints`] → The API documents that `Labels`
  and `Taints` "may return null". The Read function nil-checks both before
  flattening, so a null field yields an empty list rather than a panic.
- [No pagination support in the API] → The API has no `Offset`/`Limit`; it
  returns all matching nodes in one call. If a cluster ever exceeds the
  internal response cap, the data source will simply return whatever the API
  returns. This matches the cloud API contract and is not a regression.
- [`node_ids` cap of 100] → Enforced by the cloud API. We do not add our own
  `MaxItems` validation; passing >100 will surface as an API error to the user,
  which is acceptable and avoids provider-side magic numbers drifting from the
  API limit.
- [Reads can be flaky] → Mitigated by `resource.Retry(ReadRetryTimeout)` with
  `tccommon.RetryError`; on a truly empty result we fail with
  `NonRetryableError` rather than silently clearing state.
