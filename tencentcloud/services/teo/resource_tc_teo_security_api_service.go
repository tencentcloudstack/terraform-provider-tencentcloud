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

func ResourceTencentCloudTeoSecurityAPIService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoSecurityAPIServiceCreate,
		Read:   resourceTencentCloudTeoSecurityAPIServiceRead,
		Update: resourceTencentCloudTeoSecurityAPIServiceUpdate,
		Delete: resourceTencentCloudTeoSecurityAPIServiceDelete,
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

			"api_services": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "API service configuration. Only one service is allowed per request.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "API service name.",
						},
						"base_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "API service base path, e.g. `/tt`.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API service ID, e.g. `apisrv-xxxxxxxx`.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoSecurityAPIServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_api_service.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request      = teo.NewCreateSecurityAPIServiceRequest()
		response     = teo.NewCreateSecurityAPIServiceResponse()
		zoneId       string
		apiServiceId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("api_services"); ok {
		svcMap := v.([]interface{})[0].(map[string]interface{})
		request.APIServices = []*teo.APIService{buildAPIServiceFromMap(svcMap, "")}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateSecurityAPIServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo security api service failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo security api service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if len(response.Response.APIServiceIds) == 0 || response.Response.APIServiceIds[0] == nil {
		return fmt.Errorf("APIServiceIds is empty.")
	}

	apiServiceId = *response.Response.APIServiceIds[0]
	d.SetId(strings.Join([]string{zoneId, apiServiceId}, tccommon.FILED_SP))
	return resourceTencentCloudTeoSecurityAPIServiceRead(d, meta)
}

func resourceTencentCloudTeoSecurityAPIServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_api_service.read")()
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
	apiServiceId := idSplit[1]

	respData, err := service.DescribeTeoSecurityAPIServiceById(ctx, zoneId, apiServiceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_security_api_service` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	svcMap := map[string]interface{}{}

	if respData.Id != nil {
		svcMap["id"] = *respData.Id
	}

	if respData.Name != nil {
		svcMap["name"] = *respData.Name
	}

	if respData.BasePath != nil {
		svcMap["base_path"] = *respData.BasePath
	}

	_ = d.Set("api_services", []interface{}{svcMap})

	return nil
}

func resourceTencentCloudTeoSecurityAPIServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_api_service.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewModifySecurityAPIServiceRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	apiServiceId := idSplit[1]

	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("api_services"); ok {
		svcMap := v.([]interface{})[0].(map[string]interface{})
		request.APIServices = []*teo.APIService{buildAPIServiceFromMap(svcMap, apiServiceId)}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifySecurityAPIServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update teo security api service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudTeoSecurityAPIServiceRead(d, meta)
}

func resourceTencentCloudTeoSecurityAPIServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_api_service.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewDeleteSecurityAPIServiceRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	request.ZoneId = helper.String(idSplit[0])
	request.APIServiceIds = []*string{helper.String(idSplit[1])}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteSecurityAPIServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo security api service failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// buildAPIServiceFromMap converts a schema map block to *teo.APIService.
// Pass id="" when creating (API does not accept Id on create).
func buildAPIServiceFromMap(m map[string]interface{}, id string) *teo.APIService {
	svc := &teo.APIService{}

	if id != "" {
		svc.Id = helper.String(id)
	}

	if val, ok := m["name"].(string); ok && val != "" {
		svc.Name = helper.String(val)
	}

	if val, ok := m["base_path"].(string); ok && val != "" {
		svc.BasePath = helper.String(val)
	}

	return svc
}
