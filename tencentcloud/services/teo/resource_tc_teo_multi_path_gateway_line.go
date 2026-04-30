package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoMultiPathGatewayLine() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoMultiPathGatewayLineCreate,
		Read:   resourceTencentCloudTeoMultiPathGatewayLineRead,
		Update: resourceTencentCloudTeoMultiPathGatewayLineUpdate,
		Delete: resourceTencentCloudTeoMultiPathGatewayLineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},

			"gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Multi-path gateway ID.",
			},

			"line_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Line type. Valid values: `direct`, `proxy`, `custom`.",
			},

			"line_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Line address, format is host:port.",
			},

			"proxy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "L4 proxy instance ID, required when LineType is proxy.",
			},

			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Forwarding rule ID, required when LineType is proxy.",
			},

			// computed
			"line_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Line ID, returned after creation by cloud API.",
			},
		},
	}
}

func resourceTencentCloudTeoMultiPathGatewayLineCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway_line.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = teov20220901.NewCreateMultiPathGatewayLineRequest()
		response  = teov20220901.NewCreateMultiPathGatewayLineResponse()
		zoneId    string
		gatewayId string
		lineId    string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("gateway_id"); ok {
		request.GatewayId = helper.String(v.(string))
		gatewayId = v.(string)
	}

	if v, ok := d.GetOk("line_type"); ok {
		request.LineType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("line_address"); ok {
		request.LineAddress = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_id"); ok {
		request.ProxyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rule_id"); ok {
		request.RuleId = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateMultiPathGatewayLineWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo multi path gateway line failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo multi path gateway line failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.LineId == nil {
		return fmt.Errorf("LineId is nil.")
	}

	lineId = *response.Response.LineId
	d.SetId(strings.Join([]string{zoneId, gatewayId, lineId}, tccommon.FILED_SP))
	return resourceTencentCloudTeoMultiPathGatewayLineRead(d, meta)
}

func resourceTencentCloudTeoMultiPathGatewayLineRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway_line.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	gatewayId := idSplit[1]
	lineId := idSplit[2]

	respData, err := service.DescribeTeoMultiPathGatewayLine(ctx, zoneId, gatewayId, lineId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_multi_path_gateway_line` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("gateway_id", gatewayId)
	_ = d.Set("line_id", lineId)

	if respData.LineType != nil {
		_ = d.Set("line_type", respData.LineType)
	}

	if respData.LineAddress != nil {
		_ = d.Set("line_address", respData.LineAddress)
	}

	if respData.ProxyId != nil {
		_ = d.Set("proxy_id", respData.ProxyId)
	}

	if respData.RuleId != nil {
		_ = d.Set("rule_id", respData.RuleId)
	}

	return nil
}

func resourceTencentCloudTeoMultiPathGatewayLineUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway_line.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	gatewayId := idSplit[1]
	lineId := idSplit[2]

	needChange := false
	mutableArgs := []string{"line_type", "line_address", "proxy_id", "rule_id"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teov20220901.NewModifyMultiPathGatewayLineRequest()
		request.ZoneId = &zoneId
		request.GatewayId = &gatewayId
		request.LineId = &lineId

		if v, ok := d.GetOk("line_type"); ok {
			request.LineType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("line_address"); ok {
			request.LineAddress = helper.String(v.(string))
		}

		if v, ok := d.GetOk("proxy_id"); ok {
			request.ProxyId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("rule_id"); ok {
			request.RuleId = helper.String(v.(string))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyMultiPathGatewayLineWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo multi path gateway line failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoMultiPathGatewayLineRead(d, meta)
}

func resourceTencentCloudTeoMultiPathGatewayLineDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway_line.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDeleteMultiPathGatewayLineRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	gatewayId := idSplit[1]
	lineId := idSplit[2]

	request.ZoneId = &zoneId
	request.GatewayId = &gatewayId
	request.LineId = &lineId

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteMultiPathGatewayLineWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo multi path gateway line failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
