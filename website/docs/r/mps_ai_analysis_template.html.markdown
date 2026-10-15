---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_ai_analysis_template"
sidebar_current: "docs-tencentcloud-resource-mps_ai_analysis_template"
description: |-
  Provides a resource to create a mps ai_analysis_template
---

# tencentcloud_mps_ai_analysis_template

Provides a resource to create a mps ai_analysis_template

## Example Usage

```hcl
resource "tencentcloud_mps_ai_analysis_template" "ai_analysis_template" {
  name = "terraform-test"

  classification_configure {
    switch = "OFF"
  }

  cover_configure {
    switch = "ON"
  }

  frame_tag_configure {
    switch = "ON"
  }

  tag_configure {
    switch = "ON"
  }
}
```

## Argument Reference

The following arguments are supported:

* `classification_configure` - (Optional, List) Ai classification task control parameters.
* `comment` - (Optional, String) Ai analysis template description information, length limit: 256 characters.
* `cover_configure` - (Optional, List) Ai cover task control parameters.
* `frame_tag_configure` - (Optional, List) Ai frame tag task control parameters.
* `name` - (Optional, String) Ai analysis template name, length limit: 64 characters.
* `tag_configure` - (Optional, List) Ai tag task control parameters.

The `classification_configure` object supports the following:

* `switch` - (Required, String) Ai classification task switch, optional value:ON/OFF.

The `cover_configure` object supports the following:

* `switch` - (Required, String) Ai cover task switch, optional value:ON/OFF.

The `frame_tag_configure` object supports the following:

* `switch` - (Required, String) Ai frame tag task switch, optional value:ON/OFF.

The `tag_configure` object supports the following:

* `switch` - (Required, String) Ai tag task switch, optional value:ON/OFF.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAIAnalysisTemplate](https://cloud.tencent.com/document/api/862/40249) - 创建内容分析模板
* Read: [DescribeAIAnalysisTemplates](https://cloud.tencent.com/document/api/862/40247) - 获取内容分析模板列表
* Update: [ModifyAIAnalysisTemplate](https://cloud.tencent.com/document/api/862/40246) - 修改内容分析模板
* Delete: [DeleteAIAnalysisTemplate](https://cloud.tencent.com/document/api/862/40248) - 删除内容分析模板

## Import

mps ai_analysis_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_ai_analysis_template.ai_analysis_template ai_analysis_template_id
```

