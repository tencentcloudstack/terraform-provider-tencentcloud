package thpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudThpcWorkspacesResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccThpcWorkspaces,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_thpc_workspaces.thpc_workspaces", "id")),
		}, {
			ResourceName:      "tencentcloud_thpc_workspaces.thpc_workspaces",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccThpcWorkspaces = `

resource "tencentcloud_thpc_workspaces" "thpc_workspaces" {
  placement = {
  }
  space_charge_prepaid = {
  }
  system_disk = {
  }
  data_disks = {
  }
  virtual_private_cloud = {
  }
  internet_accessible = {
  }
  login_settings = {
  }
  enhanced_service = {
    security_service = {
    }
    monitor_service = {
    }
    automation_service = {
    }
  }
  tag_specification = {
    tags = {
    }
  }
}
`
