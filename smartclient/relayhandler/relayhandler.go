package relayhandler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"

	"poktroll/smartclient/client"
	"poktroll/utils"
	applicationTypes "poktroll/x/application/types"
	svcTypes "poktroll/x/service/types"
	sessionTypes "poktroll/x/session/types"
)

var (
	errInvalidProtocol = errors.New("invalid protocol")
	errSessionMismatch = errors.New("session mismatch")
)

type RelayHandler struct {
	ctx        context.Context
	listenAddr string

	applicationQueryClient applicationTypes.QueryClient
	sessionQueryClient     sessionTypes.QueryClient
	blockQueryClient       *client.BlocksQueryClient

	applicationAddress string
	currentSession     *sessionTypes.Session
	providedServices   []*svcTypes.ServiceId
	servicesEndpoints  ServicesEndpoints
	servicesSessions   map[string]utils.Observable[*sessionTypes.Session]
	sessionMutex       sync.RWMutex

	endpointSelectionStrategy EndpointSelectionStrategy
	signer                    Signer
}

func NewRelayHandler(
	listenAddr string,
	applicationQueryClient applicationTypes.QueryClient,
	sessionQueryClient sessionTypes.QueryClient,
	blockQueryClient *client.BlocksQueryClient,
	applicationAddress string,
	endpointSelectionStrategy EndpointSelectionStrategy,
	signer Signer,
) *RelayHandler {
	return &RelayHandler{
		listenAddr:                listenAddr,
		applicationQueryClient:    applicationQueryClient,
		sessionQueryClient:        sessionQueryClient,
		blockQueryClient:          blockQueryClient,
		applicationAddress:        applicationAddress,
		servicesEndpoints:         &map[string][]svcTypes.Endpoint{},
		servicesSessions:          make(map[string]utils.Observable[*sessionTypes.Session]),
		sessionMutex:              sync.RWMutex{},
		endpointSelectionStrategy: endpointSelectionStrategy,
		signer:                    signer,
	}
}

func (relayHandler *RelayHandler) Start(ctx context.Context) (*RelayHandler, error) {
	relayHandler.ctx = ctx

	// get services provided by the application
	applicationRequest := &applicationTypes.QueryGetApplicationRequest{
		Address: relayHandler.applicationAddress,
	}
	appInfo, err := relayHandler.applicationQueryClient.Application(ctx, applicationRequest)
	if err != nil {
		return nil, err
	}
	relayHandler.providedServices = appInfo.Application.Services

	if err := relayHandler.serviceSessionListener(ctx); err != nil {
		return nil, err
	}

	go http.ListenAndServe(relayHandler.listenAddr, relayHandler)

	return relayHandler, nil
}

func (relayHandler *RelayHandler) serviceSessionListener(ctx context.Context) error {
	for _, service := range relayHandler.providedServices {
		sessionNotifee, err := relayHandler.sessionListener(ctx, service)
		if err != nil {
			return err
		}

		relayHandler.servicesSessions[service.Id] = sessionNotifee

		go relayHandler.serviceListener(ctx, service, sessionNotifee)
	}

	return nil
}

// servicers will change over sessions, so we need to listen for changes
func (relayHandler *RelayHandler) serviceListener(
	ctx context.Context,
	service *svcTypes.ServiceId,
	sessionNotifee utils.Observable[*sessionTypes.Session],
) error {
	ch := sessionNotifee.Subscribe(ctx).Ch()
	for session := range ch {
		relayHandler.bindToNewSessionServicers(session)
	}

	return nil
}

func (relayHandler *RelayHandler) bindToNewSessionServicers(session *sessionTypes.Session) {
	serviceEndpoints := make(map[string][]svcTypes.Endpoint)
	for _, servicer := range session.Servicers {
		for _, service := range servicer.Services {
			serviceEndpoints[service.Id.Id] = append(
				serviceEndpoints[service.Id.Id],
				service.Endpoints...,
			)
		}
	}

	relayHandler.SetServiceEndpoints(&serviceEndpoints, session)
}

func (relayHandler *RelayHandler) sessionListener(
	ctx context.Context,
	service *svcTypes.ServiceId,
) (utils.Observable[*sessionTypes.Session], error) {
	sessionNotifee, sessionNotifier := utils.NewControlledObservable[*sessionTypes.Session](nil)

	// current session at start
	session, err := relayHandler.getCurrentSession(ctx, service)
	if err != nil {
		return nil, err
	}

	ch := relayHandler.blockQueryClient.BlocksNotifee().Subscribe(ctx).Ch()
	go func() {
		for block := range ch {
			if block.Height() > session.SessionBlockStartHeight+session.NumBlocksPerSession {
				session, err = relayHandler.getCurrentSession(ctx, service)
				if err != nil {
					return
				}

				sessionNotifier <- session
			}
		}
	}()

	return sessionNotifee, nil
}

func (relayHandler *RelayHandler) getCurrentSession(ctx context.Context, service *svcTypes.ServiceId) (*sessionTypes.Session, error) {
	query := &sessionTypes.QueryGetSessionRequest{
		BlockHeight: relayHandler.blockQueryClient.LatestBlock(ctx).Height(),
		AppAddress:  relayHandler.applicationAddress,
		ServiceId:   service.Id,
	}
	currentSession, err := relayHandler.sessionQueryClient.GetSession(ctx, query)
	if err != nil {
		return nil, err
	}

	return &currentSession.Session, nil
}

func (relayHandler *RelayHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	protocol := req.URL.Scheme
	serviceId := strings.Split(path, "/")[0]

	if protocol == "http" {
		relayHandler.handleHTTPRelays(w, req, serviceId)
	} else if protocol == "ws" {
		relayHandler.handleWsRelays(w, req, serviceId)
	} else {
		utils.ReplyWithHTTPError(400, errInvalidProtocol, w)
		return
	}
}

func (relayHandler *RelayHandler) SetServiceEndpoints(serviceEndpoints ServicesEndpoints, session *sessionTypes.Session) {
	relayHandler.sessionMutex.Lock()
	defer relayHandler.sessionMutex.Unlock()
	relayHandler.servicesEndpoints = serviceEndpoints
	relayHandler.currentSession = session
}

func (relayHandler *RelayHandler) getSessionConnectionInfo(serviceId string) *svcTypes.Endpoint {
	relayHandler.sessionMutex.RLock()
	defer relayHandler.sessionMutex.RUnlock()

	serviceEndpoints, ok := (*relayHandler.servicesEndpoints)[serviceId]
	if !ok {
		return nil
	}

	// service should enforce its available protocols
	endpoint := relayHandler.endpointSelectionStrategy.SelectEndpoint(serviceEndpoints)
	return endpoint
}

type ChooseFirstEndpoint struct{}

func (c *ChooseFirstEndpoint) SelectEndpoint(endpoints []svcTypes.Endpoint) *svcTypes.Endpoint {
	return &endpoints[0]
}
