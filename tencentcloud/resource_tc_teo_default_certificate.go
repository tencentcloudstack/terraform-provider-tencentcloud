/*
Provides a resource to create a teo default_certificate

Example Usage

```hcl
resource "tencentcloud_teo_default_certificate" "default_certificate" {
  zone_id = ""
  cert_info {
			cert_id = ""
			status = ""

  }
}

```
Import

teo default_certificate can be imported using the id, e.g.
```
$ terraform import tencentcloud_teo_default_certificate.default_certificate defaultCertificate_id
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

func resourceTencentCloudTeoDefaultCertificate() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoDefaultCertificateRead,
		Create: resourceTencentCloudTeoDefaultCertificateCreate,
		Update: resourceTencentCloudTeoDefaultCertificateUpdate,
		Delete: resourceTencentCloudTeoDefaultCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"cert_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
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

	var (
		zoneId string
		certId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("cert_info"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			if v, ok := dMap["cert_id"]; ok {
				certId = v.(string)
				break
			}
		}
	}

	d.SetId(zoneId + FILED_SP + certId)
	return resourceTencentCloudTeoDefaultCertificateUpdate(d, meta)
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

	defaultCertificate, err := service.DescribeTeoDefaultCertificate(ctx, zoneId, certId)

	if err != nil {
		return err
	}

	if defaultCertificate == nil {
		d.SetId("")
		return fmt.Errorf("resource `defaultCertificate` %s does not exist", certId)
	}

	_ = d.Set("zone_id", zoneId)

	if defaultCertificate != nil {
		certInfoList := []interface{}{}
		certInfoMap := map[string]interface{}{}
		if defaultCertificate.CertId != nil {
			certInfoMap["cert_id"] = defaultCertificate.CertId
		}
		if defaultCertificate.Alias != nil {
			certInfoMap["alias"] = defaultCertificate.Alias
		}
		if defaultCertificate.Type != nil {
			certInfoMap["type"] = defaultCertificate.Type
		}
		if defaultCertificate.ExpireTime != nil {
			certInfoMap["expire_time"] = defaultCertificate.ExpireTime
		}
		if defaultCertificate.EffectiveTime != nil {
			certInfoMap["effective_time"] = defaultCertificate.EffectiveTime
		}
		if defaultCertificate.CommonName != nil {
			certInfoMap["common_name"] = defaultCertificate.CommonName
		}
		if defaultCertificate.SubjectAltName != nil {
			certInfoMap["subject_alt_name"] = defaultCertificate.SubjectAltName
		}
		if defaultCertificate.Status != nil {
			certInfoMap["status"] = defaultCertificate.Status
		}
		if defaultCertificate.Message != nil {
			certInfoMap["message"] = defaultCertificate.Message
		}
		certInfoList = append(certInfoList, certInfoMap)
		_ = d.Set("cert_info", certInfoList)
	}

	return nil
}

func resourceTencentCloudTeoDefaultCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_default_certificate.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := teo.NewModifyDefaultCertificateRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	certId := idSplit[1]

	request.ZoneId = &zoneId
	request.CertId = &certId

	if d.HasChange("cert_info") {
		if v, ok := d.GetOk("cert_info"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				if v, ok := dMap["cert_id"]; ok {
					request.CertId = helper.String(v.(string))
				}
				if v, ok := dMap["status"]; ok {
					request.Status = helper.String(v.(string))
				}
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().ModifyDefaultCertificate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo defaultCertificate failed, reason:%+v", logId, err)
		return err
	}

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	err = resource.Retry(60*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeTeoDefaultCertificate(ctx, zoneId, certId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.Status == *request.Status {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("defaultCertificate status is %v, retry...", *instance.Status))
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTeoDefaultCertificateRead(d, meta)
}

func resourceTencentCloudTeoDefaultCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_default_certificate.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
