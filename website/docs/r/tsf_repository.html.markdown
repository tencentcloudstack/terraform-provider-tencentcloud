---
subcategory: "Tencent Service Framework(TSF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tsf_repository"
sidebar_current: "docs-tencentcloud-resource-tsf_repository"
description: |-
  Provides a resource to create a TSF repository
---

# tencentcloud_tsf_repository

Provides a resource to create a TSF repository

## Example Usage

```hcl
resource "tencentcloud_tsf_repository" "example" {
  repository_name = "tf-example"
  repository_type = "default"
  bucket_name     = "tf-example-1308919341"
  bucket_region   = "ap-guangzhou"
  directory       = "/"
  repository_desc = "desc."
}
```

## Argument Reference

The following arguments are supported:

* `bucket_name` - (Required, String) the name of the bucket where the warehouse is located.
* `bucket_region` - (Required, String) Bucket region where the warehouse is located.
* `repository_name` - (Required, String) warehouse name.
* `repository_type` - (Required, String) warehouse type (default warehouse: default, private warehouse: private).
* `directory` - (Optional, String) directory.
* `repository_desc` - (Optional, String) warehouse description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - warehouse creation time.
* `is_used` - whether the repository is in use.
* `repository_id` - Warehouse ID.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateRepository](https://www.tencentcloud.com/search/?wd=CreateRepository&type=document&lang=en)
* Read: [DescribeRepository](https://www.tencentcloud.com/search/?wd=DescribeRepository&type=document&lang=en)
* Update: [UpdateRepository](https://www.tencentcloud.com/search/?wd=UpdateRepository&type=document&lang=en)
* Delete: [DeleteRepository](https://www.tencentcloud.com/search/?wd=DeleteRepository&type=document&lang=en)

## Import

TSF repository can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_repository.example repository-xxxxxxxx
```

