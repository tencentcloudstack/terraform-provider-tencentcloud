/*
Provides a resource to create a live domain_cert_bindings

Example Usage

```hcl
resource "tencentcloud_live_domain_cert_bindings" "domain_cert_bindings" {
  domain_infos {
		domain_name = "abc.com"
		status = 1

  }
  cloud_cert_id = "123"
  certificate_public_key = "xxx"
  certificate_private_key = "xxx"
  certificate_alias = "adc"
}
```

Import

live domain_cert_bindings can be imported using the id, e.g.

```
terraform import tencentcloud_live_domain_cert_bindings.domain_cert_bindings domain_cert_bindings_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"log"
)

func resourceTencentCloudLiveDomainCertBindings() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveDomainCertBindingsCreate,
		Read:   resourceTencentCloudLiveDomainCertBindingsRead,
		Update: resourceTencentCloudLiveDomainCertBindingsUpdate,
		Delete: resourceTencentCloudLiveDomainCertBindingsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_infos": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Playback domain name/status information list of the certificate to be bound. If CloudCertId and certificate public key private key pair are not transferred, and the domain name list has binding rules, only batch update the enabling status of domain name https rules, and upload existing self owned certificates that have not been uploaded to Tencent Cloud SSL.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Domain Name.",
						},
						"status": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Whether to enable https rules for domain names. 1: Enable 0: Disabled -1: Keep unchanged.",
						},
					},
				},
			},

			"cloud_cert_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The certificate ID of Tencent Cloud SSL.",
			},

			"certificate_public_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Certificate public key. Choose one of CloudCertId and public key private key pair. If CloudCertId is selected, public key and private key parameters will be discarded. Otherwise, public key private key pair will be automatically uploaded to ssl to create a new certificate, and CloudCertId returned after successful upload will be used.",
			},

			"certificate_private_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Certificate private key. Choose one of CloudCertId and public key private key pair. If CloudCertId is transferred, public key and private key parameters will be discarded. Otherwise, public key private key pair will be automatically uploaded to ssl to create a new certificate, and CloudCertId returned after successful upload will be used.",
			},

			"certificate_alias": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The comments uploaded to the SSL certificate center are only valid when creating a new certificate. Ignored when transferring CloudCertId.",
			},
		},
	}
}

func resourceTencentCloudLiveDomainCertBindingsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain_cert_bindings.create")()
	defer inconsistentCheck(d, meta)()

	var domainName string
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
	}

	d.SetId(domainName)

	return resourceTencentCloudLiveDomainCertBindingsUpdate(d, meta)
}

func resourceTencentCloudLiveDomainCertBindingsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain_cert_bindings.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	domainCertBindingsId := d.Id()

	domainCertBindings, err := service.DescribeLiveDomainCertBindingsById(ctx, domainName)
	if err != nil {
		return err
	}

	if domainCertBindings == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveDomainCertBindings` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if domainCertBindings.DomainInfos != nil {
		domainInfosList := []interface{}{}
		for _, domainInfos := range domainCertBindings.DomainInfos {
			domainInfosMap := map[string]interface{}{}

			if domainCertBindings.DomainInfos.DomainName != nil {
				domainInfosMap["domain_name"] = domainCertBindings.DomainInfos.DomainName
			}

			if domainCertBindings.DomainInfos.Status != nil {
				domainInfosMap["status"] = domainCertBindings.DomainInfos.Status
			}

			domainInfosList = append(domainInfosList, domainInfosMap)
		}

		_ = d.Set("domain_infos", domainInfosList)

	}

	if domainCertBindings.CloudCertId != nil {
		_ = d.Set("cloud_cert_id", domainCertBindings.CloudCertId)
	}

	if domainCertBindings.CertificatePublicKey != nil {
		_ = d.Set("certificate_public_key", domainCertBindings.CertificatePublicKey)
	}

	if domainCertBindings.CertificatePrivateKey != nil {
		_ = d.Set("certificate_private_key", domainCertBindings.CertificatePrivateKey)
	}

	if domainCertBindings.CertificateAlias != nil {
		_ = d.Set("certificate_alias", domainCertBindings.CertificateAlias)
	}

	return nil
}

func resourceTencentCloudLiveDomainCertBindingsUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain_cert_bindings.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewModifyLiveDomainCertBindingsRequest()

	domainCertBindingsId := d.Id()

	request.DomainName = &domainName

	immutableArgs := []string{"domain_infos", "cloud_cert_id", "certificate_public_key", "certificate_private_key", "certificate_alias"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().ModifyLiveDomainCertBindings(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live domainCertBindings failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLiveDomainCertBindingsRead(d, meta)
}

func resourceTencentCloudLiveDomainCertBindingsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain_cert_bindings.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
