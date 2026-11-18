package makers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	makerssdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewMakersService(client *connectivity.TencentCloudClient) MakersService {
	return MakersService{client: client}
}

type MakersService struct {
	client *connectivity.TencentCloudClient
}

// DescribeDeploymentLog returns the build log of the deployment identified by
// zoneId, projectId and deploymentId. The three keys form an exact query, so
// no pagination is involved. A nil LogUrl with a nil error is a valid outcome:
// the cloud side documents that LogUrl may be null when no log is available.
func (me *MakersService) DescribeDeploymentLog(ctx context.Context, request *makerssdk.DescribeDeploymentLogRequest) (logUrl *string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.DescribeDeploymentLogResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().DescribeDeploymentLogWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeDeploymentLog returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	logUrl = response.Response.LogUrl
	return
}

// Deployment terminal statuses reported by the DescribeDeployments API.
const (
	// DeploymentStatusSuccess means the build finished successfully.
	DeploymentStatusSuccess = "Success"
	// DeploymentStatusFailed means the build finished with an error.
	DeploymentStatusFailed = "Failed"
	// DeploymentStatusProcess means the build is still running.
	DeploymentStatusProcess = "Process"
	// DeploymentStatusInvalid means the deployment is invalid.
	DeploymentStatusInvalid = "Invalid"
	// DeploymentStatusPending means the build is queued and not started yet.
	DeploymentStatusPending = "Pending"
	// DeploymentStatusTimeout means the build exceeded the cloud-side time limit.
	DeploymentStatusTimeout = "Timeout"
	// DeploymentStatusCancelled means the build was cancelled.
	DeploymentStatusCancelled = "Cancelled"
)

// CreateDeployment triggers a new Makers build via the CreateDeployment API.
// The API is asynchronous: it returns a DeploymentId as soon as the build is
// accepted, and the build result must be confirmed separately through
// DescribeDeployments. DeploymentId may legitimately be nil on the wire; the
// caller decides whether that is acceptable.
func (me *MakersService) CreateDeployment(ctx context.Context, request *makerssdk.CreateDeploymentRequest) (deploymentId *string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.CreateDeploymentResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().CreateDeploymentWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("CreateDeployment returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	deploymentId = response.Response.DeploymentId
	return
}

// DescribeDeploymentById queries a single deployment by its exact DeploymentId.
// ZoneId and ProjectId scope the query. It returns the deployment object, or
// (nil, nil) when no deployment matches the id.
func (me *MakersService) DescribeDeploymentById(ctx context.Context, zoneId, projectId, deploymentId string) (deployment *makerssdk.Deployment, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := makerssdk.NewDescribeDeploymentsRequest()
	request.ZoneId = helper.String(zoneId)
	request.ProjectId = helper.String(projectId)
	request.DeploymentIds = []*string{helper.String(deploymentId)}
	request.Limit = helper.Int64(1)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *makerssdk.DescribeDeploymentsResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, err := me.client.UseMakersClient().DescribeDeploymentsWithContext(ctx, request)
		if err != nil {
			return tccommon.RetryError(err)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("DescribeDeployments returned an empty response")
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	// The query filters by an exact deployment id, so at most one entry can
	// match. An empty list is a legitimate "not found" outcome, not an error.
	for _, item := range response.Response.Deployments {
		if item != nil && item.DeploymentId != nil && *item.DeploymentId == deploymentId {
			deployment = item
			break
		}
	}
	return
}

// WaitMakersDeploymentComplete polls DescribeDeployments until the deployment
// reaches a terminal state (Success / Failed / Invalid / Timeout / Cancelled),
// or the timeout budget is exhausted. It returns the last observed deployment
// snapshot so the caller can attach it to a timeout diagnostic. A nil
// deployment is treated as "not visible yet" and keeps the loop polling.
func (me *MakersService) WaitMakersDeploymentComplete(ctx context.Context, zoneId, projectId, deploymentId string, timeout time.Duration, pollInterval time.Duration) (last *makerssdk.Deployment, err error) {
	deadline := time.Now().Add(timeout)
	for {
		latest, e := me.DescribeDeploymentById(ctx, zoneId, projectId, deploymentId)
		if e != nil {
			err = e
			return
		}
		last = latest
		if latest != nil && latest.Status != nil {
			switch *latest.Status {
			case DeploymentStatusSuccess:
				return
			case DeploymentStatusFailed, DeploymentStatusInvalid, DeploymentStatusTimeout, DeploymentStatusCancelled:
				err = fmt.Errorf("deployment %s ended in status %s", deploymentId, *latest.Status)
				return
			}
		}

		if time.Now().After(deadline) {
			err = fmt.Errorf("deployment %s did not converge within %s", deploymentId, timeout)
			return
		}
		time.Sleep(pollInterval)
	}
}
