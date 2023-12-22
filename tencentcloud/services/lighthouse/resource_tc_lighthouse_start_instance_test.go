package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseStartInstanceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
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

const testAccLighthouseStartInstance = tcacctest.DefaultLighthoustVariables + `

resource "tencentcloud_lighthouse_start_instance" "start_instance" {
  instance_id = var.lighthouse_id
}

`
