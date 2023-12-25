package ssl

import (
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSslCheckCertificateChainOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslCheckCertificateChainCreate,
		Read:   resourceTencentCloudSslCheckCertificateChainRead,
		Delete: resourceTencentCloudSslCheckCertificateChainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate_chain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The certificate chain to check.",
			},
		},
	}
}

func resourceTencentCloudSslCheckCertificateChainCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_check_certificate_chain_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = ssl.NewCheckCertificateChainRequest()
	)
	if v, ok := d.GetOk("certificate_chain"); ok {
		request.CertificateChain = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient().CheckCertificateChain(request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result != nil && result.Response != nil && !*result.Response.IsValid {
			err := fmt.Errorf("[DEBUG]%s Certificate chain failed to check, IsValid [%v]\n", logId, *result.Response.IsValid)
			return tccommon.RetryError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl checkCertificateChain failed, reason:%+v", logId, err)
		return err
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	d.SetId(id)
	return resourceTencentCloudSslCheckCertificateChainRead(d, meta)
}

func resourceTencentCloudSslCheckCertificateChainRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_check_certificate_chain_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslCheckCertificateChainDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_check_certificate_chain_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
