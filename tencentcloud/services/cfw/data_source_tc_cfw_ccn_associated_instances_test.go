package cfw_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCfwCcnAssociatedInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCfwCcnAssociatedInstancesDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cfw_ccn_associated_instances.example"),
			),
		}},
	})
}

const testAccCfwCcnAssociatedInstancesDataSource = `
data "tencentcloud_cfw_ccn_associated_instances" "example" {
  ccn_id = "ccn-fkb9bo2v"
}
`
