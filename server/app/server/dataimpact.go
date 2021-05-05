package server

import (
	"DataImpactProject/server/app/db"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

type DataImpactServer struct {
	//not really usefull but we can let it here
	Database *mgo.Database
	Router   *gin.Engine
}

// NewDataImpactServer init data impact server
func NewDataImpactServer() (*DataImpactServer, error) {
	session, err := db.GetMongoSession()
	if err != nil {
		return nil, err
	}
	return &DataImpactServer{Database: session.DB("dataImpact"), Router: gin.Default()}, nil
}
