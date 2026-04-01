package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccTencentCloudCvmResourcePoolPacksResource_basic tests basic create/read/delete lifecycle
func TestAccTencentCloudCvmResourcePoolPacksResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config:    testAccCvmResourcePoolPacks,
				PreConfig: func() { tcacctest.AccStepSetRegion(t, "ap-guangzhou") },
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_resource_pool_packs.pool_packs", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_resource_pool_packs.pool_packs", "dedicated_resource_pack_id"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_resource_pool_packs.pool_packs", "zone", "ap-guangzhou-7"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_resource_pool_packs.pool_packs", "instance_type", "SA9.96XLARGE1152"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_resource_pool_packs.pool_packs", "instance_count", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_resource_pool_packs.pool_packs", "period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_resource_pool_packs.pool_packs", "resource_pool_pack_type", "EXCLUSIVE"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_resource_pool_packs.pool_packs", "auto_placement", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_resource_pool_packs.pool_packs", "renew_flag", "NOTIFY_AND_MANUAL_RENEW"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_resource_pool_packs.pool_packs", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_resource_pool_packs.pool_packs", "instance_family"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_resource_pool_packs.pool_packs", "start_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_resource_pool_packs.pool_packs", "end_time"),
				),
			},
			{
				ResourceName:      "tencentcloud_cvm_resource_pool_packs.pool_packs",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccTencentCloudCvmResourcePoolPacksResource_forceNew verifies ForceNew behavior
func TestAccTencentCloudCvmResourcePoolPacksResource_forceNew(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config:    testAccCvmResourcePoolPacks,
				PreConfig: func() { tcacctest.AccStepSetRegion(t, "ap-guangzhou") },
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_resource_pool_packs.pool_packs", "id"),
				),
			},
			{
				Config:    testAccCvmResourcePoolPacksForceNew,
				PreConfig: func() { tcacctest.AccStepSetRegion(t, "ap-guangzhou") },
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_resource_pool_packs.pool_packs", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_resource_pool_packs.pool_packs", "period", "2"),
				),
			},
		},
	})
}

const testAccCvmResourcePoolPacksBasis = `
variable "availability_zone" {
  default = "ap-guangzhou-7"
}
`

const testAccCvmResourcePoolPacks = testAccCvmResourcePoolPacksBasis + `
resource "tencentcloud_cvm_resource_pool_packs" "pool_packs" {
  zone                              = var.availability_zone
  instance_type                     = "SA9.96XLARGE1152"
  instance_count                    = 1
  period                            = 1
  resource_pool_pack_type           = "EXCLUSIVE"
  auto_placement                    = true
  dedicated_resource_pool_pack_name = "terraform-test-pool"
  renew_flag                        = "NOTIFY_AND_MANUAL_RENEW"
}
`

const testAccCvmResourcePoolPacksForceNew = testAccCvmResourcePoolPacksBasis + `
resource "tencentcloud_cvm_resource_pool_packs" "pool_packs" {
  zone                              = var.availability_zone
  instance_type                     = "SA9.96XLARGE1152"
  instance_count                    = 1
  period                            = 2
  resource_pool_pack_type           = "EXCLUSIVE"
  auto_placement                    = true
  dedicated_resource_pool_pack_name = "terraform-test-pool-new"
  renew_flag                        = "NOTIFY_AND_MANUAL_RENEW"
}
`
