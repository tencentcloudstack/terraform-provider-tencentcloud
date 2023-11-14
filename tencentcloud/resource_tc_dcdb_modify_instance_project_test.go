package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbModifyInstanceProjectResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbModifyInstanceProject,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_modify_instance_project.modify_instance_project", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_modify_instance_project.modify_instance_project",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbModifyInstanceProject = `

resource "tencentcloud_dcdb_modify_instance_project" "modify_instance_project" {
  instance_ids = 
  project_id = 
}

`
