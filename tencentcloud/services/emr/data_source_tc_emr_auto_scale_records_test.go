package emr_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudEmrAutoScaleRecordsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEmrAutoScaleRecordsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_emr_auto_scale_records.auto_scale_records")),
			},
		},
	})
}

const testAccEmrAutoScaleRecordsDataSource = `

data "tencentcloud_emr_auto_scale_records" "auto_scale_records" {
  instance_id = "emr-8j38bip0"
  filters {
    key   = "StartTime"
    value = "2006-01-02 15:04:05"
  }
}

`
