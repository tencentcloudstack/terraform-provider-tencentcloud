package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoMultiPathGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoMultiPathGatewayCreate,
		Read:   resourceTencentCloudTeoMultiPathGatewayRead,
		Update: resourceTencentCloudTeoMultiPathGatewayUpdate,
		Delete: resourceTencentCloudTeoMultiPathGatewayDelete,
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

			"gateway_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Gateway type. Valid values: `cloud` (cloud gateway managed by Tencent Cloud), `private` (self-deployed private gateway).",
			},

			"gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Gateway name, up to 16 characters, available characters (a-z, A-Z, 0-9, -, _).",
			},

			"gateway_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Gateway port, range 1-65535 (except 8888).",
			},

			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Gateway region, required when GatewayType is cloud.",
			},

			"gateway_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Gateway address, required when GatewayType is private.",
			},

			"gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Gateway ID.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Gateway status. Valid values: `creating`, `online`, `offline`, `disable`.",
			},

			"need_confirm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the gateway origin IP list change needs confirmation.",
			},

			"lines": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Line information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"line_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line ID.",
						},
						"line_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line type. Valid values: `direct`, `proxy`, `custom`.",
						},
						"line_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Line address, format is host:port.",
						},
						"proxy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "L4 proxy instance ID, returned when LineType is proxy.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forwarding rule ID, returned when LineType is proxy.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoMultiPathGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = teo.NewCreateMultiPathGatewayRequest()
		response  = teo.NewCreateMultiPathGatewayResponse()
		zoneId    string
		gatewayId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("gateway_type"); ok {
		request.GatewayType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("gateway_name"); ok {
		request.GatewayName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("gateway_port"); ok {
		request.GatewayPort = helper.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("region_id"); ok {
		request.RegionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("gateway_ip"); ok {
		request.GatewayIP = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateMultiPathGatewayWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo multi path gateway failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo multi path gateway failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.GatewayId == nil || *response.Response.GatewayId == "" {
		return fmt.Errorf("GatewayId is empty.")
	}

	gatewayId = *response.Response.GatewayId
	d.SetId(strings.Join([]string{zoneId, gatewayId}, tccommon.FILED_SP))
	return resourceTencentCloudTeoMultiPathGatewayRead(d, meta)
}

func resourceTencentCloudTeoMultiPathGatewayRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	gatewayId := idSplit[1]

	respData, err := service.DescribeTeoMultiPathGatewayById(ctx, zoneId, gatewayId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_multi_path_gateway` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if respData.GatewayId != nil {
		_ = d.Set("gateway_id", respData.GatewayId)
	}

	if respData.GatewayType != nil {
		_ = d.Set("gateway_type", respData.GatewayType)
	}

	if respData.GatewayName != nil {
		_ = d.Set("gateway_name", respData.GatewayName)
	}

	if respData.GatewayPort != nil {
		_ = d.Set("gateway_port", respData.GatewayPort)
	}

	if respData.RegionId != nil {
		_ = d.Set("region_id", respData.RegionId)
	}

	if respData.GatewayIP != nil {
		_ = d.Set("gateway_ip", respData.GatewayIP)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.NeedConfirm != nil {
		_ = d.Set("need_confirm", respData.NeedConfirm)
	}

	if respData.Lines != nil && len(respData.Lines) > 0 {
		linesList := make([]map[string]interface{}, 0, len(respData.Lines))
		for _, line := range respData.Lines {
			lineMap := map[string]interface{}{}
			if line.LineId != nil {
				lineMap["line_id"] = line.LineId
			}
			if line.LineType != nil {
				lineMap["line_type"] = line.LineType
			}
			if line.LineAddress != nil {
				lineMap["line_address"] = line.LineAddress
			}
			if line.ProxyId != nil {
				lineMap["proxy_id"] = line.ProxyId
			}
			if line.RuleId != nil {
				lineMap["rule_id"] = line.RuleId
			}
			linesList = append(linesList, lineMap)
		}
		_ = d.Set("lines", linesList)
	}

	return nil
}

func resourceTencentCloudTeoMultiPathGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	gatewayId := idSplit[1]

	needChange := false
	mutableArgs := []string{"gateway_name", "gateway_ip", "gateway_port"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teo.NewModifyMultiPathGatewayRequest()
		request.ZoneId = helper.String(zoneId)
		request.GatewayId = helper.String(gatewayId)

		if v, ok := d.GetOk("gateway_name"); ok {
			request.GatewayName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("gateway_ip"); ok {
			request.GatewayIP = helper.String(v.(string))
		}

		if v, ok := d.GetOk("gateway_port"); ok {
			request.GatewayPort = helper.Int64(int64(v.(int)))
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyMultiPathGatewayWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo multi path gateway failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoMultiPathGatewayRead(d, meta)
}

func resourceTencentCloudTeoMultiPathGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_multi_path_gateway.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	gatewayId := idSplit[1]

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := service.DeleteTeoMultiPathGateway(ctx, zoneId, gatewayId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo multi path gateway failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
