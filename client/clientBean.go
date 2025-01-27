// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

//go:generate mockgen -package $GOPACKAGE -source $GOFILE -destination clientBean_mock.go -self_package github.com/uber/cadence/client

package client

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	clientworker "go.uber.org/cadence/worker"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/api/peer"
	"go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/peer/roundrobin"
	"go.uber.org/yarpc/transport/grpc"
	"go.uber.org/yarpc/transport/tchannel"

	"github.com/uber/cadence/client/admin"
	"github.com/uber/cadence/client/frontend"
	"github.com/uber/cadence/client/history"
	"github.com/uber/cadence/client/matching"
	"github.com/uber/cadence/common"
	"github.com/uber/cadence/common/authorization"
	"github.com/uber/cadence/common/cluster"
	"github.com/uber/cadence/common/log"
	"github.com/uber/cadence/common/log/tag"
)

const (
	defaultRefreshInterval = time.Second * 10
)

type (
	// Bean in an collection of clients
	Bean interface {
		GetHistoryClient() history.Client
		SetHistoryClient(client history.Client)
		GetMatchingClient(domainIDToName DomainIDToNameFunc) (matching.Client, error)
		SetMatchingClient(client matching.Client)
		GetFrontendClient() frontend.Client
		SetFrontendClient(client frontend.Client)
		GetRemoteAdminClient(cluster string) admin.Client
		SetRemoteAdminClient(cluster string, client admin.Client)
		GetRemoteFrontendClient(cluster string) frontend.Client
		SetRemoteFrontendClient(cluster string, client frontend.Client)
	}

	DispatcherOptions struct {
		AuthProvider clientworker.AuthorizationProvider
	}

	// DispatcherProvider provides a dispatcher to a given address
	DispatcherProvider interface {
		GetTChannel(name string, address string, options *DispatcherOptions) (*yarpc.Dispatcher, error)
		GetGRPC(name string, address string, options *DispatcherOptions) (*yarpc.Dispatcher, error)
	}

	clientBeanImpl struct {
		sync.Mutex
		historyClient         history.Client
		matchingClient        atomic.Value
		frontendClient        frontend.Client
		remoteAdminClients    map[string]admin.Client
		remoteFrontendClients map[string]frontend.Client
		factory               Factory
	}

	dnsDispatcherProvider struct {
		interval time.Duration
		logger   log.Logger
	}
	dnsUpdater struct {
		interval     time.Duration
		dnsAddress   string
		port         string
		currentPeers map[string]struct{}
		list         peer.List
		logger       log.Logger
	}
	dnsRefreshResult struct {
		updates  peer.ListUpdates
		newPeers map[string]struct{}
		changed  bool
	}
	aPeer struct {
		addrPort string
	}
)

// NewClientBean provides a collection of clients
func NewClientBean(factory Factory, dispatcherProvider DispatcherProvider, clusterMetadata cluster.Metadata) (Bean, error) {

	historyClient, err := factory.NewHistoryClient()
	if err != nil {
		return nil, err
	}

	remoteAdminClients := map[string]admin.Client{}
	remoteFrontendClients := map[string]frontend.Client{}
	for clusterName, info := range clusterMetadata.GetAllClusterInfo() {
		if !info.Enabled {
			continue
		}

		var dispatcherOptions *DispatcherOptions
		if info.AuthorizationProvider.Enable {
			authProvider, err := authorization.GetAuthProviderClient(info.AuthorizationProvider.PrivateKey)
			if err != nil {
				return nil, err
			}
			dispatcherOptions = &DispatcherOptions{
				AuthProvider: authProvider,
			}
		}

		var dispatcher *yarpc.Dispatcher
		var err error
		switch info.RPCTransport {
		case tchannel.TransportName:
			dispatcher, err = dispatcherProvider.GetTChannel(info.RPCName, info.RPCAddress, dispatcherOptions)
		case grpc.TransportName:
			dispatcher, err = dispatcherProvider.GetGRPC(info.RPCName, info.RPCAddress, dispatcherOptions)
		}

		if err != nil {
			return nil, err
		}

		adminClient, err := factory.NewAdminClientWithTimeoutAndDispatcher(
			info.RPCName,
			admin.DefaultTimeout,
			admin.DefaultLargeTimeout,
			dispatcher,
		)
		if err != nil {
			return nil, err
		}

		frontendClient, err := factory.NewFrontendClientWithTimeoutAndDispatcher(
			info.RPCName,
			frontend.DefaultTimeout,
			frontend.DefaultLongPollTimeout,
			dispatcher,
		)
		if err != nil {
			return nil, err
		}

		remoteAdminClients[clusterName] = adminClient
		remoteFrontendClients[clusterName] = frontendClient
	}

	return &clientBeanImpl{
		factory:               factory,
		historyClient:         historyClient,
		frontendClient:        remoteFrontendClients[clusterMetadata.GetCurrentClusterName()],
		remoteAdminClients:    remoteAdminClients,
		remoteFrontendClients: remoteFrontendClients,
	}, nil
}

func (h *clientBeanImpl) GetHistoryClient() history.Client {
	return h.historyClient
}

func (h *clientBeanImpl) SetHistoryClient(
	client history.Client,
) {

	h.historyClient = client
}

func (h *clientBeanImpl) GetMatchingClient(domainIDToName DomainIDToNameFunc) (matching.Client, error) {
	if client := h.matchingClient.Load(); client != nil {
		return client.(matching.Client), nil
	}
	return h.lazyInitMatchingClient(domainIDToName)
}

func (h *clientBeanImpl) SetMatchingClient(
	client matching.Client,
) {

	h.matchingClient.Store(client)
}

func (h *clientBeanImpl) GetFrontendClient() frontend.Client {
	return h.frontendClient
}

func (h *clientBeanImpl) SetFrontendClient(
	client frontend.Client,
) {

	h.frontendClient = client
}

func (h *clientBeanImpl) GetRemoteAdminClient(cluster string) admin.Client {
	client, ok := h.remoteAdminClients[cluster]
	if !ok {
		panic(fmt.Sprintf(
			"Unknown cluster name: %v with given cluster client map: %v.",
			cluster,
			h.remoteAdminClients,
		))
	}
	return client
}

func (h *clientBeanImpl) SetRemoteAdminClient(
	cluster string,
	client admin.Client,
) {

	h.remoteAdminClients[cluster] = client
}

func (h *clientBeanImpl) GetRemoteFrontendClient(cluster string) frontend.Client {
	client, ok := h.remoteFrontendClients[cluster]
	if !ok {
		panic(fmt.Sprintf(
			"Unknown cluster name: %v with given cluster client map: %v.",
			cluster,
			h.remoteFrontendClients,
		))
	}
	return client
}

func (h *clientBeanImpl) SetRemoteFrontendClient(
	cluster string,
	client frontend.Client,
) {

	h.remoteFrontendClients[cluster] = client
}

func (h *clientBeanImpl) lazyInitMatchingClient(domainIDToName DomainIDToNameFunc) (matching.Client, error) {
	h.Lock()
	defer h.Unlock()
	if cached := h.matchingClient.Load(); cached != nil {
		return cached.(matching.Client), nil
	}
	client, err := h.factory.NewMatchingClient(domainIDToName)
	if err != nil {
		return nil, err
	}
	h.matchingClient.Store(client)
	return client, nil
}

// NewDNSYarpcDispatcherProvider create a dispatcher provider which handles with IP address
func NewDNSYarpcDispatcherProvider(logger log.Logger, interval time.Duration) DispatcherProvider {
	if interval <= 0 {
		interval = defaultRefreshInterval
	}
	return &dnsDispatcherProvider{
		interval: interval,
		logger:   logger,
	}
}

func (p *dnsDispatcherProvider) GetTChannel(serviceName string, address string, options *DispatcherOptions) (*yarpc.Dispatcher, error) {
	tchanTransport, err := tchannel.NewTransport(
		tchannel.ServiceName(serviceName),
		// this aim to get rid of the annoying popup about accepting incoming network connections
		tchannel.ListenAddr("127.0.0.1:0"),
	)
	if err != nil {
		return nil, err
	}

	peerList := roundrobin.New(tchanTransport)
	peerListUpdater, err := newDNSUpdater(peerList, address, p.interval, p.logger)
	if err != nil {
		return nil, err
	}
	peerListUpdater.Start()
	outbound := tchanTransport.NewOutbound(peerList)

	p.logger.Info("Creating TChannel dispatcher outbound", tag.Address(address))
	return p.createOutboundDispatcher(serviceName, outbound, options)
}

func (p *dnsDispatcherProvider) GetGRPC(serviceName string, address string, options *DispatcherOptions) (*yarpc.Dispatcher, error) {
	grpcTransport := grpc.NewTransport()

	peerList := roundrobin.New(grpcTransport)
	peerListUpdater, err := newDNSUpdater(peerList, address, p.interval, p.logger)
	if err != nil {
		return nil, err
	}
	peerListUpdater.Start()
	outbound := grpcTransport.NewOutbound(peerList)

	p.logger.Info("Creating GRPC dispatcher outbound", tag.Address(address))
	return p.createOutboundDispatcher(serviceName, outbound, options)
}

func (p *dnsDispatcherProvider) createOutboundDispatcher(serviceName string, outbound transport.UnaryOutbound, options *DispatcherOptions) (*yarpc.Dispatcher, error) {
	cfg := yarpc.Config{
		Name: crossDCCaller,
		Outbounds: yarpc.Outbounds{
			serviceName: transport.Outbounds{
				Unary:       outbound,
				ServiceName: serviceName,
			},
		},
	}
	if options != nil && options.AuthProvider != nil {
		cfg.OutboundMiddleware = yarpc.OutboundMiddleware{
			Unary: &outboundMiddleware{authProvider: options.AuthProvider},
		}
	}

	// Attach the outbound to the dispatcher (this will add middleware/logging/etc)
	dispatcher := yarpc.NewDispatcher(cfg)

	if err := dispatcher.Start(); err != nil {
		return nil, err
	}
	return dispatcher, nil
}

type outboundMiddleware struct {
	authProvider clientworker.AuthorizationProvider
}

func (om *outboundMiddleware) Call(ctx context.Context, request *transport.Request, out transport.UnaryOutbound) (*transport.Response, error) {
	if om.authProvider != nil {
		token, err := om.authProvider.GetAuthToken()
		if err != nil {
			return nil, err
		}
		request.Headers = request.Headers.
			With(common.AuthorizationTokenHeaderName, string(token))
	}
	return out.Call(ctx, request)
}

func newDNSUpdater(list peer.List, dnsPort string, interval time.Duration, logger log.Logger) (*dnsUpdater, error) {
	ss := strings.Split(dnsPort, ":")
	if len(ss) != 2 {
		return nil, fmt.Errorf("incorrect DNS:Port format")
	}
	return &dnsUpdater{
		interval:     interval,
		logger:       logger,
		list:         list,
		dnsAddress:   ss[0],
		port:         ss[1],
		currentPeers: make(map[string]struct{}),
	}, nil
}

func (d *dnsUpdater) Start() {
	go func() {
		for {
			now := time.Now()
			res, err := d.refresh()
			if err != nil {
				d.logger.Error("Failed to update DNS", tag.Error(err), tag.Address(d.dnsAddress))
			}
			if res != nil && res.changed {
				if len(res.updates.Additions) > 0 {
					d.logger.Info("Add new peers by DNS lookup", tag.Address(d.dnsAddress), tag.Addresses(identifiersToStringList(res.updates.Additions)))
				}
				if len(res.updates.Removals) > 0 {
					d.logger.Info("Remove stale peers by DNS lookup", tag.Address(d.dnsAddress), tag.Addresses(identifiersToStringList(res.updates.Removals)))
				}

				err := d.list.Update(res.updates)
				if err != nil {
					d.logger.Error("Failed to update peerList", tag.Error(err), tag.Address(d.dnsAddress))
				}
				d.currentPeers = res.newPeers
			}
			sleepDu := now.Add(d.interval).Sub(now)
			time.Sleep(sleepDu)
		}
	}()
}

func (d *dnsUpdater) refresh() (*dnsRefreshResult, error) {
	resolver := net.DefaultResolver
	ips, err := resolver.LookupHost(context.Background(), d.dnsAddress)
	if err != nil {
		return nil, err
	}
	newPeers := map[string]struct{}{}
	for _, ip := range ips {
		adr := fmt.Sprintf("%v:%v", ip, d.port)
		newPeers[adr] = struct{}{}
	}

	updates := peer.ListUpdates{
		Additions: make([]peer.Identifier, 0),
		Removals:  make([]peer.Identifier, 0),
	}
	changed := false
	// remove if it doesn't exist anymore
	for addr := range d.currentPeers {
		if _, ok := newPeers[addr]; !ok {
			changed = true
			updates.Removals = append(
				updates.Removals,
				aPeer{addrPort: addr},
			)
		}
	}

	// add if it doesn't exist before
	for addr := range newPeers {
		if _, ok := d.currentPeers[addr]; !ok {
			changed = true
			updates.Additions = append(
				updates.Additions,
				aPeer{addrPort: addr},
			)
		}
	}

	return &dnsRefreshResult{
		updates:  updates,
		newPeers: newPeers,
		changed:  changed,
	}, nil
}

func (a aPeer) Identifier() string {
	return a.addrPort
}

func identifiersToStringList(ids []peer.Identifier) []string {
	ss := make([]string, 0, len(ids))
	for _, id := range ids {
		ss = append(ss, id.Identifier())
	}
	return ss
}
