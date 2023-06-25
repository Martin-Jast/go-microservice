package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Martin-Jast/go-microservice/persistence"
	"github.com/Martin-Jast/go-microservice/utest"
	"github.com/gavv/httpexpect"
	"github.com/gorilla/mux"
)

// BaseHandlerTest basic structure of an request validation test -> Do not set properties that are not going to be used
type BaseHandlerTest struct {
	Name                  string
	Req                   interface{}
	Path                  string
	SetupPreTestDBs       func(ctx context.Context, th *testHandler) *SetupResult
	AssertPosTestDBStates func(ctx context.Context, th *testHandler, sr *SetupResult, res *httpexpect.Response) error
	ExpectedCode          int
	ExpectedBody          string
	HTTPMethod            string
	InvalidContext        bool
	MountRequest          func(sr *SetupResult) interface{}
	MountPath             func(sr *SetupResult) (string, interface{})
	RouteHandler          func(sr *SetupResult) map[string]func(http.ResponseWriter, *http.Request)
}

type SetupResult struct {
	// Add here any elements that are generated for the individual test and are needed to create path, request or to assert responses
	BaseDoc    []persistence.BaseModel
}

type singleHandlerTest struct {
	BaseHandlerTest
	utest.TestMock
}
type internalHT struct {
	base BaseHandlerTest
	t    *testing.T
	th   *testHandler
	ctx  context.Context
}

// GetTestName return test name
func (inst internalHT) GetTestName() string {
	return inst.base.Name
}

// ExecuteTest execute test and return an eventual error
func (inst internalHT) ExecuteTest() error {
	return handlerTest(inst.ctx, inst.base, inst.t, inst.th)
}

func handlerTest(ctx context.Context, tc BaseHandlerTest, t *testing.T, th *testHandler) error {
	// Setup environment
	// DB Setup
	var setupResult *SetupResult
	if tc.SetupPreTestDBs != nil {
		setupResult = tc.SetupPreTestDBs(ctx, th)
	}

	// TODO: Drop all LOCAL databases after each test -> Some of the DBs don't have the Drop method in the interface
	// defer th.dropDatabase(ctx)

	// HTTP setup
	// Create router to handle outside calls
	router := mux.NewRouter()
	if tc.RouteHandler != nil {
		routesHandlers := tc.RouteHandler(setupResult)
		for outPath := range routesHandlers {
			router.HandleFunc(outPath, routesHandlers[outPath])
		}
	}
	ts := httptest.NewServer(router)
	defer ts.Close()

	expect := th.createHTTPExpect(t)
	var httpFunc func(path string, pathargs ...interface{}) *httpexpect.Request
	switch tc.HTTPMethod {
	case "POST":
		httpFunc = expect.POST
	case "GET":
		httpFunc = expect.GET
	case "PUT":
		httpFunc = expect.PUT
	case "DELETE":
		httpFunc = expect.DELETE
	default:
		log.Panic("test: " + tc.Name + " : No HTTPMethod passed -> Please Specify the method to be called")
	}
	if tc.Path == "" && tc.MountPath == nil {
		log.Panic("No Path or MountPath given, no way to know where the request should go")
	}

	// Use or Mount Request
	var request interface{}
	if tc.MountRequest != nil {
		request = tc.MountRequest(setupResult)
	} else {
		request = tc.Req
	}
	// Use or Mount Path
	var path string
	var query interface{}
	if tc.MountPath != nil {
		path, query = tc.MountPath(setupResult)
	} else {
		path = tc.Path
	}
	// Mount request
	reqAsString, _ := json.Marshal(request)
	preRequest := httpFunc(path).
		WithText(string(reqAsString))
	if query != nil {
		preRequest = preRequest.WithQueryObject(query)
	}
	// Send
	resp := preRequest.Expect().
		Status(tc.ExpectedCode)
	if tc.ExpectedBody != "" {
		resp.Body().Equal(tc.ExpectedBody)
	}
	if resp.Body().Raw() != "" && t.Failed() {
		// Remove this if it gets too verbose
		fmt.Println("Test-" + t.Name() + ": " + resp.Body().Raw())
	}
	if tc.AssertPosTestDBStates != nil {
		err := tc.AssertPosTestDBStates(ctx, th, setupResult, resp)
		return err
	}
	return nil
}
func transformSingleBaseToHandlerTest(base BaseHandlerTest, t *testing.T) singleHandlerTest {
	th, ctx := createTestHandler()
	tMock := internalHT{base, t, th, ctx}
	valTest := singleHandlerTest{
		base, tMock,
	}
	return valTest
}

func transformAllBaseToHandlerTest(base []BaseHandlerTest, t *testing.T) []utest.TestMock {
	handlerTests := []utest.TestMock{}
	for i := range base {
		handlerTests = append(handlerTests, transformSingleBaseToHandlerTest(base[i], t))
	}
	return handlerTests
}

// ExecHandlerTest run handler tests between the expected request format and the mock input
func ExecHandlerTest(base []BaseHandlerTest, t *testing.T) []error {
	valt := transformAllBaseToHandlerTest(base, t)
	return utest.RunTest(valt, t)
}
