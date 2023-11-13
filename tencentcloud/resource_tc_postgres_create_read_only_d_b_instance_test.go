package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresCreateReadOnlyDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresCreateReadOnlyDBInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_create_read_only_d_b_instance.create_read_only_d_b_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_create_read_only_d_b_instance.create_read_only_d_b_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresCreateReadOnlyDBInstance = `

resource "tencentcloud_postgres_create_read_only_d_b_instance" "create_read_only_d_b_instance" {
  spec_code = ""
  storage = 
  instance_count = 
  period = 
  master_d_b_instance_id = ""
  zone = ""
  project_id = 
  d_b_version = ""
  instance_charge_type = ""
  auto_voucher = 
  voucher_ids = 
  auto_renew_flag = 
  vpc_id = ""
  subnet_id = ""
  activity_id = 
  name = ""
  need_support_ipv6 = 
  read_only_group_id = ""
  tag_list {
		tag_key = ""
		tag_value = ""

  }
  security_group_ids = 
}

`
