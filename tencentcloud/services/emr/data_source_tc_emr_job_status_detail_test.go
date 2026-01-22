package emr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudEmrJobStatusDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccEmrJobStatusDetailDataSource,
			Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_emr_job_status_detail.emr_job_status_detail")),
		}},
	})
}

const testAccEmrJobStatusDetailDataSource = `
data tencentcloud_emr_job_status_detail "emr_job_status_detail" {
	instance_id = "emr-byhnjsb3"
	flow_param {
		f_key = "FlowId"
		f_value = "1921228"
	}
}
`
