---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_application"
sidebar_current: "docs-tencentcloud-datasource-tsf_application"
description: |-
  Use this data source to query detailed information of tsf application
---

# tencentcloud_tsf_application

Use this data source to query detailed information of tsf application

## Example Usage

```hcl
data "tencentcloud_tsf_application" "application" {
  application_type  = "V"
  microservice_type = "N"
  # application_resource_type_list = [""]
  application_id_list = ["application-a24x29xv"]
}
```

## Argument Reference

The following arguments are supported:

* `application_id_list` - (Optional, Set: [`String`]) Id list.
* `application_resource_type_list` - (Optional, Set: [`String`]) An array of application resource types.
* `application_type` - (Optional, String) The application type. V OR C, V means VM, C means container.
* `microservice_type` - (Optional, String) The microservice type of the application.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - The application paging list information.
  * `content` - The list of application information.
    * `apigateway_service_id` - gateway service id.
    * `application_desc` - The description of the application.
    * `application_id` - The ID of the application.
    * `application_name` - The name of the application.
    * `application_remark_name` - remark name.
    * `application_resource_type` - application resource type.
    * `application_runtime_type` - application runtime type.
    * `application_type` - The type of the application.
    * `create_time` - create time.
    * `ignore_create_image_repository` - whether ignore create image repository.
    * `microservice_type` - The microservice type of the application.
    * `prog_lang` - Programming language.
    * `service_config_list` - service config list.
      * `health_check` - health check setting.
        * `path` - health check path.
      * `name` - serviceName.
      * `ports` - port list.
        * `protocol` - protocol.
        * `target_port` - service port.
    * `update_time` - update time.
  * `total_count` - The total number of applications.


