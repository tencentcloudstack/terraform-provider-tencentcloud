package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapCountryAreaMappingDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapCountryAreaMappingDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_country_area_mapping.country_area_mapping"),
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
