package mongodb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMongodbInstanceCurrentOpDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config:    testAccMongodbInstanceCurrentOpDataSource,
				Check:     resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mongodb_instance_current_op.instance_current_op")),
			},
		},
	})
}

const testAccMongodbInstanceCurrentOpDataSource = `

data "tencentcloud_mongodb_instance_current_op" "instance_current_op" {
  instance_id = "cmgo-gwqk8669"
  op = "command"
  order_by_type = "desc"
}

`
