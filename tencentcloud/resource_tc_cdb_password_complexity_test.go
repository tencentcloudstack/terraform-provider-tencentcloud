package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbPasswordComplexityResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbPasswordComplexity,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_password_complexity.password_complexity", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_password_complexity.password_complexity",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbPasswordComplexity = `

resource "tencentcloud_cdb_password_complexity" "password_complexity" {
  instance_ids = 
  param_list {
		name = ""
		current_value = ""

  }
}

`
