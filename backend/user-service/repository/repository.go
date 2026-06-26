package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"micro-admin-system/backend/proto/gen"
	"micro-admin-system/backend/user-service/model"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func normalizePage(page int32, pageSize int32) (int32, int32) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	return page, pageSize
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, username, password, nickname, email, phone, status, COALESCE(dept_id, 0), DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
		FROM sys_user WHERE username = ? AND deleted_at IS NULL`, username)
	u := &model.User{}
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Nickname, &u.Email, &u.Phone, &u.Status, &u.DeptID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *Repository) GetUserInfo(ctx context.Context, userID int64) (*pb.UserInfo, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, username, nickname, email, phone, status, COALESCE(dept_id, 0)
		FROM sys_user WHERE id = ? AND deleted_at IS NULL`, userID)
	info := &pb.UserInfo{}
	if err := row.Scan(&info.Id, &info.Username, &info.Nickname, &info.Email, &info.Phone, &info.Status, &info.DeptId); err != nil {
		return nil, err
	}
	roles, err := r.GetUserRoleCodes(ctx, userID)
	if err != nil {
		return nil, err
	}
	permissions, err := r.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}
	info.Roles = roles
	info.Permissions = permissions
	return info, nil
}

func (r *Repository) GetUserRoleCodes(ctx context.Context, userID int64) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT r.code FROM sys_role r
		JOIN sys_user_role ur ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.status = 1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		result = append(result, code)
	}
	return result, rows.Err()
}

func (r *Repository) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT DISTINCT m.permission FROM sys_menu m
		JOIN sys_role_menu rm ON rm.menu_id = m.id
		JOIN sys_user_role ur ON ur.role_id = rm.role_id
		WHERE ur.user_id = ? AND m.permission <> '' AND m.status = 1
		ORDER BY m.permission`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		result = append(result, permission)
	}
	return result, rows.Err()
}

func (r *Repository) HasPermission(ctx context.Context, userID int64, permission string) (bool, error) {
	roles, err := r.GetUserRoleCodes(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, role := range roles {
		if role == "admin" {
			return true, nil
		}
	}
	perms, err := r.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, p := range perms {
		if p == permission || p == "*:*:*" {
			return true, nil
		}
	}
	return false, nil
}

func (r *Repository) ListUsers(ctx context.Context, req *pb.ListRequest) (*pb.UserListResponse, error) {
	page, pageSize := normalizePage(req.Page, req.PageSize)
	keyword := "%" + req.Keyword + "%"
	var total int64
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM sys_user
		WHERE deleted_at IS NULL AND (? = '%%' OR username LIKE ? OR nickname LIKE ? OR phone LIKE ?)`,
		keyword, keyword, keyword, keyword).Scan(&total); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, username, nickname, email, phone, status, COALESCE(dept_id, 0), DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
		FROM sys_user
		WHERE deleted_at IS NULL AND (? = '%%' OR username LIKE ? OR nickname LIKE ? OR phone LIKE ?)
		ORDER BY id DESC LIMIT ? OFFSET ?`,
		keyword, keyword, keyword, keyword, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &pb.UserListResponse{Total: total}
	for rows.Next() {
		item := &pb.User{}
		if err := rows.Scan(&item.Id, &item.Username, &item.Nickname, &item.Email, &item.Phone, &item.Status, &item.DeptId, &item.CreatedAt); err != nil {
			return nil, err
		}
		roleIDs, err := r.GetUserRoleIDs(ctx, item.Id)
		if err != nil {
			return nil, err
		}
		item.RoleIds = roleIDs
		resp.Items = append(resp.Items, item)
	}
	return resp, rows.Err()
}

func (r *Repository) GetUserRoleIDs(ctx context.Context, userID int64) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT role_id FROM sys_user_role WHERE user_id = ? ORDER BY role_id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *Repository) CreateUser(ctx context.Context, req *pb.UserSaveRequest, passwordHash string) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, `
		INSERT INTO sys_user(username, password, nickname, email, phone, status, dept_id)
		VALUES(?, ?, ?, ?, ?, ?, ?)`,
		req.Username, passwordHash, req.Nickname, req.Email, req.Phone, req.Status, req.DeptId)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if err := assignUserRoles(ctx, tx, id, req.RoleIds); err != nil {
		return 0, err
	}
	return id, tx.Commit()
}

func (r *Repository) UpdateUser(ctx context.Context, req *pb.UserSaveRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, `
		UPDATE sys_user SET nickname = ?, email = ?, phone = ?, status = ?, dept_id = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL`,
		req.Nickname, req.Email, req.Phone, req.Status, req.DeptId, req.Id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	if req.RoleIds != nil {
		if err := assignUserRoles(ctx, tx, req.Id, req.RoleIds); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) DeleteUser(ctx context.Context, id int64) error {
	if id == 1 {
		return errors.New("admin user cannot be deleted")
	}
	res, err := r.db.ExecContext(ctx, `UPDATE sys_user SET deleted_at = NOW() WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) UpdateUserStatus(ctx context.Context, id int64, status int32) error {
	if id == 1 && status != 1 {
		return errors.New("admin user cannot be disabled")
	}
	res, err := r.db.ExecContext(ctx, `UPDATE sys_user SET status = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`, status, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) ResetPassword(ctx context.Context, id int64, passwordHash string) error {
	res, err := r.db.ExecContext(ctx, `UPDATE sys_user SET password = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL`, passwordHash, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) AssignUserRoles(ctx context.Context, id int64, roleIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := assignUserRoles(ctx, tx, id, roleIDs); err != nil {
		return err
	}
	return tx.Commit()
}

func assignUserRoles(ctx context.Context, tx *sql.Tx, userID int64, roleIDs []int64) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM sys_user_role WHERE user_id = ?`, userID); err != nil {
		return err
	}
	for _, roleID := range roleIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO sys_user_role(user_id, role_id) VALUES(?, ?)`, userID, roleID); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) ListRoles(ctx context.Context, req *pb.ListRequest) (*pb.RoleListResponse, error) {
	page, pageSize := normalizePage(req.Page, req.PageSize)
	keyword := "%" + req.Keyword + "%"
	var total int64
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM sys_role
		WHERE (? = '%%' OR name LIKE ? OR code LIKE ?)`, keyword, keyword, keyword).Scan(&total); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, code, remark, status FROM sys_role
		WHERE (? = '%%' OR name LIKE ? OR code LIKE ?)
		ORDER BY id DESC LIMIT ? OFFSET ?`, keyword, keyword, keyword, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	resp := &pb.RoleListResponse{Total: total}
	for rows.Next() {
		role := &pb.Role{}
		if err := rows.Scan(&role.Id, &role.Name, &role.Code, &role.Remark, &role.Status); err != nil {
			return nil, err
		}
		menuIDs, err := r.GetRoleMenuIDs(ctx, role.Id)
		if err != nil {
			return nil, err
		}
		role.MenuIds = menuIDs
		resp.Items = append(resp.Items, role)
	}
	return resp, rows.Err()
}

func (r *Repository) CreateRole(ctx context.Context, req *pb.RoleSaveRequest) (int64, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO sys_role(name, code, remark, status) VALUES(?, ?, ?, ?)`, req.Name, req.Code, req.Remark, req.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) UpdateRole(ctx context.Context, req *pb.RoleSaveRequest) error {
	res, err := r.db.ExecContext(ctx, `UPDATE sys_role SET name = ?, code = ?, remark = ?, status = ?, updated_at = NOW() WHERE id = ?`, req.Name, req.Code, req.Remark, req.Status, req.Id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) DeleteRole(ctx context.Context, id int64) error {
	if id == 1 {
		return errors.New("admin role cannot be deleted")
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM sys_user_role WHERE role_id = ?`, id); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM sys_role_menu WHERE role_id = ?`, id); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `DELETE FROM sys_role WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return tx.Commit()
}

func (r *Repository) GetRoleMenuIDs(ctx context.Context, roleID int64) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT menu_id FROM sys_role_menu WHERE role_id = ? ORDER BY menu_id`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *Repository) AssignRoleMenus(ctx context.Context, roleID int64, menuIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM sys_role_menu WHERE role_id = ?`, roleID); err != nil {
		return err
	}
	for _, menuID := range menuIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO sys_role_menu(role_id, menu_id) VALUES(?, ?)`, roleID, menuID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) ListMenus(ctx context.Context) (*pb.MenuTreeResponse, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, parent_id, name, path, component, permission, icon, type, sort
		FROM sys_menu WHERE status = 1 ORDER BY sort ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var menus []*pb.Menu
	for rows.Next() {
		menu := &pb.Menu{}
		if err := rows.Scan(&menu.Id, &menu.ParentId, &menu.Name, &menu.Path, &menu.Component, &menu.Permission, &menu.Icon, &menu.Type, &menu.Sort); err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &pb.MenuTreeResponse{Items: buildMenuTree(menus, 0)}, nil
}

func buildMenuTree(items []*pb.Menu, parentID int64) []*pb.Menu {
	var tree []*pb.Menu
	for _, item := range items {
		if item.ParentId == parentID {
			cp := *item
			cp.Children = buildMenuTree(items, item.Id)
			tree = append(tree, &cp)
		}
	}
	return tree
}

func (r *Repository) CreateMenu(ctx context.Context, req *pb.MenuSaveRequest) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO sys_menu(parent_id, name, path, component, permission, icon, type, sort, status)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, 1)`,
		req.ParentId, req.Name, req.Path, req.Component, req.Permission, req.Icon, req.Type, req.Sort)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) UpdateMenu(ctx context.Context, req *pb.MenuSaveRequest) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE sys_menu SET parent_id = ?, name = ?, path = ?, component = ?, permission = ?, icon = ?, type = ?, sort = ?, updated_at = NOW()
		WHERE id = ?`,
		req.ParentId, req.Name, req.Path, req.Component, req.Permission, req.Icon, req.Type, req.Sort, req.Id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) DeleteMenu(ctx context.Context, id int64) error {
	var childCount int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sys_menu WHERE parent_id = ?`, id).Scan(&childCount); err != nil {
		return err
	}
	if childCount > 0 {
		return fmt.Errorf("menu has %d child nodes", childCount)
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM sys_role_menu WHERE menu_id = ?`, id); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `DELETE FROM sys_menu WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return tx.Commit()
}

func (r *Repository) ListDepts(ctx context.Context) (*pb.DeptTreeResponse, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, parent_id, name, sort, status FROM sys_dept ORDER BY sort ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var depts []*pb.Dept
	for rows.Next() {
		dept := &pb.Dept{}
		if err := rows.Scan(&dept.Id, &dept.ParentId, &dept.Name, &dept.Sort, &dept.Status); err != nil {
			return nil, err
		}
		depts = append(depts, dept)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &pb.DeptTreeResponse{Items: buildDeptTree(depts, 0)}, nil
}

func buildDeptTree(items []*pb.Dept, parentID int64) []*pb.Dept {
	var tree []*pb.Dept
	for _, item := range items {
		if item.ParentId == parentID {
			cp := *item
			cp.Children = buildDeptTree(items, item.Id)
			tree = append(tree, &cp)
		}
	}
	return tree
}

func (r *Repository) CreateDept(ctx context.Context, req *pb.DeptSaveRequest) (int64, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO sys_dept(parent_id, name, sort, status) VALUES(?, ?, ?, ?)`, req.ParentId, req.Name, req.Sort, req.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) UpdateDept(ctx context.Context, req *pb.DeptSaveRequest) error {
	res, err := r.db.ExecContext(ctx, `UPDATE sys_dept SET parent_id = ?, name = ?, sort = ?, status = ?, updated_at = NOW() WHERE id = ?`, req.ParentId, req.Name, req.Sort, req.Status, req.Id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) DeleteDept(ctx context.Context, id int64) error {
	var childCount int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sys_dept WHERE parent_id = ?`, id).Scan(&childCount); err != nil {
		return err
	}
	if childCount > 0 {
		return fmt.Errorf("department has %d child nodes", childCount)
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM sys_dept WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
