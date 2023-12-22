package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixVpcSgSnapshotFileContentDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSgSnapshotFileContentDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_vpc_sg_snapshot_file_content.sg_snapshot_file_content")),
			},
		},
	})
}

const testAccVpcSgSnapshotFileContentDataSource = `

data "tencentcloud_vpc_sg_snapshot_file_content" "sg_snapshot_file_content" {
  snapshot_policy_id = "sspolicy-ebjofe71"
  snapshot_file_id   = "ssfile-017gepjxpr"
  security_group_id  = "sg-ntrgm89v"
}

`
