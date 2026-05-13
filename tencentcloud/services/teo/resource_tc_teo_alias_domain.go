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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
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

			"cert_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Certificate configuration. Valid values: `none` (no configuration), `hosting` (SSL hosted certificate). Default value: `none`.",
			},

			"cert_id": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Certificate ID list. Required when `cert_type` is `hosting`.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"paused": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the alias domain is disabled. `false`: enabled; `true`: disabled.",
			},

			// Computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alias domain status. Valid values: `active` (effective), `pending` (deploying), `conflict` (reclaimed), `stop` (disabled).",
			},

			"forbid_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Block mode. Valid values: `0` (not blocked), `11` (compliance blocked), `14` (not registered blocked).",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alias domain creation time.",
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alias domain modification time.",
			},
		},
	}
}

func resourceTencentCloudTeoAliasDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_alias_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewCreateAliasDomainRequest()
		zoneId  string
		alias   string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("alias_name"); ok {
		request.AliasName = helper.String(v.(string))
		alias = v.(string)
	}

	if v, ok := d.GetOk("target_name"); ok {
		request.TargetName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cert_type"); ok {
		request.CertType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cert_id"); ok {
		for _, id := range v.([]interface{}) {
			request.CertId = append(request.CertId, helper.String(id.(string)))
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateAliasDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo alias domain failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{zoneId, alias}, tccommon.FILED_SP))

	// Wait for alias domain to become active
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if _, err := (&resource.StateChangeConf{
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{"pending"},
		Target:     []string{"active"},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Refresh: func() (interface{}, string, error) {
			aliasDomain, e := service.DescribeTeoAliasDomainById(ctx, zoneId, alias)
			if e != nil {
				return nil, "", e
			}
			if aliasDomain == nil {
				return nil, "pending", nil
			}
			status := "pending"
			if aliasDomain.Status != nil {
				status = *aliasDomain.Status
			}
			return aliasDomain, status, nil
		},
	}).WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for teo alias domain (%s) to become active: %s", d.Id(), err)
	}

	// If paused=true on creation, disable right after creation
	if v, ok := d.GetOkExists("paused"); ok && v.(bool) {
		pauseRequest := teo.NewModifyAliasDomainStatusRequest()
		pauseRequest.ZoneId = helper.String(zoneId)
		pauseRequest.Paused = helper.Bool(true)
		pauseRequest.AliasNames = []*string{helper.String(alias)}

		pauseErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyAliasDomainStatusWithContext(ctx, pauseRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, pauseRequest.GetAction(), pauseRequest.ToJsonString(), result.ToJsonString())
			return nil
		})

		if pauseErr != nil {
			log.Printf("[CRITAL]%s set teo alias domain status failed, reason:%+v", logId, pauseErr)
			return pauseErr
		}
	}

	return resourceTencentCloudTeoAliasDomainRead(d, meta)
}

func resourceTencentCloudTeoAliasDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_alias_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	aliasName := idSplit[1]

	respData, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_alias_domain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("alias_name", aliasName)

	if respData.TargetName != nil {
		_ = d.Set("target_name", respData.TargetName)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.ForbidMode != nil {
		_ = d.Set("forbid_mode", respData.ForbidMode)
	}

	if respData.CreatedOn != nil {
		_ = d.Set("created_on", respData.CreatedOn)
	}

	if respData.ModifiedOn != nil {
		_ = d.Set("modified_on", respData.ModifiedOn)
	}

	// Note: cert_type and cert_id are not returned by DescribeAliasDomains,
	// so we do not overwrite them to prevent state drift.

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
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	aliasName := idSplit[1]

	// Update target_name / cert_type / cert_id via ModifyAliasDomain
	if d.HasChanges("target_name", "cert_type", "cert_id") {
		request := teo.NewModifyAliasDomainRequest()
		request.ZoneId = helper.String(zoneId)
		request.AliasName = helper.String(aliasName)

		if v, ok := d.GetOk("target_name"); ok {
			request.TargetName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("cert_type"); ok {
			request.CertType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("cert_id"); ok {
			for _, id := range v.([]interface{}) {
				request.CertId = append(request.CertId, helper.String(id.(string)))
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyAliasDomainWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo alias domain failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	// Update paused status via ModifyAliasDomainStatus
	if d.HasChange("paused") {
		paused := d.Get("paused").(bool)
		statusRequest := teo.NewModifyAliasDomainStatusRequest()
		statusRequest.ZoneId = helper.String(zoneId)
		statusRequest.Paused = helper.Bool(paused)
		statusRequest.AliasNames = []*string{helper.String(aliasName)}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyAliasDomainStatusWithContext(ctx, statusRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, statusRequest.GetAction(), statusRequest.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo alias domain status failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		// Enabling (paused=false) is async: wait for status to become active
		if !paused {
			service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
			if _, err := (&resource.StateChangeConf{
				Delay:      5 * time.Second,
				MinTimeout: 3 * time.Second,
				Pending:    []string{"stop", "pending"},
				Target:     []string{"active"},
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Refresh: func() (interface{}, string, error) {
					aliasDomain, e := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
					if e != nil {
						return nil, "", e
					}
					if aliasDomain == nil {
						return nil, "pending", nil
					}
					status := "pending"
					if aliasDomain.Status != nil {
						status = *aliasDomain.Status
					}
					return aliasDomain, status, nil
				},
			}).WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for teo alias domain (%s) to become active: %s", d.Id(), err)
			}
		}
	}

	return resourceTencentCloudTeoAliasDomainRead(d, meta)
}

func resourceTencentCloudTeoAliasDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_alias_domain.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	aliasName := idSplit[1]

	// Must ensure the alias domain is disabled (stop) before deleting.
	// If current status is not "stop", disable it first.
	aliasDomain, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		return err
	}

	if aliasDomain == nil {
		return nil
	}

	if aliasDomain.Status == nil && *aliasDomain.Status != "stop" {
		pauseRequest := teo.NewModifyAliasDomainStatusRequest()
		pauseRequest.ZoneId = helper.String(zoneId)
		pauseRequest.Paused = helper.Bool(true)
		pauseRequest.AliasNames = []*string{helper.String(aliasName)}

		pauseErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyAliasDomainStatusWithContext(ctx, pauseRequest)
			if e != nil {
				return tccommon.RetryError(e)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, pauseRequest.GetAction(), pauseRequest.ToJsonString(), result.ToJsonString())
			return nil
		})

		if pauseErr != nil {
			log.Printf("[CRITAL]%s disable teo alias domain before delete failed, reason:%+v", logId, pauseErr)
			return pauseErr
		}
	}

	request := teo.NewDeleteAliasDomainRequest()
	request.ZoneId = helper.String(zoneId)
	request.AliasNames = []*string{helper.String(aliasName)}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteAliasDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo alias domain failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
