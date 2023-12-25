package tcr_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctcr "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcr"

	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-shanghai -sweep-run=tencentcloud_tcr_token
	resource.AddTestSweepers("tencentcloud_tcr_token", &resource.Sweeper{
		Name: "tencentcloud_tcr_token",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()

			service := svctcr.NewTCRService(client)

			var filters []*tcr.Filter
			filters = append(filters, &tcr.Filter{
				Name:   helper.String("RegistryName"),
				Values: []*string{helper.String(tcacctest.DefaultTCRInstanceName)},
			})

			instances, err := service.DescribeTCRInstances(ctx, "", filters)
			if err != nil {
				return err
			}

			if len(instances) == 0 {
				return nil
			}

			instanceId := *instances[0].RegistryId

			tokens, err := service.DescribeTCRTokens(ctx, instanceId, "")
			if err != nil {
				return err
			}

			for i := range tokens {
				token := tokens[i]
				id := *token.Id
				created, err := time.Parse(time.RFC3339, *token.CreatedAt)
				if err != nil {
					created = time.Time{}
				}
				if tcacctest.IsResourcePersist("", &created) {
					continue
				}
				log.Printf("%s -> %s (%s) will delete", instanceId, id, *token.Desc)
				err = service.DeleteTCRLongTermToken(ctx, instanceId, id)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudTcrToken_basic_and_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTCRTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTCRToken_basic,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tcr_token.mytcr_token", "description", "test token"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_token.mytcr_token", "enable", "true"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_token.mytcr_token", "token_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_token.mytcr_token", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_token.mytcr_token", "token"),
					resource.TestCheckResourceAttrSet("tencentcloud_tcr_token.mytcr_token", "user_name"),
				),
				Destroy: false,
			},
			{
				ResourceName:            "tencentcloud_tcr_token.mytcr_token",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token", "user_name"},
			},
			{
				Config: testAccTCRToken_basic_update_remark,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRTokenExists("tencentcloud_tcr_token.mytcr_token"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_token.mytcr_token", "enable", "false"),
				),
			},
		},
	})
}

func testAccCheckTCRTokenDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	tcrService := svctcr.NewTCRService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_token" {
			continue
		}
		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		tokenId := items[1]
		_, has, err := tcrService.DescribeTCRLongTermTokenById(ctx, instanceId, tokenId)
		if has {
			return fmt.Errorf("TCR token still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTCRTokenExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("TCR token %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("TCR token id is not set")
		}
		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		tokenId := items[1]
		tcrService := svctcr.NewTCRService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, has, err := tcrService.DescribeTCRLongTermTokenById(ctx, instanceId, tokenId)
		if !has {
			return fmt.Errorf("TCR token %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTCRToken_basic = tcacctest.DefaultTCRInstanceData + `
resource "tencentcloud_tcr_token" "mytcr_token" {
  instance_id = local.tcr_id
  description       = "test token"
}`

const testAccTCRToken_basic_update_remark = tcacctest.DefaultTCRInstanceData + `

resource "tencentcloud_tcr_token" "mytcr_token" {
  instance_id = local.tcr_id
  description       = "test token"
  enable   = false
}`
