package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmDisasterRecoverGroupQuotaDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmDisasterRecoverGroupQuotaDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_disaster_recover_group_quota.disaster_recover_group_quota")),
			},
		},
	})
}

const testAccCvmDisasterRecoverGroupQuotaDataSource = `

data "tencentcloud_cvm_disaster_recover_group_quota" "disaster_recover_group_quota" {
}

`
