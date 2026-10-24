package ssl

import (
	"fmt"
	"log"

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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
				Required:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
				Description: "Certificate public key. If you upload the certificate public key, CertificateId does not need to be passed.",
			},

			"certificate_private_key": {
				Required:    true,
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

			"deploy_status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Update task creation status. 1 means created successfully; 0 means there is an update task in progress, and no new update task is created.",
			},

			"update_sync_progress": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Update asynchronous creation task progress details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Asynchronous update progress status: 0 pending, 1 processed, 3 processing.",
						},
						"update_sync_progress_regions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Region result list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region.",
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total count.",
									},
									"offset_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of completed executions.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Asynchronous update progress status: 0 pending, 1 processed, 3 processing.",
									},
								},
							},
						},
					},
				},
			},

			"deploy_record_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Certificate deployment record list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Record ID.",
						},
						"old_cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Original certificate ID.",
						},
						"resource_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of deployed resource types.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Deployment status.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deployment time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Last update time.",
						},
					},
				},
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
	)

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
	log.Printf("[CRUD]%s ssl_upload_certificate_instance create success, logId=%s, id=%d", logId, logId, deployRecordId)
	d.SetId(helper.UInt64ToStr(deployRecordId))

	// retain deploy_status / update_sync_progress from create response
	if response.Response.DeployStatus != nil {
		_ = d.Set("deploy_status", *response.Response.DeployStatus)
	}
	if response.Response.UpdateSyncProgress != nil && len(response.Response.UpdateSyncProgress) > 0 {
		updateSyncProgressList := make([]map[string]interface{}, 0, len(response.Response.UpdateSyncProgress))
		for _, progress := range response.Response.UpdateSyncProgress {
			progressMap := map[string]interface{}{}
			if progress.ResourceType != nil {
				progressMap["resource_type"] = progress.ResourceType
			}
			if progress.Status != nil {
				progressMap["status"] = progress.Status
			}
			if progress.UpdateSyncProgressRegions != nil && len(progress.UpdateSyncProgressRegions) > 0 {
				regionsList := make([]map[string]interface{}, 0, len(progress.UpdateSyncProgressRegions))
				for _, region := range progress.UpdateSyncProgressRegions {
					regionMap := map[string]interface{}{}
					if region.Region != nil {
						regionMap["region"] = region.Region
					}
					if region.TotalCount != nil {
						regionMap["total_count"] = region.TotalCount
					}
					if region.OffsetCount != nil {
						regionMap["offset_count"] = region.OffsetCount
					}
					if region.Status != nil {
						regionMap["status"] = region.Status
					}
					regionsList = append(regionsList, regionMap)
				}
				progressMap["update_sync_progress_regions"] = regionsList
			}
			updateSyncProgressList = append(updateSyncProgressList, progressMap)
		}
		_ = d.Set("update_sync_progress", updateSyncProgressList)
	}

	return resourceTencentCloudSslUploadCertificateInstanceRead(d, meta)
}

func resourceTencentCloudSslUploadCertificateInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_upload_certificate_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var oldCertificateId string
	if v, ok := d.GetOk("old_certificate_id"); ok {
		oldCertificateId = v.(string)
	}

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
		log.Printf("[CRUD] ssl_upload_certificate_instance id=%s", d.Id())
		d.SetId("")
		return nil
	}

	id := d.Id()
	var matchedRecord *ssl.UploadUpdateRecordInfo
	for _, record := range records {
		if record.Id != nil && helper.UInt64ToStr(*record.Id) == id {
			matchedRecord = record
			break
		}
	}

	if matchedRecord == nil {
		log.Printf("[CRUD] ssl_upload_certificate_instance id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if matchedRecord.Id != nil {
		_ = d.Set("deploy_record_id", *matchedRecord.Id)
	}
	if matchedRecord.OldCertId != nil {
		_ = d.Set("old_certificate_id", matchedRecord.OldCertId)
	}
	if matchedRecord.Status != nil {
		_ = d.Set("deploy_status", *matchedRecord.Status)
	}

	deployRecordList := make([]map[string]interface{}, 0, len(records))
	for _, record := range records {
		recordMap := map[string]interface{}{}
		if record.Id != nil {
			recordMap["id"] = *record.Id
		}
		if record.OldCertId != nil {
			recordMap["old_cert_id"] = record.OldCertId
		}
		if record.ResourceTypes != nil {
			resourceTypesList := make([]string, 0, len(record.ResourceTypes))
			for _, rt := range record.ResourceTypes {
				if rt != nil {
					resourceTypesList = append(resourceTypesList, *rt)
				}
			}
			recordMap["resource_types"] = resourceTypesList
		}
		if record.Status != nil {
			recordMap["status"] = *record.Status
		}
		if record.CreateTime != nil {
			recordMap["create_time"] = record.CreateTime
		}
		if record.UpdateTime != nil {
			recordMap["update_time"] = record.UpdateTime
		}
		deployRecordList = append(deployRecordList, recordMap)
	}
	_ = d.Set("deploy_record_list", deployRecordList)

	// call DescribeHostUploadUpdateRecordDetail for additional detail
	var detailLimit int64 = 200
	var detailOffset int64 = 0
	var detailRecords []*ssl.DeployRecord
	detailErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		for {
			detailRequest := ssl.NewDescribeHostUploadUpdateRecordDetailRequest()
			detailRequest.Limit = &detailLimit
			detailRequest.Offset = &detailOffset
			detailIdInt64 := helper.StrToInt64(d.Id())
			detailRequest.DeployRecordId = &detailIdInt64

			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient().DescribeHostUploadUpdateRecordDetail(detailRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, detailRequest.GetAction(), detailRequest.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || len(result.Response.DeployRecordDetail) == 0 {
				break
			}
			detailRecords = append(detailRecords, result.Response.DeployRecordDetail...)
			break
		}
		return nil
	})

	if detailErr != nil {
		log.Printf("[CRITAL]%s read ssl uploadCertificateInstance detail failed, reason:%+v", logId, detailErr)
		return detailErr
	}

	if len(detailRecords) > 0 {
		detail := detailRecords[0]
		log.Printf("[DEBUG]%s ssl_upload_certificate_instance detail, total=%d, success=%d, failed=%d, running=%d", logId, detail.TotalCount, detail.SuccessTotalCount, detail.FailedTotalCount, detail.RunningTotalCount)
	}

	return nil
}

func resourceTencentCloudSslUploadCertificateInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_upload_certificate_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = ssl.NewUploadUpdateCertificateInstanceRequest()
		response       = ssl.NewUploadUpdateCertificateInstanceResponse()
		deployRecordId uint64
	)

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
	log.Printf("[CRUD]%s ssl_upload_certificate_instance update success, logId=%s, id=%d", logId, logId, deployRecordId)
	d.SetId(helper.UInt64ToStr(deployRecordId))

	if response.Response.DeployStatus != nil {
		_ = d.Set("deploy_status", *response.Response.DeployStatus)
	}
	if response.Response.UpdateSyncProgress != nil && len(response.Response.UpdateSyncProgress) > 0 {
		updateSyncProgressList := make([]map[string]interface{}, 0, len(response.Response.UpdateSyncProgress))
		for _, progress := range response.Response.UpdateSyncProgress {
			progressMap := map[string]interface{}{}
			if progress.ResourceType != nil {
				progressMap["resource_type"] = progress.ResourceType
			}
			if progress.Status != nil {
				progressMap["status"] = progress.Status
			}
			if progress.UpdateSyncProgressRegions != nil && len(progress.UpdateSyncProgressRegions) > 0 {
				regionsList := make([]map[string]interface{}, 0, len(progress.UpdateSyncProgressRegions))
				for _, region := range progress.UpdateSyncProgressRegions {
					regionMap := map[string]interface{}{}
					if region.Region != nil {
						regionMap["region"] = region.Region
					}
					if region.TotalCount != nil {
						regionMap["total_count"] = region.TotalCount
					}
					if region.OffsetCount != nil {
						regionMap["offset_count"] = region.OffsetCount
					}
					if region.Status != nil {
						regionMap["status"] = region.Status
					}
					regionsList = append(regionsList, regionMap)
				}
				progressMap["update_sync_progress_regions"] = regionsList
			}
			updateSyncProgressList = append(updateSyncProgressList, progressMap)
		}
		_ = d.Set("update_sync_progress", updateSyncProgressList)
	}

	return resourceTencentCloudSslUploadCertificateInstanceRead(d, meta)
}

func resourceTencentCloudSslUploadCertificateInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ssl_upload_certificate_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	deployRecordId := helper.StrToInt64(d.Id())

	request := ssl.NewUploadUpdateCertificateRecordRollbackRequest()
	request.DeployRecordId = helper.Int64(deployRecordId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseSSLCertificateClient().UploadUpdateCertificateRecordRollback(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete ssl uploadCertificateInstance failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
