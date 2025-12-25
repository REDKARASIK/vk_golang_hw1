package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// Tут вы пишете код.
// Обратите внимание - в этом задании запрещены глобальные переменные.
// Если хочется, то для красоты можно разнести логику по разным файлам

type MicroService struct {
	UnimplementedAdminServer
	UnimplementedBizServer

	aclData       map[string][]string
	logStreams    map[chan *Event]bool
	logStreamsMux sync.RWMutex
	stats           map[string]uint64
	statsByConsumer map[string]uint64
	statsMu         sync.RWMutex
}

func (s *MicroService) Add(ctx context.Context, req *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (s *MicroService) Check(ctx context.Context, req *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (s *MicroService) Test(ctx context.Context, req *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}

func (s *MicroService) logEvent(event *Event) {
	s.logStreamsMux.RLock()
	defer s.logStreamsMux.RUnlock()

	for ch := range s.logStreams {
		select {
		case ch <- event:
		default:
		}
	}
}

func (s *MicroService) updateStats(method, consumer string) {
	s.statsMu.Lock()
	defer s.statsMu.Unlock()

	s.stats[method]++
	s.statsByConsumer[consumer]++
}

func (s *MicroService) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	consumers := md.Get("consumer")
	if len(consumers) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing consumer")
	}
	consumer := consumers[0]

	if err := s.checkACL(consumer, info.FullMethod); err != nil {
		return nil, err
	}

	p, _ := peer.FromContext(ctx)
	host := ""
	if p != nil && p.Addr != nil {
		host = p.Addr.String()
	}

	event := &Event{
		Timestamp: time.Now().Unix(),
		Consumer:  consumer,
		Method:    info.FullMethod,
		Host:      host,
	}

	s.logEvent(event)
	s.updateStats(info.FullMethod, consumer)

	return handler(ctx, req)
}

func (s *MicroService) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}

	consumers := md.Get("consumer")
	if len(consumers) == 0 {
		return status.Error(codes.Unauthenticated, "missing consumer")
	}
	consumer := consumers[0]

	if err := s.checkACL(consumer, info.FullMethod); err != nil {
		return err
	}

	p, _ := peer.FromContext(ctx)
	host := ""
	if p != nil && p.Addr != nil {
		host = p.Addr.String()
	}

	event := &Event{
		Timestamp: time.Now().Unix(),
		Consumer:  consumer,
		Method:    info.FullMethod,
		Host:      host,
	}

	s.logEvent(event)
	s.updateStats(info.FullMethod, consumer)

	return handler(srv, ss)
}

func (s *MicroService) checkACL(consumer, method string) error {
	if consumer == "" {
		return status.Error(codes.Unauthenticated, "missing consumer")
	}

	methods, ok := s.aclData[consumer]
	if !ok {
		return status.Error(codes.Unauthenticated, "unknown consumer")
	}

	for _, allowedMethod := range methods {
		if method == allowedMethod {
			return nil
		}

		if strings.HasSuffix(allowedMethod, "/*") {
			prefix := strings.TrimSuffix(allowedMethod, "/*")
			if strings.HasPrefix(method, prefix+"/") {
				return nil
			}
		}
	}

	return status.Error(codes.Unauthenticated, "access denied")
}

func (s *MicroService) Logging(req *Nothing, stream Admin_LoggingServer) error {
	logChan := make(chan *Event, 100)

	s.logStreamsMux.Lock()
	s.logStreams[logChan] = true
	s.logStreamsMux.Unlock()

	defer func() {
		s.logStreamsMux.Lock()
		delete(s.logStreams, logChan)
		s.logStreamsMux.Unlock()
		close(logChan)
	}()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case event := <-logChan:
			if err := stream.Send(event); err != nil {
				return err
			}
		}
	}
}

func (s *MicroService) Statistics(req *StatInterval, stream Admin_StatisticsServer) error {
	interval := time.Duration(req.GetIntervalSeconds()) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	s.statsMu.RLock()
	lastSent := make(map[string]uint64)
	lastSentConsumer := make(map[string]uint64)
	for k, v := range s.stats {
		lastSent[k] = v
	}
	for k, v := range s.statsByConsumer {
		lastSentConsumer[k] = v
	}
	s.statsMu.RUnlock()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			s.statsMu.RLock()
			stat := &Stat{
				Timestamp:  time.Now().Unix(),
				ByMethod:   make(map[string]uint64),
				ByConsumer: make(map[string]uint64),
			}
			for k, v := range s.stats {
				if diff := v - lastSent[k]; diff > 0 {
					stat.ByMethod[k] = diff
					lastSent[k] = v
				}
			}
			for k, v := range s.statsByConsumer {
				if diff := v - lastSentConsumer[k]; diff > 0 {
					stat.ByConsumer[k] = diff
					lastSentConsumer[k] = v
				}
			}
			s.statsMu.RUnlock()

			if err := stream.Send(stat); err != nil {
				return err
			}
		}
	}
}

func StartMyMicroservice(ctx context.Context, addr string, acl string) error {
	aclData := make(map[string][]string)
	if err := json.Unmarshal([]byte(acl), &aclData); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	microserv := &MicroService{
		aclData: aclData,
		logStreams: make(map[chan *Event]bool),
		stats: make(map[string]uint64),
		statsByConsumer: make(map[string]uint64),
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(microserv.unaryInterceptor),
		grpc.StreamInterceptor(microserv.streamInterceptor),
	)

	RegisterAdminServer(server, microserv)
	RegisterBizServer(server, microserv)

	go func() {
		<-ctx.Done()
		server.GracefulStop()
		lis.Close()
	}()

	go func() {
		server.Serve(lis)
	}()

	return nil
}
