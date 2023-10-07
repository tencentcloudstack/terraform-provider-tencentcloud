/*
Provides a resource to create a mps edit_media_operation

Example Usage

Operation through COS

```hcl
resource "tencentcloud_cos_bucket" "output" {
	bucket = "tf-bucket-mps-output-${local.app_id}"
  }

data "tencentcloud_cos_bucket_object" "object" {
	bucket = "keep-bucket-${local.app_id}"
	key    = "/mps-test/test.mov"
}

resource "tencentcloud_mps_edit_media_operation" "operation" {
  file_infos {
		input_info {
			type = "COS"
			cos_input_info {
				bucket = data.tencentcloud_cos_bucket_object.object.bucket
				region = "%s"
				object = data.tencentcloud_cos_bucket_object.object.key
			}
		}
		start_time_offset = 60
		end_time_offset = 120
  }
  output_storage {
		type = "COS"
		cos_output_storage {
			bucket = tencentcloud_cos_bucket.output.bucket
			region = "%s"
		}
  }
  output_object_path = "/output"
}
```

*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMpsEditMediaOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsEditMediaOperationCreate,
		Read:   resourceTencentCloudMpsEditMediaOperationRead,
		Delete: resourceTencentCloudMpsEditMediaOperationDelete,
		Schema: map[string]*schema.Schema{
			"file_infos": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Information of input video file.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"input_info": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Video input information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The input type. Valid values: `COS`: A COS bucket address.  `URL`: A URL.  `AWS-S3`: An AWS S3 bucket address. Currently, this type is only supported for transcoding tasks.",
									},
									"cos_input_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The information of the COS object to process. This parameter is valid and required when `Type` is `COS`.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bucket": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The COS bucket of the object to process, such as `TopRankVideo-125xxx88`.",
												},
												"region": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The region of the COS bucket, such as `ap-chongqing`.",
												},
												"object": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The path of the object to process, such as `/movie/201907/WildAnimal.mov`.",
												},
											},
										},
									},
									"url_input_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The URL of the object to process. This parameter is valid and required when `Type` is `URL`.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "URL of a video.",
												},
											},
										},
									},
									"s3_input_info": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "The information of the AWS S3 object processed. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"s3_bucket": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The AWS S3 bucket.",
												},
												"s3_region": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The region of the AWS S3 bucket.",
												},
												"s3_object": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The path of the AWS S3 object.",
												},
												"s3_secret_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The key ID required to access the AWS S3 object.",
												},
												"s3_secret_key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The key required to access the AWS S3 object.",
												},
											},
										},
									},
								},
							},
						},
						"start_time_offset": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "Start time offset of video clipping in seconds.",
						},
						"end_time_offset": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "End time offset of video clipping in seconds.",
						},
					},
				},
			},

			"output_storage": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The storage location of the media processing output file.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The storage type for a media processing output file. Valid values: `COS`: Tencent Cloud COS. `AWS-S3`: AWS S3. This type is only supported for AWS tasks, and the output bucket must be in the same region as the bucket of the source file.",
						},
						"cos_output_storage": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The location to save the output object in COS. This parameter is valid and required when `Type` is COS.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The bucket to which the output file of media processing is saved, such as `TopRankVideo-125xxx88`. If this parameter is left empty, the value of the upper layer will be inherited.",
									},
									"region": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The region of the output bucket, such as `ap-chongqing`. If this parameter is left empty, the value of the upper layer will be inherited.",
									},
								},
							},
						},
						"s3_output_storage": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The AWS S3 bucket to save the output file. This parameter is required if `Type` is `AWS-S3`.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"s3_bucket": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The AWS S3 bucket.",
									},
									"s3_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The region of the AWS S3 bucket.",
									},
									"s3_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The key ID required to upload files to the AWS S3 object.",
									},
									"s3_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The key required to upload files to the AWS S3 object.",
									},
								},
							},
						},
					},
				},
			},

			"output_object_path": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The path to save the media processing output file.",
			},

			"output_config": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Configuration for output files of video editing.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Format. Valid values: `mp4` (default), `hls`, `mov`, `flv`, `avi`.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The editing mode. Valid values are `normal` and `fast`. The default is `normal`, which indicates precise editing.",
						},
					},
				},
			},

			"task_notify_config": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Event notification information of task. If this parameter is left empty, no event notifications will be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cmq_model": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CMQ or TDMQ-CMQ model. Valid values: Queue, Topic.",
						},
						"cmq_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CMQ or TDMQ-CMQ region, such as `sh` (Shanghai) or `bj` (Beijing).",
						},
						"topic_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CMQ or TDMQ-CMQ topic to receive notifications. This parameter is valid when `CmqModel` is `Topic`.",
						},
						"queue_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CMQ or TDMQ-CMQ queue to receive notifications. This parameter is valid when `CmqModel` is `Queue`.",
						},
						"notify_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Workflow notification method. Valid values: Finish, Change. If this parameter is left empty, `Finish` will be used.",
						},
						"notify_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The notification type. Valid values: `CMQ`: This value is no longer used. Please use `TDMQ-CMQ` instead. `TDMQ-CMQ`: Message queue. `URL`: If `NotifyType` is set to `URL`, HTTP callbacks are sent to the URL specified by `NotifyUrl`. HTTP and JSON are used for the callbacks. The packet contains the response parameters of the `ParseNotification` API. `SCF`: This notification type is not recommended. You need to configure it in the SCF console. `AWS-SQS`: AWS queue. This type is only supported for AWS tasks, and the queue must be in the same region as the AWS bucket. If you do not pass this parameter or pass in an empty string, `CMQ` will be used. To use a different notification type, specify this parameter accordingly.",
						},
						"notify_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "HTTP callback URL, required if `NotifyType` is set to `URL`.",
						},
						"aws_sqs": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The AWS SQS queue. This parameter is required if `NotifyType` is `AWS-SQS`.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sqs_region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The region of the SQS queue.",
									},
									"sqs_queue_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the SQS queue.",
									},
									"s3_secret_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The key ID required to read from/write to the SQS queue.",
									},
									"s3_secret_key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The key required to read from/write to the SQS queue.",
									},
								},
							},
						},
					},
				},
			},

			"tasks_priority": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Task priority. The higher the value, the higher the priority. Value range: [-10,10]. If this parameter is left empty, 0 will be used.",
			},

			"session_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID used for deduplication. If there was a request with the same ID in the last three days, the current request will return an error. The ID can contain up to 50 characters. If this parameter is left empty or an empty string is entered, no deduplication will be performed.",
			},

			"session_context": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The source context which is used to pass through the user request information. The task flow status change callback will return the value of this field. It can contain up to 1,000 characters.",
			},
		},
	}
}

func resourceTencentCloudMpsEditMediaOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_edit_media_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mps.NewEditMediaRequest()
		response = mps.NewEditMediaResponse()
		taskId   string
	)
	if v, ok := d.GetOk("file_infos"); ok {
		for _, item := range v.([]interface{}) {
			editMediaFileInfo := mps.EditMediaFileInfo{}
			dMap := item.(map[string]interface{})
			if inputInfoMap, ok := helper.InterfaceToMap(dMap, "input_info"); ok {
				mediaInputInfo := mps.MediaInputInfo{}
				if v, ok := inputInfoMap["type"]; ok {
					mediaInputInfo.Type = helper.String(v.(string))
				}
				if cosInputInfoMap, ok := helper.InterfaceToMap(inputInfoMap, "cos_input_info"); ok {
					cosInputInfo := mps.CosInputInfo{}
					if v, ok := cosInputInfoMap["bucket"]; ok {
						cosInputInfo.Bucket = helper.String(v.(string))
					}
					if v, ok := cosInputInfoMap["region"]; ok {
						cosInputInfo.Region = helper.String(v.(string))
					}
					if v, ok := cosInputInfoMap["object"]; ok {
						cosInputInfo.Object = helper.String(v.(string))
					}
					mediaInputInfo.CosInputInfo = &cosInputInfo
				}
				if urlInputInfoMap, ok := helper.InterfaceToMap(inputInfoMap, "url_input_info"); ok {
					urlInputInfo := mps.UrlInputInfo{}
					if v, ok := urlInputInfoMap["url"]; ok {
						urlInputInfo.Url = helper.String(v.(string))
					}
					mediaInputInfo.UrlInputInfo = &urlInputInfo
				}
				if s3InputInfoMap, ok := helper.InterfaceToMap(inputInfoMap, "s3_input_info"); ok {
					s3InputInfo := mps.S3InputInfo{}
					if v, ok := s3InputInfoMap["s3_bucket"]; ok {
						s3InputInfo.S3Bucket = helper.String(v.(string))
					}
					if v, ok := s3InputInfoMap["s3_region"]; ok {
						s3InputInfo.S3Region = helper.String(v.(string))
					}
					if v, ok := s3InputInfoMap["s3_object"]; ok {
						s3InputInfo.S3Object = helper.String(v.(string))
					}
					if v, ok := s3InputInfoMap["s3_secret_id"]; ok {
						s3InputInfo.S3SecretId = helper.String(v.(string))
					}
					if v, ok := s3InputInfoMap["s3_secret_key"]; ok {
						s3InputInfo.S3SecretKey = helper.String(v.(string))
					}
					mediaInputInfo.S3InputInfo = &s3InputInfo
				}
				editMediaFileInfo.InputInfo = &mediaInputInfo
			}
			if v, ok := dMap["start_time_offset"]; ok {
				editMediaFileInfo.StartTimeOffset = helper.Float64(v.(float64))
			}
			if v, ok := dMap["end_time_offset"]; ok {
				editMediaFileInfo.EndTimeOffset = helper.Float64(v.(float64))
			}
			request.FileInfos = append(request.FileInfos, &editMediaFileInfo)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "output_storage"); ok {
		taskOutputStorage := mps.TaskOutputStorage{}
		if v, ok := dMap["type"]; ok {
			taskOutputStorage.Type = helper.String(v.(string))
		}
		if cosOutputStorageMap, ok := helper.InterfaceToMap(dMap, "cos_output_storage"); ok {
			cosOutputStorage := mps.CosOutputStorage{}
			if v, ok := cosOutputStorageMap["bucket"]; ok {
				cosOutputStorage.Bucket = helper.String(v.(string))
			}
			if v, ok := cosOutputStorageMap["region"]; ok {
				cosOutputStorage.Region = helper.String(v.(string))
			}
			taskOutputStorage.CosOutputStorage = &cosOutputStorage
		}
		if s3OutputStorageMap, ok := helper.InterfaceToMap(dMap, "s3_output_storage"); ok {
			s3OutputStorage := mps.S3OutputStorage{}
			if v, ok := s3OutputStorageMap["s3_bucket"]; ok {
				s3OutputStorage.S3Bucket = helper.String(v.(string))
			}
			if v, ok := s3OutputStorageMap["s3_region"]; ok {
				s3OutputStorage.S3Region = helper.String(v.(string))
			}
			if v, ok := s3OutputStorageMap["s3_secret_id"]; ok {
				s3OutputStorage.S3SecretId = helper.String(v.(string))
			}
			if v, ok := s3OutputStorageMap["s3_secret_key"]; ok {
				s3OutputStorage.S3SecretKey = helper.String(v.(string))
			}
			taskOutputStorage.S3OutputStorage = &s3OutputStorage
		}
		request.OutputStorage = &taskOutputStorage
	}

	if v, ok := d.GetOk("output_object_path"); ok {
		request.OutputObjectPath = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "output_config"); ok {
		editMediaOutputConfig := mps.EditMediaOutputConfig{}
		if v, ok := dMap["container"]; ok {
			editMediaOutputConfig.Container = helper.String(v.(string))
		}
		if v, ok := dMap["type"]; ok {
			editMediaOutputConfig.Type = helper.String(v.(string))
		}
		request.OutputConfig = &editMediaOutputConfig
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "task_notify_config"); ok {
		taskNotifyConfig := mps.TaskNotifyConfig{}
		if v, ok := dMap["cmq_model"]; ok {
			taskNotifyConfig.CmqModel = helper.String(v.(string))
		}
		if v, ok := dMap["cmq_region"]; ok {
			taskNotifyConfig.CmqRegion = helper.String(v.(string))
		}
		if v, ok := dMap["topic_name"]; ok {
			taskNotifyConfig.TopicName = helper.String(v.(string))
		}
		if v, ok := dMap["queue_name"]; ok {
			taskNotifyConfig.QueueName = helper.String(v.(string))
		}
		if v, ok := dMap["notify_mode"]; ok {
			taskNotifyConfig.NotifyMode = helper.String(v.(string))
		}
		if v, ok := dMap["notify_type"]; ok {
			taskNotifyConfig.NotifyType = helper.String(v.(string))
		}
		if v, ok := dMap["notify_url"]; ok {
			taskNotifyConfig.NotifyUrl = helper.String(v.(string))
		}
		if awsSQSMap, ok := helper.InterfaceToMap(dMap, "aws_sqs"); ok {
			awsSQS := mps.AwsSQS{}
			if v, ok := awsSQSMap["sqs_region"]; ok {
				awsSQS.SQSRegion = helper.String(v.(string))
			}
			if v, ok := awsSQSMap["sqs_queue_name"]; ok {
				awsSQS.SQSQueueName = helper.String(v.(string))
			}
			if v, ok := awsSQSMap["s3_secret_id"]; ok {
				awsSQS.S3SecretId = helper.String(v.(string))
			}
			if v, ok := awsSQSMap["s3_secret_key"]; ok {
				awsSQS.S3SecretKey = helper.String(v.(string))
			}
			taskNotifyConfig.AwsSQS = &awsSQS
		}
		request.TaskNotifyConfig = &taskNotifyConfig
	}

	if v, _ := d.GetOk("tasks_priority"); v != nil {
		request.TasksPriority = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("session_id"); ok {
		request.SessionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("session_context"); ok {
		request.SessionContext = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().EditMedia(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps editMediaOperation failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudMpsEditMediaOperationRead(d, meta)
}

func resourceTencentCloudMpsEditMediaOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_edit_media_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsEditMediaOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_edit_media_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
