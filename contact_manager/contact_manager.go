package contact_manager

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/g45t345rt/g45w/settings"
)

type Contact struct {
	Name      string `json:"name"`
	Addr      string `json:"addr"`
	Note      string `json:"note"`
	Timestamp int64  `json:"timestamp"`
	ListOrder int    `json:"order"`
}

type ContactManager struct {
	Contacts   map[string]Contact
	WalletAddr string
}

func NewContactManager(walletAddr string) *ContactManager {
	return &ContactManager{
		WalletAddr: walletAddr,
	}
}

func (c *ContactManager) contactsPath() string {
	walletDir := settings.Instance.WalletsDir
	return filepath.Join(walletDir, c.WalletAddr, "contacts.json")
}

func (c *ContactManager) Load() error {
	data, err := os.ReadFile(c.contactsPath())
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &c.Contacts)
	if err != nil {
		return err
	}

	return nil
}

func (c *ContactManager) AddContact(newContact Contact) error {
	c.Contacts[newContact.Addr] = newContact
	return c.saveContacts()
}

func (c *ContactManager) DelContact(addr string) error {
	delete(c.Contacts, addr)
	return c.saveContacts()
}

func (c *ContactManager) ClearContacts(newContact Contact) error {
	c.Contacts = make(map[string]Contact)
	return c.saveContacts()
}

func (c *ContactManager) saveContacts() error {
	data, err := json.Marshal(c.Contacts)
	if err != nil {
		return err
	}

	return os.WriteFile(c.contactsPath(), data, os.ModePerm)
}
