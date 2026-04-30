package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

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
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
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
				Description: "Gateway type. Valid values: `cloud`, `private`.",
			},

			"gateway_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Gateway name, up to 16 characters.",
			},

			"gateway_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Gateway port, range 1-65535 (excluding 8888).",
			},

			"region_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Gateway region, required when GatewayType is cloud.",
			},

			"gateway_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Gateway IP address, required when GatewayType is private.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Gateway status. Valid values: `online` (enable), `offline` (disable). If not set, the value is populated by the server.",
			},

			// computed
			"gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Gateway ID.",
			},

			"need_confirm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the gateway origin IP list needs reconfirmation.",
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
		request   = teov20220901.NewCreateMultiPathGatewayRequest()
		response  = teov20220901.NewCreateMultiPathGatewayResponse()
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

	// Wait for gateway to be online (CreateMultiPathGateway is async)
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		respData, e := service.DescribeTeoMultiPathGatewayById(ctx, zoneId, gatewayId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if respData == nil {
			return resource.RetryableError(fmt.Errorf("teo multi path gateway is still creating"))
		}
		if respData.Status != nil && *respData.Status == "creating" {
			return resource.RetryableError(fmt.Errorf("teo multi path gateway is still creating, current status: %s", *respData.Status))
		}
		if respData.Status != nil && *respData.Status == "online" {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("teo multi path gateway status is unexpected: %s", *respData.Status))
	})
	if err != nil {
		log.Printf("[CRITAL]%s wait for teo multi path gateway online failed, reason:%+v", logId, err)
		return err
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
		return fmt.Errorf("id is broken,%s", d.Id())
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

	if respData.GatewayName != nil {
		_ = d.Set("gateway_name", respData.GatewayName)
	}

	if respData.GatewayType != nil {
		_ = d.Set("gateway_type", respData.GatewayType)
	}

	if respData.GatewayPort != nil {
		_ = d.Set("gateway_port", int(*respData.GatewayPort))
	}

	if respData.Status != nil {
		// // disable -> offline: The modification API now supports the 'offline' status; however, since the actual result returned is 'disable', this scenario is handled for backward compatibility.
		if *respData.Status == "disable" {
			_ = d.Set("status", "offline")
		} else {
			_ = d.Set("status", respData.Status)
		}
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
		return fmt.Errorf("id is broken,%s", d.Id())
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
		request := teov20220901.NewModifyMultiPathGatewayRequest()
		request.ZoneId = &zoneId
		request.GatewayId = &gatewayId

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

		// Wait for gateway to be stable after modify (ModifyMultiPathGateway is async)
		service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			respData, e := service.DescribeTeoMultiPathGatewayById(ctx, zoneId, gatewayId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			if respData == nil {
				return resource.NonRetryableError(fmt.Errorf("teo multi path gateway not found after update"))
			}
			if respData.Status != nil && *respData.Status == "creating" {
				return resource.RetryableError(fmt.Errorf("teo multi path gateway is still being modified, current status: %s", *respData.Status))
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s wait for teo multi path gateway stable after update failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			targetStatus := v.(string)
			statusRequest := teov20220901.NewModifyMultiPathGatewayStatusRequest()
			statusRequest.ZoneId = &zoneId
			statusRequest.GatewayId = &gatewayId
			statusRequest.GatewayStatus = helper.String(targetStatus)

			reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyMultiPathGatewayStatusWithContext(ctx, statusRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, statusRequest.GetAction(), statusRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if reqErr != nil {
				log.Printf("[CRITAL]%s modify teo multi path gateway status failed, reason:%+v", logId, reqErr)
				return reqErr
			}

			// Wait for gateway status to reach the target value (ModifyMultiPathGatewayStatus is async)
			service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
			err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				respData, e := service.DescribeTeoMultiPathGatewayById(ctx, zoneId, gatewayId)
				if e != nil {
					return tccommon.RetryError(e)
				}
				if respData == nil {
					return resource.NonRetryableError(fmt.Errorf("teo multi path gateway not found after modify status"))
				}
				if respData.Status != nil && *respData.Status == "creating" {
					return resource.RetryableError(fmt.Errorf("teo multi path gateway is still transitioning, current status: %s", *respData.Status))
				}
				if respData.Status != nil && (*respData.Status == targetStatus || (targetStatus == "offline" && *respData.Status == "disable")) {
					return nil
				}
				// Still keep retrying while status has not yet converged to target value.
				return resource.RetryableError(fmt.Errorf("teo multi path gateway status has not reached target %s, current status: %v", targetStatus, respData.Status))
			})
			if err != nil {
				log.Printf("[CRITAL]%s wait for teo multi path gateway status to reach target failed, reason:%+v", logId, err)
				return err
			}
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
		request = teov20220901.NewDeleteMultiPathGatewayRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	gatewayId := idSplit[1]

	request.ZoneId = &zoneId
	request.GatewayId = &gatewayId

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
