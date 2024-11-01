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
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_events_audit_track.events_audit_track", "id")),
		}, {
			ResourceName:      "tencentcloud_events_audit_track.events_audit_track",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccEventsAuditTrack = `

resource "tencentcloud_events_audit_track" "events_audit_track" {
  storage = {
  }
  filters = {
    resource_fields = {
    }
  }
}
`
