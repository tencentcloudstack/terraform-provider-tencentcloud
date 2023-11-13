package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverGeneralCloudInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralCloudInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_cloud_instance.general_cloud_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_general_cloud_instance.general_cloud_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverGeneralCloudInstance = `

resource "tencentcloud_sqlserver_general_cloud_instance" "general_cloud_instance" {
  zone = "ap-guangzhou-1"
  memory = 
  storage = 
  cpu = 
  machine_type = "CLOUD_SSD"
  instance_charge_type = "postpaid"
  project_id = 
  goods_num = 1
  subnet_id = "subnet-bdoe83fa"
  vpc_id = "vpc-dsp338hz"
  period = 
  auto_voucher = 
  voucher_ids = 
  d_b_version = ""
  auto_renew_flag = 
  security_group_list = 
  weekly = 
  start_time = ""
  span = 
  multi_zones = 
  resource_tags {
		tag_key = ""
		tag_value = ""

  }
  collation = ""
  time_zone = ""
}

`
