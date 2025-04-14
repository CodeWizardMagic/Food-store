package usecase

import (
	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/repository"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	UserRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{UserRepo: userRepo}
}

func (uc *UserUseCase) Register(user domain.User) (string, error) {
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return "", errors.New("all fields are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	userID, err := uc.UserRepo.Save(user)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (uc *UserUseCase) Authenticate(username, password string) (string, string, error) {
	user, err := uc.UserRepo.FindByUsername(username)
	if err != nil {
		return "", "", errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid username or password")
	}

	token := uuid.New().String()
	tokenData := domain.Token{
		UserID:    user.ID,
		Token:     token,
		CreatedAt: time.Now(),
	}

	if err := uc.UserRepo.SaveToken(tokenData); err != nil {
		return "", "", err
	}

	return token, user.ID, nil
}

func (uc *UserUseCase) GetProfile(userID string) (domain.User, error) {
	user, err := uc.UserRepo.FindByID(userID)
	if err != nil {
		return domain.User{}, err
	}

	user.Password = ""
	return user, nil
}

func (uc *UserUseCase) ValidateToken(token string) (string, error) {
	userID, err := uc.UserRepo.FindUserIDByToken(token)
	if err != nil {
		return "", errors.New("invalid token")
	}
	return userID, nil
}