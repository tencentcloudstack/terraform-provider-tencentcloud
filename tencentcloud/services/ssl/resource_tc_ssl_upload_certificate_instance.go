package ssl

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudSslUploadCertificateInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslUploadCertificateInstanceCreate,
		Read:   resourceTencentCloudSslUploadCertificateInstanceRead,
		Update: resourceTencentCloudSslUploadCertificateInstanceUpdate,
		Delete: resourceTencentCloudSslUploadCertificateInstanceDelete,

		Schema: map[string]*schema.Schema{
			"old_certificate_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Update the original certificate ID. The cloud resources bound to this certificate ID are queried, and then these cloud resources are updated with the new certificate.",
			},

			"resource_types": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The resource type that needs to be deployed. The parameter value is optional: clb.",
			},

			"certificate_public_key": {
				Optional:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				Description: "Certificate public key. If you upload the certificate public key, CertificateId does not need to be passed.",
			},

			"certificate_private_key": {
				Optional:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				Description: "Certificate private key. If you upload the certificate public key, CertificateId does not need to be passed.",
			},

			"resource_types_regions": {
				Optional:    true,
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

			"certificate_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID.",
			},

			"deploy_record_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Cloud resource update task ID. 0 means the task is in progress. When DeployRecordId is greater than 0, the task is created successfully.",
			},
		},
	}
}

func resourceTencentCloudSslUploadCertificateInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_upload_certificate_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = ssl.NewUploadUpdateCertificateInstanceRequest()
		response       = ssl.NewUploadUpdateCertificateInstanceResponse()
		deployRecordId uint64
		certificateId  string
	)

	if v, ok := d.GetOk("old_certificate_id"); ok {
		request.OldCertificateId = helper.String(v.(string))
		certificateId = v.(string)
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

	if v, ok := d.GetOk("certificate_public_key"); ok {
		request.CertificatePublicKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("certificate_private_key"); ok {
		request.CertificatePrivateKey = helper.String(v.(string))
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

	if v, ok := d.GetOk("certificate_id"); ok {
		request.CertificateId = helper.String(v.(string))
		// Consult with the person in charge of the SSL interface; when `certificate_id` is provided, set `CertificatePublicKey` and `CertificatePrivateKey` to null values.
		request.CertificatePublicKey = helper.String("")
		request.CertificatePrivateKey = helper.String("")
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient().UploadUpdateCertificateInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		if result == nil || result.Response == nil || result.Response.DeployRecordId == nil {
			return resource.RetryableError(fmt.Errorf("operate ssl uploadCertificateInstance response is null"))
		}
		response = result
		if *response.Response.DeployRecordId == uint64(0) {
			return resource.RetryableError(fmt.Errorf("operate ssl uploadCertificateInstance not done"))
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate ssl uploadCertificateInstance failed, reason:%+v", logId, err)
		return err
	}

	deployRecordId = *response.Response.DeployRecordId

	// 异步任务，通过 DescribeHostUploadUpdateRecord 轮询任务状态，Status: 0 处理中, 1 完成
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var recordList []*ssl.UploadUpdateRecordInfo
		limit := uint64(200)
		offset := uint64(0)
		for {
			describeRequest := ssl.NewDescribeHostUploadUpdateRecordRequest()
			describeRequest.Limit = &limit
			describeRequest.Offset = &offset
			describeRequest.OldCertificateId = helper.String(certificateId)

			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient().DescribeHostUploadUpdateRecord(describeRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, describeRequest.GetAction(), describeRequest.ToJsonString(), result.ToJsonString())
			}
			if result == nil || result.Response == nil || len(result.Response.DeployRecordList) == 0 {
				break
			}
			recordList = append(recordList, result.Response.DeployRecordList...)
			if result.Response.TotalCount == nil || uint64(len(recordList)) >= *result.Response.TotalCount {
				break
			}
			offset += limit
		}

		var matchedRecord *ssl.UploadUpdateRecordInfo
		for _, record := range recordList {
			if record.Id != nil && *record.Id == deployRecordId {
				matchedRecord = record
				break
			}
		}
		if matchedRecord == nil || matchedRecord.Status == nil {
			return resource.RetryableError(fmt.Errorf("ssl uploadCertificateInstance deploy record %d not found", deployRecordId))
		}
		if *matchedRecord.Status == 0 {
			return resource.RetryableError(fmt.Errorf("ssl uploadCertificateInstance deploy record %d is still in progress", deployRecordId))
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s ssl uploadCertificateInstance deploy record %d check failed, reason:%+v", logId, deployRecordId, err)
		return err
	}

	log.Printf("[CRUD]%s ssl_upload_certificate_instance create success, logId=%s, id=%d", logId, logId, deployRecordId)
	d.SetId(strings.Join([]string{certificateId,
		helper.UInt64ToStr(deployRecordId)}, tccommon.FILED_SP))

	return resourceTencentCloudSslUploadCertificateInstanceRead(d, meta)
}

func resourceTencentCloudSslUploadCertificateInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_upload_certificate_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 4 {
		return fmt.Errorf("ids param is error. id:  %s", ids)
	}
	oldCertificateId := ids[0]
	recordId := ids[1]

	var (
		limit   uint64 = 200
		offset  uint64 = 0
		records []*ssl.UploadUpdateRecordInfo
	)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var innerRecords []*ssl.UploadUpdateRecordInfo
		for {
			request := ssl.NewDescribeHostUploadUpdateRecordRequest()
			request.Limit = &limit
			request.Offset = &offset
			request.OldCertificateId = helper.String(oldCertificateId)

			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient().DescribeHostUploadUpdateRecord(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || len(result.Response.DeployRecordList) == 0 {
				break
			}
			innerRecords = append(innerRecords, result.Response.DeployRecordList...)
			if result.Response.TotalCount == nil || uint64(len(innerRecords)) >= *result.Response.TotalCount {
				break
			}
			offset += limit
		}
		records = innerRecords
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read ssl uploadCertificateInstance failed, reason:%+v", logId, err)
		return err
	}

	if err != nil {
		log.Printf("[CRITAL]%s read ssl uploadCertificateInstance failed, reason:%+v", logId, err)
		return err
	}

	if len(records) == 0 {
		log.Printf("[DescribeHostUploadUpdateRecord] ssl_upload_certificate_instance id=%s", d.Id())
		d.SetId("")
		return nil
	}

	var matchedRecord *ssl.UploadUpdateRecordInfo
	for _, record := range records {
		if record.Id != nil && helper.UInt64ToStr(*record.Id) == recordId {
			matchedRecord = record
			break
		}
	}

	if matchedRecord == nil {
		log.Printf("[DescribeHostUploadUpdateRecord] ssl_upload_certificate_instance id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if matchedRecord.Id != nil {
		_ = d.Set("deploy_record_id", *matchedRecord.Id)
	}
	if matchedRecord.OldCertId != nil {
		_ = d.Set("old_certificate_id", matchedRecord.OldCertId)
	}

	return nil
}

func resourceTencentCloudSslUploadCertificateInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_upload_certificate_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	return resourceTencentCloudSslUploadCertificateInstanceRead(d, meta)
}

func resourceTencentCloudSslUploadCertificateInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_upload_certificate_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
