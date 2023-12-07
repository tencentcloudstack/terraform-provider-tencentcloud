package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudLighthouseSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseSnapshotCreate,
		Read:   resourceTencentCloudLighthouseSnapshotRead,
		Update: resourceTencentCloudLighthouseSnapshotUpdate,
		Delete: resourceTencentCloudLighthouseSnapshotDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the instance for which to create a snapshot.",
			},

			"snapshot_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Snapshot name, which can contain up to 60 characters.",
			},
		},
	}
}

func resourceTencentCloudLighthouseSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_snapshot.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = lighthouse.NewCreateInstanceSnapshotRequest()
		response   = lighthouse.NewCreateInstanceSnapshotResponse()
		snapshotId string
	)
	request.InstanceId = helper.String(d.Get("instance_id").(string))

	if v, ok := d.GetOk("snapshot_name"); ok {
		request.SnapshotName = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().CreateInstanceSnapshot(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse snapshot failed, reason:%+v", logId, err)
		return err
	}

	snapshotId = *response.Response.SnapshotId
	d.SetId(snapshotId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"NORMAL"}, 20*readRetryTimeout, time.Second, service.LighthouseSnapshotStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseSnapshotRead(d, meta)
}

func resourceTencentCloudLighthouseSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_snapshot.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	snapshotId := d.Id()

	snapshot, err := service.DescribeLighthouseSnapshotById(ctx, snapshotId)
	if err != nil {
		return err
	}

	if snapshot == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseSnapshot` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if snapshot.SnapshotName != nil {
		_ = d.Set("snapshot_name", snapshot.SnapshotName)
	}

	return nil
}

func resourceTencentCloudLighthouseSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_snapshot.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := lighthouse.NewModifySnapshotAttributeRequest()

	snapshotId := d.Id()

	request.SnapshotId = &snapshotId

	if d.HasChange("snapshot_name") {
		if v, ok := d.GetOk("snapshot_name"); ok {
			request.SnapshotName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().ModifySnapshotAttribute(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update lighthouse snapshot failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLighthouseSnapshotRead(d, meta)
}

func resourceTencentCloudLighthouseSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_snapshot.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}
	snapshotId := d.Id()

	if err := service.DeleteLighthouseSnapshotById(ctx, snapshotId); err != nil {
		return err
	}

	return nil
}
