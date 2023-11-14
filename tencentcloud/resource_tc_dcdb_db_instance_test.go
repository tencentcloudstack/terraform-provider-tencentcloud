package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbDbInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbDbInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_instance.db_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_db_instance.db_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbDbInstance = `

resource "tencentcloud_dcdb_db_instance" "db_instance" {
  zones = &lt;nil&gt;
  period = &lt;nil&gt;
  shard_memory = &lt;nil&gt;
  shard_storage = &lt;nil&gt;
  shard_node_count = &lt;nil&gt;
  shard_count = &lt;nil&gt;
  count = &lt;nil&gt;
  project_id = &lt;nil&gt;
  vpc_id = &lt;nil&gt;
  subnet_id = &lt;nil&gt;
  db_version_id = "5.7.17"
  auto_voucher = &lt;nil&gt;
  voucher_ids = &lt;nil&gt;
  security_group_id = &lt;nil&gt;
  instance_name = &lt;nil&gt;
  ipv6_flag = &lt;nil&gt;
  resource_tags {
		tag_key = &lt;nil&gt;
		tag_value = &lt;nil&gt;

  }
  init_params {
		param = &lt;nil&gt;
		value = ""

  }
  dcn_region = &lt;nil&gt;
  dcn_instance_id = &lt;nil&gt;
  auto_renew_flag = &lt;nil&gt;
  security_group_ids = &lt;nil&gt;
}

`
