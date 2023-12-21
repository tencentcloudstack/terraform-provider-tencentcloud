package mariadb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbDedicatedclusterDbInstance_basic -v
func TestAccTencentCloudNeedFixMariadbDedicatedclusterDbInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbDedicatedclusterDbInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mariadb_dedicatedcluster_db_instance.dedicatedcluster_db_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_mariadb_dedicatedcluster_db_instance.dedicatedclusterDbInstance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbDedicatedclusterDbInstance = `
resource "tencentcloud_mariadb_dedicatedcluster_db_instance" "dedicatedcluster_db_instance" {
  goods_num          = 1
  memory             = 2
  storage            = 10
  cluster_id         = "dbdc-24odnuhr"
  vpc_id             = "vpc-ii1jfbhl"
  subnet_id          = "subnet-3ku415by"
  db_version_id      = "8.0"
  instance_name      = "cluster-mariadb-test-1"
}
`
