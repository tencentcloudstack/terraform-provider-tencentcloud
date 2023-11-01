/*
Provides a resource to create a css backup_stream

~> **NOTE:** This resource is only valid when the push stream. When the push stream ends, it will be deleted.

Example Usage

```hcl
resource "tencentcloud_css_backup_stream" "backup_stream" {
  push_domain_name  = "177154.push.tlivecloud.com"
  app_name          = "live"
  stream_name       = "1308919341_test"
  upstream_sequence = "2209501773993286139"
}
```

Import

css backup_stream can be imported using the id, e.g.

```
terraform import tencentcloud_css_backup_stream.backup_stream pushDomainName#appName#streamName
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
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssBackupStream() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssBackupStreamCreate,
		Read:   resourceTencentCloudCssBackupStreamRead,
		Update: resourceTencentCloudCssBackupStreamUpdate,
		Delete: resourceTencentCloudCssBackupStreamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"push_domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Push domain.",
			},

			"app_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "App name.",
			},

			"stream_name": {
				Required:    true,
				ForceNew:    true,
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

func resourceTencentCloudCssBackupStreamCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_backup_stream.create")()
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

	return resourceTencentCloudCssBackupStreamUpdate(d, meta)
}

func resourceTencentCloudCssBackupStreamRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_backup_stream.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	pushDomainName := idSplit[0]
	appName := idSplit[1]
	streamName := idSplit[2]

	backupStream, err := service.DescribeCssBackupStreamById(ctx, pushDomainName, appName, streamName)
	if err != nil {
		return err
	}

	if backupStream == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssBackupStream` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("push_domain_name", pushDomainName)
	_ = d.Set("app_name", appName)
	_ = d.Set("stream_name", streamName)

	if backupStream.UpstreamSequence != nil {
		_ = d.Set("upstream_sequence", backupStream.UpstreamSequence)
	}

	return nil
}

func resourceTencentCloudCssBackupStreamUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_backup_stream.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := css.NewSwitchBackupStreamRequest()

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

	if v, ok := d.GetOk("upstream_sequence"); ok {
		request.UpstreamSequence = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().SwitchBackupStream(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update css backupStream failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssBackupStreamRead(d, meta)
}

func resourceTencentCloudCssBackupStreamDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_backup_stream.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
