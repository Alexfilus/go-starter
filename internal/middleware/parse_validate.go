package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/otyang/go-pkg/response"
	"github.com/otyang/go-pkg/validators"
)

// ParseBody reads either a json or xml or html forms or urls body
// and then parse to a struct. After which if there is an error it detects error type
func ParseBody(c *fiber.Ctx, payloadPtr any) error {
	if err := c.BodyParser(payloadPtr); err != nil {
		_, err = response.IsJsonErrorGetDetails(err)
		return err
	}
	return nil
}

// ValidateBody validates a struct and returns error if any found
func ValidateBody(payloadPtr any) error {
	err := validators.HelperValidateStructGoKit(payloadPtr)
	if err != nil {
		return err
	}
	return nil
}

// ParseAndValidate performs the work of both PARSEBODY and VALIDATEBODY
func ParseAndValidate(ctx *fiber.Ctx, payloadPtr any) error {
	if err := ParseBody(ctx, payloadPtr); err != nil {
		return err
	}

	if err := ValidateBody(payloadPtr); err != nil {
		return err
	}
	return nil
}

// MwValidateBody performs same function as ParseAndValidate but this is a middleware
func MwValidateBody[bodyType any](c *fiber.Ctx) error {
	body := new(bodyType)

	if err := ParseBody(c, body); err != nil {
		rsp := response.BadRequest(err.Error(), "")
		return c.Status(rsp.StatusCode).JSON(rsp)
	}

	if err := ValidateBody(body); err != nil {
		rsp := response.BadRequest(err.Error(), "")
		return c.Status(rsp.StatusCode).JSON(rsp)
	}

	c.Locals("jsonBody", body)
	return c.Next() // call next
}

// ValidatedDataFromContext(c).(*type)
func ValidatedDataFromContext(ctx *fiber.Ctx) (val any) {
	return ctx.Locals("jsonBody")
}
