package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveDescribeLivePadTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveDescribeLivePadTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_live_describe_live_pad_templates.describe_live_pad_templates")),
			},
		},
	})
}

const testAccLiveDescribeLivePadTemplatesDataSource = `

data "tencentcloud_live_describe_live_pad_templates" "describe_live_pad_templates" {
  }

`
