package passcrypto

type WithPasswordHasherable interface {
	GetHasher() PasswordHashable
	SetHasher(hashFunc PasswordHashable)
}

type WithPasswordHasher struct {
	hasher PasswordHashable
}

var _ WithPasswordHasherable = (*WithPasswordHasher)(nil)

func (w *WithPasswordHasher) GetHasher() PasswordHashable {
	if w.hasher == nil {
		w.hasher = defaultPasswordHasher{}
	}

	return w.hasher
}

func (w *WithPasswordHasher) SetHasher(hashFunc PasswordHashable) {
	if hashFunc == nil {
		hashFunc = defaultPasswordHasher{}
	}

	w.hasher = hashFunc
}
