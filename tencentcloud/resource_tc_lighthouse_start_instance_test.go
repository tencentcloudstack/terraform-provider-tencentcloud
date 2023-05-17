package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseStartInstanceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseStopInstance,
			},
			{
				Config: testAccLighthouseStartInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_start_instance.start_instance", "id")),
			},
		},
	})
}

const testAccLighthouseStartInstance = `

resource "tencentcloud_lighthouse_start_instance" "start_instance" {
  instance_id = "lhins-hwe21u91"
}

`
