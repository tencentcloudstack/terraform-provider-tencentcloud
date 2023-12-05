Provides a resource to create a cls cos shipper.

Example Usage

```hcl
resource "tencentcloud_cls_cos_shipper" "shipper" {
  bucket       = "preset-scf-bucket-1308919341"
  interval     = 300
  max_size     = 200
  partition    = "/%Y/%m/%d/%H/"
  prefix       = "ap-guangzhou-fffsasad-1649734752"
  shipper_name = "ap-guangzhou-fffsasad-1649734752"
  topic_id     = "4d07fba0-b93e-4e0b-9a7f-d58542560bbb"

  compress {
    format = "lzop"
  }

  content {
    format = "json"

    json {
      enable_tag  = true
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
$ terraform import tencentcloud_cls_cos_shipper.shipper 5d1b7b2a-c163-4c48-bb01-9ee00584d761
```