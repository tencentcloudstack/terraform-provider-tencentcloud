package emr

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudEmrBootScriptConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEmrBootScriptConfigCreate,
		Read:   resourceTencentCloudEmrBootScriptConfigRead,
		Update: resourceTencentCloudEmrBootScriptConfigUpdate,
		Delete: resourceTencentCloudEmrBootScriptConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "EMR instance ID.",
			},

			"boot_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Boot type. Valid values: `resourceAfter`, `clusterAfter`, `clusterBefore`.",
			},

			"pre_executed_file_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of pre-execution script settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Script path on COS.",
						},

						"args": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Script execution parameters.",
						},

						"bucket": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "COS bucket name.",
						},

						"region": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "COS region name.",
						},

						"domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "COS domain data.",
						},

						"run_order": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Execution order.",
						},

						"when_run": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Execution timing. Valid values: `resourceAfter`, `clusterAfter`.",
						},

						"cos_file_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Script file name.",
						},

						"cos_file_uri": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Script COS address.",
						},

						"cos_secret_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "COS SecretId.",
						},

						"cos_secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "COS SecretKey.",
						},

						"app_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "COS AppId.",
						},

						"remark": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Remark.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudEmrBootScriptConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_boot_script_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	instanceId := d.Get("instance_id").(string)
	bootType := d.Get("boot_type").(string)
	d.SetId(strings.Join([]string{instanceId, bootType}, tccommon.FILED_SP))

	return resourceTencentCloudEmrBootScriptConfigUpdate(d, meta)
}

func resourceTencentCloudEmrBootScriptConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_boot_script_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	parts := strings.SplitN(d.Id(), tccommon.FILED_SP, 2)
	if len(parts) != 2 {
		return fmt.Errorf("id is broken,%s", parts)
	}

	instanceId := parts[0]
	bootType := parts[1]

	respData, err := service.DescribeEmrBootScript(ctx, instanceId, bootType)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_emr_boot_script_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("boot_type", bootType)

	var rawSettings []*emr.PreExecuteFileSetting
	switch bootType {
	case "resourceAfter":
		rawSettings = respData.ResourceAfter
	case "clusterBefore":
		rawSettings = respData.ClusterBefore
	case "clusterAfter":
		rawSettings = respData.ClusterAfter
	}

	settingsList := make([]map[string]interface{}, 0, len(rawSettings))
	for _, s := range rawSettings {
		settingsMap := map[string]interface{}{}

		if s.Path != nil {
			settingsMap["path"] = s.Path
		}

		if s.Args != nil {
			settingsMap["args"] = s.Args
		}

		if s.Bucket != nil {
			settingsMap["bucket"] = s.Bucket
		}

		if s.Region != nil {
			settingsMap["region"] = s.Region
		}

		if s.Domain != nil {
			settingsMap["domain"] = s.Domain
		}

		if s.RunOrder != nil {
			settingsMap["run_order"] = s.RunOrder
		}

		if s.WhenRun != nil {
			settingsMap["when_run"] = s.WhenRun
		}

		if s.CosFileName != nil {
			settingsMap["cos_file_name"] = s.CosFileName
		}

		if s.CosFileURI != nil {
			settingsMap["cos_file_uri"] = s.CosFileURI
		}

		if s.CosSecretId != nil {
			settingsMap["cos_secret_id"] = s.CosSecretId
		}

		if s.CosSecretKey != nil {
			settingsMap["cos_secret_key"] = s.CosSecretKey
		}

		if s.AppId != nil {
			settingsMap["app_id"] = s.AppId
		}

		if s.Remark != nil {
			settingsMap["remark"] = s.Remark
		}

		settingsList = append(settingsList, settingsMap)
	}

	_ = d.Set("pre_executed_file_settings", settingsList)

	return nil
}

func resourceTencentCloudEmrBootScriptConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_boot_script_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = emr.NewModifyBootScriptRequest()
	)

	parts := strings.SplitN(d.Id(), tccommon.FILED_SP, 2)
	if len(parts) != 2 {
		return fmt.Errorf("id is broken,%s", parts)
	}

	request.InstanceId = helper.String(parts[0])
	request.BootType = helper.String(parts[1])

	if v, ok := d.GetOk("pre_executed_file_settings"); ok {
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			setting := &emr.PreExecuteFileSetting{}

			if v, ok := m["path"].(string); ok && v != "" {
				setting.Path = helper.String(v)
			}

			if v, ok := m["args"].(string); ok && v != "" {
				setting.Args = helper.String(v)
			}

			if v, ok := m["bucket"].(string); ok && v != "" {
				setting.Bucket = helper.String(v)
			}

			if v, ok := m["region"].(string); ok && v != "" {
				setting.Region = helper.String(v)
			}

			if v, ok := m["domain"].(string); ok && v != "" {
				setting.Domain = helper.String(v)
			}

			if v, ok := m["run_order"].(int); ok && v != 0 {
				setting.RunOrder = helper.Int64(int64(v))
			}

			if v, ok := m["when_run"].(string); ok && v != "" {
				setting.WhenRun = helper.String(v)
			}

			if v, ok := m["cos_file_name"].(string); ok && v != "" {
				setting.CosFileName = helper.String(v)
			}

			if v, ok := m["cos_file_uri"].(string); ok && v != "" {
				setting.CosFileURI = helper.String(v)
			}

			if v, ok := m["cos_secret_id"].(string); ok && v != "" {
				setting.CosSecretId = helper.String(v)
			}

			if v, ok := m["cos_secret_key"].(string); ok && v != "" {
				setting.CosSecretKey = helper.String(v)
			}

			if v, ok := m["app_id"].(string); ok && v != "" {
				setting.AppId = helper.String(v)
			}

			if v, ok := m["remark"].(string); ok && v != "" {
				setting.Remark = helper.String(v)
			}

			request.PreExecutedFileSettings = append(request.PreExecutedFileSettings, setting)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseEmrClient().ModifyBootScriptWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update emr boot script config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudEmrBootScriptConfigRead(d, meta)
}

func resourceTencentCloudEmrBootScriptConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_boot_script_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
