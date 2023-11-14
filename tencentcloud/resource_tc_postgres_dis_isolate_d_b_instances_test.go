package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresDisIsolateDBInstancesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresDisIsolateDBInstances,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_dis_isolate_d_b_instances.dis_isolate_d_b_instances", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_dis_isolate_d_b_instances.dis_isolate_d_b_instances",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresDisIsolateDBInstances = `

resource "tencentcloud_postgres_dis_isolate_d_b_instances" "dis_isolate_d_b_instances" {
  d_b_instance_id_set = &lt;nil&gt;
  period = 12
  auto_voucher = false
  voucher_ids = &lt;nil&gt;
  tags = {
    "createdBy" = "terraform"
  }
}

`
