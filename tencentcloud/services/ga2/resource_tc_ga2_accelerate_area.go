package ga2

import (
	"context"
	"fmt"
	"log"
	"time"

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
				Description: "List of accelerate area configurations for Create/Update.",
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
							Description: "ISP type, supports `BGP`, `三网`, `精品`. Default is `BGP`.",
						},
						"ip_version": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "IPv4",
							Description: "IP version, only supports `IPv4`. Default is `IPv4`.",
						},
					},
				},
			},
			"accelerate_area_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of accelerate area details returned by Read.",
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
							Description: "The accelerate area ID.",
						},
						"ip_address": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IP addresses.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ip_address_info_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IP address info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_address": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP address.",
									},
									"isp_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP type.",
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
				Description: "The last async task ID.",
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
		service             = NewGa2Service(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		globalAcceleratorId string
	)

	globalAcceleratorId = d.Get("global_accelerator_id").(string)

	areas := make([]*ga2v20250115.AcceleratorAreas, 0)
	if v, ok := d.GetOk("accelerator_areas"); ok {
		for _, item := range v.([]interface{}) {
			areaMap := item.(map[string]interface{})
			area := &ga2v20250115.AcceleratorAreas{}
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

			areas = append(areas, area)
		}
	}

	taskId, err := service.CreateAccelerateAreas(ctx, globalAcceleratorId, areas)
	if err != nil {
		log.Printf("[CRITAL]%s create ga2 accelerate area failed, reason:%+v", logId, err)
		return err
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	err = service.DescribeTaskResult(ctx, taskId, timeout)
	if err != nil {
		return fmt.Errorf("waiting for task %s to complete failed: %s", taskId, err.Error())
	}

	d.SetId(globalAcceleratorId)
	_ = d.Set("task_id", taskId)

	return resourceTencentCloudGa2AccelerateAreaRead(d, meta)
}

func resourceTencentCloudGa2AccelerateAreaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewGa2Service(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	globalAcceleratorId := d.Id()

	areaSet, err := service.DescribeAccelerateAreas(ctx, globalAcceleratorId)
	if err != nil {
		return err
	}

	if len(areaSet) == 0 {
		log.Printf("[WARN]%s resource `tencentcloud_ga2_accelerate_area` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("global_accelerator_id", globalAcceleratorId)

	accelerateAreaSetList := make([]map[string]interface{}, 0, len(areaSet))
	for _, area := range areaSet {
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
			areaMap["ip_address"] = area.IpAddress
		}

		if area.IpAddressInfoSet != nil {
			ipAddressInfoList := make([]map[string]interface{}, 0, len(area.IpAddressInfoSet))
			for _, ipInfo := range area.IpAddressInfoSet {
				ipInfoMap := map[string]interface{}{}
				if ipInfo.IpAddress != nil {
					ipInfoMap["ip_address"] = ipInfo.IpAddress
				}

				if ipInfo.IspType != nil {
					ipInfoMap["isp_type"] = ipInfo.IspType
				}

				ipAddressInfoList = append(ipAddressInfoList, ipInfoMap)
			}

			areaMap["ip_address_info_set"] = ipAddressInfoList
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
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewGa2Service(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	globalAcceleratorId := d.Id()

	if d.HasChange("accelerator_areas") {
		// First read current areas to get AcceleratorAreaIds
		currentAreas, err := service.DescribeAccelerateAreas(ctx, globalAcceleratorId)
		if err != nil {
			return err
		}

		// Build a map from region to area ID for matching
		regionToAreaId := make(map[string]string)
		for _, area := range currentAreas {
			if area.AccelerateRegion != nil && area.AcceleratorAreaId != nil {
				regionToAreaId[*area.AccelerateRegion] = *area.AcceleratorAreaId
			}
		}

		areas := make([]*ga2v20250115.AcceleratorAreas, 0)
		if v, ok := d.GetOk("accelerator_areas"); ok {
			for _, item := range v.([]interface{}) {
				areaMap := item.(map[string]interface{})
				area := &ga2v20250115.AcceleratorAreas{}
				region := ""
				if v, ok := areaMap["accelerate_region"].(string); ok && v != "" {
					area.AccelerateRegion = helper.String(v)
					region = v
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

				// Set AcceleratorAreaId if the region already exists
				if areaId, exists := regionToAreaId[region]; exists {
					area.AcceleratorAreaId = helper.String(areaId)
				}

				areas = append(areas, area)
			}
		}

		taskId, err := service.ModifyAccelerateAreas(ctx, globalAcceleratorId, areas)
		if err != nil {
			log.Printf("[CRITAL]%s modify ga2 accelerate area failed, reason:%+v", logId, err)
			return err
		}

		timeout := d.Timeout(schema.TimeoutUpdate)
		err = service.DescribeTaskResult(ctx, taskId, timeout)
		if err != nil {
			return fmt.Errorf("waiting for task %s to complete failed: %s", taskId, err.Error())
		}

		_ = d.Set("task_id", taskId)
	}

	return resourceTencentCloudGa2AccelerateAreaRead(d, meta)
}

func resourceTencentCloudGa2AccelerateAreaDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewGa2Service(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	globalAcceleratorId := d.Id()

	// First read current areas to get all AcceleratorAreaIds
	currentAreas, err := service.DescribeAccelerateAreas(ctx, globalAcceleratorId)
	if err != nil {
		return err
	}

	if len(currentAreas) == 0 {
		return nil
	}

	areaIds := make([]*string, 0, len(currentAreas))
	for _, area := range currentAreas {
		if area.AcceleratorAreaId != nil {
			areaIds = append(areaIds, area.AcceleratorAreaId)
		}
	}

	if len(areaIds) == 0 {
		return nil
	}

	taskId, err := service.DeleteAccelerateAreas(ctx, globalAcceleratorId, areaIds)
	if err != nil {
		log.Printf("[CRITAL]%s delete ga2 accelerate area failed, reason:%+v", logId, err)
		return err
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	err = service.DescribeTaskResult(ctx, taskId, timeout)
	if err != nil {
		return fmt.Errorf("waiting for task %s to complete failed: %s", taskId, err.Error())
	}

	return nil
}
