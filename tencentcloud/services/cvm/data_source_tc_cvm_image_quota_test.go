package cvm_test

import (
	"testing"

	resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmImageQuotaDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmImageQuotaDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_image_quota.image_quota"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_image_quota.image_quota", "image_num_quota", "500")),
			},
		},
	})
}

const testAccCvmImageQuotaDataSource_BasicCreate = `

data "tencentcloud_cvm_image_quota" "image_quota" {
}

`
