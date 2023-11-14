/*
Provides a resource to create a ci media_workflow

Example Usage

```hcl
resource "tencentcloud_ci_media_workflow" "media_workflow" {
  name = &lt;nil&gt;
  workflow_id = &lt;nil&gt;
  state = &lt;nil&gt;
  topology {
		dependencies {
			key = &lt;nil&gt;
			value = &lt;nil&gt;
		}
		nodes {
			key = &lt;nil&gt;
			node {
				type = &lt;nil&gt;
				input {
					queue_id = &lt;nil&gt;
					object_prefix = &lt;nil&gt;
					notify_config {
						u_r_l = &lt;nil&gt;
						event = &lt;nil&gt;
						type = &lt;nil&gt;
						result_format = &lt;nil&gt;
					}
					ext_filter {
						state = &lt;nil&gt;
						audio = &lt;nil&gt;
						custom = &lt;nil&gt;
						custom_exts = &lt;nil&gt;
						all_file = &lt;nil&gt;
					}
				}
				operation {
					template_id = &lt;nil&gt;
					output {
						region = &lt;nil&gt;
						bucket = &lt;nil&gt;
						object = &lt;nil&gt;
						au_object = &lt;nil&gt;
						sprite_object = &lt;nil&gt;
					}
					watermark_template_id = &lt;nil&gt;
					delogo_param {
						switch = &lt;nil&gt;
						dx = &lt;nil&gt;
						dy = &lt;nil&gt;
						width = &lt;nil&gt;
						height = &lt;nil&gt;
					}
					s_d_rto_h_d_r {
						hdr_mode = &lt;nil&gt;
					}
					s_c_f {
						region = &lt;nil&gt;
						function_name = &lt;nil&gt;
						namespace = &lt;nil&gt;
					}
					hls_pack_info {
						video_stream_config {
							video_stream_name = &lt;nil&gt;
							band_width = &lt;nil&gt;
						}
					}
					transcode_template_id = &lt;nil&gt;
					smart_cover {
						format = &lt;nil&gt;
						width = &lt;nil&gt;
						height = &lt;nil&gt;
						count = &lt;nil&gt;
						delete_duplicates = &lt;nil&gt;
					}
					segment_config {
						format = &lt;nil&gt;
						duration = &lt;nil&gt;
					}
					digital_watermark {
						message = &lt;nil&gt;
						type = &lt;nil&gt;
						version = &lt;nil&gt;
					}
					stream_pack_config_info {
						pack_type = &lt;nil&gt;
						ignore_failed_stream = &lt;nil&gt;
						reserve_all_stream_node = &lt;nil&gt;
					}
					stream_pack_info {
						video_stream_config {
							video_stream_name = &lt;nil&gt;
							band_width = &lt;nil&gt;
						}
					}
				}
			}
		}

  }
      bucket_id = &lt;nil&gt;
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	ci "github.com/tencentyun/cos-go-sdk-v5"
	"log"
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
				Optional:    true,
				Type:        schema.TypeString,
				Description: ".",
			},

			"state": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Original audio bit rate, unit: Kbps, Value range: [8, 1000].",
			},

			"topology": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dependencies": {
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
													Optional:    true,
													Description: ".",
												},
												"input": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
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
																Optional:    true,
																Description: ".",
															},
															"notify_config": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"u_r_l": {
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
																Optional:    true,
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
				Optional:    true,
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

	var (
		request    = ci.NewCreateMediaWorkflowRequest()
		response   = ci.NewCreateMediaWorkflowResponse()
		workflowId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("workflow_id"); ok {
		workflowId = v.(string)
		request.WorkflowId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("state"); ok {
		request.State = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "topology"); ok {
		topology := ci.Topology{}
		if v, ok := dMap["dependencies"]; ok {
			for _, item := range v.([]interface{}) {
				dependenciesMap := item.(map[string]interface{})
				dependencies := ci.Dependencies{}
				if v, ok := dependenciesMap["key"]; ok {
					dependencies.Key = helper.String(v.(string))
				}
				if v, ok := dependenciesMap["value"]; ok {
					dependencies.Value = helper.String(v.(string))
				}
				topology.Dependencies = append(topology.Dependencies, &dependencies)
			}
		}
		if v, ok := dMap["nodes"]; ok {
			for _, item := range v.([]interface{}) {
				nodesMap := item.(map[string]interface{})
				nodes := ci.Nodes{}
				if v, ok := nodesMap["key"]; ok {
					nodes.Key = helper.String(v.(string))
				}
				if nodeMap, ok := helper.InterfaceToMap(nodesMap, "node"); ok {
					nodes := ci.Nodes{}
					if v, ok := nodeMap["type"]; ok {
						nodes.Type = helper.String(v.(string))
					}
					if inputMap, ok := helper.InterfaceToMap(nodeMap, "input"); ok {
						nodeInput := ci.NodeInput{}
						if v, ok := inputMap["queue_id"]; ok {
							nodeInput.QueueId = helper.String(v.(string))
						}
						if v, ok := inputMap["object_prefix"]; ok {
							nodeInput.ObjectPrefix = helper.String(v.(string))
						}
						if notifyConfigMap, ok := helper.InterfaceToMap(inputMap, "notify_config"); ok {
							notifyConfig := ci.NotifyConfig{}
							if v, ok := notifyConfigMap["u_r_l"]; ok {
								notifyConfig.URL = helper.String(v.(string))
							}
							if v, ok := notifyConfigMap["event"]; ok {
								notifyConfig.Event = helper.String(v.(string))
							}
							if v, ok := notifyConfigMap["type"]; ok {
								notifyConfig.Type = helper.String(v.(string))
							}
							if v, ok := notifyConfigMap["result_format"]; ok {
								notifyConfig.ResultFormat = helper.String(v.(string))
							}
							nodeInput.NotifyConfig = &notifyConfig
						}
						if extFilterMap, ok := helper.InterfaceToMap(inputMap, "ext_filter"); ok {
							extFilter := ci.ExtFilter{}
							if v, ok := extFilterMap["state"]; ok {
								extFilter.State = helper.String(v.(string))
							}
							if v, ok := extFilterMap["audio"]; ok {
								extFilter.Audio = helper.String(v.(string))
							}
							if v, ok := extFilterMap["custom"]; ok {
								extFilter.Custom = helper.String(v.(string))
							}
							if v, ok := extFilterMap["custom_exts"]; ok {
								extFilter.CustomExts = helper.String(v.(string))
							}
							if v, ok := extFilterMap["all_file"]; ok {
								extFilter.AllFile = helper.String(v.(string))
							}
							nodeInput.ExtFilter = &extFilter
						}
						nodes.Input = &nodeInput
					}
					if operationMap, ok := helper.InterfaceToMap(nodeMap, "operation"); ok {
						nodeOperation := ci.NodeOperation{}
						if v, ok := operationMap["template_id"]; ok {
							nodeOperation.TemplateId = helper.String(v.(string))
						}
						if outputMap, ok := helper.InterfaceToMap(operationMap, "output"); ok {
							nodeOutput := ci.NodeOutput{}
							if v, ok := outputMap["region"]; ok {
								nodeOutput.Region = helper.String(v.(string))
							}
							if v, ok := outputMap["bucket"]; ok {
								nodeOutput.Bucket = helper.String(v.(string))
							}
							if v, ok := outputMap["object"]; ok {
								nodeOutput.Object = helper.String(v.(string))
							}
							if v, ok := outputMap["au_object"]; ok {
								nodeOutput.AuObject = helper.String(v.(string))
							}
							if v, ok := outputMap["sprite_object"]; ok {
								nodeOutput.SpriteObject = helper.String(v.(string))
							}
							nodeOperation.Output = &nodeOutput
						}
						if v, ok := operationMap["watermark_template_id"]; ok {
							watermarkTemplateIdSet := v.(*schema.Set).List()
							for i := range watermarkTemplateIdSet {
								watermarkTemplateId := watermarkTemplateIdSet[i].(string)
								nodeOperation.WatermarkTemplateId = append(nodeOperation.WatermarkTemplateId, &watermarkTemplateId)
							}
						}
						if delogoParamMap, ok := helper.InterfaceToMap(operationMap, "delogo_param"); ok {
							delogoParam := ci.DelogoParam{}
							if v, ok := delogoParamMap["switch"]; ok {
								delogoParam.Switch = helper.String(v.(string))
							}
							if v, ok := delogoParamMap["dx"]; ok {
								delogoParam.Dx = helper.String(v.(string))
							}
							if v, ok := delogoParamMap["dy"]; ok {
								delogoParam.Dy = helper.String(v.(string))
							}
							if v, ok := delogoParamMap["width"]; ok {
								delogoParam.Width = helper.String(v.(string))
							}
							if v, ok := delogoParamMap["height"]; ok {
								delogoParam.Height = helper.String(v.(string))
							}
							nodeOperation.DelogoParam = &delogoParam
						}
						if sDRtoHDRMap, ok := helper.InterfaceToMap(operationMap, "s_d_rto_h_d_r"); ok {
							nodeSDRtoHDR := ci.NodeSDRtoHDR{}
							if v, ok := sDRtoHDRMap["hdr_mode"]; ok {
								nodeSDRtoHDR.HdrMode = helper.String(v.(string))
							}
							nodeOperation.SDRtoHDR = &nodeSDRtoHDR
						}
						if sCFMap, ok := helper.InterfaceToMap(operationMap, "s_c_f"); ok {
							nodeSCF := ci.NodeSCF{}
							if v, ok := sCFMap["region"]; ok {
								nodeSCF.Region = helper.String(v.(string))
							}
							if v, ok := sCFMap["function_name"]; ok {
								nodeSCF.FunctionName = helper.String(v.(string))
							}
							if v, ok := sCFMap["namespace"]; ok {
								nodeSCF.Namespace = helper.String(v.(string))
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
										videoStreamConfig.VideoStreamName = helper.String(v.(string))
									}
									if v, ok := videoStreamConfigMap["band_width"]; ok {
										videoStreamConfig.BandWidth = helper.String(v.(string))
									}
									nodeHlsPackInfo.VideoStreamConfig = append(nodeHlsPackInfo.VideoStreamConfig, &videoStreamConfig)
								}
							}
							nodeOperation.HlsPackInfo = &nodeHlsPackInfo
						}
						if v, ok := operationMap["transcode_template_id"]; ok {
							nodeOperation.TranscodeTemplateId = helper.String(v.(string))
						}
						if smartCoverMap, ok := helper.InterfaceToMap(operationMap, "smart_cover"); ok {
							nodeSmartCover := ci.NodeSmartCover{}
							if v, ok := smartCoverMap["format"]; ok {
								nodeSmartCover.Format = helper.String(v.(string))
							}
							if v, ok := smartCoverMap["width"]; ok {
								nodeSmartCover.Width = helper.String(v.(string))
							}
							if v, ok := smartCoverMap["height"]; ok {
								nodeSmartCover.Height = helper.String(v.(string))
							}
							if v, ok := smartCoverMap["count"]; ok {
								nodeSmartCover.Count = helper.String(v.(string))
							}
							if v, ok := smartCoverMap["delete_duplicates"]; ok {
								nodeSmartCover.DeleteDuplicates = helper.String(v.(string))
							}
							nodeOperation.SmartCover = &nodeSmartCover
						}
						if segmentConfigMap, ok := helper.InterfaceToMap(operationMap, "segment_config"); ok {
							nodeSegmentConfig := ci.NodeSegmentConfig{}
							if v, ok := segmentConfigMap["format"]; ok {
								nodeSegmentConfig.Format = helper.String(v.(string))
							}
							if v, ok := segmentConfigMap["duration"]; ok {
								nodeSegmentConfig.Duration = helper.String(v.(string))
							}
							nodeOperation.SegmentConfig = &nodeSegmentConfig
						}
						if digitalWatermarkMap, ok := helper.InterfaceToMap(operationMap, "digital_watermark"); ok {
							digitalWatermark := ci.DigitalWatermark{}
							if v, ok := digitalWatermarkMap["message"]; ok {
								digitalWatermark.Message = helper.String(v.(string))
							}
							if v, ok := digitalWatermarkMap["type"]; ok {
								digitalWatermark.Type = helper.String(v.(string))
							}
							if v, ok := digitalWatermarkMap["version"]; ok {
								digitalWatermark.Version = helper.String(v.(string))
							}
							nodeOperation.DigitalWatermark = &digitalWatermark
						}
						if streamPackConfigInfoMap, ok := helper.InterfaceToMap(operationMap, "stream_pack_config_info"); ok {
							nodeStreamPackConfigInfo := ci.NodeStreamPackConfigInfo{}
							if v, ok := streamPackConfigInfoMap["pack_type"]; ok {
								nodeStreamPackConfigInfo.PackType = helper.String(v.(string))
							}
							if v, ok := streamPackConfigInfoMap["ignore_failed_stream"]; ok {
								nodeStreamPackConfigInfo.IgnoreFailedStream = helper.Bool(v.(bool))
							}
							if v, ok := streamPackConfigInfoMap["reserve_all_stream_node"]; ok {
								nodeStreamPackConfigInfo.ReserveAllStreamNode = helper.String(v.(string))
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
										videoStreamConfig.VideoStreamName = helper.String(v.(string))
									}
									if v, ok := videoStreamConfigMap["band_width"]; ok {
										videoStreamConfig.BandWidth = helper.String(v.(string))
									}
									nodeHlsPackInfo.VideoStreamConfig = append(nodeHlsPackInfo.VideoStreamConfig, &videoStreamConfig)
								}
							}
							nodeOperation.StreamPackInfo = &nodeHlsPackInfo
						}
						nodes.Operation = &nodeOperation
					}
					nodes.Node = &nodes
				}
				topology.Nodes = append(topology.Nodes, &nodes)
			}
		}
		request.Topology = &topology
	}

	if v, ok := d.GetOk("bucket_id"); ok {
		request.BucketId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().CreateMediaWorkflow(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaWorkflow failed, reason:%+v", logId, err)
		return err
	}

	workflowId = *response.Response.WorkflowId
	d.SetId(workflowId)

	return resourceTencentCloudCiMediaWorkflowRead(d, meta)
}

func resourceTencentCloudCiMediaWorkflowRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_workflow.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	mediaWorkflowId := d.Id()

	mediaWorkflow, err := service.DescribeCiMediaWorkflowById(ctx, workflowId)
	if err != nil {
		return err
	}

	if mediaWorkflow == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiMediaWorkflow` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mediaWorkflow.Name != nil {
		_ = d.Set("name", mediaWorkflow.Name)
	}

	if mediaWorkflow.WorkflowId != nil {
		_ = d.Set("workflow_id", mediaWorkflow.WorkflowId)
	}

	if mediaWorkflow.State != nil {
		_ = d.Set("state", mediaWorkflow.State)
	}

	if mediaWorkflow.Topology != nil {
		topologyMap := map[string]interface{}{}

		if mediaWorkflow.Topology.Dependencies != nil {
			dependenciesList := []interface{}{}
			for _, dependencies := range mediaWorkflow.Topology.Dependencies {
				dependenciesMap := map[string]interface{}{}

				if dependencies.Key != nil {
					dependenciesMap["key"] = dependencies.Key
				}

				if dependencies.Value != nil {
					dependenciesMap["value"] = dependencies.Value
				}

				dependenciesList = append(dependenciesList, dependenciesMap)
			}

			topologyMap["dependencies"] = []interface{}{dependenciesList}
		}

		if mediaWorkflow.Topology.Nodes != nil {
			nodesList := []interface{}{}
			for _, nodes := range mediaWorkflow.Topology.Nodes {
				nodesMap := map[string]interface{}{}

				if nodes.Key != nil {
					nodesMap["key"] = nodes.Key
				}

				if nodes.Node != nil {
					nodeMap := map[string]interface{}{}

					if nodes.Node.Type != nil {
						nodeMap["type"] = nodes.Node.Type
					}

					if nodes.Node.Input != nil {
						inputMap := map[string]interface{}{}

						if nodes.Node.Input.QueueId != nil {
							inputMap["queue_id"] = nodes.Node.Input.QueueId
						}

						if nodes.Node.Input.ObjectPrefix != nil {
							inputMap["object_prefix"] = nodes.Node.Input.ObjectPrefix
						}

						if nodes.Node.Input.NotifyConfig != nil {
							notifyConfigMap := map[string]interface{}{}

							if nodes.Node.Input.NotifyConfig.URL != nil {
								notifyConfigMap["u_r_l"] = nodes.Node.Input.NotifyConfig.URL
							}

							if nodes.Node.Input.NotifyConfig.Event != nil {
								notifyConfigMap["event"] = nodes.Node.Input.NotifyConfig.Event
							}

							if nodes.Node.Input.NotifyConfig.Type != nil {
								notifyConfigMap["type"] = nodes.Node.Input.NotifyConfig.Type
							}

							if nodes.Node.Input.NotifyConfig.ResultFormat != nil {
								notifyConfigMap["result_format"] = nodes.Node.Input.NotifyConfig.ResultFormat
							}

							inputMap["notify_config"] = []interface{}{notifyConfigMap}
						}

						if nodes.Node.Input.ExtFilter != nil {
							extFilterMap := map[string]interface{}{}

							if nodes.Node.Input.ExtFilter.State != nil {
								extFilterMap["state"] = nodes.Node.Input.ExtFilter.State
							}

							if nodes.Node.Input.ExtFilter.Audio != nil {
								extFilterMap["audio"] = nodes.Node.Input.ExtFilter.Audio
							}

							if nodes.Node.Input.ExtFilter.Custom != nil {
								extFilterMap["custom"] = nodes.Node.Input.ExtFilter.Custom
							}

							if nodes.Node.Input.ExtFilter.CustomExts != nil {
								extFilterMap["custom_exts"] = nodes.Node.Input.ExtFilter.CustomExts
							}

							if nodes.Node.Input.ExtFilter.AllFile != nil {
								extFilterMap["all_file"] = nodes.Node.Input.ExtFilter.AllFile
							}

							inputMap["ext_filter"] = []interface{}{extFilterMap}
						}

						nodeMap["input"] = []interface{}{inputMap}
					}

					if nodes.Node.Operation != nil {
						operationMap := map[string]interface{}{}

						if nodes.Node.Operation.TemplateId != nil {
							operationMap["template_id"] = nodes.Node.Operation.TemplateId
						}

						if nodes.Node.Operation.Output != nil {
							outputMap := map[string]interface{}{}

							if nodes.Node.Operation.Output.Region != nil {
								outputMap["region"] = nodes.Node.Operation.Output.Region
							}

							if nodes.Node.Operation.Output.Bucket != nil {
								outputMap["bucket"] = nodes.Node.Operation.Output.Bucket
							}

							if nodes.Node.Operation.Output.Object != nil {
								outputMap["object"] = nodes.Node.Operation.Output.Object
							}

							if nodes.Node.Operation.Output.AuObject != nil {
								outputMap["au_object"] = nodes.Node.Operation.Output.AuObject
							}

							if nodes.Node.Operation.Output.SpriteObject != nil {
								outputMap["sprite_object"] = nodes.Node.Operation.Output.SpriteObject
							}

							operationMap["output"] = []interface{}{outputMap}
						}

						if nodes.Node.Operation.WatermarkTemplateId != nil {
							operationMap["watermark_template_id"] = nodes.Node.Operation.WatermarkTemplateId
						}

						if nodes.Node.Operation.DelogoParam != nil {
							delogoParamMap := map[string]interface{}{}

							if nodes.Node.Operation.DelogoParam.Switch != nil {
								delogoParamMap["switch"] = nodes.Node.Operation.DelogoParam.Switch
							}

							if nodes.Node.Operation.DelogoParam.Dx != nil {
								delogoParamMap["dx"] = nodes.Node.Operation.DelogoParam.Dx
							}

							if nodes.Node.Operation.DelogoParam.Dy != nil {
								delogoParamMap["dy"] = nodes.Node.Operation.DelogoParam.Dy
							}

							if nodes.Node.Operation.DelogoParam.Width != nil {
								delogoParamMap["width"] = nodes.Node.Operation.DelogoParam.Width
							}

							if nodes.Node.Operation.DelogoParam.Height != nil {
								delogoParamMap["height"] = nodes.Node.Operation.DelogoParam.Height
							}

							operationMap["delogo_param"] = []interface{}{delogoParamMap}
						}

						if nodes.Node.Operation.SDRtoHDR != nil {
							sDRtoHDRMap := map[string]interface{}{}

							if nodes.Node.Operation.SDRtoHDR.HdrMode != nil {
								sDRtoHDRMap["hdr_mode"] = nodes.Node.Operation.SDRtoHDR.HdrMode
							}

							operationMap["s_d_rto_h_d_r"] = []interface{}{sDRtoHDRMap}
						}

						if nodes.Node.Operation.SCF != nil {
							sCFMap := map[string]interface{}{}

							if nodes.Node.Operation.SCF.Region != nil {
								sCFMap["region"] = nodes.Node.Operation.SCF.Region
							}

							if nodes.Node.Operation.SCF.FunctionName != nil {
								sCFMap["function_name"] = nodes.Node.Operation.SCF.FunctionName
							}

							if nodes.Node.Operation.SCF.Namespace != nil {
								sCFMap["namespace"] = nodes.Node.Operation.SCF.Namespace
							}

							operationMap["s_c_f"] = []interface{}{sCFMap}
						}

						if nodes.Node.Operation.HlsPackInfo != nil {
							hlsPackInfoMap := map[string]interface{}{}

							if nodes.Node.Operation.HlsPackInfo.VideoStreamConfig != nil {
								videoStreamConfigList := []interface{}{}
								for _, videoStreamConfig := range nodes.Node.Operation.HlsPackInfo.VideoStreamConfig {
									videoStreamConfigMap := map[string]interface{}{}

									if videoStreamConfig.VideoStreamName != nil {
										videoStreamConfigMap["video_stream_name"] = videoStreamConfig.VideoStreamName
									}

									if videoStreamConfig.BandWidth != nil {
										videoStreamConfigMap["band_width"] = videoStreamConfig.BandWidth
									}

									videoStreamConfigList = append(videoStreamConfigList, videoStreamConfigMap)
								}

								hlsPackInfoMap["video_stream_config"] = []interface{}{videoStreamConfigList}
							}

							operationMap["hls_pack_info"] = []interface{}{hlsPackInfoMap}
						}

						if nodes.Node.Operation.TranscodeTemplateId != nil {
							operationMap["transcode_template_id"] = nodes.Node.Operation.TranscodeTemplateId
						}

						if nodes.Node.Operation.SmartCover != nil {
							smartCoverMap := map[string]interface{}{}

							if nodes.Node.Operation.SmartCover.Format != nil {
								smartCoverMap["format"] = nodes.Node.Operation.SmartCover.Format
							}

							if nodes.Node.Operation.SmartCover.Width != nil {
								smartCoverMap["width"] = nodes.Node.Operation.SmartCover.Width
							}

							if nodes.Node.Operation.SmartCover.Height != nil {
								smartCoverMap["height"] = nodes.Node.Operation.SmartCover.Height
							}

							if nodes.Node.Operation.SmartCover.Count != nil {
								smartCoverMap["count"] = nodes.Node.Operation.SmartCover.Count
							}

							if nodes.Node.Operation.SmartCover.DeleteDuplicates != nil {
								smartCoverMap["delete_duplicates"] = nodes.Node.Operation.SmartCover.DeleteDuplicates
							}

							operationMap["smart_cover"] = []interface{}{smartCoverMap}
						}

						if nodes.Node.Operation.SegmentConfig != nil {
							segmentConfigMap := map[string]interface{}{}

							if nodes.Node.Operation.SegmentConfig.Format != nil {
								segmentConfigMap["format"] = nodes.Node.Operation.SegmentConfig.Format
							}

							if nodes.Node.Operation.SegmentConfig.Duration != nil {
								segmentConfigMap["duration"] = nodes.Node.Operation.SegmentConfig.Duration
							}

							operationMap["segment_config"] = []interface{}{segmentConfigMap}
						}

						if nodes.Node.Operation.DigitalWatermark != nil {
							digitalWatermarkMap := map[string]interface{}{}

							if nodes.Node.Operation.DigitalWatermark.Message != nil {
								digitalWatermarkMap["message"] = nodes.Node.Operation.DigitalWatermark.Message
							}

							if nodes.Node.Operation.DigitalWatermark.Type != nil {
								digitalWatermarkMap["type"] = nodes.Node.Operation.DigitalWatermark.Type
							}

							if nodes.Node.Operation.DigitalWatermark.Version != nil {
								digitalWatermarkMap["version"] = nodes.Node.Operation.DigitalWatermark.Version
							}

							operationMap["digital_watermark"] = []interface{}{digitalWatermarkMap}
						}

						if nodes.Node.Operation.StreamPackConfigInfo != nil {
							streamPackConfigInfoMap := map[string]interface{}{}

							if nodes.Node.Operation.StreamPackConfigInfo.PackType != nil {
								streamPackConfigInfoMap["pack_type"] = nodes.Node.Operation.StreamPackConfigInfo.PackType
							}

							if nodes.Node.Operation.StreamPackConfigInfo.IgnoreFailedStream != nil {
								streamPackConfigInfoMap["ignore_failed_stream"] = nodes.Node.Operation.StreamPackConfigInfo.IgnoreFailedStream
							}

							if nodes.Node.Operation.StreamPackConfigInfo.ReserveAllStreamNode != nil {
								streamPackConfigInfoMap["reserve_all_stream_node"] = nodes.Node.Operation.StreamPackConfigInfo.ReserveAllStreamNode
							}

							operationMap["stream_pack_config_info"] = []interface{}{streamPackConfigInfoMap}
						}

						if nodes.Node.Operation.StreamPackInfo != nil {
							streamPackInfoMap := map[string]interface{}{}

							if nodes.Node.Operation.StreamPackInfo.VideoStreamConfig != nil {
								videoStreamConfigList := []interface{}{}
								for _, videoStreamConfig := range nodes.Node.Operation.StreamPackInfo.VideoStreamConfig {
									videoStreamConfigMap := map[string]interface{}{}

									if videoStreamConfig.VideoStreamName != nil {
										videoStreamConfigMap["video_stream_name"] = videoStreamConfig.VideoStreamName
									}

									if videoStreamConfig.BandWidth != nil {
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

					nodesMap["node"] = []interface{}{nodeMap}
				}

				nodesList = append(nodesList, nodesMap)
			}

			topologyMap["nodes"] = []interface{}{nodesList}
		}

		_ = d.Set("topology", []interface{}{topologyMap})
	}

	if mediaWorkflow.CreateTime != nil {
		_ = d.Set("create_time", mediaWorkflow.CreateTime)
	}

	if mediaWorkflow.UpdateTime != nil {
		_ = d.Set("update_time", mediaWorkflow.UpdateTime)
	}

	if mediaWorkflow.BucketId != nil {
		_ = d.Set("bucket_id", mediaWorkflow.BucketId)
	}

	return nil
}

func resourceTencentCloudCiMediaWorkflowUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_workflow.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ci.NewUpdateMediaWorkflowRequest()

	mediaWorkflowId := d.Id()

	request.WorkflowId = &workflowId

	immutableArgs := []string{"name", "workflow_id", "state", "topology", "create_time", "update_time", "bucket_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("state") {
		if v, ok := d.GetOk("state"); ok {
			request.State = helper.String(v.(string))
		}
	}

	if d.HasChange("topology") {
		if dMap, ok := helper.InterfacesHeadMap(d, "topology"); ok {
			topology := ci.Topology{}
			if v, ok := dMap["dependencies"]; ok {
				for _, item := range v.([]interface{}) {
					dependenciesMap := item.(map[string]interface{})
					dependencies := ci.Dependencies{}
					if v, ok := dependenciesMap["key"]; ok {
						dependencies.Key = helper.String(v.(string))
					}
					if v, ok := dependenciesMap["value"]; ok {
						dependencies.Value = helper.String(v.(string))
					}
					topology.Dependencies = append(topology.Dependencies, &dependencies)
				}
			}
			if v, ok := dMap["nodes"]; ok {
				for _, item := range v.([]interface{}) {
					nodesMap := item.(map[string]interface{})
					nodes := ci.Nodes{}
					if v, ok := nodesMap["key"]; ok {
						nodes.Key = helper.String(v.(string))
					}
					if nodeMap, ok := helper.InterfaceToMap(nodesMap, "node"); ok {
						nodes := ci.Nodes{}
						if v, ok := nodeMap["type"]; ok {
							nodes.Type = helper.String(v.(string))
						}
						if inputMap, ok := helper.InterfaceToMap(nodeMap, "input"); ok {
							nodeInput := ci.NodeInput{}
							if v, ok := inputMap["queue_id"]; ok {
								nodeInput.QueueId = helper.String(v.(string))
							}
							if v, ok := inputMap["object_prefix"]; ok {
								nodeInput.ObjectPrefix = helper.String(v.(string))
							}
							if notifyConfigMap, ok := helper.InterfaceToMap(inputMap, "notify_config"); ok {
								notifyConfig := ci.NotifyConfig{}
								if v, ok := notifyConfigMap["u_r_l"]; ok {
									notifyConfig.URL = helper.String(v.(string))
								}
								if v, ok := notifyConfigMap["event"]; ok {
									notifyConfig.Event = helper.String(v.(string))
								}
								if v, ok := notifyConfigMap["type"]; ok {
									notifyConfig.Type = helper.String(v.(string))
								}
								if v, ok := notifyConfigMap["result_format"]; ok {
									notifyConfig.ResultFormat = helper.String(v.(string))
								}
								nodeInput.NotifyConfig = &notifyConfig
							}
							if extFilterMap, ok := helper.InterfaceToMap(inputMap, "ext_filter"); ok {
								extFilter := ci.ExtFilter{}
								if v, ok := extFilterMap["state"]; ok {
									extFilter.State = helper.String(v.(string))
								}
								if v, ok := extFilterMap["audio"]; ok {
									extFilter.Audio = helper.String(v.(string))
								}
								if v, ok := extFilterMap["custom"]; ok {
									extFilter.Custom = helper.String(v.(string))
								}
								if v, ok := extFilterMap["custom_exts"]; ok {
									extFilter.CustomExts = helper.String(v.(string))
								}
								if v, ok := extFilterMap["all_file"]; ok {
									extFilter.AllFile = helper.String(v.(string))
								}
								nodeInput.ExtFilter = &extFilter
							}
							nodes.Input = &nodeInput
						}
						if operationMap, ok := helper.InterfaceToMap(nodeMap, "operation"); ok {
							nodeOperation := ci.NodeOperation{}
							if v, ok := operationMap["template_id"]; ok {
								nodeOperation.TemplateId = helper.String(v.(string))
							}
							if outputMap, ok := helper.InterfaceToMap(operationMap, "output"); ok {
								nodeOutput := ci.NodeOutput{}
								if v, ok := outputMap["region"]; ok {
									nodeOutput.Region = helper.String(v.(string))
								}
								if v, ok := outputMap["bucket"]; ok {
									nodeOutput.Bucket = helper.String(v.(string))
								}
								if v, ok := outputMap["object"]; ok {
									nodeOutput.Object = helper.String(v.(string))
								}
								if v, ok := outputMap["au_object"]; ok {
									nodeOutput.AuObject = helper.String(v.(string))
								}
								if v, ok := outputMap["sprite_object"]; ok {
									nodeOutput.SpriteObject = helper.String(v.(string))
								}
								nodeOperation.Output = &nodeOutput
							}
							if v, ok := operationMap["watermark_template_id"]; ok {
								watermarkTemplateIdSet := v.(*schema.Set).List()
								for i := range watermarkTemplateIdSet {
									watermarkTemplateId := watermarkTemplateIdSet[i].(string)
									nodeOperation.WatermarkTemplateId = append(nodeOperation.WatermarkTemplateId, &watermarkTemplateId)
								}
							}
							if delogoParamMap, ok := helper.InterfaceToMap(operationMap, "delogo_param"); ok {
								delogoParam := ci.DelogoParam{}
								if v, ok := delogoParamMap["switch"]; ok {
									delogoParam.Switch = helper.String(v.(string))
								}
								if v, ok := delogoParamMap["dx"]; ok {
									delogoParam.Dx = helper.String(v.(string))
								}
								if v, ok := delogoParamMap["dy"]; ok {
									delogoParam.Dy = helper.String(v.(string))
								}
								if v, ok := delogoParamMap["width"]; ok {
									delogoParam.Width = helper.String(v.(string))
								}
								if v, ok := delogoParamMap["height"]; ok {
									delogoParam.Height = helper.String(v.(string))
								}
								nodeOperation.DelogoParam = &delogoParam
							}
							if sDRtoHDRMap, ok := helper.InterfaceToMap(operationMap, "s_d_rto_h_d_r"); ok {
								nodeSDRtoHDR := ci.NodeSDRtoHDR{}
								if v, ok := sDRtoHDRMap["hdr_mode"]; ok {
									nodeSDRtoHDR.HdrMode = helper.String(v.(string))
								}
								nodeOperation.SDRtoHDR = &nodeSDRtoHDR
							}
							if sCFMap, ok := helper.InterfaceToMap(operationMap, "s_c_f"); ok {
								nodeSCF := ci.NodeSCF{}
								if v, ok := sCFMap["region"]; ok {
									nodeSCF.Region = helper.String(v.(string))
								}
								if v, ok := sCFMap["function_name"]; ok {
									nodeSCF.FunctionName = helper.String(v.(string))
								}
								if v, ok := sCFMap["namespace"]; ok {
									nodeSCF.Namespace = helper.String(v.(string))
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
											videoStreamConfig.VideoStreamName = helper.String(v.(string))
										}
										if v, ok := videoStreamConfigMap["band_width"]; ok {
											videoStreamConfig.BandWidth = helper.String(v.(string))
										}
										nodeHlsPackInfo.VideoStreamConfig = append(nodeHlsPackInfo.VideoStreamConfig, &videoStreamConfig)
									}
								}
								nodeOperation.HlsPackInfo = &nodeHlsPackInfo
							}
							if v, ok := operationMap["transcode_template_id"]; ok {
								nodeOperation.TranscodeTemplateId = helper.String(v.(string))
							}
							if smartCoverMap, ok := helper.InterfaceToMap(operationMap, "smart_cover"); ok {
								nodeSmartCover := ci.NodeSmartCover{}
								if v, ok := smartCoverMap["format"]; ok {
									nodeSmartCover.Format = helper.String(v.(string))
								}
								if v, ok := smartCoverMap["width"]; ok {
									nodeSmartCover.Width = helper.String(v.(string))
								}
								if v, ok := smartCoverMap["height"]; ok {
									nodeSmartCover.Height = helper.String(v.(string))
								}
								if v, ok := smartCoverMap["count"]; ok {
									nodeSmartCover.Count = helper.String(v.(string))
								}
								if v, ok := smartCoverMap["delete_duplicates"]; ok {
									nodeSmartCover.DeleteDuplicates = helper.String(v.(string))
								}
								nodeOperation.SmartCover = &nodeSmartCover
							}
							if segmentConfigMap, ok := helper.InterfaceToMap(operationMap, "segment_config"); ok {
								nodeSegmentConfig := ci.NodeSegmentConfig{}
								if v, ok := segmentConfigMap["format"]; ok {
									nodeSegmentConfig.Format = helper.String(v.(string))
								}
								if v, ok := segmentConfigMap["duration"]; ok {
									nodeSegmentConfig.Duration = helper.String(v.(string))
								}
								nodeOperation.SegmentConfig = &nodeSegmentConfig
							}
							if digitalWatermarkMap, ok := helper.InterfaceToMap(operationMap, "digital_watermark"); ok {
								digitalWatermark := ci.DigitalWatermark{}
								if v, ok := digitalWatermarkMap["message"]; ok {
									digitalWatermark.Message = helper.String(v.(string))
								}
								if v, ok := digitalWatermarkMap["type"]; ok {
									digitalWatermark.Type = helper.String(v.(string))
								}
								if v, ok := digitalWatermarkMap["version"]; ok {
									digitalWatermark.Version = helper.String(v.(string))
								}
								nodeOperation.DigitalWatermark = &digitalWatermark
							}
							if streamPackConfigInfoMap, ok := helper.InterfaceToMap(operationMap, "stream_pack_config_info"); ok {
								nodeStreamPackConfigInfo := ci.NodeStreamPackConfigInfo{}
								if v, ok := streamPackConfigInfoMap["pack_type"]; ok {
									nodeStreamPackConfigInfo.PackType = helper.String(v.(string))
								}
								if v, ok := streamPackConfigInfoMap["ignore_failed_stream"]; ok {
									nodeStreamPackConfigInfo.IgnoreFailedStream = helper.Bool(v.(bool))
								}
								if v, ok := streamPackConfigInfoMap["reserve_all_stream_node"]; ok {
									nodeStreamPackConfigInfo.ReserveAllStreamNode = helper.String(v.(string))
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
											videoStreamConfig.VideoStreamName = helper.String(v.(string))
										}
										if v, ok := videoStreamConfigMap["band_width"]; ok {
											videoStreamConfig.BandWidth = helper.String(v.(string))
										}
										nodeHlsPackInfo.VideoStreamConfig = append(nodeHlsPackInfo.VideoStreamConfig, &videoStreamConfig)
									}
								}
								nodeOperation.StreamPackInfo = &nodeHlsPackInfo
							}
							nodes.Operation = &nodeOperation
						}
						nodes.Node = &nodes
					}
					topology.Nodes = append(topology.Nodes, &nodes)
				}
			}
			request.Topology = &topology
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().UpdateMediaWorkflow(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
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
	mediaWorkflowId := d.Id()

	if err := service.DeleteCiMediaWorkflowById(ctx, workflowId); err != nil {
		return err
	}

	return nil
}
