package tests

import (
	"bytes"
	"github.com/go-masonry/mortar/interfaces/cfg"
	mock_cfg "github.com/go-masonry/mortar/interfaces/cfg/mock"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
	"testing"
)

type middlewareSuite struct {
	suite.Suite

	ctrl         *gomock.Controller
	cfgMock      *mock_cfg.MockConfig
	app          *fxtest.App
	loggerOutput bytes.Buffer
	// populate
	logExtractor      log.ContextExtractor
	clientInterceptor grpc.UnaryClientInterceptor
	serverInterceptor grpc.UnaryServerInterceptor
}

func TestMiddleware(t *testing.T) {
	suite.Run(t, new(middlewareSuite))
}

func (s *middlewareSuite) SetupTest() {
	// This one runs before `BeforeTest`
	s.ctrl = gomock.NewController(s.T())
	s.cfgMock = mock_cfg.NewMockConfig(s.ctrl)
	s.loggerOutput = bytes.Buffer{} // init buffer
}

func (s *middlewareSuite) BeforeTest(suiteName, testName string) {
	var extraOptions fx.Option
	switch testName {
	case "TestLoggerGRPCIncomingContextExtractor":
		extraOptions = s.testLoggerGRPCIncomingContextExtractorBeforeTest()
	case "TestClientInterceptorHeaderCopier":
		extraOptions = s.testClientInterceptorHeaderCopierBeforeTest()
	case "TestLoggerGRPCInterceptor":
		extraOptions = s.testLoggerGRPCInterceptorBeforeTest()
	case "TestMonitorGRPCInterceptor":
		extraOptions = s.testMonitorGRPCInterceptorBeforeTest()
	default:
		s.T().Fatalf("no pre test logic found for %s", testName)
	}

	s.app = fxtest.New(s.T(),
		s.suiteOptions(),
		extraOptions,
	)
	s.app.RequireStart()
}

func (s *middlewareSuite) suiteOptions() fx.Option {
	return fx.Options(
		fx.Provide(func() cfg.Config {
			return s.cfgMock
		}),
	)
}

func (s *middlewareSuite) TearDownTest() {
	s.app.RequireStop()
	s.ctrl.Finish()
}