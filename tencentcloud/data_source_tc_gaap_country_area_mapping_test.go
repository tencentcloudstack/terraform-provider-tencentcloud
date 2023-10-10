package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapCountryAreaMappingDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapCountryAreaMappingDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_country_area_mapping.country_area_mapping"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_country_area_mapping.country_area_mapping", "country_area_mapping_list.#"),
				),
			},
		},
	})
}

const testAccGaapCountryAreaMappingDataSource = `

data "tencentcloud_gaap_country_area_mapping" "country_area_mapping" {
  }

`
