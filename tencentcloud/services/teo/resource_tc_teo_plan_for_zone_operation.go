package teo

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoPlanForZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoPlanForZoneCreate,
		Read:   resourceTencentCloudTeoPlanForZoneRead,
		Delete: resourceTencentCloudTeoPlanForZoneDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"plan_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Plan type to purchase. Valid values: `sta`, `sta_with_bot`, `sta_cm`, `sta_cm_with_bot`, `sta_global`, `sta_global_with_bot`, `ent`, `ent_with_bot`, `ent_cm`, `ent_cm_with_bot`, `ent_global`, `ent_global_with_bot`.",
			},
			"resource_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of purchased resource names returned by the API.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"deal_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of purchased order/deal names returned by the API.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudTeoPlanForZoneCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_plan_for_zone.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	zoneId := d.Get("zone_id").(string)
	if zoneId == "" {
		return fmt.Errorf("zone_id is required")
	}

	planType := d.Get("plan_type").(string)
	if planType == "" {
		return fmt.Errorf("plan_type is required")
	}

	request := teov20220901.NewCreatePlanForZoneRequest()
	request.ZoneId = helper.String(zoneId)
	request.PlanType = helper.String(planType)

	// Get Teo service
	service := NewTeoService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	// Call CreatePlanForZone API
	log.Printf("[DEBUG]%s api[%s] request body [%s]\n", logId, request.GetAction(), request.ToJsonString())

	var resourceNames, dealNames []*string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var e error
		resourceNames, dealNames, e = service.TeoPlanForZone(zoneId, planType)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if resourceNames == nil && dealNames == nil {
			return resource.NonRetryableError(fmt.Errorf("teo plan_for_zone purchase returned empty result, zone_id=%s", zoneId))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s teo plan_for_zone create failed, id=%s, reason[%s]\n", logId, zoneId, err.Error())
		return fmt.Errorf("CreatePlanForZone failed: %s", err.Error())
	}

	log.Printf("[DEBUG]%s api[%s] success, resource_names: %+v, deal_names: %+v\n",
		logId, request.GetAction(), resourceNames, dealNames)

	// Set ID
	d.SetId(helper.BuildToken())

	// Set resource_names
	if err := d.Set("resource_names", flattenTeoStringPtrList(resourceNames)); err != nil {
		return err
	}

	// Set deal_names
	if err := d.Set("deal_names", flattenTeoStringPtrList(dealNames)); err != nil {
		return err
	}

	return resourceTencentCloudTeoPlanForZoneRead(d, meta)
}

func resourceTencentCloudTeoPlanForZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_plan_for_zone.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoPlanForZoneDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_plan_for_zone.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

// flattenTeoStringPtrList converts a []*string into a []interface{} of strings.
func flattenTeoStringPtrList(list []*string) []interface{} {
	result := make([]interface{}, 0, len(list))
	for _, v := range list {
		if v == nil {
			result = append(result, "")
		} else {
			result = append(result, *v)
		}
	}
	return result
}
