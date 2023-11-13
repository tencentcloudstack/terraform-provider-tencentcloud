package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_instance.instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbInstance = `

resource "tencentcloud_cynosdb_instance" "instance" {
  cluster_id = "cynosdbmysql-6gtlgm5l"
  cpu = 2
  memory = 4
  read_only_count = 1
  instance_grp_id = "cynosmysql-grp-xxxxxxxx"
  vpc_id = "vpc-1ptuei0b"
  subnet_id = "subnet-1tmw9t4o"
  port = 2000
  instance_name = "cynosmysql-xxxxxxxx"
  auto_voucher = 0
  db_type = "MYSQL"
  order_source = "api"
  deal_mode = 0
  param_template_id = 0
  instance_params {
		param_name = ""
		current_value = ""
		old_value = ""

  }
  security_group_ids = 
}

`
