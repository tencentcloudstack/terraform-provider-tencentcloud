/*
Provides a resource to create a tat invocation_command_attachment

Example Usage

```hcl
resource "tencentcloud_tat_invocation_command_attachment" "invocation_command_attachment" {
  content = ""
  instance_ids =
  command_name = ""
  description = ""
  command_type = ""
  working_directory = ""
  timeout =
  save_command =
  enable_parameter =
  default_parameters = ""
  parameters = ""
  username = ""
  output_c_o_s_bucket_url = ""
  output_c_o_s_key_prefix = ""
}
```

Import

tat invocation_command_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tat_invocation_command_attachment.invocation_command_attachment invocation_command_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTatInvocationCommandAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTatInvocationCommandAttachmentCreate,
		Read:   resourceTencentCloudTatInvocationCommandAttachmentRead,
		Delete: resourceTencentCloudTatInvocationCommandAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"content": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Base64-encoded command. The maximum length is 64 KB.",
			},

			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "IDs of instances about to execute commands. Up to 100 IDs are allowed. Supported instance types:CVMLIGHTHOUSE.",
			},

			"command_name": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Command name. The name can be up to 60 bytes, and contain [a-z], [A-Z], [0-9] and [_-.].",
			},

			"description": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Command description. The maximum length is 120 characters.",
			},

			"command_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Command type. SHELL and POWERSHELL are supported. The default value is SHELL.",
			},

			"working_directory": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Command execution path. The default value is /root for SHELL commands and C:Program Filesqcloudtat_agentworkdir for POWERSHELL commands.",
			},

			"timeout": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Command timeout period. Default value: 60 seconds. Value range: [1, 86400].",
			},

			"save_command": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to save the command. Valid values:rue: SaveFalse:Do not saveThe default value is False.",
			},

			"enable_parameter": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable the custom parameter feature.This cannot be modified once created.Default value: false.",
			},

			"default_parameters": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The default value of the custom parameter value when it is enabled. The field type is JSON encoded string. For example, {varA: 222}.key is the name of the custom parameter and value is the default value. Both key and value are strings.If Parameters is not provided, the default values specified here are used.Up to 20 custom parameters are supported.The name of the custom parameter cannot exceed 64 characters and can contain [a-z], [A-Z], [0-9] and [-_].",
			},

			"parameters": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Custom parameters of Command. The field type is JSON encoded string. For example, {varA: 222}.key is the name of the custom parameter and value is the default value. Both key and value are strings.If no parameter value is provided, the DefaultParameters is used.Up to 20 custom parameters are supported.The name of the custom parameter cannot exceed 64 characters and can contain [a-z], [A-Z], [0-9] and [-_].",
			},

			"username": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The username used to execute the command on the CVM or Lighthouse instance.The principle of least privilege is the best practice for permission management. We recommend you execute TAT commands as a general user. By default, the user root is used to execute commands on Linux and the user System is used on Windows.",
			},

			"output_c_o_s_bucket_url": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The COS bucket URL for uploading logs. The URL must start with https, such as https://BucketName-123454321.cos.ap-beijing.myqcloud.com.",
			},

			"output_c_o_s_key_prefix": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The COS bucket directory where the logs are saved. Check below for the rules of the directory name.1. It must be a combination of number, letters, and visible characters. Up to 60 characters are allowed.2. Use a slash (/) to create a subdirectory.3. .. can not be used as the folder name. It cannot start with a slash (/), and cannot contain consecutive slashes.",
			},
		},
	}
}

func resourceTencentCloudTatInvocationCommandAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invocation_command_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = tat.NewRunCommandRequest()
		response     = tat.NewRunCommandResponse()
		invocationId string
	)
	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	if v, ok := d.GetOk("command_name"); ok {
		request.CommandName = helper.String(v.(string))
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

	if v, ok := d.GetOkExists("save_command"); ok {
		request.SaveCommand = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("enable_parameter"); ok {
		request.EnableParameter = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("default_parameters"); ok {
		request.DefaultParameters = helper.String(v.(string))
	}

	if v, ok := d.GetOk("parameters"); ok {
		request.Parameters = helper.String(v.(string))
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTatClient().RunCommand(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tat invocationCommandAttachment failed, reason:%+v", logId, err)
		return err
	}

	invocationId = *response.Response.InvocationId
	d.SetId(invocationId)

	return resourceTencentCloudTatInvocationCommandAttachmentRead(d, meta)
}

func resourceTencentCloudTatInvocationCommandAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invocation_command_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}

	invocationCommandAttachmentId := d.Id()

	invocationCommandAttachment, err := service.DescribeTatInvocationCommandAttachmentById(ctx, invocationId)
	if err != nil {
		return err
	}

	if invocationCommandAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TatInvocationCommandAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if invocationCommandAttachment.Content != nil {
		_ = d.Set("content", invocationCommandAttachment.Content)
	}

	if invocationCommandAttachment.InstanceIds != nil {
		_ = d.Set("instance_ids", invocationCommandAttachment.InstanceIds)
	}

	if invocationCommandAttachment.CommandName != nil {
		_ = d.Set("command_name", invocationCommandAttachment.CommandName)
	}

	if invocationCommandAttachment.Description != nil {
		_ = d.Set("description", invocationCommandAttachment.Description)
	}

	if invocationCommandAttachment.CommandType != nil {
		_ = d.Set("command_type", invocationCommandAttachment.CommandType)
	}

	if invocationCommandAttachment.WorkingDirectory != nil {
		_ = d.Set("working_directory", invocationCommandAttachment.WorkingDirectory)
	}

	if invocationCommandAttachment.Timeout != nil {
		_ = d.Set("timeout", invocationCommandAttachment.Timeout)
	}

	if invocationCommandAttachment.SaveCommand != nil {
		_ = d.Set("save_command", invocationCommandAttachment.SaveCommand)
	}

	if invocationCommandAttachment.EnableParameter != nil {
		_ = d.Set("enable_parameter", invocationCommandAttachment.EnableParameter)
	}

	if invocationCommandAttachment.DefaultParameters != nil {
		_ = d.Set("default_parameters", invocationCommandAttachment.DefaultParameters)
	}

	if invocationCommandAttachment.Parameters != nil {
		_ = d.Set("parameters", invocationCommandAttachment.Parameters)
	}

	if invocationCommandAttachment.Username != nil {
		_ = d.Set("username", invocationCommandAttachment.Username)
	}

	if invocationCommandAttachment.OutputCOSBucketUrl != nil {
		_ = d.Set("output_c_o_s_bucket_url", invocationCommandAttachment.OutputCOSBucketUrl)
	}

	if invocationCommandAttachment.OutputCOSKeyPrefix != nil {
		_ = d.Set("output_c_o_s_key_prefix", invocationCommandAttachment.OutputCOSKeyPrefix)
	}

	return nil
}

func resourceTencentCloudTatInvocationCommandAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invocation_command_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}
	invocationCommandAttachmentId := d.Id()

	if err := service.DeleteTatInvocationCommandAttachmentById(ctx, invocationId); err != nil {
		return err
	}

	return nil
}
