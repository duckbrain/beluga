package grifts

import (
	"fmt"

	"github.com/duckbrain/beluga/belugad/models"
	"github.com/gobuffalo/nulls"
	"github.com/manifoldco/promptui"
	. "github.com/markbates/grift/grift"
)

var _ = Namespace("users", func() {

	Desc("new", "Create a new user")
	Add("new", func(c *Context) error {
		user := &models.User{}

		key, err := user.GenerateKey()
		if err != nil {
			return err
		}

		username, err := (&promptui.Prompt{
			Label: "Username (leave blank for deploy-only account)",
		}).Run()
		if err != nil {
			return err
		}

		if len(username) > 0 {
			password, err := (&promptui.Prompt{
				Label: "Password",
				Mask:  '*',
			}).Run()
			if err != nil {
				return err
			}

			user.Username = nulls.NewString(username)
			err = user.SetPassword(password)
			if err != nil {
				return err
			}

			i, _, err := (&promptui.Select{
				Label: "Is Admin",
				Items: []string{"no", "yes"},
				Size:  2,
			}).Run()

			user.IsAdmin = i != 0
		}

		err = models.DB.Create(user)
		if err != nil {
			return err
		}

		fmt.Printf("Deploy Key: %v\n", key)

		return nil
	})

	Desc("clear", "Delete all users")
	Add("clear", func(c *Context) error {
		return models.DB.RawQuery("DELETE FROM users;").Exec()
	})

})
