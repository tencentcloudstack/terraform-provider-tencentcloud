package tencentcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSslUpdateCertificateInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslUpdateCertificateInstanceCreate,
		Read:   resourceTencentCloudSslUpdateCertificateInstanceRead,
		Delete: resourceTencentCloudSslUpdateCertificateInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
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
				Description: "The resource type that needs to be deployed. The parameter value is optional: clb,cdn,waf,live,ddos,teo,apigateway,vod,tke,tcb.",
			},

			"resource_types_regions": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "List of regions where cloud resources need to be deployed.",
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
		},
	}
}

func resourceTencentCloudSslUpdateCertificateInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_update_certificate_instance_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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
			resourceTypes := resourceTypesSet[i].(string)
			request.ResourceTypes = append(request.ResourceTypes, &resourceTypes)
		}
	}

	if v, ok := d.GetOk("resource_types_regions"); ok {
		for _, item := range v.([]interface{}) {
			resourceTypeRegions := ssl.ResourceTypeRegions{}
			dMap := item.(map[string]interface{})
			if v, ok := dMap["resource_type"]; ok {
				resourceTypeRegions.ResourceType = helper.String(v.(string))
			}
			if v, ok := dMap["regions"]; ok {
				regionsSet := v.(*schema.Set).List()
				for i := range regionsSet {
					regions := regionsSet[i].(string)
					resourceTypeRegions.Regions = append(resourceTypeRegions.Regions, &regions)
				}
			}
			request.ResourceTypesRegions = append(request.ResourceTypesRegions, &resourceTypeRegions)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSSLCertificateClient().UpdateCertificateInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result.Response.DeployStatus == nil {
			err := fmt.Errorf("api[%s] status is nil", request.GetAction())
			return retryError(err)
		}

		if *result.Response.DeployStatus < 0 {
			err := fmt.Errorf("status is %d, need retry", *result.Response.DeployStatus)
			return retryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl updateCertificateInstance failed, reason:%+v", logId, err)
		return err
	}

	deployRecordId = *response.Response.DeployRecordId
	d.SetId(helper.UInt64ToStr(deployRecordId))

	return resourceTencentCloudSslUpdateCertificateInstanceRead(d, meta)
}

func resourceTencentCloudSslUpdateCertificateInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_update_certificate_instance_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslUpdateCertificateInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_update_certificate_instance_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
