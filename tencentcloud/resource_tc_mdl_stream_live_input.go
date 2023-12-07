package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mdl "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/mdl/v20200326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMdlStreamLiveInput() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMdlStreamLiveInputCreate,
		Read:   resourceTencentCloudMdlStreamLiveInputRead,
		Update: resourceTencentCloudMdlStreamLiveInputUpdate,
		Delete: resourceTencentCloudMdlStreamLiveInputDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Input name, which can contain 1-32 case-sensitive letters, digits, and underscores and must be unique at the region level.",
			},

			"type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Input typeValid values: `RTMP_PUSH`, `RTP_PUSH`, `UDP_PUSH`, `RTMP_PULL`, `HLS_PULL`, `MP4_PULL`.",
			},

			"security_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "ID of the input security group to attachYou can attach only one security group to an input.",
			},

			"input_settings": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Input settings. For the type `RTMP_PUSH`, `RTMP_PULL`, `HLS_PULL`, or `MP4_PULL`, 1 or 2 inputs of the corresponding type can be configured.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Application name, which is valid if `Type` is `RTMP_PUSH` and can contain 1-32 letters and digitsNote: This field may return `null`, indicating that no valid value was found.",
						},
						"stream_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Stream name, which is valid if `Type` is `RTMP_PUSH` and can contain 1-32 letters and digitsNote: This field may return `null`, indicating that no valid value was found.",
						},
						"source_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source URL, which is valid if `Type` is `RTMP_PULL`, `HLS_PULL`, or `MP4_PULL` and can contain 1-512 charactersNote: This field may return `null`, indicating that no valid value was found.",
						},
						"input_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "RTP/UDP input address, which does not need to be entered for the input parameter.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"source_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Source type for stream pulling and relaying. To pull content from private-read COS buckets under the current account, set this parameter to `TencentCOS`; otherwise, leave it empty.Note: this field may return `null`, indicating that no valid value was found.",
						},
						"delay_time": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Delayed time (ms) for playback, which is valid if `Type` is `RTMP_PUSH`Value range: 0 (default) or 10000-600000The value must be a multiple of 1,000.Note: This field may return `null`, indicating that no valid value was found.",
						},
						"input_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The domain of an SRT_PUSH address. If this is a request parameter, you do not need to specify it.Note: This field may return `null`, indicating that no valid value was found.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The username, which is used for authentication.Note: This field may return `null`, indicating that no valid value was found.",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The password, which is used for authentication.Note: This field may return `null`, indicating that no valid value was found.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMdlStreamLiveInputCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mdl_stream_live_input.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mdl.NewCreateStreamLiveInputRequest()
		response = mdl.NewCreateStreamLiveInputResponse()
		id       string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if v, ok := d.GetOk("input_settings"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			inputSettingInfo := mdl.InputSettingInfo{}
			if v, ok := dMap["app_name"]; ok {
				inputSettingInfo.AppName = helper.String(v.(string))
			}
			if v, ok := dMap["stream_name"]; ok {
				inputSettingInfo.StreamName = helper.String(v.(string))
			}
			if v, ok := dMap["source_url"]; ok {
				inputSettingInfo.SourceUrl = helper.String(v.(string))
			}
			if v, ok := dMap["input_address"]; ok {
				inputSettingInfo.InputAddress = helper.String(v.(string))
			}
			if v, ok := dMap["source_type"]; ok {
				inputSettingInfo.SourceType = helper.String(v.(string))
			}
			if v, ok := dMap["delay_time"]; ok {
				inputSettingInfo.DelayTime = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["input_domain"]; ok {
				inputSettingInfo.InputDomain = helper.String(v.(string))
			}
			if v, ok := dMap["user_name"]; ok {
				inputSettingInfo.UserName = helper.String(v.(string))
			}
			if v, ok := dMap["password"]; ok {
				inputSettingInfo.Password = helper.String(v.(string))
			}
			request.InputSettings = append(request.InputSettings, &inputSettingInfo)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMdlClient().CreateStreamLiveInput(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mdl streamliveInput failed, reason:%+v", logId, err)
		return err
	}

	id = *response.Response.Id
	d.SetId(id)

	return resourceTencentCloudMdlStreamLiveInputRead(d, meta)
}

func resourceTencentCloudMdlStreamLiveInputRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mdl_stream_live_input.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MdlService{client: meta.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	streamliveInput, err := service.DescribeMdlStreamLiveInputById(ctx, id)
	if err != nil {
		return err
	}

	if streamliveInput == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MdlStreamliveInput` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if streamliveInput.Name != nil {
		_ = d.Set("name", streamliveInput.Name)
	}

	if streamliveInput.Type != nil {
		_ = d.Set("type", streamliveInput.Type)
	}

	if streamliveInput.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", streamliveInput.SecurityGroupIds)
	}

	if streamliveInput.InputSettings != nil {
		inputSettingsList := []interface{}{}
		for _, inputSettings := range streamliveInput.InputSettings {
			inputSettingsMap := map[string]interface{}{}

			if inputSettings.AppName != nil {
				inputSettingsMap["app_name"] = inputSettings.AppName
			}

			if inputSettings.StreamName != nil {
				inputSettingsMap["stream_name"] = inputSettings.StreamName
			}

			if inputSettings.SourceUrl != nil {
				inputSettingsMap["source_url"] = inputSettings.SourceUrl
			}

			if inputSettings.InputAddress != nil {
				inputSettingsMap["input_address"] = inputSettings.InputAddress
			}

			if inputSettings.SourceType != nil {
				inputSettingsMap["source_type"] = inputSettings.SourceType
			}

			if inputSettings.DelayTime != nil {
				inputSettingsMap["delay_time"] = inputSettings.DelayTime
			}

			if inputSettings.InputDomain != nil {
				inputSettingsMap["input_domain"] = inputSettings.InputDomain
			}

			if inputSettings.UserName != nil {
				inputSettingsMap["user_name"] = inputSettings.UserName
			}

			if inputSettings.Password != nil {
				inputSettingsMap["password"] = inputSettings.Password
			}

			inputSettingsList = append(inputSettingsList, inputSettingsMap)
		}

		_ = d.Set("input_settings", inputSettingsList)

	}

	return nil
}

func resourceTencentCloudMdlStreamLiveInputUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mdl_streamlive_input.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mdl.NewModifyStreamLiveInputRequest()

	id := d.Id()

	request.Id = &id

	needChange := false
	mutableArgs := []string{"name", "security_group_ids", "input_settings"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("security_group_ids"); ok {
			securityGroupIdsSet := v.(*schema.Set).List()
			for i := range securityGroupIdsSet {
				securityGroupIds := securityGroupIdsSet[i].(string)
				request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
			}
		}

		if v, ok := d.GetOk("input_settings"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				inputSettingInfo := mdl.InputSettingInfo{}
				if v, ok := dMap["app_name"]; ok {
					inputSettingInfo.AppName = helper.String(v.(string))
				}
				if v, ok := dMap["stream_name"]; ok {
					inputSettingInfo.StreamName = helper.String(v.(string))
				}
				if v, ok := dMap["source_url"]; ok {
					inputSettingInfo.SourceUrl = helper.String(v.(string))
				}
				if v, ok := dMap["input_address"]; ok {
					inputSettingInfo.InputAddress = helper.String(v.(string))
				}
				if v, ok := dMap["source_type"]; ok {
					inputSettingInfo.SourceType = helper.String(v.(string))
				}
				if v, ok := dMap["delay_time"]; ok {
					inputSettingInfo.DelayTime = helper.IntInt64(v.(int))
				}
				if v, ok := dMap["input_domain"]; ok {
					inputSettingInfo.InputDomain = helper.String(v.(string))
				}
				if v, ok := dMap["user_name"]; ok {
					inputSettingInfo.UserName = helper.String(v.(string))
				}
				if v, ok := dMap["password"]; ok {
					inputSettingInfo.Password = helper.String(v.(string))
				}
				request.InputSettings = append(request.InputSettings, &inputSettingInfo)
			}
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMdlClient().ModifyStreamLiveInput(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mdl streamliveInput failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMdlStreamLiveInputRead(d, meta)
}

func resourceTencentCloudMdlStreamLiveInputDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mdl_stream_live_input.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MdlService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()

	if err := service.DeleteMdlStreamLiveInputById(ctx, id); err != nil {
		return err
	}

	return nil
}
