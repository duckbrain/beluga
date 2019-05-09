package grifts

import (
	"errors"
	"fmt"

	"github.com/duckbrain/beluga/belugad/models"
	. "github.com/markbates/grift/grift"
)

var _ = Namespace("stack", func() {

	Desc("new", "Create a new stack")
	Add("new", func(c *Context) error {
		if len(c.Args) == 0 {
			return errors.New("Usage: stack:new STACK_NAME...")
		}

		fmt.Printf("name key\n")
		for _, name := range c.Args {
			key, err := models.GenerateKey(models.KeyLength)
			if err != nil {
				return err
			}
			stack := models.Stack{
				Name: name,
				Key:  key,
			}

			err = models.DB.Create(&stack)
			if err != nil {
				return err
			}
			fmt.Printf("%v %v\n", name, key)
		}
		return nil
	})

})
