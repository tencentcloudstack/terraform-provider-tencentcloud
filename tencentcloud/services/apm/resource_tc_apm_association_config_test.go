package apm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudApmAssociationConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApmAssociation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_apm_association_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_association_config.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_association_config.example", "product_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_association_config.example", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_association_config.example", "peer_id"),
				),
			},
			{
				Config: testAccApmAssociationUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_apm_association_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_association_config.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_association_config.example", "product_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_association_config.example", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_apm_association_config.example", "peer_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_apm_association_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccApmAssociation = `
resource "tencentcloud_apm_association_config" "example" {
  instance_id  = tencentcloud_apm_instance.example.id
  product_name = "Prometheus"
  status       = 1
  peer_id      = "prom-kx3eqdby"
}
`

const testAccApmAssociationUpdate = `
resource "tencentcloud_apm_association_config" "example" {
  instance_id  = tencentcloud_apm_instance.example.id
  product_name = "Prometheus"
  status       = 2
  peer_id      = "prom-kx3eqdby"
}
`
