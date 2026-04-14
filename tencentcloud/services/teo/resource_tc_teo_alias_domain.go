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

func ResourceTencentCloudTeoAliasDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoAliasDomainCreate,
		Read:   resourceTencentCloudTeoAliasDomainRead,
		Update: resourceTencentCloudTeoAliasDomainUpdate,
		Delete: resourceTencentCloudTeoAliasDomainDelete,
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

			"alias_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Alias domain name.",
			},

			"target_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target domain name.",
			},

			"paused": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the alias domain is paused. true: paused, false: active.",
			},
		},
	}
}

func resourceTencentCloudTeoAliasDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_alias_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = teo.NewCreateAliasDomainRequest()
		zoneId    string
		aliasName string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("alias_name"); ok {
		request.AliasName = helper.String(v.(string))
		aliasName = v.(string)
	}

	if v, ok := d.GetOk("target_name"); ok {
		request.TargetName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateAliasDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo alias domain failed, reason:%+v", logId, err)
		return err
	}

	// wait for creation to complete
	if err := resourceTencentCloudTeoAliasDomainCreatePostHandle(ctx, zoneId, aliasName, meta, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{zoneId, aliasName}, tccommon.FILED_SP))

	// set paused status if specified
	if v, ok := d.GetOkExists("paused"); ok && v.(bool) {
		request := teo.NewModifyAliasDomainStatusRequest()
		request.ZoneId = helper.String(zoneId)
		request.AliasNames = []*string{helper.String(aliasName)}
		request.Paused = helper.Bool(v.(bool))

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyAliasDomainStatusWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update teo alias domain paused status failed, reason:%+v", logId, err)
			return err
		}

		// wait for status update to complete
		if err := resourceTencentCloudTeoAliasDomainUpdatePostHandle(ctx, zoneId, aliasName, v.(bool), meta, d.Timeout(schema.TimeoutCreate)); err != nil {
			return err
		}
	}

	return resourceTencentCloudTeoAliasDomainRead(d, meta)
}

func resourceTencentCloudTeoAliasDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_alias_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	aliasName := idSplit[1]

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("alias_name", aliasName)

	respData, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_alias_domain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ZoneId != nil {
		_ = d.Set("zone_id", respData.ZoneId)
	}

	if respData.AliasName != nil {
		_ = d.Set("alias_name", respData.AliasName)
	}

	if respData.TargetName != nil {
		_ = d.Set("target_name", respData.TargetName)
	}

	if respData.Status != nil {
		// paused is true if status is "stop", false otherwise
		paused := *respData.Status == "stop"
		_ = d.Set("paused", paused)
	}

	return nil
}

func resourceTencentCloudTeoAliasDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_alias_domain.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	aliasName := idSplit[1]

	// update target_name
	if d.HasChange("target_name") {
		request := teo.NewModifyAliasDomainRequest()
		request.ZoneId = helper.String(zoneId)
		request.AliasName = helper.String(aliasName)

		if v, ok := d.GetOk("target_name"); ok {
			request.TargetName = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyAliasDomainWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update teo alias domain failed, reason:%+v", logId, err)
			return err
		}

		// wait for target_name update to complete
		if err := resourceTencentCloudTeoAliasDomainUpdatePostHandle(ctx, zoneId, aliasName, d.Get("paused").(bool), meta, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	// update paused status
	if d.HasChange("paused") {
		request := teo.NewModifyAliasDomainStatusRequest()
		request.ZoneId = helper.String(zoneId)
		request.AliasNames = []*string{helper.String(aliasName)}

		if v, ok := d.GetOkExists("paused"); ok {
			request.Paused = helper.Bool(v.(bool))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyAliasDomainStatusWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update teo alias domain paused status failed, reason:%+v", logId, err)
			return err
		}

		// wait for status update to complete
		paused := d.Get("paused").(bool)
		if err := resourceTencentCloudTeoAliasDomainUpdatePostHandle(ctx, zoneId, aliasName, paused, meta, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}

	return resourceTencentCloudTeoAliasDomainRead(d, meta)
}

func resourceTencentCloudTeoAliasDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_alias_domain.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewDeleteAliasDomainRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	aliasName := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.AliasNames = []*string{helper.String(aliasName)}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DeleteAliasDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete teo alias domain failed, reason:%+v", logId, err)
		return err
	}

	// wait for deletion to complete
	if err := resourceTencentCloudTeoAliasDomainDeletePostHandle(ctx, zoneId, aliasName, meta, d.Timeout(schema.TimeoutDelete)); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudTeoAliasDomainCreatePostHandle(ctx context.Context, zoneId, aliasName string, meta interface{}, timeout time.Duration) error {
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(timeout, func() *resource.RetryError {
		respData, e := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if respData == nil {
			return resource.NonRetryableError(fmt.Errorf("alias domain %s not found", aliasName))
		}

		if respData.Status != nil && (*respData.Status == "active" || *respData.Status == "stop") {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("alias domain status is %s, waiting for it to be active or stop", *respData.Status))
	})

	if err != nil {
		return fmt.Errorf("wait for alias domain creation to complete failed: %v", err)
	}

	return nil
}

func resourceTencentCloudTeoAliasDomainUpdatePostHandle(ctx context.Context, zoneId, aliasName string, paused bool, meta interface{}, timeout time.Duration) error {
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	expectedStatus := "stop"
	if !paused {
		expectedStatus = "active"
	}

	err := resource.Retry(timeout, func() *resource.RetryError {
		respData, e := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if respData == nil {
			return resource.NonRetryableError(fmt.Errorf("alias domain %s not found", aliasName))
		}

		if respData.Status != nil && *respData.Status == expectedStatus {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("alias domain status is %s, waiting for it to be %s", *respData.Status, expectedStatus))
	})

	if err != nil {
		return fmt.Errorf("wait for alias domain update to complete failed: %v", err)
	}

	return nil
}

func resourceTencentCloudTeoAliasDomainDeletePostHandle(ctx context.Context, zoneId, aliasName string, meta interface{}, timeout time.Duration) error {
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(timeout, func() *resource.RetryError {
		respData, e := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if respData == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("alias domain %s still exists", aliasName))
	})

	if err != nil {
		return fmt.Errorf("wait for alias domain deletion to complete failed: %v", err)
	}

	return nil
}
