/*
Provides a resource to create a dts migrate_check_config

Example Usage

```hcl
resource "tencentcloud_dts_migrate_check_config" "migrate_check_config" {
  job_id = ""
}
```

*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
)

func resourceTencentCloudDtsMigrateCheckConfig() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDtsMigrateCheckConfigRead,
		Update: resourceTencentCloudDtsMigrateCheckConfigUpdate,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job Id.",
			},
		},
	}
}

func resourceTencentCloudDtsMigrateCheckConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_check_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	jobId := d.Id()

	migrateCheckConfig, err := service.DescribeDtsMigrateCheckById(ctx, jobId)
	if err != nil {
		return err
	}

	if migrateCheckConfig == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("job_id", jobId)

	return nil
}

func resourceTencentCloudDtsMigrateCheckConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dts_migrate_check_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dts.NewCreateMigrateCheckJobRequest()

	jobId := d.Id()

	request.JobId = &jobId
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDtsClient().CreateMigrateCheckJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create dts migrateCheckConfig failed, reason:%+v", logId, err)
		return err
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"checkPass", "checkNotPass"}, 0*readRetryTimeout, time.Second, service.DtsMigrateCheckConfigStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDtsMigrateCheckConfigRead(d, meta)
}
