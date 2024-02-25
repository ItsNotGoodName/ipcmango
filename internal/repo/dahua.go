package repo

// func DahuaSelectFilter(ctx context.Context, sb sq.SelectBuilder, deviceIDField string, levels ...models.DahuaPermissionLevel) sq.SelectBuilder {
// 	actor := core.UseActor(ctx)
//
// 	if actor.Admin {
// 		return sb
// 	}
//
// 	var level models.DahuaPermissionLevel
// 	if len(levels) != 0 {
// 		level = levels[0]
// 	}
//
// 	return sb.
// 		Where(sq.Expr(deviceIDField+` IN (
// 			SELECT
// 				device_id
// 			FROM
// 				dahua_permissions
// 			WHERE
// 				dahua_permissions.level > ?
// 				AND (
// 					dahua_permissions.user_id = ?
// 					OR dahua_permissions.group_id IN (
// 						SELECT
// 							group_id
// 						FROM
// 							group_users
// 						WHERE
// 							group_users.user_id = ?
// 					)
// 				)
// 			)
// 		`, level, actor.UserID, actor.UserID))
// }
//
// // DahuaFatDevice
//
// type DahuaFatDevice struct {
// 	DahuaDevice
// 	Seed int64
// }
//
// type DahuaFatDeviceParams struct {
// 	IPs      []string
// 	IDs      []int64
// 	Features []models.DahuaFeature
// 	Limit    int
// }
//
// func (db DB) DahuaListFatDevices(ctx context.Context, args ...DahuaFatDeviceParams) ([]DahuaFatDevice, error) {
// 	var arg DahuaFatDeviceParams
// 	if len(args) != 0 {
// 		arg = args[0]
// 	}
//
// 	// SELECT ...
// 	sb := sq.
// 		Select(
// 			"dahua_devices.*",
// 			"coalesce(seed, id) AS seed",
// 		).
// 		From("dahua_devices").
// 		LeftJoin("dahua_seeds ON dahua_seeds.device_id = dahua_devices.id")
// 	// WHERE
// 	and := sq.And{}
//
// 	eq := sq.Eq{}
// 	if arg.IPs != nil {
// 		eq["ip"] = arg.IPs
// 	}
// 	if arg.IDs != nil {
// 		eq["id"] = arg.IDs
// 	}
// 	and = append(and, eq)
//
// 	if len(arg.Features) != 0 {
// 		var feature models.DahuaFeature
// 		for _, f := range arg.Features {
// 			feature = feature | f
// 		}
// 		and = append(and, sq.Expr("feature & ? = ?", feature, feature))
// 	}
//
// 	sb = sb.Where(and)
// 	// LIMIT
// 	if arg.Limit != 0 {
// 		sb = sb.Limit(uint64(arg.Limit))
// 	}
//
// 	var res []DahuaFatDevice
// 	return res, ssq.Query(ctx, db, &res, DahuaSelectFilter(ctx, sb, "dahua_devices.id"))
// }
//
// func (db DB) DahuaGetFatDevice(ctx context.Context, arg DahuaFatDeviceParams) (DahuaFatDevice, error) {
// 	arg.Limit = 1
// 	devices, err := db.DahuaListFatDevices(ctx, arg)
// 	if err != nil {
// 		return DahuaFatDevice{}, err
// 	}
// 	if len(devices) == 0 {
// 		return DahuaFatDevice{}, ErrNotFound
// 	}
// 	return devices[0], nil
// }
