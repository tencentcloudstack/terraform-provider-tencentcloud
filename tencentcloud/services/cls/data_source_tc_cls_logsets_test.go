package cls_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudClsLogsetsDataSource_basic -v
func TestAccTencentCloudClsLogsetsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsLogsetsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cls_logsets.logsets"),
				),
			},
		},
	})
}

const testAccClsLogsetsDataSource = `
data "tencentcloud_cls_logsets" "logsets" {
  filters {
    key    = "logsetId"
    values = ["50d499a8-c4c0-4442-aa04-e8aa8a02437d"]
  }
}
`
