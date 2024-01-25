package rpcserver

import (
	"context"

	"github.com/ItsNotGoodName/ipcmanview/internal/auth"
	"github.com/ItsNotGoodName/ipcmanview/internal/models"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/pkg/ssq"
	"github.com/ItsNotGoodName/ipcmanview/rpc"
	sq "github.com/Masterminds/squirrel"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewAdmin(db repo.DB) *Admin {
	return &Admin{
		db: db,
	}
}

type Admin struct {
	db repo.DB
}

// ---------- Device

func (*Admin) GetAdminDevicesPage(context.Context, *rpc.GetAdminDevicesPageReq) (*rpc.GetAdminDevicesPageResp, error) {
	return nil, errNotImplemented
}

func (*Admin) GetDevice(context.Context, *rpc.GetDeviceReq) (*rpc.GetDeviceResp, error) {
	return nil, errNotImplemented
}

func (*Admin) CreateDevice(context.Context, *rpc.CreateDeviceReq) (*rpc.CreateDeviceResp, error) {
	return nil, errNotImplemented
}

func (*Admin) UpdateDevice(context.Context, *rpc.UpdateDeviceReq) (*emptypb.Empty, error) {
	return nil, errNotImplemented
}

func (*Admin) DeleteDevice(context.Context, *rpc.DeleteDeviceReq) (*emptypb.Empty, error) {
	return nil, errNotImplemented
}

// ---------- User

func (a *Admin) GetAdminUsersPage(ctx context.Context, req *rpc.GetAdminUsersPageReq) (*rpc.GetAdminUsersPageResp, error) {
	page := parsePagePagination(req.Page)

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
		case "name":
			sb = sb.OrderBy(convertOrderToSQL("username", req.Sort.GetOrder()))
		case "email":
			sb = sb.OrderBy(convertOrderToSQL("email", req.Sort.GetOrder()))
		case "createdAt":
			sb = sb.OrderBy(convertOrderToSQL("users.created_at", req.Sort.GetOrder()))
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
				DisabledAtTime: timestamppb.New(row.DisabledAt.Time),
				CreatedAtTime:  timestamppb.New(row.CreatedAt.Time),
			})
		}

		return items, nil
	}()
	if err != nil {
		return nil, check(err)
	}

	count, err := func() (int64, error) {
		var row struct{ Count int64 }
		return row.Count, ssq.QueryOne(ctx, a.db, &row, sq.
			Select("COUNT(*) AS count").
			From("groups"))
	}()
	if err != nil {
		return nil, check(err)
	}

	return &rpc.GetAdminUsersPageResp{
		Items:      items,
		PageResult: convertPagePaginationResult(page.Result(int(count))),
		Sort:       req.Sort,
	}, nil

}

// ---------- Group

func (a *Admin) GetAdminGroupsPage(ctx context.Context, req *rpc.GetAdminGroupsPageReq) (*rpc.GetAdminGroupsPageResp, error) {
	page := parsePagePagination(req.Page)

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
			sb = sb.OrderBy(convertOrderToSQL("name", req.Sort.GetOrder()))
		case "userCount":
			sb = sb.OrderBy(convertOrderToSQL("user_count", req.Sort.GetOrder()))
		case "createdAt":
			sb = sb.OrderBy(convertOrderToSQL("groups.created_at", req.Sort.GetOrder()))
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
				Disabled:       row.DisabledAt.Valid,
				DisabledAtTime: timestamppb.New(row.DisabledAt.Time),
				CreatedAtTime:  timestamppb.New(row.CreatedAt.Time),
			})
		}

		return items, nil
	}()
	if err != nil {
		return nil, check(err)
	}

	count, err := func() (int64, error) {
		var row struct{ Count int64 }
		return row.Count, ssq.QueryOne(ctx, a.db, &row, sq.
			Select("COUNT(*) AS count").
			From("groups"))
	}()
	if err != nil {
		return nil, check(err)
	}

	return &rpc.GetAdminGroupsPageResp{
		Items:      items,
		PageResult: convertPagePaginationResult(page.Result(int(count))),
		Sort:       req.Sort,
	}, nil
}

func (a *Admin) GetAdminGroupIDPage(ctx context.Context, req *rpc.GetAdminGroupIDPageReq) (*rpc.GetAdminGroupIDPageResp, error) {
	dbGroup, err := a.db.GetGroup(ctx, req.Id)
	if err != nil {
		return nil, check(err)
	}
	return &rpc.GetAdminGroupIDPageResp{
		Group: &rpc.GetAdminGroupIDPageResp_Group{
			Id:             dbGroup.ID,
			Name:           dbGroup.Name,
			Description:    dbGroup.Description,
			Disabled:       dbGroup.DisabledAt.Valid,
			DisabledAtTime: timestamppb.New(dbGroup.DisabledAt.Time),
			CreatedAtTime:  timestamppb.New(dbGroup.CreatedAt.Time),
			UpdatedAtTime:  timestamppb.New(dbGroup.UpdatedAt.Time),
		},
	}, nil
}

func (a *Admin) GetGroup(ctx context.Context, req *rpc.GetGroupReq) (*rpc.GetGroupResp, error) {
	dbGroup, err := a.db.GetGroup(ctx, req.Id)
	if err != nil {
		return nil, check(err)
	}
	return &rpc.GetGroupResp{
		Id: req.Id,
		Model: &rpc.GroupModel{
			Name:        dbGroup.Name,
			Description: dbGroup.Description,
		},
	}, nil
}

func convertCreateUpdateGroupError(msg string, err error) error {
	if errs, ok := asValidationErrors(err); ok {
		return NewError(err, msg).Validation(errs, [][2]string{
			{"name", "Name"},
			{"description", "Description"},
		})
	}

	if constraintErr, ok := asConstraintError(err); ok {
		return NewError(err, msg).Constraint(constraintErr, [][3]string{
			{"name", "groups.name", "Name already taken."},
		})
	}

	return check(err)
}

func (a *Admin) CreateGroup(ctx context.Context, req *rpc.CreateGroupReq) (*rpc.CreateGroupResp, error) {
	id, err := auth.CreateGroup(ctx, a.db, models.Group{
		Name:        req.Model.GetName(),
		Description: req.Model.GetDescription(),
	})
	if err != nil {
		return nil, convertCreateUpdateGroupError("Failed to create group.", err)
	}

	return &rpc.CreateGroupResp{
		Id: id,
	}, nil
}

func (a *Admin) UpdateGroup(ctx context.Context, req *rpc.UpdateGroupReq) (*emptypb.Empty, error) {
	dbGroup, err := a.db.GetGroup(ctx, req.Id)
	if err != nil {
		return nil, check(err)
	}
	group := dbGroup.Convert()

	group.Name = req.Model.GetName()
	group.Description = req.Model.GetDescription()

	_, err = auth.UpdateGroup(ctx, a.db, group)
	if err != nil {
		return nil, convertCreateUpdateGroupError("Failed to update group.", err)
	}

	return &emptypb.Empty{}, nil
}

func (a *Admin) DeleteGroup(ctx context.Context, req *rpc.DeleteGroupReq) (*emptypb.Empty, error) {
	err := auth.DeleteGroup(ctx, a.db, req.Id)
	if err != nil {
		return nil, check(err)
	}
	return &emptypb.Empty{}, nil
}

func (a *Admin) SetGroupDisable(ctx context.Context, req *rpc.SetGroupDisableReq) (*emptypb.Empty, error) {
	if req.Disable {
		err := auth.DisableGroup(ctx, a.db, req.Id)
		if err != nil {
			return nil, check(err)
		}
	} else {
		err := auth.EnableGroup(ctx, a.db, req.Id)
		if err != nil {
			return nil, check(err)
		}
	}
	return &emptypb.Empty{}, nil
}

func (a *Admin) SetUserDisable(ctx context.Context, req *rpc.SetUserDisableReq) (*emptypb.Empty, error) {
	if req.Disable {
		err := auth.DisableUser(ctx, a.db, req.Id)
		if err != nil {
			return nil, check(err)
		}
	} else {
		err := auth.EnableUser(ctx, a.db, req.Id)
		if err != nil {
			return nil, check(err)
		}
	}
	return &emptypb.Empty{}, nil
}
