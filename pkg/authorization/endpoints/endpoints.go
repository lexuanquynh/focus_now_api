package endpoints

import (
	utils "LoveLetterProject/internal"
	"LoveLetterProject/internal/database"
	"LoveLetterProject/pkg/authorization"
	"LoveLetterProject/pkg/authorization/middleware"
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/hashicorp/go-hclog"
	"github.com/juju/ratelimit"
	"net/http"
	"strings"
)

type Set struct {
	HealthCheckEndpoint           endpoint.Endpoint
	RegisterEndpoint              endpoint.Endpoint
	VerifyMailEndpoint            endpoint.Endpoint
	LoginEndpoint                 endpoint.Endpoint
	LogoutEndpoint                endpoint.Endpoint
	GetUserEndpoint               endpoint.Endpoint
	GetProfileEndpoint            endpoint.Endpoint
	UpdateProfileEndpoint         endpoint.Endpoint
	UpdatePasswordEndpoint        endpoint.Endpoint
	GetForgetPasswordCodeEndpoint endpoint.Endpoint
	ResetPasswordEndpoint         endpoint.Endpoint
	GetMultiRatioDataEndpoint     endpoint.Endpoint
	InsertEarnScoreEndpoint       endpoint.Endpoint
	GetEarnScoreEndpoint          endpoint.Endpoint
	GenerateAccessTokenEndpoint   endpoint.Endpoint
}

func NewEndpointSet(svc authorization.Service,
	auth middleware.Authentication,
	r database.UserRepository,
	logger hclog.Logger,
	validator *database.Validation,
	tb *ratelimit.Bucket) Set {
	healthCheckEndpoint := MakeHealthCheckEndpoint(svc)
	healthCheckEndpoint = middleware.RateLimitRequest(tb, logger)(healthCheckEndpoint)

	registerEndpoint := MakeRegisterEndpoint(svc)
	registerEndpoint = middleware.ValidateParamRequest(validator, logger)(registerEndpoint)
	registerEndpoint = middleware.RateLimitRequest(tb, logger)(registerEndpoint)

	verifyMailEndpoint := MakeVerifyMailEndpoint(svc)
	verifyMailEndpoint = middleware.ValidateParamRequest(validator, logger)(verifyMailEndpoint)
	verifyMailEndpoint = middleware.RateLimitRequest(tb, logger)(verifyMailEndpoint)

	loginEndpoint := MakeLoginEndpoint(svc)
	loginEndpoint = middleware.RateLimitRequest(tb, logger)(loginEndpoint)
	loginEndpoint = middleware.ValidateParamRequest(validator, logger)(loginEndpoint)

	logoutEndpoint := MakeLogoutEndpoint(svc)
	logoutEndpoint = middleware.RateLimitRequest(tb, logger)(logoutEndpoint)
	logoutEndpoint = middleware.ValidateParamRequest(validator, logger)(logoutEndpoint)
	logoutEndpoint = middleware.ValidateRefreshToken(auth, r, logger)(logoutEndpoint)

	getUserEndpoint := MakeGetUserEndpoint(svc)
	getUserEndpoint = middleware.RateLimitRequest(tb, logger)(getUserEndpoint)
	getUserEndpoint = middleware.ValidateParamRequest(validator, logger)(getUserEndpoint)
	getUserEndpoint = middleware.ValidateAccessToken(auth, r, logger)(getUserEndpoint)

	getProfileEndpoint := MakeGetProfileEndpoint(svc)
	getProfileEndpoint = middleware.RateLimitRequest(tb, logger)(getProfileEndpoint)
	getProfileEndpoint = middleware.ValidateParamRequest(validator, logger)(getProfileEndpoint)
	getProfileEndpoint = middleware.ValidateAccessToken(auth, r, logger)(getProfileEndpoint)

	updateProfileEndpoint := MakeUpdateProfileEndpoint(svc)
	updateProfileEndpoint = middleware.RateLimitRequest(tb, logger)(updateProfileEndpoint)
	updateProfileEndpoint = middleware.ValidateParamRequest(validator, logger)(updateProfileEndpoint)
	updateProfileEndpoint = middleware.ValidateAccessToken(auth, r, logger)(updateProfileEndpoint)

	updatePasswordEndpoint := MakeUpdatePasswordEndpoint(svc)
	updatePasswordEndpoint = middleware.RateLimitRequest(tb, logger)(updatePasswordEndpoint)
	updatePasswordEndpoint = middleware.ValidateParamRequest(validator, logger)(updatePasswordEndpoint)
	updatePasswordEndpoint = middleware.ValidateAccessToken(auth, r, logger)(updatePasswordEndpoint)

	getForgetPasswordCodeEndpoint := MakeGetForgetPasswordCodeEndpoint(svc)
	getForgetPasswordCodeEndpoint = middleware.RateLimitRequest(tb, logger)(getForgetPasswordCodeEndpoint)
	getForgetPasswordCodeEndpoint = middleware.ValidateParamRequest(validator, logger)(getForgetPasswordCodeEndpoint)

	resetPasswordEndpoint := MakeCreateNewPasswordWithCodeEndpoint(svc)
	resetPasswordEndpoint = middleware.RateLimitRequest(tb, logger)(resetPasswordEndpoint)
	resetPasswordEndpoint = middleware.ValidateParamRequest(validator, logger)(resetPasswordEndpoint)

	getMultiRatioDataEndpoint := MakeGetMultiRatioDataEndpoint(svc)
	getMultiRatioDataEndpoint = middleware.RateLimitRequest(tb, logger)(getMultiRatioDataEndpoint)

	insertEarnScoreEndpoint := MakeInsertEarnScoreEndpoint(svc)
	insertEarnScoreEndpoint = middleware.RateLimitRequest(tb, logger)(insertEarnScoreEndpoint)
	insertEarnScoreEndpoint = middleware.ValidateParamRequest(validator, logger)(insertEarnScoreEndpoint)
	insertEarnScoreEndpoint = middleware.ValidateAccessToken(auth, r, logger)(insertEarnScoreEndpoint)

	getEarnScoreEndpoint := MakeGetEarnScoreEndpoint(svc)
	getEarnScoreEndpoint = middleware.RateLimitRequest(tb, logger)(getEarnScoreEndpoint)
	getEarnScoreEndpoint = middleware.ValidateParamRequest(validator, logger)(getEarnScoreEndpoint)
	getEarnScoreEndpoint = middleware.ValidateAccessToken(auth, r, logger)(getEarnScoreEndpoint)

	generateAccessTokenEndpoint := MakeGenerateAccessTokenEndpoint(svc)
	generateAccessTokenEndpoint = middleware.RateLimitRequest(tb, logger)(generateAccessTokenEndpoint)
	generateAccessTokenEndpoint = middleware.ValidateParamRequest(validator, logger)(generateAccessTokenEndpoint)
	generateAccessTokenEndpoint = middleware.ValidateRefreshToken(auth, r, logger)(generateAccessTokenEndpoint)

	return Set{
		HealthCheckEndpoint:           healthCheckEndpoint,
		RegisterEndpoint:              registerEndpoint,
		VerifyMailEndpoint:            verifyMailEndpoint,
		LoginEndpoint:                 loginEndpoint,
		LogoutEndpoint:                logoutEndpoint,
		GetUserEndpoint:               getUserEndpoint,
		GetProfileEndpoint:            getProfileEndpoint,
		UpdateProfileEndpoint:         updateProfileEndpoint,
		UpdatePasswordEndpoint:        updatePasswordEndpoint,
		GetForgetPasswordCodeEndpoint: getForgetPasswordCodeEndpoint,
		ResetPasswordEndpoint:         resetPasswordEndpoint,
		GetMultiRatioDataEndpoint:     getMultiRatioDataEndpoint,
		InsertEarnScoreEndpoint:       insertEarnScoreEndpoint,
		GetEarnScoreEndpoint:          getEarnScoreEndpoint,
		GenerateAccessTokenEndpoint:   generateAccessTokenEndpoint,
	}
}

// MakeRegisterEndpoint returns an endpoint that invokes Register on the service.
func MakeRegisterEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(authorization.RegisterRequest)
		if !ok {
			err := errors.New("invalid request")
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "invalid request")
			return nil, cusErr
		}
		message, err := svc.SignUp(ctx, &req)

		if err != nil {
			if strings.Contains(err.Error(), utils.PgDuplicateKeyMsg) {
				cusErr := utils.NewErrorWrapper(http.StatusConflict, err, "Tài khoản đã tồn tại. Vui lòng thử lại.")
				return nil, cusErr
			}
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "Không thể tạo tài khoản. Vui lòng thử lại.")
			return nil, cusErr
		}
		return message, err
	}
}

// MakeVerifyMailEndpoint returns an endpoint that invokes VerifyMail on the service.
func MakeVerifyMailEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(authorization.VerifyMailRequest)
		if !ok {
			err := errors.New("invalid request")
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "invalid request")
			return nil, cusErr
		}
		message, err := svc.VerifyMail(ctx, &req)

		if err != nil {
			return nil, err
		}
		return message, err
	}
}

// MakeLoginEndpoint returns an endpoint that invokes Login on the service.
func MakeLoginEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(authorization.LoginRequest)
		if !ok {
			err := errors.New("invalid request")
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "invalid request")
			return nil, cusErr
		}
		user, err := svc.Login(ctx, &req)

		if err != nil {
			if strings.Contains(err.Error(), utils.PgDuplicateKeyMsg) {
				cusErr := utils.NewErrorWrapper(http.StatusConflict, err, "Account already exists. Please try again.")
				return nil, cusErr
			}
			return nil, err
		}
		return user, err
	}
}

// MakeHealthCheckEndpoint returns an endpoint that invokes HealthCheck on the service.
func MakeHealthCheckEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err := svc.HealthCheck(ctx)
		return nil, err
	}
}

// MakeLogoutEndpoint returns an endpoint that invokes Logout on the service.
func MakeLogoutEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(authorization.LogoutRequest)
		if !ok {
			err := errors.New("invalid request")
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "invalid request")
			return nil, cusErr
		}
		err := svc.Logout(ctx, &req)

		if err != nil {
			cusErr := utils.NewErrorWrapper(http.StatusInternalServerError, err, "Can't logout. Please try again later.")
			return nil, cusErr
		}
		message := "Logout successfully"
		return message, nil
	}
}

// MakeGetUserEndpoint returns an endpoint that invokes GetUser on the service.
func MakeGetUserEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, ok := request.(authorization.GetUserRequest)
		if !ok {
			err := errors.New("invalid request")
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "invalid request")
			return nil, cusErr
		}
		user, err := svc.GetUser(ctx)
		if err != nil {
			cusErr := utils.NewErrorWrapper(http.StatusInternalServerError, err, "Can't get user. Please try again later.")
			return nil, cusErr
		}

		return user, err
	}
}

// MakeGetProfileEndpoint returns an endpoint that invokes GetProfile on the service.
func MakeGetProfileEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, ok := request.(authorization.GetProfileRequest)
		if !ok {
			err := errors.New("invalid request")
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "invalid request")
			return nil, cusErr
		}
		profile, err := svc.GetProfile(ctx)
		if err != nil {
			cusErr := utils.NewErrorWrapper(http.StatusInternalServerError, err, "Can't get profile. Please try again later.")
			return nil, cusErr
		}

		return profile, err
	}
}

// MakeUpdateProfileEndpoint returns an endpoint that invokes UpdateProfile on the service.
func MakeUpdateProfileEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(authorization.UpdateProfileRequest)
		if !ok {
			err := errors.New("invalid request")
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "invalid request")
			return nil, cusErr
		}
		profile, err := svc.UpdateProfile(ctx, &req)
		if err != nil {
			cusErr := utils.NewErrorWrapper(http.StatusInternalServerError, err, "Can't update profile. Please try again later.")
			return nil, cusErr
		}
		return profile, nil
	}
}

// MakeUpdatePasswordEndpoint returns an endpoint that invokes UpdatePassword on the service.
func MakeUpdatePasswordEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(authorization.UpdatePasswordRequest)
		if !ok {
			err := errors.New("invalid request")
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "invalid request")
			return nil, cusErr
		}
		message, err := svc.UpdatePassword(ctx, &req)
		if err != nil {
			cusErr := utils.NewErrorWrapper(http.StatusInternalServerError, err, message)
			return nil, cusErr
		}
		return message, nil
	}
}

// MakeGetForgetPasswordCodeEndpoint returns an endpoint that invokes GetForgetPasswordCode on the service.
func MakeGetForgetPasswordCodeEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(authorization.GetForgetPasswordCodeRequest)
		if !ok {
			err := errors.New("invalid request")
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "invalid request")
			return nil, cusErr
		}
		err := svc.GetForgetPasswordCode(ctx, req.Email)
		if err != nil {
			if strings.Contains(err.Error(), "successfully mailed password reset code") {
				return err.Error(), nil
			}
			return nil, err
		}
		return "successfully mailed password reset code. Please check your email.", nil
	}
}

// MakeCreateNewPasswordWithCodeEndpoint returns an endpoint that invokes ResetPassword on the service.
func MakeCreateNewPasswordWithCodeEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(authorization.CreateNewPasswordWithCodeRequest)
		if !ok {
			err := errors.New("invalid request")
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "invalid request")
			return nil, cusErr
		}
		err := svc.ResetPassword(ctx, &req)
		if err != nil {
			return nil, err
		}
		return "successfully updated password.", nil
	}
}

// MakeGetMultiRatioDataEndpoint returns an endpoint that invokes GetMultiRatioData on the service.
func MakeGetMultiRatioDataEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		ratioData, err := svc.GetMultiRatioData(ctx)
		if err != nil {
			return nil, err
		}
		return ratioData, nil
	}
}

// MakeInsertEarnScoreEndpoint returns an endpoint that invokes InsertEarnScore on the service.
func MakeInsertEarnScoreEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(authorization.InsertEarnScoreRequest)
		if !ok {
			err := errors.New("invalid request")
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "invalid request")
			return nil, cusErr
		}
		err := svc.InsertEarnScore(ctx, &req)
		if err != nil {
			return nil, err
		}
		return "successfully inserted earn score.", nil
	}
}

// MakeGetEarnScoreEndpoint returns an endpoint that invokes GetEarnScore on the service.
func MakeGetEarnScoreEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, ok := request.(authorization.GetEarnScoreRequest)
		if !ok {
			err := errors.New("invalid request")
			cusErr := utils.NewErrorWrapper(http.StatusBadRequest, err, "invalid request")
			return nil, cusErr
		}
		earnScore, err := svc.GetEarnScore(ctx)
		if err != nil {
			return nil, err
		}
		return earnScore, nil
	}
}

// MakeGenerateAccessTokenEndpoint returns an endpoint that invokes GenerateAccessToken on the service.
func MakeGenerateAccessTokenEndpoint(svc authorization.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, ok := request.(authorization.GenerateAccessTokenRequest)
		if !ok {
			cusErr := utils.NewErrorResponse(utils.BadRequest)
			return nil, cusErr
		}
		token, err := svc.GenerateAccessToken(ctx)
		if err != nil {
			return nil, err
		}
		return token, nil
	}
}
