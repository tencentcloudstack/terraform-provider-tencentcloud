package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDnspodSnapshotConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDnspodSnapshotConfigCreate,
		Read:   resourceTencentCloudDnspodSnapshotConfigRead,
		Update: resourceTencentCloudDnspodSnapshotConfigUpdate,
		Delete: resourceTencentCloudDnspodSnapshotConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain name.",
			},

			"period": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Backup interval: empty string - no backup, half_hour - every half hour, hourly - every hour, daily - every day, monthly - every month.",
			},
		},
	}
}

func resourceTencentCloudDnspodSnapshotConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_snapshot_config.create")()
	defer inconsistentCheck(d, meta)()

	var domain string

	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
	}

	d.SetId(domain)

	return resourceTencentCloudDnspodSnapshotConfigUpdate(d, meta)
}

func resourceTencentCloudDnspodSnapshotConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_snapshot_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}

	domain := d.Id()

	snapshotConfig, err := service.DescribeDnspodSnapshotConfigById(ctx, domain)
	if err != nil {
		return err
	}

	if snapshotConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DnspodSnapshot_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("domain", domain)

	if snapshotConfig.Config != nil {
		_ = d.Set("period", snapshotConfig.Config)
	}

	return nil
}

func resourceTencentCloudDnspodSnapshotConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_snapshot_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dnspod.NewModifySnapshotConfigRequest()

	domain := d.Id()
	request.Domain = &domain

	if v, ok := d.GetOk("period"); ok {
		request.Period = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDnsPodClient().ModifySnapshotConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dnspod snapshot_config failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDnspodSnapshotConfigRead(d, meta)
}

func resourceTencentCloudDnspodSnapshotConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dnspod_snapshot_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
