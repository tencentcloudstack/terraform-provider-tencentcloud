package ga2

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGa2AccelerateArea() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGa2AccelerateAreaCreate,
		Read:   resourceTencentCloudGa2AccelerateAreaRead,
		Update: resourceTencentCloudGa2AccelerateAreaUpdate,
		Delete: resourceTencentCloudGa2AccelerateAreaDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"global_accelerator_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Global accelerator instance ID.",
			},

			"accelerator_areas": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Accelerate area configuration list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accelerate_region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Accelerate region.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Bandwidth.",
						},
						"isp_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "BGP",
							Description: "ISP type. Supports `BGP`, `三网`, `精品`. Default is `BGP`.",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "IPv4",
							Description: "IP version. Only supports `IPv4`. Default is `IPv4`.",
						},
					},
				},
			},

			"accelerate_area_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Accelerate area information returned from the API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accelerator_area_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Accelerator area ID.",
						},
						"accelerate_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Accelerate region.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bandwidth.",
						},
						"isp_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ISP type.",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP version.",
						},
						"ip_address": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "IP address list.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudGa2AccelerateAreaCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request            = ga2v20250115.NewCreateAccelerateAreasRequest()
		response           = ga2v20250115.NewCreateAccelerateAreasResponse()
		globalAcceleratorId string
	)

	if v, ok := d.GetOk("global_accelerator_id"); ok {
		globalAcceleratorId = v.(string)
		request.GlobalAcceleratorId = helper.String(globalAcceleratorId)
	}

	if v, ok := d.GetOk("accelerator_areas"); ok {
		for _, item := range v.([]interface{}) {
			areaMap := item.(map[string]interface{})
			area := ga2v20250115.AcceleratorAreas{}
			if v, ok := areaMap["accelerate_region"].(string); ok && v != "" {
				area.AccelerateRegion = helper.String(v)
			}

			if v, ok := areaMap["bandwidth"].(int); ok {
				area.Bandwidth = helper.Uint64(uint64(v))
			}

			if v, ok := areaMap["isp_type"].(string); ok && v != "" {
				area.IspType = helper.String(v)
			}

			if v, ok := areaMap["ip_version"].(string); ok && v != "" {
				area.IpVersion = helper.String(v)
			}

			request.AcceleratorAreas = append(request.AcceleratorAreas, &area)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateAccelerateAreasWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ga2 accelerate area failed, Response is nil"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create ga2 accelerate area failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.TaskId == nil {
		return fmt.Errorf("Create ga2 accelerate area failed, TaskId is nil")
	}

	d.SetId(globalAcceleratorId)
	return resourceTencentCloudGa2AccelerateAreaRead(d, meta)
}

func resourceTencentCloudGa2AccelerateAreaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	globalAcceleratorId := d.Id()

	accelerateAreas, err := service.DescribeAccelerateAreasByGlobalAcceleratorId(ctx, globalAcceleratorId)
	if err != nil {
		return err
	}

	if len(accelerateAreas) == 0 {
		log.Printf("[WARN]%s resource `tencentcloud_ga2_accelerate_area` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("global_accelerator_id", globalAcceleratorId)

	accelerateAreaSetList := make([]map[string]interface{}, 0, len(accelerateAreas))
	for _, area := range accelerateAreas {
		areaMap := map[string]interface{}{}
		if area.AcceleratorAreaId != nil {
			areaMap["accelerator_area_id"] = area.AcceleratorAreaId
		}

		if area.AccelerateRegion != nil {
			areaMap["accelerate_region"] = area.AccelerateRegion
		}

		if area.Bandwidth != nil {
			areaMap["bandwidth"] = area.Bandwidth
		}

		if area.IspType != nil {
			areaMap["isp_type"] = area.IspType
		}

		if area.IpVersion != nil {
			areaMap["ip_version"] = area.IpVersion
		}

		if area.IpAddress != nil {
			areaMap["ip_address"] = helper.StringsInterfaces(area.IpAddress)
		}

		accelerateAreaSetList = append(accelerateAreaSetList, areaMap)
	}

	_ = d.Set("accelerate_area_set", accelerateAreaSetList)

	return nil
}

func resourceTencentCloudGa2AccelerateAreaUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	globalAcceleratorId := d.Id()

	if d.HasChange("accelerator_areas") {
		request := ga2v20250115.NewModifyAccelerateAreasRequest()
		request.GlobalAcceleratorId = helper.String(globalAcceleratorId)

		if v, ok := d.GetOk("accelerator_areas"); ok {
			for _, item := range v.([]interface{}) {
				areaMap := item.(map[string]interface{})
				area := ga2v20250115.AcceleratorAreas{}
				if v, ok := areaMap["accelerate_region"].(string); ok && v != "" {
					area.AccelerateRegion = helper.String(v)
				}

				if v, ok := areaMap["bandwidth"].(int); ok {
					area.Bandwidth = helper.Uint64(uint64(v))
				}

				if v, ok := areaMap["isp_type"].(string); ok && v != "" {
					area.IspType = helper.String(v)
				}

				if v, ok := areaMap["ip_version"].(string); ok && v != "" {
					area.IpVersion = helper.String(v)
				}

				request.AcceleratorAreas = append(request.AcceleratorAreas, &area)
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyAccelerateAreasWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify ga2 accelerate area failed, Response is nil"))
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update ga2 accelerate area failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudGa2AccelerateAreaRead(d, meta)
}

func resourceTencentCloudGa2AccelerateAreaDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	globalAcceleratorId := d.Id()

	accelerateAreas, err := service.DescribeAccelerateAreasByGlobalAcceleratorId(ctx, globalAcceleratorId)
	if err != nil {
		return err
	}

	if len(accelerateAreas) == 0 {
		return nil
	}

	acceleratorAreaIds := make([]*string, 0, len(accelerateAreas))
	for _, area := range accelerateAreas {
		if area.AcceleratorAreaId != nil {
			acceleratorAreaIds = append(acceleratorAreaIds, area.AcceleratorAreaId)
		}
	}

	if len(acceleratorAreaIds) == 0 {
		return nil
	}

	request := ga2v20250115.NewDeleteAccelerateAreasRequest()
	request.GlobalAcceleratorId = helper.String(globalAcceleratorId)
	request.AcceleratorAreaIds = acceleratorAreaIds

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteAccelerateAreasWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ga2 accelerate area failed, Response is nil"))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete ga2 accelerate area failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
