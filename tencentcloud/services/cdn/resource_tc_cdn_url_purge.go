package cdn

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudUrlPurge() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudUrlPurgeRead,
		Create: resourceTencentCloudUrlPurgeCreate,
		Update: resourceTencentCloudUrlPurgeUpdate,
		Delete: resourceTencentCloudUrlPurgeDelete,
		Schema: map[string]*schema.Schema{
			"urls": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of url to purge. NOTE: urls need include protocol prefix `http://` or `https://`.",
			},
			"redo": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Change to purge again. NOTE: this argument only works while resource update, if set to `0` or null will not be triggered.",
			},
			"area": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specify purge area. NOTE: only purge same area cache contents.",
			},
			"url_encode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to encode urls, if set to `true` will auto encode instead of manual process.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Task id of last operation.",
			},
			"purge_history": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "logs of latest purge task.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purge task id.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purge url.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purge status of `fail`, `done`, `process`.",
						},
						"purge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purge category in of `url` or `path`.",
						},
						"flush_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purge flush type of `flush` or `delete`.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Purge task create time.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudUrlPurgeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_url_purge.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := CdnService{client}

	taskId, ok := d.Get("task_id").(string)

	if !ok || taskId == "" {
		return fmt.Errorf("no task id provided")
	}

	request := cdn.NewDescribePurgeTasksRequest()
	request.TaskId = &taskId

	var (
		logs []*cdn.PurgeTask
		err  error
	)

	err = resource.Retry(tccommon.ReadRetryTimeout*2, func() *resource.RetryError {
		logs, err = service.DescribePurgeTasks(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
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
				return resource.NonRetryableError(fmt.Errorf("purge url %s failed", *item.Url))
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
	d.SetId("purges-" + GetUrlsHash(helper.InterfacesStrings(urls)))

	purgeHistory := make([]interface{}, 0)

	if len(logs) > 0 {
		for i := range logs {
			item := logs[i]
			purge := map[string]interface{}{
				"task_id":     item.TaskId,
				"url":         item.Url,
				"create_time": item.CreateTime,
				"status":      item.Status,
				"purge_type":  item.PurgeType,
				"flush_type":  item.FlushType,
			}
			purgeHistory = append(purgeHistory, purge)
		}
	}

	_ = d.Set("purge_history", purgeHistory)

	return nil
}

func resourceTencentCloudUrlPurgeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_url_purge.create")()

	taskId, err := tencentcloudCdnUrlPurge(d, meta)

	if err != nil {
		return err
	}

	_ = d.Set("task_id", taskId)

	urls := d.Get("urls").([]interface{})
	d.SetId("purges-" + GetUrlsHash(helper.InterfacesStrings(urls)))

	return resourceTencentCloudUrlPurgeRead(d, meta)
}

func resourceTencentCloudUrlPurgeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_url_purge.update")()

	redo, ok := d.GetOk("redo")

	if !d.HasChange("redo") || !ok || redo.(int) == 0 {
		return nil
	}

	taskId, err := tencentcloudCdnUrlPurge(d, meta)

	if err != nil {
		return err
	}

	_ = d.Set("task_id", taskId)

	return resourceTencentCloudUrlPurgeRead(d, meta)
}

func resourceTencentCloudUrlPurgeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cdn_url_purge.delete")()
	log.Printf("noop deleting resoruce %s", "tencentcloud_cdn_url_purge")
	return nil
}

func tencentcloudCdnUrlPurge(d *schema.ResourceData, meta interface{}) (string, error) {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	service := CdnService{client}

	urls := d.Get("urls").([]interface{})
	request := cdn.NewPurgeUrlsCacheRequest()
	request.Urls = helper.InterfacesStringsPoint(urls)

	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}

	if v, ok := d.GetOk("url_encode"); ok {
		request.UrlEncode = helper.Bool(v.(bool))
	}

	return service.PurgeUrlsCache(ctx, request)
}
