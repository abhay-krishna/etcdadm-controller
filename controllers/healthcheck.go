package controllers

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

type etcdHealthCheckConfig struct {
	etcdHttpClient *http.Client
}

type EtcdHealthCheckResponse struct {
	Health string `json:"health"`
}

type portNotOpenError struct{}

func (h *portNotOpenError) Error() string {
	return "etcd endpoint port is not open"
}

var portNotOpenErr = &portNotOpenError{}

func (r *EtcdadmClusterReconciler) performEndpointHealthCheck(ctx context.Context, cluster *clusterv1.Cluster, endpoint string) error {
	if err := r.setEtcdHttpClientIfUnset(ctx, cluster); err != nil {
		return err
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return errors.Wrapf(err, "invalid etcd endpoint url")
	}
	if !isPortOpen(ctx, u.Host) {
		return portNotOpenErr
	}

	client := r.etcdHealthCheckConfig.etcdHttpClient
	healthCheckURL := getMemberHealthCheckEndpoint(endpoint)
	r.Log.Info(fmt.Sprintf("Performing healthcheck on endpoint %s", healthCheckURL))
	req, err := http.NewRequest("GET", healthCheckURL, nil)
	if err != nil {
		return errors.Wrap(err, "error creating healthcheck request")
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error checking etcd member health")
	}
	// reuse connection
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Wrap(err, "Etcd member not ready, retry")
	}

	if err := parseEtcdHealthCheckOutput(body); err != nil {
		return errors.Wrap(err, fmt.Sprintf("etcd member %v failed healthcheck", endpoint))
	}
	r.Log.Info(fmt.Sprintf("Etcd member %v ready", endpoint))
	return nil
}

func parseEtcdHealthCheckOutput(data []byte) error {
	obj := EtcdHealthCheckResponse{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	if obj.Health == "true" {
		return nil
	}
	return fmt.Errorf("/health returned %q", obj.Health)
}

func (r *EtcdadmClusterReconciler) setEtcdHttpClientIfUnset(ctx context.Context, cluster *clusterv1.Cluster) error {
	if r.etcdHealthCheckConfig.etcdHttpClient != nil {
		return nil
	}
	caCertPool := x509.NewCertPool()
	caCert, err := r.getCACert(ctx, cluster)
	if err != nil {
		return err
	}
	caCertPool.AppendCertsFromPEM(caCert)

	clientCert, err := r.getClientCerts(ctx, cluster)
	if err != nil {
		return errors.Wrap(err, "Error getting client cert for healthcheck")
	}
	r.etcdHealthCheckConfig.etcdHttpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{clientCert},
			},
		},
	}
	return nil
}

func isPortOpen(ctx context.Context, endpoint string) bool {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return false
	}

	if conn != nil {
		conn.Close()
		return true
	}

	return false
}