package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudEbPlateformEventTemplateDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbPlateformEventTemplateDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_eb_plateform_event_template.plateform_event_template")),
			},
		},
	})
}

const testAccEbPlateformEventTemplateDataSource = `

data "tencentcloud_eb_plateform_event_template" "plateform_event_template" {
  event_type = ""
  }

`
