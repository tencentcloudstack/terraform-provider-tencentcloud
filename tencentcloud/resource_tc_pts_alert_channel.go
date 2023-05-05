/*
Provides a resource to create a pts alert_channel

~> **NOTE:** Modification is not currently supported, please go to the console to modify.

Example Usage

```hcl
resource "tencentcloud_monitor_alarm_notice" "example" {
  name                  = "test_alarm_notice_1"
  notice_type           = "ALL"
  notice_language       = "zh-CN"

  user_notices    {
      receiver_type              = "USER"
      start_time                 = 0
      end_time                   = 1
      notice_way                 = ["EMAIL", "SMS", "WECHAT"]
      user_ids                   = [10001]
      group_ids                  = []
      phone_order                = [10001]
      phone_circle_times         = 2
      phone_circle_interval      = 50
      phone_inner_interval       = 60
      need_phone_arrive_notice   = 1
      phone_call_type            = "CIRCLE"
      weekday                    =[1,2,3,4,5,6,7]
  }

  url_notices {
      url    = "https://www.mytest.com/validate"
      end_time =  0
      start_time = 1
      weekday = [1,2,3,4,5,6,7]
  }

}

resource "tencentcloud_pts_project" "project" {
  name = "ptsObjectName"
  description = "desc"
  tags {
    tag_key = "createdBy"
    tag_value = "terraform"
  }
}

resource "tencentcloud_pts_alert_channel" "alert_channel" {
  notice_id = tencentcloud_monitor_alarm_notice.example.id
  project_id = tencentcloud_pts_project.project.id
  amp_consumer_id = "Consumer-vvy1xxxxxx"
}

```
Import

pts alert_channel can be imported using the project_id#notice_id, e.g.
```
$ terraform import tencentcloud_pts_alert_channel.alert_channel project-kww5v8se#notice-kl66t6y9
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
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	pts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/pts/v20210728"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPtsAlertChannel() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudPtsAlertChannelRead,
		Create: resourceTencentCloudPtsAlertChannelCreate,
		Update: resourceTencentCloudPtsAlertChannelUpdate,
		Delete: resourceTencentCloudPtsAlertChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"notice_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Notice ID.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID.",
			},

			"amp_consumer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AMP Consumer ID.",
			},

			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Status Note: this field may return null, indicating that a valid value cannot be obtained.",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time Note: this field may return null, indicating that a valid value cannot be obtained.",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time Note: this field may return null, indicating that a valid value cannot be obtained.",
			},

			"app_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "App ID Note: this field may return null, indicating that a valid value cannot be obtained.",
			},

			"uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User ID Note: this field may return null, indicating that a valid value cannot be obtained.",
			},

			"sub_account_uin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sub-user ID Note: this field may return null, indicating that a valid value cannot be obtained.",
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

	if v, ok := d.GetOk("amp_consumer_id"); ok {
		request.AMPConsumerId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePtsClient().CreateAlertChannel(request)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation.DbRecordCreateFailed" {
					return resource.NonRetryableError(e)
				}
			}
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create pts alertChannel failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(projectId + FILED_SP + noticeId)
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
	projectId := idSplit[0]
	noticeId := idSplit[1]

	alertChannel, err := service.DescribePtsAlertChannel(ctx, noticeId, projectId)

	if err != nil {
		return err
	}

	if alertChannel == nil {
		d.SetId("")
		return fmt.Errorf("resource `alertChannel` %s does not exist", noticeId)
	}

	if alertChannel.NoticeId != nil {
		_ = d.Set("notice_id", alertChannel.NoticeId)
	}

	if alertChannel.ProjectId != nil {
		_ = d.Set("project_id", alertChannel.ProjectId)
	}

	if alertChannel.AMPConsumerId != nil {
		_ = d.Set("amp_consumer_id", alertChannel.AMPConsumerId)
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
	projectId := idSplit[0]
	noticeId := idSplit[1]

	if err := service.DeletePtsAlertChannelById(ctx, noticeId, projectId); err != nil {
		return err
	}

	return nil
}
