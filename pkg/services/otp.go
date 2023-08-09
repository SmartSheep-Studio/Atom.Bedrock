package services

import (
	"errors"
	"fmt"
	"github.com/samber/lo"
	"strconv"
	"strings"
	"time"

	"code.smartsheep.studio/atom/bedrock/pkg/datasource/models"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OTPService struct {
	db     *gorm.DB
	mailer *MailerService
}

func NewOTPService(db *gorm.DB, mailer *MailerService) *OTPService {
	return &OTPService{db, mailer}
}

func (v *OTPService) LookupOTP(user models.User, code string) (models.OTP, error) {
	var otp models.OTP
	if err := v.db.Where("code = ? AND user_id = ?", code, user.ID).First(&otp).Error; err != nil {
		return otp, err
	} else if otp.ExpiredAt != nil && time.Now().Unix() >= otp.ExpiredAt.Unix() {
		return otp, fmt.Errorf("the OTP has been expired")
	} else {
		return otp, nil
	}
}

func (v *OTPService) NewOTP(user models.User, method int, payload models.OTPPayload, expires *time.Duration) (models.OTP, error) {
	var otp models.OTP
	if err := v.db.Where("user_id = ? AND type = ?", user.ID, method).First(&otp).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return otp, err
		} else {
			code := strings.ToUpper(uuid.NewString()[:6])
			otp = models.OTP{
				Type:        method,
				Code:        code,
				Payload:     datatypes.NewJSONType(payload),
				RefreshedAt: nil,
				UserID:      user.ID,
			}
			if expires != nil {
				expired := time.Now().Add(*expires)
				otp.ExpiredAt = &expired
			}
		}
	} else {
		if otp.ExpiredAt != nil && time.Now().Unix() >= otp.ExpiredAt.Unix() {
			now := time.Now()
			code := strings.ToUpper(uuid.NewString()[:6])
			otp.RefreshedAt = &now
			otp.Code = code
		} else {
			return otp, fmt.Errorf("find a same type and isn't expired one time passcode, please wait a few minute and try again!")
		}
	}

	if err := v.db.Save(&otp).Error; err != nil {
		return otp, err
	} else {
		return otp, nil
	}
}

func (v *OTPService) ApplyOTP(otp models.OTP) error {
	var user models.User
	if err := v.db.Where("id = ?", otp.UserID).Preload("Groups").First(&user).Error; err != nil {
		return err
	}

	switch otp.Type {
	case models.OneTimeVerifyContactCode:
		id, _ := strconv.Atoi(otp.Payload.Data().Target)
		var contact models.UserContact
		if err := v.db.Where("id = ? AND user_id = ?", id, otp.UserID).First(&contact).Error; err != nil {
			return err
		} else if contact.Type != models.UserContactTypeEmail {
			return fmt.Errorf("couldn't send mail to verify contact type isn't email")
		}
		if user.VerifiedAt == nil || len(user.Groups) == 0 {
			if len(user.Groups) == 0 {
				var group models.UserGroup
				if err := v.db.Where("slug = ?", "verified_users").First(&group).Error; err != nil {
					return err
				}
				user.Groups = append(user.Groups, group)
			}
			if user.VerifiedAt == nil {
				user.VerifiedAt = lo.ToPtr(time.Now())
			}
			if err := v.db.Save(&user).Error; err != nil {
				return err
			}
		}
		contact.VerifiedAt = lo.ToPtr(time.Now())
		return v.db.Save(&contact).Error
	default:
		return fmt.Errorf("unsupported OTP type for applying")
	}
}

func (v *OTPService) SendMail(otp models.OTP) error {
	var user models.User
	if err := v.db.Where("id = ?", otp.UserID).First(&user).Error; err != nil {
		return err
	}

	switch otp.Type {
	case models.OneTimeVerifyContactCode:
		id, _ := strconv.Atoi(otp.Payload.Data().Target)
		var contact models.UserContact
		if err := v.db.Where("id = ? AND user_id = ?", id, otp.UserID).First(&contact).Error; err != nil {
			return err
		} else if contact.Type != models.UserContactTypeEmail {
			return fmt.Errorf("couldn't send mail to verify contact type isn't email")
		}
		return v.mailer.SendMail(contact.Content, "[Atom ID] Verify your email", fmt.Sprintf(`Hello, %s!
You have just initiated a verification email request in the Atom user center, here is your email verification one-time password: %s

Please note that the verification code will expire on %s, and you cannot initiate any email verification requests again before the verification code expires or is used!

This message is automated from Atom ID, please do not reply to this message, you will not get a reply.
If you did not initiate any verification, please ignore this email.

Request details:
One Time Password ID: #%d
Requesting User ID: #%d
Requester IP: %s`, user.Nickname, otp.Code, otp.ExpiredAt.UTC(), otp.ID, user.ID, otp.Payload.Data().IpAddress))
	case models.OneTimeDangerousPassCode:
		var contacts []models.UserContact
		var contact models.UserContact
		if err := v.db.Where("user_id = ?", otp.UserID).Find(&contacts).Error; err != nil {
			return err
		} else {
			for _, item := range contacts {
				if item.Type == models.UserContactTypeEmail {
					contact = item
					break
				}
			}
		}
		return v.mailer.SendMail(contact.Content, "[Atom ID] Verify that is you", fmt.Sprintf(`Hello, %s!
You were just performing dangerous operation %s, which was automatically blocked by Atom ID security center, if you want to continue this operation, please enter the following one-time password: %s

Please note that this one-time password will expire on %s, and you cannot initiate any dangerous operation requests again until the verification code expires or is used!

This message was sent automatically by Atom ID, please do not reply to this message, you will not get a reply.
If you have not initiated a dangerous operation, please check your account for suspicious sessions, log out of the session in time, and change your account password! Because your account may be controlled by others!

Request details:
One Time Password ID: #%d
Requesting User ID: #%d
Requester IP: %s`, user.Nickname, otp.Payload.Data().Target, otp.Code, otp.ExpiredAt.UTC(), otp.ID, user.ID, otp.Payload.Data().IpAddress))
	default:
		return fmt.Errorf("unsupported OTP type for sending email")
	}
}
