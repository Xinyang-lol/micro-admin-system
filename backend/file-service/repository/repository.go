package repository

import (
	"context"
	"database/sql"

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
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	return page, pageSize
}

func (r *Repository) Create(ctx context.Context, req *pb.UploadFileMetaRequest) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
		INSERT INTO file_info(original_name, stored_name, path, size, content_type, uploader_id)
		VALUES(?, ?, ?, ?, ?, ?)`,
		req.OriginalName, req.StoredName, req.Path, req.Size, req.ContentType, req.UploaderId)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) List(ctx context.Context, req *pb.FileListRequest) (*pb.FileListResponse, error) {
	page, pageSize := normalizePage(req.Page, req.PageSize)
	keyword := "%" + req.Keyword + "%"
	var total int64
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM file_info
		WHERE deleted_at IS NULL AND (? = '%%' OR original_name LIKE ? OR stored_name LIKE ?)`,
		keyword, keyword, keyword).Scan(&total); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, original_name, stored_name, path, size, content_type, uploader_id, DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
		FROM file_info
		WHERE deleted_at IS NULL AND (? = '%%' OR original_name LIKE ? OR stored_name LIKE ?)
		ORDER BY id DESC LIMIT ? OFFSET ?`,
		keyword, keyword, keyword, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	resp := &pb.FileListResponse{Total: total}
	for rows.Next() {
		item := &pb.FileInfo{}
		if err := rows.Scan(&item.Id, &item.OriginalName, &item.StoredName, &item.Path, &item.Size, &item.ContentType, &item.UploaderId, &item.CreatedAt); err != nil {
			return nil, err
		}
		resp.Items = append(resp.Items, item)
	}
	return resp, rows.Err()
}

func (r *Repository) Get(ctx context.Context, id int64) (*pb.FileInfo, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, original_name, stored_name, path, size, content_type, uploader_id, DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s')
		FROM file_info WHERE id = ? AND deleted_at IS NULL`, id)
	item := &pb.FileInfo{}
	if err := row.Scan(&item.Id, &item.OriginalName, &item.StoredName, &item.Path, &item.Size, &item.ContentType, &item.UploaderId, &item.CreatedAt); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `UPDATE file_info SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
