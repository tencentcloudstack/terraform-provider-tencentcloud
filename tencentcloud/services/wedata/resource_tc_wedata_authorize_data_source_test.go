package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataAuthorizeDataSourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataAuthorizeDataSource,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_authorize_data_source.wedata_authorize_data_source", "id")),
		}, {
			ResourceName:      "tencentcloud_wedata_authorize_data_source.wedata_authorize_data_source",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccWedataAuthorizeDataSource = `

resource "tencentcloud_wedata_authorize_data_source" "wedata_authorize_data_source" {
}
`
