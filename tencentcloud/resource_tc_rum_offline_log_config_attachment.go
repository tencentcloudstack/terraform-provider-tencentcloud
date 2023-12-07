package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudRumOfflineLogConfigAttachment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudRumOfflineLogConfigAttachmentRead,
		Create: resourceTencentCloudRumOfflineLogConfigAttachmentCreate,
		Delete: resourceTencentCloudRumOfflineLogConfigAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique project key for reporting.",
			},

			"unique_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique identifier of the user to be listened on(aid or uin).",
			},

			"msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Interface call information.",
			},
		},
	}
}

func resourceTencentCloudRumOfflineLogConfigAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_offline_log_config_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = rum.NewCreateOfflineLogConfigRequest()
		projectKey string
		uniqueId   string
	)

	if v, ok := d.GetOk("project_key"); ok {
		projectKey = v.(string)
		request.ProjectKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("unique_id"); ok {
		uniqueId = v.(string)
		request.UniqueID = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseRumClient().CreateOfflineLogConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create rum offlineLogConfigAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(projectKey + FILED_SP + uniqueId)
	return resourceTencentCloudRumOfflineLogConfigAttachmentRead(d, meta)
}

func resourceTencentCloudRumOfflineLogConfigAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_offline_log_config_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectKey := idSplit[0]
	uniqueId := idSplit[1]

	offlineLogConfig, err := service.DescribeRumOfflineLogConfigAttachment(ctx, projectKey, uniqueId)

	if err != nil {
		return err
	}

	if offlineLogConfig == nil {
		d.SetId("")
		return fmt.Errorf("resource `offlineLogConfigAttachment` %s does not exist", d.Id())
	}

	_ = d.Set("project_key", projectKey)
	_ = d.Set("unique_id", uniqueId)

	if offlineLogConfig.Msg != nil {
		_ = d.Set("msg", offlineLogConfig.Msg)
	}

	return nil
}

func resourceTencentCloudRumOfflineLogConfigAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_rum_offline_log_config_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	projectKey := idSplit[0]
	uniqueId := idSplit[1]

	if err := service.DeleteRumOfflineLogConfigAttachmentById(ctx, projectKey, uniqueId); err != nil {
		return err
	}

	return nil
}
