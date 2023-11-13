package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCloudauditTrackResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudauditTrack,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cloudaudit_track.track", "id")),
			},
			{
				ResourceName:      "tencentcloud_cloudaudit_track.track",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCloudauditTrack = `

resource "tencentcloud_cloudaudit_track" "track" {
  track_id = &lt;nil&gt;
  name = &lt;nil&gt;
  action_type = &lt;nil&gt;
  resource_type = &lt;nil&gt;
  status = &lt;nil&gt;
  event_names = &lt;nil&gt;
  storage {
		storage_type = &lt;nil&gt;
		storage_region = &lt;nil&gt;
		storage_name = &lt;nil&gt;
		storage_prefix = &lt;nil&gt;

  }
  track_for_all_members = &lt;nil&gt;
  }

`
