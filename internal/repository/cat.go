package repository

import (
	ctx2 "awesomeProject/internal/ctx"
	"awesomeProject/internal/db"
	"awesomeProject/internal/types"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type CatRepository struct {
	db *sql.DB
}

func (repo *CatRepository) checkDBInRepository() bool {
	if repo.db == nil {
		log.Println("DB pointer in repo is nil, try to get new one")

		newDB := db.GetDB()

		if newDB == nil {
			log.Println("DB pointer is still nil, check DB health")

			return false
		}

		repo.db = newDB
	}

	return true
}

func executeWithTimeout[T interface{}](cb func(args ...interface{}) T) (*T, error) {
	result := make(chan T, 1)
	rootCtx := ctx2.GetRootCtx()
	timeoutCtx, cancel := context.WithTimeout(rootCtx.Context, 5*time.Second)
	defer cancel()

	go func() { result <- cb() }()

	select {
	case r := <-result:
		return &r, nil
	case <-timeoutCtx.Done():
		return nil, errors.New("timeout expired")
	}
}

func NewCatRepository(db *sql.DB) *CatRepository {
	return &CatRepository{db}
}

func (repo *CatRepository) GetAvailableCats() []string {
	isDBAlive := repo.checkDBInRepository()

	if !isDBAlive {
		return nil
	}

	rows, err := repo.db.Query(`SELECT name FROM tag`)

	if err != nil {
		return nil
	}

	response := make([]string, 0)

	for rows.Next() {
		var cat string
		err = rows.Scan(&cat)

		if err != nil {
			return nil
		}

		response = append(response, cat)
	}

	return response
}

func (repo *CatRepository) GetCatByTag(tag string) *types.CatPhoto {
	var catPhoto types.CatPhoto
	isDBAlive := repo.checkDBInRepository()

	if !isDBAlive {
		return nil
	}

	row := repo.db.QueryRow(`SELECT p.link
    FROM photo_tag pt, photo p, tag t
    WHERE pt.tag_id = t.id
    AND t.name = $1 AND p.id = pt.photo_id
    ORDER BY random()
    LIMIT 1
    `, tag)

	if row == nil {
		return nil
	}

	err := row.Scan(&catPhoto.Link)

	if err != nil {
		return nil
	}

	return &catPhoto
}

func (repo *CatRepository) StoreCat(tag string, link string) error {
	var photoId, tagId int

	rootCtx := ctx2.GetRootCtx()
	needToInsertNewTag := false
	isDBAlive := repo.checkDBInRepository()

	if !isDBAlive {
		return errors.New("DB not alive")
	}

	ctx, cancel := context.WithTimeout(rootCtx.Context, 5*time.Second)
	defer cancel()

	tx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		return errors.New("failed to start transaction")
	}

	row := tx.QueryRowContext(ctx, `INSERT INTO photo ("link") VALUES ($1) RETURNING id`, link)
	err = row.Scan(&photoId)

	if err != nil {
		tx.Rollback()

		return errors.New("failed to insert link into photos table")
	}

	row = tx.QueryRowContext(ctx, `SELECT tag.id FROM tag WHERE tag.name = ($1)`, tag)
	err = row.Scan(&tagId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			needToInsertNewTag = true
		} else {
			tx.Rollback()

			return errors.New("failed to select tag from tags table")
		}
	}

	if needToInsertNewTag {
		row = tx.QueryRowContext(ctx, `INSERT INTO tag ("name") VALUES ($1) returning id`, tag)
		err = row.Scan(&tagId)

		if err != nil {
			tx.Rollback()

			return errors.New("failed to insert tag into tags table")
		}
	}

	row = tx.QueryRowContext(ctx, `INSERT INTO photo_tag ("photo_id", "tag_id") VALUES ($1, $2)`, photoId, tagId)
	err = row.Scan()

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()

		return errors.New("failed to insert info into photo_tag table")
	}

	tx.Commit()

	return nil
}
