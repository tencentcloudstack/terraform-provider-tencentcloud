package igtm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIgtmStrategyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccIgtmStrategy,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_igtm_strategy.igtm_strategy", "id")),
		}, {
			ResourceName:      "tencentcloud_igtm_strategy.igtm_strategy",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccIgtmStrategy = `

resource "tencentcloud_igtm_strategy" "igtm_strategy" {
  source = {
  }
  main_address_pool_set = {
    address_pools = {
    }
  }
  fallback_address_pool_set = {
    address_pools = {
    }
  }
}
`
