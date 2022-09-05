/*
Provides a resource to create a teo hostCertificate

Example Usage

```hcl
resource "tencentcloud_teo_host_certificate" "host_certificate" {
  zone_id = tencentcloud_teo_zone.zone.id
  host    = tencentcloud_teo_dns_record.dns_record.name

  cert_info {
    cert_id = "yqWPPbs7"
    status  = "deployed"
  }
}

```
Import

teo hostCertificate can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_host_certificate.host_certificate zoneId#host
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220106"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoHostCertificate() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoHostCertificateRead,
		Create: resourceTencentCloudTeoHostCertificateCreate,
		Update: resourceTencentCloudTeoHostCertificateUpdate,
		Delete: resourceTencentCloudTeoHostCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name.",
			},

			"cert_info": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Server certificate configuration.Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Server certificate ID, which is the ID of the default certificate. If you choose to upload an external certificate for SSL certificate management, a certificate ID will be generated.",
						},
						"alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alias of the certificate.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate type.- default: Default certificate.- upload: External certificate.- managed: Tencent Cloud managed certificate.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time when the certificate expires.Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Certificate deployment status.- processing: Deploying- deployed: DeployedNote: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoHostCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_host_certificate.create")()
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
	return resourceTencentCloudTeoHostCertificateUpdate(d, meta)
}

func resourceTencentCloudTeoHostCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_host_certificate.read")()
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

	hostCertificate, err := service.DescribeTeoHostCertificate(ctx, zoneId, host)

	if err != nil {
		return err
	}

	if hostCertificate == nil {
		d.SetId("")
		return fmt.Errorf("resource `hostCertificate` %s does not exist", d.Id())
	}

	if hostCertificate.Host != nil {
		_ = d.Set("host", hostCertificate.Host)
	}

	if hostCertificate.CertInfo != nil {
		certInfoList := []interface{}{}
		for _, certInfo := range hostCertificate.CertInfo {
			certInfoMap := map[string]interface{}{}
			if certInfo.CertId != nil {
				certInfoMap["cert_id"] = certInfo.CertId
			}
			if certInfo.Alias != nil {
				certInfoMap["alias"] = certInfo.Alias
			}
			if certInfo.Type != nil {
				certInfoMap["type"] = certInfo.Type
			}
			if certInfo.ExpireTime != nil {
				certInfoMap["expire_time"] = certInfo.ExpireTime
			}
			if certInfo.Status != nil {
				certInfoMap["status"] = certInfo.Status
			}

			certInfoList = append(certInfoList, certInfoMap)
		}
		_ = d.Set("cert_info", certInfoList)
	}

	return nil
}

func resourceTencentCloudTeoHostCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_host_certificate.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := teo.NewModifyHostsCertificateRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	host := idSplit[1]

	request.ZoneId = helper.String(zoneId)
	request.Hosts = []*string{helper.String(host)}

	if v, ok := d.GetOk("cert_info"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			serverCertInfo := teo.ServerCertInfo{}
			if v, ok := dMap["cert_id"]; ok {
				serverCertInfo.CertId = helper.String(v.(string))
			}
			if v, ok := dMap["status"]; ok {
				serverCertInfo.Status = helper.String(v.(string))
			}

			request.CertInfo = append(request.CertInfo, &serverCertInfo)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyHostsCertificate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo hostCertificate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTeoHostCertificateRead(d, meta)
}

func resourceTencentCloudTeoHostCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_host_certificate.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
