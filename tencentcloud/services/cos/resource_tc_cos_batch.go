package cos

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func ResourceTencentCloudCosBatch() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudCosBatchRead,
		Create: resourceTencentCloudCosBatchCreate,
		Update: resourceTencentCloudCosBatchUpdate,
		Delete: resourceTencentCloudCosBatchDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"uin": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Uin.",
			},
			"appid": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Appid.",
			},
			"confirmation_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to confirm before performing the task. The default is false.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Mission description. If you configured this information when you created the task, the content is returned. The description length ranges from 0 to 256 bytes.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Mission priority. The higher the value, the higher the priority of the task. Priority values range from 0 to 2147483647.",
			},
			"role_arn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "COS resource identifier, which is used to identify the role you created. You need this resource identifier to verify your identity.",
			},
			"manifest": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "List of objects to be processed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "The location information of the list of objects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"etag": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the etag of the object list. Length 1-1024 bytes.",
									},
									"object_arn": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the unique resource identifier of the object manifest, which is 1-1024 bytes long.",
									},
									"object_version_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the version of the object manifest ID, which is 1-1024 bytes long.",
									},
								},
							},
						},
						"spec": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Format information that describes the list of objects. If it is a CSV file, this element describes the fields contained in the manifest.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fields": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Describes the fields contained in the listing, which you need to use to specify CSV file fields when Format is COSBatchOperations_CSV_V1. Legal fields are: Ignore, Bucket, Key, VersionId.",
									},
									"format": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the format information for the list of objects. Legal fields are: COSBatchOperations_CSV_V1, COSInventoryReport_CSV_V1.",
									},
								},
							},
						},
					},
				},
			},
			"operation": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Select the action to be performed on the objects in the manifest file.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cos_put_object_copy": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the specific parameters for the batch copy operation on the objects in the list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_control_directive": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "This element specifies how ACL is copied. Valid values:\n" +
											"- Copy: inherits the source object ACL\n" +
											"- Replaced: replace source ACL\n" +
											"- Add: add a new ACL based on the source ACL.",
									},
									"access_control_grants": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Controls the specific access to the object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"display_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "User name.",
												},
												"identifier": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "User ID (UIN) in qcs format. For example: qcs::cam::uin/100000000001:uin/100000000001.",
												},
												"type_identifier": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the type of Identifier. Currently, only user ID is supported. Enumerated value: ID.",
												},
												"permission": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specify a permission to be granted. Enumerated value: READ,WRITE,FULL_CONTROL.",
												},
											},
										},
									},
									"canned_access_control_list": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Defines the ACL property of the object. Valid values: private, public-read.",
									},
									"prefix_replace": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Specifies whether the prefix of the source object needs to be replaced. A value of true indicates the replacement object prefix, which needs to be used with <ResourcesPrefix> and <TargetKeyPrefix>. Default value: false.",
									},
									"resources_prefix": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "This field is valid only when the < PrefixReplace > value is true. Specify the source object prefix to be replaced, and the replacement directory should end with `/`. Can be empty with a maximum length of 1024 bytes.",
									},
									"target_key_prefix": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "This field is valid only when the <PrefixReplace> value is true. This value represents the replaced prefix, and the replacement directory should end with /. Can be empty with a maximum length of 1024 bytes.",
									},
									"modified_since_constraint": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "When the object is modified after the specified time, the operation is performed, otherwise 412 is returned.",
									},
									"unmodified_since_constraint": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "When the object has not been modified after the specified time, the operation is performed, otherwise 412 is returned.",
									},
									"metadata_directive": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "This element specifies whether to copy object metadata from the source object or replace it with metadata in the < NewObjectMetadata > element. Valid values are: Copy, Replaced, Add. Copy: inherit source object metadata; Replaced: replace source metadata; Add: add new metadata based on source metadata.",
									},
									"new_object_metadata": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Configure the metadata for the object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cache_control": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The caching instructions defined in RFC 2616 are saved as object metadata.",
												},
												"content_disposition": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The file name defined in RFC 2616 is saved as object metadata.",
												},
												"content_encoding": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The encoding format defined in RFC 2616 is saved as object metadata.",
												},
												"content_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The content types defined in RFC 2616 are saved as object metadata.",
												},
												"http_expires_date": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The cache expiration time defined in RFC 2616 is saved as object metadata.",
												},
												"sse_algorithm": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Server encryption algorithm. Currently, only AES256 is supported.",
												},
												"user_metadata": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Includes user-defined metadata.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "key.",
															},
															"value": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "value.",
															},
														},
													},
												},
											},
										},
									},
									"tagging_directive": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "This element specifies whether to copy the object tag from the source object or replace it with the tag in the < NewObjectTagging > element. Valid values are: Copy, Replaced, Add. Copy: inherits the source object tag; Replaced: replaces the source tag; Add: adds a new tag based on the source tag.",
									},
									"new_object_tagging": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The label of the configuration object, which must be specified when the < TaggingDirective > value is Replace or Add.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "key.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "value.",
												},
											},
										},
									},
									"storage_class": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Sets the storage level of the object. Enumerated value: STANDARD,STANDARD_IA. Default value: STANDARD.",
									},
									"target_resource": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Sets the target bucket for the Copy. Use qcs to specify, for example, qcs::cos:ap-chengdu:uid/1250000000:examplebucket-1250000000.",
									},
								},
							},
						},
						"cos_initiate_restore_object": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the specific parameters for the batch restore operation for archive storage type objects in the inventory.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"expiration_in_days": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Sets the number of days after which the copy will be automatically expired and deleted, an integer in the range of 1-365.",
									},
									"job_tier": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Select archive recovery model. Available values: Bulk, Standard.",
									},
								},
							},
						},
					},
				},
			},
			"report": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Task completion report.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Delivery bucket for task completion reports.",
						},
						"enabled": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to output the task completion report.",
						},
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task completion report format information. Legal value: Report_CSV_V1.",
						},
						"prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Prefix information for the task completion report. Length 0-256 bytes.",
						},
						"report_scope": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Task completion report the task information that needs to be recorded to determine whether to record the execution information of all operations or the information of failed operations. Legal values: AllTasks, FailedTasksOnly.",
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Current status of the task.\n" +
					"Legal parameter values include Active, Cancelled, Cancelling, Complete, Completing, Failed, Failing, New, Paused, Pausing, Preparing, Ready, Suspended.\n" +
					"For Update status, when you move a task to the Ready state, COS will assume that you have confirmed the task and will perform it. When you move a task to the Cancelled state, COS cancels the task. Optional parameters include: Ready, Cancelled.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job id.",
			},
		},
	}
}

func resourceTencentCloudCosBatchRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_batch.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	uin := idSplit[0]
	appid, _ := strconv.Atoi(idSplit[1])
	jobId := idSplit[2]
	headers := &cos.BatchRequestHeaders{
		XCosAppid: appid,
	}

	result, response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.DescribeJob(ctx, jobId, headers)
	responseBody, _ := json.Marshal(response.Body)
	if err != nil {
		log.Printf("[DEBUG]%s api[DescribeJob] success, request body [%s], response body [%s], err: [%s]\n", logId, jobId, responseBody, err.Error())
		return err
	}
	if result == nil || result.Job == nil {
		return fmt.Errorf("DescribeJob response is nil!")
	}
	confirmationRequired, err := strconv.ParseBool(result.Job.ConfirmationRequired)
	if err != nil {
		return err
	}
	_ = d.Set("uin", uin)
	_ = d.Set("appid", appid)
	_ = d.Set("job_id", jobId)
	_ = d.Set("confirmation_required", confirmationRequired)
	_ = d.Set("description", result.Job.Description)
	_ = d.Set("priority", result.Job.Priority)
	_ = d.Set("role_arn", result.Job.RoleArn)
	manifestResult := make(map[string]interface{})
	locationResult := make(map[string]interface{})
	specResult := make(map[string]interface{})
	manifest := result.Job.Manifest
	location := manifest.Location
	spec := manifest.Spec
	locationResult["etag"] = location.ETag
	locationResult["object_arn"] = location.ObjectArn
	locationResult["object_version_id"] = location.ObjectVersionId
	manifestResult["location"] = []interface{}{locationResult}
	specResult["fields"] = spec.Fields
	specResult["format"] = spec.Format
	manifestResult["spec"] = []interface{}{specResult}
	_ = d.Set("manifest", []interface{}{manifestResult})

	operationResult := make(map[string]interface{})
	if result.Job.Operation.PutObjectCopy != nil {
		putObjectCopyResult := make(map[string]interface{})
		putObjectCopy := result.Job.Operation.PutObjectCopy
		putObjectCopyResult["access_control_directive"] = putObjectCopy.AccessControlDirective
		accessControlGrants := putObjectCopy.AccessControlGrants
		if accessControlGrants != nil {
			accessControlGrantsResult := make(map[string]interface{})
			accessControlGrantsResult["display_name"] = accessControlGrants.COSGrants.Grantee.DisplayName
			accessControlGrantsResult["identifier"] = accessControlGrants.COSGrants.Grantee.Identifier
			accessControlGrantsResult["type_identifier"] = accessControlGrants.COSGrants.Grantee.TypeIdentifier
			accessControlGrantsResult["permission"] = accessControlGrants.COSGrants.Permission
			putObjectCopyResult["access_control_grants"] = []interface{}{accessControlGrantsResult}

		}

		putObjectCopyResult["canned_access_control_list"] = putObjectCopy.CannedAccessControlList
		putObjectCopyResult["prefix_replace"] = putObjectCopy.PrefixReplace
		putObjectCopyResult["resources_prefix"] = putObjectCopy.ResourcesPrefix
		putObjectCopyResult["target_key_prefix"] = putObjectCopy.TargetKeyPrefix
		putObjectCopyResult["modified_since_constraint"] = putObjectCopy.ModifiedSinceConstraint
		putObjectCopyResult["unmodified_since_constraint"] = putObjectCopy.UnModifiedSinceConstraint
		putObjectCopyResult["metadata_directive"] = putObjectCopy.MetadataDirective

		newObjectMetadata := putObjectCopy.NewObjectMetadata
		if newObjectMetadata != nil {
			newObjectMetadataResult := make(map[string]interface{})
			newObjectMetadataResult["cache_control"] = newObjectMetadata.CacheControl
			newObjectMetadataResult["content_disposition"] = newObjectMetadata.ContentDisposition
			newObjectMetadataResult["content_encoding"] = newObjectMetadata.ContentEncoding
			newObjectMetadataResult["content_type"] = newObjectMetadata.ContentType
			newObjectMetadataResult["http_expires_date"] = newObjectMetadata.HttpExpiresDate
			newObjectMetadataResult["sse_algorithm"] = newObjectMetadata.SSEAlgorithm
			userMetadataResult := make([]interface{}, 0)
			userMetadata := newObjectMetadata.UserMetadata
			for _, item := range userMetadata {
				userMetadataResult = append(userMetadataResult, map[string]interface{}{
					"key":   item.Key,
					"value": item.Value,
				})
			}
			newObjectMetadataResult["user_metadata"] = userMetadataResult
			putObjectCopyResult["new_object_metadata"] = []interface{}{newObjectMetadataResult}
		}

		putObjectCopyResult["tagging_directive"] = putObjectCopy.TaggingDirective
		if putObjectCopy.NewObjectTagging != nil {
			cosTagResult := make([]interface{}, 0)
			for _, item := range putObjectCopy.NewObjectTagging.COSTag {
				cosTagResult = append(cosTagResult, map[string]interface{}{
					"key":   item.Key,
					"value": item.Value,
				})
			}
			putObjectCopyResult["new_object_tagging"] = cosTagResult
		}

		putObjectCopyResult["storage_class"] = putObjectCopy.StorageClass
		putObjectCopyResult["target_resource"] = putObjectCopy.TargetResource

		operationResult["cos_put_object_copy"] = []interface{}{putObjectCopyResult}
	}
	if result.Job.Operation.RestoreObject != nil {
		restoreObjectResult := make(map[string]interface{})
		restoreObject := result.Job.Operation.RestoreObject
		restoreObjectResult["expiration_in_days"] = restoreObject.ExpirationInDays
		restoreObjectResult["job_tier"] = restoreObject.JobTier
		operationResult["cos_initiate_restore_object"] = []interface{}{restoreObjectResult}
	}

	_ = d.Set("operation", []interface{}{operationResult})
	if result.Job.Report != nil {
		report := result.Job.Report
		reportResult := make(map[string]interface{})
		reportResult["bucket"] = report.Bucket
		reportResult["enabled"] = report.Enabled
		reportResult["format"] = report.Format
		reportResult["prefix"] = report.Prefix
		reportResult["report_scope"] = report.ReportScope
		_ = d.Set("report", []interface{}{reportResult})
	}

	_ = d.Set("status", result.Job.Status)
	return nil
}

func resourceTencentCloudCosBatchCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_batch.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	uin := d.Get("uin").(string)
	opt := &cos.BatchCreateJobOptions{
		ClientRequestToken: uuid.New().String(),
		Priority:           d.Get("priority").(int),
		RoleArn:            d.Get("role_arn").(string),
	}
	if v, ok := d.GetOk("confirmation_required"); ok && v.(bool) {
		opt.ConfirmationRequired = "true"
	} else {
		opt.ConfirmationRequired = "false"
	}

	if v, ok := d.GetOk("description"); ok {
		opt.Description = v.(string)
	}

	if manifestMap, ok := helper.InterfacesHeadMap(d, "manifest"); ok {
		batchJobManifest := &cos.BatchJobManifest{}
		locationMap := manifestMap["location"].([]interface{})[0].(map[string]interface{})
		specMap := manifestMap["spec"].([]interface{})[0].(map[string]interface{})
		location := &cos.BatchJobManifestLocation{
			ETag:      locationMap["etag"].(string),
			ObjectArn: locationMap["object_arn"].(string),
		}
		if v, ok := locationMap["object_version_id"]; ok {
			location.ObjectVersionId = v.(string)
		}
		batchJobManifest.Location = location
		spec := &cos.BatchJobManifestSpec{
			Format: specMap["format"].(string),
		}
		batchJobManifest.Spec = spec
		if v, ok := specMap["fields"]; ok {
			fields := make([]string, 0)
			for _, item := range v.([]interface{}) {
				fields = append(fields, item.(string))
			}
			spec.Fields = fields
		}
		opt.Manifest = batchJobManifest
	}

	if operationMap, ok := helper.InterfacesHeadMap(d, "operation"); ok {
		operation := &cos.BatchJobOperation{}
		if v, ok := operationMap["cos_put_object_copy"]; ok {
			cosPutObjectCopy := v.([]interface{})[0].(map[string]interface{})
			putObjectCopy := &cos.BatchJobOperationCopy{}

			if v, ok := cosPutObjectCopy["access_control_directive"]; ok {
				putObjectCopy.AccessControlDirective = v.(string)
			}
			if v, ok := cosPutObjectCopy["access_control_grants"]; ok && len(v.([]interface{})) > 0 {
				accessControlGrantMap := v.([]interface{})[0].(map[string]interface{})
				grantee := &cos.BatchGrantee{}
				grant := &cos.BatchCOSGrant{}
				if v, ok := accessControlGrantMap["display_name"]; ok {
					grantee.DisplayName = v.(string)
				}
				if v, ok := accessControlGrantMap["identifier"]; ok {
					grantee.Identifier = v.(string)
				}
				if v, ok := accessControlGrantMap["type_identifier"]; ok {
					grantee.TypeIdentifier = v.(string)
				}
				grant.Grantee = grantee
				if v, ok := accessControlGrantMap["permission"]; ok {
					grant.Permission = v.(string)
				}
				putObjectCopy.AccessControlGrants = &cos.BatchAccessControlGrants{
					COSGrants: grant,
				}
			}
			if v, ok := cosPutObjectCopy["canned_access_control_list"]; ok {
				putObjectCopy.CannedAccessControlList = v.(string)
			}
			if v, ok := cosPutObjectCopy["prefix_replace"]; ok {
				putObjectCopy.PrefixReplace = v.(bool)
			}
			if v, ok := cosPutObjectCopy["resources_prefix"]; ok {
				putObjectCopy.ResourcesPrefix = v.(string)
			}
			if v, ok := cosPutObjectCopy["target_key_prefix"]; ok {
				putObjectCopy.TargetKeyPrefix = v.(string)
			}
			if v, ok := cosPutObjectCopy["metadata_directive"]; ok {
				putObjectCopy.MetadataDirective = v.(string)
			}
			if v, ok := cosPutObjectCopy["modified_since_constraint"]; ok {
				putObjectCopy.ModifiedSinceConstraint = int64(v.(int))
			}
			if v, ok := cosPutObjectCopy["unmodified_since_constraint"]; ok {
				putObjectCopy.UnModifiedSinceConstraint = int64(v.(int))
			}

			if v, ok := cosPutObjectCopy["new_object_metadata"]; ok && len(v.([]interface{})) > 0 {
				newObjectMetadataMap := v.([]interface{})[0].(map[string]interface{})
				newObjectMetadata := &cos.BatchNewObjectMetadata{}
				if v, ok := newObjectMetadataMap["cache_control"]; ok {
					newObjectMetadata.CacheControl = v.(string)
				}
				if v, ok := newObjectMetadataMap["content_disposition"]; ok {
					newObjectMetadata.ContentDisposition = v.(string)
				}
				if v, ok := newObjectMetadataMap["content_encoding"]; ok {
					newObjectMetadata.ContentEncoding = v.(string)
				}
				if v, ok := newObjectMetadataMap["content_type"]; ok {
					newObjectMetadata.ContentType = v.(string)
				}
				if v, ok := newObjectMetadataMap["http_expires_date"]; ok {
					newObjectMetadata.HttpExpiresDate = v.(string)
				}
				if v, ok := newObjectMetadataMap["sse_algorithm"]; ok {
					newObjectMetadata.SSEAlgorithm = v.(string)
				}
				if v, ok := newObjectMetadataMap["user_metadata"]; ok {
					newObjectMetadata.UserMetadata = make([]cos.BatchMetadata, 0)
					for _, userMetadataItem := range v.([]interface{}) {
						userMetadataItemMap := userMetadataItem.(map[string]interface{})
						batchMetadata := cos.BatchMetadata{
							Key:   userMetadataItemMap["key"].(string),
							Value: userMetadataItemMap["value"].(string),
						}
						newObjectMetadata.UserMetadata = append(newObjectMetadata.UserMetadata, batchMetadata)
					}
				}
				putObjectCopy.NewObjectMetadata = newObjectMetadata
			}

			if v, ok := cosPutObjectCopy["tagging_directive"]; ok {
				putObjectCopy.TaggingDirective = v.(string)
			}
			if v, ok := cosPutObjectCopy["new_object_tagging"]; ok {
				newObjectTaggings := v.([]interface{})
				cosTags := make([]cos.BatchCOSTag, 0)
				for _, item := range newObjectTaggings {
					tag := item.(map[string]interface{})
					cosTags = append(cosTags, cos.BatchCOSTag{
						Key:   tag["key"].(string),
						Value: tag["value"].(string),
					})
				}
				putObjectCopy.NewObjectTagging = &cos.BatchNewObjectTagging{COSTag: cosTags}
			}
			if v, ok := cosPutObjectCopy["storage_class"]; ok {
				putObjectCopy.StorageClass = v.(string)
			}
			if v, ok := cosPutObjectCopy["target_resource"]; ok {
				putObjectCopy.TargetResource = v.(string)
			}

			operation.PutObjectCopy = putObjectCopy

		}

		if v, ok := operationMap["cos_initiate_restore_object"]; ok && len(v.([]interface{})) > 0 {
			restoreObject := &cos.BatchInitiateRestoreObject{}
			cosInitiateRestoreObject := v.([]interface{})[0].(map[string]interface{})
			if v, ok := cosInitiateRestoreObject["expiration_in_days"]; ok {
				restoreObject.ExpirationInDays = v.(int)
			}
			if v, ok := cosInitiateRestoreObject["job_tier"]; ok {
				restoreObject.JobTier = v.(string)
			}
			operation.RestoreObject = restoreObject
		}
		opt.Operation = operation
	}

	if reportMap, ok := helper.InterfacesHeadMap(d, "report"); ok {
		batchJobReport := &cos.BatchJobReport{}
		if v, ok := reportMap["bucket"]; ok {
			batchJobReport.Bucket = v.(string)
		}
		if v, ok := reportMap["enabled"]; ok {
			batchJobReport.Enabled = v.(string)
		}
		if v, ok := reportMap["format"]; ok {
			batchJobReport.Format = v.(string)
		}
		if v, ok := reportMap["prefix"]; ok {
			batchJobReport.Prefix = v.(string)
		}
		if v, ok := reportMap["report_scope"]; ok {
			batchJobReport.ReportScope = v.(string)
		}
		opt.Report = batchJobReport
	}
	appid := d.Get("appid").(int)
	headers := &cos.BatchRequestHeaders{
		XCosAppid: appid,
	}
	var batchCreateJobResult *cos.BatchCreateJobResult
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		req, _ := json.Marshal(opt)
		result, response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.CreateJob(ctx, opt, headers)
		responseBody, _ := json.Marshal(response.Body)
		log.Printf("[DEBUG]%s api[CreateJob], request body [%s], response body [%s]\n", logId, req, responseBody)
		if e != nil {
			return tccommon.RetryError(e)
		}

		batchCreateJobResult = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create job failed, reason:%+v", logId, err)
		return err
	}
	if v, ok := d.GetOk("status"); ok {
		opt := &cos.BatchUpdateStatusOptions{
			JobId:              batchCreateJobResult.JobId,
			RequestedJobStatus: v.(string),
		}
		req, _ := json.Marshal(opt)
		_, response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.UpdateJobStatus(ctx, opt, headers)
		responseBody, _ := json.Marshal(response.Body)
		if err != nil {
			log.Printf("[DEBUG]%s api[UpdateJobStatus] error, request body [%s], response body [%s], err: [%s]\n", logId, req, responseBody, err.Error())
			return err
		}
	}
	d.SetId(uin + tccommon.FILED_SP + strconv.Itoa(appid) + tccommon.FILED_SP + batchCreateJobResult.JobId)
	return resourceTencentCloudCosBatchRead(d, meta)
}

func resourceTencentCloudCosBatchUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_batch.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	uin := idSplit[0]
	appid, _ := strconv.Atoi(idSplit[1])
	jobId := idSplit[2]
	headers := &cos.BatchRequestHeaders{
		XCosAppid: appid,
	}
	if d.HasChange("priority") {
		opt := &cos.BatchUpdatePriorityOptions{
			JobId:    jobId,
			Priority: d.Get("priority").(int),
		}
		req, _ := json.Marshal(opt)
		_, response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.UpdateJobPriority(ctx, opt, headers)
		responseBody, _ := json.Marshal(response.Body)
		if err != nil {
			log.Printf("[DEBUG]%s api[UpdateJobPriority] error, request body [%s], response body [%s], err: [%s]\n", logId, req, responseBody, err.Error())
			return err
		}
	}
	if d.HasChange("status") {
		opt := &cos.BatchUpdateStatusOptions{
			JobId:              jobId,
			RequestedJobStatus: d.Get("status").(string),
		}
		req, _ := json.Marshal(opt)
		_, response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.UpdateJobStatus(ctx, opt, headers)
		responseBody, _ := json.Marshal(response.Body)
		if err != nil {
			log.Printf("[DEBUG]%s api[UpdateJobStatus] error, request body [%s], response body [%s], err: [%s]\n", logId, req, responseBody, err.Error())
			return err
		}
	}
	return resourceTencentCloudCosBatchRead(d, meta)
}

func resourceTencentCloudCosBatchDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_batch.delete")()
	defer tccommon.InconsistentCheck(d, meta)()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	uin := idSplit[0]
	appid, _ := strconv.Atoi(idSplit[1])
	jobId := idSplit[2]
	headers := &cos.BatchRequestHeaders{
		XCosAppid: appid,
	}
	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCosBatchClient(uin).Batch.DeleteJob(ctx, jobId, headers)
	responseBody, _ := json.Marshal(response.Body)
	if err != nil {
		log.Printf("[DEBUG]%s api[DeleteJob] success, response body [%s], err: [%s]\n", logId, responseBody, err.Error())
		return err
	}
	return nil
}
