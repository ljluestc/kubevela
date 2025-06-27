package testing

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	admissionv1 "k8s.io/api/admission/v1"
	authenticationv1 "k8s.io/api/authentication/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// WebhookTestCase defines a test case for webhook testing
type WebhookTestCase struct {
	Name           string
	Object         runtime.Object
	ExpectedStatus bool
	ExpectedError  string
	UserInfo       *authenticationv1.UserInfo
}

// SetupWebhookRequest creates an admission request for webhook testing
func SetupWebhookRequest(t *testing.T, obj runtime.Object, userInfo *authenticationv1.UserInfo) admission.Request {
	assert := assert.New(t)

	objJSON, err := json.Marshal(obj)
	assert.NoError(err, "should marshal object to JSON")

	req := admission.Request{
		AdmissionRequest: admissionv1.AdmissionRequest{
			Object: runtime.RawExtension{
				Raw: objJSON,
			},
			Operation: admissionv1.Create,
		},
	}

	if userInfo != nil {
		req.UserInfo = *userInfo
	}

	return req
}

// RunWebhookTest tests a webhook handler with the given test case
func RunWebhookTest[T admission.Handler, TestingTType interface {
	Errorf(format string, args ...interface{})
}](t TestingTType, handler T, testCase WebhookTestCase) {
	assert := assert.New(t)

	// Convert t to *testing.T for SetupWebhookRequest
	tObj, ok := any(t).(*testing.T)
	if !ok {
		panic("t must be *testing.T")
	}

	req := SetupWebhookRequest(tObj, testCase.Object, testCase.UserInfo)
	resp := handler.Handle(context.Background(), req)

	assert.Equal(testCase.ExpectedStatus, resp.Allowed, "webhook response allowed status should match expected")

	if testCase.ExpectedError != "" {
		assert.Contains(resp.Result.Message, testCase.ExpectedError, "webhook response should contain expected error")
	}
}

// CreateWebhookServer creates a test server for webhook handlers
func CreateWebhookServer(t *testing.T, pattern string, handler http.HandlerFunc) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(pattern, handler)

	return httptest.NewServer(mux)
}
