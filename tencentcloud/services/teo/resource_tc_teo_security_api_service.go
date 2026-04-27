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

func ResourceTencentCloudTeoSecurityApiService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoSecurityApiServiceCreate,
		Read:   resourceTencentCloudTeoSecurityApiServiceRead,
		Update: resourceTencentCloudTeoSecurityApiServiceUpdate,
		Delete: resourceTencentCloudTeoSecurityApiServiceDelete,
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
				ForceNew:    true,
				Description: "API service list.",
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
							Description: "Base path of the API service.",
						},
					},
				},
			},

			"api_service_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "API service ID list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"api_resources": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "API resource list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource name.",
						},
						"api_service_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "API service IDs associated with the API resource.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource path.",
						},
						"methods": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "HTTP methods. Supported values: GET, POST, PUT, HEAD, PATCH, OPTIONS, DELETE.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"request_constraint": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request content matching rule.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoSecurityApiServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_api_service.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId        string
		apiServiceIds []string
		request       = teov20220901.NewCreateSecurityAPIServiceRequest()
		response      = teov20220901.NewCreateSecurityAPIServiceResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("api_services"); ok {
		for _, item := range v.([]interface{}) {
			apiServiceMap := item.(map[string]interface{})
			apiService := teov20220901.APIService{}
			if v, ok := apiServiceMap["name"].(string); ok && v != "" {
				apiService.Name = helper.String(v)
			}
			if v, ok := apiServiceMap["base_path"].(string); ok && v != "" {
				apiService.BasePath = helper.String(v)
			}
			request.APIServices = append(request.APIServices, &apiService)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateSecurityAPIServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo security api service failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil {
		return fmt.Errorf("Create teo security api service failed, Response is nil.")
	}

	for _, id := range response.Response.APIServiceIds {
		if id != nil {
			apiServiceIds = append(apiServiceIds, *id)
		}
	}

	// Set composite ID: zone_id + comma-joined api_service_ids
	d.SetId(strings.Join([]string{zoneId, strings.Join(apiServiceIds, ",")}, tccommon.FILED_SP))

	return resourceTencentCloudTeoSecurityApiServiceRead(d, meta)
}

func resourceTencentCloudTeoSecurityApiServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_api_service.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	apiServiceIdsStr := idSplit[1]

	_ = d.Set("zone_id", zoneId)

	// Read API Services
	apiServicesResp, err := service.DescribeTeoSecurityAPIServiceById(ctx, zoneId)
	if err != nil {
		return err
	}

	if apiServicesResp == nil || len(apiServicesResp.APIServices) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_security_api_service` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	// Filter by api_service_ids
	apiServiceIds := strings.Split(apiServiceIdsStr, ",")
	apiServiceIdSet := make(map[string]bool)
	for _, id := range apiServiceIds {
		apiServiceIdSet[id] = true
	}

	filteredApiServices := make([]*teov20220901.APIService, 0)
	for _, svc := range apiServicesResp.APIServices {
		if svc.Id != nil && apiServiceIdSet[*svc.Id] {
			filteredApiServices = append(filteredApiServices, svc)
		}
	}

	if len(filteredApiServices) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_security_api_service` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	// Set api_services
	apiServicesList := make([]map[string]interface{}, 0, len(filteredApiServices))
	matchedIds := make([]string, 0, len(filteredApiServices))
	for _, svc := range filteredApiServices {
		svcMap := map[string]interface{}{}
		if svc.Name != nil {
			svcMap["name"] = svc.Name
		}
		if svc.BasePath != nil {
			svcMap["base_path"] = svc.BasePath
		}
		apiServicesList = append(apiServicesList, svcMap)
		if svc.Id != nil {
			matchedIds = append(matchedIds, *svc.Id)
		}
	}
	_ = d.Set("api_services", apiServicesList)
	_ = d.Set("api_service_ids", matchedIds)

	// Read API Resources
	apiResourcesResp, err := service.DescribeTeoSecurityAPIResourceById(ctx, zoneId)
	if err != nil {
		return err
	}

	if apiResourcesResp != nil && len(apiResourcesResp.APIResources) > 0 {
		// Filter api_resources by associated api_service_ids
		filteredApiResources := make([]*teov20220901.APIResource, 0)
		for _, res := range apiResourcesResp.APIResources {
			if res.APIServiceIds != nil {
				for _, svcId := range res.APIServiceIds {
					if svcId != nil && apiServiceIdSet[*svcId] {
						filteredApiResources = append(filteredApiResources, res)
						break
					}
				}
			}
		}

		if len(filteredApiResources) > 0 {
			apiResourcesList := make([]map[string]interface{}, 0, len(filteredApiResources))
			for _, res := range filteredApiResources {
				resMap := map[string]interface{}{}
				if res.Id != nil {
					resMap["id"] = res.Id
				}
				if res.Name != nil {
					resMap["name"] = res.Name
				}
				if res.APIServiceIds != nil {
					svcIds := make([]string, 0, len(res.APIServiceIds))
					for _, svcId := range res.APIServiceIds {
						if svcId != nil {
							svcIds = append(svcIds, *svcId)
						}
					}
					resMap["api_service_ids"] = svcIds
				}
				if res.Path != nil {
					resMap["path"] = res.Path
				}
				if res.Methods != nil {
					methods := make([]string, 0, len(res.Methods))
					for _, m := range res.Methods {
						if m != nil {
							methods = append(methods, *m)
						}
					}
					resMap["methods"] = methods
				}
				if res.RequestConstraint != nil {
					resMap["request_constraint"] = res.RequestConstraint
				}
				apiResourcesList = append(apiResourcesList, resMap)
			}
			_ = d.Set("api_resources", apiResourcesList)
		}
	}

	return nil
}

func resourceTencentCloudTeoSecurityApiServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_api_service.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId := d.Get("zone_id").(string)

	// Only api_resources is mutable; zone_id and api_services are ForceNew
	if d.HasChange("api_resources") {
		request := teov20220901.NewModifySecurityAPIResourceRequest()
		request.ZoneId = helper.String(zoneId)

		if v, ok := d.GetOk("api_resources"); ok {
			for _, item := range v.([]interface{}) {
				apiResourceMap := item.(map[string]interface{})
				apiResource := teov20220901.APIResource{}
				if v, ok := apiResourceMap["id"].(string); ok && v != "" {
					apiResource.Id = helper.String(v)
				}
				if v, ok := apiResourceMap["name"].(string); ok && v != "" {
					apiResource.Name = helper.String(v)
				}
				if v, ok := apiResourceMap["api_service_ids"].([]interface{}); ok && len(v) > 0 {
					for _, svcId := range v {
						apiResource.APIServiceIds = append(apiResource.APIServiceIds, helper.String(svcId.(string)))
					}
				}
				if v, ok := apiResourceMap["path"].(string); ok && v != "" {
					apiResource.Path = helper.String(v)
				}
				if v, ok := apiResourceMap["methods"].([]interface{}); ok && len(v) > 0 {
					for _, method := range v {
						apiResource.Methods = append(apiResource.Methods, helper.String(method.(string)))
					}
				}
				if v, ok := apiResourceMap["request_constraint"].(string); ok && v != "" {
					apiResource.RequestConstraint = helper.String(v)
				}
				request.APIResources = append(request.APIResources, &apiResource)
			}
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifySecurityAPIResourceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo security api resource failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoSecurityApiServiceRead(d, meta)
}

func resourceTencentCloudTeoSecurityApiServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_security_api_service.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	// Get zone_id and api_service_ids from d.Get() instead of parsing d.Id()
	zoneId := d.Get("zone_id").(string)
	apiServiceIdsRaw := d.Get("api_service_ids").([]interface{})

	var (
		request  = teov20220901.NewDeleteSecurityAPIServiceRequest()
		response = teov20220901.NewDeleteSecurityAPIServiceResponse()
	)

	request.ZoneId = helper.String(zoneId)

	for _, id := range apiServiceIdsRaw {
		request.APIServiceIds = append(request.APIServiceIds, helper.String(id.(string)))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteSecurityAPIServiceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo security api service failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
