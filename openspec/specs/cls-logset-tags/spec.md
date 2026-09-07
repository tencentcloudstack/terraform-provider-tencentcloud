# cls-logset-tags Specification

## Purpose
TBD - created by archiving change add-cls-logset-tags-param. Update Purpose after archive.
## Requirements
### Requirement: CLS logset native tag binding via TypeList

The `tencentcloud_cls_logset` resource SHALL expose an optional `tags` schema field of `schema.TypeList` containing nested `key` (TypeString, Required) and `value` (TypeString, Required) fields. Each list element maps to one `cls.Tag{Key, Value}` pair. The `tags` field SHALL replace the existing `TypeMap` `tags` field that was previously managed through the TencentCloud Tag Service.

#### Scenario: Create logset with native CLS API tags

- **WHEN** a user creates a `tencentcloud_cls_logset` resource with a `tags` block containing one or more `{key, value}` pairs
- **THEN** the provider SHALL populate `CreateLogsetRequest.Tags` with a `[]*cls.Tag` slice built from the schema list and send it to the CLS `CreateLogset` API
- **AND** the provider SHALL NOT call the TencentCloud Tag Service (`svctag.TagService.ModifyTags`) for tag creation

#### Scenario: Create logset without tags

- **WHEN** a user creates a `tencentcloud_cls_logset` resource without specifying the `tags` field
- **THEN** the provider SHALL NOT set the `Tags` field on `CreateLogsetRequest` (leave it nil)
- **AND** the resource SHALL be created successfully without tags

#### Scenario: Read tags from DescribeLogsets response

- **WHEN** the provider reads a `tencentcloud_cls_logset` via the `DescribeLogsets` API and the response `LogsetInfo.Tags` is non-empty
- **THEN** the provider SHALL convert the `[]*cls.Tag` response into a `[]map[string]interface{}` list (each map containing `key` and `value`) and set it on the `tags` schema field via `d.Set("tags", ...)`
- **AND** the provider SHALL NOT call `svctag.TagService.DescribeResourceTags` for tag reading

#### Scenario: Read logset with no tags

- **WHEN** the provider reads a `tencentcloud_cls_logset` and the response `LogsetInfo.Tags` is nil or empty
- **THEN** the provider SHALL set an empty list on the `tags` schema field
- **AND** no error SHALL be returned

#### Scenario: Update tags via ModifyLogset API

- **WHEN** a user updates the `tags` field of an existing `tencentcloud_cls_logset` (change detected via `d.HasChange("tags")`)
- **THEN** the provider SHALL populate `ModifyLogsetRequest.Tags` with the full new `[]*cls.Tag` list from the schema and call the CLS `ModifyLogset` API
- **AND** the provider SHALL NOT call `svctag.TagService.ModifyTags` for tag updates

#### Scenario: Update logset_name without changing tags

- **WHEN** a user updates only the `logset_name` field and the `tags` field is unchanged
- **THEN** the provider SHALL call `ModifyLogset` with only the `LogsetId` and `LogsetName` fields
- **AND** the `Tags` field on `ModifyLogsetRequest` SHALL NOT be populated (to avoid clobbering existing tags)

#### Scenario: Tag field schema structure

- **WHEN** the `tencentcloud_cls_logset` resource schema is defined
- **THEN** the `tags` field SHALL be `schema.TypeList`, `Optional: true`, with `Elem` of `schema.Resource` containing `"key"` (TypeString, Required) and `"value"` (TypeString, Required)
- **AND** each nested field SHALL have a description documenting its purpose
