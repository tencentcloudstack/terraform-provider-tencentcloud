package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbDeployGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbDeployGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_deploy_group.deploy_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_deploy_group.deploy_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbDeployGroup = `

resource "tencentcloud_cdb_deploy_group" "deploy_group" {
  deploy_group_name = &lt;nil&gt;
  description = &lt;nil&gt;
  limit_num = &lt;nil&gt;
  dev_class = 
}

`
