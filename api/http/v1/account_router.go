package v1

import (
	"context"
	"net/http"
	"starland-account/configs"
	"starland-account/internal/pkg/util"
	"starland-account/internal/service/account"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AccountHTTPServer interface {
	Auth(context.Context, *account.AccountRequest) error
	QueryAccount(context.Context, string) (*account.AccountResponse, error)
	ClaimPoints(context.Context, *account.ClaimPointsRequest) (string, error)
	SavePointsAddr(context.Context, string, string) error
}

func InitAccountRouter(app fiber.Router, service AccountHTTPServer, conf *configs.Config) {
	router := app.Group("v1")
	router.Post("/account", auth(service))
	router.Post("/account/claim_points", claimPoints(service))
	router.Get("/account/:id", queryAccounts(service))
	router.Post("/account/:id/save_points_addr", savePointsAddr(service))
}

func auth(service AccountHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				AccountID string `json:"account_id"`
				Email     string `json:"email"`
				Name      string `json:"name"`
				Provider  string `json:"provider"`
				AvatarURL string `json:"avatar_url"`
			}
		)

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		act := &account.AccountRequest{
			AccountID: req.AccountID,
			Email:     req.Email,
			Name:      req.Email,
			Provider:  req.Provider,
			AvatarURL: req.AvatarURL,
		}

		if err := service.Auth(ctx.Context(), act); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse("ok"))
	}
}

func queryAccounts(service AccountHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				ID string `params:"id"`
			}
		)

		if err := ctx.ParamsParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		response, err := service.QueryAccount(ctx.Context(), req.ID)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(response))
	}
}

func claimPoints(service AccountHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				AccountID string `json:"account_id"`
				Points    int    `json:"points"`
				IsOK      bool   `json:"is_ok"`
			}
		)

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		zap.S().Info("req:", req.AccountID)
		cpr := &account.ClaimPointsRequest{
			AccountID: req.AccountID,
			Points:    req.Points,
			IsOK:      req.IsOK,
		}

		res, err := service.ClaimPoints(ctx.Context(), cpr)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse(res))
	}
}

func savePointsAddr(service AccountHTTPServer) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var (
			req struct {
				Addr    string `json:"addr"`
				Account string `json:"account"`
			}
		)

		if err := ctx.BodyParser(&req); err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}

		err := service.SavePointsAddr(ctx.Context(), req.Account, req.Addr)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(util.MakeResponseWithMsg(err.Error()))
		}
		return ctx.Status(http.StatusOK).JSON(util.MakeResponse("ok"))
	}
}
