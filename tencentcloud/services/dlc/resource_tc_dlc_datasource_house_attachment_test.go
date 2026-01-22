package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcDatasourceHouseAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcDatasourceHouseAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_datasource_house_attachment.example", "id"),
				),
			},
			{
				Config: testAccDlcDatasourceHouseAttachmentUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_datasource_house_attachment.example", "id"),
				),
			},
		},
	})
}

const testAccDlcDatasourceHouseAttachment = `
resource "tencentcloud_dlc_datasource_house_attachment" "example" {
  datasource_connection_name = "tf-example"
  datasource_connection_type = "Mysql"
  datasource_connection_config {
    mysql {
      location {
        vpc_id            = "vpc-khkyabcd"
        vpc_cidr_block    = "192.168.0.0/16"
        subnet_id         = "subnet-o7n9eg12"
        subnet_cidr_block = "192.168.0.0/24"
      }
    }
  }
  data_engine_names       = ["engine_demo"]
  network_connection_type = 4
  network_connection_desc = "remark."
}
`

const testAccDlcDatasourceHouseAttachmentUpdate = `
resource "tencentcloud_dlc_datasource_house_attachment" "example" {
  datasource_connection_name = "tf-example"
  datasource_connection_type = "Mysql"
  datasource_connection_config {
    mysql {
      location {
        vpc_id            = "vpc-khkyabcd"
        vpc_cidr_block    = "192.168.0.0/16"
        subnet_id         = "subnet-o7n9eg12"
        subnet_cidr_block = "192.168.0.0/24"
      }
    }
  }
  data_engine_names       = ["engine_demo"]
  network_connection_type = 4
  network_connection_desc = "remark update."
}
`
