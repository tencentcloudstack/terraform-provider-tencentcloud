package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataCodePermissionsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWedataCodePermissions,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_permissions.example", "id"),
				),
			},
			{
				Config: testAccWedataCodePermissionsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_wedata_code_permissions.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_wedata_code_permissions.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWedataCodePermissions = `
resource "tencentcloud_wedata_code_permissions" "example" {
  project_id = "3108707295180644352"
  authorize_permission_objects {
    resource {
      resource_type        = ""
      resource_id          = ""
      resource_id_for_path = ""
      resource_cfs_path    = ""
    }

    authorize_subjects {
      subject_type   = ""
      subject_values = []
      privileges     = []
    }
  }
}
`

const testAccWedataCodePermissionsUpdate = `
resource "tencentcloud_wedata_code_permissions" "example" {
  project_id = "3108707295180644352"
  authorize_permission_objects {
    resource {
      resource_type        = ""
      resource_id          = ""
      resource_id_for_path = ""
      resource_cfs_path    = ""
    }

    authorize_subjects {
      subject_type   = ""
      subject_values = []
      privileges     = []
    }
  }
}
`
