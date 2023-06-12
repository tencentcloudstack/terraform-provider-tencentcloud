package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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

const testAccDcdbRenewalPriceDataSource = `

data "tencentcloud_dcdb_renewal_price" "renewal_price" {
  instance_id = ""
  period = 
  amount_unit = ""
    }

`
