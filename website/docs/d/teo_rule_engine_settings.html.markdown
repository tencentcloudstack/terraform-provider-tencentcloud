---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_rule_engine_settings"
sidebar_current: "docs-tencentcloud-datasource-teo_rule_engine_settings"
description: |-
  Use this data source to query detailed information of teo ruleEngineSettings
---

# tencentcloud_teo_rule_engine_settings

Use this data source to query detailed information of teo ruleEngineSettings

## Example Usage

```hcl
data "tencentcloud_teo_rule_engine_settings" "ruleEngineSettings" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `actions` - Detail info of actions which can be used in rule engine.
  * `action` - Action name.
  * `properties` - Action properties.
    * `choice_properties` - Associative properties of this property, they are all required. Note: This field may return null, indicating that no valid value can be obtained.
      * `choices_value` - The choices which can bse used. This list may be empty.
      * `extra_parameter` - Special parameter. Note: This field may return null, indicating that no valid value can be obtained.
        * `choices` - A list of choices which can be used when `Type` is `CHOICE`.
        * `id` - Parameter name. Valid values:- `Action`: this extra parameter is required when modify HTTP header, this action should be a `RewriteAction`.- `StatusCode`: this extra parameter is required when modify HTTP status code, this action should be a `CodeAction`.- `NULL`: this action should be a `NormalAction`.
        * `type` - Parameter value type. Valid values:- `CHOICE`: select one value from `Choices`.- `CUSTOM_NUM`: integer value.- `CUSTOM_STRING`: string value.
      * `is_allow_empty` - Whether this property is allowed to set empty.
      * `is_multiple` - Whether this property is allowed to set multiple values.
      * `max` - Max integer value can bse used when property type is `CUSTOM_NUM`. When `Min` and `Max` both are 0, this field is meaningless.
      * `min` - Min integer value can bse used when property type is `CUSTOM_NUM`. When `Min` and `Max` both are 0, this field is meaningless.
      * `name` - Property name.
      * `type` - Property value type. Valid values:- `CHOICE`: enum type, must select one of the value in `ChoicesValue`.- `TOGGLE`: switch type, must select one of the value in `ChoicesValue`.- `CUSTOM_NUM`: integer type.- `CUSTOM_STRING`: string type.
    * `choices_value` - The choices which can be used. This list may be empty.
    * `extra_parameter` - Special parameter. Note: This field may return null, indicating that no valid value can be obtained.
      * `choices` - A list of choices which can be used when `Type` is `CHOICE`.
      * `id` - Parameter name. Valid values:- `Action`: this extra parameter is required when modify HTTP header, this action should be a `RewriteAction`.- `StatusCode`: this extra parameter is required when modify HTTP status code, this action should be a `CodeAction`.- `NULL`: this action should be a `NormalAction`.
      * `type` - Parameter value type. Valid values:- `CHOICE`: select one value from `Choices`.- `CUSTOM_NUM`: integer value.- `CUSTOM_STRING`: string value.
    * `is_allow_empty` - Whether this property is allowed to set empty.
    * `is_multiple` - Whether this property is allowed to set multiple values.
    * `max` - Max integer value can bse used when property type is `CUSTOM_NUM`. When `Min` and `Max` both are 0, this field is meaningless.
    * `min` - Min integer value can bse used when property type is `CUSTOM_NUM`. When `Min` and `Max` both are 0, this field is meaningless.
    * `name` - Property name.
    * `type` - Property value type. Valid values:- `CHOICE`: enum type, must select one of the value in `ChoicesValue`.- `TOGGLE`: switch type, must select one of the value in `ChoicesValue`.- `OBJECT`: object type, the `ChoiceProperties` list all properties of the object.- `CUSTOM_NUM`: integer type.- `CUSTOM_STRING`: string type.


