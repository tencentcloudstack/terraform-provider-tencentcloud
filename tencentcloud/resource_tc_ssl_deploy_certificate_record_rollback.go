/*
Provides a resource to create a ssl deploy_certificate_record_rollback

Example Usage

```hcl
resource "tencentcloud_ssl_deploy_certificate_record_rollback" "deploy_certificate_record_rollback" {
  deploy_record_id =
}
```

Import

ssl deploy_certificate_record_rollback can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_deploy_certificate_record_rollback.deploy_certificate_record_rollback deploy_certificate_record_rollback_id
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

func resourceTencentCloudSslDeployCertificateRecordRollback() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslDeployCertificateRecordRollbackCreate,
		Read:   resourceTencentCloudSslDeployCertificateRecordRollbackRead,
		Delete: resourceTencentCloudSslDeployCertificateRecordRollbackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"deploy_record_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Deployment record ID to be rollback.",
			},
		},
	}
}

func resourceTencentCloudSslDeployCertificateRecordRollbackCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_deploy_certificate_record_rollback.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = ssl.NewDeployCertificateRecordRollbackRequest()
		response       = ssl.NewDeployCertificateRecordRollbackResponse()
		deployRecordId uint64
	)
	if v, _ := d.GetOk("deploy_record_id"); v != nil {
		request.DeployRecordId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSSLCertificateClient().DeployCertificateRecordRollback(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl deployCertificateRecordRollback failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(helper.UInt64ToStr(deployRecordId))

	return resourceTencentCloudSslDeployCertificateRecordRollbackRead(d, meta)
}

func resourceTencentCloudSslDeployCertificateRecordRollbackRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_deploy_certificate_record_rollback.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslDeployCertificateRecordRollbackDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_deploy_certificate_record_rollback.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
