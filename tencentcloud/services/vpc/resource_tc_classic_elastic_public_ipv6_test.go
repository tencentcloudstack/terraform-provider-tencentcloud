package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudClassicElasticPublicIpv6Resource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClassicElasticPublicIpv6,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_classic_elastic_public_ipv6.classic_elastic_public_ipv6", "id"),
					resource.TestCheckResourceAttr("tencentcloud_classic_elastic_public_ipv6.classic_elastic_public_ipv6", "internet_max_bandwidth_out", "1"),
					resource.TestCheckResourceAttr("tencentcloud_classic_elastic_public_ipv6.classic_elastic_public_ipv6", "internet_charge_type", "TRAFFIC_POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr("tencentcloud_classic_elastic_public_ipv6.classic_elastic_public_ipv6", "tags.test1key", "test1value"),
				),
			},
			{
				Config: testAccClassicElasticPublicIpv6Update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_classic_elastic_public_ipv6.classic_elastic_public_ipv6", "id"),
					resource.TestCheckResourceAttr("tencentcloud_classic_elastic_public_ipv6.classic_elastic_public_ipv6", "internet_max_bandwidth_out", "2"),
					resource.TestCheckResourceAttr("tencentcloud_classic_elastic_public_ipv6.classic_elastic_public_ipv6", "internet_charge_type", "TRAFFIC_POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttr("tencentcloud_classic_elastic_public_ipv6.classic_elastic_public_ipv6", "tags.test2key", "test2value"),
				),
			},
			{
				ResourceName:      "tencentcloud_classic_elastic_public_ipv6.classic_elastic_public_ipv6",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

const testAccClassicElasticPublicIpv6 = `

resource "tencentcloud_classic_elastic_public_ipv6" "classic_elastic_public_ipv6" {
  ip6_address                = "2402:4e00:101d:8300:0:9dbc:f45a:2c4f"
  tags = {
    "test1key" = "test1value"
  }
}
`

const testAccClassicElasticPublicIpv6Update = `

resource "tencentcloud_classic_elastic_public_ipv6" "classic_elastic_public_ipv6" {
  ip6_address                = "2402:4e00:101d:8300:0:9dbc:f45a:2c4f"
  internet_max_bandwidth_out = 2
  tags = {
    "test2key" = "test2value"
  }
}
`
