---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_config_template"
sidebar_current: "docs-tencentcloud-resource-tsf_config_template"
description: |-
  Provides a resource to create a tsf config_template
---

# tencentcloud_tsf_config_template

Provides a resource to create a tsf config_template

## Example Usage

```hcl
resource "tencentcloud_tsf_config_template" "config_template" {
  config_template_name  = "terraform-template-name"
  config_template_type  = "Ribbon"
  config_template_value = <<-EOT
    ribbon.ReadTimeout: 5000
    ribbon.ConnectTimeout: 2000
    ribbon.MaxAutoRetries: 0
    ribbon.MaxAutoRetriesNextServer: 1
    ribbon.OkToRetryOnAllOperations: true
  EOT
  config_template_desc  = "terraform-test"
}
```

## Argument Reference

The following arguments are supported:

* `config_template_name` - (Required, String) Configuration template name.
* `config_template_type` - (Required, String) Configure the microservice framework corresponding to the template.
* `config_template_value` - (Required, String) Configure template data.
* `config_template_desc` - (Optional, String) Configuration template description.
* `program_id_list` - (Optional, Set: [`String`]) Program id list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `config_template_id` - Template Id.
* `create_time` - creation time.
* `update_time` - update time.


