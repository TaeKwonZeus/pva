package data

import (
	"github.com/TaeKwonZeus/pva/crypt"
	"github.com/TaeKwonZeus/pva/network"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
)

// Store abstracts away cryptographic operations on data from db.
type Store struct {
	db *db
}

func NewStore(path string) (*Store, error) {
	pool, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = pool.Exec(startupQuery)
	if err != nil {
		return nil, err
	}

	return &Store{db: &db{pool}}, nil
}

func (s *Store) Close() error {
	return s.db.pool.Close()
}

func (s *Store) VerifyPassword(username string, password string) (verified bool, user User) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return
	}

	key := crypt.DeriveKey(password, user.Salt)
	if _, err = crypt.AesDecrypt(user.PrivateKeyEncrypted, key); err != nil {
		return
	}

	return true, user
}

func (s *Store) GetIndex(id int) (index Index, err error) {
	return s.db.getIndex(id)
}

func (s *Store) GetUserCount() (n int, err error) {
	return s.db.getUserCount()
}

func (s *Store) GetUser(id int) (user User, err error) {
	return s.db.getUser(id)
}

func (s *Store) GetUserByUsername(username string) (user User, err error) {
	return s.db.getUserByUsername(username)
}

func (s *Store) CreateUser(user User, password string) error {
	privateKey, publicKey, err := crypt.NewKeypair()
	if err != nil {
		return err
	}

	user.Salt, err = crypt.GenerateSalt()
	if err != nil {
		return err
	}

	key := crypt.DeriveKey(password, user.Salt)
	user.PrivateKeyEncrypted, err = crypt.AesEncrypt(privateKey, key)
	user.PublicKey = publicKey
	if err != nil {
		return err
	}

	return s.db.createUser(user)
}

func (s *Store) CreateVault(vault Vault, user User) error {
	key, err := crypt.NewAesKey()
	if err != nil {
		return err
	}

	vaultKeyEncrypted, err := crypt.RsaEncrypt(key, user.PublicKey)
	if err != nil {
		return err
	}

	if err = s.db.createVault(vault, vaultKeyEncrypted, user.ID); err != nil {
		return err
	}

	// Add a keyEncrypted for all admins
	admins, err := s.db.getAdmins()
	if err != nil {
		return err
	}

	var vaultKeys []vaultKey
	for _, admin := range admins {
		vaultKeyEncrypted, err = crypt.RsaEncrypt(key, admin.PublicKey)
		if err != nil {
			return err
		}
		vaultKeys = append(vaultKeys, vaultKey{UserId: admin.ID, VaultId: vault.ID, KeyEncrypted: vaultKeyEncrypted})
	}
	return s.db.createVaultKeys(vaultKeys...)
}

func decryptVault(vault *Vault, user User) error {
	vaultKey, err := crypt.RsaDecrypt(vault.KeyEncrypted, user.PrivateKey)
	if err != nil {
		return err
	}

	for i := range vault.Passwords {
		passwordDecrypted, err := crypt.AesDecrypt(vault.Passwords[i].PasswordEncrypted, vaultKey)
		if err != nil {
			return err
		}

		vault.Passwords[i].Password = string(passwordDecrypted)
	}

	return nil
}

func (s *Store) GetVault(id int, user User) (vault Vault, err error) {
	vault, err = s.db.getVault(id, user.ID)
	if err != nil {
		return
	}
	err = decryptVault(&vault, user)
	return
}

func (s *Store) GetVaults(user User) (vaults []Vault, err error) {
	vaults, err = s.db.getVaults(user.ID)
	if err != nil {
		return
	}
	log.Infof("Before: %v", vaults)

	for i := range vaults {
		err = decryptVault(&vaults[i], user)
		if err != nil {
			return nil, err
		}
	}
	log.Infof("After: %v", vaults)
	return
}

func (s *Store) CheckVaultOwnership(vaultId int, user User) bool {
	keyEncrypted, err := s.db.getVaultKey(vaultId, user.ID)
	if err != nil {
		return false
	}
	_, err = crypt.RsaDecrypt(keyEncrypted, user.PrivateKey)
	return err == nil
}

func (s *Store) UpdateVault(vault Vault) error {
	return s.db.updateVault(vault)
}

func (s *Store) DeleteVault(id int) error {
	return s.db.deleteVault(id)
}

func (s *Store) getDecryptedVaultKey(vaultId int, user User) ([]byte, error) {
	keyEncrypted, err := s.db.getVaultKey(vaultId, user.ID)
	if err != nil {
		return nil, err
	}
	return crypt.RsaDecrypt(keyEncrypted, user.PrivateKey)
}

func (s *Store) ShareVault(vaultId int, target User, user User) error {
	key, err := s.getDecryptedVaultKey(vaultId, user)
	if err != nil {
		return err
	}

	keyEncrypted, err := crypt.RsaEncrypt(key, target.PublicKey)
	if err != nil {
		return err
	}

	return s.db.createVaultKeys(vaultKey{
		UserId:       target.ID,
		VaultId:      vaultId,
		KeyEncrypted: keyEncrypted,
	})
}

func (s *Store) CreatePassword(password Password, vaultId int, user User) error {
	vault, err := s.db.getVault(vaultId, user.ID)
	if err != nil {
		return err
	}

	vaultKey, err := crypt.RsaDecrypt(vault.KeyEncrypted, user.PrivateKey)
	if err != nil {
		return err
	}
	password.PasswordEncrypted, err = crypt.AesEncrypt([]byte(password.Password), vaultKey)
	if err != nil {
		return err
	}

	return s.db.createPassword(password, vaultId)
}

func (s *Store) UpdatePassword(password Password, vaultId int, user User) error {
	// If password isn't being updated we can skip any cryptographic operations altogether
	if password.Password == "" {
		return s.db.updatePassword(password)
	}

	key, err := s.getDecryptedVaultKey(vaultId, user)
	if err != nil {
		return err
	}

	password.PasswordEncrypted, err = crypt.AesEncrypt([]byte(password.Password), key)
	if err != nil {
		return err
	}

	return s.db.updatePassword(password)
}

func (s *Store) DeletePassword(id int) error {
	return s.db.deletePassword(id)
}

func (s *Store) CreateDevice(device Device) error {
	return s.db.createDevice(device)
}

func (s *Store) GetDevices() (devices []Device, err error) {
	devices, err = s.db.getDevices()
	if err != nil {
		return
	}
	deviceMap := make(map[string]Device)
	for _, device := range devices {
		deviceMap[device.IP] = device
	}

	scan := network.Devices()

	// Add connected devices to the response, whether they're saved or not.
	// If they are saved but not connected, Connected will equal false.
	// If they are connected but not saved, ID will equal 0.
	for _, device := range scan {
		// Entry is a pointer to the entry in the slice so we can just edit it
		entry, ok := deviceMap[device.IP.String()]
		if ok {
			//entry.NetworkName = device.Name
			//entry.MAC = device.MAC.String()
			entry.Connected = true
		} else {
			devices = append(devices, Device{
				IP: device.IP.String(),
				//NetworkName: device.Name,
				//MAC:         device.MAC.String(),
				Connected: true,
			})
		}
	}

	return
}

func (s *Store) UpdateDevice(device Device) error {
	return s.db.updateDevice(device)
}

func (s *Store) DeleteDevice(id int) error {
	return s.db.deleteDevice(id)
}
