package tencentcloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/athom/goset"
	"github.com/hashicorp/terraform/helper/resource"
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
