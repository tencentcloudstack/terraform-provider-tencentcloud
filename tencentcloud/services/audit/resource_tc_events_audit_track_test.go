package audit_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudEventsAuditTrackResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccEventsAuditTrack,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_events_audit_track.example", "id")),
		}, {
			ResourceName:      "tencentcloud_events_audit_track.example",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccEventsAuditTrack = `

resource "tencentcloud_events_audit_track" "example" {
  name = "track_example"

  status                = 1
  track_for_all_members = 0

  storage {
    storage_name   = "393953ac-5c1b-457d-911d-376271b1b4f2"
    storage_prefix = "cloudaudit"
    storage_region = "ap-guangzhou"
    storage_type   = "cls"
  }

  filters {
    resource_fields {
      resource_type = "cam"
      action_type   = "*"
      event_names   = ["AddSubAccount", "AddSubAccountCheckingMFA"]
    }
    resource_fields {
      resource_type = "cvm"
      action_type   = "*"
      event_names   = ["*"]
    }
    resource_fields {
      resource_type = "tke"
      action_type   = "*"
      event_names   = ["*"]
    }
  }
}
`
