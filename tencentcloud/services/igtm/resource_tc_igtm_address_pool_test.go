package igtm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIgtmAddressPoolResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccIgtmAddressPool,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_igtm_address_pool.igtm_address_pool", "id")),
		}, {
			ResourceName:      "tencentcloud_igtm_address_pool.igtm_address_pool",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccIgtmAddressPool = `

resource "tencentcloud_igtm_address_pool" "igtm_address_pool" {
  address_set = {
  }
}
`
