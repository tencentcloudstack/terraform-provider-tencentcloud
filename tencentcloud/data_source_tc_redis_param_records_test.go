package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRedisParamRecordsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisParamRecordsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_redis_param_records.param_records")),
			},
		},
	})
}

const testAccRedisParamRecordsDataSource = `

data "tencentcloud_redis_param_records" "param_records" {
  instance_id = "crs-c1nl9rpv"
  limit = &lt;nil&gt;
  offset = &lt;nil&gt;
  }

`
