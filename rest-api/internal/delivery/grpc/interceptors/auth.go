package interceptors

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
	"workshops/rest-api/internal/config"
	"workshops/rest-api/internal/delivery/http"
	"workshops/rest-api/internal/entities"
)

type AuthMD struct {
	JwtWrapper *entities.JwtWrapper
	Service    http.AuthService
}

func NewAuthMD(jwtWrapper *entities.JwtWrapper, service http.AuthService) AuthMD {
	return AuthMD{
		JwtWrapper: jwtWrapper,
		Service:    service,
	}
}

const (
	authHeader = "authorization"
	bearerAuth = "bearer"
)

func (a *AuthMD) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		clientToken := a.getAuthCredentials(ctx)
		if clientToken == "" {
			return nil, status.Error(codes.Unauthenticated, "unauthenticated")
		}

		claims, err := a.JwtWrapper.ValidateToken(clientToken)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "can't parse credentials")
		}

		c, ok := claims.(*entities.JwtClaim)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "can't parse credentials")
		}

		var user *entities.User
		if user, err = a.Service.FindByToken(ctx, clientToken); err != nil || user.Email != c.Email {
			return nil, status.Error(codes.Unauthenticated, "Token was canceled")
		}

		ctx = context.WithValue(ctx, config.CtxAuthId, user.Id)
		ctx = context.WithValue(ctx, config.CtxToken, clientToken)

		return handler(ctx, req)
	}
}

func (a *AuthMD) getAuthCredentials(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := md.Get(authHeader)
	if len(values) == 0 {
		return ""
	}

	fields := strings.SplitN(values[0], " ", 2)
	if len(fields) < 2 {
		return ""
	}

	if !strings.EqualFold(fields[0], bearerAuth) {
		return ""
	}
	return fields[1]
}
