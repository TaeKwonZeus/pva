package data

import (
	"database/sql"
	"github.com/TaeKwonZeus/pva/crypt"
	"github.com/TaeKwonZeus/pva/network"
	"time"
)

// Store abstracts away cryptographic operations on data from db.
type Store struct {
	db *db
}

func NewStore(path string) (*Store, error) {
	pool, err := sql.Open("sqlite3", path)
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

func (s *Store) VerifyPassword(username string, password string) (verified bool, user *User) {
	user, err := s.GetUserByUsername(username)
	if err != nil {
		return false, nil
	}

	key := crypt.DeriveKey(password, user.salt)
	if _, err := crypt.AesDecrypt(user.privateKeyEncrypted, key); err != nil {
		return false, nil
	}

	return true, user
}

func (s *Store) GetIndex(id int) (index *Index, err error) {
	return s.db.getIndex(id)
}

func (s *Store) GetUserCount() (n int, err error) {
	return s.db.getUserCount()
}

func (s *Store) GetUser(id int) (*User, error) {
	return s.db.getUser(id)
}

func (s *Store) GetUserByUsername(username string) (*User, error) {
	return s.db.getUserByUsername(username)
}

func (s *Store) CreateUser(user *User, password string) error {
	privateKey, publicKey, err := crypt.NewKeypair()
	if err != nil {
		return err
	}

	user.salt, err = crypt.GenerateSalt()
	if err != nil {
		return err
	}

	key := crypt.DeriveKey(password, user.salt)
	user.privateKeyEncrypted, err = crypt.AesEncrypt(privateKey, key)
	user.publicKey = publicKey
	if err != nil {
		return err
	}

	return s.db.createUser(user)
}

func (s *Store) CreateVault(vault *Vault, owner *User) error {
	key, err := crypt.NewAesKey()
	if err != nil {
		return err
	}

	vaultKeyEncrypted, err := crypt.RsaEncrypt(key, owner.publicKey)
	if err != nil {
		return err
	}

	vault.OwnerId = owner.ID
	if err = s.db.createVault(vault, vaultKeyEncrypted); err != nil {
		return err
	}

	// Add a key for all admins
	admins, err := s.db.getAdmins()
	if err != nil {
		return err
	}

	var vaultKeys []*vaultKey
	for _, admin := range admins {
		vaultKeyEncrypted, err = crypt.RsaEncrypt(key, admin.publicKey)
		if err != nil {
			return err
		}
		vaultKeys = append(vaultKeys, &vaultKey{userId: admin.ID, vaultId: vault.ID, keyEncrypted: vaultKeyEncrypted})
	}
	return s.db.createVaultKeys(vaultKeys...)
}

func decryptVault(vnk *vaultAndKey, user *User) (*Vault, error) {
	vaultKey, err := crypt.RsaDecrypt(vnk.keyEncrypted, user.privateKey)
	if err != nil {
		return nil, err
	}

	for _, p := range vnk.vault.Passwords {
		passwordDecrypted, err := crypt.AesDecrypt(p.passwordEncrypted, vaultKey)
		if err != nil {
			return nil, err
		}

		p.Password = string(passwordDecrypted)
	}

	return vnk.vault, nil
}

func (s *Store) GetVault(id int, user *User) (*Vault, error) {
	vnk, err := s.db.getVault(id, user.ID)
	if err != nil {
		return nil, err
	}

	return decryptVault(vnk, user)
}

func (s *Store) GetVaults(user *User) ([]*Vault, error) {
	vnks, err := s.db.getVaults(user.ID)
	if err != nil {
		return nil, err
	}

	out := make([]*Vault, len(vnks))

	for i, vnk := range vnks {
		out[i], err = decryptVault(vnk, user)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func (s *Store) CheckVaultOwnership(vaultId int, user *User) bool {
	keyEncrypted, err := s.db.getVaultKey(vaultId, user.ID)
	if err != nil {
		return false
	}
	_, err = crypt.RsaDecrypt(keyEncrypted, user.privateKey)
	return err == nil
}

func (s *Store) UpdateVault(vault *Vault) error {
	return s.db.updateVault(vault)
}

func (s *Store) DeleteVault(id int) error {
	return s.db.deleteVault(id)
}

func (s *Store) getDecryptedVaultKey(vaultId int, user *User) ([]byte, error) {
	keyEncrypted, err := s.db.getVaultKey(vaultId, user.ID)
	if err != nil {
		return nil, err
	}
	return crypt.RsaDecrypt(keyEncrypted, user.privateKey)
}

func (s *Store) ShareVault(vaultId int, target *User, user *User) error {
	key, err := s.getDecryptedVaultKey(vaultId, user)
	if err != nil {
		return err
	}

	keyEncrypted, err := crypt.RsaEncrypt(key, target.publicKey)
	if err != nil {
		return err
	}

	return s.db.createVaultKeys(&vaultKey{
		userId:       target.ID,
		vaultId:      vaultId,
		keyEncrypted: keyEncrypted,
	})
}

func (s *Store) CreatePassword(password *Password, vaultId int, user *User) error {
	password.CreatedAt = time.Now()
	password.UpdatedAt = time.Now()

	vnk, err := s.db.getVault(vaultId, user.ID)
	if err != nil {
		return err
	}

	vaultKey, err := crypt.RsaDecrypt(vnk.keyEncrypted, user.privateKey)
	if err != nil {
		return err
	}
	password.passwordEncrypted, err = crypt.AesEncrypt([]byte(password.Password), vaultKey)
	if err != nil {
		return err
	}

	return s.db.createPassword(password, vaultId)
}

func (s *Store) UpdatePassword(password *Password, vaultId int, user *User) error {
	password.UpdatedAt = time.Now()

	// If password isn't being updated we can skip any cryptographic operations altogether
	if password.Password == "" {
		return s.db.updatePassword(password)
	}

	key, err := s.getDecryptedVaultKey(vaultId, user)
	if err != nil {
		return err
	}

	password.passwordEncrypted, err = crypt.AesEncrypt([]byte(password.Password), key)
	if err != nil {
		return err
	}

	return s.db.updatePassword(password)
}

func (s *Store) DeletePassword(id int) error {
	return s.db.deletePassword(id)
}

func (s *Store) CreateDevice(device *Device) error {
	return s.db.createDevice(device)
}

func (s *Store) GetDevices() (devices []*Device, err error) {
	devices, err = s.db.getDevices()
	if err != nil {
		return
	}
	deviceMap := make(map[string]*Device)
	for _, device := range devices {
		deviceMap[device.IP] = device
	}

	scan, err := network.Devices()
	if err != nil {
		return nil, err
	}

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
			devices = append(devices, &Device{
				IP: device.IP.String(),
				//NetworkName: device.Name,
				//MAC:         device.MAC.String(),
				Connected: true,
			})
		}
	}

	return
}

func (s *Store) UpdateDevice(device *Device) error {
	return s.db.updateDevice(device)
}

func (s *Store) DeleteDevice(id int) error {
	return s.db.deleteDevice(id)
}
