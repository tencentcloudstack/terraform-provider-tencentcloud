package ssl

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSslUpdateCertificateRecordRetryOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslUpdateCertificateRecordRetryCreate,
		Read:   resourceTencentCloudSslUpdateCertificateRecordRetryRead,
		Delete: resourceTencentCloudSslUpdateCertificateRecordRetryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"deploy_record_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Deployment record ID to be retried.",
			},

			"deploy_record_detail_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Deployment record details ID to be retried.",
			},
		},
	}
}

func resourceTencentCloudSslUpdateCertificateRecordRetryCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_update_certificate_record_retry_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = ssl.NewUpdateCertificateRecordRetryRequest()
		deployRecordId int
	)
	if v, _ := d.GetOk("deploy_record_id"); v != nil {
		deployRecordId = v.(int)
		request.DeployRecordId = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("deploy_record_detail_id"); v != nil {
		request.DeployRecordDetailId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient().UpdateCertificateRecordRetry(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl updateCertificateRecordRetry failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(helper.IntToStr(deployRecordId))

	return resourceTencentCloudSslUpdateCertificateRecordRetryRead(d, meta)
}

func resourceTencentCloudSslUpdateCertificateRecordRetryRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_update_certificate_record_retry_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslUpdateCertificateRecordRetryDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_update_certificate_record_retry_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
