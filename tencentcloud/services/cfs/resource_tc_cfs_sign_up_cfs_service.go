package cfs

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
)

func ResourceTencentCloudCfsSignUpCfsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfsSignUpCfsServiceCreate,
		Read:   resourceTencentCloudCfsSignUpCfsServiceRead,
		Delete: resourceTencentCloudCfsSignUpCfsServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cfs_service_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Current status of the CFS service for this user. Valid values: creating (activating); created (activated).",
			},
		},
	}
}

func resourceTencentCloudCfsSignUpCfsServiceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfs_sign_up_cfs_service.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request          = cfs.NewSignUpCfsServiceRequest()
		response         = cfs.NewSignUpCfsServiceResponse()
		cfsServiceStatus string
	)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfsClient().SignUpCfsService(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cfs signUpCfsService failed, reason:%+v", logId, err)
		return nil
	}

	cfsServiceStatus = *response.Response.CfsServiceStatus
	d.SetId(cfsServiceStatus)

	return resourceTencentCloudCfsSignUpCfsServiceRead(d, meta)
}

func resourceTencentCloudCfsSignUpCfsServiceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_sign_up_cfs_service.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = cfs.NewDescribeCfsServiceStatusRequest()
		response = cfs.NewDescribeCfsServiceStatusResponse()
	)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfsClient().DescribeCfsServiceStatus(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cfs signUpCfsService failed, reason:%+v", logId, err)
		return nil
	}

	_ = d.Set("cfs_service_status", response.Response.CfsServiceStatus)

	return nil
}

func resourceTencentCloudCfsSignUpCfsServiceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfs_sign_up_cfs_service.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
