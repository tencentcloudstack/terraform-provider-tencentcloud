package emr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudEmrServiceNodeInfosDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccEmrServiceNodeInfosDataSource,
			Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_emr_service_node_infos.emr_service_node_infos")),
		}},
	})
}

const testAccEmrServiceNodeInfosDataSource = `
data "tencentcloud_emr_service_node_infos" "emr_service_node_infos" {
  instance_id = "emr-rzrochgp"
  offset = 1
  limit = 10
  search_text = ""
  conf_status = 2
  maintain_state_id = 2
  operator_state_id = 1
  health_state_id = "2"
  service_name = "YARN"
  node_type_name = "master"
  data_node_maintenance_id = 0
}
`
