## Context

The `tencentcloud_clb_log_topic` resource manages CLB (Cloud Load Balancer) log topics. It currently supports `log_set_id`, `topic_name`, `status`, and the computed `create_time` fields. The underlying cloud APIs already support tags, but the Terraform resource does not expose this capability.

A key technical detail is that the Create and Read/Update operations span two different SDK packages with different tag structures:
- **CreateTopic** (`clb/v20180317`): uses `[]*TagInfo` with fields `TagKey` / `TagValue`
- **ModifyTopic** (`cls/v20201016`): uses `[]*Tag` with fields `Key` / `Value`
- **DescribeTopics** (`cls/v20201016`): returns `TopicInfo.Tags` as `[]*Tag` with fields `Key` / `Value`

The existing resource implementation already crosses these two SDK packages: `CreateTopic` is called via `ClbService.CreateTopic` (clb client), while `ModifyTopic` and `DescribeTopics` are called via the cls client / `ClsService.DescribeClsTopicById`.

## Goals / Non-Goals

**Goals:**
- Add a single new `tags` parameter to `tencentcloud_clb_log_topic` that lets users configure tags on CLB log topics.
- Ensure tags are passed on Create (via `TagInfo`), updated on Update (via `Tag`), and read back from Describe (via `Tag`).
- Maintain full backward compatibility: the new parameter is optional.

**Non-Goals:**
- Refactoring the existing cross-package (clb/cls) structure of the resource.
- Adding any other parameters besides `tags`.
- Changing the existing `status`, `topic_name`, or `log_set_id` behavior.

## Decisions

### Decision 1: `tags` schema field as a `TypeMap` of string key/value pairs

The Terraform schema will define `tags` as a `TypeMap` (`Elem: &schema.Schema{Type: schema.TypeString}`), where each map key is the tag key and each map value is the tag value.

**Rationale**: Tags are naturally represented as a key/value map. The Read (`DescribeTopics`) and Update (`ModifyTopic`) APIs both use the `Tag` structure with `Key`/`Value` fields, and Create (`CreateTopic`) uses `TagInfo` with `TagKey`/`TagValue`. Using a `TypeMap` keeps state aligned with the API responses and avoids the extra `TypeList` wrapper. The Create and Update code paths map the map entries to `TagInfo.TagKey`/`TagInfo.TagValue` and `Tag.Key`/`Tag.Value` respectively.

**Alternative considered**: Use a `TypeList` of objects with `key`/`value` element fields. Rejected because a map is simpler and matches the key/value nature of tags, eliminating the redundant list wrapper.

### Decision 2: Tags are updatable (not ForceNew)

The `tags` parameter will be `Optional` without `ForceNew`, and the Update method will call `ModifyTopic` with the new tags when `d.HasChange("tags")`.

**Rationale**: The `ModifyTopic` API (`cls/v20201016`) accepts `Tags []*Tag`, confirming tags can be modified after creation without recreating the topic.

### Decision 3: Create passes tags through the existing `ClbService.CreateTopic` params map

The existing `CreateTopic` service method takes a `params map[string]interface{}`. Tags will be passed through this map and converted to `[]*clb.TagInfo` inside the service method, keeping the resource Create function consistent with the existing pattern (which already passes `topic_name` and `partition_count` via the params map).

### Decision 4: Read maps `TopicInfo.Tags` ([]*Tag) back to the schema map

The existing `resourceTencentCloudClbInstanceTopicRead` calls `ClsService.DescribeClsTopicById`, which returns a `*cls.TopicInfo`. The `TopicInfo.Tags` field (`[]*cls.Tag` with `Key`/`Value`) will be flattened into the `tags` schema map (`map[string]string`), with nil-safety checks before setting.

## Risks / Trade-offs

- **[Tag structure mismatch between Create and Update/Read]** → Mitigated by Decision 1: a single `TypeMap` schema field with explicit mapping to the correct SDK structures (`TagInfo` on Create, `Tag` on Update/Read).
- **[Existing `CreateTopic` service method uses a loosely-typed params map]** → Mitigated by passing tags as a `map[string]interface{}` value in the map and handling it explicitly inside the service method; type assertions are guarded.
- **[Read returns empty tags vs. nil]** → Mitigated by nil-checking `res.Tags` before flattening; an empty/nil tag list results in an empty map in state.
