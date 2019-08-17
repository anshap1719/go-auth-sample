package auth

import (
	"gigglesearch.org/giggle-auth/auth/app"
	"gigglesearch.org/giggle-auth/auth/database"
	"gigglesearch.org/giggle-auth/utils/auth"
	"github.com/globalsign/mgo"
	"github.com/goadesign/goa"
)

// NewsletterController implements the newsletter resource.
type NewsletterController struct {
	*goa.Controller
	*auth.JWTSecurity
}

// NewNewsletterController creates a newsletter controller.
func NewNewsletterController(service *goa.Service, jwtSec auth.JWTSecurity) *NewsletterController {
	return &NewsletterController{
		Controller:  service.NewController("NewsletterController"),
		JWTSecurity: &jwtSec,
	}
}

func (c *NewsletterController) GetSubscribers(ctx *app.GetSubscribersNewsletterContext) error {
	if c.IsAdmin(ctx.Request) {
		subs, err := database.GetNewsletterSubscribers()
		if err != nil {
			return ctx.InternalServerError(goa.ErrInternal(err))
		}

		return ctx.OK(database.SubscribersToAppSubscribers(subs))
	} else {
		return ctx.NotFound(goa.ErrNotFound(ctx.RequestURI))
	}
}

func (c *NewsletterController) GetSubscriberByEmail(ctx *app.GetSubscriberByEmailNewsletterContext) error {
	sub, err := database.GetNewsletterSubscriber(*ctx.Email)
	if err == mgo.ErrNotFound {
		return ctx.NotFound(goa.ErrNotFound(err))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK(database.SubscriberToAppSubscriber(sub))
}

// AddSubscriber runs the add-subscriber action.
func (c *NewsletterController) AddSubscriber(ctx *app.AddSubscriberNewsletterContext) error {
	// NewsletterController_AddSubscriber: start_implement

	if err := database.AddSubscriber(database.AppSubscriberToSubscriber(*ctx.Payload)); err == database.ErrAlreadySubscribed {
		return ctx.BadRequest(goa.ErrBadRequest("user is already subscribed"))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(""))

	// NewsletterController_AddSubscriber: end_implement
}

// RemoveSubscriber runs the remove-subscriber action.
func (c *NewsletterController) RemoveSubscriber(ctx *app.RemoveSubscriberNewsletterContext) error {
	// NewsletterController_RemoveSubscriber: start_implement

	if err := database.RemoveSubscriber(*ctx.Email); err != nil {
		return ctx.BadRequest(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(""))
	// NewsletterController_RemoveSubscriber: end_implement
}

func (c *NewsletterController) UpdateSubscriber(ctx *app.UpdateSubscriberNewsletterContext) error {
	var updatedSub = database.AppSubscriberToSubscriber(*ctx.Payload)

	if err := database.UpdateSubscriber(updatedSub); err == mgo.ErrNotFound {
		return ctx.NotFound(goa.ErrNotFound(err))
	} else if err != nil {
		return ctx.InternalServerError(goa.ErrInternal(err))
	}

	return ctx.OK([]byte(""))
}
