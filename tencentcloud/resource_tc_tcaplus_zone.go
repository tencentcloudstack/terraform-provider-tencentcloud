/*
/*
Use this resource to create tcaplus zone

Example Usage

```hcl
resource "tencentcloud_tcaplus_application" "test" {
  idl_type                 = "PROTO"
  app_name                 = "tf_tcaplus_app_test"
  vpc_id                   = "vpc-7k6gzox6"
  subnet_id                = "subnet-akwgvfa3"
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_zone" "zone" {
  app_id = tencentcloud_tcaplus_application.test.id
  zone_name      = "tf_test_zone_name"
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudTcaplusZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcaplusZoneCreate,
		Read:   resourceTencentCloudTcaplusZoneRead,
		Update: resourceTencentCloudTcaplusZoneUpdate,
		Delete: resourceTencentCloudTcaplusZoneDelete,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Application of the tcapplus zone belongs.",
			},
			"zone_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 30),
				Description:  "Name of the tcapplus zone. length should between 1 and 30.",
			},
			// Computed values.
			"table_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of tables.",
			},
			"total_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total storage(MB).",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the tcapplus zone.",
			},
		},
	}
}

func resourceTencentCloudTcaplusZoneCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_zone.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		applicationId = d.Get("app_id").(string)
		zoneName      = d.Get("zone_name").(string)
	)
	zoneId, err := tcaplusService.CreateZone(ctx, applicationId, zoneName)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s:%s", applicationId, zoneId))
	return resourceTencentCloudTcaplusZoneRead(d, meta)
}

func resourceTencentCloudTcaplusZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_zone.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	applicationId := d.Get("app_id").(string)
	zoneId := d.Id()

	info, has, err := tcaplusService.DescribeZone(ctx, applicationId, zoneId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			info, has, err = tcaplusService.DescribeZone(ctx, applicationId, zoneId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_name", info.ZoneName)
	_ = d.Set("table_count", int(*info.TableCount))
	_ = d.Set("total_size", int(*info.TotalSize))
	_ = d.Set("create_time", info.CreatedTime)

	return nil
}

func resourceTencentCloudTcaplusZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_zone.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	applicationId := d.Get("app_id").(string)
	zoneId := d.Id()

	if d.HasChange("zone_name") {
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			err := tcaplusService.ModifyZoneName(ctx, applicationId, zoneId, d.Get("zone_name").(string))
			if err != nil {
				return retryError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return resourceTencentCloudTcaplusZoneRead(d, meta)
}

func resourceTencentCloudTcaplusZoneDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_zone.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	applicationId := d.Get("app_id").(string)
	zoneId := d.Id()

	err := tcaplusService.DeleteZone(ctx, applicationId, zoneId)
	if err != nil {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			err = tcaplusService.DeleteZone(ctx, applicationId, zoneId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	_, has, err := tcaplusService.DescribeZone(ctx, applicationId, zoneId)
	if err != nil || has {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = tcaplusService.DescribeZone(ctx, applicationId, zoneId)
			if err != nil {
				return retryError(err)
			}
			if has {
				err = fmt.Errorf("delete zone fail, zone still exist from sdk DescribeZones")
				return resource.RetryableError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	if !has {
		return nil
	} else {
		return errors.New("delete zone fail, zone still exist from sdk DescribeZones")
	}
}
