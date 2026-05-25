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
				Description: "Global accelerator instance ID.",
			},
			"accelerator_areas": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Accelerate area configurations.",
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
				Description: "Accelerate area set returned from the API.",
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
							Description: "IP address list.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ip_address_info_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IP address info set.",
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
										Description: "ISP type.",
									},
								},
							},
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

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := NewGa2Service(client)

	globalAcceleratorId := d.Get("global_accelerator_id").(string)
	areas := buildAcceleratorAreas(d)

	taskId, err := service.CreateAccelerateAreas(ctx, globalAcceleratorId, areas)
	if err != nil {
		return err
	}

	// Poll task result until success
	err = waitForTaskResult(ctx, &service, taskId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	d.SetId(globalAcceleratorId)

	return resourceTencentCloudGa2AccelerateAreaRead(d, meta)
}

func resourceTencentCloudGa2AccelerateAreaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := NewGa2Service(client)

	globalAcceleratorId := d.Id()

	areas, err := service.DescribeAccelerateAreas(ctx, globalAcceleratorId)
	if err != nil {
		return err
	}

	if len(areas) == 0 {
		d.SetId("")
		return nil
	}

	_ = d.Set("global_accelerator_id", globalAcceleratorId)

	accelerateAreaSet := make([]map[string]interface{}, 0, len(areas))
	for _, area := range areas {
		areaMap := map[string]interface{}{}
		if area.AcceleratorAreaId != nil {
			areaMap["accelerator_area_id"] = *area.AcceleratorAreaId
		}
		if area.AccelerateRegion != nil {
			areaMap["accelerate_region"] = *area.AccelerateRegion
		}
		if area.Bandwidth != nil {
			areaMap["bandwidth"] = int(*area.Bandwidth)
		}
		if area.IspType != nil {
			areaMap["isp_type"] = *area.IspType
		}
		if area.IpVersion != nil {
			areaMap["ip_version"] = *area.IpVersion
		}
		if area.IpAddress != nil {
			ipAddresses := make([]string, 0, len(area.IpAddress))
			for _, ip := range area.IpAddress {
				if ip != nil {
					ipAddresses = append(ipAddresses, *ip)
				}
			}
			areaMap["ip_address"] = ipAddresses
		}
		if area.IpAddressInfoSet != nil {
			ipInfoSet := make([]map[string]interface{}, 0, len(area.IpAddressInfoSet))
			for _, ipInfo := range area.IpAddressInfoSet {
				ipInfoMap := map[string]interface{}{}
				if ipInfo.IpAddress != nil {
					ipInfoMap["ip_address"] = *ipInfo.IpAddress
				}
				if ipInfo.IspType != nil {
					ipInfoMap["isp_type"] = *ipInfo.IspType
				}
				ipInfoSet = append(ipInfoSet, ipInfoMap)
			}
			areaMap["ip_address_info_set"] = ipInfoSet
		}
		accelerateAreaSet = append(accelerateAreaSet, areaMap)
	}

	_ = d.Set("accelerate_area_set", accelerateAreaSet)

	return nil
}

func resourceTencentCloudGa2AccelerateAreaUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := NewGa2Service(client)

	globalAcceleratorId := d.Id()

	if d.HasChange("accelerator_areas") {
		areas := buildAcceleratorAreasForModify(d)

		taskId, err := service.ModifyAccelerateAreas(ctx, globalAcceleratorId, areas)
		if err != nil {
			return err
		}

		err = waitForTaskResult(ctx, &service, taskId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudGa2AccelerateAreaRead(d, meta)
}

func resourceTencentCloudGa2AccelerateAreaDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := NewGa2Service(client)

	globalAcceleratorId := d.Id()

	// First read current areas to get their IDs
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

	taskId, err := service.DeleteAccelerateAreas(ctx, globalAcceleratorId, areaIds)
	if err != nil {
		return err
	}

	err = waitForTaskResult(ctx, &service, taskId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return err
	}

	return nil
}

func buildAcceleratorAreas(d *schema.ResourceData) []*ga2v20250115.AcceleratorAreas {
	rawAreas := d.Get("accelerator_areas").([]interface{})
	areas := make([]*ga2v20250115.AcceleratorAreas, 0, len(rawAreas))
	for _, rawArea := range rawAreas {
		areaMap := rawArea.(map[string]interface{})
		area := &ga2v20250115.AcceleratorAreas{}
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
		areas = append(areas, area)
	}
	return areas
}

func buildAcceleratorAreasForModify(d *schema.ResourceData) []*ga2v20250115.AcceleratorAreas {
	rawAreas := d.Get("accelerator_areas").([]interface{})
	areas := make([]*ga2v20250115.AcceleratorAreas, 0, len(rawAreas))

	// Get existing area IDs from accelerate_area_set to match with new config
	existingAreas := d.Get("accelerate_area_set").([]interface{})
	existingAreaMap := make(map[string]string) // region -> area_id
	for _, existing := range existingAreas {
		if existingMap, ok := existing.(map[string]interface{}); ok {
			region := ""
			areaId := ""
			if v, ok := existingMap["accelerate_region"].(string); ok {
				region = v
			}
			if v, ok := existingMap["accelerator_area_id"].(string); ok {
				areaId = v
			}
			if region != "" && areaId != "" {
				existingAreaMap[region] = areaId
			}
		}
	}

	for _, rawArea := range rawAreas {
		areaMap := rawArea.(map[string]interface{})
		area := &ga2v20250115.AcceleratorAreas{}
		if v, ok := areaMap["accelerate_region"].(string); ok && v != "" {
			area.AccelerateRegion = helper.String(v)
			// Try to match with existing area ID
			if areaId, exists := existingAreaMap[v]; exists {
				area.AcceleratorAreaId = helper.String(areaId)
			}
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
		areas = append(areas, area)
	}
	return areas
}

func waitForTaskResult(ctx context.Context, service *Ga2Service, taskId string, timeout time.Duration) error {
	logId := tccommon.GetLogId(ctx)

	err := resource.Retry(timeout, func() *resource.RetryError {
		status, e := service.DescribeTaskResult(ctx, taskId)
		if e != nil {
			return resource.NonRetryableError(e)
		}

		if status == "SUCCESS" {
			return nil
		}

		if status == "FAILED" {
			return resource.NonRetryableError(fmt.Errorf("task %s failed", taskId))
		}

		log.Printf("[DEBUG]%s task %s status: %s, retrying...\n", logId, taskId, status)
		return resource.RetryableError(fmt.Errorf("task %s is still in progress, status: %s", taskId, status))
	})

	return err
}
