package stat

import (
	"admin/pkg/logger"
	"context"

	pb "github.com/ShopOnGO/admin-proto/pkg/service"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

// StatServiceDeps содержит зависимости сервиса статистики
type StatServiceDeps struct {
	StatRepository *StatRepository
}

// StatService реализует StatServiceServer (из protobuf)
type StatService struct {
	pb.UnimplementedStatServiceServer // Встраиваем, чтобы обеспечить forward compatibility
	StatRepository                    *StatRepository
}

// NewStatService создаёт новый сервис статистики
func NewStatService(deps *StatServiceDeps) *StatService {
	return &StatService{
		StatRepository: deps.StatRepository,
	}
}

// AddClick обрабатывает gRPC-запрос на добавление клика
// обработка ошибок фиктивна потому что иначе grpc не хочет
func (s *StatService) AddClick(ctx context.Context, req *pb.ClickRequest) (*pb.ClickResponse, error) {
	linkId := req.LinkId
	logger.Infof("Обрабатываем клик по ссылке ID %d", linkId)
	s.StatRepository.AddClick(uint(linkId))

	return &pb.ClickResponse{}, nil
}
