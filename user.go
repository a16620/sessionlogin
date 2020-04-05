package main

import (
	"container/list"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"
)

//User : ...
type User struct {
	name       string
	password   [32]byte
	permission uint8
}

//MarshalJSON as Implement
func (u *User) MarshalJSON() ([]byte, error) {
	var tmp struct {
		Uname string `json:"name"`
	}
	tmp.Uname = u.name
	return json.Marshal(tmp)
}

//UserPtr : ...
type UserPtr struct {
	Expired    time.Time
	user       *User
	sessionKey string
}

//Refresh : Refresh session lifetime
func (uptr *UserPtr) Refresh() {
	uptr.Expired = time.Now().Add(time.Minute * 5)
}

//IsValid : Check the session is not expired
func (uptr *UserPtr) IsValid() bool {
	if uptr.Expired.After(time.Now()) {
		return true
	}
	return false
}

//MakeSessionKey is make hash from time and salt
func MakeSessionKey(salt string) (string, error) {
	biTime, e := time.Now().MarshalBinary()
	if e != nil {
		return "", e
	}
	bf := append(biTime, []byte(salt)...)
	hash := sha256.Sum256(bf)
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

//AccountManager : Managing accounts
type AccountManager struct {
	accounts  map[string]*User
	accLock   sync.RWMutex
	verifieds map[string]*UserPtr
	Lock      sync.RWMutex
}

//CollectExpired : delete expireds
func (mgr *AccountManager) CollectExpired() {
	eList := list.New()
	mgr.Lock.RLock()
	for k, v := range mgr.verifieds {
		if !v.IsValid() {
			eList.PushBack(k)
		}
	}
	mgr.Lock.RUnlock()

	if eList.Len() == 0 {
		return
	}

	mgr.Lock.Lock()
	for k := eList.Front(); k != nil; k = k.Next() {
		kc, _ := (k.Value).(string)
		delete(mgr.verifieds, kc)
	}
	mgr.Lock.Unlock()
}

//GetAccounts ...
func (mgr *AccountManager) GetAccounts() ([]byte, error) {
	mgr.Lock.RLock()
	defer mgr.Lock.RUnlock()
	return json.Marshal(mgr.accounts)
}

//LoadAccount : load information from file
func (mgr *AccountManager) LoadAccount(path string) {
	user := new(User)
	user.name = "a16620"
	user.password = sha256.Sum256([]byte("passw"))
	mgr.accounts["a16620"] = user
}

//Logout : delete session key on the table
func (mgr *AccountManager) Logout(key string) {
	mgr.Lock.Lock()
	delete(mgr.verifieds, key)
	mgr.Lock.Unlock()
}

//Verify : register on table
func (mgr *AccountManager) Verify(id, pw string) (*UserPtr, bool) {
	mgr.accLock.RLock()
	user, ext := mgr.accounts[id]
	mgr.accLock.RUnlock()
	if !ext {
		return nil, false
	}

	if user.password == sha256.Sum256([]byte(pw)) {
		uptr := new(UserPtr)
		uptr.user = user
		uptr.Refresh()

		sk, err := MakeSessionKey(id)

		if err != nil {
			//오류처리
		}
		sessionKey := string(sk)
		uptr.sessionKey = sessionKey
		mgr.Lock.Lock()
		mgr.verifieds[sessionKey] = uptr
		mgr.Lock.Unlock()

		return uptr, true
	}

	return nil, false
}

//CheckSession : Check and refresh session
func (mgr *AccountManager) CheckSession(key string) (*UserPtr, bool) {
	mgr.Lock.RLock()
	user, ext := mgr.verifieds[key]
	mgr.Lock.RUnlock()
	if !ext {
		return nil, false
	}

	if !user.IsValid() {
		mgr.Lock.Lock()
		delete(mgr.verifieds, key)
		mgr.Lock.Unlock()
		return nil, false
	}

	user.Refresh()

	return user, true
}

//SessionKeyUser : Identifier
const SessionKeyUser string = "session_user_key"
