/*
Provide a resource to invoke a Url Push request.

Example Usage

```hcl
resource "tencentcloud_cdn_url_push" "foo" {
  urls = ["https://www.example.com/b"]
}
```

Change `redo` argument to request new push task with same urls

```hcl
resource "tencentcloud_cdn_url_push" "foo" {
  urls = [
    "https://www.example.com/a"
  ]
  redo = 1
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudUrlPush() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudUrlPushRead,
		Create: resourceTencentCloudUrlPushCreate,
		Update: resourceTencentCloudUrlPushUpdate,
		Delete: resourceTencentCloudUrlPushDelete,
		Schema: map[string]*schema.Schema{
			"urls": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of url to push. NOTE: urls need include protocol prefix `http://` or `https://`.",
			},
			"redo": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Change to push again. NOTE: this argument only works while resource update, if set to `0` or null will not be triggered.",
			},
			"area": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify push area. NOTE: only push same area cache contents.",
			},
			"user_agent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify `User-Agent` HTTP header, default: `TencentCdn`.",
			},
			"layer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Layer to push.",
			},
			"parse_m3u8": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to recursive parse m3u8 files.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Push task id.",
			},
			// Plan to support
			//"disable_range": {
			//	Type:        schema.TypeBool,
			//	Optional:    true,
			//	Description: "Whether to disable range origin pull.",
			//},
			"push_history": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "logs of latest push task.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Push task id.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Push url.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Push status of `fail`, `done`, `process` or `invalid` (4xx, 5xx response).",
						},
						"percent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Push progress in percent.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Push task create time.",
						},
						"area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Push tag area in `mainland`, `overseas` or `global`.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Push task update time.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudUrlPushRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdn_url_push.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := CdnService{client}

	taskId, ok := d.Get("task_id").(string)

	if !ok || taskId == "" {
		return fmt.Errorf("no task id provided")
	}

	request := cdn.NewDescribePushTasksRequest()
	request.TaskId = &taskId

	var (
		logs []*cdn.PushTask
		err  error
	)

	err = resource.Retry(readRetryTimeout*2, func() *resource.RetryError {
		logs, err = service.DescribePushTasks(ctx, request)
		if err != nil {
			return retryError(err)
		}
		if len(logs) == 0 {
			return resource.RetryableError(fmt.Errorf("task %s returns nil logs, retrying", taskId))
		}
		for i := range logs {
			item := logs[i]
			status := item.Status
			if status == nil {
				continue
			}
			switch *status {
			case "process":
				return resource.RetryableError(fmt.Errorf("processing %s", *item.Url))
			case "fail":
				return resource.NonRetryableError(fmt.Errorf("push url %s failed", *item.Url))
			default:
				continue
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	urls := d.Get("urls").([]interface{})
	d.SetId("pushes-" + GetUrlsHash(helper.InterfacesStrings(urls)))

	pushHistory := make([]interface{}, 0)

	if len(logs) > 0 {
		for i := range logs {
			item := logs[i]
			pushHistory = append(pushHistory, map[string]interface{}{
				"task_id":     item.TaskId,
				"url":         item.Url,
				"status":      item.Status,
				"percent":     item.Percent,
				"create_time": item.CreateTime,
				"area":        item.Area,
				"update_time": item.UpdateTime,
			})
		}
	}

	if err := d.Set("push_history", pushHistory); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudUrlPushCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdn_url_push.create")()

	taskId, err := tencentcloudCdnUrlPush(d, meta)

	if err != nil {
		return err
	}

	_ = d.Set("task_id", taskId)

	urls := d.Get("urls").([]interface{})
	d.SetId("pushes-" + GetUrlsHash(helper.InterfacesStrings(urls)))

	return resourceTencentCloudUrlPushRead(d, meta)
}

func resourceTencentCloudUrlPushUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdn_url_push.update")()

	redo, ok := d.GetOk("redo")

	if !d.HasChange("redo") || !ok || redo.(int) == 0 {
		return nil
	}

	taskId, err := tencentcloudCdnUrlPush(d, meta)

	if err != nil {
		return err
	}

	_ = d.Set("task_id", taskId)

	return resourceTencentCloudUrlPushRead(d, meta)
}

func resourceTencentCloudUrlPushDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdn_url_push.delete")()

	log.Printf("noop deleting resoruce %s", "tencentcloud_cdn_url_push")
	return nil
}

func tencentcloudCdnUrlPush(d *schema.ResourceData, meta interface{}) (string, error) {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	service := CdnService{client}

	urls := d.Get("urls").([]interface{})
	request := cdn.NewPushUrlsCacheRequest()
	request.Urls = helper.InterfacesStringsPoint(urls)

	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}

	if v, ok := d.GetOk("layer"); ok {
		request.Layer = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_agent"); ok {
		request.UserAgent = helper.String(v.(string))
	}

	if v, ok := d.GetOk("parse_m3u8"); ok {
		request.ParseM3U8 = helper.Bool(v.(bool))
	}

	return service.PushUrlsCache(ctx, request)
}
