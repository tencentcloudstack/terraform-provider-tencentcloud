package apigateway

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apigateway "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/apigateway/v20180808"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAPIGatewayApiAppAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAPIGatewayApiAppAttachmentCreate,
		Read:   resourceTencentCloudAPIGatewayApiAppAttachmentRead,
		Delete: resourceTencentCloudAPIGatewayApiAppAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"api_app_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the application to be bound.",
			},
			"environment": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The environment to be bound.",
			},
			"service_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the service to be bound.",
			},
			"api_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The unique ID of the API to be bound.",
			},
		},
	}
}

func resourceTencentCloudAPIGatewayApiAppAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_api_app_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = apigateway.NewBindApiAppRequest()
		apiAppId    string
		environment string
		serviceId   string
		apiId       string
	)

	if v, ok := d.GetOk("api_app_id"); ok {
		request.ApiAppId = helper.String(v.(string))
		apiAppId = v.(string)
	}

	if v, ok := d.GetOk("environment"); ok {
		request.Environment = helper.String(v.(string))
		environment = v.(string)
	}

	if v, ok := d.GetOk("service_id"); ok {
		request.ServiceId = helper.String(v.(string))
		serviceId = v.(string)
	}

	if v, ok := d.GetOk("api_id"); ok {
		request.ApiId = helper.String(v.(string))
		apiId = v.(string)
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAPIGatewayClient().BindApiApp(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if !*result.Response.Result {
			e = fmt.Errorf(" create apigateway apiAppAttachment result: false.")
			return resource.NonRetryableError(e)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create apigateway apiAppAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{apiAppId, environment, serviceId, apiId}, tccommon.FILED_SP))
	return resourceTencentCloudAPIGatewayApiAppAttachmentRead(d, meta)
}

func resourceTencentCloudAPIGatewayApiAppAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_api_app_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("api_gateway_api_app_attachment id is broken, id is %s", d.Id())
	}
	apiAppId := idSplit[0]
	environment := idSplit[1]
	serviceId := idSplit[2]
	apiId := idSplit[3]

	apiAppAttachment, err := service.DescribeAPIGatewayApiAppAttachmentById(ctx, apiAppId, environment, serviceId, apiId)
	if err != nil {
		return err
	}

	if apiAppAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ApigatewayApiAppAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if apiAppAttachment.ApiAppId != nil {
		_ = d.Set("api_app_id", apiAppAttachment.ApiAppId)
	}

	if apiAppAttachment.EnvironmentName != nil {
		_ = d.Set("environment", apiAppAttachment.EnvironmentName)
	}

	if apiAppAttachment.ServiceId != nil {
		_ = d.Set("service_id", apiAppAttachment.ServiceId)
	}

	if apiAppAttachment.ApiId != nil {
		_ = d.Set("api_id", apiAppAttachment.ApiId)
	}

	return nil
}

func resourceTencentCloudAPIGatewayApiAppAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_api_gateway_api_app_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = APIGatewayService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 4 {
		return fmt.Errorf("api_gateway_api_app_attachment id is broken, id is %s", d.Id())
	}
	apiAppId := idSplit[0]
	environment := idSplit[1]
	serviceId := idSplit[2]
	apiId := idSplit[3]

	if err := service.DeleteAPIGatewayApiAppAttachmentById(ctx, apiAppId, environment, serviceId, apiId); err != nil {
		return err
	}

	return nil
}
