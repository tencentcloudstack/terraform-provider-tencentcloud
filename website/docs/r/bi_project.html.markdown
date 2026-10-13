---
subcategory: "Business Intelligence(BI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bi_project"
sidebar_current: "docs-tencentcloud-resource-bi_project"
description: |-
  Provides a resource to create a bi project
---

# tencentcloud_bi_project

Provides a resource to create a bi project

## Example Usage

```hcl
resource "tencentcloud_bi_project" "project" {
  name       = "terraform_test"
  color_code = "#7BD936"
  logo       = "TF-test"
  mark       = "project mark"
}
```

## Argument Reference

The following arguments are supported:

* `color_code` - (Required, String) Logo background color.
* `name` - (Required, String) Project name.
* `logo` - (Optional, String) Project logo.
* `mark` - (Optional, String) Remark.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateProject](https://cloud.tencent.com/document/api/590/98838) - 创建项目
* Read: [DescribeProjectInfo](https://cloud.tencent.com/document/api/590/98836) - 项目详情
* Update: [ModifyProject](https://cloud.tencent.com/document/api/590/98834) - 修改项目
* Delete: [DeleteProject](https://cloud.tencent.com/document/api/590/98837) - 删除项目

## Import

bi project can be imported using the id, e.g.

```
terraform import tencentcloud_bi_project.project project_id
```

