package polls

import (
	"net/http"
	"strconv"

	"github.com/InTeam-Russia/go-backend-template/internal/apierr"
	"github.com/InTeam-Russia/go-backend-template/internal/auth/session"
	"github.com/InTeam-Russia/go-backend-template/internal/ml"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AnswerRequest struct {
	Answer string `json:"answer" binding:"required"`
}

func SetupRoutes(
	r *gin.Engine,
	pollRepo Repo,
	sessionRepo session.Repo,
	mlService ml.Service,
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

	r.POST("/polls/:id/answer", func(c *gin.Context) {
		session, err := session.CheckHTTPReq(c, sessionRepo, logger)
		if err != nil {
			return
		}

		var request AnswerRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, apierr.InvalidJSON)
			return
		}

		pollIdStr := c.Param("id")
		pollId, err := strconv.ParseInt(pollIdStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, apierr.InvalidID)
			return
		}

		err = pollRepo.AddAnswer(session.UserId, pollId, request.Answer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServer)
			logger.Error(err.Error())
			return
		}

		err = mlService.OnAnswer(session.UserId, request.Answer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServer)
			logger.Error(err.Error())
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "OK",
		})
	})

	r.POST("/polls", func(c *gin.Context) {
		session, err := session.CheckHTTPReq(c, sessionRepo, logger)
		if err != nil {
			return
		}

		var request CreateModel
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, apierr.InvalidJSON)
			return
		}

		if request.Type != FREE && request.Type != RADIO {
			c.JSON(http.StatusBadRequest, apierr.InvalidPollType)
			return
		}

		if request.Type == RADIO && len(request.Answers) == 0 {
			c.JSON(http.StatusBadRequest, apierr.NoRadioAnswers)
			return
		}

		cluster, err := mlService.OnQuestion(request.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServer)
			logger.Error(err.Error())
			return
		}

		err = pollRepo.CreatePoll(&request, session.UserId, cluster)
		if err != nil {
			c.JSON(http.StatusInternalServerError, apierr.InternalServer)
			logger.Error(err.Error())
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "OK",
		})
	})
}
