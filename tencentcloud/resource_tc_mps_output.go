/*
Provides a resource to create a mps output

Example Usage

```hcl
resource "tencentcloud_mps_output" "output" {
  flow_id = ""
  output {
		output_name = ""
		description = ""
		protocol = ""
		output_region = ""
		s_r_t_settings {
			destinations {
				ip = ""
				port =
			}
			stream_id = ""
			latency =
			recv_latency =
			peer_latency =
			peer_idle_timeout =
			passphrase = ""
			pb_key_len =
			mode = ""
		}
		r_t_m_p_settings {
			destinations {
				url = ""
				stream_key = ""
			}
			chunk_size =
		}
		r_t_p_settings {
			destinations {
				ip = ""
				port =
			}
			f_e_c = ""
			idle_timeout =
		}
		allow_ip_list =
		max_concurrent =

  }
}
```

Import

mps output can be imported using the id, e.g.

```
terraform import tencentcloud_mps_output.output output_id
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

func resourceTencentCloudMpsOutput() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsOutputCreate,
		Read:   resourceTencentCloudMpsOutputRead,
		Update: resourceTencentCloudMpsOutputUpdate,
		Delete: resourceTencentCloudMpsOutputDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flow_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Flow ID.",
			},

			"output": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Output configuration of the transport stream.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"output_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the output.",
						},
						"description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Output description.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Output protocol, optional [SRT|RTP|RTMP|RTMP_PULL].",
						},
						"output_region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Output region.",
						},
						"s_r_t_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Configuration of the output SRT.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destinations": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The target address of the relay is required when Mode is CALLER, and only one group can be filled in.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Output IP.",
												},
												"port": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Output port.",
												},
											},
										},
									},
									"stream_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Relay the stream ID of SRT. You can choose uppercase and lowercase letters, numbers and special characters (.#!:&amp;amp;,=_-). The length is 0~512.",
									},
									"latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The total delay of relaying SRT, the default is 0, the unit is ms, and the range is [0, 3000].",
									},
									"recv_latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The reception delay of relay SRT, the default is 120, the unit is ms, the range is [0, 3000].",
									},
									"peer_latency": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The peer delay of relaying SRT, the default is 0, the unit is ms, and the range is [0, 3000].",
									},
									"peer_idle_timeout": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The peer idle timeout for relaying SRT, the default is 5000, the unit is ms, and the range is [1000, 10000].",
									},
									"passphrase": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The encryption key for relaying SRT, which is empty by default, indicating no encryption. Only ascii code values can be filled in, and the length is [10, 79].",
									},
									"pb_key_len": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The key length of relay SRT, the default is 0, optional [0|16|24|32].",
									},
									"mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "SRT mode, optional [LISTENER|CALLER], default is CALLER.",
									},
								},
							},
						},
						"r_t_m_p_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Output RTMP configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destinations": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The target address of the relay can be filled in 1~2.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Relayed URL, the format is: rtmp://domain/live.",
												},
												"stream_key": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Relayed StreamKey, in the format: stream?key=value.",
												},
											},
										},
									},
									"chunk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "RTMP Chunk size, range is [4096, 40960].",
									},
								},
							},
						},
						"r_t_p_settings": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Output RTP configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destinations": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The target address of the relay can be filled in 1~2.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ip": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The target IP of the relay.",
												},
												"port": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Destination port for relays.",
												},
											},
										},
									},
									"f_e_c": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "You can only fill in none.",
									},
									"idle_timeout": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Idle timeout, unit ms.",
									},
								},
							},
						},
						"allow_ip_list": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "IP whitelist list, the format is CIDR, such as 0.0.0.0/0. When the Protocol is RTMP_PULL, it is valid, and if it is empty, it means that the client IP is not limited.",
						},
						"max_concurrent": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum number of concurrent pull streams, the maximum is 4, and the default is 4.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMpsOutputCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_output.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mps.NewCreateStreamLinkOutputInfoRequest()
		response = mps.NewCreateStreamLinkOutputInfoResponse()
		outputId string
	)
	if v, ok := d.GetOk("flow_id"); ok {
		request.FlowId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "output"); ok {
		createOutputInfo := mps.CreateOutputInfo{}
		if v, ok := dMap["output_name"]; ok {
			createOutputInfo.OutputName = helper.String(v.(string))
		}
		if v, ok := dMap["description"]; ok {
			createOutputInfo.Description = helper.String(v.(string))
		}
		if v, ok := dMap["protocol"]; ok {
			createOutputInfo.Protocol = helper.String(v.(string))
		}
		if v, ok := dMap["output_region"]; ok {
			createOutputInfo.OutputRegion = helper.String(v.(string))
		}
		if sRTSettingsMap, ok := helper.InterfaceToMap(dMap, "s_r_t_settings"); ok {
			createOutputSRTSettings := mps.CreateOutputSRTSettings{}
			if v, ok := sRTSettingsMap["destinations"]; ok {
				for _, item := range v.([]interface{}) {
					destinationsMap := item.(map[string]interface{})
					createOutputSRTSettingsDestinations := mps.CreateOutputSRTSettingsDestinations{}
					if v, ok := destinationsMap["ip"]; ok {
						createOutputSRTSettingsDestinations.Ip = helper.String(v.(string))
					}
					if v, ok := destinationsMap["port"]; ok {
						createOutputSRTSettingsDestinations.Port = helper.IntInt64(v.(int))
					}
					createOutputSRTSettings.Destinations = append(createOutputSRTSettings.Destinations, &createOutputSRTSettingsDestinations)
				}
			}
			if v, ok := sRTSettingsMap["stream_id"]; ok {
				createOutputSRTSettings.StreamId = helper.String(v.(string))
			}
			if v, ok := sRTSettingsMap["latency"]; ok {
				createOutputSRTSettings.Latency = helper.IntInt64(v.(int))
			}
			if v, ok := sRTSettingsMap["recv_latency"]; ok {
				createOutputSRTSettings.RecvLatency = helper.IntInt64(v.(int))
			}
			if v, ok := sRTSettingsMap["peer_latency"]; ok {
				createOutputSRTSettings.PeerLatency = helper.IntInt64(v.(int))
			}
			if v, ok := sRTSettingsMap["peer_idle_timeout"]; ok {
				createOutputSRTSettings.PeerIdleTimeout = helper.IntInt64(v.(int))
			}
			if v, ok := sRTSettingsMap["passphrase"]; ok {
				createOutputSRTSettings.Passphrase = helper.String(v.(string))
			}
			if v, ok := sRTSettingsMap["pb_key_len"]; ok {
				createOutputSRTSettings.PbKeyLen = helper.IntInt64(v.(int))
			}
			if v, ok := sRTSettingsMap["mode"]; ok {
				createOutputSRTSettings.Mode = helper.String(v.(string))
			}
			createOutputInfo.SRTSettings = &createOutputSRTSettings
		}
		if rTMPSettingsMap, ok := helper.InterfaceToMap(dMap, "r_t_m_p_settings"); ok {
			createOutputRTMPSettings := mps.CreateOutputRTMPSettings{}
			if v, ok := rTMPSettingsMap["destinations"]; ok {
				for _, item := range v.([]interface{}) {
					destinationsMap := item.(map[string]interface{})
					createOutputRtmpSettingsDestinations := mps.CreateOutputRtmpSettingsDestinations{}
					if v, ok := destinationsMap["url"]; ok {
						createOutputRtmpSettingsDestinations.Url = helper.String(v.(string))
					}
					if v, ok := destinationsMap["stream_key"]; ok {
						createOutputRtmpSettingsDestinations.StreamKey = helper.String(v.(string))
					}
					createOutputRTMPSettings.Destinations = append(createOutputRTMPSettings.Destinations, &createOutputRtmpSettingsDestinations)
				}
			}
			if v, ok := rTMPSettingsMap["chunk_size"]; ok {
				createOutputRTMPSettings.ChunkSize = helper.IntInt64(v.(int))
			}
			createOutputInfo.RTMPSettings = &createOutputRTMPSettings
		}
		if rTPSettingsMap, ok := helper.InterfaceToMap(dMap, "r_t_p_settings"); ok {
			createOutputInfoRTPSettings := mps.CreateOutputInfoRTPSettings{}
			if v, ok := rTPSettingsMap["destinations"]; ok {
				for _, item := range v.([]interface{}) {
					destinationsMap := item.(map[string]interface{})
					createOutputRTPSettingsDestinations := mps.CreateOutputRTPSettingsDestinations{}
					if v, ok := destinationsMap["ip"]; ok {
						createOutputRTPSettingsDestinations.Ip = helper.String(v.(string))
					}
					if v, ok := destinationsMap["port"]; ok {
						createOutputRTPSettingsDestinations.Port = helper.IntInt64(v.(int))
					}
					createOutputInfoRTPSettings.Destinations = append(createOutputInfoRTPSettings.Destinations, &createOutputRTPSettingsDestinations)
				}
			}
			if v, ok := rTPSettingsMap["f_e_c"]; ok {
				createOutputInfoRTPSettings.FEC = helper.String(v.(string))
			}
			if v, ok := rTPSettingsMap["idle_timeout"]; ok {
				createOutputInfoRTPSettings.IdleTimeout = helper.IntInt64(v.(int))
			}
			createOutputInfo.RTPSettings = &createOutputInfoRTPSettings
		}
		if v, ok := dMap["allow_ip_list"]; ok {
			allowIpListSet := v.(*schema.Set).List()
			for i := range allowIpListSet {
				allowIpList := allowIpListSet[i].(string)
				createOutputInfo.AllowIpList = append(createOutputInfo.AllowIpList, &allowIpList)
			}
		}
		if v, ok := dMap["max_concurrent"]; ok {
			createOutputInfo.MaxConcurrent = helper.IntUint64(v.(int))
		}
		request.Output = &createOutputInfo
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateStreamLinkOutputInfo(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps output failed, reason:%+v", logId, err)
		return err
	}

	outputId = *response.Response.OutputId
	d.SetId(outputId)

	return resourceTencentCloudMpsOutputRead(d, meta)
}

func resourceTencentCloudMpsOutputRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_output.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	outputId := d.Id()

	output, err := service.DescribeMpsOutputById(ctx, outputId)
	if err != nil {
		return err
	}

	if output == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsOutput` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if output.FlowId != nil {
		_ = d.Set("flow_id", output.FlowId)
	}

	if output.Output != nil {
		outputMap := map[string]interface{}{}

		if output.Output.OutputName != nil {
			outputMap["output_name"] = output.Output.OutputName
		}

		if output.Output.Description != nil {
			outputMap["description"] = output.Output.Description
		}

		if output.Output.Protocol != nil {
			outputMap["protocol"] = output.Output.Protocol
		}

		if output.Output.OutputRegion != nil {
			outputMap["output_region"] = output.Output.OutputRegion
		}

		if output.Output.SRTSettings != nil {
			sRTSettingsMap := map[string]interface{}{}

			if output.Output.SRTSettings.Destinations != nil {
				destinationsList := []interface{}{}
				for _, destinations := range output.Output.SRTSettings.Destinations {
					destinationsMap := map[string]interface{}{}

					if destinations.Ip != nil {
						destinationsMap["ip"] = destinations.Ip
					}

					if destinations.Port != nil {
						destinationsMap["port"] = destinations.Port
					}

					destinationsList = append(destinationsList, destinationsMap)
				}

				sRTSettingsMap["destinations"] = []interface{}{destinationsList}
			}

			if output.Output.SRTSettings.StreamId != nil {
				sRTSettingsMap["stream_id"] = output.Output.SRTSettings.StreamId
			}

			if output.Output.SRTSettings.Latency != nil {
				sRTSettingsMap["latency"] = output.Output.SRTSettings.Latency
			}

			if output.Output.SRTSettings.RecvLatency != nil {
				sRTSettingsMap["recv_latency"] = output.Output.SRTSettings.RecvLatency
			}

			if output.Output.SRTSettings.PeerLatency != nil {
				sRTSettingsMap["peer_latency"] = output.Output.SRTSettings.PeerLatency
			}

			if output.Output.SRTSettings.PeerIdleTimeout != nil {
				sRTSettingsMap["peer_idle_timeout"] = output.Output.SRTSettings.PeerIdleTimeout
			}

			if output.Output.SRTSettings.Passphrase != nil {
				sRTSettingsMap["passphrase"] = output.Output.SRTSettings.Passphrase
			}

			if output.Output.SRTSettings.PbKeyLen != nil {
				sRTSettingsMap["pb_key_len"] = output.Output.SRTSettings.PbKeyLen
			}

			if output.Output.SRTSettings.Mode != nil {
				sRTSettingsMap["mode"] = output.Output.SRTSettings.Mode
			}

			outputMap["s_r_t_settings"] = []interface{}{sRTSettingsMap}
		}

		if output.Output.RTMPSettings != nil {
			rTMPSettingsMap := map[string]interface{}{}

			if output.Output.RTMPSettings.Destinations != nil {
				destinationsList := []interface{}{}
				for _, destinations := range output.Output.RTMPSettings.Destinations {
					destinationsMap := map[string]interface{}{}

					if destinations.Url != nil {
						destinationsMap["url"] = destinations.Url
					}

					if destinations.StreamKey != nil {
						destinationsMap["stream_key"] = destinations.StreamKey
					}

					destinationsList = append(destinationsList, destinationsMap)
				}

				rTMPSettingsMap["destinations"] = []interface{}{destinationsList}
			}

			if output.Output.RTMPSettings.ChunkSize != nil {
				rTMPSettingsMap["chunk_size"] = output.Output.RTMPSettings.ChunkSize
			}

			outputMap["r_t_m_p_settings"] = []interface{}{rTMPSettingsMap}
		}

		if output.Output.RTPSettings != nil {
			rTPSettingsMap := map[string]interface{}{}

			if output.Output.RTPSettings.Destinations != nil {
				destinationsList := []interface{}{}
				for _, destinations := range output.Output.RTPSettings.Destinations {
					destinationsMap := map[string]interface{}{}

					if destinations.Ip != nil {
						destinationsMap["ip"] = destinations.Ip
					}

					if destinations.Port != nil {
						destinationsMap["port"] = destinations.Port
					}

					destinationsList = append(destinationsList, destinationsMap)
				}

				rTPSettingsMap["destinations"] = []interface{}{destinationsList}
			}

			if output.Output.RTPSettings.FEC != nil {
				rTPSettingsMap["f_e_c"] = output.Output.RTPSettings.FEC
			}

			if output.Output.RTPSettings.IdleTimeout != nil {
				rTPSettingsMap["idle_timeout"] = output.Output.RTPSettings.IdleTimeout
			}

			outputMap["r_t_p_settings"] = []interface{}{rTPSettingsMap}
		}

		if output.Output.AllowIpList != nil {
			outputMap["allow_ip_list"] = output.Output.AllowIpList
		}

		if output.Output.MaxConcurrent != nil {
			outputMap["max_concurrent"] = output.Output.MaxConcurrent
		}

		_ = d.Set("output", []interface{}{outputMap})
	}

	return nil
}

func resourceTencentCloudMpsOutputUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_output.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyStreamLinkOutputInfoRequest()

	outputId := d.Id()

	request.OutputId = &outputId

	immutableArgs := []string{"flow_id", "output"}

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

	if d.HasChange("output") {
		if dMap, ok := helper.InterfacesHeadMap(d, "output"); ok {
			createOutputInfo := mps.CreateOutputInfo{}
			if v, ok := dMap["output_name"]; ok {
				createOutputInfo.OutputName = helper.String(v.(string))
			}
			if v, ok := dMap["description"]; ok {
				createOutputInfo.Description = helper.String(v.(string))
			}
			if v, ok := dMap["protocol"]; ok {
				createOutputInfo.Protocol = helper.String(v.(string))
			}
			if v, ok := dMap["output_region"]; ok {
				createOutputInfo.OutputRegion = helper.String(v.(string))
			}
			if sRTSettingsMap, ok := helper.InterfaceToMap(dMap, "s_r_t_settings"); ok {
				createOutputSRTSettings := mps.CreateOutputSRTSettings{}
				if v, ok := sRTSettingsMap["destinations"]; ok {
					for _, item := range v.([]interface{}) {
						destinationsMap := item.(map[string]interface{})
						createOutputSRTSettingsDestinations := mps.CreateOutputSRTSettingsDestinations{}
						if v, ok := destinationsMap["ip"]; ok {
							createOutputSRTSettingsDestinations.Ip = helper.String(v.(string))
						}
						if v, ok := destinationsMap["port"]; ok {
							createOutputSRTSettingsDestinations.Port = helper.IntInt64(v.(int))
						}
						createOutputSRTSettings.Destinations = append(createOutputSRTSettings.Destinations, &createOutputSRTSettingsDestinations)
					}
				}
				if v, ok := sRTSettingsMap["stream_id"]; ok {
					createOutputSRTSettings.StreamId = helper.String(v.(string))
				}
				if v, ok := sRTSettingsMap["latency"]; ok {
					createOutputSRTSettings.Latency = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["recv_latency"]; ok {
					createOutputSRTSettings.RecvLatency = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["peer_latency"]; ok {
					createOutputSRTSettings.PeerLatency = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["peer_idle_timeout"]; ok {
					createOutputSRTSettings.PeerIdleTimeout = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["passphrase"]; ok {
					createOutputSRTSettings.Passphrase = helper.String(v.(string))
				}
				if v, ok := sRTSettingsMap["pb_key_len"]; ok {
					createOutputSRTSettings.PbKeyLen = helper.IntInt64(v.(int))
				}
				if v, ok := sRTSettingsMap["mode"]; ok {
					createOutputSRTSettings.Mode = helper.String(v.(string))
				}
				createOutputInfo.SRTSettings = &createOutputSRTSettings
			}
			if rTMPSettingsMap, ok := helper.InterfaceToMap(dMap, "r_t_m_p_settings"); ok {
				createOutputRTMPSettings := mps.CreateOutputRTMPSettings{}
				if v, ok := rTMPSettingsMap["destinations"]; ok {
					for _, item := range v.([]interface{}) {
						destinationsMap := item.(map[string]interface{})
						createOutputRtmpSettingsDestinations := mps.CreateOutputRtmpSettingsDestinations{}
						if v, ok := destinationsMap["url"]; ok {
							createOutputRtmpSettingsDestinations.Url = helper.String(v.(string))
						}
						if v, ok := destinationsMap["stream_key"]; ok {
							createOutputRtmpSettingsDestinations.StreamKey = helper.String(v.(string))
						}
						createOutputRTMPSettings.Destinations = append(createOutputRTMPSettings.Destinations, &createOutputRtmpSettingsDestinations)
					}
				}
				if v, ok := rTMPSettingsMap["chunk_size"]; ok {
					createOutputRTMPSettings.ChunkSize = helper.IntInt64(v.(int))
				}
				createOutputInfo.RTMPSettings = &createOutputRTMPSettings
			}
			if rTPSettingsMap, ok := helper.InterfaceToMap(dMap, "r_t_p_settings"); ok {
				createOutputInfoRTPSettings := mps.CreateOutputInfoRTPSettings{}
				if v, ok := rTPSettingsMap["destinations"]; ok {
					for _, item := range v.([]interface{}) {
						destinationsMap := item.(map[string]interface{})
						createOutputRTPSettingsDestinations := mps.CreateOutputRTPSettingsDestinations{}
						if v, ok := destinationsMap["ip"]; ok {
							createOutputRTPSettingsDestinations.Ip = helper.String(v.(string))
						}
						if v, ok := destinationsMap["port"]; ok {
							createOutputRTPSettingsDestinations.Port = helper.IntInt64(v.(int))
						}
						createOutputInfoRTPSettings.Destinations = append(createOutputInfoRTPSettings.Destinations, &createOutputRTPSettingsDestinations)
					}
				}
				if v, ok := rTPSettingsMap["f_e_c"]; ok {
					createOutputInfoRTPSettings.FEC = helper.String(v.(string))
				}
				if v, ok := rTPSettingsMap["idle_timeout"]; ok {
					createOutputInfoRTPSettings.IdleTimeout = helper.IntInt64(v.(int))
				}
				createOutputInfo.RTPSettings = &createOutputInfoRTPSettings
			}
			if v, ok := dMap["allow_ip_list"]; ok {
				allowIpListSet := v.(*schema.Set).List()
				for i := range allowIpListSet {
					allowIpList := allowIpListSet[i].(string)
					createOutputInfo.AllowIpList = append(createOutputInfo.AllowIpList, &allowIpList)
				}
			}
			if v, ok := dMap["max_concurrent"]; ok {
				createOutputInfo.MaxConcurrent = helper.IntUint64(v.(int))
			}
			request.Output = &createOutputInfo
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ModifyStreamLinkOutputInfo(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps output failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsOutputRead(d, meta)
}

func resourceTencentCloudMpsOutputDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_output.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	outputId := d.Id()

	if err := service.DeleteMpsOutputById(ctx, outputId); err != nil {
		return err
	}

	return nil
}
