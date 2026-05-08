package tco

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOrganizationIPWhitelistConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationIPWhitelistConfigCreate,
		Read:   resourceTencentCloudOrganizationIPWhitelistConfigRead,
		Update: resourceTencentCloudOrganizationIPWhitelistConfigUpdate,
		Delete: resourceTencentCloudOrganizationIPWhitelistConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone ID.",
			},
			"ip_whitelist": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "IP whitelist entries.",
			},
		},
	}
}

func resourceTencentCloudOrganizationIPWhitelistConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_ip_whitelist_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	d.SetId(zoneId)

	return resourceTencentCloudOrganizationIPWhitelistConfigUpdate(d, meta)
}

func resourceTencentCloudOrganizationIPWhitelistConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_ip_whitelist_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	zoneId := d.Id()
	if zoneId == "" {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	ipWhitelist, err := service.DescribeOrganizationIPWhitelistConfigById(ctx, zoneId)
	if err != nil {
		return err
	}

	if ipWhitelist == nil {
		log.Printf("[WARN]%s resource `tencentcloud_organization_ip_whitelist_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("ip_whitelist", ipWhitelist)

	return nil
}

func resourceTencentCloudOrganizationIPWhitelistConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_ip_whitelist_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	zoneId := d.Id()

	request := organization.NewUpdateIPWhitelistRequest()
	request.ZoneId = &zoneId

	if v, ok := d.GetOk("ip_whitelist"); ok {
		ipWhitelist := v.([]interface{})
		request.IpWhitelist = make([]*string, 0, len(ipWhitelist))
		for _, item := range ipWhitelist {
			request.IpWhitelist = append(request.IpWhitelist, helper.String(item.(string)))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().UpdateIPWhitelistWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update organization IP whitelist config failed, reason: %+v", logId, err)
		return err
	}

	return resourceTencentCloudOrganizationIPWhitelistConfigRead(d, meta)
}

func resourceTencentCloudOrganizationIPWhitelistConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_ip_whitelist_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	zoneId := d.Id()
	if zoneId == "" {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	request := organization.NewUpdateIPWhitelistRequest()
	request.ZoneId = &zoneId
	request.IpWhitelist = []*string{}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().UpdateIPWhitelistWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete organization IP whitelist config failed, reason: %+v", logId, err)
		return err
	}

	return nil
}
