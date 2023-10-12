/*
Provides a resource to create a ssl update_certificate_record_retry

Example Usage

```hcl
resource "tencentcloud_ssl_update_certificate_record_retry_operation" "update_certificate_record_retry" {
  deploy_record_id = "1603"
}
```

Import

ssl update_certificate_record_retry can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_update_certificate_record_retry_operation.update_certificate_record_retry update_certificate_record_retry_id
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

func resourceTencentCloudSslUpdateCertificateRecordRetryOperation() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_ssl_update_certificate_record_retry_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSSLCertificateClient().UpdateCertificateRecordRetry(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_ssl_update_certificate_record_retry_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslUpdateCertificateRecordRetryDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_update_certificate_record_retry_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
