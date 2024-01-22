package service

import (
	"github.com/LeftUnion/theName/service/TheName/internal/database"
	"github.com/LeftUnion/theName/service/TheName/internal/service/client"
)

type Service struct {
	client.IClient
}

func NewService(db database.IDatabase) IService {
	return &Service{
		IClient: client.New(db),
	}
}
