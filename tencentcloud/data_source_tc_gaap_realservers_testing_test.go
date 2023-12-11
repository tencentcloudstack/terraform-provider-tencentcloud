package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudTestingGaapRealservers_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudTestingGaapRealserversBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_realservers.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudTestingGaapRealservers_domain(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudTestingGaapRealserversDomain,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_realservers.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudTestingGaapRealservers_name(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudTestingGaapRealserversName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_realservers.foo"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudTestingGaapRealserversBasic = `

data tencentcloud_gaap_realservers "foo" {
  ip = "1.1.1.10"
}
`

const TestAccDataSourceTencentCloudTestingGaapRealserversDomain = `

data tencentcloud_gaap_realservers "foo" {
  domain = "www.qq.com"
}
`

const TestAccDataSourceTencentCloudTestingGaapRealserversName = `

data tencentcloud_gaap_realservers "foo" {
  name = "www.tencent.com"
}
`
