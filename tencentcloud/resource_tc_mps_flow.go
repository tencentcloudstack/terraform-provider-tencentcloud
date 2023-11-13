/*
Provides a resource to create a mps flow

Example Usage

```hcl
resource "tencentcloud_mps_flow" "flow" {
  flow_name = ""
  max_bandwidth =
  input_group {
		input_name = ""
		protocol = ""
		description = ""
		allow_ip_list =
		s_r_t_settings {
			mode = ""
			stream_id = ""
			latency =
			recv_latency =
			peer_latency =
			peer_idle_timeout =
			passphrase = ""
			pb_key_len =
			source_addresses {
				ip = ""
				port =
			}
		}
		r_t_p_settings {
			f_e_c = ""
			idle_timeout =
		}
		fail_over = ""
		r_t_m_p_pull_settings {
			source_addresses {
				tc_url = ""
				stream_key = ""
			}
		}
		r_t_s_p_pull_settings {
			source_addresses {
				url = ""
			}
		}
		h_l_s_pull_settings {
			source_addresses {
				url = ""
			}
		}
		resilient_stream {
			enable =
			buffer_time =
		}

  }
  event_id = ""
}
```

Import

mps flow can be imported using the id, e.g.

```
terraform import tencentcloud_mps_flow.flow flow_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMpsFlow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsFlowCreate,
		Read:   resourceTencentCloudMpsFlowRead,
		Update: resourceTencentCloudMpsFlowUpdate,
		Delete: resourceTencentCloudMpsFlowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flow_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Flow name。.",
			},

			"max_bandwidth": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Maximum bandwidth, unit bps, optional [10000000, 20000000, 50000000].",
			},

			"input_group": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The input group for the flow.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"input_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Input name, you can fill in uppercase and lowercase letters, numbers and underscores, and the length is [1, 32].",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Input protocol, optional [SRT|RTP|RTMP|RTMP_PULL].",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Input description with a length of [0, 255].",
						},
						"allow_ip_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "The input IP whitelist, the format is CIDR.",
						},
						"s_r_t_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The input SRT configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SRT mode, optional [LISTENER|CALLER], default is LISTENER.",
									},
									"stream_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Stream ID, optional uppercase and lowercase letters, numbers and special characters (.#!:&amp;amp;,=_-), length 0~512. For specific format, please refer to:https://github.com/Haivision/srt/blob/master/docs/features/access-control.md#standard-keys。.",
									},
									"latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Delay, default 0, unit ms, range [0, 3000].",
									},
									"recv_latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Receiving delay, default is 120, unit ms, range is [0, 3000].",
									},
									"peer_latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Peer delay, the default is 0, the unit is ms, and the range is [0, 3000].",
									},
									"peer_idle_timeout": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Peer timeout, default is 5000, unit ms, range is [1000, 10000].",
									},
									"passphrase": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The decryption key, which is empty by default, means no encryption. Only ascii code values can be filled in, and the length is [10, 79].",
									},
									"pb_key_len": {
										Type:        schema.TypeInt,
										Optional:    true,
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
						"r_t_p_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "RTP configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"f_e_c": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Defaults to &amp;#39;none&amp;#39;, optional values[&amp;#39;none&amp;#39;].",
									},
									"idle_timeout": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Idle timeout, the default is 5000, the unit is ms, and the range is [1000, 10000].",
									},
								},
							},
						},
						"fail_over": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The active/standby switch of the input, [OPEN|CLOSE] is optional, and the default is CLOSE.",
						},
						"r_t_m_p_pull_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
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
						"r_t_s_p_pull_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
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
						"h_l_s_pull_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
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

			"event_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The event ID associated with this Flow. Each flow can only be associated with one Event.",
			},
		},
	}
}

func resourceTencentCloudMpsFlowCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_flow.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mps.NewCreateStreamLinkFlowRequest()
		response = mps.NewCreateStreamLinkFlowResponse()
		flowId   string
	)
	if v, ok := d.GetOk("flow_name"); ok {
		request.FlowName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_bandwidth"); ok {
		request.MaxBandwidth = helper.IntInt64(v.(int))
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
			}
			if v, ok := dMap["description"]; ok {
				createInput.Description = helper.String(v.(string))
			}
			if v, ok := dMap["allow_ip_list"]; ok {
				allowIpListSet := v.(*schema.Set).List()
				for i := range allowIpListSet {
					allowIpList := allowIpListSet[i].(string)
					createInput.AllowIpList = append(createInput.AllowIpList, &allowIpList)
				}
			}
			if sRTSettingsMap, ok := helper.InterfaceToMap(dMap, "s_r_t_settings"); ok {
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
			if rTPSettingsMap, ok := helper.InterfaceToMap(dMap, "r_t_p_settings"); ok {
				createInputRTPSettings := mps.CreateInputRTPSettings{}
				if v, ok := rTPSettingsMap["f_e_c"]; ok {
					createInputRTPSettings.FEC = helper.String(v.(string))
				}
				if v, ok := rTPSettingsMap["idle_timeout"]; ok {
					createInputRTPSettings.IdleTimeout = helper.IntInt64(v.(int))
				}
				createInput.RTPSettings = &createInputRTPSettings
			}
			if v, ok := dMap["fail_over"]; ok {
				createInput.FailOver = helper.String(v.(string))
			}
			if rTMPPullSettingsMap, ok := helper.InterfaceToMap(dMap, "r_t_m_p_pull_settings"); ok {
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
			if rTSPPullSettingsMap, ok := helper.InterfaceToMap(dMap, "r_t_s_p_pull_settings"); ok {
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
			if hLSPullSettingsMap, ok := helper.InterfaceToMap(dMap, "h_l_s_pull_settings"); ok {
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

	if v, ok := d.GetOk("event_id"); ok {
		request.EventId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateStreamLinkFlow(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps flow failed, reason:%+v", logId, err)
		return err
	}

	flowId = *response.Response.FlowId
	d.SetId(flowId)

	return resourceTencentCloudMpsFlowRead(d, meta)
}

func resourceTencentCloudMpsFlowRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_flow.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	flowId := d.Id()

	flow, err := service.DescribeMpsFlowById(ctx, flowId)
	if err != nil {
		return err
	}

	if flow == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsFlow` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if flow.FlowName != nil {
		_ = d.Set("flow_name", flow.FlowName)
	}

	if flow.MaxBandwidth != nil {
		_ = d.Set("max_bandwidth", flow.MaxBandwidth)
	}

	if flow.InputGroup != nil {
		inputGroupList := []interface{}{}
		for _, inputGroup := range flow.InputGroup {
			inputGroupMap := map[string]interface{}{}

			if flow.InputGroup.InputName != nil {
				inputGroupMap["input_name"] = flow.InputGroup.InputName
			}

			if flow.InputGroup.Protocol != nil {
				inputGroupMap["protocol"] = flow.InputGroup.Protocol
			}

			if flow.InputGroup.Description != nil {
				inputGroupMap["description"] = flow.InputGroup.Description
			}

			if flow.InputGroup.AllowIpList != nil {
				inputGroupMap["allow_ip_list"] = flow.InputGroup.AllowIpList
			}

			if flow.InputGroup.SRTSettings != nil {
				sRTSettingsMap := map[string]interface{}{}

				if flow.InputGroup.SRTSettings.Mode != nil {
					sRTSettingsMap["mode"] = flow.InputGroup.SRTSettings.Mode
				}

				if flow.InputGroup.SRTSettings.StreamId != nil {
					sRTSettingsMap["stream_id"] = flow.InputGroup.SRTSettings.StreamId
				}

				if flow.InputGroup.SRTSettings.Latency != nil {
					sRTSettingsMap["latency"] = flow.InputGroup.SRTSettings.Latency
				}

				if flow.InputGroup.SRTSettings.RecvLatency != nil {
					sRTSettingsMap["recv_latency"] = flow.InputGroup.SRTSettings.RecvLatency
				}

				if flow.InputGroup.SRTSettings.PeerLatency != nil {
					sRTSettingsMap["peer_latency"] = flow.InputGroup.SRTSettings.PeerLatency
				}

				if flow.InputGroup.SRTSettings.PeerIdleTimeout != nil {
					sRTSettingsMap["peer_idle_timeout"] = flow.InputGroup.SRTSettings.PeerIdleTimeout
				}

				if flow.InputGroup.SRTSettings.Passphrase != nil {
					sRTSettingsMap["passphrase"] = flow.InputGroup.SRTSettings.Passphrase
				}

				if flow.InputGroup.SRTSettings.PbKeyLen != nil {
					sRTSettingsMap["pb_key_len"] = flow.InputGroup.SRTSettings.PbKeyLen
				}

				if flow.InputGroup.SRTSettings.SourceAddresses != nil {
					sourceAddressesList := []interface{}{}
					for _, sourceAddresses := range flow.InputGroup.SRTSettings.SourceAddresses {
						sourceAddressesMap := map[string]interface{}{}

						if sourceAddresses.Ip != nil {
							sourceAddressesMap["ip"] = sourceAddresses.Ip
						}

						if sourceAddresses.Port != nil {
							sourceAddressesMap["port"] = sourceAddresses.Port
						}

						sourceAddressesList = append(sourceAddressesList, sourceAddressesMap)
					}

					sRTSettingsMap["source_addresses"] = []interface{}{sourceAddressesList}
				}

				inputGroupMap["s_r_t_settings"] = []interface{}{sRTSettingsMap}
			}

			if flow.InputGroup.RTPSettings != nil {
				rTPSettingsMap := map[string]interface{}{}

				if flow.InputGroup.RTPSettings.FEC != nil {
					rTPSettingsMap["f_e_c"] = flow.InputGroup.RTPSettings.FEC
				}

				if flow.InputGroup.RTPSettings.IdleTimeout != nil {
					rTPSettingsMap["idle_timeout"] = flow.InputGroup.RTPSettings.IdleTimeout
				}

				inputGroupMap["r_t_p_settings"] = []interface{}{rTPSettingsMap}
			}

			if flow.InputGroup.FailOver != nil {
				inputGroupMap["fail_over"] = flow.InputGroup.FailOver
			}

			if flow.InputGroup.RTMPPullSettings != nil {
				rTMPPullSettingsMap := map[string]interface{}{}

				if flow.InputGroup.RTMPPullSettings.SourceAddresses != nil {
					sourceAddressesList := []interface{}{}
					for _, sourceAddresses := range flow.InputGroup.RTMPPullSettings.SourceAddresses {
						sourceAddressesMap := map[string]interface{}{}

						if sourceAddresses.TcUrl != nil {
							sourceAddressesMap["tc_url"] = sourceAddresses.TcUrl
						}

						if sourceAddresses.StreamKey != nil {
							sourceAddressesMap["stream_key"] = sourceAddresses.StreamKey
						}

						sourceAddressesList = append(sourceAddressesList, sourceAddressesMap)
					}

					rTMPPullSettingsMap["source_addresses"] = []interface{}{sourceAddressesList}
				}

				inputGroupMap["r_t_m_p_pull_settings"] = []interface{}{rTMPPullSettingsMap}
			}

			if flow.InputGroup.RTSPPullSettings != nil {
				rTSPPullSettingsMap := map[string]interface{}{}

				if flow.InputGroup.RTSPPullSettings.SourceAddresses != nil {
					sourceAddressesList := []interface{}{}
					for _, sourceAddresses := range flow.InputGroup.RTSPPullSettings.SourceAddresses {
						sourceAddressesMap := map[string]interface{}{}

						if sourceAddresses.Url != nil {
							sourceAddressesMap["url"] = sourceAddresses.Url
						}

						sourceAddressesList = append(sourceAddressesList, sourceAddressesMap)
					}

					rTSPPullSettingsMap["source_addresses"] = []interface{}{sourceAddressesList}
				}

				inputGroupMap["r_t_s_p_pull_settings"] = []interface{}{rTSPPullSettingsMap}
			}

			if flow.InputGroup.HLSPullSettings != nil {
				hLSPullSettingsMap := map[string]interface{}{}

				if flow.InputGroup.HLSPullSettings.SourceAddresses != nil {
					sourceAddressesList := []interface{}{}
					for _, sourceAddresses := range flow.InputGroup.HLSPullSettings.SourceAddresses {
						sourceAddressesMap := map[string]interface{}{}

						if sourceAddresses.Url != nil {
							sourceAddressesMap["url"] = sourceAddresses.Url
						}

						sourceAddressesList = append(sourceAddressesList, sourceAddressesMap)
					}

					hLSPullSettingsMap["source_addresses"] = []interface{}{sourceAddressesList}
				}

				inputGroupMap["h_l_s_pull_settings"] = []interface{}{hLSPullSettingsMap}
			}

			if flow.InputGroup.ResilientStream != nil {
				resilientStreamMap := map[string]interface{}{}

				if flow.InputGroup.ResilientStream.Enable != nil {
					resilientStreamMap["enable"] = flow.InputGroup.ResilientStream.Enable
				}

				if flow.InputGroup.ResilientStream.BufferTime != nil {
					resilientStreamMap["buffer_time"] = flow.InputGroup.ResilientStream.BufferTime
				}

				inputGroupMap["resilient_stream"] = []interface{}{resilientStreamMap}
			}

			inputGroupList = append(inputGroupList, inputGroupMap)
		}

		_ = d.Set("input_group", inputGroupList)

	}

	if flow.EventId != nil {
		_ = d.Set("event_id", flow.EventId)
	}

	return nil
}

func resourceTencentCloudMpsFlowUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_flow.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyStreamLinkFlowRequest()

	flowId := d.Id()

	request.FlowId = &flowId

	immutableArgs := []string{"flow_name", "max_bandwidth", "input_group", "event_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("flow_name") {
		if v, ok := d.GetOk("flow_name"); ok {
			request.FlowName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ModifyStreamLinkFlow(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps flow failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsFlowRead(d, meta)
}

func resourceTencentCloudMpsFlowDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_flow.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	flowId := d.Id()

	if err := service.DeleteMpsFlowById(ctx, flowId); err != nil {
		return err
	}

	return nil
}
