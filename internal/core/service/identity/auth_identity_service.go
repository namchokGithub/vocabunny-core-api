package identity

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
	"golang.org/x/crypto/bcrypt"
)

type AuthIdentityServiceDependencies struct {
	AuthIdentityRepository port.AuthIdentityRepository
	UserRepository         port.UserRepository
	TxManager              port.TransactionManager
	TokenManager           port.TokenManager
}

type authIdentityService struct {
	authIdentityRepository port.AuthIdentityRepository
	userRepository         port.UserRepository
	txManager              port.TransactionManager
	tokenManager           port.TokenManager
}

func NewAuthIdentityService(deps AuthIdentityServiceDependencies) port.AuthIdentityService {
	return &authIdentityService{
		authIdentityRepository: deps.AuthIdentityRepository,
		userRepository:         deps.UserRepository,
		txManager:              deps.TxManager,
		tokenManager:           deps.TokenManager,
	}
}

func (s *authIdentityService) Create(ctx context.Context, input domain.AuthIdentityCreateInput) (domain.AuthIdentity, error) {
	input.ProviderUserID = strings.TrimSpace(input.ProviderUserID)
	input.Password = strings.TrimSpace(input.Password)

	if _, err := s.userRepository.FindByID(ctx, input.UserID); err != nil {
		return domain.AuthIdentity{}, err
	}

	if !isValidAuthProvider(input.Provider) {
		return domain.AuthIdentity{}, helper.BadRequest("invalid_provider", "provider is invalid", nil)
	}

	var created domain.AuthIdentity
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.validateIdentityInput(txCtx, input, nil); err != nil {
			return err
		}

		passwordHash, err := hashPasswordIfNeeded(input.Provider, input.Password)
		if err != nil {
			return err
		}

		identity, err := s.authIdentityRepository.Create(txCtx, domain.AuthIdentity{
			UserID:         input.UserID,
			Provider:       input.Provider,
			ProviderUserID: input.ProviderUserID,
			PasswordHash:   passwordHash,
			AuditFields: domain.AuditFields{
				CreatedBy: input.ActorID,
				UpdatedBy: input.ActorID,
			},
		})
		if err != nil {
			return err
		}

		created = identity
		return nil
	})
	if err != nil {
		return domain.AuthIdentity{}, err
	}

	return created, nil
}

func (s *authIdentityService) Update(ctx context.Context, input domain.AuthIdentityUpdateInput) (domain.AuthIdentity, error) {
	var updated domain.AuthIdentity
	err := s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		current, err := s.authIdentityRepository.FindByID(txCtx, input.ID)
		if err != nil {
			return err
		}

		if input.ProviderUserID.Set {
			input.ProviderUserID.Value = strings.TrimSpace(input.ProviderUserID.Value)
			if current.Provider != domain.AuthProviderPassword && input.ProviderUserID.Value == "" {
				return helper.BadRequest("invalid_provider_user_id", "provider_user_id is required", nil)
			}

			if current.Provider != domain.AuthProviderPassword {
				exists, err := s.authIdentityRepository.ExistsByProviderIdentity(txCtx, current.Provider, input.ProviderUserID.Value, &input.ID)
				if err != nil {
					return err
				}
				if exists {
					return helper.Conflict("duplicate_provider_identity", "provider identity already exists", nil)
				}
			}
		}

		if input.Password.Set {
			password := strings.TrimSpace(input.Password.Value)
			if current.Provider != domain.AuthProviderPassword {
				return helper.BadRequest("invalid_password_update", "password can only be updated for PASSWORD provider", nil)
			}
			if len(password) < 8 {
				return helper.BadRequest("invalid_password", "password must be at least 8 characters", nil)
			}
			input.Password.Value = password
		}

		updated, err = s.authIdentityRepository.Update(txCtx, input)
		return err
	})
	if err != nil {
		return domain.AuthIdentity{}, err
	}

	return updated, nil
}

func (s *authIdentityService) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	return s.txManager.RunInTx(ctx, func(txCtx context.Context) error {
		if _, err := s.authIdentityRepository.FindByID(txCtx, id); err != nil {
			return err
		}
		return s.authIdentityRepository.Delete(txCtx, id, actorID)
	})
}

func (s *authIdentityService) FindByID(ctx context.Context, id uuid.UUID) (domain.AuthIdentity, error) {
	return s.authIdentityRepository.FindByID(ctx, id)
}

func (s *authIdentityService) FindAll(ctx context.Context, query domain.AuthIdentityQuery) (domain.PageResult[domain.AuthIdentity], error) {
	return s.authIdentityRepository.FindAll(ctx, query)
}

func (s *authIdentityService) LoginWithPassword(ctx context.Context, input domain.PasswordLoginInput) (domain.AuthToken, error) {
	login := strings.TrimSpace(input.EmailOrUsername)
	password := strings.TrimSpace(input.Password)
	if login == "" || password == "" {
		return domain.AuthToken{}, helper.BadRequest("invalid_login", "email_or_username and password are required", nil)
	}

	identity, err := s.authIdentityRepository.FindPasswordIdentityByLogin(ctx, login)
	if err != nil {
		return domain.AuthToken{}, helper.Unauthorized("invalid_credentials", "invalid credentials", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(identity.PasswordHash), []byte(password)); err != nil {
		return domain.AuthToken{}, helper.Unauthorized("invalid_credentials", "invalid credentials", err)
	}

	user, err := s.userRepository.FindByID(ctx, identity.UserID)
	if err != nil {
		return domain.AuthToken{}, err
	}

	if user.Status != domain.UserStatusActive {
		return domain.AuthToken{}, helper.Unauthorized("user_inactive", "user is not active", nil)
	}

	token, err := s.tokenManager.GenerateAccessToken(user.ID.String())
	if err != nil {
		return domain.AuthToken{}, helper.Internal("generate_token_failed", "failed to generate access token", err)
	}

	return domain.AuthToken{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.tokenManager.AccessTokenTTLSeconds(),
		User:        user,
	}, nil
}

func (s *authIdentityService) validateIdentityInput(ctx context.Context, input domain.AuthIdentityCreateInput, excludeID *uuid.UUID) error {
	switch input.Provider {
	case domain.AuthProviderPassword:
		if len(input.Password) < 8 {
			return helper.BadRequest("invalid_password", "password must be at least 8 characters", nil)
		}
	default:
		if input.ProviderUserID == "" {
			return helper.BadRequest("invalid_provider_user_id", "provider_user_id is required", nil)
		}
		exists, err := s.authIdentityRepository.ExistsByProviderIdentity(ctx, input.Provider, input.ProviderUserID, excludeID)
		if err != nil {
			return err
		}
		if exists {
			return helper.Conflict("duplicate_provider_identity", "provider identity already exists", nil)
		}
	}
	return nil
}

func hashPasswordIfNeeded(provider domain.AuthProvider, password string) (string, error) {
	if provider != domain.AuthProviderPassword {
		return "", nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", helper.Internal("hash_password_failed", "failed to hash password", err)
	}

	return string(hash), nil
}

func isValidAuthProvider(provider domain.AuthProvider) bool {
	switch provider {
	case domain.AuthProviderPassword, domain.AuthProviderGoogle, domain.AuthProviderApple:
		return true
	default:
		return false
	}
}
