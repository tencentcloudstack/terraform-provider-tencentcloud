package cls_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudClsDlcDeliverResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsDlcDeliver,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_dlc_deliver.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_dlc_deliver.example", "task_id"),
					resource.TestCheckResourceAttr("tencentcloud_cls_dlc_deliver.example", "name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_cls_dlc_deliver.example", "deliver_type", "0"),
				),
			},
			{
				Config: testAccClsDlcDeliverUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_dlc_deliver.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_dlc_deliver.example", "task_id"),
					resource.TestCheckResourceAttr("tencentcloud_cls_dlc_deliver.example", "name", "tf-example-update"),
					resource.TestCheckResourceAttr("tencentcloud_cls_dlc_deliver.example", "deliver_type", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_dlc_deliver.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsDlcDeliver = `
resource "tencentcloud_cls_dlc_deliver" "example" {
  topic_id     = "5ba3b3eb-7459-4807-82d9-c98236d2e100"
  name         = "tf-example"
  deliver_type = 0
  start_time   = 1775118742

  dlc_info {
    table_info {
      data_directory = "DataLakeCategary"
      database_name  = "alexhliu_test"
      table_name     = "alexhliu_inner_table"
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
`

const testAccClsDlcDeliverUpdate = `
resource "tencentcloud_cls_dlc_deliver" "example" {
  topic_id     = "5ba3b3eb-7459-4807-82d9-c98236d2e100"
  name         = "tf-example-update"
  deliver_type = 0
  start_time   = 1775118742

  dlc_info {
    table_info {
      data_directory = "DataLakeCategary"
      database_name  = "alexhliu_test"
      table_name     = "alexhliu_inner_table"
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
`
