package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfReservedConcurrencyConfigDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfReservedConcurrencyConfigDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_reserved_concurrency_config.reserved_concurrency_config")),
			},
		},
	})
}

const testAccScfReservedConcurrencyConfigDataSource = `

data "tencentcloud_scf_reserved_concurrency_config" "reserved_concurrency_config" {
  function_name = ""
  namespace = ""
  }

`
