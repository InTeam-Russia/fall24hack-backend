package recommendations

import (
	"net/http"
	"strconv"

	"github.com/InTeam-Russia/go-backend-template/internal/apierr"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/session"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/user"
	"github.com/InTeam-Russia/go-backend-template/internal/ml"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/exp/constraints"
)

func SetupRoutes(
	r *gin.Engine,
	sessionRepo session.Repo,
	mlService ml.Service,
	userRepo user.Repo,
	logger *zap.Logger,
) {
	r.GET("/users", func(c *gin.Context) {
		pageIndex, err := strconv.Atoi(c.DefaultQuery("pageIndex", "0"))
		if err != nil {
			c.JSON(http.StatusBadRequest, apierr.InvalidPageIndex)
			return
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "5"))
		if err != nil {
			c.JSON(http.StatusBadRequest, apierr.InvalidPageSize)
			return
		}

		searchType := c.DefaultQuery("searchType", string(ml.CODIRECTIONAL))
		if searchType != string(ml.CODIRECTIONAL) && searchType != string(ml.OPPOSITE) {
			c.JSON(http.StatusBadRequest, apierr.InvalidPageSize)
			return
		}

		session, err := session.CheckHTTPReq(c, sessionRepo, logger)
		if err != nil {
			return
		}

		mlUsers, err := mlService.UsersANN(session.UserId, (pageIndex+1)*pageSize, ml.SearchType(searchType))
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServer)
			logger.Error(err.Error())
			return
		}

		ids := make([]int64, 0, len(mlUsers))
		for _, mlu := range mlUsers {
			ids = append(ids, mlu.Id)
		}

		mlUsers = mlUsers[pageIndex*pageSize:]

		dbUsers, err := userRepo.GetByIds(ids)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServer)
			logger.Error(err.Error())
			return
		}

		users := make([]UserResponse, 0, len(dbUsers))
		for i := 0; i < min(len(dbUsers), len(mlUsers)); i++ {
			mlu := mlUsers[i]
			dbu := dbUsers[i]

			if mlu.Id == session.UserId {
				continue
			}

			user := UserResponse{
				Id:                    dbu.Id,
				FirstName:             dbu.FirstName,
				LastName:              dbu.LastName,
				Username:              dbu.Username,
				Email:                 dbu.Email,
				TgLink:                dbu.TgLink,
				OverlappingPercentage: mlu.OverlappingPercentage,
			}

			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	})
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
