/*
Provides a resource to create a tat invocation_command_attachment

Example Usage

```hcl
resource "tencentcloud_tat_invocation_command_attachment" "invocation_command_attachment" {
  content = base64encode("pwd")
  instance_id = "ins-881b1c8w"
  command_name = "terraform-test"
  description = "shell test"
  command_type = "SHELL"
  working_directory = "/root"
  timeout = 100
  save_command = false
  enable_parameter = false
  # default_parameters = "{\"varA\": \"222\"}"
  # parameters = "{\"varA\": \"222\"}"
  username = "root"
  output_cos_bucket_url = "https://BucketName-123454321.cos.ap-beijing.myqcloud.com"
  output_cos_key_prefix = "log"
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTatInvocationCommandAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTatInvocationCommandAttachmentCreate,
		Read:   resourceTencentCloudTatInvocationCommandAttachmentRead,
		Delete: resourceTencentCloudTatInvocationCommandAttachmentDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"content": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Base64-encoded command. The maximum length is 64 KB.",
			},

			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of instances about to execute commands. Supported instance types:  CVM  LIGHTHOUSE.",
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
				Default:     "/root",
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

			"output_cos_bucket_url": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The COS bucket URL for uploading logs; The URL must start with https, such as https://BucketName-123454321.cos.ap-beijing.myqcloud.com.",
			},

			"output_cos_key_prefix": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The COS bucket directory where the logs are saved; Check below for the rules of the directory name: 1 It must be a combination of number, letters, and visible characters, Up to 60 characters are allowed; 2 Use a slash (/) to create a subdirectory; 3 can not be used as the folder name; It cannot start with a slash (/), and cannot contain consecutive slashes.",
			},

			"command_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Command ID.",
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
		instanceId   string
	)
	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = []*string{helper.String(v.(string))}
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

	if v, ok := d.GetOk("output_cos_bucket_url"); ok {
		request.OutputCOSBucketUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("output_cos_key_prefix"); ok {
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
	d.SetId(invocationId + FILED_SP + instanceId)

	return resourceTencentCloudTatInvocationCommandAttachmentRead(d, meta)
}

func resourceTencentCloudTatInvocationCommandAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invocation_command_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	invocationId := idSplit[0]
	instanceId := idSplit[1]

	invocation, err := service.DescribeTatInvocationById(ctx, invocationId)
	if err != nil {
		return err
	}

	if invocation == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TatInvocationCommandAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if invocation.CommandContent != nil {
		_ = d.Set("content", invocation.CommandContent)
	}

	_ = d.Set("instance_id", instanceId)

	// if invocation.CommandName != nil {
	// 	_ = d.Set("command_name", invocation.CommandName)
	// }

	// if invocation.Description != nil {
	// 	_ = d.Set("description", invocation.Description)
	// }

	if invocation.CommandType != nil {
		_ = d.Set("command_type", invocation.CommandType)
	}

	if invocation.WorkingDirectory != nil {
		_ = d.Set("working_directory", invocation.WorkingDirectory)
	}

	if invocation.Timeout != nil {
		_ = d.Set("timeout", invocation.Timeout)
	}

	// if invocation.SaveCommand != nil {
	// 	_ = d.Set("save_command", invocation.SaveCommand)
	// }

	// if invocation.EnableParameter != nil {
	// 	_ = d.Set("enable_parameter", invocation.EnableParameter)
	// }

	if invocation.DefaultParameters != nil {
		_ = d.Set("default_parameters", invocation.DefaultParameters)
	}

	if invocation.Parameters != nil {
		_ = d.Set("parameters", invocation.Parameters)
	}

	if invocation.Username != nil {
		_ = d.Set("username", invocation.Username)
	}

	if invocation.OutputCOSBucketUrl != nil {
		_ = d.Set("output_cos_bucket_url", invocation.OutputCOSBucketUrl)
	}

	if invocation.OutputCOSKeyPrefix != nil {
		_ = d.Set("output_cos_key_prefix", invocation.OutputCOSKeyPrefix)
	}

	if invocation.CommandId != nil {
		_ = d.Set("command_id", invocation.CommandId)
	}

	return nil
}

func resourceTencentCloudTatInvocationCommandAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invocation_command_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	invocationId := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteTatInvocationById(ctx, invocationId, instanceId); err != nil {
		return err
	}

	return nil
}
