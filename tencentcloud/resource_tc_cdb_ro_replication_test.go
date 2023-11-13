package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbRoReplicationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbRoReplication,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_ro_replication.ro_replication", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_ro_replication.ro_replication",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbRoReplication = `

resource "tencentcloud_cdb_ro_replication" "ro_replication" {
  instance_id = ""
}

`
