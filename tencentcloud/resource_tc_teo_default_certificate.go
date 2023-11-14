/*
Provides a resource to create a teo default_certificate

Example Usage

```hcl
resource "tencentcloud_teo_default_certificate" "default_certificate" {
  zone_id = &lt;nil&gt;
  cert_info {
		cert_id = &lt;nil&gt;
		status = &lt;nil&gt;

  }
}
```

Import

teo default_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_teo_default_certificate.default_certificate default_certificate_id
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
	"time"
)

func resourceTencentCloudTeoDefaultCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoDefaultCertificateCreate,
		Read:   resourceTencentCloudTeoDefaultCertificateRead,
		Delete: resourceTencentCloudTeoDefaultCertificateDelete,
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

			"cert_info": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "List of default certificates. Note: This field may return null, indicating that no valid value can be obtained.",
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
							Description: "Certificate alias. Note: This field may return null, indicating that no valid value can be obtained.",
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
						"common_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate common name. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"subject_alt_name": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Domain names added to the SAN certificate. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Certificate status.- `applying`: Application in progress.- `failed`: Application failed.- `processing`: Deploying certificate.- `deployed`: Certificate deployed.- `disabled`: Certificate disabled. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Returns a message to display failure causes when `Status` is failed. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoDefaultCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_default_certificate.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = teo.NewModifyDefaultCertificateRequest()
		response = teo.NewModifyDefaultCertificateResponse()
		zoneId   string
		certId   string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cert_info"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			defaultServerCertInfo := teo.DefaultServerCertInfo{}
			if v, ok := dMap["cert_id"]; ok {
				defaultServerCertInfo.CertId = helper.String(v.(string))
			}
			if v, ok := dMap["status"]; ok {
				defaultServerCertInfo.Status = helper.String(v.(string))
			}
			request.CertInfo = append(request.CertInfo, &defaultServerCertInfo)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyDefaultCertificate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo defaultCertificate failed, reason:%+v", logId, err)
		return err
	}

	zoneId = *response.Response.ZoneId
	d.SetId(strings.Join([]string{zoneId, certId}, FILED_SP))

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"deployed"}, 60*readRetryTimeout, time.Second, service.TeoDefaultCertificateStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudTeoDefaultCertificateRead(d, meta)
}

func resourceTencentCloudTeoDefaultCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_default_certificate.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	certId := idSplit[1]

	defaultCertificate, err := service.DescribeTeoDefaultCertificateById(ctx, zoneId, certId)
	if err != nil {
		return err
	}

	if defaultCertificate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoDefaultCertificate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if defaultCertificate.ZoneId != nil {
		_ = d.Set("zone_id", defaultCertificate.ZoneId)
	}

	if defaultCertificate.CertInfo != nil {
		certInfoList := []interface{}{}
		for _, certInfo := range defaultCertificate.CertInfo {
			certInfoMap := map[string]interface{}{}

			if defaultCertificate.CertInfo.CertId != nil {
				certInfoMap["cert_id"] = defaultCertificate.CertInfo.CertId
			}

			if defaultCertificate.CertInfo.Alias != nil {
				certInfoMap["alias"] = defaultCertificate.CertInfo.Alias
			}

			if defaultCertificate.CertInfo.Type != nil {
				certInfoMap["type"] = defaultCertificate.CertInfo.Type
			}

			if defaultCertificate.CertInfo.ExpireTime != nil {
				certInfoMap["expire_time"] = defaultCertificate.CertInfo.ExpireTime
			}

			if defaultCertificate.CertInfo.EffectiveTime != nil {
				certInfoMap["effective_time"] = defaultCertificate.CertInfo.EffectiveTime
			}

			if defaultCertificate.CertInfo.CommonName != nil {
				certInfoMap["common_name"] = defaultCertificate.CertInfo.CommonName
			}

			if defaultCertificate.CertInfo.SubjectAltName != nil {
				certInfoMap["subject_alt_name"] = defaultCertificate.CertInfo.SubjectAltName
			}

			if defaultCertificate.CertInfo.Status != nil {
				certInfoMap["status"] = defaultCertificate.CertInfo.Status
			}

			if defaultCertificate.CertInfo.Message != nil {
				certInfoMap["message"] = defaultCertificate.CertInfo.Message
			}

			certInfoList = append(certInfoList, certInfoMap)
		}

		_ = d.Set("cert_info", certInfoList)

	}

	return nil
}

func resourceTencentCloudTeoDefaultCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_default_certificate.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	certId := idSplit[1]

	if err := service.DeleteTeoDefaultCertificateById(ctx, zoneId, certId); err != nil {
		return err
	}

	return nil
}
