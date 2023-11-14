package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcCrossBorderComplianceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcCrossBorderComplianceDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_cross_border_compliance.cross_border_compliance")),
			},
		},
	})
}

const testAccVpcCrossBorderComplianceDataSource = `

data "tencentcloud_vpc_cross_border_compliance" "cross_border_compliance" {
  service_provider = "UNICOM"
  compliance_id = 10002
  company = "腾讯科技（广州）有限公司"
  uniform_social_credit_code = "91440101327598294H"
  legal_person = "张颖"
  issuing_authority = "广州市海珠区市场监督管理局"
  business_address = "广州市海珠区新港中路397号自编72号(商业街F5-1)"
  post_code = 510320
  manager = "李四"
  manager_id = "360732199007108888"
  manager_address = "广州市海珠区新港中路8888号"
  manager_telephone = "020-81167888"
  email = "test@tencent.com"
  service_start_date = "2020-07-29"
  service_end_date = "2021-07-29"
  state = "APPROVED"
  offset = 1
  limit = 2
}

`
