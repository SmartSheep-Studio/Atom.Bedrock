package services

import (
	models2 "code.smartsheep.studio/atom/bedrock/pkg/server/datasource/models"
	"fmt"
	"github.com/IGLOU-EU/go-wildcard"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type AuthService struct {
	db    *gorm.DB
	users *UserService
}

func NewAuthService(db *gorm.DB, users *UserService) *AuthService {
	return &AuthService{db, users}
}

type AuthRequire2FAError struct {
	message string
}

func (err *AuthRequire2FAError) Error() string {
	return err.message
}

func (v *AuthService) AuthUser(id string, password string) (models2.User, error) {
	user, err := v.users.LookupUser(id)
	if err != nil {
		return user, fmt.Errorf("couldn't find user with %s", id)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return user, fmt.Errorf("invalid password")
	} else {
		return user, nil
	}
}

func (v *AuthService) NewSession(user models2.User, item *models2.UserSession) error {
	item.UserID = user.ID

	// TODO Add security check when log in at a new place(ip address)

	return v.db.Save(&item).Error
}

func (v *AuthService) NewJwt(session models2.UserSession, flag string, audience ...string) (models2.UserClaims, string, error) {
	var expires *jwt.NumericDate
	if flag == models2.UserClaimsTypeRefresh && session.ExpiredAt != nil {
		exp := session.ExpiredAt.Add(24 * 7 * time.Hour)
		expires = jwt.NewNumericDate(exp)
	} else if flag == models2.UserClaimsTypeAccess && session.ExpiredAt != nil {
		expires = jwt.NewNumericDate(*session.ExpiredAt)
	}

	audience = append(audience, viper.GetString("name"))
	claims := models2.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    viper.GetString("base_url"),
			Subject:   strconv.Itoa(int(session.UserID)),
			Audience:  audience,
			ExpiresAt: expires,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},

		Type:            flag,
		ClientID:        session.ClientID,
		SessionID:       session.ID,
		PersonalTokenID: lo.ToPtr(uint(0)),
	}

	session.Access = claims.ID

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString([]byte(viper.GetString("security.secret")))
	return claims, token, err
}

func (v *AuthService) ReadJwt(token string) (*models2.UserClaims, error) {
	res, err := jwt.ParseWithClaims(token, &models2.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(viper.GetString("security.secret")), nil
	})

	return res.Claims.(*models2.UserClaims), err
}

func (v *AuthService) ReadClaims(claims models2.UserClaims) (models2.UserSession, models2.User, error) {
	var session models2.UserSession
	if err := v.db.Where("id = ?", claims.SessionID).First(&session).Error; err != nil {
		return session, models2.User{}, fmt.Errorf("could not found session: #%d, because %s", claims.SessionID, err.Error())
	} else if session.ExpiredAt != nil && session.ExpiredAt.Unix() < time.Now().Unix() {
		return session, models2.User{}, fmt.Errorf("invalid session")
	}

	var user models2.User
	if err := v.db.Where("id = ?", claims.Subject).Preload("Contacts").Preload("Groups").First(&user).Error; err != nil {
		return session, user, fmt.Errorf("could not found user: #%s, because %s", claims.Subject, err.Error())
	}

	return session, user, nil
}

func (v *AuthService) HasUserPermissions(user models2.User, requires ...string) error {
	perms, err := user.GetPermissions()
	if err != nil {
		return err
	}

	for _, require := range requires {
		passed := false
		for key, val := range perms {
			if wildcard.Match(key, require) && (val != nil || val != false) {
				passed = true
				break
			}
		}

		if !passed {
			return fmt.Errorf("missing permission: %s", require)
		}
	}

	return nil
}

func (v *AuthService) HasSessionScope(session models2.UserSession, requires ...string) error {
	for _, require := range requires {
		passed := false
		for _, perm := range session.Scope {
			if wildcard.Match(perm, require) {
				passed = true
				break
			}
		}

		if !passed {
			return fmt.Errorf("missing scope: %s", requires)
		}
	}

	return nil
}
