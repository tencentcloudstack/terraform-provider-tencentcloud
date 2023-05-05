package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTatInvokerRecordsDataSource_basic -v
func TestAccTencentCloudNeedFixTatInvokerRecordsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTatInvokerRecordsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tat_invoker_records.invoker_records")),
			},
		},
	})
}

const testAccTatInvokerRecordsDataSource = `

data "tencentcloud_tat_invoker_records" "invoker_records" {
  invoker_ids = ""
}

`
