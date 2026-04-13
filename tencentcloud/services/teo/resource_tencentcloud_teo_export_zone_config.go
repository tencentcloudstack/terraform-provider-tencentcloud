package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var exportTimeout time.Duration

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
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the site.",
			},
			"export_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"all", "basic", "cache", "https", "origin", "waf", "rate_limit", "rule_engine"}),
				Description:  "Export type. Valid values: `all`: all configurations; `basic`: basic configurations; `cache`: cache configurations; `https`: HTTPS configurations; `origin`: origin configurations; `waf`: WAF configurations; `rate_limit`: rate limit configurations; `rule_engine`: rule engine configurations.",
			},
			"config_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Exported configuration content in JSON format.",
			},
		},
	}
}

func resourceTencentCloudTeoExportZoneConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.WithValue(context.Background(), exportTimeout, d.Timeout(schema.TimeoutCreate)), logId, d, meta)
		zoneId     string
		exportType string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("export_type"); ok {
		exportType = v.(string)
	}

	request := teo.NewDescribeZoneConfigRequest()
	request.ZoneId = helper.String(zoneId)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DescribeZoneConfigWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Export teo zone config failed, Response is nil."))
		}

		// Set the config content
		_ = d.Set("config_content", result.ToJsonString())

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s export teo zone config failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{zoneId, exportType}, tccommon.FILED_SP))

	return resourceTencentCloudTeoExportZoneConfigRead(d, meta)
}

func resourceTencentCloudTeoExportZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.WithValue(context.Background(), exportTimeout, d.Timeout(schema.TimeoutRead)), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	exportType := idSplit[1]

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("export_type", exportType)

	request := teo.NewDescribeZoneConfigRequest()
	request.ZoneId = helper.String(zoneId)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DescribeZoneConfigWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe teo zone config failed, Response is nil."))
		}

		// Set the config content
		_ = d.Set("config_content", result.ToJsonString())

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read teo zone config failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func resourceTencentCloudTeoExportZoneConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(
			context.WithValue(context.Background(), exportTimeout, d.Timeout(schema.TimeoutUpdate)), logId, d, meta)
	)

	// For export resource, we just need to refresh the config content
	if d.HasChange("zone_id") {
		return fmt.Errorf("zone_id cannot be changed for export resource")
	}

	if d.HasChange("export_type") {
		idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
		zoneId := idSplit[0]
		exportType := d.Get("export_type").(string)
		d.SetId(strings.Join([]string{zoneId, exportType}, tccommon.FILED_SP))
	}

	return resourceTencentCloudTeoExportZoneConfigRead(d, meta)
}

func resourceTencentCloudTeoExportZoneConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_export_zone_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	log.Printf("[INFO]%s delete teo export zone config, resource id [%s]\n", logId, d.Id())

	// Export resource does not actually delete anything from the cloud
	// It just removes the resource from Terraform state
	d.SetId("")

	return nil
}
