package tcr_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctcr "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tcr"

	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-shanghai -sweep-run=tencentcloud_tcr_namespace
	resource.AddTestSweepers("tencentcloud_tcr_namespace", &resource.Sweeper{
		Name: "tencentcloud_tcr_namespace",
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

			namespaces, err := service.DescribeTCRNameSpaces(ctx, instanceId, "test")
			if err != nil {
				return err
			}

			for i := range namespaces {
				n := namespaces[i]
				if tcacctest.IsResourcePersist(*n.Name, nil) {
					continue
				}
				err = service.DeleteTCRNameSpace(ctx, instanceId, *n.Name)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudTcrNamespace_basic_and_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTCRNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTCRNamespace_basic,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "name", "test_ns"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "is_public", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "is_auto_scan", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "is_prevent_vul", "true"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "severity", "medium"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "cve_whitelist_items.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "cve_whitelist_items.0.cve_id", "cve-xxxxx"),
				),
			},
			{
				ResourceName:      "tencentcloud_tcr_namespace.mytcr_namespace",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTCRNamespace_basic_update_remark,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRNamespaceExists("tencentcloud_tcr_namespace.mytcr_namespace"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "name", "test2_ns"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "is_public", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "is_auto_scan", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "is_prevent_vul", "false"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "severity", "high"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "cve_whitelist_items.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_tcr_namespace.mytcr_namespace", "cve_whitelist_items.0.cve_id", "cve-xxxx"),
				),
			},
		},
	})
}

func testAccCheckTCRNamespaceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	tcrService := svctcr.NewTCRService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_tcr_namespace" {
			continue
		}
		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		namespaceName := items[1]
		_, has, err := tcrService.DescribeTCRNameSpaceById(ctx, instanceId, namespaceName)
		if has {
			return fmt.Errorf("TCR namespace still exists")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckTCRNamespaceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("TCR namespace %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("TCR namespace id is not set")
		}
		items := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(items) != 2 {
			return fmt.Errorf("invalid ID %s", rs.Primary.ID)
		}

		instanceId := items[0]
		namespaceName := items[1]
		tcrService := svctcr.NewTCRService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, has, err := tcrService.DescribeTCRNameSpaceById(ctx, instanceId, namespaceName)
		if !has {
			return fmt.Errorf("TCR namespace %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccTCRNamespace_basic = tcacctest.DefaultTCRInstanceData + `

resource "tencentcloud_tcr_namespace" "mytcr_namespace" {
  instance_id 	 = local.tcr_id
  name			 = "test_ns"
  is_public		 = true
  is_auto_scan	 = true
  is_prevent_vul = true
  severity		 = "medium"
  cve_whitelist_items	{
    cve_id = "cve-xxxxx"
  }
}`

const testAccTCRNamespace_basic_update_remark = tcacctest.DefaultTCRInstanceData + `

resource "tencentcloud_tcr_namespace" "mytcr_namespace" {
  instance_id 	 = local.tcr_id
  name        	 = "test2_ns"
  is_public   	 = false
  is_auto_scan	 = false
  is_prevent_vul = false
  severity		 = "high"
  cve_whitelist_items	{
    cve_id = "cve-xxxx"
  }
}`
