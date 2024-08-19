package cdc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdc/v20201214"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

func ResourceTencentCloudCdcSite() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTencentCloudCdcSiteCreate,
		Read:   ResourceTencentCloudCdcSiteRead,
		Update: ResourceTencentCloudCdcSiteUpdate,
		Delete: ResourceTencentCloudCdcSiteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site Name.",
			},
			"country": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site Country.",
			},
			"province": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site Province.",
			},
			"city": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site City.",
			},
			"address_line": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site Detail Address.",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Site Description.",
			},
			//"note": {
			//	Optional:    true,
			//	Type:        schema.TypeString,
			//	Description: "Site Note.",
			//},
			"fiber_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Site Fiber Type. Using optical fiber type to connect the CDC device to the network SM(Single-Mode) or MM(Multi-Mode) fibers are available.",
			},
			"optical_standard": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Site Optical Standard. Optical standard used to connect the CDC device to the network This field depends on the uplink speed, optical fiber type, and distance to upstream equipment. Allow value: `SM`, `MM`.",
			},
			"power_connectors": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Site Power Connectors. Example: 380VAC3P.",
			},
			"power_feed_drop": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Site Power Feed Drop. Whether power is supplied from above or below the rack. Allow value: `UP`, `DOWN`.",
			},
			"max_weight": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Site Max Weight capacity (KG).",
			},
			"power_draw_kva": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Site Power DrawKva (KW).",
			},
			"uplink_speed_gbps": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Uplink speed from the network to Tencent Cloud Region.",
			},
			"uplink_count": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of uplinks used by each CDC device (2 devices per rack) when connected to the network.",
			},
			"condition_requirement": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether the following environmental conditions are met: n1. There are no material requirements or the acceptance standard on site that will affect the delivery and installation of the CDC device. n2. The following conditions are met for finalized rack positions: Temperature ranges from 41 to 104 degrees F (5 to 40 degrees C). Humidity ranges from 10 degrees F (-12 degrees C) to 70 degrees F (21 degrees C) and relative humidity ranges from 8% RH to 80% RH. Air flows from front to back at the rack position and there is sufficient air in CFM (cubic feet per minute). The air quantity in CFM must be 145.8 times the power consumption (in KVA) of CDC.",
			},
			"dimension_requirement": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether the following dimension conditions are met: Your loading dock can accommodate one rack container (H x W x D = 94 x 54 x 48). You can provide a clear route from the delivery point of your rack (H x W x D = 80 x 24 x 48) to its final installation location. You should consider platforms, corridors, doors, turns, ramps, freight elevators as well as other access restrictions when measuring the depth. There shall be a 48 or greater front clearance and a 24 or greater rear clearance where the CDC is finally installed.",
			},
			"redundant_networking": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether redundant upstream equipment (switch or router) is provided so that both network devices can be connected to the network.",
			},
			//"postal_code": {
			//	Optional:    true,
			//	Type:        schema.TypeInt,
			//	Description: "Postal code of the site area.",
			//},
			"optional_address_line": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Detailed address of the site area (to be added).",
			},
			"need_help": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether you need help from Tencent Cloud for rack installation.",
			},
			"redundant_power": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether there is power redundancy.",
			},
			"breaker_requirement": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether there is an upstream circuit breaker.",
			},
		},
	}
}

func ResourceTencentCloudCdcSiteCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_site.create")()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		request  = cdc.NewCreateSiteRequest()
		response = cdc.NewCreateSiteResponse()
		siteId   string
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("country"); ok {
		request.Country = helper.String(v.(string))
	}

	if v, ok := d.GetOk("province"); ok {
		request.Province = helper.String(v.(string))
	}

	if v, ok := d.GetOk("city"); ok {
		request.City = helper.String(v.(string))
	}

	if v, ok := d.GetOk("address_line"); ok {
		request.AddressLine = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	//if v, ok := d.GetOk("note"); ok {
	//	request.Note = helper.String(v.(string))
	//}

	if v, ok := d.GetOk("fiber_type"); ok {
		request.FiberType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("optical_standard"); ok {
		request.OpticalStandard = helper.String(v.(string))
	}

	if v, ok := d.GetOk("power_connectors"); ok {
		request.PowerConnectors = helper.String(v.(string))
	}

	if v, ok := d.GetOk("power_feed_drop"); ok {
		request.PowerFeedDrop = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_weight"); ok {
		request.MaxWeight = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("power_draw_kva"); ok {
		request.PowerDrawKva = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("uplink_speed_gbps"); ok {
		request.UplinkSpeedGbps = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("uplink_count"); ok {
		request.UplinkCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("condition_requirement"); ok {
		request.ConditionRequirement = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("dimension_requirement"); ok {
		request.DimensionRequirement = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("redundant_networking"); ok {
		request.RedundantNetworking = helper.Bool(v.(bool))
	}

	//if v, ok := d.GetOkExists("postal_code"); ok {
	//	request.PostalCode = helper.IntInt64(v.(int))
	//}

	if v, ok := d.GetOk("optional_address_line"); ok {
		request.OptionalAddressLine = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("need_help"); ok {
		request.NeedHelp = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("redundant_power"); ok {
		request.RedundantPower = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("breaker_requirement"); ok {
		request.BreakerRequirement = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdcClient().CreateSite(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil {
			e = fmt.Errorf("create cdc site failed")
			return resource.NonRetryableError(e)
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cdc site failed, reason:%+v", logId, err)
		return err
	}

	siteId = *response.Response.SiteId
	d.SetId(siteId)

	return ResourceTencentCloudCdcSiteRead(d, meta)
}

func ResourceTencentCloudCdcSiteRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_site.read")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		siteId  = d.Id()
	)

	siteDetail, err := service.DescribeCdcSiteDetailById(ctx, siteId)
	if err != nil {
		return err
	}

	if siteDetail == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdcSite` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if siteDetail.Name != nil {
		_ = d.Set("name", siteDetail.Name)
	}

	if siteDetail.Country != nil {
		_ = d.Set("country", siteDetail.Country)
	}

	if siteDetail.Province != nil {
		_ = d.Set("province", siteDetail.Province)
	}

	if siteDetail.City != nil {
		_ = d.Set("city", siteDetail.City)
	}

	if siteDetail.AddressLine != nil {
		_ = d.Set("address_line", siteDetail.AddressLine)
	}

	if siteDetail.Description != nil {
		_ = d.Set("description", siteDetail.Description)
	}

	//if siteDetail.Note != nil {
	//	_ = d.Set("note", siteDetail.Note)
	//}

	if siteDetail.FiberType != nil {
		_ = d.Set("fiber_type", siteDetail.FiberType)
	}

	if siteDetail.OpticalStandard != nil {
		_ = d.Set("optical_standard", siteDetail.OpticalStandard)
	}

	if siteDetail.PowerConnectors != nil {
		_ = d.Set("power_connectors", siteDetail.PowerConnectors)
	}

	if siteDetail.PowerFeedDrop != nil {
		_ = d.Set("power_feed_drop", siteDetail.PowerFeedDrop)
	}

	if siteDetail.MaxWeight != nil {
		_ = d.Set("max_weight", siteDetail.MaxWeight)
	}

	if siteDetail.PowerDrawKva != nil {
		_ = d.Set("power_draw_kva", siteDetail.PowerDrawKva)
	}

	if siteDetail.UplinkSpeedGbps != nil {
		_ = d.Set("uplink_speed_gbps", siteDetail.UplinkSpeedGbps)
	}

	if siteDetail.UplinkCount != nil {
		_ = d.Set("uplink_count", siteDetail.UplinkCount)
	}

	if siteDetail.ConditionRequirement != nil {
		_ = d.Set("condition_requirement", siteDetail.ConditionRequirement)
	}

	if siteDetail.DimensionRequirement != nil {
		_ = d.Set("dimension_requirement", siteDetail.DimensionRequirement)
	}

	if siteDetail.RedundantNetworking != nil {
		_ = d.Set("redundant_networking", siteDetail.RedundantNetworking)
	}

	//if siteDetail.PostalCode != nil {
	//	_ = d.Set("postal_code", siteDetail.PostalCode)
	//}

	if siteDetail.OptionalAddressLine != nil {
		_ = d.Set("optional_address_line", siteDetail.OptionalAddressLine)
	}

	if siteDetail.NeedHelp != nil {
		_ = d.Set("need_help", siteDetail.NeedHelp)
	}

	if siteDetail.RedundantPower != nil {
		_ = d.Set("redundant_power", siteDetail.RedundantPower)
	}

	if siteDetail.BreakerRequirement != nil {
		_ = d.Set("breaker_requirement", siteDetail.BreakerRequirement)
	}

	return nil
}

func ResourceTencentCloudCdcSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_site.update")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = cdc.NewModifySiteInfoRequest()
		siteId  = d.Id()
	)

	immutableArgs := []string{"fiber_type", "optical_standard", "power_connectors", "power_feed_drop", "max_weight", "power_draw_kva", "uplink_speed_gbps", "uplink_count", "condition_requirement", "dimension_requirement", "redundant_networking", "optional_address_line", "need_help", "redundant_power", "breaker_requirement"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.SiteId = &siteId
	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("country") {
		if v, ok := d.GetOk("country"); ok {
			request.Country = helper.String(v.(string))
		}
	}

	if d.HasChange("province") {
		if v, ok := d.GetOk("province"); ok {
			request.Province = helper.String(v.(string))
		}
	}

	if d.HasChange("city") {
		if v, ok := d.GetOk("city"); ok {
			request.City = helper.String(v.(string))
		}
	}

	if d.HasChange("address_line") {
		if v, ok := d.GetOk("address_line"); ok {
			request.AddressLine = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	//if d.HasChange("note") {
	//	if v, ok := d.GetOk("note"); ok {
	//		request.Note = helper.String(v.(string))
	//	}
	//}

	//if d.HasChange("postal_code") {
	//	if v, ok := d.GetOkExists("postal_code"); ok {
	//		request.PostalCode = helper.String(v.(string))
	//	}
	//}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCdcClient().ModifySiteInfo(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update cdc site failed, reason:%+v", logId, err)
		return err
	}

	return ResourceTencentCloudCdcSiteRead(d, meta)
}

func ResourceTencentCloudCdcSiteDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdc_site.delete")()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		siteId  = d.Id()
	)

	if err := service.DeleteCdcSiteById(ctx, siteId); err != nil {
		return err
	}

	return nil
}
