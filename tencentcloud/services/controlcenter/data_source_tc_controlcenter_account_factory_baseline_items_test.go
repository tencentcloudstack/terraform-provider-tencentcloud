package controlcenter_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudControlcenterAccountFactoryBaselineItemsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccControlcenterAccountFactoryBaselineItemsDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_controlcenter_account_factory_baseline_items.example")),
		}},
	})
}

const testAccControlcenterAccountFactoryBaselineItemsDataSource = `
data "tencentcloud_controlcenter_account_factory_baseline_items" "example" {}
`
