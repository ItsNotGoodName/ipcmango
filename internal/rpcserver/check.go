package rpcserver

// import (
// 	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
// 	"github.com/twitchtv/twirp"
// )
//
// func check(err error) twirp.Error {
// 	if core.IsNotFound(err) {
// 		return NewError(err, "Not found.").NotFound()
// 	}
// 	return NewError(err, "Something went wrong.").Internal()
// }
//
// func checkCreateUpdateUser(err error, msg string) error {
// 	if errs, ok := asValidationErrors(err); ok {
// 		return NewError(err, msg).Validation(errs, [][2]string{
// 			{"email", "Email"},
// 			{"username", "Username"},
// 			{"password", "Password"},
// 		})
// 	}
//
// 	if constraintErr, ok := asConstraintError(err); ok {
// 		return NewError(err, "Failed to sign up.").Constraint(constraintErr, [][3]string{
// 			{"username", "users.username", "Name already taken."},
// 			{"email", "users.email", "Email already taken."},
// 		})
// 	}
//
// 	return check(err)
// }
//
// func checkCreateUpdateGroup(err error, msg string) error {
// 	if errs, ok := asValidationErrors(err); ok {
// 		return NewError(err, msg).Validation(errs, [][2]string{
// 			{"name", "Name"},
// 			{"description", "Description"},
// 		})
// 	}
//
// 	if constraintErr, ok := asConstraintError(err); ok {
// 		return NewError(err, msg).Constraint(constraintErr, [][3]string{
// 			{"name", "groups.name", "Name already taken."},
// 		})
// 	}
//
// 	return check(err)
// }
//
// func checkCreateUpdateDevice(err error, msg string) error {
// 	if errs, ok := asValidationErrors(err); ok {
// 		return NewError(err, msg).Validation(errs, [][2]string{
// 			{"name", "Name"},
// 			{"description", "Description"},
// 			{"location", "Location"},
// 		})
// 	}
//
// 	if constraintErr, ok := asConstraintError(err); ok {
// 		return NewError(err, msg).Constraint(constraintErr, [][3]string{
// 			{"name", "dahua_devices.name", "Name already taken."},
// 			{"url", "dahua_devices.ip", "URL already taken."},
// 		})
// 	}
//
// 	return check(err)
// }
