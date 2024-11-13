package server

import (
	"context"
	"server/internal/dto"
	"server/internal/service"
	"server/internal/transport/grpc/proto"
)

type GRPCServer struct {
	proto.UnimplementedMetricServiceServer
	service *service.Service
}

func NewGRPCServer(s *service.Service) *GRPCServer {
	return &GRPCServer{service: s}
}

func (s *GRPCServer) SetUpdates(ctx context.Context, req *proto.SetUpdatesRequest) (*proto.SetUpdatesResponse, error) {
	var metrics []dto.Metric
	for _, m := range req.Metrics {
		metrics = append(metrics, dto.Metric{
			Name:  m.Name,
			Type:  m.Type,
			Value: &m.Value,
			Delta: &m.Delta,
		})
	}

	updatedMetrics, err := s.service.SetUpdates(metrics)
	if err != nil {
		return nil, err
	}

	var response proto.SetUpdatesResponse
	for _, m := range updatedMetrics {
		response.Metrics = append(response.Metrics, &proto.MetricModel{
			Id:        uint32(m.ID),
			Name:      m.Name,
			Type:      m.Type,
			Value:     *m.Value,
			Delta:     *m.Delta,
			CreatedAt: m.CreatedAt.String(),
		})
	}

	return &response, nil
}

func (s *GRPCServer) SetMetric(ctx context.Context, req *proto.SetMetricRequest) (*proto.SetMetricResponse, error) {
	metric := dto.Metric{
		Name:  req.Name,
		Type:  req.Type,
		Value: &req.Value,
	}

	updatedMetric, err := s.service.SetMetric(metric)
	if err != nil {
		return nil, err
	}

	return &proto.SetMetricResponse{
		Metric: &proto.MetricModel{
			Id:        uint32(updatedMetric.ID),
			Name:      updatedMetric.Name,
			Type:      updatedMetric.Type,
			Value:     *updatedMetric.Value,
			Delta:     *updatedMetric.Delta,
			CreatedAt: updatedMetric.CreatedAt.String(),
		},
	}, nil
}

func (s *GRPCServer) GetMetricValue(ctx context.Context, req *proto.GetMetricValueRequest) (*proto.GetMetricValueResponse, error) {
	value, err := s.service.GetMetricValue(req.Name, req.Type)
	if err != nil {
		return nil, err
	}

	return &proto.GetMetricValueResponse{Value: *value}, nil
}

func (s *GRPCServer) GetHTML(ctx context.Context, req *proto.GetHTMLRequest) (*proto.GetHTMLResponse, error) {
	responseWriter := NewCustomResponseWriter()

	err := s.service.GetHTML(responseWriter)
	if err != nil {
		return nil, err
	}

	htmlContent := responseWriter.Buffer.String()

	return &proto.GetHTMLResponse{HtmlContent: htmlContent}, nil
}
