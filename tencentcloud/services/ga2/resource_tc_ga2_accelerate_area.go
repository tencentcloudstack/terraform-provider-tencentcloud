package ga2

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// ResourceTencentCloudGa2AccelerateArea manages a single Tencent Cloud GA2 acceleration region.
//
// All write APIs (CreateAccelerateAreas / ModifyAccelerateAreas / DeleteAccelerateAreas) are
// asynchronous: each returns a TaskId that must be polled via DescribeTaskResult until
// Status == "SUCCESS". CreateAccelerateAreas does NOT return the generated AcceleratorAreaId, so
// after the create task succeeds the resource resolves it via DescribeAccelerateAreas keyed on the
// natural key (GlobalAcceleratorId, AccelerateRegion).
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
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_accelerator_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Global accelerator instance ID this acceleration region belongs to.",
			},
			"accelerate_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Acceleration region. Serves as the natural key used to resolve the acceleration region ID " +
					"after creation. Cannot be modified after creation; modifying it forces a new resource.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Acceleration bandwidth in Mbps.",
			},
			"isp_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "ISP type. Valid values: `BGP` (BGP), `STATIC_IP` (multi-ISP static IP), `QUALITY_BGP` (premium BGP). " +
					"Default: `BGP`.",
			},
			"ip_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "IP version. Only `IPv4` is supported. Default: `IPv4`.",
			},
			"ip_address": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Bound IP address list. Treated as an unordered set; HCL element order has no semantic meaning.",
			},

			// Computed
			"accelerator_area_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Acceleration region ID.",
			},
			"ip_address_info_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IP address information list.",
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
							Description: "ISP type of the IP address.",
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
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = ga2v20250115.NewCreateAccelerateAreasRequest()
		response = ga2v20250115.NewCreateAccelerateAreasResponse()
		gaId     string
		region   string
		taskId   string
	)

	if v, ok := d.GetOk("global_accelerator_id"); ok {
		gaId = v.(string)
		request.GlobalAcceleratorId = helper.String(gaId)
	}

	if v, ok := d.GetOk("accelerate_region"); ok {
		region = v.(string)
	}

	// The resource manages a single acceleration region, so the request carries a one-element list.
	request.AcceleratorAreas = []*ga2v20250115.AcceleratorAreas{buildGa2AcceleratorArea(d, "", "create")}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateAccelerateAreasWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ga2 accelerate areas failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create ga2 accelerate areas failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.TaskId == nil || *response.Response.TaskId == "" {
		return fmt.Errorf("TaskId is nil.")
	}
	taskId = *response.Response.TaskId

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	// CreateAccelerateAreas does not return AcceleratorAreaId, so resolve it from the natural key.
	area, err := service.DescribeGa2AccelerateAreaByRegion(ctx, gaId, region)
	if err != nil {
		return err
	}

	if area == nil || area.AcceleratorAreaId == nil || *area.AcceleratorAreaId == "" {
		return fmt.Errorf("Create ga2 accelerate area succeeded but acceleration region [%s] under accelerator [%s] could not be resolved.", region, gaId)
	}

	d.SetId(strings.Join([]string{gaId, *area.AcceleratorAreaId}, tccommon.FILED_SP))
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

	gaId, areaId, err := parseGa2AccelerateAreaId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeGa2AccelerateAreaById(ctx, gaId, areaId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_ga2_accelerate_area` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	// global_accelerator_id is the first segment of the composite ID and is not echoed by the API item.
	_ = d.Set("global_accelerator_id", gaId)

	if respData.AcceleratorAreaId != nil {
		_ = d.Set("accelerator_area_id", respData.AcceleratorAreaId)
	}

	if respData.AccelerateRegion != nil {
		_ = d.Set("accelerate_region", respData.AccelerateRegion)
	}

	if respData.Bandwidth != nil {
		_ = d.Set("bandwidth", int(*respData.Bandwidth))
	}

	if respData.IspType != nil {
		_ = d.Set("isp_type", respData.IspType)
	}

	if respData.IpVersion != nil {
		_ = d.Set("ip_version", respData.IpVersion)
	}

	if len(respData.IpAddress) > 0 {
		_ = d.Set("ip_address", helper.PStrings(respData.IpAddress))
	}

	if len(respData.IpAddressInfoSet) > 0 {
		_ = d.Set("ip_address_info_set", flattenGa2IpAddressInfoSet(respData.IpAddressInfoSet))
	}

	return nil
}

func resourceTencentCloudGa2AccelerateAreaUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	gaId, _, err := parseGa2AccelerateAreaId(d.Id())
	if err != nil {
		return err
	}

	// ModifyAccelerateAreas accepts the mutable area fields; accelerate_region is ForceNew.
	modifyFields := []string{"bandwidth", "isp_type", "ip_version", "ip_address"}
	needModify := false
	for _, f := range modifyFields {
		if d.HasChange(f) {
			needModify = true
			break
		}
	}

	if !needModify {
		return resourceTencentCloudGa2AccelerateAreaRead(d, meta)
	}

	request := ga2v20250115.NewModifyAccelerateAreasRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.AcceleratorAreas = []*ga2v20250115.AcceleratorAreas{buildGa2AcceleratorArea(d, "", "update")}

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyAccelerateAreasWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify ga2 accelerate areas failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update ga2 accelerate areas failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return err
	}

	return resourceTencentCloudGa2AccelerateAreaRead(d, meta)
}

func resourceTencentCloudGa2AccelerateAreaDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_accelerate_area.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = ga2v20250115.NewDeleteAccelerateAreasRequest()
	)

	gaId, areaId, err := parseGa2AccelerateAreaId(d.Id())
	if err != nil {
		return err
	}

	request.GlobalAcceleratorId = helper.String(gaId)
	request.AcceleratorAreaIds = []*string{helper.String(areaId)}

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteAccelerateAreasWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ga2 accelerate areas failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete ga2 accelerate areas failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return err
	}

	return nil
}

// parseGa2AccelerateAreaId splits the composite resource ID into its two components.
func parseGa2AccelerateAreaId(id string) (gaId, areaId string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid resource id %q, expected format <global_accelerator_id>%s<accelerator_area_id>", id, tccommon.FILED_SP)
		return
	}
	gaId, areaId = parts[0], parts[1]
	if gaId == "" || areaId == "" {
		err = fmt.Errorf("invalid resource id %q, components must all be non-empty", id)
		return
	}
	return
}

// buildGa2AcceleratorArea assembles the single AcceleratorAreas element from the schema. When areaId
// is non-empty (Update path) it is set so the API can correlate the entry being modified.
func buildGa2AcceleratorArea(d *schema.ResourceData, areaId, step string) *ga2v20250115.AcceleratorAreas {
	area := &ga2v20250115.AcceleratorAreas{}
	if step == "create" {
		if areaId != "" {
			area.AcceleratorAreaId = helper.String(areaId)
		}

		if v, ok := d.GetOk("accelerate_region"); ok {
			area.AccelerateRegion = helper.String(v.(string))
		}

		if v, ok := d.GetOk("bandwidth"); ok {
			area.Bandwidth = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("isp_type"); ok {
			area.IspType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("ip_version"); ok {
			area.IpVersion = helper.String(v.(string))
		}

		if v, ok := d.GetOk("ip_address"); ok {
			area.IpAddress = buildGa2AccelerateAreaStringSet(v.(*schema.Set))
		}
	} else if step == "update" {
		if v, ok := d.GetOk("accelerate_region"); ok {
			area.AccelerateRegion = helper.String(v.(string))
		}

		if v, ok := d.GetOk("bandwidth"); ok {
			area.Bandwidth = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("isp_type"); ok {
			area.IspType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("ip_version"); ok {
			area.IpVersion = helper.String(v.(string))
		}
	}

	return area
}

// buildGa2AccelerateAreaStringSet converts a TypeSet of strings into a []*string suitable for the SDK.
func buildGa2AccelerateAreaStringSet(set *schema.Set) []*string {
	if set == nil || set.Len() == 0 {
		return nil
	}
	result := make([]*string, 0, set.Len())
	for _, item := range set.List() {
		s, ok := item.(string)
		if !ok || s == "" {
			continue
		}
		v := s
		result = append(result, &v)
	}
	return result
}

// flattenGa2IpAddressInfoSet maps the SDK IpAddressInfoSet slice into the computed nested block payload.
func flattenGa2IpAddressInfoSet(infoSet []*ga2v20250115.IpAddressInfoSet) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(infoSet))
	for _, info := range infoSet {
		if info == nil {
			continue
		}
		m := map[string]interface{}{}
		if info.IpAddress != nil {
			m["ip_address"] = *info.IpAddress
		}
		if info.IspType != nil {
			m["isp_type"] = *info.IspType
		}
		result = append(result, m)
	}
	return result
}
