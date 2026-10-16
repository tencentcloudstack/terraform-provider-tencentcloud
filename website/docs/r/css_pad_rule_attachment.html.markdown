---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_pad_rule_attachment"
sidebar_current: "docs-tencentcloud-resource-css_pad_rule_attachment"
description: |-
  Provides a resource to create a css pad_rule_attachment
---

# tencentcloud_css_pad_rule_attachment

Provides a resource to create a css pad_rule_attachment

## Example Usage

```hcl
resource "tencentcloud_css_pad_rule_attachment" "pad_rule_attachment" {
  domain_name = "177154.push.tlivecloud.com"
  template_id = 17067
  app_name    = "qqq"
  stream_name = "ppp"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, String, ForceNew) Push domain.
* `template_id` - (Required, Int, ForceNew) Template id.
* `app_name` - (Optional, String, ForceNew) Push path, must same with play path, default is live.
* `stream_name` - (Optional, String, ForceNew) Stream name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateLivePadRule](https://www.tencentcloud.com/document/api/267/71832)
* Read: [DescribeLivePadRules](https://www.tencentcloud.com/search/?wd=DescribeLivePadRules&type=document&lang=en)
* Delete: [DeleteLivePadRule](https://www.tencentcloud.com/search/?wd=DeleteLivePadRule&type=document&lang=en)
* Delete: [DescribeLivePadRules](https://www.tencentcloud.com/search/?wd=DescribeLivePadRules&type=document&lang=en)

## Import

css pad_rule_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_css_pad_rule_attachment.pad_rule_attachment templateId#domainName
```

