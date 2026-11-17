package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/mail"
	"net/smtp"
	"regexp"
	"strings"

	dbsqlc "github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/db/sqlc"
	"github.com/dinesh24murali/vibecode-framework-version1-e-commerce-codex/backend/internal/repositories"
)

// ErrContactInvalid is returned when the enquiry fails validation.
var ErrContactInvalid = errors.New("invalid contact enquiry")

var phoneRe = regexp.MustCompile(`^\+?[0-9][0-9 ()-]{5,18}$`)

// ContactInput is a contact form submission.
type ContactInput struct {
	Phone                         *string
	Name, Email, Subject, Message string
}

// ContactNotifier emails the business about a new enquiry.
type ContactNotifier interface {
	Notify(ctx context.Context, q dbsqlc.ContactQuery) error
}

// ContactService stores enquiries and notifies the business.
type ContactService struct {
	repo     repositories.ContactRepository
	notifier ContactNotifier
}

// NewContactService returns a ContactService. notifier may be nil (no email).
func NewContactService(repo repositories.ContactRepository, notifier ContactNotifier) *ContactService {
	return &ContactService{repo: repo, notifier: notifier}
}

// Submit validates and persists the enquiry, then sends the notification email.
// A notification failure is logged but does not fail the request: the enquiry is stored.
func (s *ContactService) Submit(ctx context.Context, in ContactInput) error {
	name, subject, message := strings.TrimSpace(in.Name), strings.TrimSpace(in.Subject), strings.TrimSpace(in.Message)
	email := strings.TrimSpace(in.Email)
	if name == "" || subject == "" || message == "" {
		return fmt.Errorf("%w: name, subject and message are required", ErrContactInvalid)
	}
	if _, err := mail.ParseAddress(email); err != nil || strings.ContainsAny(email, "\r\n") {
		return fmt.Errorf("%w: invalid email", ErrContactInvalid)
	}
	phone := sql.NullString{}
	if in.Phone != nil && strings.TrimSpace(*in.Phone) != "" {
		p := strings.TrimSpace(*in.Phone)
		if !phoneRe.MatchString(p) {
			return fmt.Errorf("%w: invalid phone", ErrContactInvalid)
		}
		phone = sql.NullString{String: p, Valid: true}
	}

	row, err := s.repo.Create(ctx, dbsqlc.CreateContactQueryParams{
		Name: name, Email: email, Phone: phone,
		Subject: sql.NullString{String: subject, Valid: true}, Message: message,
	})
	if err != nil {
		return err
	}
	if s.notifier != nil {
		if err := s.notifier.Notify(ctx, row); err != nil {
			slog.Error("contact: notification email failed", "id", row.ID, "error", err)
		}
	}
	return nil
}

// SMTPContactNotifier sends the enquiry to a fixed business address over SMTP
// (MailHog locally, SES SMTP in production).
type SMTPContactNotifier struct {
	Addr, From, To string
}

// header strips line breaks so user input cannot inject headers.
func header(s string) string {
	return strings.NewReplacer("\r", " ", "\n", " ").Replace(s)
}

// Notify implements ContactNotifier.
func (n SMTPContactNotifier) Notify(_ context.Context, q dbsqlc.ContactQuery) error {
	phone := "-"
	if q.Phone.Valid {
		phone = q.Phone.String
	}
	body := fmt.Sprintf("From: %s\r\nTo: %s\r\nReply-To: %s\r\nSubject: New enquiry: %s\r\n\r\nName: %s\r\nEmail: %s\r\nPhone: %s\r\n\r\n%s\r\n",
		header(n.From), header(n.To), header(q.Email), header(q.Subject.String), header(q.Name), header(q.Email), header(phone), q.Message)
	return smtp.SendMail(n.Addr, nil, n.From, []string{n.To}, []byte(body))
}
