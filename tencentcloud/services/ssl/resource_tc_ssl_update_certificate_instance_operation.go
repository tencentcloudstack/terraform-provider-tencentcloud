package ssl

import (
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSslUpdateCertificateInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslUpdateCertificateInstanceOperationCreate,
		Read:   resourceTencentCloudSslUpdateCertificateInstanceOperationRead,
		Delete: resourceTencentCloudSslUpdateCertificateInstanceOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Update new certificate ID.",
			},

			"old_certificate_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Update the original certificate ID.",
			},

			"resource_types": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The resource type that needs to be deployed. The parameter value is optional: clb, cdn, waf, live, ddos, teo, apigateway, vod, tke, tcb.",
			},

			"resource_types_regions": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "List of regions where cloud resources need to be deploye.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cloud resource type.",
						},
						"regions": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Region list.",
						},
					},
				},
			},

			"certificate_public_key": {
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				Description: "Certificate public key. If you upload the certificate public key, CertificateId does not need to be passed.",
			},

			"certificate_private_key": {
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				Description: "Certificate private key. If you upload the certificate public key, CertificateId does not need to be passed.",
			},

			"expiring_notification_switch": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Whether to ignore expiration reminders for old certificates 0: Do not ignore notifications. 1: Ignore the notification and ignore the OldCertificateId expiration reminder.",
			},

			"repeatable": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether the same certificate is allowed to be uploaded repeatedly. If you choose to upload the certificate, you can configure this parameter.",
			},

			"allow_download": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to allow downloading, if you choose to upload the certificate, you can configure this parameter.",
			},

			"project_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Project ID, if you choose to upload the certificate, you can configure this parameter.",
			},
		},
	}
}

func resourceTencentCloudSslUpdateCertificateInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_update_certificate_instance_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = ssl.NewUpdateCertificateInstanceRequest()
		response       = ssl.NewUpdateCertificateInstanceResponse()
		deployRecordId uint64
	)
	if v, ok := d.GetOk("certificate_id"); ok {
		request.CertificateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("old_certificate_id"); ok {
		request.OldCertificateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_types"); ok {
		resourceTypesSet := v.(*schema.Set).List()
		for i := range resourceTypesSet {
			if resourceTypesSet[i] != nil {
				resourceTypes := resourceTypesSet[i].(string)
				request.ResourceTypes = append(request.ResourceTypes, &resourceTypes)
			}
		}
	}

	if v, ok := d.GetOk("resource_types_regions"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			resourceTypeRegions := ssl.ResourceTypeRegions{}
			if v, ok := dMap["resource_type"]; ok {
				resourceTypeRegions.ResourceType = helper.String(v.(string))
			}
			if v, ok := dMap["regions"]; ok {
				regionsSet := v.(*schema.Set).List()
				for i := range regionsSet {
					if regionsSet[i] != nil {
						regions := regionsSet[i].(string)
						resourceTypeRegions.Regions = append(resourceTypeRegions.Regions, &regions)
					}
				}
			}
			request.ResourceTypesRegions = append(request.ResourceTypesRegions, &resourceTypeRegions)
		}
	}

	if v, ok := d.GetOk("certificate_public_key"); ok {
		request.CertificatePublicKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("certificate_private_key"); ok {
		request.CertificatePrivateKey = helper.String(v.(string))
	}

	if v, _ := d.GetOkExists("expiring_notification_switch"); v != nil {
		request.ExpiringNotificationSwitch = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOkExists("repeatable"); v != nil {
		request.Repeatable = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOkExists("allow_download"); v != nil {
		request.AllowDownload = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOkExists("project_id"); v != nil {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient().UpdateCertificateInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		if response == nil || response.Response == nil || response.Response.DeployRecordId == nil {
			return resource.RetryableError(fmt.Errorf("operate ssl updateCertificateInstanceOperation response is null"))
		}
		if *response.Response.DeployRecordId == uint64(0) {
			return resource.RetryableError(fmt.Errorf("operate ssl updateCertificateInstanceOperation not done"))
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate ssl updateCertificateInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	deployRecordId = *response.Response.DeployRecordId
	d.SetId(helper.UInt64ToStr(deployRecordId))

	return resourceTencentCloudSslUpdateCertificateInstanceOperationRead(d, meta)
}

func resourceTencentCloudSslUpdateCertificateInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_update_certificate_instance_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslUpdateCertificateInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_update_certificate_instance_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
