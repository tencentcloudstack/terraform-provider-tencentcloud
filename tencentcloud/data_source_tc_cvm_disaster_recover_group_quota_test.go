package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCvmDisasterRecoverGroupQuotaDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmDisasterRecoverGroupQuotaDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_disaster_recover_group_quota.disaster_recover_group_quota")),
			},
		},
	})
}

const testAccCvmDisasterRecoverGroupQuotaDataSource = `

data "tencentcloud_cvm_disaster_recover_group_quota" "disaster_recover_group_quota" {
}

`
