package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudAuditTrackResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAuditTrack,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_audit_track.track", "id"),
					resource.TestCheckResourceAttr("tencentcloud_audit_track.track", "name", "terraform_track"),
					resource.TestCheckResourceAttr("tencentcloud_audit_track.track", "action_type", "Read"),
					resource.TestCheckResourceAttr("tencentcloud_audit_track.track", "status", "1"),
					resource.TestCheckResourceAttr("tencentcloud_audit_track.track", "track_for_all_members", "0"),
				),
			},
			{
				ResourceName:      "tencentcloud_audit_track.track",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccAuditTrack = `

resource "tencentcloud_audit_track" "track" {
  action_type           = "Read"
  event_names           = [
    "*",
  ]
  name                  = "terraform_track"
  resource_type         = "*"
  status                = 1
  track_for_all_members = 0

  storage {
    storage_name   = "keep-bucket"
    storage_prefix = "cloudaudit"
    storage_region = "ap-guangzhou"
    storage_type   = "cos"
  }
}
`
