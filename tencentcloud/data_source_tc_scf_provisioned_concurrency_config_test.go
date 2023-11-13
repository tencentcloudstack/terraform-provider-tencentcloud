package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfProvisionedConcurrencyConfigDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfProvisionedConcurrencyConfigDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_provisioned_concurrency_config.provisioned_concurrency_config")),
			},
		},
	})
}

const testAccScfProvisionedConcurrencyConfigDataSource = `

data "tencentcloud_scf_provisioned_concurrency_config" "provisioned_concurrency_config" {
  function_name = ""
  namespace = ""
  qualifier = ""
    }

`
