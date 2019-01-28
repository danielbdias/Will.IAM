package repositories

import (
	"fmt"

	"github.com/ghostec/Will.IAM/models"
)

// ServiceAccounts repository
type ServiceAccounts interface {
	Get(string) (*models.ServiceAccount, error)
	List() ([]models.ServiceAccount, error)
	Search(string) ([]models.ServiceAccount, error)
	ForEmail(string) (*models.ServiceAccount, error)
	ForKeyPair(string, string) (*models.ServiceAccount, error)
	Create(*models.ServiceAccount) error
}

type serviceAccounts struct {
	storage *Storage
}

func (sas serviceAccounts) Get(id string) (*models.ServiceAccount, error) {
	sa := new(models.ServiceAccount)
	if _, err := sas.storage.PG.DB.Query(
		sa,
		`SELECT id, name, key_id, key_secret, email, base_role_id
		FROM service_accounts
		WHERE id = ?`,
		id,
	); err != nil {
		return nil, err
	}
	return sa, nil
}

func (sas serviceAccounts) List() ([]models.ServiceAccount, error) {
	var saSl []models.ServiceAccount
	if _, err := sas.storage.PG.DB.Query(
		&saSl,
		`SELECT id, name, email, picture FROM service_accounts
		ORDER BY created_at DESC`,
	); err != nil {
		return nil, err
	}
	return saSl, nil
}

func (sas serviceAccounts) Search(
	term string,
) ([]models.ServiceAccount, error) {
	saSl := []models.ServiceAccount{}
	if _, err := sas.storage.PG.DB.Query(
		&saSl,
		`SELECT id, name, email, picture FROM service_accounts
		WHERE name ILIKE '%?0%' OR email ILIKE '%?0%'
		ORDER BY created_at DESC`,
		term,
	); err != nil {
		return nil, err
	}
	return saSl, nil
}

// ForEmail retrieves Service Account corresponding
func (sas serviceAccounts) ForEmail(
	email string,
) (*models.ServiceAccount, error) {
	sa := new(models.ServiceAccount)
	if _, err := sas.storage.PG.DB.Query(
		sa, `SELECT id, name, key_id, key_secret, email, base_role_id
		FROM service_accounts WHERE email = ? LIMIT 1`, email,
	); err != nil {
		return nil, err
	}
	if sa.ID == "" {
		return nil, fmt.Errorf("Service Account not found for email %s", email)
	}
	return sa, nil
}

// ForKeyPair retrieves Service Account corresponding
func (sas serviceAccounts) ForKeyPair(
	keyID, keySecret string,
) (*models.ServiceAccount, error) {
	sa := []*models.ServiceAccount{}
	if _, err := sas.storage.PG.DB.Query(
		&sa, `SELECT id, name, key_id, key_secret, email, base_role_id
		FROM service_accounts WHERE key_id = ? AND key_secret = ?`,
		keyID, keySecret,
	); err != nil {
		return nil, err
	}
	if len(sa) == 0 {
		return nil, fmt.Errorf("service account not found")
	}
	return sa[0], nil
}

func (sas serviceAccounts) Create(sa *models.ServiceAccount) error {
	_, err := sas.storage.PG.DB.Query(
		sa, `INSERT INTO service_accounts (id, name, email, key_id, key_secret,
		base_role_id) VALUES (?id, ?name, ?email, ?key_id, ?key_secret,
		?base_role_id) ON CONFLICT (email) DO UPDATE
		SET picture = ?picture, updated_at = now() RETURNING id`, sa,
	)
	return err
}

// NewServiceAccounts serviceAccounts ctor
func NewServiceAccounts(s *Storage) ServiceAccounts {
	return &serviceAccounts{storage: s}
}
