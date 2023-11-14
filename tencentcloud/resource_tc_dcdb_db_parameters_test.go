package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbDbParametersResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbDbParameters,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_parameters.db_parameters", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_db_parameters.db_parameters",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbDbParameters = `

resource "tencentcloud_dcdb_db_parameters" "db_parameters" {
  instance_id = ""
  sync_mode = 
}

`
