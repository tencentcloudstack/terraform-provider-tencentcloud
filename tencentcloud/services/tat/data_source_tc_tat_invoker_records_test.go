package tat_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTatInvokerRecordsDataSource_basic -v
func TestAccTencentCloudNeedFixTatInvokerRecordsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTatInvokerRecordsDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tat_invoker_records.invoker_records")),
			},
		},
	})
}

const testAccTatInvokerRecordsDataSource = `

data "tencentcloud_tat_invoker_records" "invoker_records" {
  invoker_ids = ""
}

`
