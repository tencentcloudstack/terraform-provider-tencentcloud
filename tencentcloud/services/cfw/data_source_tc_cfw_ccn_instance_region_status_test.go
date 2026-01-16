package cfw_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCfwCcnInstanceRegionStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCfwCcnInstanceRegionStatusDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cfw_ccn_instance_region_status.example"),
			),
		}},
	})
}

const testAccCfwCcnInstanceRegionStatusDataSource = `
data "tencentcloud_cfw_ccn_instance_region_status" "example" {
  ccn_id = "ccn-fkb9bo2v"
  instance_ids = [
    "vpc-axbsvrrg"
  ]
  routing_mode = 1
}
`
