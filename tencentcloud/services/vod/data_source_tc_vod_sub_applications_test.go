package vod_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudVodSubApplications_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSubApplicationsBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vod_sub_applications.all"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_sub_applications.all", "sub_application_info_set.#"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudVodSubApplications_NameFilter(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSubApplicationsNameFilter,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vod_sub_applications.by_name"),
					resource.TestCheckResourceAttr("data.tencentcloud_vod_sub_applications.by_name", "sub_application_info_set.0.name", "terraform-test-app"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_sub_applications.by_name", "sub_application_info_set.0.sub_app_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_sub_applications.by_name", "sub_application_info_set.0.sub_app_id_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_sub_applications.by_name", "sub_application_info_set.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_vod_sub_applications.by_name", "sub_application_info_set.0.create_time"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudVodSubApplications_Pagination(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSubApplicationsPagination,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vod_sub_applications.paginated"),
				),
			},
		},
	})
}

const testAccDataSourceVodSubApplication = `
resource "tencentcloud_vod_sub_application" "test" {
  name        = "terraform-test-app"
  status      = "On"
  description = "Test application for terraform"
}
`

const testAccVodSubApplicationsBasic = `
data "tencentcloud_vod_sub_applications" "all" {}
`

const testAccVodSubApplicationsNameFilter = testAccDataSourceVodSubApplication + `
data "tencentcloud_vod_sub_applications" "by_name" {
  name = tencentcloud_vod_sub_application.test.name
}
`

const testAccVodSubApplicationsPagination = `
data "tencentcloud_vod_sub_applications" "paginated" {
  offset = 0
  limit  = 10
}
`
