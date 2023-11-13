/*
Provides a resource to create a teo host_certificate

Example Usage

```hcl
resource "tencentcloud_teo_host_certificate" "host_certificate" {
  zone_id = &lt;nil&gt;
  host = &lt;nil&gt;
  cert_info {
		cert_id = &lt;nil&gt;
		status = &lt;nil&gt;

  }
}
```

Import

teo host_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_teo_host_certificate.host_certificate host_certificate_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTeoHostCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoHostCertificateCreate,
		Read:   resourceTencentCloudTeoHostCertificateRead,
		Delete: resourceTencentCloudTeoHostCertificateDelete,
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
				Description: "Domain name.",
			},

			"cert_info": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Server certificate configuration. Note: This field may return null, indicating that no valid value can be obtained.",
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
							Description: "Alias of the certificate. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate type.- `default`: Default certificate.- `upload`: External certificate.- `managed`: Tencent Cloud managed certificate. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time when the certificate expires. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"effective_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time when the certificate takes effect. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Certificate deployment status.- `processing`: Deploying- `deployed`: Deployed Note: This field may return null, indicating that no valid value can be obtained.",
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

	logId := getLogId(contextNil)

	var (
		request  = teo.NewModifyHostsCertificateRequest()
		response = teo.NewModifyHostsCertificateResponse()
		zoneId   string
		host     string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host"); ok {
		host = v.(string)
		request.Host = helper.String(v.(string))
	}

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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo hostCertificate failed, reason:%+v", logId, err)
		return err
	}

	zoneId = *response.Response.ZoneId
	d.SetId(strings.Join([]string{zoneId, host}, FILED_SP))

	return resourceTencentCloudTeoHostCertificateRead(d, meta)
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

	hostCertificate, err := service.DescribeTeoHostCertificateById(ctx, zoneId, host)
	if err != nil {
		return err
	}

	if hostCertificate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoHostCertificate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if hostCertificate.ZoneId != nil {
		_ = d.Set("zone_id", hostCertificate.ZoneId)
	}

	if hostCertificate.Host != nil {
		_ = d.Set("host", hostCertificate.Host)
	}

	if hostCertificate.CertInfo != nil {
		certInfoList := []interface{}{}
		for _, certInfo := range hostCertificate.CertInfo {
			certInfoMap := map[string]interface{}{}

			if hostCertificate.CertInfo.CertId != nil {
				certInfoMap["cert_id"] = hostCertificate.CertInfo.CertId
			}

			if hostCertificate.CertInfo.Alias != nil {
				certInfoMap["alias"] = hostCertificate.CertInfo.Alias
			}

			if hostCertificate.CertInfo.Type != nil {
				certInfoMap["type"] = hostCertificate.CertInfo.Type
			}

			if hostCertificate.CertInfo.ExpireTime != nil {
				certInfoMap["expire_time"] = hostCertificate.CertInfo.ExpireTime
			}

			if hostCertificate.CertInfo.EffectiveTime != nil {
				certInfoMap["effective_time"] = hostCertificate.CertInfo.EffectiveTime
			}

			if hostCertificate.CertInfo.Status != nil {
				certInfoMap["status"] = hostCertificate.CertInfo.Status
			}

			certInfoList = append(certInfoList, certInfoMap)
		}

		_ = d.Set("cert_info", certInfoList)

	}

	return nil
}

func resourceTencentCloudTeoHostCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_host_certificate.delete")()
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

	if err := service.DeleteTeoHostCertificateById(ctx, zoneId, host); err != nil {
		return err
	}

	return nil
}
