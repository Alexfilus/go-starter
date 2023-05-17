package entity

// GooKit flags
type ProfileRequest struct {
	Name  string `validate:"required|min_len:7" conform:"trim"`
	Email string `validate:"email"  message:"email is invalid." label:"User Email"`
	Phone string `validate:"required|isE164PhoneNumber"  label:"User Phone"`
	Age   int    `validate:"required|int|min:1|max:99"`
}
