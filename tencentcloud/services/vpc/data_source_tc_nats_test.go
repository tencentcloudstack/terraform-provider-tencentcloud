package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNatsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudNatsDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_nats.multi_nat"),
					resource.TestCheckResourceAttr("data.tencentcloud_nats.multi_nat", "nats.#", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_nats.multi_nat", "nats.0.name", "terraform_test_nats"),
					resource.TestCheckResourceAttr("data.tencentcloud_nats.multi_nat", "nats.1.bandwidth", "500"),
				),
			},
			{
				Config: testAccTencentCloudNatsDataSourceConfig_withVerboseLevel,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_nats.verbose_nat"),
					resource.TestCheckResourceAttr("data.tencentcloud_nats.verbose_nat", "verbose_level", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudNatsDataSourceConfig_basic = `
resource "tencentcloud_vpc" "main" {
  name       = "terraform_test_nats"
  cidr_block = "10.6.0.0/16"
}

resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}

resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}

resource "tencentcloud_nat_gateway" "dev_nat" {
  vpc_id         = tencentcloud_vpc.main.id
  name           = "terraform_test_nats"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
  ]
}

resource "tencentcloud_nat_gateway" "test_nat" {
  vpc_id         = tencentcloud_vpc.main.id
  name           = "terraform_test_nats"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
}

data "tencentcloud_nats" "multi_nat" {
  state          = 0
  name           = tencentcloud_nat_gateway.dev_nat.name
  vpc_id         = tencentcloud_vpc.main.id
  max_concurrent = tencentcloud_nat_gateway.test_nat.max_concurrent
  bandwidth      = tencentcloud_nat_gateway.test_nat.bandwidth
}
`

const testAccTencentCloudNatsDataSourceConfig_withVerboseLevel = `
resource "tencentcloud_vpc" "main" {
  name       = "terraform_test_nats_verbose"
  cidr_block = "10.7.0.0/16"
}

resource "tencentcloud_eip" "eip_test_verbose" {
  name = "terraform_test_verbose"
}

resource "tencentcloud_nat_gateway" "test_nat_verbose" {
  vpc_id         = tencentcloud_vpc.main.id
  name           = "terraform_test_nats_verbose"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_test_verbose.public_ip,
  ]
}

data "tencentcloud_nats" "verbose_nat" {
  vpc_id        = tencentcloud_vpc.main.id
  verbose_level = 1
}
`
