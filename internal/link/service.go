package link

import (
	"admin/pkg/logger"
	"context"
	"math/rand"

	pb "github.com/ShopOnGO/admin-proto/pkg/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type LinkService struct {
	pb.UnimplementedLinkServiceServer
	LinkRepository *LinkRepository
}

func NewLinkService(linkRepository *LinkRepository) *LinkService {
	return &LinkService{
		LinkRepository: linkRepository,
	}
}

func (s *LinkService) Create(ctx context.Context, req *pb.CreateLinkRequest) (*pb.CreateLinkResponse, error) {
	link := NewLink(req.Url)
	for {
		existedLink, _ := s.LinkRepository.GetByHash(link.Hash)
		if existedLink == nil {
			break
		}
		link.GenerateHash()
	}
	newLink, err := s.LinkRepository.Create(link)
	if err != nil {
		logger.Errorf("failed to create link: %v ", err)
		return nil, status.Errorf(codes.Internal, ErrCreateLink, err)
	}
	return &pb.CreateLinkResponse{Link: ConvertToProtoLink(newLink)}, nil
}

func (s *LinkService) Update(ctx context.Context, req *pb.UpdateLinkRequest) (*pb.UpdateLinkResponse, error) {
	updatedLink, err := s.LinkRepository.Update(&Link{
		Model: gorm.Model{ID: uint(req.Id)},
		Url:   req.Url,
		Hash:  req.Hash,
	})
	if err != nil {
		logger.Errorf("failed to update link: %v", err)
		return nil, status.Errorf(codes.Internal, ErrUpdateLink, err)
	}
	return &pb.UpdateLinkResponse{Link: ConvertToProtoLink(updatedLink)}, nil
}

func (s *LinkService) Delete(ctx context.Context, req *pb.DeleteLinkRequest) (*pb.DeleteLinkResponse, error) {
	var err error
	err = s.LinkRepository.Delete(uint(req.Id), req.Unscoped)
	if err != nil {
		logger.Errorf("Failed to delete link: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete link: %v", err)
	}

	return &pb.DeleteLinkResponse{}, nil
}

func (s *LinkService) GetLinkByHash(ctx context.Context, req *pb.GetLinkByHashRequest) (*pb.GetLinkByHashResponse, error) {
	link, err := s.LinkRepository.GetByHash(req.Hash)
	if err != nil {
		logger.Errorf("failed to get by hash: %v", err)
		return nil, status.Errorf(codes.NotFound, ErrLinkNotFound, err)
	}
	return &pb.GetLinkByHashResponse{Link: ConvertToProtoLink(link)}, nil
}

func (s *LinkService) GetById(ctx context.Context, req *pb.GetLinkByIDRequest) (*pb.GetLinkByIDResponse, error) {
	link, err := s.LinkRepository.GetById(uint(req.Id))
	if err != nil {
		logger.Errorf("failed to get by id: %v", err)
		return nil, status.Errorf(codes.NotFound, ErrLinkNotFound, err)
	}
	return &pb.GetLinkByIDResponse{Link: ConvertToProtoLink(link)}, nil
}

func (s *LinkService) GetAllLinks(ctx context.Context, req *pb.GetAllLinksRequest) (*pb.GetAllLinksResponse, error) {
	links, err := s.LinkRepository.GetAll(int(req.Limit), int(req.Offset), req.IncludeDeleted)
	if err != nil {
		logger.Errorf("failed to get all links: %v", err)
		return nil, status.Errorf(codes.Internal, ErrGetLinks, err)
	}

	var grpcLinks []*pb.Link
	for _, l := range links {
		grpcLinks = append(grpcLinks, ConvertToProtoLink(&l))
	}
	count := s.LinkRepository.Count(req.IncludeDeleted)
	return &pb.GetAllLinksResponse{
		Links: grpcLinks,
		Count: count,
	}, nil
}

func (s *LinkService) CountLinks(ctx context.Context, req *pb.CountLinksRequest) (*pb.CountLinksResponse, error) {
	count := s.LinkRepository.Count(false)
	return &pb.CountLinksResponse{
		Count: count,
	}, nil
}

func NewLink(url string) *Link {
	link := &Link{
		Url: url,
	}
	link.GenerateHash()
	return link
}

func (link *Link) GenerateHash() {
	link.Hash = RandStringRunes(10)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func convertDeletedAt(d gorm.DeletedAt) *timestamppb.Timestamp {
	if d.Valid {
		return timestamppb.New(d.Time)
	}
	return nil
}
func ConvertToProtoLink(link *Link) *pb.Link {
	if link == nil {
		return nil
	}
	return &pb.Link{
		Model: &pb.Model{
			Id:        uint32(link.ID),
			CreatedAt: timestamppb.New(link.CreatedAt),
			UpdatedAt: timestamppb.New(link.UpdatedAt),
			DeletedAt: convertDeletedAt(link.DeletedAt),
		},
		Url:  link.Url,
		Hash: link.Hash,
	}
}
