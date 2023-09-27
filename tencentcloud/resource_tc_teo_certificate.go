/*
Provides a resource to create a teo certificate

Example Usage

```hcl
resource "tencentcloud_teo_certificate" "certificate" {
  zone_id = ""
  host = ""
  cert_id = ""
  alias = ""
  type = ""
  expire_time = ""
  deploy_time = ""
  sign_algo = ""
  common_name = ""
  mode = ""
}
```

Import

teo certificate can be imported using the id, e.g.

```
terraform import tencentcloud_teo_certificate.certificate certificate_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoCertificateCreate,
		Read:   resourceTencentCloudTeoCertificateRead,
		Update: resourceTencentCloudTeoCertificateUpdate,
		Delete: resourceTencentCloudTeoCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"hosts": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Acceleration domain name that needs to modify the certificate configuration.",
			},

			"cert_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the server certificate.Note: This field may return null, indicating that no valid values can be obtained.",
			},
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alias of the certificate.Note: This field may return null, indicating that no valid values can be obtained.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the certificate. Values: `default`: Default certificate; `upload`: Specified certificate; `managed`: Tencent Cloud-managed certificate; Note: This field may return `null`, indicating that no valid value can be obtained.",
			},
			"expire_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Time when the certificate expires.Note: This field may return null, indicating that no valid values can be obtained.",
			},
			"deploy_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Time when the certificate is deployed.Note: This field may return null, indicating that no valid values can be obtained.",
			},
			"sign_algo": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Signature algorithm.Note: This field may return null, indicating that no valid values can be obtained.",
			},
			"common_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain name of the certificate.Note: This field may return `null`, indicating that no valid value can be obtained.",
			},

			"mode": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Mode of configuring the certificate, the values are: `disable`: Do not configure the certificate; `eofreecert`: Configure EdgeOne free certificate; `sslcert`: Configure SSL certificate. If not filled in, the default value is disable.",
			},
		},
	}
}

func resourceTencentCloudTeoCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_certificate.create")()
	defer inconsistentCheck(d, meta)()

	var (
		zoneId string
		host   string
		certId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
	}
	if v, ok := d.GetOk("cert_id"); ok {
		certId = v.(string)
	}

	d.SetId(zoneId + FILED_SP + host + FILED_SP + certId)

	return resourceTencentCloudTeoCertificateUpdate(d, meta)
}

func resourceTencentCloudTeoCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_certificate.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	host := idSplit[1]
	// certId := idSplit[2]

	certificate, err := service.DescribeTeoAccelerationDomainById(ctx, zoneId, host)
	if err != nil {
		return err
	}

	if certificate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoCertificate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if certificate.ZoneId != nil {
		_ = d.Set("zone_id", certificate.ZoneId)
	}

	if certificate.DomainName != nil {
		_ = d.Set("hosts", certificate.DomainName)
	}

	// if certificate.ServerCertInfo != nil {
	// 	for _, serverCertInfo := range certificate.ServerCertInfo {
	// 		if serverCertInfo.CertId != nil && serverCertInfo.CertId == certId {
	// 			_ = d.Set("cert_id", serverCertInfo.CertId)

	// 			if serverCertInfo.Alias != nil {
	// 				_ = d.Set("alias", serverCertInfo.Alias)
	// 			}

	// 			if serverCertInfo.Type != nil {
	// 				_ = d.Set("type", serverCertInfo.Type)
	// 			}

	// 			if serverCertInfo.ExpireTime != nil {
	// 				_ = d.Set("expire_time", serverCertInfo.ExpireTime)
	// 			}

	// 			if serverCertInfo.DeployTime != nil {
	// 				_ = d.Set("deploy_time", serverCertInfo.DeployTime)
	// 			}

	// 			if serverCertInfo.SignAlgo != nil {
	// 				_ = d.Set("sign_algo", serverCertInfo.SignAlgo)
	// 			}

	// 			if serverCertInfo.CommonName != nil {
	// 				_ = d.Set("common_name", serverCertInfo.CommonName)
	// 			}
	// 		}
	// 	}
	// }

	// if certificate.Mode != nil {
	// 	_ = d.Set("mode", certificate.Mode)
	// }

	return nil
}

func resourceTencentCloudTeoCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_certificate.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := teo.NewModifyHostsCertificateRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	host := idSplit[1]
	certId := idSplit[2]

	request.ZoneId = &zoneId
	request.Hosts = []*string{&host}

	serverCertInfo := teo.ServerCertInfo{}
	serverCertInfo.CertId = &certId

	if v, ok := d.GetOk("alias"); ok {
		serverCertInfo.Alias = helper.String(v.(string))
	}
	if v, ok := d.GetOk("type"); ok {
		serverCertInfo.Type = helper.String(v.(string))
	}
	if v, ok := d.GetOk("expire_time"); ok {
		serverCertInfo.ExpireTime = helper.String(v.(string))
	}
	if v, ok := d.GetOk("deploy_time"); ok {
		serverCertInfo.DeployTime = helper.String(v.(string))
	}
	if v, ok := d.GetOk("sign_algo"); ok {
		serverCertInfo.SignAlgo = helper.String(v.(string))
	}
	if v, ok := d.GetOk("common_name"); ok {
		serverCertInfo.CommonName = helper.String(v.(string))
	}
	request.ServerCertInfo = append(request.ServerCertInfo, &serverCertInfo)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyHostsCertificate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update teo certificate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoCertificateRead(d, meta)
}

func resourceTencentCloudTeoCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_certificate.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
