package mongodb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMongodbInstanceParamsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { tcacctest.AccStepPreConfigSetTempAKSK(t, tcacctest.ACCOUNT_TYPE_COMMON) },
				Config:    testAccMongodbInstanceParamsDataSource,
				Check:     resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mongodb_instance_params.instance_params")),
			},
		},
	})
}

const testAccMongodbInstanceParamsDataSource = `

data "tencentcloud_mongodb_instance_params" "instance_params" {
  instance_id = "cmgo-gwqk8669"
}

`
