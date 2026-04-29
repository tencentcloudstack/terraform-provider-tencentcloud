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

func ResourceTencentCloudTeoPrefetchOriginLimit() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoPrefetchOriginLimitConfigCreate,
		Read:   resourceTencentCloudTeoPrefetchOriginLimitConfigRead,
		Update: resourceTencentCloudTeoPrefetchOriginLimitConfigUpdate,
		Delete: resourceTencentCloudTeoPrefetchOriginLimitConfigDelete,
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
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Accelerated domain name.",
			},
			"area": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Acceleration area for prefetch origin limit. Valid values: `Overseas`, `MainlandChina`.",
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Prefetch origin bandwidth limit. Value range: 100-100000, in Mbps.",
			},
			"enabled": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Prefetch origin limit switch. Valid values: `on`, `off`.",
			},
		},
	}
}

func resourceTencentCloudTeoPrefetchOriginLimitConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_prefetch_origin_limit.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		zoneId     = d.Get("zone_id").(string)
		domainName = d.Get("domain_name").(string)
		area       = d.Get("area").(string)
		bandwidth  = d.Get("bandwidth").(int)
		enabled    = d.Get("enabled").(string)
		request    = teov20220901.NewModifyPrefetchOriginLimitRequest()
	)

	request.ZoneId = &zoneId
	request.DomainName = &domainName
	request.Area = &area
	request.Bandwidth = helper.Int64(int64(bandwidth))
	request.Enabled = &enabled

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyPrefetchOriginLimitWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo prefetch origin limit failed, reason:%+v", logId, err)
		return err
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("domain_name", domainName)
	_ = d.Set("area", area)

	d.SetId(strings.Join([]string{zoneId, domainName, area}, tccommon.FILED_SP))

	// Read to confirm creation
	respData, err := service.DescribeTeoPrefetchOriginLimitById(ctx, zoneId, domainName, area)
	if err != nil {
		return err
	}
	if respData == nil {
		return fmt.Errorf("teo prefetch origin limit not found after create")
	}

	return resourceTencentCloudTeoPrefetchOriginLimitConfigRead(d, meta)
}

func resourceTencentCloudTeoPrefetchOriginLimitConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_prefetch_origin_limit.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		id      = d.Id()
		idParts = strings.Split(id, tccommon.FILED_SP)
	)

	if len(idParts) != 3 {
		return fmt.Errorf("resource ID format error, expected zone_id%sdomain_name%sarea", tccommon.FILED_SP, tccommon.FILED_SP)
	}

	zoneId := idParts[0]
	domainName := idParts[1]
	area := idParts[2]

	respData, err := service.DescribeTeoPrefetchOriginLimitById(ctx, zoneId, domainName, area)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_prefetch_origin_limit` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if respData.DomainName != nil {
		_ = d.Set("domain_name", respData.DomainName)
	}
	if respData.Area != nil {
		_ = d.Set("area", respData.Area)
	}
	if respData.Bandwidth != nil {
		_ = d.Set("bandwidth", int(*respData.Bandwidth))
	}

	return nil
}

func resourceTencentCloudTeoPrefetchOriginLimitConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_prefetch_origin_limit.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		id      = d.Id()
		idParts = strings.Split(id, tccommon.FILED_SP)
	)

	if len(idParts) != 3 {
		return fmt.Errorf("resource ID format error, expected zone_id%sdomain_name%sarea", tccommon.FILED_SP, tccommon.FILED_SP)
	}

	zoneId := idParts[0]
	domainName := idParts[1]
	area := idParts[2]

	if d.HasChange("bandwidth") || d.HasChange("enabled") {
		request := teov20220901.NewModifyPrefetchOriginLimitRequest()
		request.ZoneId = &zoneId
		request.DomainName = &domainName
		request.Area = &area
		bandwidth := d.Get("bandwidth").(int)
		enabled := d.Get("enabled").(string)
		request.Bandwidth = helper.Int64(int64(bandwidth))
		request.Enabled = &enabled

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyPrefetchOriginLimitWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo prefetch origin limit failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoPrefetchOriginLimitConfigRead(d, meta)
}

func resourceTencentCloudTeoPrefetchOriginLimitConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_prefetch_origin_limit.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		id      = d.Id()
		idParts = strings.Split(id, tccommon.FILED_SP)
	)

	if len(idParts) != 3 {
		return fmt.Errorf("resource ID format error, expected zone_id%sdomain_name%sarea", tccommon.FILED_SP, tccommon.FILED_SP)
	}

	zoneId := idParts[0]
	domainName := idParts[1]
	area := idParts[2]

	request := teov20220901.NewModifyPrefetchOriginLimitRequest()
	request.ZoneId = &zoneId
	request.DomainName = &domainName
	request.Area = &area
	bandwidth := d.Get("bandwidth").(int)
	request.Bandwidth = helper.Int64(int64(bandwidth))
	enabled := "off"
	request.Enabled = &enabled

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyPrefetchOriginLimitWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo prefetch origin limit failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
