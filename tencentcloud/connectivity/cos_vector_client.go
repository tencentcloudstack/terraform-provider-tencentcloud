package connectivity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"syscall"
	"time"

	cos "github.com/tencentyun/cos-go-sdk-v5"
)

var cosVectorRegionPattern = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)+$`)

type cosVectorClientState struct {
	client    *cos.Client
	transport *http.Transport
	auth      *cos.AuthorizationTransport
	region    string
}

// UseTencentCosVectorClient 在同一配置实例内复用向量连接，不跨账号缓存。
func (me *TencentCloudClient) UseTencentCosVectorClient() (*cos.Client, error) {
	me.cosVectorMu.Lock()
	defer me.cosVectorMu.Unlock()
	if me.Credential == nil || me.Credential.SecretId == "" || me.Credential.SecretKey == "" {
		return nil, errors.New("COS Vectors requires configured TencentCloud credentials")
	}
	if !cosVectorRegionPattern.MatchString(me.Region) {
		return nil, errors.New("COS Vectors requires a valid provider region")
	}
	if cached := me.cosVectorConn; cached != nil && cached.region == me.Region {
		id, key, token, _ := cached.auth.GetCredential()
		if id == me.Credential.SecretId && key == me.Credential.SecretKey && token == me.Credential.Token {
			return cached.client, nil
		}
	}
	u, err := cos.NewVectorURL(me.Region, true)
	if err != nil {
		return nil, errors.New("cannot construct COS Vectors endpoint")
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	transport.DialContext = dialCosVectorPublicEndpoint
	auth := &cos.AuthorizationTransport{
		SecretID: me.Credential.SecretId, SecretKey: me.Credential.SecretKey,
		SessionToken: me.Credential.Token, Transport: &CosVectorResponseTransport{Transport: transport},
	}
	client := cos.NewClient(&cos.BaseURL{VectorURL: u}, &http.Client{
		Timeout:       30 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse },
		Transport:     auth,
	})
	// 建桶结果可能未知；只轮询 Get，不重复发送创建请求。
	client.Conf.RetryOpt.Count = 1
	client.UserAgent = fmt.Sprintf("Terraform-%s", GetReqClientVersion())
	if previous := me.cosVectorConn; previous != nil {
		previous.transport.CloseIdleConnections()
	}
	me.cosVectorConn = &cosVectorClientState{client: client, transport: transport, auth: auth, region: me.Region}
	return client, nil
}

const CosVectorInvalidResponseHeader = "X-Terraform-Cos-Vector-Invalid-Response"

// CosVectorResponseTransport 防止 SDK 将损坏的错误体当成资源不存在。
type CosVectorResponseTransport struct{ Transport http.RoundTripper }

func (t *CosVectorResponseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	data, readErr := io.ReadAll(io.LimitReader(resp.Body, (1<<20)+1))
	_ = resp.Body.Close()
	if resp.Header == nil {
		resp.Header = make(http.Header)
	}
	resp.Header.Del(CosVectorInvalidResponseHeader)
	invalid := readErr != nil || len(data) > 1<<20
	if resp.StatusCode != http.StatusOK {
		var body struct {
			Message string `json:"message"`
		}
		invalid = invalid || json.Unmarshal(data, &body) != nil || body.Message == ""
	} else if req.URL.Path != "/DeleteVectorBucket" && req.URL.Path != "/DeleteIndex" && req.URL.Path != "/PutVectorBucketPolicy" && req.URL.Path != "/DeleteVectorBucketPolicy" {
		invalid = invalid || !json.Valid(data)
	}
	if invalid {
		if resp.Header == nil {
			resp.Header = make(http.Header)
		}
		problem := "InvalidResponse"
		var netErr net.Error
		switch {
		case errors.Is(readErr, context.Canceled):
			problem = "Canceled"
		case errors.Is(readErr, context.DeadlineExceeded), errors.As(readErr, &netErr) && netErr.Timeout():
			problem = "Timeout"
		case errors.Is(readErr, io.EOF), errors.Is(readErr, io.ErrUnexpectedEOF), errors.Is(readErr, syscall.ECONNRESET), errors.Is(readErr, syscall.EPIPE):
			problem = "ConnectionFailure"
		}
		resp.Header.Set(CosVectorInvalidResponseHeader, problem)
		data = nil
	}
	resp.Body = io.NopCloser(bytes.NewReader(data))
	return resp, nil
}

func dialCosVectorPublicEndpoint(ctx context.Context, network, address string) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, errors.New("COS Vectors endpoint DNS lookup failed")
	}
	if len(ips) == 0 {
		return nil, errors.New("COS Vectors endpoint has no addresses")
	}
	for _, addr := range ips {
		if !cosVectorPublicIP(addr.IP) {
			return nil, errors.New("COS Vectors endpoint resolved to a non-public address; private endpoints require explicit authorization")
		}
	}
	var dialErr error
	for _, addr := range ips {
		conn, err := (&net.Dialer{Timeout: 10 * time.Second}).DialContext(ctx, network, net.JoinHostPort(addr.IP.String(), port))
		if err == nil {
			return conn, nil
		}
		dialErr = err
	}
	return nil, dialErr
}

func cosVectorPublicIP(ip net.IP) bool {
	if !ip.IsGlobalUnicast() || ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
		return false
	}
	if v4 := ip.To4(); v4 != nil {
		switch v4[0] {
		case 0, 9, 10, 11, 21, 30, 127:
			return false
		}
		if v4[0] == 100 && v4[1] >= 64 && v4[1] <= 127 {
			return false
		}
	}
	return true
}
