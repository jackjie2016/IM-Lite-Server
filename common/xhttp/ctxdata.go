package xhttp

import "context"

// CtxKeyUserId get uid from ctx
var CtxKeyUserId = "userId"

func GetUidFromCtx(ctx context.Context) string {
	var uid string
	if jsonUid, ok := ctx.Value(CtxKeyUserId).(string); ok {
		uid = jsonUid
	}
	return uid
}

// CtxKeyPlatform get platform from ctx
var CtxKeyPlatform = "platform"

func GetPlatformFromCtx(ctx context.Context) string {
	var platform string
	if jsonPlatform, ok := ctx.Value(CtxKeyPlatform).(string); ok {
		platform = jsonPlatform
	}
	return platform
}

// CtxKeyAppVersion get appVersion from ctx
var CtxKeyAppVersion = "appVersion"

func GetAppVersionFromCtx(ctx context.Context) string {
	var appVersion string
	if jsonAppVersion, ok := ctx.Value(CtxKeyAppVersion).(string); ok {
		appVersion = jsonAppVersion
	}
	return appVersion
}

// CtxKeyUserAgent get appVersion from ctx
var CtxKeyUserAgent = "userAgent"

func GetUserAgentFromCtx(ctx context.Context) string {
	var userAgent string
	if jsonUserAgent, ok := ctx.Value(CtxKeyUserAgent).(string); ok {
		userAgent = jsonUserAgent
	}
	return userAgent
}

// CtxKeyIp get appVersion from ctx
var CtxKeyIp = "ip"

func GetIpFromCtx(ctx context.Context) string {
	var ip string
	if jsonIp, ok := ctx.Value(CtxKeyIp).(string); ok {
		ip = jsonIp
	}
	return ip
}

// CtxKeyToken get appVersion from ctx
var CtxKeyToken = "token"

func GetTokenFromCtx(ctx context.Context) string {
	var token string
	if jsonToken, ok := ctx.Value(CtxKeyToken).(string); ok {
		token = jsonToken
	}
	return token
}
