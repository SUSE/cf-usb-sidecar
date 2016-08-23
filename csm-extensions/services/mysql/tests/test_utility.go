package IntegrationTest

import (
	"flag"
	"strconv"
	"strings"
	"testing"

	swaggerClient "github.com/go-swagger/go-swagger/client"
	httpClient "github.com/go-swagger/go-swagger/httpkit/client"
	"github.com/go-swagger/go-swagger/strfmt"
	csmClient "github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager-client/client"
)

var (
	host                = flag.String("host", "0.0.0.0", "HTTP listen port")
	port                = flag.Int("port", 8081, "HTTP listen port")
	runIntegrationTests = flag.Bool("integration", false, "Run the api tests")
	apiKey              = flag.String("apikey", "", "API KEY")
	transport           *httpClient.Runtime
	client              *csmClient.CatlogServiceManager
	authFunc            swaggerClient.AuthInfoWriter

	serviceRoot = "../../CATALOGS/hpe-catalog/services"
)

type errorResponse struct {
	Code    int
	Message string
}

func errorMessage(err error) string {
	if err != nil {
		errorMessage := strings.Split(err.Error(), "Message:")
		if len(errorMessage) == 2 {
			return strings.Trim(errorMessage[1], "}")
		}
	}
	return ""
}

func errorCode(err error) string {
	if err != nil {
		errorMessage := strings.Split(err.Error(), "][")
		if len(errorMessage) == 2 {
			return strings.Split(errorMessage[1], "]")[0]
		}
	}
	return ""
}

func skipWhenUnitTesting(test *testing.T) {
	if !*runIntegrationTests {
		test.Skip("To run this test, use: go test -integration=true")
	}
}

func init() {
	flag.Parse()
	transportHost := *host + ":" + strconv.Itoa(*port)
	transport = httpClient.New(transportHost, "", []string{"http"})
	client = csmClient.New(transport, strfmt.Default)
	authFunc = httpClient.APIKeyAuth("x-sidecar-token", "header", "sidecar-auth-token")
}
