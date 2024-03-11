package rpcserver

import (
	"context"
	"net/url"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewAdmin(db sqlite.DB) *Admin {
	return &Admin{
		db: db,
	}
}

type Admin struct {
	db sqlite.DB
}

// ---------- Device

func (a *Admin) GetAdminDevicesPage(ctx context.Context, req *rpc.GetAdminDevicesPageReq) (*rpc.GetAdminDevicesPageResp, error) {
	page := decodePagePagination(req.Page)

	items, err := func() ([]*rpc.GetAdminDevicesPageResp_Device, error) {
		var row struct {
			repo.DahuaDevice
		}
		// SELECT ...
		sb := sq.
			Select(
				"dahua_devices.*",
			).
			From("dahua_devices")
		// ORDER BY
		switch req.Sort.GetField() {
		case "name":
			sb = sb.OrderBy(decodeOrderSQL("name", req.Sort.GetOrder()))
		case "url":
			sb = sb.OrderBy(decodeOrderSQL("url", req.Sort.GetOrder()))
		case "createdAt":
			sb = sb.OrderBy(decodeOrderSQL("created_at", req.Sort.GetOrder()))
		}
		// OFFSET ...
		sb = sb.
			Offset(uint64(page.Offset())).
			Limit(uint64(page.Limit()))

		rows, scanner, err := ssq.QueryRows(ctx, a.db, sb)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var items []*rpc.GetAdminDevicesPageResp_Device
		for rows.Next() {
			err := scanner.Scan(&row)
			if err != nil {
				return nil, err
			}

			items = append(items, &rpc.GetAdminDevicesPageResp_Device{
				Id:             row.ID,
				Name:           row.Name,
				Url:            row.Url.String(),
				Username:       row.Username,
				Disabled:       row.DisabledAt.Valid,
				DisabledAtTime: timestamppb.New(row.DisabledAt.Time.Time),
				CreatedAtTime:  timestamppb.New(row.CreatedAt.Time),
			})
		}

		return items, nil
	}()
	if err != nil {
		return nil, err
	}

	count, err := func() (int64, error) {
		var row struct{ Count int64 }
		return row.Count, ssq.QueryOne(ctx, a.db, &row, sq.
			Select("COUNT(*) AS count").
			From("dahua_devices"))
	}()
	if err != nil {
		return nil, err
	}

	return &rpc.GetAdminDevicesPageResp{
		Items:      items,
		PageResult: encodePagePaginationResult(page.Result(int(count))),
		Sort:       req.Sort,
	}, nil
}

func (a *Admin) GetAdminDevicesIDPage(ctx context.Context, req *rpc.GetAdminDevicesIDPageReq) (*rpc.GetAdminDevicesIDPageResp, error) {
	v, err := dahua.GetDevice(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &rpc.GetAdminDevicesIDPageResp{
		Device: &rpc.GetAdminDevicesIDPageResp_Device{
			Id:             v.ID,
			Name:           v.Name,
			Url:            v.Url.String(),
			Username:       v.Username,
			Disabled:       v.DisabledAt.Valid,
			Location:       v.Location.String(),
			CreatedAtTime:  timestamppb.New(v.CreatedAt.Time),
			UpdatedAtTime:  timestamppb.New(v.UpdatedAt.Time),
			DisabledAtTime: timestamppb.New(v.DisabledAt.Time.Time),
			Features:       dahua.FeatureToStrings(v.Feature),
		},
	}, nil

}

func (a *Admin) GetDevice(ctx context.Context, req *rpc.GetDeviceReq) (*rpc.GetDeviceResp, error) {
	v, err := dahua.GetDevice(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &rpc.GetDeviceResp{
		Id:       v.ID,
		Name:     v.Name,
		Url:      v.Url.String(),
		Username: v.Username,
		Location: v.Location.String(),
		Features: dahua.FeatureToStrings(v.Feature),
		Email:    v.Email.String,
	}, nil
}

func (a *Admin) CreateDevice(ctx context.Context, req *rpc.CreateDeviceReq) (*rpc.CreateDeviceResp, error) {
	urL, err := url.Parse(req.Url)
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation(req.Location)
	if err != nil {
		return nil, err
	}

	id, err := dahua.CreateDevice(ctx, dahua.CreateDeviceParams{
		Name:     req.Name,
		URL:      urL,
		Username: req.Username,
		Password: req.Password,
		Location: loc,
		Feature:  dahua.FeatureFromStrings(req.Features),
		Email:    req.Email,
	})
	if err != nil {
		if errs, ok := core.AsFieldErrors(err); ok {
			return nil, newInvalidArgument(errs,
				keymap("name", "Name"),
				keymap("url", "URL"),
				keymap("email", "Email"),
			)
		}
		return nil, err
	}

	return &rpc.CreateDeviceResp{
		Id: id,
	}, nil
}

func (a *Admin) UpdateDevice(ctx context.Context, req *rpc.UpdateDeviceReq) (*emptypb.Empty, error) {
	urL, err := url.Parse(req.Url)
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation(req.Location)
	if err != nil {
		return nil, err
	}

	err = dahua.UpdateDevice(ctx, dahua.UpdateDeviceParams{
		ID:          req.Id,
		Name:        req.Name,
		URL:         urL,
		Username:    req.Username,
		NewPassword: req.NewPassword,
		Location:    loc,
		Feature:     dahua.FeatureFromStrings(req.Features),
		Email:       req.Email,
	})
	if err != nil {
		if errs, ok := core.AsFieldErrors(err); ok {
			return nil, newInvalidArgument(errs,
				keymap("name", "Name"),
				keymap("url", "URL"),
				keymap("email", "Email"),
			)
		}
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) DeleteDevice(ctx context.Context, req *rpc.DeleteDeviceReq) (*emptypb.Empty, error) {
	for _, id := range req.Ids {
		err := dahua.DeleteDevice(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}

func (a *Admin) SetDeviceDisable(ctx context.Context, req *rpc.SetDeviceDisableReq) (*emptypb.Empty, error) {
	for _, v := range req.Items {
		err := dahua.UpdateDeviceDisabled(ctx, v.Id, v.Disable)
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}

// ---------- User

func (a *Admin) GetAdminUsersPage(ctx context.Context, req *rpc.GetAdminUsersPageReq) (*rpc.GetAdminUsersPageResp, error) {
	page := decodePagePagination(req.Page)

	items, err := func() ([]*rpc.GetAdminUsersPageResp_User, error) {
		var row struct {
			repo.User
			Admin bool
		}
		// SELECT ...
		sb := sq.
			Select(
				"users.*",
				"admins.user_id IS NOT NULL as 'admin'",
			).
			From("users").
			LeftJoin("admins ON admins.user_id = users.id")
		// ORDER BY
		switch req.Sort.GetField() {
		case "username":
			sb = sb.OrderBy(decodeOrderSQL("username", req.Sort.GetOrder()))
		case "email":
			sb = sb.OrderBy(decodeOrderSQL("email", req.Sort.GetOrder()))
		case "createdAt":
			sb = sb.OrderBy(decodeOrderSQL("users.created_at", req.Sort.GetOrder()))
		default:
			sb = sb.OrderBy("admin DESC")
		}
		// OFFSET ...
		sb = sb.
			Offset(uint64(page.Offset())).
			Limit(uint64(page.Limit()))

		rows, scanner, err := ssq.QueryRows(ctx, a.db, sb)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var items []*rpc.GetAdminUsersPageResp_User
		for rows.Next() {
			err := scanner.Scan(&row)
			if err != nil {
				return nil, err
			}

			items = append(items, &rpc.GetAdminUsersPageResp_User{
				Id:             row.ID,
				Username:       row.Username,
				Email:          row.Email,
				Disabled:       row.DisabledAt.Valid,
				Admin:          row.Admin,
				DisabledAtTime: timestamppb.New(row.DisabledAt.Time.Time),
				CreatedAtTime:  timestamppb.New(row.CreatedAt.Time),
			})
		}

		return items, nil
	}()
	if err != nil {
		return nil, err
	}

	count, err := func() (int64, error) {
		var row struct{ Count int64 }
		return row.Count, ssq.QueryOne(ctx, a.db, &row, sq.
			Select("COUNT(*) AS count").
			From("users"))
	}()
	if err != nil {
		return nil, err
	}

	return &rpc.GetAdminUsersPageResp{
		Items:      items,
		PageResult: encodePagePaginationResult(page.Result(int(count))),
		Sort:       req.Sort,
	}, nil

}

func (a *Admin) CreateUser(ctx context.Context, req *rpc.CreateUserReq) (*emptypb.Empty, error) {
	cfg, err := system.GetConfig()
	if err != nil {
		return nil, err
	}

	_, err = auth.CreateUser(ctx, cfg, auth.CreateUserParams{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Admin:    req.Admin,
		Disabled: req.Disabled,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) GetUser(ctx context.Context, req *rpc.GetUserReq) (*rpc.GetUserResp, error) {
	v, err := a.db.C().AuthGetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &rpc.GetUserResp{
		Id:       v.ID,
		Username: v.Username,
		Email:    v.Email,
	}, nil
}

func (a *Admin) UpdateUser(ctx context.Context, req *rpc.UpdateUserReq) (*emptypb.Empty, error) {
	err := auth.UpdateUser(ctx, auth.UpdateUserParams{
		ID:       req.Id,
		Email:    req.Email,
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) DeleteUser(ctx context.Context, req *rpc.DeleteUserReq) (*emptypb.Empty, error) {
	for _, id := range req.Ids {
		err := auth.DeleteUser(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) SetUserDisable(ctx context.Context, req *rpc.SetUserDisableReq) (*emptypb.Empty, error) {
	for _, item := range req.Items {
		err := auth.UpdateUserDisabled(ctx, item.Id, item.Disable)
		if err != nil {
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) SetUserAdmin(ctx context.Context, req *rpc.SetUserAdminReq) (*emptypb.Empty, error) {
	err := auth.UpdateUserAdmin(ctx, req.Id, req.Admin)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) ResetUserPassword(ctx context.Context, req *rpc.ResetUserPasswordReq) (*emptypb.Empty, error) {
	session := useAuthSession(ctx)

	if err := auth.UpdateUserPassword(ctx, auth.UpdateUserPasswordParams{
		UserID:           req.Id,
		OldPasswordSkip:  true,
		NewPassword:      req.NewPassword,
		CurrentSessionID: session.SessionID,
	}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ---------- Group

func (a *Admin) GetAdminGroupsPage(ctx context.Context, req *rpc.GetAdminGroupsPageReq) (*rpc.GetAdminGroupsPageResp, error) {
	page := decodePagePagination(req.Page)

	items, err := func() ([]*rpc.GetAdminGroupsPageResp_Group, error) {
		var row struct {
			repo.Group
			UserCount int64
		}
		// SELECT ...
		sb := sq.
			Select(
				"groups.*",
				"COUNT(group_users.group_id) AS user_count",
			).
			From("groups").
			LeftJoin("group_users ON group_users.group_id = groups.id").
			GroupBy("groups.id")
		// ORDER BY
		switch req.Sort.GetField() {
		case "name":
			sb = sb.OrderBy(decodeOrderSQL("name", req.Sort.GetOrder()))
		case "userCount":
			sb = sb.OrderBy(decodeOrderSQL("user_count", req.Sort.GetOrder()))
		case "createdAt":
			sb = sb.OrderBy(decodeOrderSQL("groups.created_at", req.Sort.GetOrder()))
		}
		// OFFSET ...
		sb = sb.
			Offset(uint64(page.Offset())).
			Limit(uint64(page.Limit()))

		rows, scanner, err := ssq.QueryRows(ctx, a.db, sb)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var items []*rpc.GetAdminGroupsPageResp_Group
		for rows.Next() {
			err := scanner.Scan(&row)
			if err != nil {
				return nil, err
			}

			items = append(items, &rpc.GetAdminGroupsPageResp_Group{
				Id:             row.ID,
				Name:           row.Name,
				UserCount:      row.UserCount,
				Disabled:       !row.DisabledAt.Valid,
				DisabledAtTime: timestamppb.New(row.DisabledAt.Time.Time),
				CreatedAtTime:  timestamppb.New(row.CreatedAt.Time),
			})
		}

		return items, nil
	}()
	if err != nil {
		return nil, err
	}

	count, err := func() (int64, error) {
		var row struct{ Count int64 }
		return row.Count, ssq.QueryOne(ctx, a.db, &row, sq.
			Select("COUNT(*) AS count").
			From("groups"))
	}()
	if err != nil {
		return nil, err
	}

	return &rpc.GetAdminGroupsPageResp{
		Items:      items,
		PageResult: encodePagePaginationResult(page.Result(int(count))),
		Sort:       req.Sort,
	}, nil
}

func (a *Admin) GetAdminGroupsIDPage(ctx context.Context, req *rpc.GetAdminGroupsIDPageReq) (*rpc.GetAdminGroupsIDPageResp, error) {
	dbUsers, err := a.db.C().AuthListUsersByGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	users := make([]*rpc.GetAdminGroupsIDPageResp_User, 0, len(dbUsers))
	for _, v := range dbUsers {
		users = append(users, &rpc.GetAdminGroupsIDPageResp_User{
			Id:       v.ID,
			Username: v.Username,
		})
	}

	v, err := a.db.C().AuthGetGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &rpc.GetAdminGroupsIDPageResp{
		Group: &rpc.GetAdminGroupsIDPageResp_Group{
			Id:             v.ID,
			Name:           v.Name,
			Description:    v.Description,
			Disabled:       v.DisabledAt.Valid,
			DisabledAtTime: timestamppb.New(v.DisabledAt.Time.Time),
			CreatedAtTime:  timestamppb.New(v.CreatedAt.Time),
			UpdatedAtTime:  timestamppb.New(v.UpdatedAt.Time),
		},
		Users: users,
	}, nil
}

func (a *Admin) GetGroup(ctx context.Context, req *rpc.GetGroupReq) (*rpc.GetGroupResp, error) {
	v, err := a.db.C().AuthGetGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &rpc.GetGroupResp{
		Id:          v.ID,
		Name:        v.Name,
		Description: v.Description,
	}, nil
}

func (a *Admin) CreateGroup(ctx context.Context, req *rpc.CreateGroupReq) (*rpc.CreateGroupResp, error) {
	id, err := auth.CreateGroup(ctx, auth.CreateGroupParams{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &rpc.CreateGroupResp{
		Id: id,
	}, nil
}

func (a *Admin) UpdateGroup(ctx context.Context, req *rpc.UpdateGroupReq) (*emptypb.Empty, error) {
	err := auth.UpdateGroup(ctx, auth.UpdateGroupParams{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) DeleteGroup(ctx context.Context, req *rpc.DeleteGroupReq) (*emptypb.Empty, error) {
	for _, id := range req.Ids {
		err := auth.DeleteGroup(ctx, id)
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}

func (a *Admin) GetAdminEventsPage(ctx context.Context, _ *emptypb.Empty) (*rpc.GetAdminEventsPageResp, error) {
	eventCount, err := dahua.CountEvents(ctx)
	if err != nil {
		return nil, err
	}

	return &rpc.GetAdminEventsPageResp{
		EventCount: eventCount,
	}, nil
}

func (a *Admin) SetGroupDisable(ctx context.Context, req *rpc.SetGroupDisableReq) (*emptypb.Empty, error) {
	for _, item := range req.Items {
		err := auth.UpdateGroupDisable(ctx, item.Id, item.Disable)
		if err != nil {
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}

func (a *Admin) CreateEventRule(ctx context.Context, req *rpc.CreateEventRuleReq) (*rpc.CreateEventRuleResp, error) {
	id, err := dahua.CreateEventRule(ctx, repo.DahuaCreateEventRuleParams{
		Code:       req.Code,
		IgnoreDb:   req.IgnoreDb,
		IgnoreLive: req.IgnoreLive,
		IgnoreMqtt: req.IgnoreMqtt,
	})
	if err != nil {
		if errs, ok := core.AsFieldErrors(err); ok {
			return nil, newInvalidArgument(errs,
				keymap("code", "Code"),
			)
		}
		return nil, err
	}

	return &rpc.CreateEventRuleResp{
		Id: id,
	}, nil
}

func (a *Admin) UpdateEventRule(ctx context.Context, req *rpc.UpdateEventRuleReq) (*emptypb.Empty, error) {
	for _, v := range req.Items {
		err := dahua.UpdateEventRule(ctx, repo.DahuaUpdateEventRuleParams{
			Code:       v.Code,
			IgnoreDb:   v.IgnoreDb,
			IgnoreLive: v.IgnoreLive,
			IgnoreMqtt: v.IgnoreMqtt,
			ID:         v.Id,
		})
		if err != nil {
			if core.IsNotFound(err) {
				continue
			}
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) ListEventRules(ctx context.Context, _ *emptypb.Empty) (*rpc.ListEventRulesResp, error) {
	v, err := dahua.ListEventRules(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]*rpc.ListEventRulesResp_Item, 0, len(v))
	for _, v := range v {
		items = append(items, &rpc.ListEventRulesResp_Item{
			Id:         v.ID,
			Code:       v.Code,
			IgnoreDb:   v.IgnoreDb,
			IgnoreLive: v.IgnoreLive,
			IgnoreMqtt: v.IgnoreMqtt,
		})
	}

	return &rpc.ListEventRulesResp{
		Items: items,
	}, nil
}

func (a *Admin) DeleteEventRules(ctx context.Context, req *rpc.DeleteEventRulesReq) (*emptypb.Empty, error) {
	for _, id := range req.Ids {
		if err := dahua.DeleteEventRule(ctx, id); err != nil {
			if core.IsNotFound(err) {
				continue
			}
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}

func (*Admin) ListLocations(context.Context, *emptypb.Empty) (*rpc.ListLocationsResp, error) {
	return &rpc.ListLocationsResp{
		Locations: core.Locations,
	}, nil
}

var listDeviceFeaturesResp *rpc.ListDeviceFeaturesResp

func init() {
	features := make([]*rpc.ListDeviceFeaturesResp_Item, 0, len(dahua.FeatureList))
	for _, v := range dahua.FeatureList {
		features = append(features, &rpc.ListDeviceFeaturesResp_Item{
			Name:        v.Name,
			Value:       v.Value,
			Description: v.Description,
		})
	}
	listDeviceFeaturesResp = &rpc.ListDeviceFeaturesResp{Features: features}
}

func (*Admin) ListDeviceFeatures(context.Context, *emptypb.Empty) (*rpc.ListDeviceFeaturesResp, error) {
	return listDeviceFeaturesResp, nil
}

func (a *Admin) UpdateConfig(cfg context.Context, req *rpc.UpdateConfigReq) (*emptypb.Empty, error) {
	err := system.UpdateConfig(system.UpdateConfigParams{
		SiteName:     req.SiteName,
		EnableSignUp: req.EnableSignUp,
	})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) DeleteEvents(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := dahua.DeleteEvents(ctx)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
