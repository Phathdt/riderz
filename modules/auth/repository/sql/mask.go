package authRepo

func (u *User) Mask() {
	u.Password = ""
}
