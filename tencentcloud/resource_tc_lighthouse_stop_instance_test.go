package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseStopInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseStopInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_stop_instance.stop_instance", "id")),
			},
			{
				Config: testAccLighthouseStartInstance,
			},
		},
	})
}

const testAccLighthouseStopInstance = `
resource "tencentcloud_lighthouse_stop_instance" "stop_instance" {
  instance_id = "lhins-hwe21u91"
}
`
