package cdwch_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	cdwch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwch/v20200915"

	localcdwch "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdwch"
)

func TestAccTencentCloudClickhouseKeyvalConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckClickhouseKeyvalConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClickhouseKeyvalConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_keyval_config.keyval_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_keyval_config.keyval_config", "instance_id", "cdwch-pcap78rz"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_keyval_config.keyval_config", "items.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_keyval_config.keyval_config", "items.0.conf_key", "max_open_files"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_keyval_config.keyval_config", "items.0.conf_value", "50000"),
				),
			},
			{
				Config: testAccClickhouseKeyvalConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_clickhouse_keyval_config.keyval_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_keyval_config.keyval_config", "instance_id", "cdwch-pcap78rz"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_keyval_config.keyval_config", "items.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_keyval_config.keyval_config", "items.0.conf_key", "max_open_files"),
					resource.TestCheckResourceAttr("tencentcloud_clickhouse_keyval_config.keyval_config", "items.0.conf_value", "60000"),
				),
			},
			{
				ResourceName:      "tencentcloud_clickhouse_keyval_config.keyval_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClickhouseKeyvalConfigDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := localcdwch.NewCdwchService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_clickhouse_keyval_config" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		configKey := idSplit[1]

		configItems, err := service.DescribeClickhouseKeyvalConfigById(ctx, instanceId)
		if err != nil {
			return err
		}

		resultMap := make(map[string]*cdwch.InstanceConfigInfo)
		for _, item := range configItems {
			resultMap[*item.ConfKey] = item
		}

		item := resultMap[configKey]
		if item != nil {
			return fmt.Errorf("instance keyval config %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

const testAccClickhouseKeyvalConfig = `
resource "tencentcloud_clickhouse_keyval_config" "keyval_config" {
  instance_id = "cdwch-pcap78rz"
  items {
    conf_key   = "max_open_files"
    conf_value = "50000"
  }
}
`

const testAccClickhouseKeyvalConfigUpdate = `
resource "tencentcloud_clickhouse_keyval_config" "keyval_config" {
  instance_id = "cdwch-pcap78rz"
  items {
    conf_key   = "max_open_files"
    conf_value = "60000"
  }
}
`
