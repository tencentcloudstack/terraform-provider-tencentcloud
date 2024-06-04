package cvm_test

import (
	"testing"

	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmHpcClusterResource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmHpcClusterResource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_hpc_cluster.hpc_cluster", "id"), resource.TestCheckResourceAttr("tencentcloud_cvm_hpc_cluster.hpc_cluster", "name", "terraform-test"), resource.TestCheckResourceAttr("tencentcloud_cvm_hpc_cluster.hpc_cluster", "remark", "create for test")),
			},
			{
				Config: testAccCvmHpcClusterResource_BasicChange1,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("tencentcloud_cvm_hpc_cluster.hpc_cluster", "remark", "create for e2e test"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_hpc_cluster.hpc_cluster", "id"), resource.TestCheckResourceAttr("tencentcloud_cvm_hpc_cluster.hpc_cluster", "name", "terraform-test1")),
			},
			{
				ResourceName:      "tencentcloud_cvm_hpc_cluster.hpc_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmHpcClusterResource_BasicCreate = `

resource "tencentcloud_cvm_hpc_cluster" "hpc_cluster" {
    zone = "ap-guangzhou-7"
    name = "terraform-test"
    remark = "create for test"
}

`
const testAccCvmHpcClusterResource_BasicChange1 = `

resource "tencentcloud_cvm_hpc_cluster" "hpc_cluster" {
    zone = "ap-guangzhou-7"
    name = "terraform-test1"
    remark = "create for e2e test"
}

`

func TestAccTencentCloudCvmHpcClusterResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config:    testAccCvmHpcCluster,
				PreConfig: func() { tcacctest.AccStepSetRegion(t, "ap-guangzhou") },
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_hpc_cluster.hpc_cluster", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_hpc_cluster.hpc_cluster", "name", "terraform-test"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_hpc_cluster.hpc_cluster", "remark", "create for test"),
				),
			},
			{
				Config:    testAccCvmHpcClusterUpdate,
				PreConfig: func() { tcacctest.AccStepSetRegion(t, "ap-guangzhou") },
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_hpc_cluster.hpc_cluster", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_hpc_cluster.hpc_cluster", "name", "terraform-test1"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_hpc_cluster.hpc_cluster", "remark", "create for e2e test"),
				),
			},
			{
				ResourceName:      "tencentcloud_cvm_hpc_cluster.hpc_cluster",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmHpcClusterBasis = `
variable "availability_zone" {
  default = "ap-guangzhou-7"
}
`

const testAccCvmHpcCluster = testAccCvmHpcClusterBasis + `
resource "tencentcloud_cvm_hpc_cluster" "hpc_cluster" {
  zone   = var.availability_zone
  name   = "terraform-test"
  remark = "create for test"
}
`

const testAccCvmHpcClusterUpdate = testAccCvmHpcClusterBasis + `
resource "tencentcloud_cvm_hpc_cluster" "hpc_cluster" {
  zone   = var.availability_zone
  name   = "terraform-test1"
  remark = "create for e2e test"
}
`
