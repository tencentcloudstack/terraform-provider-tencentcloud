---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_parameter_template"
sidebar_current: "docs-tencentcloud-resource-postgresql_parameter_template"
description: |-
  Provides a resource to create a postgresql parameter_template
---

# tencentcloud_postgresql_parameter_template

Provides a resource to create a postgresql parameter_template

## Example Usage

```hcl
resource "tencentcloud_postgresql_parameter_template" "parameter_template" {
  template_name        = "your_temp_name"
  db_major_version     = "13"
  db_engine            = "postgresql"
  template_description = "For_tf_test"

  modify_param_entry_set {
    name           = "timezone"
    expected_value = "UTC"
  }
  modify_param_entry_set {
    name           = "lock_timeout"
    expected_value = "123"
  }

  delete_param_set = ["lc_time"]
}
```

## Argument Reference

The following arguments are supported:

* `db_engine` - (Required, String) Database engine, such as postgresql, mssql_compatible.
* `db_major_version` - (Required, String) The major database version number, such as 11, 12, 13.
* `template_name` - (Required, String) Template name, which can contain 1-60 letters, digits, and symbols (-_./()+=:@).
* `delete_param_set` - (Optional, Set: [`String`]) The set of parameters that need to be deleted.
* `modify_param_entry_set` - (Optional, List) The set of parameters that need to be modified or added. Note: the same parameter cannot appear in the set of modifying and adding and deleting at the same time.
* `template_description` - (Optional, String) Parameter template description, which can contain 1-60 letters, digits, and symbols (-_./()+=:@).

The `modify_param_entry_set` object supports the following:

* `expected_value` - (Required, String) Modify the parameter value. The input parameters are passed in the form of strings, for example: decimal `0.1`, integer `1000`, enumeration `replica`.
* `name` - (Required, String) The parameter name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `param_info_set` - Parameter information contained in the parameter template. Note: This field may return null, indicating that no valid value can be obtained.
  * `advanced` - Whether it is a key parameter. (true is the key parameter, the modification needs to be paid attention to, which may affect the performance of the instance) Note: This field may return null, indicating that no valid value can be obtained.
  * `classification_cn` - Parameter Chinese classification. Note: This field may return null, indicating that no valid value can be obtained.
  * `classification_en` - Parameter English Classification. Note: This field may return null, indicating that no valid value can be obtained.
  * `current_value` - The current running value of the parameter. returned as a string. Note: This field may return null, indicating that no valid value can be obtained.
  * `default_value` - Parameter default value. returned as a string. Note: This field may return null, indicating that no valid value can be obtained.
  * `enum_value` - Enumeration type parameter, value range. Note: This field may return null, indicating that no valid value can be obtained.
  * `id` - parameter ID. Note: This field may return null, indicating that no valid value can be obtained.
  * `last_modify_time` - parameter last modified time. Note: This field may return null, indicating that no valid value can be obtained.
  * `max` - Numerical type (integer, real) parameters, value lower bound. Note: This field may return null, indicating that no valid value can be obtained.
  * `min` - Numeric type (integer, real) parameters, value upper bound. Note: This field may return null, indicating that no valid value can be obtained.
  * `name` - parameter name. Note: This field may return null, indicating that no valid value can be obtained.
  * `need_reboot` - Parameter modification, whether to restart to take effect. (true is required, false is not required) Note: This field may return null, indicating that no valid value can be obtained.
  * `param_description_ch` - Parameter Chinese description. Note: This field may return null, indicating that no valid value can be obtained.
  * `param_description_en` - Parameter English description. Note: This field may return null, indicating that no valid value can be obtained.
  * `param_value_type` - Parameter value type: integer (integer), real (floating point), bool (Boolean), enum (enumeration type), mutil_enum (enumeration type, support multiple choices).When the parameter type is For integer (integer type) and real (floating point type), the value range of the parameter is determined according to the Max and Min of the return value; When the parameter type is bool (Boolean type), the value range of the parameter setting value is true | false; When the parameter type is enum (enumeration type) or mutil_enum (multi-enumeration type), the value range of the parameter is determined by the EnumValue in the return value. Note: This field may return null, indicating that it cannot be fetched to a valid value.
  * `spec_related` - Is it related to the specification. (true means related, false means don&#39;t want to close) Note: This field may return null, indicating that no valid value can be obtained.
  * `spec_relation_set` - Parameter specification related information, storing specific parameter information under specific specifications. Note: This field may return null, indicating that no valid value can be obtained.
    * `enum_value` - Enumeration type parameter, value range. Note: This field may return null, indicating that no valid value can be obtained.
    * `max` - Numeric type (integer, real) parameters, value upper bound. Note: This field may return null, indicating that no valid value can be obtained.
    * `memory` - The parameter information belongs to the specification. Note: This field may return null, indicating that no valid value can be obtained.
    * `min` - Numerical type (integer, real) parameters, value lower bound. Note: This field may return null, indicating that no valid value can be obtained.
    * `name` - parameter name. Note: This field may return null, indicating that no valid value can be obtained.
    * `unit` - The unit of the parameter value. When the parameter has no unit, the field returns empty. Note: This field may return null, indicating that no valid value can be obtained.
    * `value` - The default value of the parameter under this specification. Note: This field may return null, indicating that no valid value can be obtained.
  * `standby_related` - There are master-standby constraints on the parameters, 0: no master-standby constraint relationship, 1: the parameter value of the standby machine must be greater than that of the master machine, 2: the parameter value of the master machine must be greater than that of the standby machine. Note: This field may return null, indicating that no valid value can be obtained.
  * `unit` - Parameter Value Unit. When the parameter has no unit, the field returns empty. Note: This field may return null, indicating that no valid value can be obtained.
  * `version_relation_set` - Parameter version association information, storing specific parameter information under a specific kernel version. Note: This field may return null, indicating that no valid value can be obtained.
    * `db_kernel_version` - The kernel version of the parameter information. Note: This field may return null, indicating that no valid value can be obtained.
    * `enum_value` - Enumeration type parameter, value range. Note: This field may return null, indicating that no valid value can be obtained.
    * `max` - Numeric type (integer, real) parameters, value upper bound. Note: This field may return null, indicating that no valid value can be obtained.
    * `min` - Numerical type (integer, real) parameters, value lower bound. Note: This field may return null, indicating that no valid value can be obtained.
    * `name` - parameter name. Note: This field may return null, indicating that no valid value can be obtained.
    * `unit` - The unit of the parameter value. When the parameter has no unit, the field returns empty. Note: This field may return null, indicating that no valid value can be obtained.
    * `value` - The default value of the parameter in this version and this specification. Note: This field may return null, indicating that no valid value can be obtained.


## Import

postgresql parameter_template can be imported using the id, e.g.

Notice: `modify_param_entry_set` and `delete_param_set` do not support import.

```
terraform import tencentcloud_postgresql_parameter_template.parameter_template parameter_template_id
```

