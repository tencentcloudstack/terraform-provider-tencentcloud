---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_backup"
sidebar_current: "docs-tencentcloud-resource-clickhouse_backup"
description: |-
  Provides a resource to open clickhouse backup
---

# tencentcloud_clickhouse_backup

Provides a resource to open clickhouse backup

## Example Usage

```hcl
resource "tencentcloud_clickhouse_backup" "backup" {
  instance_id     = "cdwch-xxxxxx"
  cos_bucket_name = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `cos_bucket_name` - (Required, String) COS bucket name.
* `instance_id` - (Required, String) Instance id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [OpenBackUp](https://www.tencentcloud.com/search/?wd=OpenBackUp&type=document&lang=en)
* Read: [DescribeBackUpSchedule](https://www.tencentcloud.com/search/?wd=DescribeBackUpSchedule&type=document&lang=en)
* Delete: [OpenBackUp](https://www.tencentcloud.com/search/?wd=OpenBackUp&type=document&lang=en)

## Import

clickhouse backup can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_backup.backup instance_id
```

