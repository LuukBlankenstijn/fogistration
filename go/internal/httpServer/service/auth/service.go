package auth

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/LuukBlankenstijn/fogistration/internal/httpServer/models"
	"github.com/LuukBlankenstijn/fogistration/internal/httpServer/utils/auth"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/database"
	"github.com/LuukBlankenstijn/fogistration/internal/shared/logging"
	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidCredentials = errors.New("invalid username or password")
var tokenIssuer = "fogistration"

type claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type JWTSigner struct {
	Secret []byte
	TTL    time.Duration
}

type service struct {
	q              *database.Queries
	jwt            *JWTSigner
	now            func() time.Time
	passwordParams auth.Params
}

func New(q *database.Queries, signer *JWTSigner) Service {
	return &service{
		q:              q,
		jwt:            signer,
		now:            time.Now,
		passwordParams: auth.Default,
	}
}

func (s *service) Authenticate(ctx context.Context, username, password string) (AuthResult, error) {
	fail := func() (AuthResult, error) { return AuthResult{}, ErrInvalidCredentials }

	u, err := s.q.GetUserByUsernameCI(ctx, username)
	if err != nil {
		return fail()
	}

	sec, err := s.q.GetAuthSecret(ctx, u.ID)
	if err != nil {
		return fail()
	}

	ok, err := auth.Verify(password, sec.Salt, sec.PasswordHash, s.passwordParams)
	if err != nil || !ok {
		return fail()
	}

	// best-effort; ignore error
	_, _ = s.q.TouchLastLogin(ctx, u.ID)

	token := ""
	if s.jwt != nil {
		t, err := s.issueJWT(u)
		if err != nil {
			return AuthResult{}, err
		}
		token = t
	}

	return AuthResult{User: models.MapUser(u), Token: token}, nil
}

func (s *service) DummyAuthenticate(ctx context.Context) (AuthResult, error) {
	u, err := s.q.GetUserByID(ctx, 1)
	if err != nil {
		return AuthResult{}, err
	}

	token := ""
	if s.jwt != nil {
		t, err := s.issueJWT(u)
		if err != nil {
			return AuthResult{}, err
		}
		token = t
	}

	return AuthResult{User: models.MapUser(u), Token: token}, nil
}

func (s *service) Validate(cookie *http.Cookie, sec string) (int64, error) {
	secret := []byte(sec)
	var cl claims
	tok, err := jwt.ParseWithClaims(
		cookie.Value, &cl,
		func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, jwt.ErrTokenUnverifiable
			}
			return secret, nil
		},
		jwt.WithIssuer(tokenIssuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithLeeway(30*time.Second), // tolerate small clock skew
	)
	if err != nil || !tok.Valid {
		// v5: no ValidationError. Check specific errors with errors.Is(...)
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			logging.Error("jwt malformed", err)
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			logging.Error("jwt bad signature", err)
		case errors.Is(err, jwt.ErrTokenExpired):
			logging.Error("jwt expired", err)
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			logging.Error("jwt not valid yet", err)
		case err != nil:
			logging.Error("jwt parse/verify error", err)
		default:
			logging.Error("jwt invalid", nil)
		}
		return 0, errors.New("invalid token")
	}
	// Extra issuer guard (defense in depth)
	if iss := cl.Issuer; !strings.EqualFold(iss, tokenIssuer) {
		return 0, errors.New("bad issuer")
	}

	// Extract subject (user id)
	var uid int64
	if cl.Subject != "" {
		if n, perr := strconv.ParseInt(cl.Subject, 10, 64); perr == nil {
			uid = n
		}
	}
	if uid == 0 {
		return 0, errors.New("invalid subject")
	}

	return uid, nil

}

func (s *service) issueJWT(u database.AppUser) (string, error) {
	claims := jwt.MapClaims{
		"sub":  strconv.Itoa(int(u.ID)),
		"name": u.Username,
		"role": u.Role,
		"iss":  tokenIssuer,
		"iat":  s.now().Unix(),
		"exp":  s.now().Add(s.jwt.TTL).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.jwt.Secret)
}
