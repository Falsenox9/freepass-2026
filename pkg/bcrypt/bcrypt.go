package bcrypt

import (
    "golang.org/x/crypto/bcrypt"
)

type Interface interface {
    GenerateFromPassword(password string) (string, error)
    CompareAndHashPassword(hashedPassword, password string) error
}

type bcryptHash struct {
    cost int
}

func Init() Interface {
    return &bcryptHash{
        cost: 12,
    }
}

func (b *bcryptHash) GenerateFromPassword(password string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
    if err != nil {
        return "", err
    }
    return string(hashedBytes), nil
}

func (b *bcryptHash) CompareAndHashPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}