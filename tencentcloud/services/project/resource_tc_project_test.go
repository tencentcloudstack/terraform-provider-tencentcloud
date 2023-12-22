package project_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixProjectResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProject,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_project.project", "id")),
			},
			{
				ResourceName:      "tencentcloud_project.project",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccProject = `

resource "tencentcloud_project" "project" {
  project_name = "terraform-test"
  info         = "for terraform test"
}

`
