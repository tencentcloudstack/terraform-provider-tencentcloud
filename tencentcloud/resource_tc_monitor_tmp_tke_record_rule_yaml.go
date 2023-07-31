/*
Provides a resource to create a tke tmpRecordRule

Example Usage

```hcl
variable "default_instance_type" {
  default = "SA1.MEDIUM2"
}

variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "availability_zone_second" {
  default = "ap-guangzhou-4"
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

locals {
  first_vpc_id     = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.vpc_id
  first_subnet_id  = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.subnet_id
  second_vpc_id    = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.vpc_id
  second_subnet_id = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.subnet_id
  sg_id            = tencentcloud_security_group.sg.id
  image_id         = data.tencentcloud_images.default.image_id
}

data "tencentcloud_vpc_subnets" "vpc_one" {
  is_default        = true
  availability_zone = var.availability_zone_first
}

data "tencentcloud_vpc_subnets" "vpc_two" {
  is_default        = true
  availability_zone = var.availability_zone_second
}

resource "tencentcloud_security_group" "sg" {
  name = "tf-example-sg"
}

resource "tencentcloud_security_group_lite_rule" "sg_rule" {
  security_group_id = tencentcloud_security_group.sg.id

  ingress = [
    "ACCEPT#10.0.0.0/16#ALL#ALL",
    "ACCEPT#172.16.0.0/22#ALL#ALL",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]

  egress = [
    "ACCEPT#172.16.0.0/22#ALL#ALL",
  ]
}

data "tencentcloud_images" "default" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                          = local.first_vpc_id
  cluster_cidr                    = var.example_cluster_cidr
  cluster_max_pod_num             = 32
  cluster_name                    = "tf_example_cluster"
  cluster_desc                    = "example for tke cluster"
  cluster_max_service_num         = 32
  cluster_internet                = false
  cluster_internet_security_group = local.sg_id
  cluster_version                 = "1.22.5"
  cluster_deploy_type             = "MANAGED_CLUSTER"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_first
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.first_subnet_id
    img_id                     = local.image_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # key_ids                   = ["skey-11112222"]
    password = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_second
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.second_subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # key_ids                   = ["skey-11112222"]
    cam_role_name = "CVM_QcsRole"
    password      = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}


# create monitor
variable "zone" {
  default = "ap-guangzhou"
}

variable "cluster_type" {
  default = "tke"
}

resource "tencentcloud_monitor_tmp_instance" "foo" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = local.first_vpc_id
  subnet_id           = local.first_subnet_id
  data_retention_time = 30
  zone                = var.availability_zone_second
  tags = {
    "createdBy" = "terraform"
  }
}


# tmp tke bind
resource "tencentcloud_monitor_tmp_tke_cluster_agent" "foo" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id

  agents {
    region          = var.zone
    cluster_type    = var.cluster_type
    cluster_id      = tencentcloud_kubernetes_cluster.example.id
    enable_external = false
  }
}

# create record rule
resource "tencentcloud_monitor_tmp_tke_record_rule_yaml" "foo" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id
  content     = <<-EOT
        apiVersion: monitoring.coreos.com/v1
        kind: PrometheusRule
        metadata:
          name: example-record
        spec:
          groups:
            - name: kube-apiserver.rules
              rules:
                - expr: sum(metrics_test)
                  labels:
                    verb: read
                  record: 'apiserver_request:burnrate1d'
    EOT

  depends_on = [tencentcloud_monitor_tmp_tke_cluster_agent.foo]
}
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
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpTkeRecordRuleYaml() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTkeTmpRecordRuleYamlRead,
		Create: resourceTencentCloudTkeTmpRecordRuleYamlCreate,
		Update: resourceTencentCloudTkeTmpRecordRuleYamlUpdate,
		Delete: resourceTencentCloudTkeTmpRecordRuleYamlDelete,
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance Id.",
			},

			"content": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateYaml,
				Description:  "Contents of record rules in yaml format.",
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the instance.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time of record rule.",
			},

			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Used for the argument, if the configuration comes to the template, the template id.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An ID identify the cluster, like cls-xxxxxx.",
			},
		},
	}
}

func resourceTencentCloudTkeTmpRecordRuleYamlCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_record_rule_yaml.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewCreatePrometheusRecordRuleYamlRequest()

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	tmpRecordRuleName := ""
	if v, ok := d.GetOk("content"); ok {
		m, _ := YamlParser(v.(string))
		metadata := m["metadata"]
		if metadata != nil {
			if metadata.(map[interface{}]interface{})["name"] != nil {
				tmpRecordRuleName = metadata.(map[interface{}]interface{})["name"].(string)
			}
		}
		request.Content = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreatePrometheusRecordRuleYaml(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tke tmpRecordRule failed, reason:%+v", logId, err)
		return err
	}

	instanceId := *request.InstanceId
	d.SetId(strings.Join([]string{instanceId, tmpRecordRuleName}, FILED_SP))
	return resourceTencentCloudTkeTmpRecordRuleYamlRead(d, meta)
}

func resourceTencentCloudTkeTmpRecordRuleYamlRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_record_rule_yaml.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	instanceId := ids[0]
	name := ids[1]

	recordRuleService := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	request, err := recordRuleService.DescribePrometheusRecordRuleByName(ctx, instanceId, name)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			request, err = recordRuleService.DescribePrometheusRecordRuleByName(ctx, instanceId, name)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	recordRules := request.Response.Records
	if len(recordRules) == 0 {
		d.SetId("")
		return nil
	}

	recordRule := recordRules[0]
	if recordRule == nil {
		return nil
	}
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("name", recordRule.Name)
	_ = d.Set("update_time", recordRule.UpdateTime)
	_ = d.Set("template_id", recordRule.TemplateId)
	//_ = d.Set("content", recordRule.Content)
	_ = d.Set("cluster_id", recordRule.ClusterId)

	return nil
}

func resourceTencentCloudTkeTmpRecordRuleYamlUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_record_rule_yaml.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewModifyPrometheusRecordRuleYamlRequest()

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	request.InstanceId = &ids[0]
	request.Name = &ids[1]

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))

			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().ModifyPrometheusRecordRuleYaml(request)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})

			if err != nil {
				return err
			}

			return resourceTencentCloudTkeTmpRecordRuleYamlRead(d, meta)
		}
	}

	return nil
}

func resourceTencentCloudTkeTmpRecordRuleYamlDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_record_rule_yaml.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	if err := service.DeletePrometheusRecordRuleYaml(ctx, ids[0], ids[1]); err != nil {
		return err
	}

	return nil
}
