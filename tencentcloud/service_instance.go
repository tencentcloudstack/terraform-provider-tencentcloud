package tencentcloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/athom/goset"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/zqfan/tencentcloud-sdk-go/client"
)

func waitInstanceReachTargetStatus(client *client.Client, instanceIds []string, targetStatus string) (instanceStatusMap map[string]string, err error) {
	return waitInstanceReachOneOfTargetStatusList(client, instanceIds, []string{targetStatus})
}

func waitInstanceReachOneOfTargetStatusList(client *client.Client, instanceIds []string, targetStatuses []string) (instanceStatusMap map[string]string, err error) {
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		instanceStatusMap, _, err = queryInstancesStatus(client, instanceIds)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		var targetStatusInstanceIds []string
		for instanceId, instanceStatus := range instanceStatusMap {
			if goset.IsIncluded(targetStatuses, instanceStatus) {
				targetStatusInstanceIds = append(targetStatusInstanceIds, instanceId)
			}
		}
		if len(targetStatusInstanceIds) == len(instanceStatusMap) {
			return nil
		}
		_, _, penddingInstanceIds, _ := goset.Difference(instanceIds, targetStatusInstanceIds)
		return resource.RetryableError(fmt.Errorf("query instances status, pendding instanceIds: %v, not all instances are ready, retry...", penddingInstanceIds))
	})
	return
}

func waitInstanceOperationReachTargetStatus(client *client.Client, instanceIds []string, targetStatus string) (operationStatusMap map[string]string, err error) {
	return waitInstanceOperationReachOneOfTargetStatusList(client, instanceIds, []string{targetStatus})
}

func waitInstanceOperationReachOneOfTargetStatusList(client *client.Client, instanceIds []string, targetStatuses []string) (operationStatusMap map[string]string, err error) {
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, operationStatusMap, err = queryInstancesStatus(client, instanceIds)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		var targetStatusInstanceIds []string
		for instanceId, instanceStatus := range operationStatusMap {
			if goset.IsIncluded(targetStatuses, instanceStatus) {
				targetStatusInstanceIds = append(targetStatusInstanceIds, instanceId)
			}
		}
		if len(targetStatusInstanceIds) == len(operationStatusMap) {
			return nil
		}
		_, _, penddingInstanceIds, _ := goset.Difference(instanceIds, targetStatusInstanceIds)
		return resource.RetryableError(fmt.Errorf("query instances status, pendding instanceIds: %v, not all instances are ready, retry...", penddingInstanceIds))
	})
	return
}

func queryInstancesStatus(client *client.Client, instanceIds []string) (instanceStatusMap map[string]string, operationStatusMap map[string]string, err error) {
	if len(instanceIds) == 0 {
		err = fmt.Errorf("queryInstancesStatus, empty instanceIds")
		return
	}

	params := map[string]string{
		"Version": "2017-03-12",
		"Action":  "DescribeInstances",
	}
	for i, instanceId := range instanceIds {
		paramsKey := fmt.Sprintf("InstanceIds.%v", i+1)
		params[paramsKey] = instanceId
	}

	var response string
	response, err = client.SendRequest("cvm", params)
	if err != nil {
		return
	}

	var jsonresp struct {
		Response struct {
			Error struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			}
			TotalCount  int
			InstanceSet []struct {
				InstanceId           string
				InstanceState        string
				LatestOperation      string
				LatestOperationState string
			}
			RequestId string
		}
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		return
	}
	if jsonresp.Response.Error.Code != "" {
		err = fmt.Errorf(
			"call DescribeInstancesStatus got error, code:%v, message:%v",
			jsonresp.Response.Error.Code,
			jsonresp.Response.Error.Message,
		)
		return
	}
	if jsonresp.Response.TotalCount == 0 && len(jsonresp.Response.InstanceSet) == 0 {
		err = fmt.Errorf(instanceNotFoundErrorMsg(instanceIds))
		return
	}

	instanceStatusMap = make(map[string]string)
	operationStatusMap = make(map[string]string)
	instanceStatusList := jsonresp.Response.InstanceSet
	for _, instanceStatus := range instanceStatusList {
		instanceStatusMap[instanceStatus.InstanceId] = instanceStatus.InstanceState
		if instanceStatus.LatestOperation != "" {
			operationStatusMap[instanceStatus.InstanceId] = instanceStatus.LatestOperationState
		}
	}

	return
}

func instanceNotFoundErrorMsg(instanceIds []string) string {
	return fmt.Sprintf("no such instances: %v", instanceIds)
}

func stopInstance(client *client.Client, instanceId string) error {
	return operateInstance(client, instanceId, "StopInstances")
}

func startInstance(client *client.Client, instanceId string) error {
	return operateInstance(client, instanceId, "StartInstances")
}

func operateInstance(client *client.Client, instanceId string, action string) error {
	params := map[string]string{
		"Version":       "2017-03-12",
		"Action":        action,
		"InstanceIds.0": instanceId,
	}
	if action == "StopInstances" {
		params["ForceStop"] = "TRUE"
	}

	return runBasicActionWithRetry(client, params)
}

func modifyInstancesProject(client *client.Client, instanceIds []string, newId int) error {
	params := map[string]string{
		"Version":   "2017-03-12",
		"Action":    "ModifyInstancesProject",
		"ProjectId": fmt.Sprintf("%v", newId),
	}
	for i, instanceId := range instanceIds {
		paramKey := fmt.Sprintf("InstanceIds.%v", i)
		paramValue := instanceId
		params[paramKey] = paramValue
	}

	return runBasicActionWithRetry(client, params)
}

func renameInstancesName(client *client.Client, instanceIds []string, newName string) error {
	_, errs := validateInstanceName(interface{}(newName), "")
	if len(errs) > 0 {
		return errs[0]
	}

	params := map[string]string{
		"Version":      "2017-03-12",
		"Action":       "ModifyInstancesAttribute",
		"InstanceName": newName,
	}
	for i, instanceId := range instanceIds {
		paramKey := fmt.Sprintf("InstanceIds.%v", i+1)
		paramValue := instanceId
		params[paramKey] = paramValue
	}

	return runBasicActionWithRetry(client, params)
}

func resetInstancePassword(client *client.Client, instanceId string, newPassword string) error {
	params := map[string]string{
		"Version":       "2017-03-12",
		"Action":        "ResetInstancesPassword",
		"Password":      newPassword,
		"InstanceIds.0": instanceId,
	}
	return operateInstanceBetweenStopAndStart(client, params, instanceId, func() error {
		// just wait, 95% change password will be reset
		// a little stupid, but it's ok for now
		// need to improve it when there is API for reseting proccess query
		time.Sleep(42 * time.Second)
		return nil
	})
}

func resetInstanceSystem(client *client.Client, d *schema.ResourceData, instanceId string, newImageId string) error {
	params := map[string]string{
		"Version":    "2017-03-12",
		"Action":     "ResetInstance",
		"ImageId":    newImageId,
		"InstanceId": instanceId,
	}

	if v, ok := d.GetOk("disable_security_service"); ok {
		disable := v.(bool)
		if disable {
			params["EnhancedService.SecurityService.Enabled"] = "FALSE"
		}
	}
	if v, ok := d.GetOk("disable_monitor_service"); ok {
		disable := v.(bool)
		if disable {
			params["EnhancedService.MonitorService.Enabled"] = "FALSE"
		}
	}
	if v, ok := d.GetOk("password"); ok {
		password := v.(string)
		params["LoginSettings.Password"] = password
	}
	if v, ok := d.GetOk("key_pair"); ok {
		keyId := v.(string)
		params["LoginSettings.KeyIds.0"] = keyId
	}
	if d.HasChange("system_disk_size") {
		if d.HasChange("system_disk_type") {
			return fmt.Errorf("system_disk_type is not allowed to change when reinstall system")
		}
		if systemDiskType, ok := d.GetOk("system_disk_type"); ok {
			params["SystemDisk.DiskType"] = systemDiskType.(string)
		}
		if systemDiskSize, ok := d.GetOk("system_disk_size"); ok {
			diskSize := systemDiskSize.(int)
			params["SystemDisk.DiskSize"] = fmt.Sprintf("%v", diskSize)
		}
	}

	if err := runBasicActionWithRetry(client, params); err != nil {
		return err
	}

	if _, err := waitInstanceOperationReachTargetStatus(client, []string{instanceId}, "OPERATING"); err != nil {
		return err
	}

	if _, err := waitInstanceOperationReachTargetStatus(client, []string{instanceId}, "SUCCESS"); err != nil {
		return err
	}

	return nil
}

func operateInstanceBetweenStopAndStart(client *client.Client, params map[string]string, instanceId string, wait func() error) error {
	// make sure instance status is STOPPED before bind/unbind key pair
	if err := stopInstance(client, instanceId); err != nil {
		if !errAlreadyStopped(err, instanceId) {
			return err
		}
	}

	if _, err := waitInstanceReachTargetStatus(client, []string{instanceId}, "STOPPED"); err != nil {
		return err
	}

	err := runBasicActionWithRetry(client, params)
	if err != nil {
		return err
	}

	// wait until operation done
	if err := wait(); err != nil {
		return err
	}

	// recover instance to running
	if err := startInstance(client, instanceId); err != nil {
		return err
	}
	if _, err := waitInstanceReachTargetStatus(client, []string{instanceId}, "RUNNING"); err != nil {
		return err
	}

	return nil
}

func runBasicActionWithRetry(client *client.Client, params map[string]string) error {
	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err := client.SendRequest("cvm", params)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		var jsonresp struct {
			Response struct {
				Error struct {
					Code    string `json:"Code"`
					Message string `json:"Message"`
				}
				RequestId string
			}
		}
		err = json.Unmarshal([]byte(response), &jsonresp)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if retryable(jsonresp.Response.Error.Code, jsonresp.Response.Error.Message) {
			return resource.RetryableError(fmt.Errorf(jsonresp.Response.Error.Message))
		}
		if jsonresp.Response.Error.Code != "" {
			return resource.NonRetryableError(fmt.Errorf(jsonresp.Response.Error.Message))
		}
		return nil
	})
}

func bindInstanceWithSgIds(client *client.Client, instanceId string, sgIds []string) (err error) {
	params := map[string]string{
		"Version":       "2017-03-12",
		"Action":        "ModifyInstancesAttribute",
		"InstanceIds.0": instanceId,
	}
	for i, sgId := range sgIds {
		paramKey := fmt.Sprintf("SecurityGroups.%v", i)
		params[paramKey] = sgId
	}
	return runBasicActionWithRetry(client, params)
}
