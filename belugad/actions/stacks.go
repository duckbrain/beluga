package actions

import (
	"strings"

	"github.com/duckbrain/beluga/belugad/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/harlow/authtoken"
)

func jsonError(c buffalo.Context, s int, m string) error {
	return c.Render(s, r.JSON(map[string]string{
		"error": m,
	}))
}

func DeployAuthMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		token, err := authtoken.FromRequest(c.Request())
		if err != nil {
			return jsonError(c, 400, "Key not provided")
		}
		domain := c.Param("domain")
		if len(domain) == 0 {
			return jsonError(c, 404, "Stack name not provided")
		}
		if strings.HasPrefix(domain, "beluga") {
			return jsonError(c, 404, "\"beluga\" prefix in stack name is reserved for internal use")
		}
		c.Logger().Debug("stack: ", domain)
		c.Logger().Debug("token: ", token)

		tx := c.Value("tx").(*pop.Connection)
		stack := models.Stack{}
		err = tx.Where("domain = ?", domain).First(&stack)
		if err != nil {
			return jsonError(c, 404, "Stack not found")
		}

		c.Set("stack", stack)

		return next(c)
	}
}

func StackDeploy(c buffalo.Context) error {
	stack := c.Value("stack").(models.Stack)
	deployment := models.Deployment{}
	if err := c.Bind(&deployment); err != nil {
		return err
	}
	return deployer.Deploy(stack.Name, deployment)
}
func StackDestroy(c buffalo.Context) error {
	stack := c.Value("stack").(models.Stack)
	deployment := models.Deployment{}
	if err := c.Bind(deployment); err != nil {
		return err
	}
	return deployer.Deploy(stack.Name, deployment)
}
