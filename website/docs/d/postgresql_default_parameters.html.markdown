---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_default_parameters"
sidebar_current: "docs-tencentcloud-datasource-postgresql_default_parameters"
description: |-
  Use this data source to query detailed information of postgresql default_parameters
---

# tencentcloud_postgresql_default_parameters

Use this data source to query detailed information of postgresql default_parameters

## Example Usage

```hcl
data "tencentcloud_postgresql_default_parameters" "default_parameters" {
  db_major_version = "13"
  db_engine        = "postgresql"
}
```

## Argument Reference

The following arguments are supported:

* `db_engine` - (Required, String) Database engine, such as postgresql, mssql_compatible.
* `db_major_version` - (Required, String) The major database version number, such as 11, 12, 13.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `param_info_set` - Parameter informationNote: This field may return null, indicating that no valid values can be obtained.
  * `advanced` - Whether it is a key parameter. Valid values: `true` (yes, and modifying it may affect instance performance), `false` (no)Note: this field may return `null`, indicating that no valid values can be obtained.
  * `classification_cn` - Parameter category in ChineseNote: this field may return `null`, indicating that no valid values can be obtained.
  * `classification_en` - Parameter category in EnglishNote: this field may return `null`, indicating that no valid values can be obtained.
  * `current_value` - Current value of the parameter, which is returned as a stringNote: this field may return `null`, indicating that no valid values can be obtained.
  * `default_value` - Default value of the parameter, which is returned as a stringNote: this field may return `null`, indicating that no valid values can be obtained.
  * `enum_value` - Value range of the enum parameterNote: this field may return `null`, indicating that no valid values can be obtained.
  * `id` - Parameter IDNote: this field may return `null`, indicating that no valid values can be obtained.
  * `last_modify_time` - The last modified time of the parameterNote: this field may return `null`, indicating that no valid values can be obtained.
  * `max` - The maximum value of the `integer` or `real` parameterNote: this field may return `null`, indicating that no valid values can be obtained.
  * `min` - The minimum value of the `integer` or `real` parameterNote: this field may return `null`, indicating that no valid values can be obtained.
  * `name` - Parameter nameNote: this field may return `null`, indicating that no valid values can be obtained.
  * `need_reboot` - Whether to restart the instance for the modified parameter to take effect. Valid values: `true` (yes), `false` (no)Note: this field may return `null`, indicating that no valid values can be obtained.
  * `param_description_ch` - Parameter description in ChineseNote: this field may return `null`, indicating that no valid values can be obtained.
  * `param_description_en` - Parameter description in EnglishNote: this field may return `null`, indicating that no valid values can be obtained.
  * `param_value_type` - Value type of the parameter. Valid values: `integer`, `real` (floating-point), `bool`, `enum`, `mutil_enum` (this type of parameter can be set to multiple enumerated values).For an `integer` or `real` parameter, the `Min` field represents the minimum value and the `Max` field the maximum value. For a `bool` parameter, the valid values include `true` and `false`; For an `enum` or `mutil_enum` parameter, the `EnumValue` field represents the valid values.Note: this field may return `null`, indicating that no valid values can be obtained.
  * `spec_related` - Whether the parameter is related to specifications. Valid values: `true` (yes), `false` (no)Note: this field may return `null`, indicating that no valid values can be obtained.
  * `spec_relation_set` - Associated parameter specification information, which refers to the detailed parameter information of the specifications.Note: This field may return null, indicating that no valid values can be obtained.
    * `enum_value` - Value range of the enum parameterNote: This field may return null, indicating that no valid values can be obtained.
    * `max` - The maximum value of the `integer` or `real` parameterNote: This field may return null, indicating that no valid values can be obtained.
    * `memory` - The specification that corresponds to the parameter informationNote: This field may return null, indicating that no valid values can be obtained.
    * `min` - The minimum value of the `integer` or `real` parameterNote: This field may return null, indicating that no valid values can be obtained.
    * `name` - Parameter nameNote: This field may return null, indicating that no valid values can be obtained.
    * `unit` - Unit of the parameter value. If the parameter has no unit, this field will return null.Note: This field may return null, indicating that no valid values can be obtained.
    * `value` - The default parameter value under this specificationNote: This field may return null, indicating that no valid values can be obtained.
  * `standby_related` - Primary-standby constraint. Valid values: `0` (no constraint), `1` (The parameter value of the standby server must be greater than that of the primary server), `2` (The parameter value of the primary server must be greater than that of the standby server.)Note: This field may return null, indicating that no valid values can be obtained.
  * `unit` - Unit of the parameter value. If the parameter has no unit, this field will return null.Note: This field may return null, indicating that no valid values can be obtained.
  * `version_relation_set` - Associated parameter version information, which refers to the detailed parameter information of the kernel version.Note: This field may return null, indicating that no valid values can be obtained.
    * `db_kernel_version` - The kernel version that corresponds to the parameter informationNote: This field may return null, indicating that no valid values can be obtained.
    * `enum_value` - Value range of the enum parameterNote: This field may return null, indicating that no valid values can be obtained.
    * `max` - The maximum value of the `integer` or `real` parameterNote: This field may return null, indicating that no valid values can be obtained.
    * `min` - The minimum value of the `integer` or `real` parameterNote: This field may return null, indicating that no valid values can be obtained.
    * `name` - Parameter nameNote: This field may return null, indicating that no valid values can be obtained.
    * `unit` - Unit of the parameter value. If the parameter has no unit, this field will return null.Note: This field may return null, indicating that no valid values can be obtained.
    * `value` - Default parameter value under the kernel version and specification of the instanceNote: This field may return null, indicating that no valid values can be obtained.


