package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudTkeClusterAttachment() *schema.Resource {
	schemaBody := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "ID of the cluster.",
		},
		"instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "ID of the CVM instance, this machine will be reloaded.",
		},
		"hostname": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "When the system is reinstalled, you can specify the HostName of the modified instance. The name should be a combination of 2 to 15 characters comprised of letters (case insensitive), numbers, and hyphens (-). Period (.) is not supported, and the name cannot be a string of pure numbers. Other types (such as Linux) of instances: The name should be a combination of 2 to 60 characters, supporting multiple periods (.). The piece between two periods is composed of letters (case insensitive), numbers, and hyphens (-).",
		},
		"password": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validateAsConfigPassword,
			Description:  "Password to access, should be set if `key_ids` not set.",
		},
		"key_ids": {
			MaxItems:    1,
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "The key pair to use for the instance, it looks like skey-16jig7tx, it should be set if `password` not set.",
		},

		//computer
		"security_groups": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "A list of security group ids after attach to cluster.",
		},
	}

	return &schema.Resource{
		Create: resourceTencentCloudTkeClusterAttachmentCreate,
		Read:   resourceTencentCloudTkeClusterAttachmentRead,
		Delete: resourceTencentCloudTkeClusterAttachmentDelete,
		Schema: schemaBody,
	}
}

func resourceTencentCloudTkeClusterAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_attachment.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	tkeService := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	cvmService := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId, clusterId := "", ""

	if items := strings.Split(d.Id(), "_"); len(items) != 2 {
		return fmt.Errorf("the resource id is corrupted")
	} else {
		instanceId, clusterId = items[0], items[1]
	}

	/*tke has been deleted*/
	info, has, err := tkeService.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			info, has, err = tkeService.DescribeCluster(ctx, clusterId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return nil
	}
	if !has {
		d.SetId("")
		return nil
	}

	/*cvm has been deleted*/
	var instance *cvm.Instance
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, err = cvmService.DescribeInstanceById(ctx, instanceId)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if instance == nil {
		d.SetId("")
		return nil
	}

	/*attachment has been  deleted*/
	_, workers, err := tkeService.DescribeClusterInstances(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, workers, err = tkeService.DescribeClusterInstances(ctx, clusterId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	has = false
	for _, worker := range workers {
		if worker.InstanceId == instanceId {
			has = true
		}
	}

	if !has {
		d.SetId("")
		return nil
	}

	if len(instance.LoginSettings.KeyIds) > 0 {
		_ = d.Set("key_ids", instance.LoginSettings.KeyIds)
	}
	_ = d.Set("security_groups", helper.StringsInterfaces(instance.SecurityGroupIds))
	return nil
}

func resourceTencentCloudTkeClusterAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_attachment.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	tkeService := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	cvmService := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	request := tke.NewAddExistedInstancesRequest()

	instanceId := helper.String(d.Get("instance_id").(string))
	request.ClusterId = helper.String(d.Get("cluster_id").(string))
	request.InstanceIds = []*string{instanceId}
	request.HostName = helper.String(d.Get("hostname").(string))
	request.LoginSettings = &tke.LoginSettings{}

	var loginSettingsNumbers = 0

	if v, ok := d.GetOk("key_ids"); ok {
		request.LoginSettings.KeyIds = helper.Strings(helper.InterfacesStrings(v.([]interface{})))
		loginSettingsNumbers++
	}

	if v, ok := d.GetOk("password"); ok {
		request.LoginSettings.Password = helper.String(v.(string))
		loginSettingsNumbers++
	}

	if loginSettingsNumbers != 1 {
		return fmt.Errorf("parameters `key_ids` and `password` must set and only set one")
	}

	/*cvm has been  attached*/
	var err error
	_, workers, err := tkeService.DescribeClusterInstances(ctx, *request.ClusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, workers, err = tkeService.DescribeClusterInstances(ctx, *request.ClusterId)
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	has := false
	for _, worker := range workers {
		if worker.InstanceId == *instanceId {
			has = true
		}
	}
	if has {
		return fmt.Errorf("instance %s has been attached to cluster %s,can not attach again", *instanceId, *request.ClusterId)
	}

	var response *tke.AddExistedInstancesResponse

	if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = tkeService.client.UseTkeClient().AddExistedInstances(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("add existed instance %s to cluster %s error,reason %v", *instanceId, *request.ClusterId, err)
	}
	var success = false
	for _, v := range response.Response.SuccInstanceIds {
		if *v == *instanceId {
			d.SetId(*instanceId + "_" + *request.ClusterId)
			success = true
		}
	}

	if !success {
		return fmt.Errorf("add existed instance %s to cluster %s error, instance not in success instanceIds", *instanceId, *request.ClusterId)
	}

	/*wait for cvm status*/
	if err = resource.Retry(4*writeRetryTimeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, *instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if instance != nil && *instance.InstanceState == CVM_STATUS_RUNNING {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cvm instance %s status is %s, retry...", *instanceId, *instance.InstanceState))
	}); err != nil {
		return err
	}

	/*wait for tke init ok */
	err = resource.Retry(4*writeRetryTimeout, func() *resource.RetryError {
		_, workers, err = tkeService.DescribeClusterInstances(ctx, *request.ClusterId)
		if err != nil {
			return retryError(err, InternalError)
		}
		has := false
		for _, worker := range workers {
			if worker.InstanceId == *instanceId {
				has = true
				if worker.InstanceState == "failed" {
					return resource.NonRetryableError(fmt.Errorf("cvm instance %s attach to cluster %s fail,reason:%s",
						*instanceId, *request.ClusterId, worker.FailedReason))
				}

				if worker.InstanceState != "running" {
					return resource.RetryableError(fmt.Errorf("cvm instance  %s in tke status is %s, retry...",
						*instanceId, worker.InstanceState))
				}

			}
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("cvm instance %s not exist in tke instance list", *instanceId))
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudTkeClusterAttachmentRead(d, meta)
}

func resourceTencentCloudTkeClusterAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster_attachment.delete")()

	tkeService := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId, clusterId := "", ""

	if items := strings.Split(d.Id(), "_"); len(items) != 2 {
		return fmt.Errorf("the resource id is corrupted")
	} else {
		instanceId, clusterId = items[0], items[1]
	}

	request := tke.NewDeleteClusterInstancesRequest()

	request.ClusterId = &clusterId
	request.InstanceIds = []*string{
		&instanceId,
	}
	request.InstanceDeleteMode = helper.String("retain")

	var err error

	if err = resource.Retry(4*writeRetryTimeout, func() *resource.RetryError {
		_, err := tkeService.client.UseTkeClient().DeleteClusterInstances(request)
		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
			if e.GetCode() == "InternalError.Param" &&
				strings.Contains(e.GetMessage(), `PARAM_ERROR[some instances []is not in right state`) {
				return nil
			}
		}

		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
