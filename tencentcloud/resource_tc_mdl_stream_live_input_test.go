package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMdlStreamLiveInputResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_INTERNATIONAL) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMdlStreamLiveInput,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mdl_stream_live_input.stream_live_input", "id")),
			},
			{
				Config: testAccMdlStreamLiveInputUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mdl_stream_live_input.stream_live_input", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mdl_stream_live_input.stream_live_input", "name", "terraform-for-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_mdl_stream_live_input.stream_live_input",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMdlStreamLiveInput = `

resource "tencentcloud_mdl_stream_live_input" "stream_live_input" {
  name               = "terraform_test"
  type               = "RTP_PUSH"
  security_group_ids = [
    "6405DF9D000007DFB4EC"
  ]
}

`

const testAccMdlStreamLiveInputUpdate = `

resource "tencentcloud_mdl_stream_live_input" "stream_live_input" {
  name               = "terraform-for-test"
  type               = "RTP_PUSH"
  security_group_ids = [
    "6405DF9D000007DFB4EC"
  ]
}

`
