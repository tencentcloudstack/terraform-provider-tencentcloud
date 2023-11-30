package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSslCompleteCertificateOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslCompleteCertificateCreate,
		Read:   resourceTencentCloudSslCompleteCertificateRead,
		Delete: resourceTencentCloudSslCompleteCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID.",
			},
		},
	}
}

func resourceTencentCloudSslCompleteCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_complete_certificate_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = ssl.NewCompleteCertificateRequest()
		response      = ssl.NewCompleteCertificateResponse()
		certificateId string
	)
	if v, ok := d.GetOk("certificate_id"); ok {
		request.CertificateId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSSLCertificateClient().CompleteCertificate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl completeCertificate failed, reason:%+v", logId, err)
		return err
	}
	if response != nil && response.Response != nil && response.Response.CertificateId != nil {
		certificateId = *response.Response.CertificateId
	}
	d.SetId(certificateId)

	return resourceTencentCloudSslCompleteCertificateRead(d, meta)
}

func resourceTencentCloudSslCompleteCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_complete_certificate_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslCompleteCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_complete_certificate_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
