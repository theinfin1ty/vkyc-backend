package seeds

import (
	"fmt"
	"net/mail"
	"vkyc-backend/configs"
	"vkyc-backend/models"
	"vkyc-backend/utils"
)

func CreateDefaultAdmin() {
	email := configs.GetEnvVariable("ADMIN_EMAIL")
	var user models.User

	if email == "" {
		fmt.Println("Admin email not set in env vars")
		return
	}

	_, err := mail.ParseAddress(email)

	if err != nil {
		fmt.Println(err)
		return
	}

	result := configs.DB.Find(&user, models.User{
		Email: email,
	}).Limit(1)

	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	if result.RowsAffected > 0 {
		fmt.Println("Default admin already exists")
		return
	}

	password := utils.GeneratePassword(8, true)

	passwordHash, err := utils.GeneratePasswordHash(password)

	if err != nil {
		fmt.Println(err)
		return
	}

	user = models.User{
		FirstName: "Admin",
		LastName:  "Admin",
		Email:     email,
		Password:  passwordHash,
		Role:      utils.AdminRole,
		// DOB:       time.Now(),
	}

	result = configs.DB.Create(&user)

	if result.Error != nil {
		fmt.Println(result.Error)
		return
	}

	err = utils.SendInviteEmail(&user, password)

	if err != nil {
		fmt.Println(err)
		return
	}

	// err = utils.SendInviteSMS(&user, password)

	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	fmt.Println("Default admin created")
}
