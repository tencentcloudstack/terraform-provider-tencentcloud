package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfReservedConcurrencyConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfReservedConcurrencyConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_scf_reserved_concurrency_config.reserved_concurrency_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_scf_reserved_concurrency_config.reserved_concurrency_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccScfReservedConcurrencyConfig = `

resource "tencentcloud_scf_reserved_concurrency_config" "reserved_concurrency_config" {
  function_name = "test_function"
  reserved_concurrency_mem = 128000
  namespace = "test_namespace"
}

`
