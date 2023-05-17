package handler_http

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/otyang/go-starter/config"
	"github.com/otyang/go-starter/internal/middleware"
	"github.com/otyang/go-starter/internal/zample/entity"
	"github.com/otyang/go-starter/pkg/datastore"
	"github.com/otyang/go-starter/pkg/logger"
	"github.com/otyang/go-starter/pkg/response"
	"github.com/uptrace/bun"
)

type Handler struct {
	Log    logger.Interface
	Config *config.Config
	Repo   entity.IRepository
}

func NewHandler(repo entity.IRepository, config *config.Config, Log logger.Interface) *Handler {
	return &Handler{
		Log:    Log,
		Config: config,
		Repo:   repo,
	}
}

func (h *Handler) Welcome(c *fiber.Ctx) error {
	resp := response.Ok("", "Hello, Zample Welcome!")
	return c.Status(resp.StatusCode).JSON(resp)
}

func (h *Handler) Profile(c *fiber.Ctx) error {
	body, ok := middleware.ValidatedDataFromContext(c).(*entity.ProfileRequest)

	if !ok {
		resp := response.InternalServerError("error from context", "")
		return c.Status(resp.StatusCode).JSON(resp)
	}

	resp := response.Ok("", body)
	return c.Status(resp.StatusCode).JSON(resp)
}

func (h *Handler) Profile2(c *fiber.Ctx) error {
	body := new(entity.ProfileRequest)

	err := middleware.ParseAndValidate(c, body)
	if err != nil {
		resp := response.BadRequest(err.Error(), "")
		return c.Status(resp.StatusCode).JSON(resp)
	}

	resp := response.Ok("", body)
	return c.Status(resp.StatusCode).JSON(resp)
}

func (h *Handler) ListViaDbHelpers(db *bun.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users := []*entity.User{}
		con := datastore.NewDBHelper(db)
		err := con.ListAll(context.TODO(), &users)
		if err != nil {
			_r := response.InternalServerError(err.Error(), "")
			return c.Status(_r.StatusCode).JSON(_r)
		}

		rsp := response.Ok("", users)
		return c.Status(200).JSON(rsp)
	}
}

func (h *Handler) ListViaRepository(c *fiber.Ctx) error {
	users, err := h.Repo.GetUsers(c.Context())
	if err != nil {
		_r := response.InternalServerError(err.Error(), "")
		return c.Status(_r.StatusCode).JSON(_r)
	}

	rsp := response.Ok("", users)
	return c.Status(200).JSON(rsp)
}

func (h *Handler) Error500(c *fiber.Ctx) error {
	err := errors.New("deliberate error for testing")
	_r := response.InternalServerError(err.Error(), "")
	return c.Status(_r.StatusCode).JSON(_r)
}

// Error returned this way in a handler are resolved
// in the errorHandler middleware
func (h *Handler) ErrorReturned(c *fiber.Ctx) error {
	err := errors.New("deliberate error for testing")
	return response.InternalServerError(err.Error(), "")
}
