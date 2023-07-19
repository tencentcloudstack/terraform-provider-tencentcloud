package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverCompleteExpansionResource_basic -v
// go test -v -run TestAccTencentCloudSqlserverCompleteExpansionResource_basic -timeout=0 ./tencentcloud/
func TestAccTencentCloudSqlserverCompleteExpansionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewSqlserverInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists(testSqlserverInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "id"),
				),
			},
			{
				Config: testAccUpdateNewSqlserverInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists(testSqlserverInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testSqlserverInstanceResourceKey, "id"),
				),
			},
		},
	})
}

const testAccNewSqlserverInstance string = testAccSqlserverBasicInstanceNetwork + `
resource "tencentcloud_sqlserver_instance" "test" {
  name                   = "tf_sqlserver_instance"
  availability_zone      = "ap-guangzhou-7"
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = "vpc-1yg5ua6l"
  subnet_id              = "subnet-h7av55g8"
  security_groups        = ["sg-mayqdlt1"]
  project_id             = 0
  memory                 = 2
  storage                = 20
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  tags                   = {
    "test" = "test"
  }
}
`

const testAccUpdateNewSqlserverInstance string = testAccSqlserverBasicInstanceNetwork + `
resource "tencentcloud_sqlserver_instance" "test" {
  name                   = "tf_sqlserver_instance"
  availability_zone      = "ap-guangzhou-7"
  charge_type            = "POSTPAID_BY_HOUR"
  vpc_id                 = "vpc-1yg5ua6l"
  subnet_id              = "subnet-h7av55g8"
  security_groups        = ["sg-mayqdlt1"]
  project_id             = 0
  memory                 = 2
  storage                = 40
  maintenance_week_set   = [1, 2, 3]
  maintenance_start_time = "09:00"
  maintenance_time_span  = 3
  wait_switch            = 1
  tags                   = {
    "test" = "test"
  }
}

resource "tencentcloud_sqlserver_complete_expansion" "complete_expansion" {
  instance_id = tencentcloud_sqlserver_instance.test.id
}
`
