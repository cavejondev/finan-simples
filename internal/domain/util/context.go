package contextutil

import (
	"context"

	"github.com/google/uuid"
)

type key string

const (
	RequestIDKey key = "request_id"
	UserIDKey    key = "user_id"
	MethodKey    key = "method"
	PathKey      key = "path"
)

func SetRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, RequestIDKey, id)
}

func GetRequestID(ctx context.Context) *string {
	if v, ok := ctx.Value(RequestIDKey).(string); ok {
		return &v
	}
	return nil
}

func SetUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIDKey, id)
}

func GetUserID(ctx context.Context) *uuid.UUID {
	if v, ok := ctx.Value(UserIDKey).(uuid.UUID); ok {
		return &v
	}
	return nil
}

func SetMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, MethodKey, method)
}

func GetMethod(ctx context.Context) *Method {
	if v, ok := ctx.Value(MethodKey).(string); ok {
		m := Method(v)
		return &m
	}
	return nil
}

func SetPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, PathKey, path)
}

func GetPath(ctx context.Context) *string {
	if v, ok := ctx.Value(PathKey).(string); ok {
		return &v
	}
	return nil
}
