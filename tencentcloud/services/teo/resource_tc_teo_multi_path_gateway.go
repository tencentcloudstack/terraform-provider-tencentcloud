package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
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
				Description: "Gateway type. Valid values: `cloud` (cloud gateway managed by Tencent Cloud), `private` (private gateway deployed by user).",
			},

			"gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Gateway name, up to 16 characters, available characters (a-z, A-Z, 0-9, -, _).",
			},

			"gateway_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Gateway port, range 1~65535 (except 8888).",
			},

			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Gateway region, required when GatewayType is cloud. You can get RegionId list from DescribeMultiPathGatewayRegions API.",
			},

			"gateway_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Gateway IP address, required when GatewayType is private. Please ensure the address has been registered in Tencent Cloud Multi-Path Gateway system.",
			},

			// Computed
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
				Description: "Whether the gateway origin IP list needs to be confirmed. Valid values: `true` (origin IP list changed, need confirmation), `false` (no change, no need to confirm).",
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
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

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

	if response.Response.GatewayId == nil {
		return fmt.Errorf("GatewayId is nil.")
	}

	gatewayId = *response.Response.GatewayId
	d.SetId(strings.Join([]string{zoneId, gatewayId}, tccommon.FILED_SP))

	// Wait for gateway to become online
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if _, err := (&resource.StateChangeConf{
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{"creating"},
		Target:     []string{"online"},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Refresh: func() (interface{}, string, error) {
			gateway, e := service.DescribeTeoMultiPathGatewayById(ctx, zoneId, gatewayId)
			if e != nil {
				return nil, "", e
			}
			if gateway == nil {
				return nil, "creating", nil
			}
			status := "creating"
			if gateway.Status != nil {
				status = *gateway.Status
			}
			return gateway, status, nil
		},
	}).WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for teo multi path gateway (%s) to become online: %s", d.Id(), err)
	}

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
	_ = d.Set("gateway_id", gatewayId)

	if respData.GatewayType != nil {
		_ = d.Set("gateway_type", respData.GatewayType)
	}

	if respData.GatewayName != nil {
		_ = d.Set("gateway_name", respData.GatewayName)
	}

	if respData.GatewayPort != nil {
		_ = d.Set("gateway_port", respData.GatewayPort)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.GatewayIP != nil {
		_ = d.Set("gateway_ip", respData.GatewayIP)
	}

	if respData.RegionId != nil {
		_ = d.Set("region_id", respData.RegionId)
	}

	if respData.NeedConfirm != nil {
		_ = d.Set("need_confirm", respData.NeedConfirm)
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

	if d.HasChanges("gateway_name", "gateway_ip", "gateway_port") {
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
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
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

	// Check if gateway exists before deleting
	gateway, err := service.DescribeTeoMultiPathGatewayById(ctx, zoneId, gatewayId)
	if err != nil {
		return err
	}

	if gateway == nil {
		return nil
	}

	request := teo.NewDeleteMultiPathGatewayRequest()
	request.ZoneId = helper.String(zoneId)
	request.GatewayId = helper.String(gatewayId)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteMultiPathGatewayWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo multi path gateway failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
