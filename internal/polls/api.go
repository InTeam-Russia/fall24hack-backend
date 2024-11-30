package polls

import (
	"net/http"
	"strconv"

	"github.com/InTeam-Russia/go-backend-template/internal/apierr"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/session"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRoutes(
	r *gin.Engine,
	pollRepo Repo,
	sessionRepo session.Repo,
	logger *zap.Logger,
) {
	r.GET("/polls", func(c *gin.Context) {
		session, err := session.CheckHTTPReq(c, sessionRepo, logger)
		if err != nil {
			return
		}

		pageIndex := c.DefaultQuery("pageIndex", "0")
		pageSize := c.DefaultQuery("pageSize", "15")

		pageIndexInt, err := strconv.Atoi(pageIndex)
		if err != nil {
			c.JSON(http.StatusBadRequest, apierr.InvalidPageIndex)
			return
		}

		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			c.JSON(http.StatusBadRequest, apierr.InvalidPageSize)
			return
		}

		polls, err := pollRepo.GetUncompletedPolls(pageIndexInt, pageSizeInt, session.UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServer)
			logger.Error(err.Error())
			return
		}

		response := make([]OutModel, 0)
		for _, p := range polls {
			if p.Type != "RADIO" {
				p.Answers = nil
			}

			response = append(response, OutModel{
				Id:      p.Id,
				Text:    p.Text,
				Type:    p.Type,
				Answers: p.Answers,
			})
		}

		c.JSON(http.StatusOK, response)
	})
}
