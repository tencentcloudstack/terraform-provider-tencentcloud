/*
Provides a resource to create a tse cngw_certificate

Example Usage

```hcl
resource "tencentcloud_tse_cngw_certificate" "cngw_certificate" {
  gateway_id = ""
  bind_domains =
  cert_id = ""
  name = ""
  key = ""
  crt = ""
}
```

Import

tse cngw_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_tse_cngw_certificate.cngw_certificate cngw_certificate_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTseCngwCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTseCngwCertificateCreate,
		Read:   resourceTencentCloudTseCngwCertificateRead,
		Update: resourceTencentCloudTseCngwCertificateUpdate,
		Delete: resourceTencentCloudTseCngwCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway ID.",
			},

			"bind_domains": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Domains of the binding.",
			},

			"cert_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID of ssl platform.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Certificate name.",
			},

			"key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Private key of certificate.",
			},

			"crt": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Pem format of certificate.",
			},
		},
	}
}

func resourceTencentCloudTseCngwCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_certificate.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tse.NewCreateCloudNativeAPIGatewayCertificateRequest()
		response  = tse.NewCreateCloudNativeAPIGatewayCertificateResponse()
		gatewayId string
		id        string
	)
	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
		request.GatewayId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bind_domains"); ok {
		bindDomainsSet := v.(*schema.Set).List()
		for i := range bindDomainsSet {
			bindDomains := bindDomainsSet[i].(string)
			request.BindDomains = append(request.BindDomains, &bindDomains)
		}
	}

	if v, ok := d.GetOk("cert_id"); ok {
		request.CertId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("key"); ok {
		request.Key = helper.String(v.(string))
	}

	if v, ok := d.GetOk("crt"); ok {
		request.Crt = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().CreateCloudNativeAPIGatewayCertificate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tse cngwCertificate failed, reason:%+v", logId, err)
		return err
	}

	gatewayId = *response.Response.GatewayId
	d.SetId(strings.Join([]string{gatewayId, id}, FILED_SP))

	return resourceTencentCloudTseCngwCertificateRead(d, meta)
}

func resourceTencentCloudTseCngwCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_certificate.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	id := idSplit[1]

	cngwCertificate, err := service.DescribeTseCngwCertificateById(ctx, gatewayId, id)
	if err != nil {
		return err
	}

	if cngwCertificate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TseCngwCertificate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cngwCertificate.GatewayId != nil {
		_ = d.Set("gateway_id", cngwCertificate.GatewayId)
	}

	if cngwCertificate.BindDomains != nil {
		_ = d.Set("bind_domains", cngwCertificate.BindDomains)
	}

	if cngwCertificate.CertId != nil {
		_ = d.Set("cert_id", cngwCertificate.CertId)
	}

	if cngwCertificate.Name != nil {
		_ = d.Set("name", cngwCertificate.Name)
	}

	if cngwCertificate.Key != nil {
		_ = d.Set("key", cngwCertificate.Key)
	}

	if cngwCertificate.Crt != nil {
		_ = d.Set("crt", cngwCertificate.Crt)
	}

	return nil
}

func resourceTencentCloudTseCngwCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_certificate.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tse.NewUpdateCloudNativeAPIGatewayCertificateInfoRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	id := idSplit[1]

	request.GatewayId = &gatewayId
	request.Id = &id

	immutableArgs := []string{"gateway_id", "bind_domains", "cert_id", "name", "key", "crt"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("gateway_id") {
		if v, ok := d.GetOk("gateway_id"); ok {
			request.GatewayId = helper.String(v.(string))
		}
	}

	if d.HasChange("bind_domains") {
		if v, ok := d.GetOk("bind_domains"); ok {
			bindDomainsSet := v.(*schema.Set).List()
			for i := range bindDomainsSet {
				bindDomains := bindDomainsSet[i].(string)
				request.BindDomains = append(request.BindDomains, &bindDomains)
			}
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTseClient().UpdateCloudNativeAPIGatewayCertificateInfo(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tse cngwCertificate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTseCngwCertificateRead(d, meta)
}

func resourceTencentCloudTseCngwCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tse_cngw_certificate.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	gatewayId := idSplit[0]
	id := idSplit[1]

	if err := service.DeleteTseCngwCertificateById(ctx, gatewayId, id); err != nil {
		return err
	}

	return nil
}
