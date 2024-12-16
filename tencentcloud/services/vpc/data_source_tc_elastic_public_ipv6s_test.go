package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudElasticPublicIpv6sDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccElasticPublicIpv6sDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_elastic_public_ipv6s.elastic_public_ipv6s"),
					resource.TestCheckResourceAttr("data.tencentcloud_elastic_public_ipv6s.elastic_public_ipv6s", "address_set.#", "1"),
				),
			},
		},
	})
}

const testAccElasticPublicIpv6sDataSource = testAccElasticPublicIpv6 + `

data "tencentcloud_elastic_public_ipv6s" "elastic_public_ipv6s" {
  ipv6_address_ids = [tencentcloud_elastic_public_ipv6.elastic_public_ipv6.id]
}
`
