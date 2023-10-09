package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudEbPlateformEventTemplateDataSource_basic -v
func TestAccTencentCloudEbPlateformEventTemplateDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-chongqing")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbPlateformEventTemplateDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_eb_plateform_event_template.plateform_event_template"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_plateform_event_template.plateform_event_template", "event_template"),
				),
			},
		},
	})
}

const testAccEbPlateformEventTemplateDataSource = `

data "tencentcloud_eb_plateform_event_template" "plateform_event_template" {
  event_type = "eb_platform_test:TEST:ALL"
}

`
