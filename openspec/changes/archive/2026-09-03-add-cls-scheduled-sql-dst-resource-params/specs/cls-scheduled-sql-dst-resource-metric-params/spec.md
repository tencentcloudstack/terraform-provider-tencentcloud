## ADDED Requirements

### Requirement: metric_names parameter in dst_resource block of tencentcloud_cls_scheduled_sql
The `tencentcloud_cls_scheduled_sql` resource SHALL support an optional `metric_names` parameter of type list of strings within the `dst_resource` block. This parameter maps to `DstResource.MetricNames` in the CreateScheduledSql and ModifyScheduledSql APIs. When `biz_type` is 1 (metric topic), multiple metric names can be specified in this field.

#### Scenario: Create scheduled SQL task with metric_names
- **WHEN** a user creates a `tencentcloud_cls_scheduled_sql` resource with `dst_resource.metric_names` set to a list of metric name strings
- **THEN** the CreateScheduledSql API SHALL be called with `DstResource.MetricNames` populated as a string array, and the resource SHALL be created successfully

#### Scenario: Read scheduled SQL task populates metric_names
- **WHEN** the Read method is called for a resource whose `DstResource.MetricNames` is non-nil
- **THEN** the `metric_names` field in the `dst_resource` block SHALL be populated with the list of metric name strings from the DescribeScheduledSqlInfo response

#### Scenario: Update scheduled SQL task with metric_names
- **WHEN** a user updates `dst_resource.metric_names` in an existing `tencentcloud_cls_scheduled_sql` resource
- **THEN** the ModifyScheduledSql API SHALL be called with `DstResource.MetricNames` populated with the new list of metric name strings

#### Scenario: Create without metric_names
- **WHEN** a user creates a `tencentcloud_cls_scheduled_sql` resource without specifying `dst_resource.metric_names`
- **THEN** the CreateScheduledSql API SHALL be called without the `MetricNames` field in `DstResource`, and no error SHALL occur

### Requirement: metric_labels parameter in dst_resource block of tencentcloud_cls_scheduled_sql
The `tencentcloud_cls_scheduled_sql` resource SHALL support an optional `metric_labels` parameter of type list of strings within the `dst_resource` block. This parameter maps to `DstResource.MetricLabels` in the CreateScheduledSql and ModifyScheduledSql APIs. It specifies metric dimensions and does not accept time-type labels.

#### Scenario: Create scheduled SQL task with metric_labels
- **WHEN** a user creates a `tencentcloud_cls_scheduled_sql` resource with `dst_resource.metric_labels` set to a list of label strings
- **THEN** the CreateScheduledSql API SHALL be called with `DstResource.MetricLabels` populated as a string array

#### Scenario: Read scheduled SQL task populates metric_labels
- **WHEN** the Read method is called for a resource whose `DstResource.MetricLabels` is non-nil
- **THEN** the `metric_labels` field in the `dst_resource` block SHALL be populated with the list of label strings from the DescribeScheduledSqlInfo response

#### Scenario: Update scheduled SQL task with metric_labels
- **WHEN** a user updates `dst_resource.metric_labels` in an existing resource
- **THEN** the ModifyScheduledSql API SHALL be called with `DstResource.MetricLabels` populated

### Requirement: custom_time parameter in dst_resource block of tencentcloud_cls_scheduled_sql
The `tencentcloud_cls_scheduled_sql` resource SHALL support an optional `custom_time` parameter of type string within the `dst_resource` block. This parameter maps to `DstResource.CustomTime` in the CreateScheduledSql and ModifyScheduledSql APIs. It specifies the metric timestamp field; the default is the left boundary time of the SQL query range.

#### Scenario: Create scheduled SQL task with custom_time
- **WHEN** a user creates a `tencentcloud_cls_scheduled_sql` resource with `dst_resource.custom_time` set to a field name string
- **THEN** the CreateScheduledSql API SHALL be called with `DstResource.CustomTime` set to that string

#### Scenario: Read scheduled SQL task populates custom_time
- **WHEN** the Read method is called for a resource whose `DstResource.CustomTime` is non-nil
- **THEN** the `custom_time` field in the `dst_resource` block SHALL be populated with the string value from the DescribeScheduledSqlInfo response

#### Scenario: Update scheduled SQL task with custom_time
- **WHEN** a user updates `dst_resource.custom_time` in an existing resource
- **THEN** the ModifyScheduledSql API SHALL be called with `DstResource.CustomTime` set to the new value

### Requirement: custom_metric_labels parameter in dst_resource block of tencentcloud_cls_scheduled_sql
The `tencentcloud_cls_scheduled_sql` resource SHALL support an optional `custom_metric_labels` parameter of type list of objects within the `dst_resource` block. Each object SHALL contain `key` (string) and `value` (string) fields. This parameter maps to `DstResource.CustomMetricLabels` (a list of `MetricLabel` structs) in the CreateScheduledSql and ModifyScheduledSql APIs. It allows users to add static dimensions to metrics beyond the `metric_labels` parameter.

#### Scenario: Create scheduled SQL task with custom_metric_labels
- **WHEN** a user creates a `tencentcloud_cls_scheduled_sql` resource with `dst_resource.custom_metric_labels` set to a list of objects, each with `key` and `value` string fields
- **THEN** the CreateScheduledSql API SHALL be called with `DstResource.CustomMetricLabels` populated as an array of `MetricLabel` structs, each containing `Key` and `Value`

#### Scenario: Read scheduled SQL task populates custom_metric_labels
- **WHEN** the Read method is called for a resource whose `DstResource.CustomMetricLabels` is non-nil
- **THEN** the `custom_metric_labels` field in the `dst_resource` block SHALL be populated with a list of objects, each containing `key` and `value` from the DescribeScheduledSqlInfo response

#### Scenario: Update scheduled SQL task with custom_metric_labels
- **WHEN** a user updates `dst_resource.custom_metric_labels` in an existing resource
- **THEN** the ModifyScheduledSql API SHALL be called with `DstResource.CustomMetricLabels` populated with the new list of `MetricLabel` structs

#### Scenario: Create without custom_metric_labels
- **WHEN** a user creates a `tencentcloud_cls_scheduled_sql` resource without specifying `dst_resource.custom_metric_labels`
- **THEN** the CreateScheduledSql API SHALL be called without the `CustomMetricLabels` field in `DstResource`, and no error SHALL occur
