package vcube

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vcubev20220410 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vcube/v20220410"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVcubeApplicationAndWebPlayerLicense() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVcubeApplicationAndWebPlayerLicenseCreate,
		Read:   resourceTencentCloudVcubeApplicationAndWebPlayerLicenseRead,
		Delete: resourceTencentCloudVcubeApplicationAndWebPlayerLicenseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Application name.",
			},

			"domain_list": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Description: "Domain list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// computed
			"license_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "License ID.",
			},

			"app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Account App ID.",
			},

			"app_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Application type: formal: formal application, test: test application.",
			},

			"application_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "User Application ID.",
			},

			"license_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "License key.",
			},

			"license_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "License url.",
			},
		},
	}
}

func resourceTencentCloudVcubeApplicationAndWebPlayerLicenseCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vcube_application_and_web_player_license.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = vcubev20220410.NewCreateApplicationAndWebPlayerLicenseRequest()
		response  = vcubev20220410.NewCreateApplicationAndWebPlayerLicenseResponse()
		licenseId string
	)

	if v, ok := d.GetOk("app_name"); ok {
		request.AppName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("domain_list"); ok {
		domainListSet := v.(*schema.Set).List()
		for i := range domainListSet {
			domainList := domainListSet[i].(string)
			request.DomainList = append(request.DomainList, helper.String(domainList))
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVcubeV20220410Client().CreateApplicationAndWebPlayerLicenseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create vcube application and web player license failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create vcube application and web player license failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.LicenseId == nil {
		return fmt.Errorf("LicenseId is nil.")
	}

	licenseId = helper.UInt64ToStr(*response.Response.LicenseId)
	d.SetId(licenseId)
	return resourceTencentCloudVcubeApplicationAndWebPlayerLicenseRead(d, meta)
}

func resourceTencentCloudVcubeApplicationAndWebPlayerLicenseRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vcube_application_and_web_player_license.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = VcubeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		licenseId = d.Id()
	)

	respData, err := service.DescribeVcubeApplicationAndWebPlayerLicenseById(ctx, licenseId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vcube_application_and_web_player_license` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.AppName != nil {
		_ = d.Set("app_name", respData.AppName)
	}

	if respData.DomainList != nil {
		_ = d.Set("domain_list", respData.DomainList)
	}

	if respData.Licenses != nil && len(respData.Licenses) == 1 {
		if respData.Licenses[0].LicenseId != nil {
			_ = d.Set("license_id", respData.Licenses[0].LicenseId)
		}
	}

	if respData.AppId != nil {
		_ = d.Set("app_id", respData.AppId)
	}

	if respData.AppType != nil {
		_ = d.Set("app_type", respData.AppType)
	}

	if respData.ApplicationId != nil {
		_ = d.Set("application_id", respData.ApplicationId)
	}

	if respData.LicenseKey != nil {
		_ = d.Set("license_key", respData.LicenseKey)
	}

	if respData.LicenseUrl != nil {
		_ = d.Set("license_url", respData.LicenseUrl)
	}

	return nil
}

func resourceTencentCloudVcubeApplicationAndWebPlayerLicenseDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vcube_application_and_web_player_license.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = vcubev20220410.NewDeleteApplicationAndWebPlayerLicenseRequest()
		licenseId = d.Id()
	)

	request.LicenseId = helper.StrToUint64Point(licenseId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVcubeV20220410Client().DeleteApplicationAndWebPlayerLicenseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})
	if reqErr != nil {
		log.Printf("[CRITAL]%s delete vcube application and web player license failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
