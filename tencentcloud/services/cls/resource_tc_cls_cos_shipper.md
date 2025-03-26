Provides a resource to create a cls cos shipper.

Example Usage

```hcl
data "tencentcloud_user_info" "info" {}

locals {
  app_id = data.tencentcloud_user_info.info.app_id
}

resource "tencentcloud_cos_bucket" "example" {
  bucket = "private-bucket-${local.app_id}"
  acl    = "private"
}

resource "tencentcloud_cls_logset" "example" {
  logset_name = "tf-example"
  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cls_topic" "example" {
  topic_name           = "tf-example"
  logset_id            = tencentcloud_cls_logset.example.id
  auto_split           = false
  max_split_partitions = 20
  partition_count      = 1
  period               = 10
  storage_type         = "hot"
  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cls_cos_shipper" "example" {
  bucket       = tencentcloud_cos_bucket.example.id
  topic_id     = tencentcloud_cls_topic.example.id
  interval     = 300
  max_size     = 200
  partition    = "/%Y/%m/%d/%H/"
  prefix       = "ap-guangzhou-fffsasad-1649734752"
  shipper_name = "ap-guangzhou-fffsasad-1649734752"

  compress {
    format = "lzop"
  }

  content {
    format = "json"

    json {
      enable_tag = true
      meta_fields = [
        "__FILENAME__",
        "__SOURCE__",
        "__TIMESTAMP__",
      ]
    }
  }
}
```

Import

cls cos shipper can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_cos_shipper.example 5d1b7b2a-c163-4c48-bb01-9ee00584d761
```