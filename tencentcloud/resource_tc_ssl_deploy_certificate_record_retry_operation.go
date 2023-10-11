/*
Provides a resource to create a ssl deploy_certificate_record_retry

Example Usage

```hcl
resource "tencentcloud_ssl_deploy_certificate_record_retry_operation" "deploy_certificate_record_retry" {
  deploy_record_id = 35474
}
```

Import

ssl deploy_certificate_record_retry can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_deploy_certificate_record_retry_operation.deploy_certificate_record_retry deploy_certificate_record_retry_id
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

func resourceTencentCloudSslDeployCertificateRecordRetryOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslDeployCertificateRecordRetryCreate,
		Read:   resourceTencentCloudSslDeployCertificateRecordRetryRead,
		Delete: resourceTencentCloudSslDeployCertificateRecordRetryDelete,
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

func resourceTencentCloudSslDeployCertificateRecordRetryCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_deploy_certificate_record_retry_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = ssl.NewDeployCertificateRecordRetryRequest()
		deployRecordId uint64
	)
	if v, _ := d.GetOk("deploy_record_id"); v != nil {
		deployRecordId = (uint64)(v.(int))
		request.DeployRecordId = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("deploy_record_detail_id"); v != nil {
		request.DeployRecordDetailId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSSLCertificateClient().DeployCertificateRecordRetry(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl deployCertificateRecordRetry failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(helper.UInt64ToStr(deployRecordId))

	return resourceTencentCloudSslDeployCertificateRecordRetryRead(d, meta)
}

func resourceTencentCloudSslDeployCertificateRecordRetryRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_deploy_certificate_record_retry_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslDeployCertificateRecordRetryDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_deploy_certificate_record_retry_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
