package ga2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_accelerator_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The global accelerator instance ID.",
			},
			"accelerator_areas": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The accelerate area configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accelerate_region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The accelerate region.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The bandwidth in Mbps.",
						},
						"isp_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "BGP",
							Description: "ISP type, supports `BGP`, `三网`, `精品`. Defaults to `BGP`.",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "IPv4",
							Description: "IP version, only supports `IPv4`. Defaults to `IPv4`.",
						},
						"accelerator_area_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The accelerator area ID, assigned by the backend.",
						},
						"ip_address": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The assigned IP addresses.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ip_address_info_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detailed IP address information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address.",
									},
									"isp_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ISP type of the IP.",
									},
								},
							},
						},
					},
				},
			},
			"accelerate_area_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The full accelerate area set returned by Read.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accelerate_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The accelerate region.",
						},
						"bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The bandwidth in Mbps.",
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
						"accelerator_area_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The accelerator area ID.",
						},
						"ip_address": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The assigned IP addresses.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ip_address_info_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Detailed IP address information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address.",
									},
									"isp_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ISP type of the IP.",
									},
								},
							},
						},
					},
				},
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The async task ID from the last write operation.",
			},
		},
	}
}

func resourceTencentCloudGa2AccelerateAreaCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request             = ga2v20250115.NewCreateAccelerateAreasRequest()
		response            *ga2v20250115.CreateAccelerateAreasResponse
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
				area.Bandwidth = helper.IntUint64(v)
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
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateAccelerateAreas(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ga2 accelerate area failed, Response is nil"))
		}

		if result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ga2 accelerate area failed, TaskId is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create ga2 accelerate area failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(globalAcceleratorId)
	_ = d.Set("task_id", response.Response.TaskId)

	// Poll until areas appear
	service := NewGa2Service(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		areas, e := service.DescribeAccelerateAreas(ctx, globalAcceleratorId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if len(areas) == 0 {
			return resource.RetryableError(fmt.Errorf("waiting for ga2 accelerate areas to be created"))
		}

		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudGa2AccelerateAreaRead(d, meta)
}

func resourceTencentCloudGa2AccelerateAreaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service             = NewGa2Service(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		globalAcceleratorId = d.Id()
	)

	areas, err := service.DescribeAccelerateAreas(ctx, globalAcceleratorId)
	if err != nil {
		return err
	}

	if len(areas) == 0 {
		log.Printf("[WARN]%s resource `tencentcloud_ga2_accelerate_area` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("global_accelerator_id", globalAcceleratorId)

	accelerateAreaSet := make([]map[string]interface{}, 0, len(areas))
	for _, area := range areas {
		areaMap := map[string]interface{}{}
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

		if area.AcceleratorAreaId != nil {
			areaMap["accelerator_area_id"] = area.AcceleratorAreaId
		}

		if area.IpAddress != nil {
			ipAddressList := make([]string, 0, len(area.IpAddress))
			for _, ip := range area.IpAddress {
				if ip != nil {
					ipAddressList = append(ipAddressList, *ip)
				}
			}
			areaMap["ip_address"] = ipAddressList
		}

		if area.IpAddressInfoSet != nil {
			ipAddressInfoSetList := make([]map[string]interface{}, 0, len(area.IpAddressInfoSet))
			for _, ipInfo := range area.IpAddressInfoSet {
				ipInfoMap := map[string]interface{}{}
				if ipInfo.IpAddress != nil {
					ipInfoMap["ip_address"] = ipInfo.IpAddress
				}

				if ipInfo.IspType != nil {
					ipInfoMap["isp_type"] = ipInfo.IspType
				}

				ipAddressInfoSetList = append(ipAddressInfoSetList, ipInfoMap)
			}
			areaMap["ip_address_info_set"] = ipAddressInfoSetList
		}

		accelerateAreaSet = append(accelerateAreaSet, areaMap)
	}

	_ = d.Set("accelerate_area_set", accelerateAreaSet)

	return nil
}

func resourceTencentCloudGa2AccelerateAreaUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request             = ga2v20250115.NewModifyAccelerateAreasRequest()
		globalAcceleratorId = d.Id()
	)

	request.GlobalAcceleratorId = helper.String(globalAcceleratorId)

	if v, ok := d.GetOk("accelerator_areas"); ok {
		for _, item := range v.([]interface{}) {
			areaMap := item.(map[string]interface{})
			area := ga2v20250115.AcceleratorAreas{}
			if v, ok := areaMap["accelerate_region"].(string); ok && v != "" {
				area.AccelerateRegion = helper.String(v)
			}

			if v, ok := areaMap["bandwidth"].(int); ok {
				area.Bandwidth = helper.IntUint64(v)
			}

			if v, ok := areaMap["isp_type"].(string); ok && v != "" {
				area.IspType = helper.String(v)
			}

			if v, ok := areaMap["ip_version"].(string); ok && v != "" {
				area.IpVersion = helper.String(v)
			}

			if v, ok := areaMap["accelerator_area_id"].(string); ok && v != "" {
				area.AcceleratorAreaId = helper.String(v)
			}

			request.AcceleratorAreas = append(request.AcceleratorAreas, &area)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyAccelerateAreas(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify ga2 accelerate area failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		if result.Response.TaskId != nil {
			_ = d.Set("task_id", result.Response.TaskId)
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update ga2 accelerate area failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Poll until changes are reflected
	service := NewGa2Service(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, e := service.DescribeAccelerateAreas(ctx, globalAcceleratorId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudGa2AccelerateAreaRead(d, meta)
}

func resourceTencentCloudGa2AccelerateAreaDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service             = NewGa2Service(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		globalAcceleratorId = d.Id()
	)

	// First, get all area IDs
	areas, err := service.DescribeAccelerateAreas(ctx, globalAcceleratorId)
	if err != nil {
		return err
	}

	if len(areas) == 0 {
		return nil
	}

	areaIds := make([]*string, 0, len(areas))
	for _, area := range areas {
		if area.AcceleratorAreaId != nil {
			areaIds = append(areaIds, area.AcceleratorAreaId)
		}
	}

	if len(areaIds) == 0 {
		return nil
	}

	request := ga2v20250115.NewDeleteAccelerateAreasRequest()
	request.GlobalAcceleratorId = helper.String(globalAcceleratorId)
	request.AcceleratorAreaIds = areaIds

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteAccelerateAreas(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ga2 accelerate area failed, Response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		if result.Response.TaskId != nil {
			_ = d.Set("task_id", result.Response.TaskId)
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete ga2 accelerate area failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Poll until areas are removed
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		areas, e := service.DescribeAccelerateAreas(ctx, globalAcceleratorId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if len(areas) > 0 {
			return resource.RetryableError(fmt.Errorf("waiting for ga2 accelerate areas to be deleted"))
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
