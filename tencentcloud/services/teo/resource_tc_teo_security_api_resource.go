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

func ResourceTencentCloudTeoSecurityAPIResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoSecurityAPIResourceCreate,
		Read:   resourceTencentCloudTeoSecurityAPIResourceRead,
		Update: resourceTencentCloudTeoSecurityAPIResourceUpdate,
		Delete: resourceTencentCloudTeoSecurityAPIResourceDelete,
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

			"api_resources": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "API resource configuration. Only one resource is allowed per request.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "API resource name.",
						},
						"path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "API resource path, e.g. `/ava`.",
						},
						"api_service_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Associated API service ID list.",
						},
						"methods": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Allowed HTTP methods. Valid values: `GET`, `POST`, `PUT`, `HEAD`, `PATCH`, `OPTIONS`, `DELETE`.",
						},
						"request_constraint": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request content matching rule expression.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "API resource ID, e.g. `apires-xxxxxxxx`.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoSecurityAPIResourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_api_resource.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request       = teo.NewCreateSecurityAPIResourceRequest()
		response      = teo.NewCreateSecurityAPIResourceResponse()
		zoneId        string
		apiResourceId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("api_resources"); ok {
		apiRes := buildSecurityAPIResourceFromMap(v.([]interface{})[0].(map[string]interface{}), "")
		request.APIResources = []*teo.APIResource{apiRes}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateSecurityAPIResourceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo security api resource failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo security api resource failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if len(response.Response.APIResourceIds) == 0 || response.Response.APIResourceIds[0] == nil {
		return fmt.Errorf("APIResourceIds is empty.")
	}

	apiResourceId = *response.Response.APIResourceIds[0]
	d.SetId(strings.Join([]string{zoneId, apiResourceId}, tccommon.FILED_SP))
	return resourceTencentCloudTeoSecurityAPIResourceRead(d, meta)
}

func resourceTencentCloudTeoSecurityAPIResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_api_resource.read")()
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
	apiResourceId := idSplit[1]

	respData, err := service.DescribeTeoSecurityAPIResourceById(ctx, zoneId, apiResourceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_security_api_resource` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	apiResMap := map[string]interface{}{}

	if respData.Id != nil {
		apiResMap["id"] = *respData.Id
	}

	if respData.Name != nil {
		apiResMap["name"] = *respData.Name
	}

	if respData.Path != nil {
		apiResMap["path"] = *respData.Path
	}

	if len(respData.APIServiceIds) > 0 {
		ids := make([]string, 0, len(respData.APIServiceIds))
		for _, id := range respData.APIServiceIds {
			if id != nil {
				ids = append(ids, *id)
			}
		}
		apiResMap["api_service_ids"] = ids
	}

	if len(respData.Methods) > 0 {
		methods := make([]string, 0, len(respData.Methods))
		for _, m := range respData.Methods {
			if m != nil {
				methods = append(methods, *m)
			}
		}
		apiResMap["methods"] = methods
	}

	if respData.RequestConstraint != nil {
		apiResMap["request_constraint"] = *respData.RequestConstraint
	}

	_ = d.Set("api_resources", []interface{}{apiResMap})

	return nil
}

func resourceTencentCloudTeoSecurityAPIResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_api_resource.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewModifySecurityAPIResourceRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	apiResourceId := idSplit[1]

	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("api_resources"); ok {
		apiRes := buildSecurityAPIResourceFromMap(v.([]interface{})[0].(map[string]interface{}), apiResourceId)
		request.APIResources = []*teo.APIResource{apiRes}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifySecurityAPIResourceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update teo security api resource failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudTeoSecurityAPIResourceRead(d, meta)
}

func resourceTencentCloudTeoSecurityAPIResourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_api_resource.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewDeleteSecurityAPIResourceRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	request.ZoneId = helper.String(idSplit[0])
	request.APIResourceIds = []*string{helper.String(idSplit[1])}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteSecurityAPIResourceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo security api resource failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// buildSecurityAPIResourceFromMap converts a schema map block into *teo.APIResource.
// Pass id="" when creating (API does not accept Id on create).
func buildSecurityAPIResourceFromMap(m map[string]interface{}, id string) *teo.APIResource {
	res := &teo.APIResource{}

	if id != "" {
		res.Id = helper.String(id)
	}

	if val, ok := m["name"].(string); ok && val != "" {
		res.Name = helper.String(val)
	}

	if val, ok := m["path"].(string); ok && val != "" {
		res.Path = helper.String(val)
	}

	if val, ok := m["api_service_ids"].([]interface{}); ok {
		for _, s := range val {
			res.APIServiceIds = append(res.APIServiceIds, helper.String(s.(string)))
		}
	}

	if val, ok := m["methods"].([]interface{}); ok {
		for _, mv := range val {
			res.Methods = append(res.Methods, helper.String(mv.(string)))
		}
	}

	if val, ok := m["request_constraint"].(string); ok && val != "" {
		res.RequestConstraint = helper.String(val)
	}

	return res
}
