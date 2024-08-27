package authcomp

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/imroc/req/v3"
	sctx "github.com/phathdt/service-context"
	"github.com/pkg/errors"
)

type AuthenComp interface {
	ValidateToken(ctx context.Context, token string) error
}

type authenComp struct {
	id       string
	endpoint string
	client   *req.Client
	logger   sctx.Logger
}
type authenResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (c *authenComp) ValidateToken(ctx context.Context, token string) error {
	var response authenResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		Get("/auth/valid")

	if err = json.Unmarshal(resp.Bytes(), &response); err != nil {
		return errors.WithStack(err)
	}

	if response.Code != 200 {
		return errors.WithStack(errors.New(response.Message))
	}

	return nil
}

func (c *authenComp) ID() string {
	return c.id
}

func (c *authenComp) InitFlags() {
	flag.StringVar(&c.endpoint, "authen-endpoint", "http://localhost:24000", "authen endpoint")
}

func (c *authenComp) Activate(sc sctx.ServiceContext) error {
	c.logger = sc.Logger(c.id)

	c.logger.Info("init authen comp")

	c.client = req.C().SetBaseURL(c.endpoint).DevMode()

	return nil
}

func (c *authenComp) Stop() error {
	return nil
}

func New(id string) *authenComp {
	return &authenComp{id: id}
}
