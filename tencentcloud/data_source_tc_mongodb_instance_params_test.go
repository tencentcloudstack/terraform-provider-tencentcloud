package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMongodbInstanceParamsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_COMMON) },
				Config:    testAccMongodbInstanceParamsDataSource,
				Check:     resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mongodb_instance_params.instance_params")),
			},
		},
	})
}

const testAccMongodbInstanceParamsDataSource = `

data "tencentcloud_mongodb_instance_params" "instance_params" {
  instance_id = "cmgo-gwqk8669"
}

`
