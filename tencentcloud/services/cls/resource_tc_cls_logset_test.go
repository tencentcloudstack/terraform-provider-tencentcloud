package cls_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	localcls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"

	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("tencentcloud_cls_logset", &resource.Sweeper{
		Name: "tencentcloud_cls_logset",
		F:    testSweepClsLogset,
	})
}

func testSweepClsLogset(region string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	sharedClient, err := tcacctest.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
	}
	client := sharedClient.(tccommon.ProviderMeta)

	clsService := localcls.NewClsService(client.GetAPIV3Conn())

	instances, err := clsService.DescribeClsLogsetByFilter(ctx, nil)
	if err != nil {
		return fmt.Errorf("get instance list error: %s", err.Error())
	}

	for _, v := range instances {
		instanceId := v.LogsetId
		instanceName := v.LogsetName

		now := time.Now()

		createTime := tccommon.StringToTime(*v.CreateTime)
		interval := now.Sub(createTime).Minutes()
		if strings.HasPrefix(*instanceName, tcacctest.KeepResource) || strings.HasPrefix(*instanceName, tcacctest.DefaultResource) {
			continue
		}
		// less than 30 minute, not delete
		if tccommon.NeedProtect == 1 && int64(interval) < 30 {
			continue
		}

		if err = clsService.DeleteClsLogsetById(ctx, *instanceId); err != nil {
			log.Printf("[ERROR] sweep instance %s error: %s", *instanceId, err.Error())
		}
	}

	return nil
}

func TestAccTencentCloudClsLogset_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsLogset_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClsLogsetExists("tencentcloud_cls_logset.logset"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_logset.logset", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_cls_logset.logset", "logset_name", "tf-logset-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_logset.logset",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckClsLogsetExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("[CHECK][CLS logset][Exists] check: CLS logset %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[CHECK][CLS logset][Exists] check: CLS logset id is not set")
		}
		service := localcls.NewClsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		resourceId := rs.Primary.ID
		instance, err := service.DescribeClsLogset(ctx, resourceId)
		if err != nil {
			return err
		}
		if instance == nil {
			return fmt.Errorf("[CHECK][CLS logset][Exists] id %s is not exist", rs.Primary.ID)
		}
		return nil
	}
}

const testAccClsLogset_basic = `
resource "tencentcloud_cls_logset" "logset" {
  logset_name = "tf-logset-test"
  tags        = {
    "test" = "test"
  }
}
`
