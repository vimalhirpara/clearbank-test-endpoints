package model

type AuthProfile struct {
	Token          string
	PrivateKeyPath string
	PublicKeyPath  string
}

func InitAuthProfile() AuthProfile {
	return AuthProfile{
		Token: "OTMwYzM0ZDBkNmI5NDcxOTg2NDczODEzYTZhY2YxZDk=.eyJpbnN0aXR1dGlvbklkIjoiNmJhZmEwNjItZWI1MC00YmRlLWI2MWYtN2JmOGVkYTZhYzlmIiwibmFtZSI6IlRlc3QiLCJ0b2tlbiI6IkEwQTY4NzZDNEEwMDQ3MTBBNDlDREU0MTY1NTRDQkQ3MDUwREZDNzM5QTU1NDVBODk3MUVBMUE2Mzg1RkExRkU1QTYwMDE0NjU1NjY0NTYxQThDNTk3QUZGMDZFMEU4QiJ9",
		//Token:          "M2FlMjM3MjFlZjJiNDc0ZTlkZmJkM2ZjZmVmYTI2NjU=.eyJpbnN0aXR1dGlvbklkIjoiNmJhZmEwNjItZWI1MC00YmRlLWI2MWYtN2JmOGVkYTZhYzlmIiwibmFtZSI6IkJldGEyMDIxVXB0bzIwMjQiLCJ0b2tlbiI6IjVDMzE5NkJCRjBERDQ0NjdCMkM0NUY0M0ZGQjY4NjczNTU5NzQ1MTI3RkQzNDg5MjkxODNBOUQ2RkFGMDVCQzIyNDcxRkQ0MkVDN0E0NTM3QURCMjBDMEYxQTE5NTlBQyJ9",
		PrivateKeyPath: "../../GoClearBank.prv",
		PublicKeyPath:  "../../GoClearBank.pub"}
}
