package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudElasticPublicIpv6Resource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticPublicIpv6,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_elastic_public_ipv6.elastic_public_ipv6", "id"),
					resource.TestCheckResourceAttr("tencentcloud_elastic_public_ipv6.elastic_public_ipv6", "address_name", "test"),
					resource.TestCheckResourceAttr("tencentcloud_elastic_public_ipv6.elastic_public_ipv6", "internet_max_bandwidth_out", "1"),
					resource.TestCheckResourceAttr("tencentcloud_elastic_public_ipv6.elastic_public_ipv6", "tags.test1key", "test1value"),
				),
			},
			{
				Config: testAccElasticPublicIpv6Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_elastic_public_ipv6.elastic_public_ipv6", "id"),
					resource.TestCheckResourceAttr("tencentcloud_elastic_public_ipv6.elastic_public_ipv6", "address_name", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_elastic_public_ipv6.elastic_public_ipv6", "internet_max_bandwidth_out", "2"),
					resource.TestCheckResourceAttr("tencentcloud_elastic_public_ipv6.elastic_public_ipv6", "tags.test2key", "test2value"),
				),
			},
			{
				ResourceName:      "tencentcloud_elastic_public_ipv6.elastic_public_ipv6",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccElasticPublicIpv6 = `
resource "tencentcloud_elastic_public_ipv6" "elastic_public_ipv6" {
    address_name = "test"
    internet_max_bandwidth_out = 1
    tags = {
        "test1key" = "test1value"
    }
}
`

const testAccElasticPublicIpv6Update = `
resource "tencentcloud_elastic_public_ipv6" "elastic_public_ipv6" {
    address_name = "test-update"
    internet_max_bandwidth_out = 2
    tags = {
        "test2key" = "test2value"
    }
}
`
