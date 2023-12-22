package tcmg_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMonitorGrafanaSsoAccount_basic -v
func TestAccTencentCloudMonitorGrafanaSsoAccount_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckSsoAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorGrafanaSsoAccount,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSsoAccountExists("tencentcloud_monitor_grafana_sso_account.ssoAccount"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_sso_account.ssoAccount", "user_id", "100027012454"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_sso_account.ssoAccount", "notes", "desc-100027012454"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_sso_account.ssoAccount", "role.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_sso_account.ssoAccount", "role.0.role", "Admin"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_grafana_sso_account.ssoAccount", "role.0.organization", "Main Org."),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_grafana_sso_account.ssoAccount",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSsoAccountDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_grafana_sso_account" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		userId := idSplit[1]

		ssoAccount, err := service.DescribeMonitorSsoAccount(ctx, instanceId, userId)
		if err != nil {
			return err
		}

		if ssoAccount != nil {
			return fmt.Errorf("SsoAccount %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckSsoAccountExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		userId := idSplit[1]

		service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		ssoAccount, err := service.DescribeMonitorSsoAccount(ctx, instanceId, userId)
		if err != nil {
			return err
		}

		if ssoAccount == nil {
			return fmt.Errorf("SsoAccount %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testMonitorGrafanaSsoAccountVar = `
variable "instance_id" {
  default = "` + tcacctest.DefaultGrafanaInstanceId + `"
}
`

const testAccMonitorGrafanaSsoAccount = testMonitorGrafanaSsoAccountVar + `

resource "tencentcloud_monitor_grafana_sso_account" "ssoAccount" {
  instance_id = var.instance_id
  user_id     = "100027012454"
  notes       = "desc-100027012454"
  role {
    organization  = "Main Org."
    role          = "Admin"
  }
}

`
