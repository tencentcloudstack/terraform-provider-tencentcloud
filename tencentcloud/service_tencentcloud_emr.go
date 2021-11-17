package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type EMRService struct {
	client *connectivity.TencentCloudClient
}

func (me *EMRService) CreateInstance(ctx context.Context, d *schema.ResourceData) error {
	logId := getLogId(ctx)
	request := emr.NewCreateInstanceRequest()
	if v, ok := d.GetOk("product_id"); ok {
		request.ProductId = common.Uint64Ptr((uint64)(v.(int)))
	}

	if v, ok := d.GetOk("vpc_settings"); ok {
		value := v.(map[string]interface{})
		var vpcId string
		var subnetId string

		if subV, ok := value["vpc_id"]; ok {
			vpcId = subV.(string)
		}
		if subV, ok := value["subnet_id"]; ok {
			subnetId = subV.(string)
		}
		vpcSettings := &emr.VPCSettings{VpcId: &vpcId, SubnetId: &subnetId}
		request.VPCSettings = vpcSettings
	}

	if v, ok := d.GetOk("softwares"); ok {
		softwares := v.([]interface{})
		request.Software = make([]*string, 0)
		for _, software := range softwares {
			request.Software = append(request.Software, common.StringPtr(software.(string)))
		}
	}

	if v, ok := d.GetOk("resource_spec"); ok {
		tmpResourceSpec := v.([]interface{})
		resourceSpec := tmpResourceSpec[0].(map[string]interface{})
		request.ResourceSpec = &emr.NewResourceSpec{}
		for k, v := range resourceSpec {
			if k == "master_resource_spec" {
				if len(v.([]interface{})) > 0 {
					request.ResourceSpec.MasterResourceSpec = helper.ParseResource(v.([]interface{})[0].(map[string]interface{}))
				}
			} else if k == "core_resource_spec" {
				if len(v.([]interface{})) > 0 {
					request.ResourceSpec.CoreResourceSpec = helper.ParseResource(v.([]interface{})[0].(map[string]interface{}))
				}
			} else if k == "task_resource_spec" {
				if len(v.([]interface{})) > 0 {
					request.ResourceSpec.TaskResourceSpec = helper.ParseResource(v.([]interface{})[0].(map[string]interface{}))
				}
			} else if k == "master_count" {
				request.ResourceSpec.MasterCount = common.Int64Ptr((int64)(v.(int)))
			} else if k == "core_count" {
				request.ResourceSpec.CoreCount = common.Int64Ptr((int64)(v.(int)))
			} else if k == "task_count" {
				request.ResourceSpec.TaskCount = common.Int64Ptr((int64)(v.(int)))
			} else if k == "common_resource_spec" {
				if len(v.([]interface{})) > 0 {
					request.ResourceSpec.CommonResourceSpec = helper.ParseResource(v.([]interface{})[0].(map[string]interface{}))
				}
			} else if k == "common_count" {
				request.ResourceSpec.CommonCount = common.Int64Ptr((int64)(v.(int)))
			}
		}
	}

	if v, ok := d.GetOk("support_ha"); ok {
		request.SupportHA = common.Uint64Ptr((uint64)(v.(int)))
	} else {
		request.SupportHA = common.Uint64Ptr(0)
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = common.StringPtr(v.(string))
	}

	if v, ok := d.GetOk("pay_mode"); ok {
		request.PayMode = common.Uint64Ptr((uint64)(v.(int)))
	}
	if v, ok := d.GetOk("placement"); ok {
		request.Placement = &emr.Placement{}
		placement := v.(map[string]interface{})

		if projectId, ok := placement["project_id"]; ok {
			projectIdInt64, _ := strconv.ParseInt(projectId.(string), 10, 64)
			request.Placement.ProjectId = common.Int64Ptr(projectIdInt64)
		} else {
			request.Placement.ProjectId = common.Int64Ptr(0)
		}
		if zone, ok := placement["zone"]; ok {
			request.Placement.Zone = common.StringPtr(zone.(string))
		}
	}

	if v, ok := d.GetOk("time_span"); ok {
		request.TimeSpan = common.Uint64Ptr((uint64)(v.(int)))
	}
	if v, ok := d.GetOk("time_unit"); ok {
		request.TimeUnit = common.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("login_settings"); ok {
		request.LoginSettings = &emr.LoginSettings{}
		loginSettings := v.(map[string]interface{})
		if password, ok := loginSettings["password"]; ok {
			request.LoginSettings.Password = common.StringPtr(password.(string))
		}
		if publicKeyId, ok := loginSettings["public_key_id"]; ok {
			request.LoginSettings.PublicKeyId = common.StringPtr(publicKeyId.(string))
		}
	}

	instanceId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		//API: https://cloud.tencent.com/document/api/589/34261
		response, err := me.client.UseEmrClient().CreateInstance(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err)
		}
		instanceId = *response.Response.InstanceId
		return nil
	})
	if err != nil {
		return err
	}
	d.Set("instance_id", instanceId)
	return nil
}

func (me *EMRService) DescribeInstances(ctx context.Context, filters map[string]interface{}) (clusters []*emr.ClusterInstancesInfo, errRet error) {
	logId := getLogId(ctx)
	request := emr.NewDescribeInstancesRequest()

	ratelimit.Check(request.GetAction())
	// API: https://cloud.tencent.com/document/api/589/41707
	if v, ok := filters["display_strategy"]; ok {
		request.DisplayStrategy = common.StringPtr(v.(string))
	}
	if v, ok := filters["prefix_instance_ids"]; ok {
		request.InstanceIds = common.StringPtrs(v.([]string))
	}
	if v, ok := filters["project_id"]; ok {
		request.ProjectId = common.Int64Ptr(v.(int64))
	}

	response, err := me.client.UseEmrClient().DescribeInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	clusters = response.Response.ClusterList
	return
}
