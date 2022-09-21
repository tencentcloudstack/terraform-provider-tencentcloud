/*
Provides a resource to create a teo host_certificate

Example Usage

```hcl
resource "tencentcloud_teo_host_certificate" "vstest_sfurnace_work" {
 zone_id = tencentcloud_teo_zone.sfurnace_work.id
 host    = tencentcloud_teo_dns_record.vstest_sfurnace_work.name

 cert_info {
   cert_id = "yqWPPbs7"
   status  = "deployed"
 }
}

```
Import

teo host_certificate can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_host_certificate.host_certificate hostCertificate_id
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
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
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

	if v, ok := d.GetOk("cert_info"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			if v, ok := dMap["cert_id"]; ok {
				certId = v.(string)
			}
			if v, ok := dMap["status"]; ok {
				if v.(string) != "" {
					return fmt.Errorf("[CRITAL] create teo hostCertificate status error")
				}
			}
		}
	}

	err := resourceTencentCloudTeoHostCertificateUpdate(d, meta)
	if err != nil {
		log.Printf("[CRITAL]%s create teo hostCertificate failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(zoneId + FILED_SP + host + FILED_SP + certId)
	return resourceTencentCloudTeoHostCertificateRead(d, meta)
}

func resourceTencentCloudTeoHostCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_host_certificate.read")()
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
	cateId := idSplit[2]

	hostCertificate, err := service.DescribeTeoHostCertificate(ctx, zoneId, host, cateId)

	if err != nil {
		return err
	}

	if hostCertificate == nil {
		d.SetId("")
		return fmt.Errorf("resource `hostCertificate` %s does not exist", cateId)
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("host", host)

	if hostCertificate != nil {
		certInfoList := []interface{}{}
		for _, certificate := range hostCertificate {
			certInfo := certificate.HostCertInfo
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
			//if certInfo.EffectiveTime != nil {
			//	certInfoMap["effective_time"] = certInfo.EffectiveTime
			//}
			if certInfo.DeployTime != nil {
				certInfoMap["deploy_time"] = certInfo.DeployTime
			}
			if certInfo.SignAlgo != nil {
				certInfoMap["sign_algo"] = certInfo.SignAlgo
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
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	host := idSplit[1]
	//cateId := idSplit[2]

	request.ZoneId = &zoneId
	request.Hosts = []*string{&host}

	if d.HasChange("zone_id") {
		return fmt.Errorf("`zone_id` do not support change now.")
	}

	if d.HasChange("host") {
		return fmt.Errorf("`host` do not support change now.")
	}

	if d.HasChange("cert_info") {
		if v, ok := d.GetOk("cert_info"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				serverCertInfo := teo.ServerCertInfo{}
				if v, ok := dMap["cert_id"]; ok {
					serverCertInfo.CertId = helper.String(v.(string))
				}
				//if v, ok := dMap["status"]; ok {
				//	serverCertInfo.Status = helper.String(v.(string))
				//}
				request.ServerCertInfo = append(request.ServerCertInfo, &serverCertInfo)
			}
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
