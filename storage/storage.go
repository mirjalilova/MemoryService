package storage

import pb "github.com/mirjalilova/MemoryService/genproto/memory"

type StorageI interface {
	Comment() CommentI
	Media() MediaI
	Memory() MemoryI
	Share() ShareI
}

type CommentI interface {
	Create(*pb.CommentCreate) (*pb.Void, error)
	Update(*pb.GetById) (*pb.Void, error)
	Delete(*pb.GetById) (*pb.Void, error)
	Get(*pb.GetById) (*pb.Comment, error)
}

type MediaI interface {
	Create(*pb.MediaCreate) (*pb.Void, error)
	Delete(*pb.GetById) (*pb.Void, error)
	Get(*pb.GetById) (*pb.Media, error)
}

type MemoryI interface {
	Create(*pb.MemoryCreate) (*pb.Void, error)
	Update(*pb.MemoryUpdate) (*pb.Void, error)
	Delete(*pb.GetById) (*pb.Void, error)
	Get(*pb.GetById) (*pb.MemoryRes, error)
	GetAll(*pb.GetAllReq) (*pb.GetAllRes, error)
}

type ShareI interface {
	Share(*pb.ShareCreate) (*pb.Void, error)
	Unshare(*pb.ShareDelete) (*pb.Void, error)
	Update(*pb.GetById) (*pb.Void, error)
	Get(*pb.GetById) (*pb.ShareRes, error)
}


