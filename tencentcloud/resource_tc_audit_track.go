/*
Provides a resource to create a audit track

Example Usage

```hcl
resource "tencentcloud_audit_track" "track" {
  action_type           = "Read"
  event_names           = [
    "*",
  ]
  name                  = "terraform_track"
  resource_type         = "*"
  status                = 1
  track_for_all_members = 0

  storage {
    storage_name   = "db90b92c-91d2-46b0-94ac-debbbb21dc4e"
    storage_prefix = "cloudaudit"
    storage_region = "ap-guangzhou"
    storage_type   = "cls"
  }
}

```
Import

audit track can be imported using the id, e.g.
```
$ terraform import tencentcloud_audit_track.track track_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAuditTrack() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudAuditTrackRead,
		Create: resourceTencentCloudAuditTrackCreate,
		Update: resourceTencentCloudAuditTrackUpdate,
		Delete: resourceTencentCloudAuditTrackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Track name.",
			},

			"action_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Track interface type, optional:- `Read`: Read interface- `Write`: Write interface- `*`: All interface.",
			},

			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Track product, optional:- `*`: All product- Single product, such as `cos`.",
			},

			"status": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Track status, optional:- `0`: Close- `1`: Open.",
			},

			"event_names": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Track interface name list:- when ResourceType is `*`, EventNames is must `[&amp;quot;*&amp;quot;]`- when ResourceType is a single product, EventNames support all interface:`[&amp;quot;*&amp;quot;]`- when ResourceType is a single product, EventNames support some interface, up to 10.",
			},

			"storage": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
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
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to enable the delivery of group member operation logs to the group management account or trusted service management account, optional:- `0`: Close- `1`: Open.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Track create time.",
			},
		},
	}
}

func resourceTencentCloudAuditTrackCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_audit_track.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = audit.NewCreateAuditTrackRequest()
		response *audit.CreateAuditTrackResponse
		trackId  uint64
	)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("action_type"); ok {
		request.ActionType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_type"); ok {
		request.ResourceType = helper.String(v.(string))
	}

	if v, _ := d.GetOk("status"); v != nil {
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
		storage := audit.Storage{}
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

	if v, _ := d.GetOk("track_for_all_members"); v != nil {
		request.TrackForAllMembers = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAuditClient().CreateAuditTrack(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create audit track failed, reason:%+v", logId, err)
		return err
	}

	trackId = *response.Response.TrackId

	d.SetId(helper.UInt64ToStr(trackId))
	return resourceTencentCloudAuditTrackRead(d, meta)
}

func resourceTencentCloudAuditTrackRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_audit_track.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AuditService{client: meta.(*TencentCloudClient).apiV3Conn}

	trackId := d.Id()

	track, err := service.DescribeAuditTrackById(ctx, trackId)

	if err != nil {
		return err
	}

	if track == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", trackId)
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

func resourceTencentCloudAuditTrackUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_audit_track.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := audit.NewModifyAuditTrackRequest()

	trackId := d.Id()

	request.TrackId = helper.StrToUint64Point(trackId)

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
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
		if v, _ := d.GetOk("status"); v != nil {
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
			storage := audit.Storage{}
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
		if v, _ := d.GetOk("track_for_all_members"); v != nil {
			request.TrackForAllMembers = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAuditClient().ModifyAuditTrack(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create audit track failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudAuditTrackRead(d, meta)
}

func resourceTencentCloudAuditTrackDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_audit_track.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := AuditService{client: meta.(*TencentCloudClient).apiV3Conn}

	trackId := d.Id()

	if err := service.DeleteAuditTrackById(ctx, trackId); err != nil {
		return err
	}

	return nil
}
