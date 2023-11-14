/*
Provides a resource to create a tat command

Example Usage

```hcl
resource "tencentcloud_tat_command" "command" {
  command_name = &lt;nil&gt;
  content = &lt;nil&gt;
  description = &lt;nil&gt;
  command_type = &lt;nil&gt;
  working_directory = &lt;nil&gt;
  timeout = &lt;nil&gt;
  enable_parameter = &lt;nil&gt;
  default_parameters = &lt;nil&gt;
  tags {
		key = &lt;nil&gt;
		value = &lt;nil&gt;

  }
  username = &lt;nil&gt;
  output_c_o_s_bucket_url = &lt;nil&gt;
  output_c_o_s_key_prefix = &lt;nil&gt;
        }
```

Import

tat command can be imported using the id, e.g.

```
terraform import tencentcloud_tat_command.command command_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTatCommand() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTatCommandCreate,
		Read:   resourceTencentCloudTatCommandRead,
		Update: resourceTencentCloudTatCommandUpdate,
		Delete: resourceTencentCloudTatCommandDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"command_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Command name. The name can be up to 60 bytes, and contain [a-z], [A-Z], [0-9] and [_-.].",
			},

			"content": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Base64-encoded command. The maximum length is 64 KB.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Command description. The maximum length is 120 characters.",
			},

			"command_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Command type. `SHELL` and `POWERSHELL` are supported. The default value is `SHELL`.",
			},

			"working_directory": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Command execution path. The default value is /root for `SHELL` commands and C:/Program Files/qcloudtat_agent/workdir for `POWERSHELL` commands.",
			},

			"timeout": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Command timeout period. Default value: 60 seconds. Value range: [1, 86400].",
			},

			"enable_parameter": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable the custom parameter feature.This cannot be modified once created.Default value: `false`.",
			},

			"default_parameters": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The default value of the custom parameter value when it is enabled. The field type is JSON encoded string. For example, {&amp;amp;#39;varA&amp;amp;#39;: &amp;amp;#39;222&amp;amp;#39;}.`key` is the name of the custom parameter and value is the default value. Both `key` and `value` are strings.If no parameter value is provided in the `InvokeCommand` API, the default value is used.Up to 20 custom parameters are supported.The name of the custom parameter cannot exceed 64 characters and can contain [a-z], [A-Z], [0-9] and [-_].",
			},

			"tags": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Tags bound to the command. At most 10 tags are allowed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"username": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The username used to execute the command on the CVM or Lighthouse instance.The principle of least privilege is the best practice for permission management. We recommend you execute TAT commands as a general user. By default, the root user is used to execute commands on Linux and the System user is used on Windows.",
			},

			"output_c_o_s_bucket_url": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The COS bucket URL for uploading logs. The URL must start with `https`, such as `https://BucketName-123454321.cos.ap-beijing.myqcloud.com`.",
			},

			"output_c_o_s_key_prefix": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The COS bucket directory where the logs are saved. Check below for the rules of the directory name.1. It must be a combination of number, letters, and visible characters. Up to 60 characters are allowed.2. Use a slash (/) to create a subdirectory.3. Consecutive dots (.) and slashes (/) are not allowed. It can not start with a slash (/).",
			},

			"created_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Command creation time.",
			},

			"updated_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Command update time.",
			},

			"formatted_description": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Formatted description of the command. This parameter is an empty string for user commands and contains values for public commands.",
			},

			"created_by": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Command creator. `TAT` indicates a public command and `USER` indicates a personal command.",
			},
		},
	}
}

func resourceTencentCloudTatCommandCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_command.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tat.NewCreateCommandRequest()
		response  = tat.NewCreateCommandResponse()
		commandId string
	)
	if v, ok := d.GetOk("command_name"); ok {
		request.CommandName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("command_type"); ok {
		request.CommandType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("working_directory"); ok {
		request.WorkingDirectory = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("timeout"); ok {
		request.Timeout = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("enable_parameter"); ok {
		request.EnableParameter = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("default_parameters"); ok {
		request.DefaultParameters = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tag := tat.Tag{}
			if v, ok := dMap["key"]; ok {
				tag.Key = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				tag.Value = helper.String(v.(string))
			}
			request.Tags = append(request.Tags, &tag)
		}
	}

	if v, ok := d.GetOk("username"); ok {
		request.Username = helper.String(v.(string))
	}

	if v, ok := d.GetOk("output_c_o_s_bucket_url"); ok {
		request.OutputCOSBucketUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("output_c_o_s_key_prefix"); ok {
		request.OutputCOSKeyPrefix = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTatClient().CreateCommand(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tat command failed, reason:%+v", logId, err)
		return err
	}

	commandId = *response.Response.CommandId
	d.SetId(commandId)

	return resourceTencentCloudTatCommandRead(d, meta)
}

func resourceTencentCloudTatCommandRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_command.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}

	commandId := d.Id()

	command, err := service.DescribeTatCommandById(ctx, commandId)
	if err != nil {
		return err
	}

	if command == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TatCommand` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if command.CommandName != nil {
		_ = d.Set("command_name", command.CommandName)
	}

	if command.Content != nil {
		_ = d.Set("content", command.Content)
	}

	if command.Description != nil {
		_ = d.Set("description", command.Description)
	}

	if command.CommandType != nil {
		_ = d.Set("command_type", command.CommandType)
	}

	if command.WorkingDirectory != nil {
		_ = d.Set("working_directory", command.WorkingDirectory)
	}

	if command.Timeout != nil {
		_ = d.Set("timeout", command.Timeout)
	}

	if command.EnableParameter != nil {
		_ = d.Set("enable_parameter", command.EnableParameter)
	}

	if command.DefaultParameters != nil {
		_ = d.Set("default_parameters", command.DefaultParameters)
	}

	if command.Tags != nil {
		tagsList := []interface{}{}
		for _, tags := range command.Tags {
			tagsMap := map[string]interface{}{}

			if command.Tags.Key != nil {
				tagsMap["key"] = command.Tags.Key
			}

			if command.Tags.Value != nil {
				tagsMap["value"] = command.Tags.Value
			}

			tagsList = append(tagsList, tagsMap)
		}

		_ = d.Set("tags", tagsList)

	}

	if command.Username != nil {
		_ = d.Set("username", command.Username)
	}

	if command.OutputCOSBucketUrl != nil {
		_ = d.Set("output_c_o_s_bucket_url", command.OutputCOSBucketUrl)
	}

	if command.OutputCOSKeyPrefix != nil {
		_ = d.Set("output_c_o_s_key_prefix", command.OutputCOSKeyPrefix)
	}

	if command.CreatedTime != nil {
		_ = d.Set("created_time", command.CreatedTime)
	}

	if command.UpdatedTime != nil {
		_ = d.Set("updated_time", command.UpdatedTime)
	}

	if command.FormattedDescription != nil {
		_ = d.Set("formatted_description", command.FormattedDescription)
	}

	if command.CreatedBy != nil {
		_ = d.Set("created_by", command.CreatedBy)
	}

	return nil
}

func resourceTencentCloudTatCommandUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_command.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tat.NewModifyCommandRequest()

	commandId := d.Id()

	request.CommandId = &commandId

	immutableArgs := []string{"command_name", "content", "description", "command_type", "working_directory", "timeout", "enable_parameter", "default_parameters", "tags", "username", "output_c_o_s_bucket_url", "output_c_o_s_key_prefix", "created_time", "updated_time", "formatted_description", "created_by"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("command_name") {
		if v, ok := d.GetOk("command_name"); ok {
			request.CommandName = helper.String(v.(string))
		}
	}

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("command_type") {
		if v, ok := d.GetOk("command_type"); ok {
			request.CommandType = helper.String(v.(string))
		}
	}

	if d.HasChange("working_directory") {
		if v, ok := d.GetOk("working_directory"); ok {
			request.WorkingDirectory = helper.String(v.(string))
		}
	}

	if d.HasChange("timeout") {
		if v, ok := d.GetOkExists("timeout"); ok {
			request.Timeout = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("enable_parameter") {
		if v, ok := d.GetOkExists("enable_parameter"); ok {
			request.EnableParameter = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("default_parameters") {
		if v, ok := d.GetOk("default_parameters"); ok {
			request.DefaultParameters = helper.String(v.(string))
		}
	}

	if d.HasChange("username") {
		if v, ok := d.GetOk("username"); ok {
			request.Username = helper.String(v.(string))
		}
	}

	if d.HasChange("output_c_o_s_bucket_url") {
		if v, ok := d.GetOk("output_c_o_s_bucket_url"); ok {
			request.OutputCOSBucketUrl = helper.String(v.(string))
		}
	}

	if d.HasChange("output_c_o_s_key_prefix") {
		if v, ok := d.GetOk("output_c_o_s_key_prefix"); ok {
			request.OutputCOSKeyPrefix = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTatClient().ModifyCommand(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tat command failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTatCommandRead(d, meta)
}

func resourceTencentCloudTatCommandDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_command.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}
	commandId := d.Id()

	if err := service.DeleteTatCommandById(ctx, commandId); err != nil {
		return err
	}

	return nil
}
