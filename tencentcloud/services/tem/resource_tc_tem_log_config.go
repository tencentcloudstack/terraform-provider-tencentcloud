package tem

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTemLogConfig() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTemLogConfigRead,
		Create: resourceTencentCloudTemLogConfigCreate,
		Update: resourceTencentCloudTemLogConfigUpdate,
		Delete: resourceTencentCloudTemLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "environment ID.",
			},

			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "application ID.",
			},

			"workload_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "application ID, which is combined by environment ID and application ID, like `en-o5edaepv#app-3j29aa2p`.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "appConfig name.",
			},

			"logset_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "logset.",
			},

			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "topic.",
			},

			"input_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "container_stdout or container_file.",
			},

			"log_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "minimalist_log or multiline_log.",
			},

			"beginning_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "regex pattern.",
			},

			"log_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "directory if container_file.",
			},

			"file_pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "file name pattern if container_file.",
			},
		},
	}
}

func resourceTencentCloudTemLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_log_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request       = tem.NewCreateLogConfigRequest()
		environmentId string
		applicationId string
		name          string
	)

	if v, ok := d.GetOk("environment_id"); ok {
		environmentId = v.(string)
		request.EnvironmentId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_id"); ok {
		applicationId = v.(string)
		request.ApplicationId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("workload_id"); ok {
		workloadId := v.(string)
		if workloadId != environmentId+tccommon.FILED_SP+applicationId {
			return fmt.Errorf("workloadId is error, it should be %s", environmentId+tccommon.FILED_SP+applicationId)
		}
	}

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("logset_id"); ok {
		request.LogsetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("input_type"); ok {
		request.InputType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_type"); ok {
		request.LogType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("beginning_regex"); ok {
		request.BeginningRegex = helper.String(v.(string))
	}

	if v, ok := d.GetOk("log_path"); ok {
		request.LogPath = helper.String(v.(string))
	}

	if v, ok := d.GetOk("file_pattern"); ok {
		request.FilePattern = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTemClient().CreateLogConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tem logConfig failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(environmentId + tccommon.FILED_SP + applicationId + tccommon.FILED_SP + name)
	return resourceTencentCloudTemLogConfigRead(d, meta)
}

func resourceTencentCloudTemLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_logConfig.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TemService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]
	name := idSplit[2]

	logConfig, err := service.DescribeTemLogConfig(ctx, environmentId, applicationId, name)

	if err != nil {
		return err
	}

	if logConfig == nil {
		d.SetId("")
		return fmt.Errorf("resource `logConfig` %s does not exist", name)
	}

	_ = d.Set("environment_id", environmentId)
	_ = d.Set("application_id", logConfig.ApplicationId)
	_ = d.Set("workload_id", environmentId+tccommon.FILED_SP+applicationId)

	if logConfig.Name != nil {
		_ = d.Set("name", logConfig.Name)
	}

	if logConfig.LogsetId != nil {
		_ = d.Set("logset_id", logConfig.LogsetId)
	}

	if logConfig.TopicId != nil {
		_ = d.Set("topic_id", logConfig.TopicId)
	}

	if logConfig.InputType != nil {
		_ = d.Set("input_type", logConfig.InputType)
	}

	if logConfig.LogType != nil {
		_ = d.Set("log_type", logConfig.LogType)
	}

	if logConfig.BeginningRegex != nil {
		_ = d.Set("beginning_regex", logConfig.BeginningRegex)
	}

	if logConfig.LogPath != nil {
		_ = d.Set("log_path", logConfig.LogPath)
	}

	if logConfig.FilePattern != nil {
		_ = d.Set("file_pattern", logConfig.FilePattern)
	}

	return nil
}

func resourceTencentCloudTemLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_log_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tem.NewModifyLogConfigRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]
	name := idSplit[2]

	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId
	request.Name = &name

	if d.HasChange("logset_id") || d.HasChange("topic_id") || d.HasChange("input_type") || d.HasChange("log_type") ||
		d.HasChange("beginning_regex") || d.HasChange("log_path") || d.HasChange("file_pattern") {
		data := tem.LogConfig{}
		if v, ok := d.GetOk("logset_id"); ok {
			data.LogsetId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("topic_id"); ok {
			data.TopicId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("input_type"); ok {
			data.InputType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("log_type"); ok {
			data.LogType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("beginning_regex"); ok {
			data.BeginningRegex = helper.String(v.(string))
		}

		if v, ok := d.GetOk("log_path"); ok {
			data.LogPath = helper.String(v.(string))
		}

		if v, ok := d.GetOk("file_pattern"); ok {
			data.FilePattern = helper.String(v.(string))
		}
		request.Data = &data
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTemClient().ModifyLogConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTemLogConfigRead(d, meta)
}

func resourceTencentCloudTemLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tem_log_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TemService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]
	name := idSplit[2]

	if err := service.DeleteTemLogConfigById(ctx, environmentId, applicationId, name); err != nil {
		return err
	}

	return nil
}
