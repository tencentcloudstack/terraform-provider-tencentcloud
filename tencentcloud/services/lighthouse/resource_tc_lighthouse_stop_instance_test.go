package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseStopInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
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

const testAccLighthouseStopInstance = tcacctest.DefaultLighthoustVariables + `
resource "tencentcloud_lighthouse_stop_instance" "stop_instance" {
  instance_id = var.lighthouse_id
}
`
