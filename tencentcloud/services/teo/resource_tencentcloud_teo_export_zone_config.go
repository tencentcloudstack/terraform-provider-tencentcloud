package teo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoExportZoneConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoExportZoneConfigCreate,
		Read:   resourceTencentCloudTeoExportZoneConfigRead,
		Update: resourceTencentCloudTeoExportZoneConfigUpdate,
		Delete: resourceTencentCloudTeoExportZoneConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the site.",
			},
			"types": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of configuration types to export. If not specified, all configuration types will be exported. Currently supported types: L7AccelerationConfig (seven-layer acceleration configuration).",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The exported configuration content in JSON format, encoded in UTF-8.",
			},
		},
	}
}

func resourceTencentCloudTeoExportZoneConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(
		context.WithValue(context.Background(), tccommon.LogIdKey, logId), logId, d, meta)

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	// Build request
	request := teo.NewExportZoneConfigRequest()
	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("types"); ok {
		types := helper.InterfacesStrings(v.([]interface{}))
		request.Types = helper.Strings(types)
	}

	// Call API
	response := teo.NewExportZoneConfigResponse()
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ExportZoneConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "ExportZoneConfig", request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "ExportZoneConfig", request.ToJsonString(), response.ToJsonString())

	// Set resource ID and attributes
	d.SetId(zoneId)
	_ = d.Set("zone_id", zoneId)

	if response.Response != nil && response.Response.Content != nil {
		_ = d.Set("content", *response.Response.Content)
	}

	// Cache content in state for Read operation
	_ = d.Set("types", d.Get("types"))

	return resourceTencentCloudTeoExportZoneConfigRead(ctx, d, meta)
}

func resourceTencentCloudTeoExportZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(
		context.WithValue(context.Background(), tccommon.LogIdKey, logId), logId, d, meta)

	zoneId := d.Id()

	// For export zone config, we just return the cached content from state
	// Since this is an export operation, we don't need to call the API again
	log.Printf("[DEBUG]%s read export zone config, zone_id [%s]\n", logId, zoneId)

	return nil
}

func resourceTencentCloudTeoExportZoneConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(
		context.WithValue(context.Background(), tccommon.LogIdKey, logId), logId, d, meta)

	zoneId := d.Id()

	// Build request
	request := teo.NewExportZoneConfigRequest()
	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("types"); ok {
		types := helper.InterfacesStrings(v.([]interface{}))
		request.Types = helper.Strings(types)
	}

	// Call API
	response := teo.NewExportZoneConfigResponse()
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ExportZoneConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, "ExportZoneConfig", request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, "ExportZoneConfig", request.ToJsonString(), response.ToJsonString())

	// Update attributes
	_ = d.Set("zone_id", zoneId)

	if response.Response != nil && response.Response.Content != nil {
		_ = d.Set("content", *response.Response.Content)
	}

	return resourceTencentCloudTeoExportZoneConfigRead(ctx, d, meta)
}

func resourceTencentCloudTeoExportZoneConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	// This is an export resource, so we don't need to delete anything from the API
	// Just remove from Terraform state
	log.Printf("[DEBUG]%s delete export zone config, zone_id [%s]\n", logId, d.Id())

	return nil
}
