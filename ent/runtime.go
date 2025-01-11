// Code generated by ent, DO NOT EDIT.

package ent

import (
	"awesomeProject1/ent/schema"
	"awesomeProject1/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescPwd is the schema descriptor for pwd field.
	userDescPwd := userFields[3].Descriptor()
	// user.PwdValidator is a validator for the "pwd" field. It is called by the builders before save.
	user.PwdValidator = func() func(string) error {
		validators := userDescPwd.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
			validators[2].(func(string) error),
		}
		return func(pwd string) error {
			for _, fn := range fns {
				if err := fn(pwd); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() string)
}
