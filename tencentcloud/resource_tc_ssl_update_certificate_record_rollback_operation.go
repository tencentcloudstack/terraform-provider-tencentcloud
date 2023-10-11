/*
Provides a resource to create a ssl update_certificate_record_rollback

Example Usage

```hcl
resource "tencentcloud_ssl_update_certificate_record_rollback_operation" "update_certificate_record_rollback" {
  deploy_record_id = "1603"
}
```

Import

ssl update_certificate_record_rollback can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_update_certificate_record_rollback_operation.update_certificate_record_rollback update_certificate_record_rollback_id
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSslUpdateCertificateRecordRollbackOperation() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_ssl_update_certificate_record_rollback_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = ssl.NewUpdateCertificateRecordRollbackRequest()
		response       = ssl.NewUpdateCertificateRecordRollbackResponse()
		deployRecordId int64
	)
	if v, _ := d.GetOk("deploy_record_id"); v != nil {
		request.DeployRecordId = helper.StrToInt64Point(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSSLCertificateClient().UpdateCertificateRecordRollback(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_ssl_update_certificate_record_rollback_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslUpdateCertificateRecordRollbackDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_update_certificate_record_rollback_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
