/*
Provides a resource to create a cloudaudit track

Example Usage

```hcl
resource "tencentcloud_cloudaudit_track" "track" {
  track_id = &lt;nil&gt;
  name = &lt;nil&gt;
  action_type = &lt;nil&gt;
  resource_type = &lt;nil&gt;
  status = &lt;nil&gt;
  event_names = &lt;nil&gt;
  storage {
		storage_type = &lt;nil&gt;
		storage_region = &lt;nil&gt;
		storage_name = &lt;nil&gt;
		storage_prefix = &lt;nil&gt;

  }
  track_for_all_members = &lt;nil&gt;
  }
```

Import

cloudaudit track can be imported using the id, e.g.

```
terraform import tencentcloud_cloudaudit_track.track track_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cloudaudit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCloudauditTrack() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCloudauditTrackCreate,
		Read:   resourceTencentCloudCloudauditTrackRead,
		Update: resourceTencentCloudCloudauditTrackUpdate,
		Delete: resourceTencentCloudCloudauditTrackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"track_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Track Id.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Track name.",
			},

			"action_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Track interface type, optional:- `Read`: Read interface- `Write`: Write interface- `*`: All interface.",
			},

			"resource_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Track product, optional:- `*`: All product- Single product, such as `cos`.",
			},

			"status": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Track status, optional:- `0`: Close- `1`: Open.",
			},

			"event_names": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Track interface name list:- when ResourceType is `*`, EventNames is must `[*]`- when ResourceType is a single product, EventNames support all interface:`[*]`- when ResourceType is a single product, EventNames support some interface, up to 10.",
			},

			"storage": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Track Storage, support `cos` and `cls`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Track Storage type, optional:- `cos`- `cls`.",
						},
						"storage_region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Storage region.",
						},
						"storage_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Track Storage name:- when StorageType is `cls`, StorageName is cls topicId- when StorageType is `cos`, StorageName is cos bucket name that does not contain `-APPID`.",
						},
						"storage_prefix": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Storage path prefix.",
						},
					},
				},
			},

			"track_for_all_members": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to enable the delivery of group member operation logs to the group management account or trusted service management account, optional:- `0`: Close- `1`: Open.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Track create time.",
			},
		},
	}
}

func resourceTencentCloudCloudauditTrackCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cloudaudit_track.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cloudaudit.NewCreateAuditTrackRequest()
		response = cloudaudit.NewCreateAuditTrackResponse()
		trackId  int
	)
	if v, ok := d.GetOkExists("track_id"); ok {
		trackId = v.(int)
		request.TrackId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("action_type"); ok {
		request.ActionType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_type"); ok {
		request.ResourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("status"); ok {
		request.Status = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("event_names"); ok {
		eventNamesSet := v.(*schema.Set).List()
		for i := range eventNamesSet {
			eventNames := eventNamesSet[i].(string)
			request.EventNames = append(request.EventNames, &eventNames)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "storage"); ok {
		storage := cloudaudit.Storage{}
		if v, ok := dMap["storage_type"]; ok {
			storage.StorageType = helper.String(v.(string))
		}
		if v, ok := dMap["storage_region"]; ok {
			storage.StorageRegion = helper.String(v.(string))
		}
		if v, ok := dMap["storage_name"]; ok {
			storage.StorageName = helper.String(v.(string))
		}
		if v, ok := dMap["storage_prefix"]; ok {
			storage.StoragePrefix = helper.String(v.(string))
		}
		request.Storage = &storage
	}

	if v, ok := d.GetOkExists("track_for_all_members"); ok {
		request.TrackForAllMembers = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCloudauditClient().CreateAuditTrack(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cloudaudit track failed, reason:%+v", logId, err)
		return err
	}

	trackId = *response.Response.TrackId
	d.SetId(helper.Int64ToStr(int64(trackId)))

	return resourceTencentCloudCloudauditTrackRead(d, meta)
}

func resourceTencentCloudCloudauditTrackRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cloudaudit_track.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CloudauditService{client: meta.(*TencentCloudClient).apiV3Conn}

	trackId := d.Id()

	track, err := service.DescribeCloudauditTrackById(ctx, trackId)
	if err != nil {
		return err
	}

	if track == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CloudauditTrack` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if track.TrackId != nil {
		_ = d.Set("track_id", track.TrackId)
	}

	if track.Name != nil {
		_ = d.Set("name", track.Name)
	}

	if track.ActionType != nil {
		_ = d.Set("action_type", track.ActionType)
	}

	if track.ResourceType != nil {
		_ = d.Set("resource_type", track.ResourceType)
	}

	if track.Status != nil {
		_ = d.Set("status", track.Status)
	}

	if track.EventNames != nil {
		_ = d.Set("event_names", track.EventNames)
	}

	if track.Storage != nil {
		storageMap := map[string]interface{}{}

		if track.Storage.StorageType != nil {
			storageMap["storage_type"] = track.Storage.StorageType
		}

		if track.Storage.StorageRegion != nil {
			storageMap["storage_region"] = track.Storage.StorageRegion
		}

		if track.Storage.StorageName != nil {
			storageMap["storage_name"] = track.Storage.StorageName
		}

		if track.Storage.StoragePrefix != nil {
			storageMap["storage_prefix"] = track.Storage.StoragePrefix
		}

		_ = d.Set("storage", []interface{}{storageMap})
	}

	if track.TrackForAllMembers != nil {
		_ = d.Set("track_for_all_members", track.TrackForAllMembers)
	}

	if track.CreateTime != nil {
		_ = d.Set("create_time", track.CreateTime)
	}

	return nil
}

func resourceTencentCloudCloudauditTrackUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cloudaudit_track.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cloudaudit.NewModifyAuditTrackRequest()

	trackId := d.Id()

	request.TrackId = &trackId

	immutableArgs := []string{"track_id", "name", "action_type", "resource_type", "status", "event_names", "storage", "track_for_all_members", "create_time"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("action_type") {
		if v, ok := d.GetOk("action_type"); ok {
			request.ActionType = helper.String(v.(string))
		}
	}

	if d.HasChange("resource_type") {
		if v, ok := d.GetOk("resource_type"); ok {
			request.ResourceType = helper.String(v.(string))
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOkExists("status"); ok {
			request.Status = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("event_names") {
		if v, ok := d.GetOk("event_names"); ok {
			eventNamesSet := v.(*schema.Set).List()
			for i := range eventNamesSet {
				eventNames := eventNamesSet[i].(string)
				request.EventNames = append(request.EventNames, &eventNames)
			}
		}
	}

	if d.HasChange("storage") {
		if dMap, ok := helper.InterfacesHeadMap(d, "storage"); ok {
			storage := cloudaudit.Storage{}
			if v, ok := dMap["storage_type"]; ok {
				storage.StorageType = helper.String(v.(string))
			}
			if v, ok := dMap["storage_region"]; ok {
				storage.StorageRegion = helper.String(v.(string))
			}
			if v, ok := dMap["storage_name"]; ok {
				storage.StorageName = helper.String(v.(string))
			}
			if v, ok := dMap["storage_prefix"]; ok {
				storage.StoragePrefix = helper.String(v.(string))
			}
			request.Storage = &storage
		}
	}

	if d.HasChange("track_for_all_members") {
		if v, ok := d.GetOkExists("track_for_all_members"); ok {
			request.TrackForAllMembers = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCloudauditClient().ModifyAuditTrack(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cloudaudit track failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCloudauditTrackRead(d, meta)
}

func resourceTencentCloudCloudauditTrackDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cloudaudit_track.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CloudauditService{client: meta.(*TencentCloudClient).apiV3Conn}
	trackId := d.Id()

	if err := service.DeleteCloudauditTrackById(ctx, trackId); err != nil {
		return err
	}

	return nil
}
