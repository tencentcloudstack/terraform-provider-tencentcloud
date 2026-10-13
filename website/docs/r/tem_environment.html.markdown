---
subcategory: "TencentCloud Elastic Microservice(TEM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tem_environment"
sidebar_current: "docs-tencentcloud-resource-tem_environment"
description: |-
  Provides a resource to create a tem environment
---

# tencentcloud_tem_environment

Provides a resource to create a tem environment

## Example Usage

```hcl
resource "tencentcloud_tem_environment" "environment" {
  environment_name = "demo"
  description      = "demo for test"
  vpc              = "vpc-2hfyray3"
  subnet_ids       = ["subnet-rdkj0agk", "subnet-r1c4pn5m", "subnet-02hcj95c"]
  tags = {
    "created" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `environment_name` - (Required, String) environment name.
* `subnet_ids` - (Required, Set: [`String`]) subnet IDs.
* `vpc` - (Required, String) vpc ID.
* `description` - (Optional, String) environment description.
* `tags` - (Optional, Map) environment tag list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateEnvironment](https://cloud.tencent.com/document/api/1371/60150) - 创建环境
* Create: [DescribeEnvironmentStatus](https://cloud.tencent.com/document/api/1371/78168) - 获取环境状态
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Read: [DescribeEnvironment](https://cloud.tencent.com/document/api/1371/78842) - 获取环境基础信息
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - 批量查看资源关联的标签
* Update: [ModifyEnvironment](https://cloud.tencent.com/document/api/1371/60148) - 编辑环境
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Delete: [DestroyEnvironment](https://cloud.tencent.com/document/api/1371/77632) - 销毁环境

## Import

tem environment can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_environment.environment environment_id
```

