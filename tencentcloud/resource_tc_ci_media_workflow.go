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
				Description: "Workflow name, supports Chinese, English, numbers, — and _, the length limit is 128 characters.",
			},

			"workflow_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "workflow id.",
			},

			"state": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "workflow state, Paused/Active.",
			},

			"topology": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "topology information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dependencies": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "node dependencies.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "key.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "value.",
									},
								},
							},
						},
						"nodes": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "node list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "node key.",
									},
									"node": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Required:    true,
										Description: "node.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "node type.",
												},
												"input": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Computed:    true,
													Description: "node input.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"queue_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "queue id.",
															},
															"object_prefix": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Object prefix.",
															},
															"notify_config": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Callback information, if not set, use the callback information of the queue.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"url": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Callback address, cannot be an intranet address.",
																		},
																		"event": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Callback information `TaskFinish`: task completed; `WorkflowFinish`: workflow completed; multiple events are supported, separated by commas.",
																		},
																		"type": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Callback type `Url`: url callback; `TDMQ`: tdmq message callback.",
																		},
																		"result_format": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Callback format `XML`: xml format; `JSON`: json format, default: `XML`.",
																		},
																	},
																},
															},
															"ext_filter": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "file extension filter.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"state": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "switch, On/Off, default: `Off`.",
																		},
																		"audio": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Open audio suffix restriction, false/true, default: `false`.",
																		},
																		"custom": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Open custom suffix limit, false/true, default: `false`.",
																		},
																		"custom_exts": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Custom suffix, multiple file suffixes are separated by /, and the number of suffixes does not exceed 10. When Custom is true, this parameter is required.",
																		},
																		"all_file": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "All files, false/true, default: `false`.",
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
													Computed:    true,
													Description: "operating rules.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"template_id": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Video turntable template id.",
															},
															"output": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "output address.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"region": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "The region of the bucket.",
																		},
																		"bucket": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Bucket name.",
																		},
																		"object": {
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Result file name, when the workflow node type is Snapshot or SmartCover, and there are more than one result file, it must contain ${Number}$; when the workflow node type is Segment, Duration is set, and the Format is not HLS or m3u8, Must contain ${Number}.",
																		},
																		"au_object": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Voice result file name, this field is required when the workflow node type is VoiceSeparate and there is a voice output.",
																		},
																		"sprite_object": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The name of the sprite, when the workflow node type is Snapshot and the sprite is opened, this field is required.",
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
																Computed:    true,
																Description: "Watermark template ID, multiple watermark templates can be used, no more than 3.",
															},
															"delogo_param": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "delogo param.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"switch": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "switch.",
																		},
																		"dx": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The x offset from the origin of the upper left corner, value range: [0, 4096], unit: px.",
																		},
																		"dy": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The y offset from the origin of the upper left corner, value range: [0, 4096], unit: px.",
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "watermark width, value range: [0, 4096], unit: px.",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The height of the watermark, value range: [0, 4096], unit: px.",
																		},
																	},
																},
															},
															"sdr_to_hdr": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "SDRtoHDR configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"hdr_mode": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "HDR standard, `HLG`, `HDR10`.",
																		},
																	},
																},
															},
															"scf": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "SCF function information.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"region": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Function region.",
																		},
																		"function_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Function name.",
																		},
																		"namespace": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Namespaces.",
																		},
																	},
																},
															},
															"hls_pack_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "Pack info.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"video_stream_config": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Video substream configuration.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"video_stream_name": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "video substream name, must correspond to an existing video node.",
																					},
																					"band_width": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Video substream bandwidth limit, unit b/s, range [0, 2000000000], 0 means unlimited, greater than or equal to 0, default value is 0.",
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
																Computed:    true,
																Description: "Audio and video transcoding template ID.",
															},
															"smart_cover": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"format": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Image format, `jpg`,`png`,`webp`, default: `jpg`.",
																		},
																		"width": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Width, value range: [128, 4096], unit: px, if only Width is set, Height is calculated according to the original ratio of the video.",
																		},
																		"height": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Height, value range: [128, 4096], unit: px, if only Height is set, Width is calculated according to the original video ratio.",
																		},
																		"count": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Number of screenshots, [1,10], default: `3`.",
																		},
																		"delete_duplicates": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "cover deduplication, false/true, default: `false`.",
																		},
																	},
																},
															},
															"segment_config": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "Audio and video conversion package parameters.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"format": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Package format, `aac`, `mp3`, `flac`, `mp4`, `ts`, `mkv`, `avi`, `hls`, `m3u8`.",
																		},
																		"duration": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Transpackage duration, unit: second, an integer not less than 5.",
																		},
																	},
																},
															},
															"digital_watermark": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: ".",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"message": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Watermark information embedded with digital watermark, the length does not exceed 64 characters, only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
																		},
																		"type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Digital watermark type, currently only can be set to Text.",
																		},
																		"version": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Digital watermark version, currently only can be set to V1.",
																		},
																	},
																},
															},
															"stream_pack_config_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "Packaging configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"pack_type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Packaging type. Default: HLS, HLS/DASH.",
																		},
																		"ignore_failed_stream": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Ignore substreams that fail to be transcoded and continue packaging. Default: true, true/false.",
																		},
																		"reserve_all_stream_node": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Reserve all stream node.",
																		},
																	},
																},
															},
															"stream_pack_info": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Computed:    true,
																Description: "Packing rules.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"video_stream_config": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Video substream configuration.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"video_stream_name": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "video substream name, must correspond to an existing video node.",
																					},
																					"band_width": {
																						Type:        schema.TypeString,
																						Optional:    true,
																						Description: "Video substream bandwidth limit, unit b/s, range [0, 2000000000], 0 means unlimited, greater than or equal to 0, default value is 0.",
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
				Description: "Create time.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},

			"bucket_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Bucket id.",
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
		dependencie := map[string]string{}
		if v, ok := dMap["dependencies"]; ok {
			for _, item := range v.([]interface{}) {
				if item != nil {
					dependenciesMap := item.(map[string]interface{})
					key := ""
					value := ""
					if v, ok := dependenciesMap["key"]; ok {
						key = v.(string)
					}
					if v, ok := dependenciesMap["value"]; ok {
						value = v.(string)
					}
					if key != "" {
						dependencie[key] = value
					}
				}
			}
			topology.Dependencies = dependencie
		}
		topology.Dependencies = dependencie
		if v, ok := dMap["nodes"]; ok {
			for _, item := range v.(*schema.Set).List() {
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
						if sCFMap, ok := helper.InterfaceToMap(operationMap, "scf"); ok {
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
	request.MediaWorkflow = &mediaWorkflow
	log.Printf("[DEBUG]======%s api[%s] success, request body [%+v]\n", logId, "CreateMediaWorkflow", mediaWorkflow.Topology)
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
			for k, v := range mediaWorkflow.Topology.Dependencies {
				dependenciesMap["key"] = k
				dependenciesMap["value"] = v
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

							operationMap["sdr_to_hdr"] = []interface{}{sDRtoHDRMap}
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

							operationMap["scf"] = []interface{}{sCFMap}
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

					nodesMap["key"] = key
					nodesMap["node"] = []interface{}{nodeMap}
				}

				nodesList = append(nodesList, nodesMap)
			}

			topologyMap["nodes"] = nodesList
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

	if v, ok := d.GetOk("name"); ok {
		mediaWorkflow.Name = v.(string)
	}

	if d.HasChange("state") {
		if v, ok := d.GetOk("state"); ok {
			mediaWorkflow.State = v.(string)
		}
	}

	if d.HasChange("topology") {
		if dMap, ok := helper.InterfacesHeadMap(d, "topology"); ok {
			topology := ci.Topology{}
			if v, ok := dMap["dependencies"]; ok {
				dependencie := map[string]string{}
				for _, item := range v.([]interface{}) {
					dependenciesMap := item.(map[string]interface{})
					key := ""
					value := ""
					if v, ok := dependenciesMap["key"]; ok {
						key = v.(string)
					}
					if v, ok := dependenciesMap["value"]; ok {
						value = v.(string)
					}
					if key != "" {
						dependencie[key] = value
					}
				}
				topology.Dependencies = dependencie
			}
			if v, ok := dMap["nodes"]; ok {
				for _, item := range v.(*schema.Set).List() {
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
