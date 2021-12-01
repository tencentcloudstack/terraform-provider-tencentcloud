/*
Provide a resource to configure kubernetes cluster app addons.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.16.0.0/16"
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = "10.31.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "keep"
  cluster_desc            = "test cluster desc"
  cluster_version         = "1.20.6"
  cluster_max_service_num = 32

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"
}

resource "tencentcloud_kubernetes_addon_attachment" "addon_cbs" {
  cluster_id = "cls-xxxxxxxx"
  name = "cbs"
  version = "1.0.0"
}
```

Install new addon by passing spec json to `req_body` directly
```
resource "tencentcloud_kubernetes_addon_attachment" "addon_cbs" {
  cluster_id = "cls-xxxxxxxx"
  req_body = {\"spec\":{\"chart\":{\"chartName\":\"cbs\",\"chartVersion\":\"1.0.0\"},\"values\":{\"rawValuesType\":\"yaml\",\"values\":[]}}}
}
```

Import

Addon can be imported by using cluster_id#addon_name
```
$ terraform import tencentcloud_kubernetes_addon_attachment.addon_cos cls-xxxxxxxx#cos
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"strings"
)

func resourceTencentCloudTkeAddonAttachment() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of cluster.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of chart.",
			},
			"version": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Chart version, default latest version. Conflict with `request_body`.",
				ConflictsWith: []string{"request_body"},
			},
			"values": {
				Type:          schema.TypeList,
				Optional:      true,
				Description:   "Values the addon passthroughs. Conflict with `request_body`.",
				ConflictsWith: []string{"request_body"},
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"request_body": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Serialized json string as request body of addon spec. If set, will ignore `version` and `values`.",
				ConflictsWith: []string{"version", "values"},
			},
			"response_body": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Addon response body.",
			},
		},
		Create: resourceTencentCloudTkeAddonAttachmentCreate,
		Update: resourceTencentCloudTkeAddonAttachmentUpdate,
		Read:   resourceTencentCloudTkeAddonAttachmentRead,
		Delete: resourceTencentCloudTkeAddonAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceTencentCloudTkeAddonAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_addon_attachment.create")()
	logId := getLogId(contextNil)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client: client}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		clusterId = d.Get("cluster_id").(string)
		addonName = d.Get("name").(string)
		version   = d.Get("version").(string)
		values    = d.Get("values").([]interface{})
		reqBody   = d.Get("request_body").(string)
	)

	if reqBody == "" {
		var reqErr error
		v := helper.InterfacesStringsPoint(values)
		reqBody, reqErr = service.GetAddonPostReqBody(addonName, version, v)
		if reqErr != nil {
			return reqErr
		}
	}

	err := service.CreateExtensionAddon(ctx, clusterId, reqBody)

	if err != nil {
		return err
	}

	d.SetId(clusterId + FILED_SP + addonName)
	return resourceTencentCloudTkeAddonAttachmentRead(d, meta)
}

func resourceTencentCloudTkeAddonAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_addon_attachment.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	id := d.Id()
	split := strings.Split(id, FILED_SP)
	clusterId := split[0]
	addonName := split[1]

	response, err := service.DescribeExtensionAddon(ctx, clusterId, addonName)

	if err != nil {
		return err
	}

	_ = d.Set("response_body", response)

	addonResponseData := &AddonResponseData{}

	if err := json.Unmarshal([]byte(response), addonResponseData); err != nil {
		return err
	}

	spec := addonResponseData.Spec

	if spec != nil {
		_ = d.Set("cluster_id", clusterId)
		_ = d.Set("name", spec.Chart.ChartName)
		_ = d.Set("version", spec.Chart.ChartVersion)
		if spec.Values != nil && len(spec.Values.Values) > 0 {
			_ = d.Set("values", helper.StringsInterfaces(spec.Values.Values))
		}
	}

	d.SetId(id)

	return nil
}

func resourceTencentCloudTkeAddonAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_addon_attachment.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		id        = d.Id()
		split     = strings.Split(id, FILED_SP)
		clusterId = split[0]
		addonName = split[1]
		version   = d.Get("version").(string)
		values    = d.Get("values").([]interface{})
		reqBody   = d.Get("request_body").(string)
		err       error
	)

	if d.HasChange("request_body") && reqBody == "" || d.HasChange("version") || d.HasChange("values") {
		reqBody, err = service.GetAddonPatchReqBody(addonName, version, helper.InterfacesStringsPoint(values))
	}

	if err != nil {
		return err
	}

	err = service.UpdateExtensionAddon(ctx, clusterId, addonName, reqBody)

	if err != nil {
		return err
	}

	d.SetPartial("version")
	d.SetPartial("values")
	d.SetPartial("request_body")

	return resourceTencentCloudTkeAddonAttachmentRead(d, meta)
}

func resourceTencentCloudTkeAddonAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.resource_tc_kubernetes_addon_attachment.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		id        = d.Id()
		split     = strings.Split(id, FILED_SP)
		clusterId = split[0]
		addonName = split[1]
	)

	if err := service.DeleteExtensionAddon(ctx, clusterId, addonName); err != nil {
		return err
	}

	return nil
}
