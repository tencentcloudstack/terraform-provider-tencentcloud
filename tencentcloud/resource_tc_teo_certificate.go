/*
Provides a resource to create a teo certificate

Example Usage

```hcl
resource "tencentcloud_teo_certificate" "certificate" {
  host    = "test.tencentcloud-terraform-provider.cn"
  mode    = "eofreecert"
  zone_id = "zone-2o1t24kgy362"
}
```

Configure SSL certificate

```hcl
resource "tencentcloud_teo_certificate" "certificate" {
    host    = "test.tencentcloud-terraform-provider.cn"
    mode    = "sslcert"
    zone_id = "zone-2o1t24kgy362"

    server_cert_info {
        alias       = "EdgeOne default"
        cert_id     = "teo-2o1tfutpnb6l"
        common_name = "tencentcloud-terraform-provider.cn"
        deploy_time = "2023-09-27T11:54:47Z"
        expire_time = "2023-12-26T06:38:47Z"
        sign_algo   = "RSA 2048"
        type        = "default"
    }
}
```

Import

teo certificate can be imported using the id, e.g.

```
terraform import tencentcloud_teo_certificate.certificate zone_id#host#cert_id
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
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"host": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Acceleration domain name that needs to modify the certificate configuration.",
			},

			"server_cert_info": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Description: "SSL certificate configuration, this parameter takes effect only when mode = sslcert, just enter the corresponding CertId. You can go to the SSL certificate list to view the CertId.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description: "Type of the certificate. Values:&amp;lt;li&amp;gt;`default`: Default certificate&amp;lt;/lil&amp;gt;&amp;lt;li&amp;gt;`upload`: Specified certificate&amp;lt;/li&amp;gt;&amp;lt;li&amp;gt;`managed`: Tencent Cloud-managed certificate&amp;lt;/li&amp;gt;Note: This field may return `null`, indicating that no valid value can be obtained.",
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
					},
				},
			},

			"mode": {
				Optional:    true,
				Computed:    true,
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
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
	}

	d.SetId(zoneId + FILED_SP + host)

	return resourceTencentCloudTeoCertificateUpdate(d, meta)
}

func resourceTencentCloudTeoCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_certificate.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	host := idSplit[1]

	accelerationDomain, err := service.DescribeTeoAccelerationDomainById(ctx, zoneId, host)
	if err != nil {
		return err
	}

	if accelerationDomain == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoCertificate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if accelerationDomain.ZoneId != nil {
		_ = d.Set("zone_id", accelerationDomain.ZoneId)
	}

	if accelerationDomain.DomainName != nil {
		_ = d.Set("host", accelerationDomain.DomainName)
	}

	if accelerationDomain.Certificate != nil {
		certificate := accelerationDomain.Certificate
		zone, err := service.DescribeTeoZone(ctx, zoneId)
		if err != nil {
			return err
		}

		serverCertInfoList := []interface{}{}
		for _, serverCertInfo := range certificate.List {
			serverCertInfoMap := map[string]interface{}{}

			if serverCertInfo.CertId != nil {
				serverCertInfoMap["cert_id"] = serverCertInfo.CertId
			}

			if serverCertInfo.Alias != nil {
				serverCertInfoMap["alias"] = serverCertInfo.Alias
			}

			if serverCertInfo.Type != nil {
				serverCertInfoMap["type"] = serverCertInfo.Type
			}

			if serverCertInfo.ExpireTime != nil {
				serverCertInfoMap["expire_time"] = serverCertInfo.ExpireTime
			}

			if serverCertInfo.DeployTime != nil {
				serverCertInfoMap["deploy_time"] = serverCertInfo.DeployTime
			}

			if serverCertInfo.SignAlgo != nil {
				serverCertInfoMap["sign_algo"] = serverCertInfo.SignAlgo
			}

			if zone.ZoneName != nil {
				serverCertInfoMap["common_name"] = zone.ZoneName
			}

			serverCertInfoList = append(serverCertInfoList, serverCertInfoMap)
		}

		_ = d.Set("server_cert_info", serverCertInfoList)

		if certificate.Mode != nil {
			_ = d.Set("mode", certificate.Mode)
		}
	}

	return nil
}

func resourceTencentCloudTeoCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_certificate.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := teo.NewModifyHostsCertificateRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	host := idSplit[1]

	request.ZoneId = &zoneId
	request.Hosts = []*string{&host}

	if v, ok := d.GetOk("server_cert_info"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			serverCertInfo := teo.ServerCertInfo{}
			if v, ok := dMap["cert_id"]; ok {
				serverCertInfo.CertId = helper.String(v.(string))
			}
			if v, ok := dMap["alias"]; ok {
				serverCertInfo.Alias = helper.String(v.(string))
			}
			if v, ok := dMap["type"]; ok {
				serverCertInfo.Type = helper.String(v.(string))
			}
			if v, ok := dMap["expire_time"]; ok {
				serverCertInfo.ExpireTime = helper.String(v.(string))
			}
			if v, ok := dMap["deploy_time"]; ok {
				serverCertInfo.DeployTime = helper.String(v.(string))
			}
			if v, ok := dMap["sign_algo"]; ok {
				serverCertInfo.SignAlgo = helper.String(v.(string))
			}
			if v, ok := dMap["common_name"]; ok {
				serverCertInfo.CommonName = helper.String(v.(string))
			}
			request.ServerCertInfo = append(request.ServerCertInfo, &serverCertInfo)
		}
	}

	if v, ok := d.GetOk("mode"); ok {
		request.Mode = helper.String(v.(string))
	}

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

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = service.CheckAccelerationDomainStatus(ctx, zoneId, host, "")
	if err != nil {
		return err
	}

	return resourceTencentCloudTeoCertificateRead(d, meta)
}

func resourceTencentCloudTeoCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_certificate.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
