package tests

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	httphelper "github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestALBASGService(t *testing.T) {
	t.Parallel()
	const (
		expectedMinSize          = 2
		expectedDesiredCapacity  = 2
		serverPort               = 8080
		expectedBody             = "Hello, World!"
		instanceIDCapturePattern = `Instance ID: (i-[a-z0-9]+)`
	)
	opts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]any{
			"min_size":         expectedMinSize,
			"desired_capacity": expectedDesiredCapacity,
			"server_port":      serverPort,
		},
	})
	defer terraform.Destroy(t, opts)
	_, err := terraform.InitAndApplyE(t, opts)
	require.NoError(t, err)

	t.Run("outputs", func(t *testing.T) {
		tests := []struct {
			name       string
			outputName string
			assertion  func(t *testing.T, val string)
		}{
			{
				name:       "alb dns name format",
				outputName: "alb_dns_name",
				assertion: func(t *testing.T, val string) {
					assert.NotEmpty(t, val)
					assert.Contains(t, val, ".elb.amazonaws.com")
				},
			},
			{
				name:       "alb url validity",
				outputName: "alb_url",
				assertion: func(t *testing.T, val string) {
					u, err := url.Parse(val)
					require.NoError(t, err)
					assert.Equal(t, "http", u.Scheme)
					assert.Equal(t, terraform.Output(t, opts, "alb_dns_name"), u.Host)
				},
			},
			{
				name:       "asg name prefix",
				outputName: "asg_name",
				assertion: func(t *testing.T, val string) {
					assert.True(t, strings.HasPrefix(val, "alb-asg-service"))
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				got := terraform.Output(t, opts, tt.outputName)
				tt.assertion(t, got)
			})
		}
	})

	t.Run("http cluster connectivity and load balancing", func(t *testing.T) {
		albURL := terraform.Output(t, opts, "alb_url")
		maxRetries := 30
		timeBetweenRetries := 10 * time.Second
		httphelper.HttpGetWithRetryWithCustomValidation(
			t,
			albURL,
			nil,
			maxRetries,
			timeBetweenRetries,
			func(statusCode int, body string) bool {
				return statusCode == http.StatusOK && strings.Contains(body, expectedBody)
			},
		)

		t.Run("verify traffic distribution", func(t *testing.T) {
			instanceIDs := make(map[string]struct{})
			re := regexp.MustCompile(instanceIDCapturePattern)
			for range 10 {
				_, body := httphelper.HttpGet(t, albURL, nil)
				match := re.FindStringSubmatch(body)
				if len(match) > 1 {
					instanceIDs[match[1]] = struct{}{}
				}
				time.Sleep(500 * time.Millisecond)
			}
			got := len(instanceIDs)
			assert.GreaterOrEqual(t, got, expectedDesiredCapacity)
		})
	})
}
