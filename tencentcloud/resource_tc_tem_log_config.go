/*
Provides a resource to create a tem log_config

Example Usage

```hcl
resource "tencentcloud_tem_log_config" "log_config" {
  environment_id = "en-xxx"
  application_id = "en-xxx"
  name = "xxx"
  logset_id = "xxx"
  topic_id = "xxx"
  input_type = "container_stdout"
  log_type = "minimalist_log"
  beginning_regex = "**.log"
  log_path = "/xxx"
  file_pattern = "*.log"
}
```

Import

tem log_config can be imported using the id, e.g.

```
terraform import tencentcloud_tem_log_config.log_config log_config_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTemLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTemLogConfigCreate,
		Read:   resourceTencentCloudTemLogConfigRead,
		Update: resourceTencentCloudTemLogConfigUpdate,
		Delete: resourceTencentCloudTemLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Environment ID.",
			},

			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application ID.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "AppConfig name.",
			},

			"logset_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Logset.",
			},

			"topic_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Topic.",
			},

			"input_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Container_stdout or container_file.",
			},

			"log_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Minimalist_log or multiline_log.",
			},

			"beginning_regex": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Regex pattern.",
			},

			"log_path": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Directory if container_file.",
			},

			"file_pattern": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "File name pattern if container_file.",
			},
		},
	}
}

func resourceTencentCloudTemLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_log_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tem.NewCreateLogConfigRequest()
		response      = tem.NewCreateLogConfigResponse()
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().CreateLogConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tem logConfig failed, reason:%+v", logId, err)
		return err
	}

	environmentId = *response.Response.EnvironmentId
	d.SetId(strings.Join([]string{environmentId, applicationId, name}, FILED_SP))

	return resourceTencentCloudTemLogConfigRead(d, meta)
}

func resourceTencentCloudTemLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_log_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]
	name := idSplit[2]

	logConfig, err := service.DescribeTemLogConfigById(ctx, environmentId, applicationId, name)
	if err != nil {
		return err
	}

	if logConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TemLogConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if logConfig.EnvironmentId != nil {
		_ = d.Set("environment_id", logConfig.EnvironmentId)
	}

	if logConfig.ApplicationId != nil {
		_ = d.Set("application_id", logConfig.ApplicationId)
	}

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
	defer logElapsed("resource.tencentcloud_tem_log_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tem.NewModifyLogConfigRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]
	name := idSplit[2]

	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId
	request.Name = &name

	immutableArgs := []string{"environment_id", "application_id", "name", "logset_id", "topic_id", "input_type", "log_type", "beginning_regex", "log_path", "file_pattern"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyLogConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tem logConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTemLogConfigRead(d, meta)
}

func resourceTencentCloudTemLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_log_config.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
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
