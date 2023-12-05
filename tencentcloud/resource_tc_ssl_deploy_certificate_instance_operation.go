package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSslDeployCertificateInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslDeployCertificateInstanceCreate,
		Read:   resourceTencentCloudSslDeployCertificateInstanceRead,
		Delete: resourceTencentCloudSslDeployCertificateInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the certificate to be deployed.",
			},

			"instance_id_list": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Need to deploy instance list.",
			},

			"resource_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Deployed cloud resource type.",
			},

			"status": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Deployment cloud resource status: Live: -1: The domain name is not associated with a certificate.1:  Domain name https is enabled.0:  Domain name https is closed.",
			},
		},
	}
}

func resourceTencentCloudSslDeployCertificateInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_deploy_certificate_instance_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = ssl.NewDeployCertificateInstanceRequest()
		response       = ssl.NewDeployCertificateInstanceResponse()
		deployRecordId uint64
	)
	if v, ok := d.GetOk("certificate_id"); ok {
		request.CertificateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id_list"); ok {
		instanceIdListSet := v.(*schema.Set).List()
		for i := range instanceIdListSet {
			instanceIdList := instanceIdListSet[i].(string)
			request.InstanceIdList = append(request.InstanceIdList, &instanceIdList)
		}
	}

	if v, ok := d.GetOk("resource_type"); ok {
		request.ResourceType = helper.String(v.(string))
	}

	if v, _ := d.GetOk("status"); v != nil {
		request.Status = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSSLCertificateClient().DeployCertificateInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl deployCertificateInstance failed, reason:%+v", logId, err)
		return err
	}

	deployRecordId = *response.Response.DeployRecordId
	d.SetId(helper.UInt64ToStr(deployRecordId))

	return resourceTencentCloudSslDeployCertificateInstanceRead(d, meta)
}

func resourceTencentCloudSslDeployCertificateInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_deploy_certificate_instance_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslDeployCertificateInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_deploy_certificate_instance_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
