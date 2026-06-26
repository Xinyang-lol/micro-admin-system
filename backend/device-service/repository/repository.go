package repository

import (
	"context"
	"database/sql"
	"strings"

	pb "micro-admin-system/backend/proto/gen"
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
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 10
	}
	return page, pageSize
}

func (r *Repository) ListDevices(ctx context.Context, req *pb.DeviceListRequest) (*pb.DeviceListResponse, error) {
	page, pageSize := normalizePage(req.Page, req.PageSize)
	conditions := []string{"1 = 1"}
	args := []any{}
	if req.Keyword != "" {
		conditions = append(conditions, "(d.name LIKE ? OR d.code LIKE ? OR d.location LIKE ?)")
		keyword := "%" + req.Keyword + "%"
		args = append(args, keyword, keyword, keyword)
	}
	if req.TypeId > 0 {
		conditions = append(conditions, "d.type_id = ?")
		args = append(args, req.TypeId)
	}
	if req.Status != "" {
		conditions = append(conditions, "d.status = ?")
		args = append(args, req.Status)
	}
	where := strings.Join(conditions, " AND ")

	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM device d WHERE "+where, args...).Scan(&total); err != nil {
		return nil, err
	}

	queryArgs := append([]any{}, args...)
	queryArgs = append(queryArgs, pageSize, (page-1)*pageSize)
	rows, err := r.db.QueryContext(ctx, `
		SELECT d.id, d.name, d.code, d.type_id, COALESCE(t.name, ''), d.status, d.location, d.remark, DATE_FORMAT(d.created_at, '%Y-%m-%d %H:%i:%s')
		FROM device d
		LEFT JOIN device_type t ON t.id = d.type_id
		WHERE `+where+`
		ORDER BY d.id DESC LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &pb.DeviceListResponse{Total: total}
	for rows.Next() {
		item := &pb.Device{}
		if err := rows.Scan(&item.Id, &item.Name, &item.Code, &item.TypeId, &item.TypeName, &item.Status, &item.Location, &item.Remark, &item.CreatedAt); err != nil {
			return nil, err
		}
		resp.Items = append(resp.Items, item)
	}
	return resp, rows.Err()
}

func (r *Repository) CreateDevice(ctx context.Context, req *pb.DeviceSaveRequest) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO device(name, code, type_id, status, location, remark)
		VALUES(?, ?, ?, ?, ?, ?)`,
		req.Name, req.Code, req.TypeId, req.Status, req.Location, req.Remark)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) UpdateDevice(ctx context.Context, req *pb.DeviceSaveRequest) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE device SET name = ?, code = ?, type_id = ?, status = ?, location = ?, remark = ?, updated_at = NOW()
		WHERE id = ?`,
		req.Name, req.Code, req.TypeId, req.Status, req.Location, req.Remark, req.Id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) DeleteDevice(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM device WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) CountAll(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM device`).Scan(&total)
	return total, err
}

func (r *Repository) CountByStatus(ctx context.Context, status string) (int64, error) {
	var total int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM device WHERE status = ?`, status).Scan(&total)
	return total, err
}

func (r *Repository) CountByType(ctx context.Context) (map[string]int64, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT COALESCE(t.name, '未分类'), COUNT(d.id)
		FROM device d LEFT JOIN device_type t ON t.id = d.type_id
		GROUP BY COALESCE(t.name, '未分类')`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := map[string]int64{}
	for rows.Next() {
		var name string
		var count int64
		if err := rows.Scan(&name, &count); err != nil {
			return nil, err
		}
		result[name] = count
	}
	return result, rows.Err()
}

func (r *Repository) ListDeviceTypes(ctx context.Context) (*pb.DeviceTypeListResponse, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, code, remark FROM device_type ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	resp := &pb.DeviceTypeListResponse{}
	for rows.Next() {
		item := &pb.DeviceType{}
		if err := rows.Scan(&item.Id, &item.Name, &item.Code, &item.Remark); err != nil {
			return nil, err
		}
		resp.Items = append(resp.Items, item)
	}
	return resp, rows.Err()
}

func (r *Repository) CreateDeviceType(ctx context.Context, req *pb.DeviceTypeSaveRequest) (int64, error) {
	res, err := r.db.ExecContext(ctx, `INSERT INTO device_type(name, code, remark) VALUES(?, ?, ?)`, req.Name, req.Code, req.Remark)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) UpdateDeviceType(ctx context.Context, req *pb.DeviceTypeSaveRequest) error {
	res, err := r.db.ExecContext(ctx, `UPDATE device_type SET name = ?, code = ?, remark = ?, updated_at = NOW() WHERE id = ?`, req.Name, req.Code, req.Remark, req.Id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) DeleteDeviceType(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM device_type WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
