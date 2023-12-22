package tat

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTatInvocationInvokeAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTatInvocationInvokeAttachmentCreate,
		Read:   resourceTencentCloudTatInvocationInvokeAttachmentRead,
		Delete: resourceTencentCloudTatInvocationInvokeAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of instances about to execute commands. Supported instance types:  CVM  LIGHTHOUSE.",
			},

			"working_directory": {
				Optional:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Default:     "/root",
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

			"output_cos_bucket_url": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The COS bucket URL for uploading logs. The URL must start with https, such as https://BucketName-123454321.cos.ap-beijing.myqcloud.com.",
			},

			"output_cos_key_prefix": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The COS bucket directory where the logs are saved; Check below for the rules of the directory name: 1 It must be a combination of number, letters, and visible characters, Up to 60 characters are allowed; 2 Use a slash (/) to create a subdirectory; 3 can not be used as the folder name; It cannot start with a slash (/), and cannot contain consecutive slashes.",
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
	defer tccommon.LogElapsed("resource.tencentcloud_tat_invocation_invoke_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = tat.NewInvokeCommandRequest()
		response     = tat.NewInvokeCommandResponse()
		invocationId string
		instanceId   string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = []*string{helper.String(v.(string))}
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

	if v, ok := d.GetOk("output_cos_bucket_url"); ok {
		request.OutputCOSBucketUrl = helper.String(v.(string))
	}

	if v, ok := d.GetOk("output_cos_key_prefix"); ok {
		request.OutputCOSKeyPrefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("command_id"); ok {
		request.CommandId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTatClient().InvokeCommand(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tat invocation failed, reason:%+v", logId, err)
		return err
	}

	invocationId = *response.Response.InvocationId
	d.SetId(invocationId + tccommon.FILED_SP + instanceId)

	return resourceTencentCloudTatInvocationInvokeAttachmentRead(d, meta)
}

func resourceTencentCloudTatInvocationInvokeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_invocation_invoke_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
		log.Printf("[WARN]%s resource `TatInvocation` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if invocation.WorkingDirectory != nil {
		_ = d.Set("working_directory", invocation.WorkingDirectory)
	}

	if invocation.Timeout != nil {
		_ = d.Set("timeout", invocation.Timeout)
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

func resourceTencentCloudTatInvocationInvokeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_invocation_invoke_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
