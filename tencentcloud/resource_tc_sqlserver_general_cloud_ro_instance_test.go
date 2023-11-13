package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSqlserverGeneralCloudRoInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralCloudRoInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_cloud_ro_instance.general_cloud_ro_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_general_cloud_ro_instance.general_cloud_ro_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSqlserverGeneralCloudRoInstance = `

resource "tencentcloud_sqlserver_general_cloud_ro_instance" "general_cloud_ro_instance" {
  instance_id = ""
  zone = ""
  read_only_group_type = 
  memory = 
  storage = 
  cpu = 
  machine_type = ""
  read_only_group_forced_upgrade = 
  read_only_group_id = ""
  read_only_group_name = ""
  read_only_group_is_offline_delay = 
  read_only_group_max_delay_time = 
  read_only_group_min_in_group = 
  instance_charge_type = ""
  goods_num = 
  subnet_id = ""
  vpc_id = ""
  period = 
  security_group_list = 
  auto_voucher = 
  voucher_ids = 
  resource_tags {
		tag_key = ""
		tag_value = ""

  }
  collation = ""
  time_zone = ""
}

`
