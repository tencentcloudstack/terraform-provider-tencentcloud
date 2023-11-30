package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMpsInput() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsInputCreate,
		Read:   resourceTencentCloudMpsInputRead,
		Update: resourceTencentCloudMpsInputUpdate,
		Delete: resourceTencentCloudMpsInputDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flow_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Flow ID.",
			},

			"input_group": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "The input group for the input. Only support one group for one `tencentcloud_mps_input`. Use `for_each` to create multiple inputs Scenario.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"input_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The input name, you can fill in uppercase and lowercase letters, numbers and underscores, and the length is [1, 32].",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Input protocol, optional [SRT|RTP|RTMP|RTMP_PULL].",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The input description with a length of [0, 255].",
						},
						"allow_ip_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							Description: "The input IP whitelist, the format is CIDR.",
						},
						"srt_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "The input SRT configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "SRT mode, optional [LISTENER|CALLER], default is LISTENER.",
									},
									"stream_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Stream ID, optional uppercase and lowercase letters, numbers and special characters (.#!:&amp;,=_-), length 0~512. Specific format can refer to:https://github.com/Haivision/srt/blob/master/docs/features/access-control.md#standard-keys.",
									},
									"latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Delay, default 0, unit ms, range [0, 3000].",
									},
									"recv_latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Receiving delay, default is 120, unit ms, range is [0, 3000].",
									},
									"peer_latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Peer delay, the default is 0, the unit is ms, and the range is [0, 3000].",
									},
									"peer_idle_timeout": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Peer timeout, default is 5000, unit ms, range is [1000, 10000].",
									},
									"passphrase": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The decryption key, which is empty by default, means no encryption. Only ascii code values can be filled in, and the length is [10, 79].",
									},
									"pb_key_len": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Key length, default is 0, optional [0|16|24|32].",
									},
									"source_addresses": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "SRT peer address, required when Mode is CALLER, and only 1 set can be filled in.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Peer IP.",
												},
												"port": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Peer port.",
												},
											},
										},
									},
								},
							},
						},
						"rtp_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Input RTP configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fec": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Defaults to &#39;none&#39;, optional values[&#39;none&#39;].",
									},
									"idle_timeout": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Idle timeout, the default is 5000, the unit is ms, and the range is [1000, 10000].",
									},
								},
							},
						},
						"fail_over": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The active/standby switch of the input, [OPEN|CLOSE] is optional, and the default is CLOSE.",
						},
						"rtmp_pull_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Input RTMP_PULL configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_addresses": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The source site address of the RTMP source site, there can only be one.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"tc_url": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "TcUrl address of the RTMP source server.",
												},
												"stream_key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "StreamKey information of the RTMP source site.",
												},
											},
										},
									},
								},
							},
						},
						"rtsp_pull_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Input RTSP_PULL configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_addresses": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The source site address of the RTSP source site, there can only be one.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The URL address of the RTSP source site.",
												},
											},
										},
									},
								},
							},
						},
						"hls_pull_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							Description: "Input HLS_PULL configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_addresses": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "There is only one origin address of the HLS origin station.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The URL address of the HLS origin site.",
												},
											},
										},
									},
								},
							},
						},
						"resilient_stream": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Delay broadcast smooth streaming configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to enable the delayed broadcast smooth spit stream, true is enabled, false is not enabled, and the default is not enabled. Note: This field may return null, indicating that no valid value can be obtained.",
									},
									"buffer_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Delay time, in seconds, currently supports a range of 10 to 300 seconds. Note: This field may return null, indicating that no valid value can be obtained.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMpsInputCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_input.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mps.NewCreateStreamLinkInputRequest()
		response = mps.NewCreateStreamLinkInputResponse()
		inputId  string
		flowId   string
		protocol string
	)
	if v, ok := d.GetOk("flow_id"); ok {
		request.FlowId = helper.String(v.(string))
		flowId = v.(string)
	}

	if v, ok := d.GetOk("input_group"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			createInput := mps.CreateInput{}
			if v, ok := dMap["input_name"]; ok {
				createInput.InputName = helper.String(v.(string))
			}
			if v, ok := dMap["protocol"]; ok {
				createInput.Protocol = helper.String(v.(string))
				protocol = v.(string)
			}
			if v, ok := dMap["description"]; ok {
				createInput.Description = helper.String(v.(string))
			}
			if v, ok := dMap["allow_ip_list"]; ok {
				allowIpListSet := v.(*schema.Set).List()
				for i := range allowIpListSet {
					if allowIpListSet[i] != nil {
						allowIpList := allowIpListSet[i].(string)
						createInput.AllowIpList = append(createInput.AllowIpList, &allowIpList)
					}
				}
			}
			if protocol == PROTOCOL_SRT {
				if sRTSettingsMap, ok := helper.InterfaceToMap(dMap, "srt_settings"); ok {
					createInputSRTSettings := mps.CreateInputSRTSettings{}
					if v, ok := sRTSettingsMap["mode"]; ok {
						createInputSRTSettings.Mode = helper.String(v.(string))
					}
					if v, ok := sRTSettingsMap["stream_id"]; ok {
						createInputSRTSettings.StreamId = helper.String(v.(string))
					}
					if v, ok := sRTSettingsMap["latency"]; ok {
						createInputSRTSettings.Latency = helper.IntInt64(v.(int))
					}
					if v, ok := sRTSettingsMap["recv_latency"]; ok {
						createInputSRTSettings.RecvLatency = helper.IntInt64(v.(int))
					}
					if v, ok := sRTSettingsMap["peer_latency"]; ok {
						createInputSRTSettings.PeerLatency = helper.IntInt64(v.(int))
					}
					if v, ok := sRTSettingsMap["peer_idle_timeout"]; ok {
						createInputSRTSettings.PeerIdleTimeout = helper.IntInt64(v.(int))
					}
					if v, ok := sRTSettingsMap["passphrase"]; ok {
						createInputSRTSettings.Passphrase = helper.String(v.(string))
					}
					if v, ok := sRTSettingsMap["pb_key_len"]; ok {
						createInputSRTSettings.PbKeyLen = helper.IntInt64(v.(int))
					}
					if v, ok := sRTSettingsMap["source_addresses"]; ok {
						for _, item := range v.([]interface{}) {
							sourceAddressesMap := item.(map[string]interface{})
							sRTSourceAddressReq := mps.SRTSourceAddressReq{}
							if v, ok := sourceAddressesMap["ip"]; ok {
								sRTSourceAddressReq.Ip = helper.String(v.(string))
							}
							if v, ok := sourceAddressesMap["port"]; ok {
								sRTSourceAddressReq.Port = helper.IntInt64(v.(int))
							}
							createInputSRTSettings.SourceAddresses = append(createInputSRTSettings.SourceAddresses, &sRTSourceAddressReq)
						}
					}
					createInput.SRTSettings = &createInputSRTSettings
				}
			}
			if protocol == PROTOCOL_RTP {
				if rTPSettingsMap, ok := helper.InterfaceToMap(dMap, "rtp_settings"); ok {
					createInputRTPSettings := mps.CreateInputRTPSettings{}
					if v, ok := rTPSettingsMap["fec"]; ok {
						createInputRTPSettings.FEC = helper.String(v.(string))
					}
					if v, ok := rTPSettingsMap["idle_timeout"]; ok {
						createInputRTPSettings.IdleTimeout = helper.IntInt64(v.(int))
					}
					createInput.RTPSettings = &createInputRTPSettings
				}
			}
			if v, ok := dMap["fail_over"]; ok {
				createInput.FailOver = helper.String(v.(string))
			}
			if protocol == PROTOCOL_RTMP || protocol == PROTOCOL_RTMP_PULL {
				if rTMPPullSettingsMap, ok := helper.InterfaceToMap(dMap, "rtmp_pull_settings"); ok {
					createInputRTMPPullSettings := mps.CreateInputRTMPPullSettings{}
					if v, ok := rTMPPullSettingsMap["source_addresses"]; ok {
						for _, item := range v.([]interface{}) {
							sourceAddressesMap := item.(map[string]interface{})
							rTMPPullSourceAddress := mps.RTMPPullSourceAddress{}
							if v, ok := sourceAddressesMap["tc_url"]; ok {
								rTMPPullSourceAddress.TcUrl = helper.String(v.(string))
							}
							if v, ok := sourceAddressesMap["stream_key"]; ok {
								rTMPPullSourceAddress.StreamKey = helper.String(v.(string))
							}
							createInputRTMPPullSettings.SourceAddresses = append(createInputRTMPPullSettings.SourceAddresses, &rTMPPullSourceAddress)
						}
					}
					createInput.RTMPPullSettings = &createInputRTMPPullSettings
				}
			}
			if protocol == PROTOCOL_RTSP_PULL {
				if rTSPPullSettingsMap, ok := helper.InterfaceToMap(dMap, "rtsp_pull_settings"); ok {
					createInputRTSPPullSettings := mps.CreateInputRTSPPullSettings{}
					if v, ok := rTSPPullSettingsMap["source_addresses"]; ok {
						for _, item := range v.([]interface{}) {
							sourceAddressesMap := item.(map[string]interface{})
							rTSPPullSourceAddress := mps.RTSPPullSourceAddress{}
							if v, ok := sourceAddressesMap["url"]; ok {
								rTSPPullSourceAddress.Url = helper.String(v.(string))
							}
							createInputRTSPPullSettings.SourceAddresses = append(createInputRTSPPullSettings.SourceAddresses, &rTSPPullSourceAddress)
						}
					}
					createInput.RTSPPullSettings = &createInputRTSPPullSettings
				}
			}
			if protocol == PROTOCOL_HLS || protocol == PROTOCOL_HLS_PULL {
				if hLSPullSettingsMap, ok := helper.InterfaceToMap(dMap, "hls_pull_settings"); ok {
					createInputHLSPullSettings := mps.CreateInputHLSPullSettings{}
					if v, ok := hLSPullSettingsMap["source_addresses"]; ok {
						for _, item := range v.([]interface{}) {
							sourceAddressesMap := item.(map[string]interface{})
							hLSPullSourceAddress := mps.HLSPullSourceAddress{}
							if v, ok := sourceAddressesMap["url"]; ok {
								hLSPullSourceAddress.Url = helper.String(v.(string))
							}
							createInputHLSPullSettings.SourceAddresses = append(createInputHLSPullSettings.SourceAddresses, &hLSPullSourceAddress)
						}
					}
					createInput.HLSPullSettings = &createInputHLSPullSettings
				}
			}
			if resilientStreamMap, ok := helper.InterfaceToMap(dMap, "resilient_stream"); ok {
				resilientStreamConf := mps.ResilientStreamConf{}
				if v, ok := resilientStreamMap["enable"]; ok {
					resilientStreamConf.Enable = helper.Bool(v.(bool))
				}
				if v, ok := resilientStreamMap["buffer_time"]; ok {
					resilientStreamConf.BufferTime = helper.IntUint64(v.(int))
				}
				createInput.ResilientStream = &resilientStreamConf
			}
			request.InputGroup = append(request.InputGroup, &createInput)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateStreamLinkInput(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps input failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.Info != nil && len(response.Response.Info.InputGroup) > 0 {
		inputId = *response.Response.Info.InputGroup[0].InputId
	}

	d.SetId(strings.Join([]string{flowId, inputId}, FILED_SP))

	return resourceTencentCloudMpsInputRead(d, meta)
}

func resourceTencentCloudMpsInputRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_input.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	var protocol string

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	flowId := idSplit[0]
	inputId := idSplit[1]

	input, err := service.DescribeMpsInputById(ctx, flowId, inputId)
	if err != nil {
		return err
	}

	if input == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource Mps input group [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("flow_id", flowId)

	inputGroupMap := map[string]interface{}{}

	if input.InputName != nil {
		inputGroupMap["input_name"] = input.InputName
	}

	if input.Protocol != nil {
		inputGroupMap["protocol"] = input.Protocol
		protocol = *input.Protocol
	}

	if input.Description != nil {
		inputGroupMap["description"] = input.Description
	}

	if input.AllowIpList != nil {
		inputGroupMap["allow_ip_list"] = input.AllowIpList
	}

	if protocol == PROTOCOL_SRT && input.SRTSettings != nil {
		sRTSettingsMap := map[string]interface{}{}

		if input.SRTSettings.Mode != nil {
			sRTSettingsMap["mode"] = input.SRTSettings.Mode
		}

		if input.SRTSettings.StreamId != nil {
			sRTSettingsMap["stream_id"] = input.SRTSettings.StreamId
		}

		if input.SRTSettings.Latency != nil {
			sRTSettingsMap["latency"] = input.SRTSettings.Latency
		}

		if input.SRTSettings.RecvLatency != nil {
			sRTSettingsMap["recv_latency"] = input.SRTSettings.RecvLatency
		}

		// if input.SRTSettings.PeerLatency != nil {
		// 	sRTSettingsMap["peer_latency"] = input.SRTSettings.PeerLatency
		// }
		//cannot be imported
		if input.SRTSettings.PeerLatency != nil {
			index := fmt.Sprintf("input_group.%d.srt_settings.0.peer_latency", 0)
			oldValue := d.Get(index).(int)
			if *input.SRTSettings.PeerLatency == 0 {
				// need fix: the SDK has bug that cannot return the real value for peer_latency.
				sRTSettingsMap["peer_latency"] = helper.IntInt64(oldValue)
			} else {
				sRTSettingsMap["peer_latency"] = input.SRTSettings.PeerLatency
			}
		}

		if input.SRTSettings.PeerIdleTimeout != nil {
			sRTSettingsMap["peer_idle_timeout"] = input.SRTSettings.PeerIdleTimeout
		}

		if input.SRTSettings.Passphrase != nil {
			sRTSettingsMap["passphrase"] = input.SRTSettings.Passphrase
		}

		if input.SRTSettings.PbKeyLen != nil {
			sRTSettingsMap["pb_key_len"] = input.SRTSettings.PbKeyLen
		}

		if input.SRTSettings.SourceAddresses != nil {
			sourceAddressesList := []interface{}{}
			for _, sourceAddresses := range input.SRTSettings.SourceAddresses {
				sourceAddressesMap := map[string]interface{}{}

				if sourceAddresses.Ip != nil {
					sourceAddressesMap["ip"] = sourceAddresses.Ip
				}

				if sourceAddresses.Port != nil {
					sourceAddressesMap["port"] = sourceAddresses.Port
				}

				sourceAddressesList = append(sourceAddressesList, sourceAddressesMap)
			}

			sRTSettingsMap["source_addresses"] = sourceAddressesList
		}

		inputGroupMap["srt_settings"] = []interface{}{sRTSettingsMap}
	}

	if protocol == PROTOCOL_RTP && input.RTPSettings != nil {
		rTPSettingsMap := map[string]interface{}{}

		if input.RTPSettings.FEC != nil {
			rTPSettingsMap["fec"] = input.RTPSettings.FEC
		}

		if input.RTPSettings.IdleTimeout != nil {
			rTPSettingsMap["idle_timeout"] = input.RTPSettings.IdleTimeout
		}

		inputGroupMap["rtp_settings"] = []interface{}{rTPSettingsMap}
	}

	if input.FailOver != nil {
		inputGroupMap["fail_over"] = input.FailOver
	}

	if (protocol == PROTOCOL_RTMP || protocol == PROTOCOL_RTMP_PULL) && input.RTMPPullSettings != nil {
		rTMPPullSettingsMap := map[string]interface{}{}

		if input.RTMPPullSettings.SourceAddresses != nil {
			sourceAddressesList := []interface{}{}
			for _, sourceAddresses := range input.RTMPPullSettings.SourceAddresses {
				sourceAddressesMap := map[string]interface{}{}

				if sourceAddresses.TcUrl != nil {
					sourceAddressesMap["tc_url"] = sourceAddresses.TcUrl
				}

				if sourceAddresses.StreamKey != nil {
					sourceAddressesMap["stream_key"] = sourceAddresses.StreamKey
				}

				sourceAddressesList = append(sourceAddressesList, sourceAddressesMap)
			}

			rTMPPullSettingsMap["source_addresses"] = sourceAddressesList
		}

		inputGroupMap["rtmp_pull_settings"] = []interface{}{rTMPPullSettingsMap}
	}

	if protocol == PROTOCOL_RTSP_PULL && input.RTSPPullSettings != nil {
		rTSPPullSettingsMap := map[string]interface{}{}

		if input.RTSPPullSettings.SourceAddresses != nil {
			sourceAddressesList := []interface{}{}
			for _, sourceAddresses := range input.RTSPPullSettings.SourceAddresses {
				sourceAddressesMap := map[string]interface{}{}

				if sourceAddresses.Url != nil {
					sourceAddressesMap["url"] = sourceAddresses.Url
				}

				sourceAddressesList = append(sourceAddressesList, sourceAddressesMap)
			}

			rTSPPullSettingsMap["source_addresses"] = sourceAddressesList
		}

		inputGroupMap["rtsp_pull_settings"] = []interface{}{rTSPPullSettingsMap}
	}

	if (protocol == PROTOCOL_HLS || protocol == PROTOCOL_HLS_PULL) && input.HLSPullSettings != nil {
		hLSPullSettingsMap := map[string]interface{}{}

		if input.HLSPullSettings.SourceAddresses != nil {
			sourceAddressesList := []interface{}{}
			for _, sourceAddresses := range input.HLSPullSettings.SourceAddresses {
				sourceAddressesMap := map[string]interface{}{}

				if sourceAddresses.Url != nil {
					sourceAddressesMap["url"] = sourceAddresses.Url
				}

				sourceAddressesList = append(sourceAddressesList, sourceAddressesMap)
			}

			hLSPullSettingsMap["source_addresses"] = sourceAddressesList
		}

		inputGroupMap["hls_pull_settings"] = []interface{}{hLSPullSettingsMap}
	}

	if input.ResilientStream != nil {
		resilientStreamMap := map[string]interface{}{}

		if input.ResilientStream.Enable != nil {
			resilientStreamMap["enable"] = input.ResilientStream.Enable
		}

		if input.ResilientStream.BufferTime != nil {
			resilientStreamMap["buffer_time"] = input.ResilientStream.BufferTime
		}

		inputGroupMap["resilient_stream"] = []interface{}{resilientStreamMap}
	}

	_ = d.Set("input_group", []interface{}{inputGroupMap})

	return nil
}

func resourceTencentCloudMpsInputUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_input.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyStreamLinkInputRequest()
	var protocol string

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	flowId := idSplit[0]
	inputId := idSplit[1]

	request.FlowId = &flowId

	immutableArgs := []string{"flow_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("input_group") {
		if v, ok := d.GetOk("input_group"); ok {
			for _, item := range v.([]interface{}) {
				modifyInput := mps.ModifyInput{}
				modifyInput.InputId = &inputId
				dMap := item.(map[string]interface{})
				if v, ok := dMap["input_name"]; ok {
					modifyInput.InputName = helper.String(v.(string))
				}
				if v, ok := dMap["protocol"]; ok {
					modifyInput.Protocol = helper.String(v.(string))
					protocol = v.(string)
				}
				if v, ok := dMap["description"]; ok {
					modifyInput.Description = helper.String(v.(string))
				}
				if v, ok := dMap["allow_ip_list"]; ok {
					allowIpListSet := v.(*schema.Set).List()
					for i := range allowIpListSet {
						if allowIpListSet[i] != nil {
							allowIpList := allowIpListSet[i].(string)
							modifyInput.AllowIpList = append(modifyInput.AllowIpList, &allowIpList)
						}
					}
				}
				if protocol == PROTOCOL_SRT {
					if sRTSettingsMap, ok := helper.InterfaceToMap(dMap, "srt_settings"); ok {
						createInputSRTSettings := mps.CreateInputSRTSettings{}
						if v, ok := sRTSettingsMap["mode"]; ok {
							createInputSRTSettings.Mode = helper.String(v.(string))
						}
						if v, ok := sRTSettingsMap["stream_id"]; ok {
							createInputSRTSettings.StreamId = helper.String(v.(string))
						}
						if v, ok := sRTSettingsMap["latency"]; ok {
							createInputSRTSettings.Latency = helper.IntInt64(v.(int))
						}
						if v, ok := sRTSettingsMap["recv_latency"]; ok {
							createInputSRTSettings.RecvLatency = helper.IntInt64(v.(int))
						}
						if v, ok := sRTSettingsMap["peer_latency"]; ok {
							createInputSRTSettings.PeerLatency = helper.IntInt64(v.(int))
						}
						if v, ok := sRTSettingsMap["peer_idle_timeout"]; ok {
							createInputSRTSettings.PeerIdleTimeout = helper.IntInt64(v.(int))
						}
						if v, ok := sRTSettingsMap["passphrase"]; ok {
							createInputSRTSettings.Passphrase = helper.String(v.(string))
						}
						if v, ok := sRTSettingsMap["pb_key_len"]; ok {
							createInputSRTSettings.PbKeyLen = helper.IntInt64(v.(int))
						}
						if v, ok := sRTSettingsMap["source_addresses"]; ok {
							for _, item := range v.([]interface{}) {
								sourceAddressesMap := item.(map[string]interface{})
								sRTSourceAddressReq := mps.SRTSourceAddressReq{}
								if v, ok := sourceAddressesMap["ip"]; ok {
									sRTSourceAddressReq.Ip = helper.String(v.(string))
								}
								if v, ok := sourceAddressesMap["port"]; ok {
									sRTSourceAddressReq.Port = helper.IntInt64(v.(int))
								}
								createInputSRTSettings.SourceAddresses = append(createInputSRTSettings.SourceAddresses, &sRTSourceAddressReq)
							}
						}
						modifyInput.SRTSettings = &createInputSRTSettings
					}
				}
				if protocol == PROTOCOL_RTP {
					if rTPSettingsMap, ok := helper.InterfaceToMap(dMap, "rtp_settings"); ok {
						createInputRTPSettings := mps.CreateInputRTPSettings{}
						if v, ok := rTPSettingsMap["fec"]; ok {
							createInputRTPSettings.FEC = helper.String(v.(string))
						}
						if v, ok := rTPSettingsMap["idle_timeout"]; ok {
							createInputRTPSettings.IdleTimeout = helper.IntInt64(v.(int))
						}
						modifyInput.RTPSettings = &createInputRTPSettings
					}
				}
				if v, ok := dMap["fail_over"]; ok {
					modifyInput.FailOver = helper.String(v.(string))
				}
				if protocol == PROTOCOL_RTMP_PULL {
					if rTMPPullSettingsMap, ok := helper.InterfaceToMap(dMap, "rtmp_pull_settings"); ok {
						createInputRTMPPullSettings := mps.CreateInputRTMPPullSettings{}
						if v, ok := rTMPPullSettingsMap["source_addresses"]; ok {
							for _, item := range v.([]interface{}) {
								sourceAddressesMap := item.(map[string]interface{})
								rTMPPullSourceAddress := mps.RTMPPullSourceAddress{}
								if v, ok := sourceAddressesMap["tc_url"]; ok {
									rTMPPullSourceAddress.TcUrl = helper.String(v.(string))
								}
								if v, ok := sourceAddressesMap["stream_key"]; ok {
									rTMPPullSourceAddress.StreamKey = helper.String(v.(string))
								}
								createInputRTMPPullSettings.SourceAddresses = append(createInputRTMPPullSettings.SourceAddresses, &rTMPPullSourceAddress)
							}
						}
						modifyInput.RTMPPullSettings = &createInputRTMPPullSettings
					}
				}
				if protocol == PROTOCOL_RTSP_PULL {
					if rTSPPullSettingsMap, ok := helper.InterfaceToMap(dMap, "rtsp_pull_settings"); ok {
						createInputRTSPPullSettings := mps.CreateInputRTSPPullSettings{}
						if v, ok := rTSPPullSettingsMap["source_addresses"]; ok {
							for _, item := range v.([]interface{}) {
								sourceAddressesMap := item.(map[string]interface{})
								rTSPPullSourceAddress := mps.RTSPPullSourceAddress{}
								if v, ok := sourceAddressesMap["url"]; ok {
									rTSPPullSourceAddress.Url = helper.String(v.(string))
								}
								createInputRTSPPullSettings.SourceAddresses = append(createInputRTSPPullSettings.SourceAddresses, &rTSPPullSourceAddress)
							}
						}
						modifyInput.RTSPPullSettings = &createInputRTSPPullSettings
					}
				}
				if hLSPullSettingsMap, ok := helper.InterfaceToMap(dMap, "hls_pull_settings"); ok {
					createInputHLSPullSettings := mps.CreateInputHLSPullSettings{}
					if v, ok := hLSPullSettingsMap["source_addresses"]; ok {
						for _, item := range v.([]interface{}) {
							sourceAddressesMap := item.(map[string]interface{})
							hLSPullSourceAddress := mps.HLSPullSourceAddress{}
							if v, ok := sourceAddressesMap["url"]; ok {
								hLSPullSourceAddress.Url = helper.String(v.(string))
							}
							createInputHLSPullSettings.SourceAddresses = append(createInputHLSPullSettings.SourceAddresses, &hLSPullSourceAddress)
						}
					}
					modifyInput.HLSPullSettings = &createInputHLSPullSettings
				}
				if resilientStreamMap, ok := helper.InterfaceToMap(dMap, "resilient_stream"); ok {
					resilientStreamConf := mps.ResilientStreamConf{}
					if v, ok := resilientStreamMap["enable"]; ok {
						resilientStreamConf.Enable = helper.Bool(v.(bool))
					}
					if v, ok := resilientStreamMap["buffer_time"]; ok {
						resilientStreamConf.BufferTime = helper.IntUint64(v.(int))
					}
					modifyInput.ResilientStream = &resilientStreamConf
				}
				//  modify api only support to modify one input one time
				request.Input = &modifyInput
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ModifyStreamLinkInput(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps input failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsInputRead(d, meta)
}

func resourceTencentCloudMpsInputDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_input.delete")()
	defer inconsistentCheck(d, meta)()
	// deleted through `tencentcloud_mps_flow`
	return nil
}
