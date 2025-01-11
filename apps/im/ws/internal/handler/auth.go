package handler

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/token"
	"im/apps/im/ws/internal/svc"
	"im/pkg/ctxdata"
	"net/http"
)

type JwtAuth struct {
	svc    *svc.ServiceContext
	parser *token.TokenParser
	logx.Logger
}

func NewJwtAuth(svc *svc.ServiceContext) *JwtAuth {
	return &JwtAuth{
		svc:    svc,
		parser: token.NewTokenParser(),
		Logger: logx.WithContext(context.Background()),
	}
}

func (j *JwtAuth) Auth(w http.ResponseWriter, r *http.Request) bool {
	//TODO implement me
	token, err := j.parser.ParseToken(r, j.svc.Config.JwtAuth.AccessSecret, "")
	if err != nil {
		fmt.Println(err)
		return false
	}

	if !token.Valid {
		fmt.Println("token is not valid!")
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	*r = *r.WithContext(context.WithValue(r.Context(), ctxdata.Identify, claims[ctxdata.Identify]))
	return true
}

func (j *JwtAuth) UserId(r *http.Request) string {
	//TODO implement me
	return ctxdata.GetUId(r.Context())
}
