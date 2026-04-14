package cls_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudClsNoticeContentsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsNoticeContentsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cls_notice_contents.example"),
				),
			},
		},
	})
}

func TestAccTencentCloudClsNoticeContentsDataSource_withFilters(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsNoticeContentsDataSourceWithFilters,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cls_notice_contents.example_with_filters"),
				),
			},
		},
	})
}

const testAccClsNoticeContentsDataSource = `
data "tencentcloud_cls_notice_contents" "example" {}
`

const testAccClsNoticeContentsDataSourceWithFilters = `
data "tencentcloud_cls_notice_contents" "example_with_filters" {
  filters {
    key    = "name"
    values = ["tf-example"]
  }
}
`
