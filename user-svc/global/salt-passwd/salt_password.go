package salt_passwd

import (
	"fmt"
	"strings"

	"crypto/sha512"
	"github.com/anaskhan96/go-password-encoder"
)

/*
	不在User结构体中定义salt,也就是不在数据表中存储salt,而是在密码中同意记录加密算法和salt和密码
*/

// 加密
func SaltPassword(passwd string) string {
	// Using the default options
	//salt, encodedPwd := password.Encode("generic password", nil)
	//check := password.Verify("generic password", salt, encodedPwd, nil)
	//fmt.Println(check) // true

	// Using custom options
	options := &password.Options{16, 100, 30, sha512.New}
	salt, encodedPwd := password.Encode(passwd, options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	return newPassword

	// 注意生成的这个密码长度不要过长,User表中定义的是varchar(100),过长会截断
	//fmt.Println(len(newPassword))

	// 切割后前面有个空的字符串
	//passwordInfo := strings.Split(newPassword, "$")
	// salt password
	//check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
	//fmt.Println(check) // true
}

// 判断密码
func ParsePassword(passwd, encryptedPasswd string) bool {
	// passwd原始的密码 encryptedPasswd加密后的密码

	// 和加密的option一致
	options := &password.Options{16, 100, 30, sha512.New}
	encryptedPasswdInfo := strings.Split(encryptedPasswd, "$")
	check := password.Verify(passwd, encryptedPasswdInfo[2], encryptedPasswdInfo[3], options)
	return check
}
