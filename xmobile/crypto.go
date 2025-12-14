package xmobile

import "go.olapie.com/x/xsecurity"

func Encrypt(data []byte, passphrase string) []byte {
	content, _ := xsecurity.Encrypt(data, passphrase)
	return content
}

func Decrypt(data []byte, passphrase string) []byte {
	content, _ := xsecurity.Encrypt(data, passphrase)
	return content
}

func EncryptFile(src, dst, passphrase string) bool {
	err := xsecurity.EncryptFile(xsecurity.SF(src), xsecurity.DF(dst), passphrase)
	return err == nil
}

func DecryptFile(src, dst, passphrase string) bool {
	err := xsecurity.DecryptFile(xsecurity.SF(src), xsecurity.DF(dst), passphrase)
	return err == nil
}
