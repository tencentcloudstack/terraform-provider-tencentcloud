package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbRenewalPriceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbRenewalPriceDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_renewal_price.renewal_price")),
			},
		},
	})
}

const testAccDcdbRenewalPriceDataSource = CommonPresetDcdb + `

data "tencentcloud_dcdb_renewal_price" "renewal_price" {
	instance_id = local.dcdb_id
	period      = 1
	amount_unit = "pent"
}

`
