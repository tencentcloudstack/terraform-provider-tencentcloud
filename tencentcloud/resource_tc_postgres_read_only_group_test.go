package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresReadOnlyGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresReadOnlyGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgres_read_only_group.read_only_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgres_read_only_group.read_only_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresReadOnlyGroup = `

resource "tencentcloud_postgres_read_only_group" "read_only_group" {
  master_d_b_instance_id = "postgres-xxxx"
  name = "test-rg"
  project_id = 0
  vpc_id = "vpc-e0tfm161"
  subnet_id = "subnet-443a3lv6"
  replay_lag_eliminate = 0
  replay_latency_eliminate = 0
  max_replay_lag = 5000
  max_replay_latency = 32
  min_delay_eliminate_reserve = 1
  security_group_ids = 
}

`
