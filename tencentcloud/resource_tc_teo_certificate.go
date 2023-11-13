/*
Provides a resource to create a teo certificate

Example Usage

```hcl
resource "tencentcloud_teo_certificate" "certificate" {
  zone_id = ""
  hosts =
  server_cert_info {
		cert_id = ""
		alias = ""
		type = ""
		expire_time = ""
		deploy_time = ""
		sign_algo = ""
		common_name = ""

  }
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
	teo "github.com/TencentCloud/tencentcloud-sdk-go-intl-en/tencentcloud/teo/v20220901"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
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
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Acceleration domain name that needs to modify the certificate configuration.",
			},

			"server_cert_info": {
				Optional:    true,
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
				Type:        schema.TypeString,
				Description: "Mode of configuring the certificate, the values are:&amp;amp;lt;li&amp;amp;gt;disable: Do not configure the certificate;&amp;amp;lt;/lil&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;eofreecert: Configure EdgeOne free certificate;&amp;amp;lt;/lil&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;sslcert: Configure SSL certificate.&amp;amp;lt;/lil&amp;amp;gt;If not filled in, the default value is disable.",
			},
		},
	}
}

func resourceTencentCloudTeoCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_certificate.create")()
	defer inconsistentCheck(d, meta)()

	var zoneId string
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	var hosts string
	if v, ok := d.GetOk("hosts"); ok {
		hosts = v.(string)
	}

	d.SetId(strings.Join([]string{zoneId, hosts}, FILED_SP))

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
	hosts := idSplit[1]

	certificate, err := service.DescribeTeoCertificateById(ctx, zoneId, hosts)
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

	if certificate.Hosts != nil {
		_ = d.Set("hosts", certificate.Hosts)
	}

	if certificate.ServerCertInfo != nil {
		serverCertInfoList := []interface{}{}
		for _, serverCertInfo := range certificate.ServerCertInfo {
			serverCertInfoMap := map[string]interface{}{}

			if certificate.ServerCertInfo.CertId != nil {
				serverCertInfoMap["cert_id"] = certificate.ServerCertInfo.CertId
			}

			if certificate.ServerCertInfo.Alias != nil {
				serverCertInfoMap["alias"] = certificate.ServerCertInfo.Alias
			}

			if certificate.ServerCertInfo.Type != nil {
				serverCertInfoMap["type"] = certificate.ServerCertInfo.Type
			}

			if certificate.ServerCertInfo.ExpireTime != nil {
				serverCertInfoMap["expire_time"] = certificate.ServerCertInfo.ExpireTime
			}

			if certificate.ServerCertInfo.DeployTime != nil {
				serverCertInfoMap["deploy_time"] = certificate.ServerCertInfo.DeployTime
			}

			if certificate.ServerCertInfo.SignAlgo != nil {
				serverCertInfoMap["sign_algo"] = certificate.ServerCertInfo.SignAlgo
			}

			if certificate.ServerCertInfo.CommonName != nil {
				serverCertInfoMap["common_name"] = certificate.ServerCertInfo.CommonName
			}

			serverCertInfoList = append(serverCertInfoList, serverCertInfoMap)
		}

		_ = d.Set("server_cert_info", serverCertInfoList)

	}

	if certificate.Mode != nil {
		_ = d.Set("mode", certificate.Mode)
	}

	return nil
}

func resourceTencentCloudTeoCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_certificate.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := teo.NewModifyHostsCertificateRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	hosts := idSplit[1]

	request.ZoneId = &zoneId
	request.Hosts = &hosts

	immutableArgs := []string{"zone_id", "hosts", "server_cert_info", "mode"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
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

	return resourceTencentCloudTeoCertificateRead(d, meta)
}

func resourceTencentCloudTeoCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_certificate.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
