package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresCloneDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresCloneDBInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_clone_d_b_instance.clone_d_b_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_clone_d_b_instance.clone_d_b_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresCloneDBInstance = `

resource "tencentcloud_postgres_clone_d_b_instance" "clone_d_b_instance" {
  d_b_instance_id = ""
  spec_code = ""
  storage = 
  period = 
  auto_renew_flag = 
  vpc_id = ""
  subnet_id = ""
  name = ""
  instance_charge_type = ""
  security_group_ids = 
  project_id = 
  tag_list {
		tag_key = ""
		tag_value = ""

  }
  d_b_node_set {
		role = ""
		zone = ""

  }
  auto_voucher = 
  voucher_ids = ""
  activity_id = 
  backup_set_id = ""
  recovery_target_time = ""
}

`
