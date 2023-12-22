package ssl

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSslUpdateCertificateRecordRollbackOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslUpdateCertificateRecordRollbackCreate,
		Read:   resourceTencentCloudSslUpdateCertificateRecordRollbackRead,
		Delete: resourceTencentCloudSslUpdateCertificateRecordRollbackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"deploy_record_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Deployment record ID to be rolled back.",
			},
		},
	}
}

func resourceTencentCloudSslUpdateCertificateRecordRollbackCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_update_certificate_record_rollback_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = ssl.NewUpdateCertificateRecordRollbackRequest()
		response       = ssl.NewUpdateCertificateRecordRollbackResponse()
		deployRecordId int64
	)
	if v, _ := d.GetOk("deploy_record_id"); v != nil {
		request.DeployRecordId = helper.StrToInt64Point(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient().UpdateCertificateRecordRollback(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl updateCertificateRecordRollback failed, reason:%+v", logId, err)
		return err
	}
	if response != nil && response.Response != nil && response.Response.DeployRecordId != nil {
		deployRecordId = *response.Response.DeployRecordId
	}
	d.SetId(helper.Int64ToStr(deployRecordId))

	return resourceTencentCloudSslUpdateCertificateRecordRollbackRead(d, meta)
}

func resourceTencentCloudSslUpdateCertificateRecordRollbackRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_update_certificate_record_rollback_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslUpdateCertificateRecordRollbackDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_update_certificate_record_rollback_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
