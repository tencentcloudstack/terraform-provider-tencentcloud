/*
Provides a resource to create a monitor tmp_manage_grafana_attachment

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}


resource "tencentcloud_monitor_tmp_instance" "foo" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  data_retention_time = 30
  zone                = var.availability_zone
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_monitor_grafana_instance" "foo" {
  instance_name         = "tf-grafana"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_ids            = [tencentcloud_subnet.subnet.id]
  grafana_init_password = "1234567890"
  enable_internet 		= false
  is_destroy 			= true

  tags = {
    "createdBy" = "test"
  }
}

resource "tencentcloud_monitor_tmp_manage_grafana_attachment" "foo" {
    grafana_id  = tencentcloud_monitor_grafana_instance.foo.id
    instance_id = tencentcloud_monitor_tmp_instance.foo.id
}
```

Import

monitor tmp_manage_grafana_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_manage_grafana_attachment.manage_grafana_attachment prom-xxxxxxxx
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpManageGrafanaAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpManageGrafanaAttachmentCreate,
		Read:   resourceTencentCloudMonitorTmpManageGrafanaAttachmentRead,
		Delete: resourceTencentCloudMonitorTmpManageGrafanaAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Prometheus instance ID.",
			},

			"grafana_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance ID.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpManageGrafanaAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_manage_grafana_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = monitor.NewBindPrometheusManagedGrafanaRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("grafana_id"); ok {
		request.GrafanaId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().BindPrometheusManagedGrafana(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create monitor manageGrafanaAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorTmpManageGrafanaAttachmentRead(d, meta)
}

func resourceTencentCloudMonitorTmpManageGrafanaAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_manage_grafana_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	manageGrafanaAttachment, err := service.DescribeMonitorManageGrafanaAttachmentById(ctx, instanceId)
	if err != nil {
		return err
	}

	if manageGrafanaAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorTmpManageGrafanaAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if manageGrafanaAttachment.InstanceId != nil {
		_ = d.Set("instance_id", manageGrafanaAttachment.InstanceId)
	}

	if manageGrafanaAttachment.GrafanaInstanceId != nil {
		_ = d.Set("grafana_id", manageGrafanaAttachment.GrafanaInstanceId)
	}

	return nil
}

func resourceTencentCloudMonitorTmpManageGrafanaAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_manage_grafana_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()

	if err := service.DeleteMonitorManageGrafanaAttachmentById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
