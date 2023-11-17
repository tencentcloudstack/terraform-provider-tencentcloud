package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudEmrCvmQuotaDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEmrCvmQuotaDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_emr_cvm_quota.cvm_quota")),
			},
		},
	})
}

const testAccEmrCvmQuotaDataSource = `

data "tencentcloud_emr_cvm_quota" "cvm_quota" {
  cluster_id = "emr-gmz8tdmv"
  zone_id = 800007
      }

`
