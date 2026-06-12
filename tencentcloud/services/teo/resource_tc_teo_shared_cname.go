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

func ResourceTencentCloudTeoSharedCname() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoSharedCnameCreate,
		Read:   resourceTencentCloudTeoSharedCnameRead,
		Update: resourceTencentCloudTeoSharedCnameUpdate,
		Delete: resourceTencentCloudTeoSharedCnameDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The zone ID of the shared CNAME.",
			},
			"shared_cname_prefix": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The shared CNAME prefix. Please enter a valid domain prefix, for example `test-api` or `test-api.com`, limited to 50 characters.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description. You can enter 1-50 characters.",
			},
			"shared_cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The full shared CNAME returned by the API.",
			},
			"ipssl_setting": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "IP SSL setting for the shared CNAME.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operation": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Operation type. Valid values: `bind`, `unbind`.",
						},
						"associated_domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The associated domain for IP SSL.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoSharedCnameCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_shared_cname.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = teo.NewCreateSharedCNAMERequest()
		response = teo.NewCreateSharedCNAMEResponse()
		zoneId   string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("shared_cname_prefix"); ok {
		request.SharedCNAMEPrefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateSharedCNAMEWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo_shared_cname failed, Response is nil"))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo_shared_cname failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	log.Printf("[DEBUG]%s create teo_shared_cname, logId: %s, current d.Id(): %s", logId, logId, d.Id())

	if response.Response.SharedCNAME == nil || *response.Response.SharedCNAME == "" {
		return fmt.Errorf("Create teo_shared_cname failed, SharedCNAME is nil or empty")
	}

	sharedCname := *response.Response.SharedCNAME
	d.SetId(strings.Join([]string{zoneId, sharedCname}, tccommon.FILED_SP))

	return resourceTencentCloudTeoSharedCnameRead(d, meta)
}

func resourceTencentCloudTeoSharedCnameRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_shared_cname.read")()
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
	sharedCname := idSplit[1]

	request := teo.NewDescribeSharedCNAMERequest()
	request.ZoneId = helper.String(zoneId)
	request.Filters = []*teo.AdvancedFilter{
		{
			Name:   helper.String("shared-cname"),
			Values: []*string{helper.String(sharedCname)},
		},
	}
	request.Limit = helper.Int64(200)

	var respData *teo.SharedCNAMEInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DescribeSharedCNAMEWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Read teo_shared_cname failed, Response is nil"))
		}

		if len(result.Response.SharedCNAMEInfo) == 0 {
			return nil
		}

		respData = result.Response.SharedCNAMEInfo[0]
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read teo_shared_cname failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `teo_shared_cname` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if respData.SharedCNAME != nil {
		_ = d.Set("shared_cname", respData.SharedCNAME)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	return nil
}

func resourceTencentCloudTeoSharedCnameUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_shared_cname.update")()
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
	sharedCname := idSplit[1]

	needChange := false
	mutableArgs := []string{"description", "ipssl_setting"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teo.NewModifySharedCNAMERequest()
		request.ZoneId = helper.String(zoneId)
		request.SharedCNAME = helper.String(sharedCname)

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOk("ipssl_setting"); ok {
			ipsslSettingList := v.([]interface{})
			if len(ipsslSettingList) > 0 {
				ipsslSettingMap := ipsslSettingList[0].(map[string]interface{})
				ipsslSetting := &teo.IPSSLSetting{}
				if v, ok := ipsslSettingMap["operation"].(string); ok && v != "" {
					ipsslSetting.Operation = helper.String(v)
				}
				if v, ok := ipsslSettingMap["associated_domain"].(string); ok && v != "" {
					ipsslSetting.AssociatedDomain = helper.String(v)
				}
				request.IPSSLSetting = ipsslSetting
			}
		}

		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifySharedCNAMEWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}

			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update teo_shared_cname failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudTeoSharedCnameRead(d, meta)
}

func resourceTencentCloudTeoSharedCnameDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_shared_cname.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewDeleteSharedCNAMERequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	sharedCname := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.SharedCNAME = helper.String(sharedCname)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DeleteSharedCNAMEWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo_shared_cname failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
