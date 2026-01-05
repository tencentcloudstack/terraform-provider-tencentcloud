package cos

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

func ResourceTencentCloudCosBucketReplication() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCosBucketReplicationCreate,
		Read:   resourceTencentCloudCosBucketReplicationRead,
		Update: resourceTencentCloudCosBucketReplicationUpdate,
		Delete: resourceTencentCloudCosBucketReplicationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Bucket format should be [custom name]-[appid], for example `mycos-1258798060`.",
			},

			"role": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Initiator identity identifier: qcs::cam::uin/<OwnerUin>:uin/<SubUin>, where <OwnerUin> is the primary account UIN, and <SubUin> can be the primary account UIN or an authorized sub-account UIN.",
			},

			"rule": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specific configuration information; supports up to 1000 entries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Used to specify the name of a particular rule.",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Indicates whether the Rule is enabled or disabled. Possible values: Enabled, Disabled.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Execution priority, used to handle scenarios where the target storage buckets are the same and multiple replication rules match the same object.\nNote: Supports setting positive integers in the range of 1-1000. The Priority values of different rules cannot be duplicated. Storage bucket replication rules must either all have Priority set or all not have Priority set. When all rules have Priority set, overlapping prefixes are allowed for different rules when the target storage buckets are the same. When different rules match the same object, the rule with the smallest Priority value will be triggered first. When none of the rules have Priority set, overlapping prefixes are not allowed for different rules.",
						},
						"filter": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Filter the objects to be copied. The bucket feature will copy objects that match the prefix and tags specified in the Filter settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"and": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "When filtering objects to be copied, if both prefix and object tag conditions are required simultaneously, or if multiple object tag conditions are needed, they must be wrapped in an `and` statement.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"prefix": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The prefix of the objects to be copied.",
												},
												"tag": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "When filtering objects to be copied, you can use object tags (multiple tags are supported) as filtering criteria, with a maximum of 10 tags allowed. After adding tags as filtering criteria, the `synchronize deletion` option must be set to false.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Tag key.",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Tag value.",
															},
														},
													},
												},
											},
										},
									},
									"prefix": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The prefix of the objects to be copied.",
									},
								},
							},
						},
						"destination": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Target storage bucket information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Resource Identifier: qcs::cos:<Region>::<BucketName-APPID>.",
									},
									"storage_class": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Target storage type: This represents the storage type used when storing data in the target storage bucket. Examples include STANDARD, STANDARD_IA, etc.\nNote: The `storage_class` parameter is mandatory if any of the following conditions are met: The source bucket and the target bucket have different availability zone configurations (one is a multi-AZ bucket, the other is a single-AZ bucket). The source bucket and the target bucket have different intelligent tiering settings (one has intelligent tiering enabled, the other does not). If the StorageClass parameter is not specified, the storage class of the objects delivered to the target bucket will default to the same as the source bucket.",
									},
									// "encryption_configuration": {
									// 	Type:        schema.TypeList,
									// 	Optional:    true,
									// 	Description: "This field must be included when `source_selection_criteria.sse_kms_encrypted_objects.status` is set to Enabled. It is used to specify the KMS key used for KMS-encrypted objects copied to the destination bucket.",
									// 	Elem: &schema.Resource{
									// 		Schema: map[string]*schema.Schema{
									// 			"replica_kms_key_id": {
									// 				Type:        schema.TypeString,
									// 				Optional:    true,
									// 				Description: "KMS key ID.",
									// 			},
									// 		},
									// 	},
									// },
								},
							},
						},
						"delete_marker_replication": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Synchronized deletion marker.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to synchronously delete the tag, supports Disabled or Enabled. The default value is Enabled, meaning the tag will be deleted synchronously.",
									},
								},
							},
						},
						// "source_selection_criteria": {
						// 	Type:        schema.TypeList,
						// 	Optional:    true,
						// 	Description: "This is used to specify additional conditions for objects supported by bucket replication rules. Currently, only the option to replicate KMS-encrypted objects is supported.",
						// 	Elem: &schema.Resource{
						// 		Schema: map[string]*schema.Schema{
						// 			"sse_kms_encrypted_objects": {
						// 				Type:        schema.TypeList,
						// 				Optional:    true,
						// 				Description: "Choose whether to copy the KMS-encrypted objects.",
						// 				Elem: &schema.Resource{
						// 					Schema: map[string]*schema.Schema{
						// 						"status": {
						// 							Type:        schema.TypeString,
						// 							Optional:    true,
						// 							Description: "Choose whether to copy KMS encrypted objects; supported values ​​are Enabled and Disabled.",
						// 						},
						// 					},
						// 				},
						// 			},
						// 		},
						// 	},
						// },
					},
				},
			},
		},
	}
}

func resourceTencentCloudCosBucketReplicationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_replication.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var bucket string
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	}

	d.SetId(bucket)

	return resourceTencentCloudCosBucketReplicationUpdate(d, meta)
}

func resourceTencentCloudCosBucketReplicationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_replication.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = CosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		bucket  = d.Id()
	)

	bucketReplication, err := service.DescribeCosBucketReplicationById(ctx, bucket)
	if err != nil {
		return err
	}

	if bucketReplication == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cos_bucket_replication` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("bucket", bucket)

	if bucketReplication.Role != "" {
		_ = d.Set("role", bucketReplication.Role)
	}

	if bucketReplication.Rule != nil {
		ruleList := make([]map[string]interface{}, 0, len(bucketReplication.Rule))
		for _, rule := range bucketReplication.Rule {
			ruleMap := map[string]interface{}{
				"id":       rule.ID,
				"status":   rule.Status,
				"priority": rule.Priority,
			}

			if rule.Filter != nil {
				filterMap := map[string]interface{}{
					"prefix": rule.Filter.Prefix,
				}

				if rule.Filter.And != nil {
					andMap := map[string]interface{}{
						"prefix": rule.Filter.And.Prefix,
					}

					if rule.Filter.And.Tag != nil {
						tagList := make([]map[string]interface{}, 0)
						for _, tag := range rule.Filter.And.Tag {
							tagMap := map[string]interface{}{
								"key":   tag.Key,
								"value": tag.Value,
							}

							tagList = append(tagList, tagMap)
						}

						andMap["tag"] = tagList
					}

					filterMap["and"] = []interface{}{andMap}
				}

				ruleMap["filter"] = []interface{}{filterMap}
			}

			if rule.Destination != nil {
				destinationMap := map[string]interface{}{
					"bucket":        rule.Destination.Bucket,
					"storage_class": rule.Destination.StorageClass,
				}

				ruleMap["destination"] = []interface{}{destinationMap}
			}

			if rule.DeleteMarkerReplication != nil {
				deleteMarkerReplicationMap := map[string]interface{}{
					"status": rule.DeleteMarkerReplication.Status,
				}

				ruleMap["delete_marker_replication"] = []interface{}{deleteMarkerReplicationMap}
			}

			ruleList = append(ruleList, ruleMap)
		}

		_ = d.Set("rule", ruleList)
	}

	return nil
}

func resourceTencentCloudCosBucketReplicationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_replication.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		bucket = d.Id()
	)

	request := cos.PutBucketReplicationOptions{}
	if v, ok := d.GetOk("role"); ok {
		request.Role = v.(string)
	}

	if v, ok := d.GetOk("rule"); ok {
		rules := make([]cos.BucketReplicationRule, 0)
		for _, item := range v.(*schema.Set).List() {
			rule := cos.BucketReplicationRule{}
			dMap := item.(map[string]interface{})
			if v, ok := dMap["id"]; ok && v.(string) != "" {
				rule.ID = v.(string)
			}

			if v, ok := dMap["staus"]; ok && v.(string) != "" {
				rule.Status = v.(string)
			}

			if v, ok := dMap["priority"]; ok {
				rule.Priority = v.(int)
			}

			if v, ok := dMap["filter"]; ok {
				for _, filterItem := range v.(*schema.Set).List() {
					filter := cos.ReplicationFilter{}
					filterMap := filterItem.(map[string]interface{})
					if v, ok := filterMap["and"]; ok {
						for _, andItem := range v.(*schema.Set).List() {
							and := cos.ReplicationFilterAnd{}
							andMap := andItem.(map[string]interface{})
							if v, ok := andMap["prefix"]; ok && v.(string) != "" {
								and.Prefix = v.(string)
							}

							if v, ok := andMap["tag"]; ok {
								for _, tagItem := range v.(*schema.Set).List() {
									tag := cos.ObjectTaggingTag{}
									tagMap := tagItem.(map[string]interface{})
									if v, ok := tagMap["key"]; ok && v.(string) != "" {
										tag.Key = v.(string)
									}

									if v, ok := tagMap["value"]; ok && v.(string) != "" {
										tag.Value = v.(string)
									}

									and.Tag = append(and.Tag, tag)
								}
							}

							filter.And = &and
						}
					}

					if v, ok := filterMap["prefix"]; ok && v.(string) != "" {
						filter.Prefix = v.(string)
					}

					rule.Filter = &filter
				}
			}

			if v, ok := dMap["destination"]; ok {
				for _, destinationItem := range v.(*schema.Set).List() {
					destination := cos.ReplicationDestination{}
					destinationMap := destinationItem.(map[string]interface{})
					if v, ok := destinationMap["bucket"]; ok && v.(string) != "" {
						destination.Bucket = v.(string)
					}

					if v, ok := destinationMap["storage_class"]; ok && v.(string) != "" {
						destination.StorageClass = v.(string)
					}

					rule.Destination = &destination
				}
			}

			if v, ok := d.GetOk("delete_marker_replication"); ok {
				for _, deleteMarkerReplicationItem := range v.(*schema.Set).List() {
					deleteMarkerReplication := cos.DeleteMarkerReplication{}
					deleteMarkerReplicationMap := deleteMarkerReplicationItem.(map[string]interface{})
					if v, ok := deleteMarkerReplicationMap["status"]; ok && v.(string) != "" {
						deleteMarkerReplication.Status = v.(string)
					}

					rule.DeleteMarkerReplication = &deleteMarkerReplication
				}
			}

			rules = append(rules, rule)
		}

		request.Rule = rules
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTencentCosClient(bucket).Bucket.PutBucketReplication(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%+v], response status [%s]\n", logId, "PutBucketReplication", request, result.Status)
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s cos versioning failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCosBucketReplicationRead(d, meta)
}

func resourceTencentCloudCosBucketReplicationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cos_bucket_replication.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
