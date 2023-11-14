/*
Provides a resource to create a tat invocation_invoke_attachment

Example Usage

```hcl
resource "tencentcloud_tat_invocation_invoke_attachment" "invocation_invoke_attachment" {
  instance_ids =
  working_directory = ""
  timeout =
  parameters = ""
  username = ""
  output_c_o_s_bucket_url = ""
  output_c_o_s_key_prefix = ""
  command_id = ""
}
```

Import

tat invocation_invoke_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tat_invocation_invoke_attachment.invocation_invoke_attachment invocation_invoke_attachment_id
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

func resourceTencentCloudTatInvocationInvokeAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTatInvocationInvokeAttachmentCreate,
		Read:   resourceTencentCloudTatInvocationInvokeAttachmentRead,
		Delete: resourceTencentCloudTatInvocationInvokeAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "IDs of instances about to execute commands. Up to 100 IDs are allowed. Supported instance types:  CVM  LIGHTHOUSE.",
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

			"command_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Command ID.",
			},
		},
	}
}

func resourceTencentCloudTatInvocationInvokeAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invocation_invoke_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = tat.NewInvokeCommandRequest()
		response     = tat.NewInvokeCommandResponse()
		invocationId string
	)
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	if v, ok := d.GetOk("working_directory"); ok {
		request.WorkingDirectory = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("timeout"); ok {
		request.Timeout = helper.IntUint64(v.(int))
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

	if v, ok := d.GetOk("command_id"); ok {
		request.CommandId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTatClient().InvokeCommand(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tat invocationInvokeAttachment failed, reason:%+v", logId, err)
		return err
	}

	invocationId = *response.Response.InvocationId
	d.SetId(invocationId)

	return resourceTencentCloudTatInvocationInvokeAttachmentRead(d, meta)
}

func resourceTencentCloudTatInvocationInvokeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invocation_invoke_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}

	invocationInvokeAttachmentId := d.Id()

	invocationInvokeAttachment, err := service.DescribeTatInvocationInvokeAttachmentById(ctx, invocationId)
	if err != nil {
		return err
	}

	if invocationInvokeAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TatInvocationInvokeAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if invocationInvokeAttachment.InstanceIds != nil {
		_ = d.Set("instance_ids", invocationInvokeAttachment.InstanceIds)
	}

	if invocationInvokeAttachment.WorkingDirectory != nil {
		_ = d.Set("working_directory", invocationInvokeAttachment.WorkingDirectory)
	}

	if invocationInvokeAttachment.Timeout != nil {
		_ = d.Set("timeout", invocationInvokeAttachment.Timeout)
	}

	if invocationInvokeAttachment.Parameters != nil {
		_ = d.Set("parameters", invocationInvokeAttachment.Parameters)
	}

	if invocationInvokeAttachment.Username != nil {
		_ = d.Set("username", invocationInvokeAttachment.Username)
	}

	if invocationInvokeAttachment.OutputCOSBucketUrl != nil {
		_ = d.Set("output_c_o_s_bucket_url", invocationInvokeAttachment.OutputCOSBucketUrl)
	}

	if invocationInvokeAttachment.OutputCOSKeyPrefix != nil {
		_ = d.Set("output_c_o_s_key_prefix", invocationInvokeAttachment.OutputCOSKeyPrefix)
	}

	if invocationInvokeAttachment.CommandId != nil {
		_ = d.Set("command_id", invocationInvokeAttachment.CommandId)
	}

	return nil
}

func resourceTencentCloudTatInvocationInvokeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invocation_invoke_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}
	invocationInvokeAttachmentId := d.Id()

	if err := service.DeleteTatInvocationInvokeAttachmentById(ctx, invocationId); err != nil {
		return err
	}

	return nil
}
