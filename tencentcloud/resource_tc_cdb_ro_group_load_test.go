package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbRoGroupLoadResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbRoGroupLoad,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_ro_group_load.ro_group_load", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_ro_group_load.ro_group_load",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbRoGroupLoad = `

resource "tencentcloud_cdb_ro_group_load" "ro_group_load" {
  ro_group_id = ""
}

`
