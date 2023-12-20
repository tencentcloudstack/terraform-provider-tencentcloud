package dcdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbRenewalPriceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbRenewalPriceDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_renewal_price.renewal_price")),
			},
		},
	})
}

const testAccDcdbRenewalPriceDataSource = tcacctest.CommonPresetDcdb + `

data "tencentcloud_dcdb_renewal_price" "renewal_price" {
	instance_id = local.dcdb_id
	period      = 1
	amount_unit = "pent"
}

`
