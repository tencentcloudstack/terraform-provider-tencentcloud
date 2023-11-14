/*
Provides a resource to create a live backup_stream

Example Usage

```hcl
resource "tencentcloud_live_backup_stream" "backup_stream" {
  push_domain_name = ""
  app_name = ""
  stream_name = ""
  upstream_sequence = ""
}
```

Import

live backup_stream can be imported using the id, e.g.

```
terraform import tencentcloud_live_backup_stream.backup_stream backup_stream_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudLiveBackupStream() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveBackupStreamCreate,
		Read:   resourceTencentCloudLiveBackupStreamRead,
		Update: resourceTencentCloudLiveBackupStreamUpdate,
		Delete: resourceTencentCloudLiveBackupStreamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"push_domain_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Push domain.",
			},

			"app_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "App name.",
			},

			"stream_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Stream id.",
			},

			"upstream_sequence": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sequence.",
			},
		},
	}
}

func resourceTencentCloudLiveBackupStreamCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_backup_stream.create")()
	defer inconsistentCheck(d, meta)()

	var pushDomainName string
	if v, ok := d.GetOk("push_domain_name"); ok {
		pushDomainName = v.(string)
	}

	var appName string
	if v, ok := d.GetOk("app_name"); ok {
		appName = v.(string)
	}

	var streamName string
	if v, ok := d.GetOk("stream_name"); ok {
		streamName = v.(string)
	}

	d.SetId(strings.Join([]string{pushDomainName, appName, streamName}, FILED_SP))

	return resourceTencentCloudLiveBackupStreamUpdate(d, meta)
}

func resourceTencentCloudLiveBackupStreamRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_backup_stream.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	pushDomainName := idSplit[0]
	appName := idSplit[1]
	streamName := idSplit[2]

	backupStream, err := service.DescribeLiveBackupStreamById(ctx, pushDomainName, appName, streamName)
	if err != nil {
		return err
	}

	if backupStream == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveBackupStream` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if backupStream.PushDomainName != nil {
		_ = d.Set("push_domain_name", backupStream.PushDomainName)
	}

	if backupStream.AppName != nil {
		_ = d.Set("app_name", backupStream.AppName)
	}

	if backupStream.StreamName != nil {
		_ = d.Set("stream_name", backupStream.StreamName)
	}

	if backupStream.UpstreamSequence != nil {
		_ = d.Set("upstream_sequence", backupStream.UpstreamSequence)
	}

	return nil
}

func resourceTencentCloudLiveBackupStreamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_backup_stream.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewSwitchBackupStreamRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	pushDomainName := idSplit[0]
	appName := idSplit[1]
	streamName := idSplit[2]

	request.PushDomainName = &pushDomainName
	request.AppName = &appName
	request.StreamName = &streamName

	immutableArgs := []string{"push_domain_name", "app_name", "stream_name", "upstream_sequence"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("push_domain_name") {
		if v, ok := d.GetOk("push_domain_name"); ok {
			request.PushDomainName = helper.String(v.(string))
		}
	}

	if d.HasChange("app_name") {
		if v, ok := d.GetOk("app_name"); ok {
			request.AppName = helper.String(v.(string))
		}
	}

	if d.HasChange("stream_name") {
		if v, ok := d.GetOk("stream_name"); ok {
			request.StreamName = helper.String(v.(string))
		}
	}

	if d.HasChange("upstream_sequence") {
		if v, ok := d.GetOk("upstream_sequence"); ok {
			request.UpstreamSequence = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().SwitchBackupStream(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live backupStream failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLiveBackupStreamRead(d, meta)
}

func resourceTencentCloudLiveBackupStreamDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_backup_stream.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
