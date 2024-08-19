package cdc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCdcDedicatedClusterResource_basic -v
func TestAccTencentCloudNeedFixCdcDedicatedClusterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdcDedicatedCluster,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster.example", "site_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster.example", "zone"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster.example", "description"),
				),
			},
			{
				Config: testAccCdcDedicatedClusterUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster.example", "site_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster.example", "zone"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdc_dedicated_cluster.example", "description"),
				),
			},
			{
				ResourceName:      "tencentcloud_cdc_dedicated_cluster.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdcDedicatedCluster = `
resource "tencentcloud_cdc_site" "example" {
  name         = "tf-example"
  country      = "China"
  province     = "Guangdong Province"
  city         = "Guangzhou"
  address_line = "Tencent Building"
  description  = "desc."
}

resource "tencentcloud_cdc_dedicated_cluster" "example" {
  site_id     = tencentcloud_cdc_site.example.id
  name        = "tf-example"
  zone        = "ap-guangzhou-6"
  description = "desc."
}
`

const testAccCdcDedicatedClusterUpdate = `
resource "tencentcloud_cdc_site" "example" {
  name         = "tf-example"
  country      = "China"
  province     = "Guangdong Province"
  city         = "Guangzhou"
  address_line = "Tencent Building"
  description  = "desc."
}

resource "tencentcloud_cdc_dedicated_cluster" "example" {
  site_id     = tencentcloud_cdc_site.example.id
  name        = "tf-example-update"
  zone        = "ap-guangzhou-3"
  description = "desc update."
}
`
