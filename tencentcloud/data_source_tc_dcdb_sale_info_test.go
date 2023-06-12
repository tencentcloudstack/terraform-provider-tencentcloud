package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbSaleInfoDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbSaleInfoDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_sale_info.sale_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.region_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.region_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.zone_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.zone_list.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.zone_list.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.zone_list.0.zone_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.zone_list.0.on_sale"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.available_choice.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.available_choice.0.master_zone.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.available_choice.0.master_zone.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.available_choice.0.master_zone.0.zone_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.available_choice.0.slave_zones.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.available_choice.0.slave_zones.0.zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_sale_info.sale_info", "region_list.0.available_choice.0.slave_zones.0.zone_id"),
				),
			},
		},
	})
}

const testAccDcdbSaleInfoDataSource = `

data "tencentcloud_dcdb_sale_info" "sale_info" {}

`
