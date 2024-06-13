package cvm_test

import (
	"testing"

	resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmDisasterRecoverGroupQuotaDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmDisasterRecoverGroupQuotaDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_disaster_recover_group_quota.disaster_recover_group_quota"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_disaster_recover_group_quota.disaster_recover_group_quota", "cvm_in_sw_group_quota", "20"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_disaster_recover_group_quota.disaster_recover_group_quota", "cvm_in_rack_group_quota", "30"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_disaster_recover_group_quota.disaster_recover_group_quota", "group_quota", "1000"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_disaster_recover_group_quota.disaster_recover_group_quota", "current_num"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_disaster_recover_group_quota.disaster_recover_group_quota", "cvm_in_host_group_quota", "50")),
			},
		},
	})
}

const testAccCvmDisasterRecoverGroupQuotaDataSource_BasicCreate = `

data "tencentcloud_cvm_disaster_recover_group_quota" "disaster_recover_group_quota" {
}

`
