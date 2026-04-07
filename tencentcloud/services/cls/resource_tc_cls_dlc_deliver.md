
Provides a resource to create CLS dcl deliver

Example Usage

```hcl
resource "tencentcloud_cls_dlc_deliver" "example" {
  topic_id     = "0b2e4b58-38aa-4178-9227-c98236d2e100"
  name         = "tf-example"
  deliver_type = 0
  start_time   = 1775118742

  dlc_info {
    table_info {
      data_directory = "DataLakeCategary"
      database_name  = "tf_example_db"
      table_name     = "tf_example_table"
    }

    field_infos {
      cls_field      = "info"
      dlc_field      = "info"
      dlc_field_type = "string"
      disable        = false
    }

    field_infos {
      cls_field      = "int_key"
      dlc_field      = "int_key"
      dlc_field_type = "int"
      disable        = false
    }

    field_infos {
      cls_field      = "bool_key"
      dlc_field      = "bool_key"
      dlc_field_type = "boolean"
      disable        = false
    }

    field_infos {
      cls_field      = "float_key"
      dlc_field      = "float_key"
      dlc_field_type = "float"
      disable        = false
    }

    field_infos {
      cls_field      = "double_key"
      dlc_field      = "double_key"
      dlc_field_type = "double"
      disable        = false
    }


    partition_infos {
      cls_field      = "__TIMESTAMP__"
      dlc_field      = "date_key"
      dlc_field_type = "date"
    }

    partition_extra {
      time_format = "/%Y/%m/%d/%H"
      time_zone   = "UTC+08:00"
    }
  }

  max_size         = 128
  interval         = 300
  has_services_log = 2
}
```

Import

CLS dcl deliver can be imported using the id (topicId#taskId), e.g.

```
terraform import tencentcloud_cls_dlc_deliver.example 0b2e4b58-38aa-4178-9227-c98236d2e100#756a4cbf-b8ae-4866-bbda-1e90aef31760
```
