/*
Provides a resource to create a ssl update_certificate_record_rollback

Example Usage

```hcl
resource "tencentcloud_ssl_update_certificate_record_rollback" "update_certificate_record_rollback" {
  deploy_record_id =
}
```

Import

ssl update_certificate_record_rollback can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_update_certificate_record_rollback.update_certificate_record_rollback update_certificate_record_rollback_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSslUpdateCertificateRecordRollback() *schema.Resource {
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
				Type:        schema.TypeInt,
				Description: "Deployment record ID to be rolled back.",
			},
		},
	}
}

func resourceTencentCloudSslUpdateCertificateRecordRollbackCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_update_certificate_record_rollback.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = ssl.NewUpdateCertificateRecordRollbackRequest()
		response       = ssl.NewUpdateCertificateRecordRollbackResponse()
		deployRecordId int
	)
	if v, _ := d.GetOk("deploy_record_id"); v != nil {
		deployRecordId = v.(int64)
		request.DeployRecordId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSslClient().UpdateCertificateRecordRollback(request)
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

	deployRecordId = *response.Response.DeployRecordId
	d.SetId(helper.Int64ToStr(deployRecordId))

	return resourceTencentCloudSslUpdateCertificateRecordRollbackRead(d, meta)
}

func resourceTencentCloudSslUpdateCertificateRecordRollbackRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_update_certificate_record_rollback.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslUpdateCertificateRecordRollbackDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_update_certificate_record_rollback.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
