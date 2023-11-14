/*
Provides a resource to create a live domain_cert

Example Usage

```hcl
resource "tencentcloud_live_domain_cert" "domain_cert" {
  domain_name = "5000.livepush.play.com"
  type = "Formal"
}
```

Import

live domain_cert can be imported using the id, e.g.

```
terraform import tencentcloud_live_domain_cert.domain_cert domain_cert_id
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

func resourceTencentCloudLiveDomainCert() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveDomainCertCreate,
		Read:   resourceTencentCloudLiveDomainCertRead,
		Update: resourceTencentCloudLiveDomainCertUpdate,
		Delete: resourceTencentCloudLiveDomainCertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain Name.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Gray: unbind gray level rules Formal (default): Unbind formal rules.",
			},
		},
	}
}

func resourceTencentCloudLiveDomainCertCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain_cert.create")()
	defer inconsistentCheck(d, meta)()

	var domainName string
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
	}

	d.SetId(domainName)

	return resourceTencentCloudLiveDomainCertUpdate(d, meta)
}

func resourceTencentCloudLiveDomainCertRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain_cert.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	domainCertId := d.Id()

	domainCert, err := service.DescribeLiveDomainCertById(ctx, domainName)
	if err != nil {
		return err
	}

	if domainCert == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveDomainCert` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if domainCert.DomainName != nil {
		_ = d.Set("domain_name", domainCert.DomainName)
	}

	if domainCert.Type != nil {
		_ = d.Set("type", domainCert.Type)
	}

	return nil
}

func resourceTencentCloudLiveDomainCertUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain_cert.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewUnBindLiveDomainCertRequest()

	domainCertId := d.Id()

	request.DomainName = &domainName

	immutableArgs := []string{"domain_name", "type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().UnBindLiveDomainCert(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live domainCert failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLiveDomainCertRead(d, meta)
}

func resourceTencentCloudLiveDomainCertDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_domain_cert.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
