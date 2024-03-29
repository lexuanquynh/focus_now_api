package database

import "context"

type UserRepository interface {
	// CreateUser Create  new user
	CreateUser(ctx context.Context, user *User) error
	// UpdateUserVerificationStatus Update user verification status
	UpdateUserVerificationStatus(ctx context.Context, email string, status bool) error
	// StoreVerificationData Save verification data into database
	StoreVerificationData(ctx context.Context, verificationData *VerificationData, isInsert bool) error
	// StoreProfileData Save profile data into database
	StoreProfileData(ctx context.Context, profileData *ProfileData) error
	// GetVerificationData Get verification data from database
	GetVerificationData(ctx context.Context, email string, verificationDataType VerificationDataType) (*VerificationData, error)
	//DeleteVerificationData Delete verification data from database
	DeleteVerificationData(ctx context.Context, email string, verificationDataType VerificationDataType) error
	// UpdateProfileData Update profile data into database
	UpdateProfileData(ctx context.Context, profileData *ProfileData) error
	// GetUserByEmail Get user by email
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	// GetUserByID Get user by id
	GetUserByID(ctx context.Context, id string) (*User, error)
	// UpdateUser Update user
	UpdateUser(ctx context.Context, user *User) error
	// GetProfileByID Get profile by user id
	GetProfileByID(ctx context.Context, userId string) (*ProfileData, error)
	// UpdateProfile Update profile
	UpdateProfile(ctx context.Context, profile *ProfileData) error
	// UpdatePassword Update password
	UpdatePassword(ctx context.Context, userID string, password string, tokenHash string) error
	// GetListOfPasswords Get list of passwords
	GetListOfPasswords(ctx context.Context, userID string) ([]string, error)
	// InsertListOfPasswords Update password into list of passwords
	InsertListOfPasswords(ctx context.Context, passwordUsers *PassworUsers) error
	// GetLimitData Get limit table data
	GetLimitData(ctx context.Context, userID string) (*LimitData, error)
	// InsertOrUpdateLimitData Insert or update limit data
	InsertOrUpdateLimitData(ctx context.Context, limitData *LimitData, isInsert bool) error
	// ClearAllLimitData Clear all limit data
	ClearAllLimitData(ctx context.Context) error
	// GetMultiRatioData Get multi ratio data
	GetMultiRatioData(ctx context.Context) (*MultiRatioData, error)
	// InsertEarnScore Insert earn score
	InsertEarnScore(ctx context.Context, earnScore *EarnScore) error
	// GetEarnScore Get earn score
	GetEarnScore(ctx context.Context, userID string) (*EarnScore, error)
}
