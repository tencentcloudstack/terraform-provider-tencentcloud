package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresDBInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_d_b_instance.d_b_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_d_b_instance.d_b_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresDBInstance = `

resource "tencentcloud_postgres_d_b_instance" "d_b_instance" {
  spec_code = "cdb.pg.sh1.2g"
  storage = 10
  instance_count = 1
  period = 1
  zone = "ap-guangzhou-7"
  charset = "UTF8"
  admin_name = "user"
  admin_password = "password!@#123ABCabc"
  project_id = 0
  d_b_version = ""
  instance_charge_type = "POSTPAID_BY_HOUR"
  auto_voucher = 
  voucher_ids = 
  vpc_id = "vpc-xxxx"
  subnet_id = "subnet-xxxx"
  auto_renew_flag = 
  activity_id = 
  name = ""
  need_support_ipv6 = 
  tag_list {
		tag_key = ""
		tag_value = ""

  }
  security_group_ids = 
  d_b_major_version = ""
  d_b_kernel_version = ""
  d_b_node_set {
		role = ""
		zone = ""

  }
  need_support_t_d_e = 
  k_m_s_key_id = ""
  k_m_s_region = ""
  d_b_engine = ""
  d_b_engine_config = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
