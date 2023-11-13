/*
Provides a resource to create a mps input

Example Usage

```hcl
resource "tencentcloud_mps_input" "input" {
  flow_id = ""
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
}
```

Import

mps input can be imported using the id, e.g.

```
terraform import tencentcloud_mps_input.input input_id
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
				Description: "The input group for the input.",
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
										Description: "Stream ID, optional uppercase and lowercase letters, numbers and special characters (.#!:&amp;amp;,=_-), length 0~512. Specific format can refer to:https://github.com/Haivision/srt/blob/master/docs/features/access-control.md#standard-keys。.",
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
							Description: "Input RTP configuration information.",
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
		},
	}
}

func resourceTencentCloudMpsInputCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_input.create")()
	defer inconsistentCheck(d, meta)()

	var inputId string
	if v, ok := d.GetOk("input_id"); ok {
		inputId = v.(string)
	}

	d.SetId(inputId)

	return resourceTencentCloudMpsInputUpdate(d, meta)
}

func resourceTencentCloudMpsInputRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_input.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	inputId := d.Id()

	input, err := service.DescribeMpsInputById(ctx, inputId)
	if err != nil {
		return err
	}

	if input == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsInput` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if input.FlowId != nil {
		_ = d.Set("flow_id", input.FlowId)
	}

	if input.InputGroup != nil {
		inputGroupList := []interface{}{}
		for _, inputGroup := range input.InputGroup {
			inputGroupMap := map[string]interface{}{}

			if input.InputGroup.InputName != nil {
				inputGroupMap["input_name"] = input.InputGroup.InputName
			}

			if input.InputGroup.Protocol != nil {
				inputGroupMap["protocol"] = input.InputGroup.Protocol
			}

			if input.InputGroup.Description != nil {
				inputGroupMap["description"] = input.InputGroup.Description
			}

			if input.InputGroup.AllowIpList != nil {
				inputGroupMap["allow_ip_list"] = input.InputGroup.AllowIpList
			}

			if input.InputGroup.SRTSettings != nil {
				sRTSettingsMap := map[string]interface{}{}

				if input.InputGroup.SRTSettings.Mode != nil {
					sRTSettingsMap["mode"] = input.InputGroup.SRTSettings.Mode
				}

				if input.InputGroup.SRTSettings.StreamId != nil {
					sRTSettingsMap["stream_id"] = input.InputGroup.SRTSettings.StreamId
				}

				if input.InputGroup.SRTSettings.Latency != nil {
					sRTSettingsMap["latency"] = input.InputGroup.SRTSettings.Latency
				}

				if input.InputGroup.SRTSettings.RecvLatency != nil {
					sRTSettingsMap["recv_latency"] = input.InputGroup.SRTSettings.RecvLatency
				}

				if input.InputGroup.SRTSettings.PeerLatency != nil {
					sRTSettingsMap["peer_latency"] = input.InputGroup.SRTSettings.PeerLatency
				}

				if input.InputGroup.SRTSettings.PeerIdleTimeout != nil {
					sRTSettingsMap["peer_idle_timeout"] = input.InputGroup.SRTSettings.PeerIdleTimeout
				}

				if input.InputGroup.SRTSettings.Passphrase != nil {
					sRTSettingsMap["passphrase"] = input.InputGroup.SRTSettings.Passphrase
				}

				if input.InputGroup.SRTSettings.PbKeyLen != nil {
					sRTSettingsMap["pb_key_len"] = input.InputGroup.SRTSettings.PbKeyLen
				}

				if input.InputGroup.SRTSettings.SourceAddresses != nil {
					sourceAddressesList := []interface{}{}
					for _, sourceAddresses := range input.InputGroup.SRTSettings.SourceAddresses {
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

			if input.InputGroup.RTPSettings != nil {
				rTPSettingsMap := map[string]interface{}{}

				if input.InputGroup.RTPSettings.FEC != nil {
					rTPSettingsMap["f_e_c"] = input.InputGroup.RTPSettings.FEC
				}

				if input.InputGroup.RTPSettings.IdleTimeout != nil {
					rTPSettingsMap["idle_timeout"] = input.InputGroup.RTPSettings.IdleTimeout
				}

				inputGroupMap["r_t_p_settings"] = []interface{}{rTPSettingsMap}
			}

			if input.InputGroup.FailOver != nil {
				inputGroupMap["fail_over"] = input.InputGroup.FailOver
			}

			if input.InputGroup.RTMPPullSettings != nil {
				rTMPPullSettingsMap := map[string]interface{}{}

				if input.InputGroup.RTMPPullSettings.SourceAddresses != nil {
					sourceAddressesList := []interface{}{}
					for _, sourceAddresses := range input.InputGroup.RTMPPullSettings.SourceAddresses {
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

			if input.InputGroup.RTSPPullSettings != nil {
				rTSPPullSettingsMap := map[string]interface{}{}

				if input.InputGroup.RTSPPullSettings.SourceAddresses != nil {
					sourceAddressesList := []interface{}{}
					for _, sourceAddresses := range input.InputGroup.RTSPPullSettings.SourceAddresses {
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

			if input.InputGroup.HLSPullSettings != nil {
				hLSPullSettingsMap := map[string]interface{}{}

				if input.InputGroup.HLSPullSettings.SourceAddresses != nil {
					sourceAddressesList := []interface{}{}
					for _, sourceAddresses := range input.InputGroup.HLSPullSettings.SourceAddresses {
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

			if input.InputGroup.ResilientStream != nil {
				resilientStreamMap := map[string]interface{}{}

				if input.InputGroup.ResilientStream.Enable != nil {
					resilientStreamMap["enable"] = input.InputGroup.ResilientStream.Enable
				}

				if input.InputGroup.ResilientStream.BufferTime != nil {
					resilientStreamMap["buffer_time"] = input.InputGroup.ResilientStream.BufferTime
				}

				inputGroupMap["resilient_stream"] = []interface{}{resilientStreamMap}
			}

			inputGroupList = append(inputGroupList, inputGroupMap)
		}

		_ = d.Set("input_group", inputGroupList)

	}

	return nil
}

func resourceTencentCloudMpsInputUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_input.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyStreamLinkInputRequest()

	inputId := d.Id()

	request.InputId = &inputId

	immutableArgs := []string{"flow_id", "input_group"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("flow_id") {
		if v, ok := d.GetOk("flow_id"); ok {
			request.FlowId = helper.String(v.(string))
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

	return nil
}
