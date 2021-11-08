package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudAvailabilityZonesByProductDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudAvailabilityZonesByProductByProductDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_availability_zones_by_product.all"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_zones_by_product.all", "zones.#"),
				),
			},
			{
				Config: testAccTencentCloudAvailabilityZonesByProductByProductDataSourceConfigFilterWithName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_availability_zones_by_product.filter"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_zones_by_product.filter", "zones.#"),
				),
			},
			{
				Config: testAccTencentCloudAvailabilityZonesByProductDataSourceConfigIncludeUnavailable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_availability_zones_by_product.unavailable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_zones_by_product.unavailable", "zones.#"),
				),
			},
			{
				Config: testAccTencentCloudAvailabilityZonesByProductDataSourceConfigIncludeUnavailable,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_availability_zones_by_product.unavailable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_availability_zones_by_product.unavailable", "zones.#"),
				),
			},
		},
	})
}

const testAccTencentCloudAvailabilityZonesByProductByProductDataSourceConfigBasic = `
data "tencentcloud_availability_zones_by_product" "basic" {
	product = var.cvm_product
}
`

const testAccTencentCloudAvailabilityZonesByProductByProductDataSourceConfigFilterWithName = defaultVpcVariable + `
data "tencentcloud_availability_zones_by_product" "filter" {
	product = var.cvm_product
  	name = var.availability_zone
}
`

const testAccTencentCloudAvailabilityZonesByProductDataSourceConfigIncludeUnavailable = `
data "tencentcloud_availability_zones_by_product" "unavailable" {
	product = var.cvm_product
  	include_unavailable = true
}
`
