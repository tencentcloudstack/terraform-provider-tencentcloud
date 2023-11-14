package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbInstanceNameResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbInstanceName,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_instance_name.instance_name", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_instance_name.instance_name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbInstanceName = `

resource "tencentcloud_cynosdb_instance_name" "instance_name" {
  instance_id = "cynosdb-ins-dokydbam"
  instance_name = "newName"
}

`
