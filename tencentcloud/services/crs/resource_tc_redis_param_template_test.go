package crs_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccrs "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/crs"

	"context"
	"testing"
	"time"

	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_redis_param_template
	resource.AddTestSweepers("tencentcloud_redis_param_template", &resource.Sweeper{
		Name: "tencentcloud_redis_param_template",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svccrs.NewRedisService(client)

			request := redis.NewDescribeParamTemplatesRequest()
			params, err := service.DescribeParamTemplates(ctx, request)
			if err != nil {
				return err
			}

			for i := range params {
				item := params[i]
				created := time.Time{}
				if tcacctest.IsResourcePersist(*item.Name, &created) {
					continue
				}
				dReq := redis.NewDeleteParamTemplateRequest()
				dReq.TemplateId = item.TemplateId
				err = service.DeleteParamTemplate(ctx, dReq)
				if err != nil {
					return err
				}
			}
			return nil
		},
	})
}

func TestAccTencentCloudRedisParamTemplateResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRedisParamTemplate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_param_template.param_template", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_redis_param_template.copied", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param_template.param_template", "name", "test-tf-template"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param_template.param_template", "description", "test tf template"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param_template.param_template", "product_type", "9"),
				),
			},
			{
				ResourceName:            "tencentcloud_redis_param_template.param_template",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"params_override", "product_type"},
			},
			{
				Config: testAccRedisParamTemplate_update,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_redis_param_template.param_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param_template.param_template", "name", "test-tf-template2"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param_template.param_template", "description", "test tf template 22"),
					resource.TestCheckResourceAttr("tencentcloud_redis_param_template.param_template", "product_type", "9"),
				),
			},
		},
	})
}

const testAccRedisParamTemplate = `
resource "tencentcloud_redis_param_template" "param_template" {
  name = "test-tf-template"
  description = "test tf template"
  product_type = 9
  params_override {
    key = "auto-failback"
    value = "no"
  }
  params_override {
    key = "hz"
    value = 20
  }
}

resource "tencentcloud_redis_param_template" "copied" {
  name = "test-tf-copied"
  description = "test tf copied"
  template_id = tencentcloud_redis_param_template.param_template.id
}
`

const testAccRedisParamTemplate_update = `
resource "tencentcloud_redis_param_template" "param_template" {
  name = "test-tf-template2"
  description = "test tf template 22"
  product_type = 9
  params_override {
    key = "hz"
    value = 30
  }
  params_override {
    key = "timeout"
    value = "3600"
  }
}

resource "tencentcloud_redis_param_template" "copied" {
  name = "test-tf-copied"
  description = "test tf copied"
  template_id = tencentcloud_redis_param_template.param_template.id
}
`
