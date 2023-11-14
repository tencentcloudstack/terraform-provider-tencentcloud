package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCwpLicenseOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCwpLicenseOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cwp_license_operation.license_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_cwp_license_operation.license_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCwpLicenseOperation = `

resource "tencentcloud_cwp_license_operation" "license_operation" {
  resource_id = ""
  license_type = 
  is_all = 
  quuid_list = 
  task_id = 
  filters {
		name = ""
		values = 

  }
}

`
