package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/LuukBlankenstijn/fogistration/internal/http-server/models"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/config"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	dbModels "github.com/LuukBlankenstijn/fogistration/internal/shared/database/models"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type oidcService struct {
	cfg         config.OIDCConfig
	tokenIssuer TokenIssuer
	pool        *pgxpool.Pool

	once      sync.Once
	initErr   error
	verifier  *oidc.IDTokenVerifier
	oauth2Cfg *oauth2.Config
	endpoint  oauth2.Endpoint
}

func newOIDCService(cfg *config.OIDCConfig, ti TokenIssuer, pool *pgxpool.Pool) *oidcService {
	if cfg == nil {
		return nil
	}
	return &oidcService{
		cfg:         *cfg,
		tokenIssuer: ti,
		pool:        pool,
	}
}

func (s *oidcService) SetEndpoint(ep oauth2.Endpoint) {
	s.endpoint = ep
}

func (s *oidcService) ensureInit(ctx context.Context) error {
	s.once.Do(func() {
		p, err := oidc.NewProvider(ctx, s.cfg.IssuerURL)
		if err != nil {
			s.initErr = fmt.Errorf("oidc discovery: %w", err)
			logging.Error("OIDC discovery failed", err, "issuer", s.cfg.IssuerURL)
			return
		}

		if s.endpoint.AuthURL == "" || s.endpoint.TokenURL == "" {
			s.endpoint = p.Endpoint()
		}

		s.verifier = p.Verifier(&oidc.Config{ClientID: s.cfg.ClientID})

		s.oauth2Cfg = &oauth2.Config{
			ClientID:     s.cfg.ClientID,
			ClientSecret: s.cfg.ClientSecret,
			RedirectURL:  s.cfg.RedirectURL,
			Scopes:       s.cfg.Scopes,
			Endpoint:     s.endpoint,
		}
	})
	return s.initErr
}

func (s *oidcService) AuthURL(ctx context.Context, state, nonce string) (string, error) {
	if err := s.ensureInit(ctx); err != nil {
		return "", err
	}
	return s.oauth2Cfg.AuthCodeURL(state, oidc.Nonce(nonce)), nil
}

func (s *oidcService) HandleCallback(ctx context.Context, code, expectedNonce string) (AuthResult, error) {
	if err := s.ensureInit(ctx); err != nil {
		return AuthResult{}, err
	}

	tok, err := s.oauth2Cfg.Exchange(ctx, code)
	if err != nil {
		logging.Error("OIDC token exchange failed", err)
		return AuthResult{}, fmt.Errorf("exchange: %w", err)
	}

	rawID, ok := tok.Extra("id_token").(string)
	if !ok {
		err := fmt.Errorf("no id_token")
		logging.Error("OIDC id_token missing", err)
		return AuthResult{}, err
	}

	idt, err := s.verifier.Verify(ctx, rawID)
	if err != nil {
		logging.Error("OIDC ID token verify failed", err)
		return AuthResult{}, fmt.Errorf("verify: %w", err)
	}
	if idt.Nonce != expectedNonce {
		err := fmt.Errorf("bad nonce")
		logging.Error("OIDC nonce mismatch", err)
		return AuthResult{}, err
	}

	var claims struct {
		Sub               string `json:"sub"`
		Email             string `json:"email"`
		Name              string `json:"name"`
		PreferredUsername string `json:"preferred_username"`
	}
	if err := idt.Claims(&claims); err != nil {
		logging.Error("OIDC claims decode failed", err)
		return AuthResult{}, fmt.Errorf("claims: %w", err)
	}

	var user database.User
	err = withTx(ctx, s.pool, func(ctx context.Context, q *database.Queries) error {
		user, err = q.GetExternaluserById(ctx, database.PgTextFromString(&claims.Sub))
		if err != pgx.ErrNoRows {
			return err
		}

		user, err = q.CreateExternalUser(ctx, database.CreateExternalUserParams{
			Username:   claims.PreferredUsername,
			Email:      claims.Email,
			Role:       dbModels.User,
			ExternalID: database.PgTextFromString(&claims.Sub),
		})

		return nil
	})
	if err != nil {
		logging.Error("DB upsert external user failed", err)
		return AuthResult{}, err
	}

	jwt, err := s.tokenIssuer.IssueJWT(user)
	if err != nil {
		logging.Error("JWT issue failed", err)
		return AuthResult{}, err
	}

	return AuthResult{User: models.MapUsers(user)[0], Token: jwt}, nil
}
