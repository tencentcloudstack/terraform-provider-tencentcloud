package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMdlStreamliveInputResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMdlStreamliveInput,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mdl_streamlive_input.streamlive_input", "id")),
			},
			{
				ResourceName:      "tencentcloud_mdl_streamlive_input.streamlive_input",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMdlStreamliveInput = `

resource "tencentcloud_mdl_streamlive_input" "streamlive_input" {
  name = ""
  type = ""
  security_group_ids = 
  input_settings {
		app_name = ""
		stream_name = ""
		source_url = ""
		input_address = ""
		source_type = ""
		delay_time = 
		input_domain = ""
		user_name = ""
		password = ""

  }
}

`
