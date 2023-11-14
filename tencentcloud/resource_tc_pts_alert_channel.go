/*
Provides a resource to create a pts alert_channel

Example Usage

```hcl
resource "tencentcloud_pts_alert_channel" "alert_channel" {
  notice_id = &lt;nil&gt;
  project_id = &lt;nil&gt;
  a_m_p_consumer_id = &lt;nil&gt;
            }
```

Import

pts alert_channel can be imported using the id, e.g.

```
terraform import tencentcloud_pts_alert_channel.alert_channel alert_channel_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudPtsAlertChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPtsAlertChannelCreate,
		Read:   resourceTencentCloudPtsAlertChannelRead,
		Update: resourceTencentCloudPtsAlertChannelUpdate,
		Delete: resourceTencentCloudPtsAlertChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"notice_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Notice ID.",
			},

			"project_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Project ID.",
			},

			"a_m_p_consumer_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "AMP Consumer ID.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Status.",
			},

			"created_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creation time.",
			},

			"updated_at": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},

			"app_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "App ID.",
			},

			"uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "User ID.",
			},

			"sub_account_uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Sub-user ID.",
			},
		},
	}
}

func resourceTencentCloudPtsAlertChannelCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_alert_channel.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = pts.NewCreateAlertChannelRequest()
		response  = pts.NewCreateAlertChannelResponse()
		noticeId  string
		projectId string
	)
	if v, ok := d.GetOk("notice_id"); ok {
		noticeId = v.(string)
		request.NoticeId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectId = v.(string)
		request.ProjectId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("a_m_p_consumer_id"); ok {
		request.AMPConsumerId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().CreateAlertChannel(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create pts alertChannel failed, reason:%+v", logId, err)
		return err
	}

	noticeId = *response.Response.NoticeId
	d.SetId(strings.Join([]string{noticeId, projectId}, FILED_SP))

	return resourceTencentCloudPtsAlertChannelRead(d, meta)
}

func resourceTencentCloudPtsAlertChannelRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_alert_channel.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	noticeId := idSplit[0]
	projectId := idSplit[1]

	alertChannel, err := service.DescribePtsAlertChannelById(ctx, noticeId, projectId)
	if err != nil {
		return err
	}

	if alertChannel == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `PtsAlertChannel` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if alertChannel.NoticeId != nil {
		_ = d.Set("notice_id", alertChannel.NoticeId)
	}

	if alertChannel.ProjectId != nil {
		_ = d.Set("project_id", alertChannel.ProjectId)
	}

	if alertChannel.AMPConsumerId != nil {
		_ = d.Set("a_m_p_consumer_id", alertChannel.AMPConsumerId)
	}

	if alertChannel.Status != nil {
		_ = d.Set("status", alertChannel.Status)
	}

	if alertChannel.CreatedAt != nil {
		_ = d.Set("created_at", alertChannel.CreatedAt)
	}

	if alertChannel.UpdatedAt != nil {
		_ = d.Set("updated_at", alertChannel.UpdatedAt)
	}

	if alertChannel.AppId != nil {
		_ = d.Set("app_id", alertChannel.AppId)
	}

	if alertChannel.Uin != nil {
		_ = d.Set("uin", alertChannel.Uin)
	}

	if alertChannel.SubAccountUin != nil {
		_ = d.Set("sub_account_uin", alertChannel.SubAccountUin)
	}

	return nil
}

func resourceTencentCloudPtsAlertChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_alert_channel.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := pts.NewRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	noticeId := idSplit[0]
	projectId := idSplit[1]

	request.NoticeId = &noticeId
	request.ProjectId = &projectId

	immutableArgs := []string{"notice_id", "project_id", "a_m_p_consumer_id", "status", "created_at", "updated_at", "app_id", "uin", "sub_account_uin"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update pts alertChannel failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudPtsAlertChannelRead(d, meta)
}

func resourceTencentCloudPtsAlertChannelDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_pts_alert_channel.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := PtsService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	noticeId := idSplit[0]
	projectId := idSplit[1]

	if err := service.DeletePtsAlertChannelById(ctx, noticeId, projectId); err != nil {
		return err
	}

	return nil
}
