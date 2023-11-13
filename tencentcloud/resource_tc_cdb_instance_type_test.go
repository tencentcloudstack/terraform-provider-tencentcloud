package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbInstanceTypeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbInstanceType,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_instance_type.instance_type", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_instance_type.instance_type",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbInstanceType = `

resource "tencentcloud_cdb_instance_type" "instance_type" {
  instance_id = ""
}

`
