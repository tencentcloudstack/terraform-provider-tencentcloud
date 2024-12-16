package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudClassicElasticPublicIpv6sDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccClassicElasticPublicIpv6sDataSource,
			Check: resource.ComposeTestCheckFunc(
				tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_classic_elastic_public_ipv6s.classic_elastic_public_ipv6s"),
				resource.TestCheckResourceAttr("data.tencentcloud_classic_elastic_public_ipv6s.classic_elastic_public_ipv6s", "address_set.#", "1"),
			),
		}},
	})
}

const testAccClassicElasticPublicIpv6sDataSource = testAccClassicElasticPublicIpv6 + `

data "tencentcloud_classic_elastic_public_ipv6s" "classic_elastic_public_ipv6s" {
  ip6_address_ids = [tencentcloud_classic_elastic_public_ipv6.classic_elastic_public_ipv6.id]
}
`
