package apiv1

import (
	"fmt"
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
)

// ListUsers returns all registered users
func ListUsers() web.HandlerFunc {
	return func(c *web.Context) error {
		message := "ListUsers...."
		log.Debug(c, message)
		allUsers := &query.GetAllUsers{}
		if err := bus.Dispatch(c, allUsers); err != nil {
			return c.Failure(err)
		}
		return c.Ok(allUsers.Result)
	}
}

// FindUser returns based on uid, or reference
func SearchUsers() web.HandlerFunc {
	return func(c *web.Context) error {
		ref := fmt.Sprintf("SearchUsers ref: %s", c.QueryParam("ref"))
		log.Debug(c, ref)
		getByReference := &query.GetUserByProvider{
			Provider: c.QueryParam("ref"),
			UID:      c.QueryParam("uid"),
		}
		if err := bus.Dispatch(c, getByReference); err != nil {
			return c.Failure(err)
		}
		return c.Ok(getByReference.Result)
	}
}

// CreateUser is used to create new users
func CreateUser() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.CreateUser)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		var user *models.User

		getByReference := &query.GetUserByProvider{Provider: "reference", UID: input.Model.Reference}
		err := bus.Dispatch(c, getByReference)
		user = getByReference.Result

		if err != nil && errors.Cause(err) == app.ErrNotFound {
			if input.Model.Email != "" {
				getByEmail := &query.GetUserByEmail{Email: input.Model.Email}
				err = bus.Dispatch(c, getByEmail)
				user = getByEmail.Result
			}
			if err != nil && errors.Cause(err) == app.ErrNotFound {
				user = &models.User{
					Tenant: c.Tenant(),
					Name:   input.Model.Name,
					Email:  input.Model.Email,
					Role:   enum.RoleVisitor,
				}
				err = bus.Dispatch(c, &cmd.RegisterUser{User: user})
			}
		}

		if err != nil {
			return c.Failure(err)
		}

		if input.Model.Reference != "" && !user.HasProvider("reference") {
			if err := bus.Dispatch(c, &cmd.RegisterUserProvider{
				UserID:       user.ID,
				ProviderName: "reference",
				ProviderUID:  input.Model.Reference,
			}); err != nil {
				return c.Failure(err)
			}
		}

		return c.Ok(web.Map{
			"id": user.ID,
		})
	}
}

func ChangeUserRole() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.CreateUser)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		return c.Ok(web.Map{
			"id": 2,
		})
	}
}

// DeleteUser erases current user personal data and sign them out
func DeleteUser() web.HandlerFunc {
	return func(c *web.Context) error {
		if err := bus.Dispatch(c, &cmd.DeleteCurrentUser{}); err != nil {
			return c.Failure(err)
		}

		c.RemoveCookie(web.CookieAuthName)
		return c.Ok(web.Map{})
	}
}
