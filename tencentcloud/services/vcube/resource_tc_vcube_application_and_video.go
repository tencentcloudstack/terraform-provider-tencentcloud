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

func ResourceTencentCloudVcubeApplicationAndVideo() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVcubeApplicationAndVideoCreate,
		Read:   resourceTencentCloudVcubeApplicationAndVideoRead,
		Delete: resourceTencentCloudVcubeApplicationAndVideoDelete,
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

			"bundle_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "IOS bundle ID. Choose at least one of `bundle_id` and `package_name`.",
			},

			"package_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Android package name. Choose at least one of `bundle_id` and `package_name`.",
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

func resourceTencentCloudVcubeApplicationAndVideoCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vcube_application_and_video.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = vcubev20220410.NewCreateApplicationAndVideoRequest()
		response  = vcubev20220410.NewCreateApplicationAndVideoResponse()
		licenseId string
	)

	if v, ok := d.GetOk("app_name"); ok {
		request.AppName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bundle_id"); ok {
		request.BundleId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("package_name"); ok {
		request.PackageName = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVcubeV20220410Client().CreateApplicationAndVideoWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create vcube application and video failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create vcube application and video failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.LicenseId == nil {
		return fmt.Errorf("LicenseId is nil.")
	}

	licenseId = helper.UInt64ToStr(*response.Response.LicenseId)
	d.SetId(licenseId)
	return resourceTencentCloudVcubeApplicationAndVideoRead(d, meta)
}

func resourceTencentCloudVcubeApplicationAndVideoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vcube_application_and_video.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = VcubeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		licenseId = d.Id()
	)

	respData, err := service.DescribeVcubeApplicationAndVideoById(ctx, licenseId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vcube_application_and_video` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.AppName != nil {
		_ = d.Set("app_name", respData.AppName)
	}

	if respData.BundleId != nil {
		_ = d.Set("bundle_id", respData.BundleId)
	}

	if respData.PackageName != nil {
		_ = d.Set("package_name", respData.PackageName)
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

func resourceTencentCloudVcubeApplicationAndVideoDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vcube_application_and_video.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = vcubev20220410.NewDeleteApplicationAndVideoLicenseRequest()
		licenseId = d.Id()
	)

	request.LicenseId = helper.StrToUint64Point(licenseId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVcubeV20220410Client().DeleteApplicationAndVideoLicenseWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete vcube application and video failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
