// A separate runtime owner: never linked into the Terraform Provider.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	runtimecode "github.com/TencentCloudAgentRuntime/ags-go-sdk/sandbox/code"
	sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tcerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type ledger struct {
	RunID       string `json:"run_id"`
	Region      string `json:"region"`
	ToolID      string `json:"tool_id"`
	ToolName    string `json:"tool_name"`
	InstanceID  string `json:"instance_id,omitempty"`
	InstanceEnd bool   `json:"instance_stop_confirmed"`
	ToolEnd     bool   `json:"tool_delete_ack_and_empty"`
	ToolACK     bool   `json:"tool_delete_ack"`
}

func ptr[T any](v T) *T { return &v }
func val[T any](v *T) (zero T) {
	if v != nil {
		return *v
	}
	return zero
}

func save(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0600)
}

func safeError(err error) string {
	var api *tcerr.TencentCloudSDKError
	if errors.As(err, &api) {
		return fmt.Sprintf("code=%s request_id=%s", api.GetCode(), api.GetRequestId())
	}
	// Runtime v0.1.5 does not guarantee secret-safe formatted errors.
	return fmt.Sprintf("%T (raw error suppressed)", err)
}

func main() {
	mode := flag.String("mode", "read", "read, code, cleanup-instance or cleanup-tool")
	path := flag.String("ledger", "", "private ownership ledger")
	output := flag.String("output", "", "private observation file")
	cloud := flag.Bool("cloud", false, "explicitly enable operations on this run's dedicated resources")
	flag.Parse()
	if err := run(*mode, *path, *output, *cloud); err != nil {
		fmt.Fprintln(os.Stderr, safeError(err))
		os.Exit(1)
	}
}

func run(mode, path, output string, cloud bool) (runErr error) {
	if !cloud || path == "" || output == "" {
		return errors.New("cloud opt-in, ledger and output required")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var owned ledger
	if err := json.Unmarshal(b, &owned); err != nil {
		return err
	}
	if owned.RunID == "" || owned.Region == "" || !strings.HasPrefix(owned.ToolID, "sdt-") || !strings.HasPrefix(owned.ToolName, "tf-ags-"+owned.RunID+"-") {
		return errors.New("ledger does not identify a dedicated run-owned Tool")
	}
	cred := common.NewTokenCredential(os.Getenv("TENCENTCLOUD_SECRET_ID"), os.Getenv("TENCENTCLOUD_SECRET_KEY"), os.Getenv("TENCENTCLOUD_SECURITY_TOKEN"))
	if cred.GetSecretId() == "" || cred.GetSecretKey() == "" {
		return errors.New("credentials required through the environment")
	}
	p := profile.NewClientProfile()
	p.HttpProfile.ReqTimeout = 15
	c, err := sdk.NewClient(cred, owned.Region, p)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	read := sdk.NewDescribeSandboxToolListRequest()
	read.ToolIds = []*string{ptr(owned.ToolID)}
	response, err := c.DescribeSandboxToolListWithContext(ctx, read)
	if err != nil {
		return err
	}
	if response.Response == nil {
		return errors.New("missing point-query response")
	}
	if len(response.Response.SandboxToolSet) == 0 && mode == "cleanup-tool" && (owned.ToolEnd || owned.ToolACK) {
		owned.ToolEnd = true
		if err := save(path, owned); err != nil {
			return err
		}
		return save(output, map[string]any{"cleanup": "PASS", "previous_delete_ack": true, "point_query_empty": true, "read_request_id": response.Response.RequestId})
	}
	if len(response.Response.SandboxToolSet) != 1 || val(response.Response.SandboxToolSet[0].ToolName) != owned.ToolName || val(response.Response.SandboxToolSet[0].ToolId) != owned.ToolID {
		return errors.New("exact Tool identity/name ownership check failed; no mutation sent")
	}
	tool := response.Response.SandboxToolSet[0]
	if mode == "read" {
		return save(output, map[string]any{"tool": tool, "request_id": response.Response.RequestId})
	}
	if mode == "cleanup-tool" {
		if owned.InstanceID != "" && !owned.InstanceEnd {
			return errors.New("release the owned instance before deleting the Tool")
		}
		req := sdk.NewDeleteSandboxToolRequest()
		req.ToolId = ptr(owned.ToolID)
		deleted, err := c.DeleteSandboxToolWithContext(ctx, req)
		if err != nil {
			return err
		}
		for {
			check, err := c.DescribeSandboxToolListWithContext(ctx, read)
			if err != nil {
				return err
			}
			if check.Response != nil && len(check.Response.SandboxToolSet) == 0 {
				owned.ToolEnd = true
				if err := save(path, owned); err != nil {
					return err
				}
				return save(output, map[string]any{"cleanup": "PASS", "delete_ack": true, "point_query_empty": true, "delete_request_id": deleted.Response.RequestId, "read_request_id": check.Response.RequestId, "provider_state_removed": false})
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Second):
			}
		}
	}
	if mode == "cleanup-instance" {
		return stopInstance(c, &owned, path, output)
	}
	if mode != "code" || val(tool.ToolType) != "code-interpreter" || owned.InstanceID != "" {
		return errors.New("code mode requires a code Tool and a ledger without a previous instance")
	}
	// Public Cloud Start accepts the exact Terraform ID without a timeout or
	// custom override; v0.1.5's convenience Create accepts only a name and
	// supplies its own timeout, so use Start + public runtime Connect instead.
	start := sdk.NewStartSandboxInstanceRequest()
	start.ToolId, start.ClientToken = ptr(owned.ToolID), ptr(owned.RunID+"-code")
	started, err := c.StartSandboxInstanceWithContext(ctx, start)
	if started != nil && started.Response != nil && started.Response.Instance != nil {
		owned.InstanceID = val(started.Response.Instance.InstanceId)
	}
	if owned.InstanceID != "" {
		defer func() {
			if e := stopInstance(c, &owned, path, output+".cleanup.json"); e != nil {
				fmt.Fprintln(os.Stderr, "INDEPENDENT_CLEANUP_REQUIRED:", safeError(e))
				runErr = e
			}
		}()
		if e := save(path, owned); e != nil {
			return e
		}
	}
	if err != nil {
		return err
	}
	if owned.InstanceID == "" {
		return errors.New("accepted instance identity missing; reconcile by dedicated Tool before another Start")
	}
	for {
		req := sdk.NewDescribeSandboxInstanceListRequest()
		req.InstanceIds = []*string{ptr(owned.InstanceID)}
		instances, err := c.DescribeSandboxInstanceListWithContext(ctx, req)
		if err != nil {
			return err
		}
		if instances.Response != nil && len(instances.Response.InstanceSet) == 1 {
			instance := instances.Response.InstanceSet[0]
			if val(instance.ToolId) != owned.ToolID {
				return errors.New("runtime instance does not belong to the Terraform Tool")
			}
			if val(instance.Status) == "RUNNING" {
				break
			}
			if val(instance.Status) == "FAILED" {
				return errors.New("instance entered FAILED")
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
	sandbox, err := runtimecode.Connect(ctx, owned.InstanceID, runtimecode.WithClient(c))
	if err != nil {
		return err
	}
	result, err := sandbox.Code.RunCode(ctx, "print(6 * 7)", nil, nil)
	if err != nil {
		return err
	}
	stdout := strings.Join(result.Logs.Stdout, "")
	if err := save(output, map[string]any{"kind": "SDK_CODE", "sdk_version": "v0.1.5", "instance_id": owned.InstanceID, "observed_tool_id": owned.ToolID, "stdout": stdout, "execution_error": result.Error, "start_request_id": started.Response.RequestId}); err != nil {
		return err
	}
	if result.Error != nil || strings.TrimSpace(stdout) != "42" {
		return errors.New("code smoke output was not 42 without an execution error")
	}
	return nil
}

func stopInstance(c *sdk.Client, owned *ledger, path, output string) error {
	if owned.InstanceID == "" || owned.InstanceEnd {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	check := sdk.NewDescribeSandboxInstanceListRequest()
	check.InstanceIds = []*string{ptr(owned.InstanceID)}
	current, err := c.DescribeSandboxInstanceListWithContext(ctx, check)
	if err != nil {
		return err
	}
	if current.Response == nil || len(current.Response.InstanceSet) != 1 || val(current.Response.InstanceSet[0].ToolId) != owned.ToolID {
		return errors.New("instance ownership could not be confirmed; no Stop sent")
	}
	req := sdk.NewStopSandboxInstanceRequest()
	req.InstanceId = ptr(owned.InstanceID)
	stopped, err := c.StopSandboxInstanceWithContext(ctx, req)
	if err != nil {
		return err
	}
	for {
		req := sdk.NewDescribeSandboxInstanceListRequest()
		req.InstanceIds = []*string{ptr(owned.InstanceID)}
		response, err := c.DescribeSandboxInstanceListWithContext(ctx, req)
		if err != nil {
			return err
		}
		if response.Response != nil && (len(response.Response.InstanceSet) == 0 || len(response.Response.InstanceSet) == 1 && val(response.Response.InstanceSet[0].Status) == "STOPPED") {
			owned.InstanceEnd = true
			if err := save(path, owned); err != nil {
				return err
			}
			return save(output, map[string]any{"cleanup": "PASS", "stop_ack": true, "stop_request_id": stopped.Response.RequestId})
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
	}
}
