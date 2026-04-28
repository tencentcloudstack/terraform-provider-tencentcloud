package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

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
				Description: "Certificate configuration type. Valid values: `none` (no certificate), `hosting` (SSL managed certificate). Default is `none`. When modifying, `apply` (apply for free certificate) is also supported.",
			},

			"cert_id": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Certificate ID list. Required when `cert_type` is `hosting`.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Alias domain status. Valid values: `active`, `pending`, `conflict`, `stop`.",
			},

			"forbid_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Forbid mode. Valid values: `0` (not forbidden), `11` (compliance forbid), `14` (no-ICP forbid).",
			},

			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the alias domain.",
			},

			"modified_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Modification time of the alias domain.",
			},
		},
	}
}

func resourceTencentCloudTeoAliasDomainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_alias_domain.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		request = teo.NewCreateAliasDomainRequest()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("alias_name"); ok {
		request.AliasName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("target_name"); ok {
		request.TargetName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cert_type"); ok {
		request.CertType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cert_id"); ok {
		certIdList := make([]*string, 0)
		for _, item := range v.([]interface{}) {
			certIdList = append(certIdList, helper.String(item.(string)))
		}
		request.CertId = certIdList
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateAliasDomainWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo alias domain failed, response is nil."))
		}

		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo alias domain failed, reason:%+v", logId, err)
		return err
	}

	zoneId := d.Get("zone_id").(string)
	aliasName := d.Get("alias_name").(string)
	d.SetId(strings.Join([]string{zoneId, aliasName}, tccommon.FILED_SP))

	return resourceTencentCloudTeoAliasDomainRead(d, meta)
}

func resourceTencentCloudTeoAliasDomainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_alias_domain.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	aliasName := idSplit[1]

	_ = d.Set("zone_id", zoneId)

	respData, err := service.DescribeTeoAliasDomainById(ctx, zoneId, aliasName)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_alias_domain` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.TargetName != nil {
		_ = d.Set("target_name", respData.TargetName)
	}

	if respData.AliasName != nil {
		_ = d.Set("alias_name", respData.AliasName)
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

	return nil
}

func resourceTencentCloudTeoAliasDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_alias_domain.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	aliasName := idSplit[1]

	needChange := false
	mutableArgs := []string{"target_name", "cert_type", "cert_id"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
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
			certIdList := make([]*string, 0)
			for _, item := range v.([]interface{}) {
				certIdList = append(certIdList, helper.String(item.(string)))
			}
			request.CertId = certIdList
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
	}

	return resourceTencentCloudTeoAliasDomainRead(d, meta)
}

func resourceTencentCloudTeoAliasDomainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_alias_domain.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	aliasName := idSplit[1]

	var (
		request = teo.NewDeleteAliasDomainRequest()
	)

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

	return nil
}
