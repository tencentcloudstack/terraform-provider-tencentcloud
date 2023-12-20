package cdwch_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudClickhouseInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseInstanceBasic,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_instance.cdwch_instance", "id")),
			},
			{
				ResourceName:            "tencentcloud_clickhouse_instance.cdwch_instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"time_span"},
			},
		},
	})
}

func TestAccTencentCloudClickhouseInstanceResource_prepaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseInstancePrepaid,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_instance.cdwch_instance_prepaid", "id")),
			},
			{
				ResourceName:            "tencentcloud_clickhouse_instance.cdwch_instance_prepaid",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"time_span"},
			},
		},
	})
}

const testAccClickhouseInstanceBasic = `
resource "tencentcloud_clickhouse_instance" "cdwch_instance" {
  zone="ap-guangzhou-6"
  ha_flag=true
  vpc_id="vpc-j4u8r51f"
  subnet_id="subnet-nvb6lfti"
  product_version="21.8.12.29"
  data_spec {
    spec_name="SCH6"
    count=2
    disk_size=300
  }
  common_spec {
    spec_name="SCH6"
    count=3
    disk_size=300
  }
  charge_type="POSTPAID_BY_HOUR"
  instance_name="tf-test-clickhouse"
}
`

const testAccClickhouseInstancePrepaid = `
resource "tencentcloud_clickhouse_instance" "cdwch_instance_prepaid" {
  zone="ap-guangzhou-6"
  ha_flag=true
  vpc_id="vpc-j4u8r51f"
  subnet_id="subnet-nvb6lfti"
  product_version="21.8.12.29"
  data_spec {
    spec_name="SCH6"
    count=2
    disk_size=300
  }
  common_spec {
    spec_name="SCH6"
    count=3
    disk_size=300
  }
  charge_type="PREPAID"
  renew_flag=1
  time_span=1
  instance_name="tf-test-clickhouse-prepaid"
}
`
