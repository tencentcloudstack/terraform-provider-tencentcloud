package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmAccountQuotaDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmAccountQuotaDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_account_quota.quota"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_account_quota.quota", "app_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_account_quota.quota", "account_quota_overview.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_account_quota.quota", "account_quota_overview.0.region"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_account_quota.quota", "account_quota_overview.0.account_quota.#"),
				),
			},
		},
	})
}

func TestAccTencentCloudCvmAccountQuotaDataSource_filterByZone(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmAccountQuotaDataSource_filterByZone,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_account_quota.quota_zone"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_account_quota.quota_zone", "app_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_account_quota.quota_zone", "account_quota_overview.#"),
				),
			},
		},
	})
}

func TestAccTencentCloudCvmAccountQuotaDataSource_filterByQuotaType(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmAccountQuotaDataSource_filterByQuotaType,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_account_quota.quota_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_account_quota.quota_type", "app_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_account_quota.quota_type", "account_quota_overview.#"),
				),
			},
		},
	})
}

const testAccCvmAccountQuotaDataSource_basic = `
data "tencentcloud_cvm_account_quota" "quota" {}
`

const testAccCvmAccountQuotaDataSource_filterByZone = `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cvm"
}

data "tencentcloud_cvm_account_quota" "quota_zone" {
  zone = [data.tencentcloud_availability_zones_by_product.zones.zones.0.name]
}
`

const testAccCvmAccountQuotaDataSource_filterByQuotaType = `
data "tencentcloud_cvm_account_quota" "quota_type" {
  quota_type = "PostPaidQuotaSet"
}
`
