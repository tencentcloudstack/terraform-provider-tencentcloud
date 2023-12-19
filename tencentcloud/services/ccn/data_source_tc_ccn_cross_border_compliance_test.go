package ccn_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCcnCrossBorderComplianceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCcnCrossBorderComplianceDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_ccn_cross_border_compliance.cross_border_compliance")),
			},
		},
	})
}

const testAccCcnCrossBorderComplianceDataSource = `

data "tencentcloud_ccn_cross_border_compliance" "cross_border_compliance" {
  service_provider = "UNICOM"
  compliance_id = 10002
  email = "test@tencent.com"
  service_start_date = "2020-07-29"
  service_end_date = "2021-07-29"
  state = "APPROVED"
}

`
