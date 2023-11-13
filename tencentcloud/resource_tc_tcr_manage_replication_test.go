package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrManageReplicationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrManageReplication,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcr_manage_replication.manage_replication", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcr_manage_replication.manage_replication",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrManageReplication = `

resource "tencentcloud_tcr_manage_replication" "manage_replication" {
  source_registry_id = "tcr-xxx"
  destination_registry_id = "tcr-xxx"
  rule {
		name = "test"
		dest_namespace = "test"
		override = false
		filters {
			type = "tag"
			value = ""
		}

  }
  description = "this is the tcr rule"
  destination_region_id = 1
  peer_replication_option {
		peer_registry_uin = "113498"
		peer_registry_token = "xxx"
		enable_peer_replication = true

  }
}

`
