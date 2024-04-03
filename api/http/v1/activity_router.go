package v1

import (
	"context"
	"errors"
	"net/http"
	"starland-account/configs"
	"starland-account/internal/pkg/util"
	"starland-account/internal/service/activity"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ActivityHTTPServer interface {
	Play(context.Context, int, string) error
	QueryActivityLogs(context.Context, string, int, int) ([]*activity.ActivityLogResponse, int64, error)
	QueryActivitys(ctx context.Context) ([]*activity.ActivityResponse, error)
	QueryIsLimit(context.Context, int, string) (bool, error)
}

func InitActivityRouter(app fiber.Router, service ActivityHTTPServer, conf *configs.Config) {
	router := app.Group("v1")
	router.Post("/activity", play(service))
	router.Get("/activity/Limit", queryIsLimit(service))
	router.Get("/activity", queryActivitys(service))
	router.Get("/activity/log/:account", queryActivityLogs(service))
}

func play(service ActivityHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				ActivityCode int    `json:"activity_code"`
				Account      string `json:"account"`
			}
		)

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		if err := service.Play(ctx.Context(), req.ActivityCode, req.Account); err != nil {
			if errors.Is(err, activity.LimitError) {
				return ctx.Status(http.StatusOK).JSON(util.MakeResponse(err.Error()).SetCode("100"))
			}
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse("ok"))
	}
}

func queryActivityLogs(service ActivityHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				Account string `params:"account"`
				Page    int    `query:"page"`
				Limit   int    `query:"limit"`
			}
			res struct {
				Data  []*activity.ActivityLogResponse `json:"data"`
				Count int64                           `json:"count"`
			}
		)

		if err := ctx.ParamsParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		if err := ctx.QueryParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		zap.S().Infof("queryActivityLogs: req: %+v", req)
		response, count, err := service.QueryActivityLogs(ctx.Context(), req.Account, req.Page, req.Limit)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		res.Count = count
		res.Data = response
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}

func queryActivitys(service ActivityHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			res struct {
				Data []*activity.ActivityResponse `json:"data"`
			}
		)
		response, err := service.QueryActivitys(ctx.Context())
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		res.Data = response
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}

func queryIsLimit(service ActivityHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				ActivityCode int    `query:"activity_code"`
				Account      string `query:"account"`
			}

			res struct {
				IsLimit bool `json:"is_limit"`
			}
		)

		if err := ctx.QueryParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		b, err := service.QueryIsLimit(ctx.Context(), req.ActivityCode, req.Account)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		res.IsLimit = b
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}
