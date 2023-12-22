package tpulsar_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTdmqProfessionalClusterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqProfessionalCluster,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_professional_cluster.professional_cluster", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_professional_cluster.professional_cluster", "cluster_name", "single_zone_cluster"),
				),
			},
			{
				Config: testAccTdmqProfessionalClusterUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdmq_professional_cluster.professional_cluster", "id"),
					resource.TestCheckResourceAttr("tencentcloud_tdmq_professional_cluster.professional_cluster", "cluster_name", "single_zone_cluster_test"),
				),
			},
			{
				ResourceName:            "tencentcloud_tdmq_professional_cluster.professional_cluster",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"time_span", "auto_voucher"},
			},
		},
	})
}

const testAccTdmqProfessionalCluster = `

resource "tencentcloud_tdmq_professional_cluster" "professional_cluster" {
  auto_renew_flag = 1
  cluster_name    = "single_zone_cluster"
  product_name    = "PULSAR.P1.MINI2"
  storage_size    = 600
  tags            = {
    "createby" = "terrafrom"
  }
  zone_ids = [
    100004,
  ]

  vpc {
    subnet_id = "subnet-2txtpql8"
    vpc_id    = "vpc-njbzmzyd"
  }
}

`

const testAccTdmqProfessionalClusterUpdate = `

resource "tencentcloud_tdmq_professional_cluster" "professional_cluster" {
  auto_renew_flag = 1
  cluster_name    = "single_zone_cluster_test"
  product_name    = "PULSAR.P1.MINI2"
  storage_size    = 600
  tags            = {
    "createby" = "terrafrom"
  }
  zone_ids = [
    100004,
  ]

  vpc {
    subnet_id = "subnet-2txtpql8"
    vpc_id    = "vpc-njbzmzyd"
  }
}

`
