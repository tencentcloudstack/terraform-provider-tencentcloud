/*
Provides a resource to create a ci media_workflow

Example Usage

```hcl
resource "tencentcloud_ci_media_workflow" "media_workflow" {
  name = ""
  workflow_id = ""
  state = ""
  topology {
		dependencies {
			key = ""
			value = ""
		}
		nodes {
			key = ""
			node {
				type = ""
				input {
					queue_id = ""
					object_prefix = ""
					notify_config {
						url = ""
						event = ""
						type = ""
						result_format = ""
					}
					ext_filter {
						state = ""
						audio = ""
						custom = ""
						custom_exts = ""
						all_file = ""
					}
				}
				operation {
					template_id = ""
				}
			}
		}

  }
  create_time = ""
  update_time = ""
  bucket_id = ""
}
```

Import

ci media_workflow can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_workflow.media_workflow media_workflow_id
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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	ci "github.com/tencentyun/cos-go-sdk-v5"
)

func resourceTencentCloudCiMediaWorkflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaWorkflowCreate,
		Read:   resourceTencentCloudCiMediaWorkflowRead,
		Update: resourceTencentCloudCiMediaWorkflowUpdate,
		Delete: resourceTencentCloudCiMediaWorkflowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Codec format, value aac, mp3, flac, amr.",
			},

			"workflow_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: ".",
			},

			"state": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Original audio bit rate, unit: Kbps, Value range: [8, 1000].",
			},

			"topology": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dependencies": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: ".",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: ".",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: ".",
									},
								},
							},
						},
						"nodes": {
							Type:        schema.TypeList,
							Required:    true,
							Description: ".",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: ".",
									},
									"node": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: ".",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: ".",
												},
												"input": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Required:    true,
													Description: ".",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"queue_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: ".",
															},
															"object_prefix": {
																Type:        schema.TypeString,
																Required:    true,
																Description: ".",
															},
															"notify_config": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Required:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"event": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"result_format": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																	},
																},
															},
															"ext_filter": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Required:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"state": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"audio": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"custom": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"custom_exts": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"all_file": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																	},
																},
															},
														},
													},
												},
												"operation": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: ".",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"template_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: ".",
															},
															"output": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"region": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"bucket": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"object": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"au_object": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"sprite_object": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																	},
																},
															},
															"watermark_template_id": {
																Type: schema.TypeSet,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
																Optional:    true,
																Description: ".",
															},
															"delogo_param": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"switch": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"dx": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"dy": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																	},
																},
															},
															"s_d_rto_h_d_r": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"hdr_mode": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																	},
																},
															},
															"s_c_f": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"region": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"function_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"namespace": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																	},
																},
															},
															"hls_pack_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"video_stream_config": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: ".",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"video_stream_name": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: ".",
																					},
																					"band_width": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: ".",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"transcode_template_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: ".",
															},
															"smart_cover": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"format": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"count": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"delete_duplicates": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																	},
																},
															},
															"segment_config": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"format": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																	},
																},
															},
															"digital_watermark": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"message": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"version": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																	},
																},
															},
															"stream_pack_config_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"pack_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																		"ignore_failed_stream": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: ".",
																		},
																		"reserve_all_stream_node": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: ".",
																		},
																	},
																},
															},
															"stream_pack_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"video_stream_config": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: ".",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"video_stream_name": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: ".",
																					},
																					"band_width": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: ".",
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: ".",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: ".",
			},

			"bucket_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: ".",
			},
		},
	}
}

func resourceTencentCloudCiMediaWorkflowCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_workflow.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request       = ci.CreateMediaWorkflowOptions{}
		response      = ci.CreateMediaWorkflowResult{}
		mediaWorkflow = ci.MediaWorkflow{}
		workflowId    string
		bucketId      string
	)
	if v, ok := d.GetOk("name"); ok {
		mediaWorkflow.Name = v.(string)
	}

	if v, ok := d.GetOk("workflow_id"); ok {
		mediaWorkflow.WorkflowId = v.(string)
	}

	if v, ok := d.GetOk("state"); ok {
		mediaWorkflow.State = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "topology"); ok {
		topology := ci.Topology{}
		if dependenciesMap, ok := helper.InterfaceToMap(dMap, "dependencies"); ok {
			dependencies := make(map[string]string)
			if v, ok := dependenciesMap["key"]; ok {
				dependencies["key"] = v.(string)
			}
			if v, ok := dependenciesMap["value"]; ok {
				dependencies["value"] = v.(string)
			}
			topology.Dependencies = dependencies
		}
		if v, ok := dMap["nodes"]; ok {
			for _, item := range v.([]interface{}) {
				nodesMap := item.(map[string]interface{})
				nodes := map[string]ci.Node{}
				key := ""
				if v, ok := nodesMap["key"]; ok {
					key = v.(string)
				}
				if nodeMap, ok := helper.InterfaceToMap(nodesMap, "node"); ok {
					node := ci.Node{}
					if v, ok := nodeMap["type"]; ok {
						node.Type = v.(string)
					}
					if inputMap, ok := helper.InterfaceToMap(nodeMap, "input"); ok {
						nodeInput := ci.NodeInput{}
						if v, ok := inputMap["queue_id"]; ok {
							nodeInput.QueueId = v.(string)
						}
						if v, ok := inputMap["object_prefix"]; ok {
							nodeInput.ObjectPrefix = v.(string)
						}
						if notifyConfigMap, ok := helper.InterfaceToMap(inputMap, "notify_config"); ok {
							notifyConfig := ci.NotifyConfig{}
							if v, ok := notifyConfigMap["url"]; ok {
								notifyConfig.URL = v.(string)
							}
							if v, ok := notifyConfigMap["event"]; ok {
								notifyConfig.Event = v.(string)
							}
							if v, ok := notifyConfigMap["type"]; ok {
								notifyConfig.Type = v.(string)
							}
							if v, ok := notifyConfigMap["result_format"]; ok {
								notifyConfig.ResultFormat = v.(string)
							}
							nodeInput.NotifyConfig = &notifyConfig
						}
						if extFilterMap, ok := helper.InterfaceToMap(inputMap, "ext_filter"); ok {
							extFilter := ci.ExtFilter{}
							if v, ok := extFilterMap["state"]; ok {
								extFilter.State = v.(string)
							}
							if v, ok := extFilterMap["audio"]; ok {
								extFilter.Audio = v.(string)
							}
							if v, ok := extFilterMap["custom"]; ok {
								extFilter.Custom = v.(string)
							}
							if v, ok := extFilterMap["custom_exts"]; ok {
								extFilter.CustomExts = v.(string)
							}
							if v, ok := extFilterMap["all_file"]; ok {
								extFilter.AllFile = v.(string)
							}
							nodeInput.ExtFilter = &extFilter
						}
						node.Input = &nodeInput
					}
					if operationMap, ok := helper.InterfaceToMap(nodeMap, "operation"); ok {
						nodeOperation := ci.NodeOperation{}
						if v, ok := operationMap["template_id"]; ok {
							nodeOperation.TemplateId = v.(string)
						}
						if outputMap, ok := helper.InterfaceToMap(operationMap, "output"); ok {
							nodeOutput := ci.NodeOutput{}
							if v, ok := outputMap["region"]; ok {
								nodeOutput.Region = v.(string)
							}
							if v, ok := outputMap["bucket"]; ok {
								nodeOutput.Bucket = v.(string)
							}
							if v, ok := outputMap["object"]; ok {
								nodeOutput.Object = v.(string)
							}
							if v, ok := outputMap["au_object"]; ok {
								nodeOutput.AuObject = v.(string)
							}
							if v, ok := outputMap["sprite_object"]; ok {
								nodeOutput.SpriteObject = v.(string)
							}
							nodeOperation.Output = &nodeOutput
						}
						if v, ok := operationMap["watermark_template_id"]; ok {
							watermarkTemplateIdSet := v.(*schema.Set).List()
							for i := range watermarkTemplateIdSet {
								watermarkTemplateId := watermarkTemplateIdSet[i].(string)
								nodeOperation.WatermarkTemplateId = append(nodeOperation.WatermarkTemplateId, watermarkTemplateId)
							}
						}
						if delogoParamMap, ok := helper.InterfaceToMap(operationMap, "delogo_param"); ok {
							delogoParam := ci.DelogoParam{}
							if v, ok := delogoParamMap["switch"]; ok {
								delogoParam.Switch = v.(string)
							}
							if v, ok := delogoParamMap["dx"]; ok {
								delogoParam.Dx = v.(string)
							}
							if v, ok := delogoParamMap["dy"]; ok {
								delogoParam.Dy = v.(string)
							}
							if v, ok := delogoParamMap["width"]; ok {
								delogoParam.Width = v.(string)
							}
							if v, ok := delogoParamMap["height"]; ok {
								delogoParam.Height = v.(string)
							}
							nodeOperation.DelogoParam = &delogoParam
						}
						if sDRtoHDRMap, ok := helper.InterfaceToMap(operationMap, "sdr_to_hdr"); ok {
							nodeSDRtoHDR := ci.NodeSDRtoHDR{}
							if v, ok := sDRtoHDRMap["hdr_mode"]; ok {
								nodeSDRtoHDR.HdrMode = v.(string)
							}
							nodeOperation.SDRtoHDR = &nodeSDRtoHDR
						}
						if sCFMap, ok := helper.InterfaceToMap(operationMap, "s_c_f"); ok {
							nodeSCF := ci.NodeSCF{}
							if v, ok := sCFMap["region"]; ok {
								nodeSCF.Region = v.(string)
							}
							if v, ok := sCFMap["function_name"]; ok {
								nodeSCF.FunctionName = v.(string)
							}
							if v, ok := sCFMap["namespace"]; ok {
								nodeSCF.Namespace = v.(string)
							}
							nodeOperation.SCF = &nodeSCF
						}
						if hlsPackInfoMap, ok := helper.InterfaceToMap(operationMap, "hls_pack_info"); ok {
							nodeHlsPackInfo := ci.NodeHlsPackInfo{}
							if v, ok := hlsPackInfoMap["video_stream_config"]; ok {
								for _, item := range v.([]interface{}) {
									videoStreamConfigMap := item.(map[string]interface{})
									videoStreamConfig := ci.VideoStreamConfig{}
									if v, ok := videoStreamConfigMap["video_stream_name"]; ok {
										videoStreamConfig.VideoStreamName = v.(string)
									}
									if v, ok := videoStreamConfigMap["band_width"]; ok {
										videoStreamConfig.BandWidth = v.(string)
									}
									nodeHlsPackInfo.VideoStreamConfig = append(nodeHlsPackInfo.VideoStreamConfig, videoStreamConfig)
								}
							}
							nodeOperation.HlsPackInfo = &nodeHlsPackInfo
						}
						if v, ok := operationMap["transcode_template_id"]; ok {
							nodeOperation.TranscodeTemplateId = v.(string)
						}
						if smartCoverMap, ok := helper.InterfaceToMap(operationMap, "smart_cover"); ok {
							nodeSmartCover := ci.NodeSmartCover{}
							if v, ok := smartCoverMap["format"]; ok {
								nodeSmartCover.Format = v.(string)
							}
							if v, ok := smartCoverMap["width"]; ok {
								nodeSmartCover.Width = v.(string)
							}
							if v, ok := smartCoverMap["height"]; ok {
								nodeSmartCover.Height = v.(string)
							}
							if v, ok := smartCoverMap["count"]; ok {
								nodeSmartCover.Count = v.(string)
							}
							if v, ok := smartCoverMap["delete_duplicates"]; ok {
								nodeSmartCover.DeleteDuplicates = v.(string)
							}
							nodeOperation.SmartCover = &nodeSmartCover
						}
						if segmentConfigMap, ok := helper.InterfaceToMap(operationMap, "segment_config"); ok {
							nodeSegmentConfig := ci.NodeSegmentConfig{}
							if v, ok := segmentConfigMap["format"]; ok {
								nodeSegmentConfig.Format = v.(string)
							}
							if v, ok := segmentConfigMap["duration"]; ok {
								nodeSegmentConfig.Duration = v.(string)
							}
							nodeOperation.SegmentConfig = &nodeSegmentConfig
						}
						if digitalWatermarkMap, ok := helper.InterfaceToMap(operationMap, "digital_watermark"); ok {
							digitalWatermark := ci.DigitalWatermark{}
							if v, ok := digitalWatermarkMap["message"]; ok {
								digitalWatermark.Message = v.(string)
							}
							if v, ok := digitalWatermarkMap["type"]; ok {
								digitalWatermark.Type = v.(string)
							}
							if v, ok := digitalWatermarkMap["version"]; ok {
								digitalWatermark.Version = v.(string)
							}
							nodeOperation.DigitalWatermark = &digitalWatermark
						}
						if streamPackConfigInfoMap, ok := helper.InterfaceToMap(operationMap, "stream_pack_config_info"); ok {
							nodeStreamPackConfigInfo := ci.NodeStreamPackConfigInfo{}
							if v, ok := streamPackConfigInfoMap["pack_type"]; ok {
								nodeStreamPackConfigInfo.PackType = v.(string)
							}
							if v, ok := streamPackConfigInfoMap["ignore_failed_stream"]; ok {
								nodeStreamPackConfigInfo.IgnoreFailedStream = v.(bool)
							}
							if v, ok := streamPackConfigInfoMap["reserve_all_stream_node"]; ok {
								nodeStreamPackConfigInfo.ReserveAllStreamNode = v.(string)
							}
							nodeOperation.StreamPackConfigInfo = &nodeStreamPackConfigInfo
						}
						if streamPackInfoMap, ok := helper.InterfaceToMap(operationMap, "stream_pack_info"); ok {
							nodeHlsPackInfo := ci.NodeHlsPackInfo{}
							if v, ok := streamPackInfoMap["video_stream_config"]; ok {
								for _, item := range v.([]interface{}) {
									videoStreamConfigMap := item.(map[string]interface{})
									videoStreamConfig := ci.VideoStreamConfig{}
									if v, ok := videoStreamConfigMap["video_stream_name"]; ok {
										videoStreamConfig.VideoStreamName = v.(string)
									}
									if v, ok := videoStreamConfigMap["band_width"]; ok {
										videoStreamConfig.BandWidth = v.(string)
									}
									nodeHlsPackInfo.VideoStreamConfig = append(nodeHlsPackInfo.VideoStreamConfig, videoStreamConfig)
								}
							}
							nodeOperation.StreamPackInfo = &nodeHlsPackInfo
						}

						node.Operation = &nodeOperation
					}
					nodes[key] = node
				}
				topology.Nodes = nodes
			}
		}
		mediaWorkflow.Topology = &topology
	}

	if v, ok := d.GetOk("bucket_id"); ok {
		bucketId = v.(string)
		mediaWorkflow.BucketId = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucketId).CI.CreateMediaWorkflow(ctx, &request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaWorkflow", request, result)
		}
		response = *result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaWorkflow failed, reason:%+v", logId, err)
		return err
	}

	workflowId = response.MediaWorkflow.WorkflowId
	d.SetId(bucketId + FILED_SP + workflowId)

	return resourceTencentCloudCiMediaWorkflowRead(d, meta)
}

func resourceTencentCloudCiMediaWorkflowRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_workflow.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	workflowId := idSplit[1]

	mediaWorkflow, err := service.DescribeCiMediaWorkflowById(ctx, bucket, workflowId)
	if err != nil {
		return err
	}

	if mediaWorkflow == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiMediaWorkflow` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mediaWorkflow.Name != "" {
		_ = d.Set("name", mediaWorkflow.Name)
	}

	if mediaWorkflow.WorkflowId != "" {
		_ = d.Set("workflow_id", mediaWorkflow.WorkflowId)
	}

	if mediaWorkflow.State != "" {
		_ = d.Set("state", mediaWorkflow.State)
	}

	if mediaWorkflow.Topology != nil {
		topologyMap := map[string]interface{}{}

		if mediaWorkflow.Topology.Dependencies != nil {
			dependenciesMap := map[string]interface{}{}

			if mediaWorkflow.Topology.Dependencies["key"] != "" {
				dependenciesMap["key"] = mediaWorkflow.Topology.Dependencies["key"]
			}

			if mediaWorkflow.Topology.Dependencies["value"] != "" {
				dependenciesMap["value"] = mediaWorkflow.Topology.Dependencies["value"]
			}

			topologyMap["dependencies"] = []interface{}{dependenciesMap}
		}

		if mediaWorkflow.Topology.Nodes != nil {
			nodesList := []interface{}{}
			for key, nodes := range mediaWorkflow.Topology.Nodes {
				nodesMap := map[string]interface{}{}

				if nodes != (ci.Node{}) {
					nodeMap := map[string]interface{}{}

					if nodes.Type != "" {
						nodeMap["type"] = nodes.Type
					}

					if nodes.Input != nil {
						inputMap := map[string]interface{}{}

						if nodes.Input.QueueId != "" {
							inputMap["queue_id"] = nodes.Input.QueueId
						}

						if nodes.Input.ObjectPrefix != "" {
							inputMap["object_prefix"] = nodes.Input.ObjectPrefix
						}

						if nodes.Input.NotifyConfig != nil {
							notifyConfigMap := map[string]interface{}{}

							if nodes.Input.NotifyConfig.URL != "" {
								notifyConfigMap["url"] = nodes.Input.NotifyConfig.URL
							}

							if nodes.Input.NotifyConfig.Event != "" {
								notifyConfigMap["event"] = nodes.Input.NotifyConfig.Event
							}

							if nodes.Input.NotifyConfig.Type != "" {
								notifyConfigMap["type"] = nodes.Input.NotifyConfig.Type
							}

							if nodes.Input.NotifyConfig.ResultFormat != "" {
								notifyConfigMap["result_format"] = nodes.Input.NotifyConfig.ResultFormat
							}

							inputMap["notify_config"] = []interface{}{notifyConfigMap}
						}

						if nodes.Input.ExtFilter != nil {
							extFilterMap := map[string]interface{}{}

							if nodes.Input.ExtFilter.State != "" {
								extFilterMap["state"] = nodes.Input.ExtFilter.State
							}

							if nodes.Input.ExtFilter.Audio != "" {
								extFilterMap["audio"] = nodes.Input.ExtFilter.Audio
							}

							if nodes.Input.ExtFilter.Custom != "" {
								extFilterMap["custom"] = nodes.Input.ExtFilter.Custom
							}

							if nodes.Input.ExtFilter.CustomExts != "" {
								extFilterMap["custom_exts"] = nodes.Input.ExtFilter.CustomExts
							}

							if nodes.Input.ExtFilter.AllFile != "" {
								extFilterMap["all_file"] = nodes.Input.ExtFilter.AllFile
							}

							inputMap["ext_filter"] = []interface{}{extFilterMap}
						}

						nodeMap["input"] = []interface{}{inputMap}
					}

					if nodes.Operation != nil {
						operationMap := map[string]interface{}{}

						if nodes.Operation.TemplateId != "" {
							operationMap["template_id"] = nodes.Operation.TemplateId
						}

						if nodes.Operation.Output != nil {
							outputMap := map[string]interface{}{}

							if nodes.Operation.Output.Region != "" {
								outputMap["region"] = nodes.Operation.Output.Region
							}

							if nodes.Operation.Output.Bucket != "" {
								outputMap["bucket"] = nodes.Operation.Output.Bucket
							}

							if nodes.Operation.Output.Object != "" {
								outputMap["object"] = nodes.Operation.Output.Object
							}

							if nodes.Operation.Output.AuObject != "" {
								outputMap["au_object"] = nodes.Operation.Output.AuObject
							}

							if nodes.Operation.Output.SpriteObject != "" {
								outputMap["sprite_object"] = nodes.Operation.Output.SpriteObject
							}

							operationMap["output"] = []interface{}{outputMap}
						}

						if nodes.Operation.WatermarkTemplateId != nil {
							operationMap["watermark_template_id"] = nodes.Operation.WatermarkTemplateId
						}

						if nodes.Operation.DelogoParam != nil {
							delogoParamMap := map[string]interface{}{}

							if nodes.Operation.DelogoParam.Switch != "" {
								delogoParamMap["switch"] = nodes.Operation.DelogoParam.Switch
							}

							if nodes.Operation.DelogoParam.Dx != "" {
								delogoParamMap["dx"] = nodes.Operation.DelogoParam.Dx
							}

							if nodes.Operation.DelogoParam.Dy != "" {
								delogoParamMap["dy"] = nodes.Operation.DelogoParam.Dy
							}

							if nodes.Operation.DelogoParam.Width != "" {
								delogoParamMap["width"] = nodes.Operation.DelogoParam.Width
							}

							if nodes.Operation.DelogoParam.Height != "" {
								delogoParamMap["height"] = nodes.Operation.DelogoParam.Height
							}

							operationMap["delogo_param"] = []interface{}{delogoParamMap}
						}

						if nodes.Operation.SDRtoHDR != nil {
							sDRtoHDRMap := map[string]interface{}{}

							if nodes.Operation.SDRtoHDR.HdrMode != "" {
								sDRtoHDRMap["hdr_mode"] = nodes.Operation.SDRtoHDR.HdrMode
							}

							operationMap["s_d_rto_h_d_r"] = []interface{}{sDRtoHDRMap}
						}

						if nodes.Operation.SCF != nil {
							sCFMap := map[string]interface{}{}

							if nodes.Operation.SCF.Region != "" {
								sCFMap["region"] = nodes.Operation.SCF.Region
							}

							if nodes.Operation.SCF.FunctionName != "" {
								sCFMap["function_name"] = nodes.Operation.SCF.FunctionName
							}

							if nodes.Operation.SCF.Namespace != "" {
								sCFMap["namespace"] = nodes.Operation.SCF.Namespace
							}

							operationMap["s_c_f"] = []interface{}{sCFMap}
						}

						if nodes.Operation.HlsPackInfo != nil {
							hlsPackInfoMap := map[string]interface{}{}

							if nodes.Operation.HlsPackInfo.VideoStreamConfig != nil {
								videoStreamConfigList := []interface{}{}
								for _, videoStreamConfig := range nodes.Operation.HlsPackInfo.VideoStreamConfig {
									videoStreamConfigMap := map[string]interface{}{}

									if videoStreamConfig.VideoStreamName != "" {
										videoStreamConfigMap["video_stream_name"] = videoStreamConfig.VideoStreamName
									}

									if videoStreamConfig.BandWidth != "" {
										videoStreamConfigMap["band_width"] = videoStreamConfig.BandWidth
									}

									videoStreamConfigList = append(videoStreamConfigList, videoStreamConfigMap)
								}

								hlsPackInfoMap["video_stream_config"] = []interface{}{videoStreamConfigList}
							}

							operationMap["hls_pack_info"] = []interface{}{hlsPackInfoMap}
						}

						if nodes.Operation.TranscodeTemplateId != "" {
							operationMap["transcode_template_id"] = nodes.Operation.TranscodeTemplateId
						}

						if nodes.Operation.SmartCover != nil {
							smartCoverMap := map[string]interface{}{}

							if nodes.Operation.SmartCover.Format != "" {
								smartCoverMap["format"] = nodes.Operation.SmartCover.Format
							}

							if nodes.Operation.SmartCover.Width != "" {
								smartCoverMap["width"] = nodes.Operation.SmartCover.Width
							}

							if nodes.Operation.SmartCover.Height != "" {
								smartCoverMap["height"] = nodes.Operation.SmartCover.Height
							}

							if nodes.Operation.SmartCover.Count != "" {
								smartCoverMap["count"] = nodes.Operation.SmartCover.Count
							}

							if nodes.Operation.SmartCover.DeleteDuplicates != "" {
								smartCoverMap["delete_duplicates"] = nodes.Operation.SmartCover.DeleteDuplicates
							}

							operationMap["smart_cover"] = []interface{}{smartCoverMap}
						}

						if nodes.Operation.SegmentConfig != nil {
							segmentConfigMap := map[string]interface{}{}

							if nodes.Operation.SegmentConfig.Format != "" {
								segmentConfigMap["format"] = nodes.Operation.SegmentConfig.Format
							}

							if nodes.Operation.SegmentConfig.Duration != "" {
								segmentConfigMap["duration"] = nodes.Operation.SegmentConfig.Duration
							}

							operationMap["segment_config"] = []interface{}{segmentConfigMap}
						}

						if nodes.Operation.DigitalWatermark != nil {
							digitalWatermarkMap := map[string]interface{}{}

							if nodes.Operation.DigitalWatermark.Message != "" {
								digitalWatermarkMap["message"] = nodes.Operation.DigitalWatermark.Message
							}

							if nodes.Operation.DigitalWatermark.Type != "" {
								digitalWatermarkMap["type"] = nodes.Operation.DigitalWatermark.Type
							}

							if nodes.Operation.DigitalWatermark.Version != "" {
								digitalWatermarkMap["version"] = nodes.Operation.DigitalWatermark.Version
							}

							operationMap["digital_watermark"] = []interface{}{digitalWatermarkMap}
						}

						if nodes.Operation.StreamPackConfigInfo != nil {
							streamPackConfigInfoMap := map[string]interface{}{}

							if nodes.Operation.StreamPackConfigInfo.PackType != "" {
								streamPackConfigInfoMap["pack_type"] = nodes.Operation.StreamPackConfigInfo.PackType
							}

							if nodes.Operation.StreamPackConfigInfo.IgnoreFailedStream {
								streamPackConfigInfoMap["ignore_failed_stream"] = nodes.Operation.StreamPackConfigInfo.IgnoreFailedStream
							} else {
								streamPackConfigInfoMap["ignore_failed_stream"] = false
							}

							if nodes.Operation.StreamPackConfigInfo.ReserveAllStreamNode != "" {
								streamPackConfigInfoMap["reserve_all_stream_node"] = nodes.Operation.StreamPackConfigInfo.ReserveAllStreamNode
							}

							operationMap["stream_pack_config_info"] = []interface{}{streamPackConfigInfoMap}
						}

						if nodes.Operation.StreamPackInfo != nil {
							streamPackInfoMap := map[string]interface{}{}

							if nodes.Operation.StreamPackInfo.VideoStreamConfig != nil {
								videoStreamConfigList := []interface{}{}
								for _, videoStreamConfig := range nodes.Operation.StreamPackInfo.VideoStreamConfig {
									videoStreamConfigMap := map[string]interface{}{}

									if videoStreamConfig.VideoStreamName != "" {
										videoStreamConfigMap["video_stream_name"] = videoStreamConfig.VideoStreamName
									}

									if videoStreamConfig.BandWidth != "" {
										videoStreamConfigMap["band_width"] = videoStreamConfig.BandWidth
									}

									videoStreamConfigList = append(videoStreamConfigList, videoStreamConfigMap)
								}

								streamPackInfoMap["video_stream_config"] = []interface{}{videoStreamConfigList}
							}

							operationMap["stream_pack_info"] = []interface{}{streamPackInfoMap}
						}

						nodeMap["operation"] = []interface{}{operationMap}
					}

					nodesMap[key] = []interface{}{nodeMap}
				}

				nodesList = append(nodesList, nodesMap)
			}

			topologyMap["nodes"] = []interface{}{nodesList}
		}

		_ = d.Set("topology", []interface{}{topologyMap})
	}

	if mediaWorkflow.CreateTime != "" {
		_ = d.Set("create_time", mediaWorkflow.CreateTime)
	}

	if mediaWorkflow.UpdateTime != "" {
		_ = d.Set("update_time", mediaWorkflow.UpdateTime)
	}

	if mediaWorkflow.BucketId != "" {
		_ = d.Set("bucket_id", mediaWorkflow.BucketId)
	}

	return nil
}

func resourceTencentCloudCiMediaWorkflowUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_workflow.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := ci.CreateMediaWorkflowOptions{}
	mediaWorkflow := ci.MediaWorkflow{}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	workflowId := idSplit[1]

	immutableArgs := []string{"workflow_id", "create_time", "update_time", "bucket_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			mediaWorkflow.Name = v.(string)
		}
	}

	if d.HasChange("state") {
		if v, ok := d.GetOk("state"); ok {
			mediaWorkflow.State = v.(string)
		}
	}

	if d.HasChange("topology") {
		if dMap, ok := helper.InterfacesHeadMap(d, "topology"); ok {
			topology := ci.Topology{}
			if dependenciesMap, ok := helper.InterfaceToMap(dMap, "dependencies"); ok {
				dependencies := make(map[string]string)
				if v, ok := dependenciesMap["key"]; ok {
					dependencies["key"] = v.(string)
				}
				if v, ok := dependenciesMap["value"]; ok {
					dependencies["value"] = v.(string)
				}
				topology.Dependencies = dependencies
			}
			if v, ok := dMap["nodes"]; ok {
				for _, item := range v.([]interface{}) {
					nodesMap := item.(map[string]interface{})
					nodes := map[string]ci.Node{}
					key := ""
					if v, ok := nodesMap["key"]; ok {
						key = v.(string)
					}
					if nodeMap, ok := helper.InterfaceToMap(nodesMap, "node"); ok {
						node := ci.Node{}
						if v, ok := nodeMap["type"]; ok {
							node.Type = v.(string)
						}
						if inputMap, ok := helper.InterfaceToMap(nodeMap, "input"); ok {
							nodeInput := ci.NodeInput{}
							if v, ok := inputMap["queue_id"]; ok {
								nodeInput.QueueId = v.(string)
							}
							if v, ok := inputMap["object_prefix"]; ok {
								nodeInput.ObjectPrefix = v.(string)
							}
							if notifyConfigMap, ok := helper.InterfaceToMap(inputMap, "notify_config"); ok {
								notifyConfig := ci.NotifyConfig{}
								if v, ok := notifyConfigMap["url"]; ok {
									notifyConfig.URL = v.(string)
								}
								if v, ok := notifyConfigMap["event"]; ok {
									notifyConfig.Event = v.(string)
								}
								if v, ok := notifyConfigMap["type"]; ok {
									notifyConfig.Type = v.(string)
								}
								if v, ok := notifyConfigMap["result_format"]; ok {
									notifyConfig.ResultFormat = v.(string)
								}
								nodeInput.NotifyConfig = &notifyConfig
							}
							if extFilterMap, ok := helper.InterfaceToMap(inputMap, "ext_filter"); ok {
								extFilter := ci.ExtFilter{}
								if v, ok := extFilterMap["state"]; ok {
									extFilter.State = v.(string)
								}
								if v, ok := extFilterMap["audio"]; ok {
									extFilter.Audio = v.(string)
								}
								if v, ok := extFilterMap["custom"]; ok {
									extFilter.Custom = v.(string)
								}
								if v, ok := extFilterMap["custom_exts"]; ok {
									extFilter.CustomExts = v.(string)
								}
								if v, ok := extFilterMap["all_file"]; ok {
									extFilter.AllFile = v.(string)
								}
								nodeInput.ExtFilter = &extFilter
							}
							node.Input = &nodeInput
						}
						if operationMap, ok := helper.InterfaceToMap(nodeMap, "operation"); ok {
							nodeOperation := ci.NodeOperation{}
							if v, ok := operationMap["template_id"]; ok {
								nodeOperation.TemplateId = v.(string)
							}
							node.Operation = &nodeOperation
						}
						nodes[key] = node
					}
					topology.Nodes = nodes
				}
			}
			mediaWorkflow.Topology = &topology
		}
	}
	request.MediaWorkflow = &mediaWorkflow

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.UpdateMediaWorkflow(ctx, &request, workflowId)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaWorkflow", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ci mediaWorkflow failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaWorkflowRead(d, meta)
}

func resourceTencentCloudCiMediaWorkflowDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_workflow.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	workflowId := idSplit[1]

	if err := service.DeleteCiMediaWorkflowById(ctx, bucket, workflowId); err != nil {
		return err
	}

	return nil
}
