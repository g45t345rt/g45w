package wallet_manager

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/g45t345rt/g45w/sc"
)

type TokenFolder struct {
	ID       int64
	Name     string
	ParentId sql.NullInt64
}

type Token struct {
	SCID              string
	Name              string
	MaxSupply         sql.NullInt64 // 1 is an NFT
	TotalSupply       sql.NullInt64
	Decimals          int64 // native dero decimals is 5
	StandardType      sc.SCType
	Metadata          sql.NullString
	IsFavorite        sql.NullBool
	ListOrderFavorite sql.NullInt64
	Image             sql.NullString
	Symbol            sql.NullString
	FolderId          sql.NullInt64
}

func initDatabaseTokens(db *sql.DB) error {
	dbTx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = dbTx.Exec(`
			CREATE TABLE IF NOT EXISTS token_folders (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name VARCHAR NOT NULL,
				parent_id INTEGER
			);

			CREATE TRIGGER IF NOT EXISTS delete_token_folders
			AFTER DELETE ON token_folders
			BEGIN
				DELETE FROM token_folders WHERE parent_id = OLD.id;
				DELETE FROM tokens WHERE folder_id = OLD.id;
			END;

			CREATE TABLE IF NOT EXISTS tokens (
				sc_id VARCHAR,
				name VARCHAR NOT NULL,
				max_supply BIGINT,
				total_supply BIGINT,
				decimals INTEGER NOT NULL,
				standard_type VARCHAR NOT NULL,
				metadata VARCHAR,
				is_favorite BOOL,
				list_order_favorite INTEGER,
				image VARCHAR,
				symbol VARCHAR,
				folder_id INTEGER,
				PRIMARY KEY(sc_id,folder_id)
			);
		`)
	if err != nil {
		return err
	}

	return handleDatabaseCommit(dbTx)
}

func (w *Wallet) GetTokenFolder(id int64) (*TokenFolder, error) {
	query := sq.Select("*").From("token_folders").Where(sq.Eq{"id": id})

	rows, err := query.RunWith(w.DB).Query()
	if err != nil {
		return nil, err
	}

	var folder *TokenFolder
	for rows.Next() {
		folder = &TokenFolder{}
		err = rows.Scan(
			&folder.ID,
			&folder.Name,
			&folder.ParentId,
		)
		if err != nil {
			return nil, err
		}
	}

	return folder, nil
}

func (w *Wallet) GetTokenFolderPath(id sql.NullInt64) (string, error) {
	if !id.Valid {
		return "root", nil
	}

	query := sq.Select("*").From("token_folders").Where(sq.Eq{"id": id})

	rows, err := query.RunWith(w.DB).Query()
	if err != nil {
		return "", err
	}

	var folder TokenFolder
	for rows.Next() {
		err = rows.Scan(
			&folder.ID,
			&folder.Name,
			&folder.ParentId,
		)
		if err != nil {
			return "", err
		}
	}

	parentName, err := w.GetTokenFolderPath(folder.ParentId)
	if err != nil {
		return "", err
	}

	return parentName + "/" + folder.Name, nil
}

func (w *Wallet) GetTokenFolderFolders(id sql.NullInt64) ([]TokenFolder, error) {
	query := sq.Select("*").From("token_folders").Where(sq.Eq{"parent_id": id})

	rows, err := query.RunWith(w.DB).Query()
	if err != nil {
		return nil, err
	}

	var folders []TokenFolder
	for rows.Next() {
		var folder TokenFolder
		err := rows.Scan(
			&folder.ID,
			&folder.Name,
			&folder.ParentId,
		)
		if err != nil {
			return folders, err
		}

		folders = append(folders, folder)
	}

	return folders, nil
}

func (w *Wallet) StoreFolderToken(folder *TokenFolder) error {
	// can't use UNIQUE() constraint because null does not count as towards uniqueness
	// https://stackoverflow.com/questions/22699409/sqlite-null-and-unique
	// we check manually instead

	query := sq.Select("COUNT(*)").From("token_folders").Where(sq.Eq{"name": folder.Name})

	if folder.ParentId.Valid {
		query = query.Where(sq.Eq{"parent_id": folder.ParentId.Int64})
	} else {
		query = query.Where(sq.Eq{"parent_id": nil})
	}

	row := query.RunWith(w.DB).QueryRow()

	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count >= 1 {
		return fmt.Errorf("folder already exists")
	}

	tx, err := w.DB.Begin()
	if err != nil {
		return err
	}

	result, err := tx.Exec(`
		INSERT INTO token_folders (name,parent_id)
		VALUES (?,?);
	`, folder.Name, folder.ParentId)
	if err != nil {
		return err
	}

	err = handleDatabaseCommit(tx)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	folder.ID = id
	return nil
}

func (w *Wallet) GetToken(scId string, folderId sql.NullInt64) (*Token, error) {
	query := sq.Select("*").From("tokens").Where(sq.Eq{"sc_id": scId})

	if folderId.Valid {
		query = query.Where(sq.Eq{"folder_id": folderId.Int64})
	} else {
		query = query.Where(sq.Eq{"folder_id": nil})
	}

	row := query.RunWith(w.DB).QueryRow()

	var token Token
	err := row.Scan(
		&token.SCID,
		&token.Name,
		&token.MaxSupply,
		&token.TotalSupply,
		&token.Decimals,
		&token.StandardType,
		&token.Metadata,
		&token.IsFavorite,
		&token.ListOrderFavorite,
		&token.Image,
		&token.Symbol,
		&token.FolderId,
	)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

type GetTokensParams struct {
	Descending bool
	OrderBy    string
	IsFavorite sql.NullBool
	FolderId   sql.NullInt64
	IsNFT      sql.NullBool
}

func (w *Wallet) GetTokens(params GetTokensParams) ([]Token, error) {
	query := sq.Select("*").From("tokens")

	if params.IsFavorite.Valid {
		query = query.Where(sq.Eq{"is_favorite": params.IsFavorite.Bool})
	}

	if params.FolderId.Valid {
		query = query.Where(sq.Eq{"folder_id": params.FolderId.Int64})
	}

	if params.IsNFT.Valid {
		if params.IsNFT.Bool {
			query = query.Where(sq.Eq{"max_supply": 1})
		} else {
			query = query.Where(sq.Gt{"max_supply": 1})
		}
	}

	if len(params.OrderBy) > 0 {
		direction := "ASC"
		if params.Descending {
			direction = "DESC"
		}

		query = query.OrderBy(params.OrderBy, direction)
	}

	rows, err := query.RunWith(w.DB).Query()
	if err != nil {
		return nil, err
	}

	var tokens []Token
	for rows.Next() {
		var token Token
		err := rows.Scan(
			&token.SCID,
			&token.Name,
			&token.MaxSupply,
			&token.TotalSupply,
			&token.Decimals,
			&token.StandardType,
			&token.Metadata,
			&token.IsFavorite,
			&token.ListOrderFavorite,
			&token.Image,
			&token.Symbol,
			&token.FolderId,
		)
		if err != nil {
			return tokens, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (w *Wallet) StoreToken(token Token) error {
	tx, err := w.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO tokens (sc_id,name,max_supply,total_supply,decimals,standard_type,metadata,is_favorite,list_order_favorite,image,symbol,folder_id)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(sc_id,folder_id) DO UPDATE SET
		name = excluded.name,
		max_supply = excluded.max_supply,
		total_supply = excluded.total_supply,
		decimals = excluded.decimals,
		standard_type = excluded.standard_type,
		metadata = excluded.metadata,
		is_favorite = excluded.is_favorite,
		list_order_favorite = excluded.list_order_favorite,
		image = excluded.image,
		symbol = excluded.symbol,
		folder_id = excluded.folder_id;
	`, token.SCID, token.Name, token.MaxSupply, token.TotalSupply, token.Decimals,
		token.StandardType, token.Metadata, token.IsFavorite,
		token.ListOrderFavorite, token.Image, token.Symbol, token.FolderId)
	if err != nil {
		return err
	}

	return handleDatabaseCommit(tx)
}

func (w *Wallet) DelTokenFolder(id int64) error {
	tx, err := w.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		PRAGMA recursive_triggers = ON;
		DELETE FROM token_folders
		WHERE id = ?;
		PRAGMA recursive_triggers = OFF;
	`, id)
	if err != nil {
		return err
	}

	return handleDatabaseCommit(tx)
}

func (w *Wallet) DelToken(scId string) error {
	tx, err := w.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM tokens
		WHERE sc_id = ?;
	`, scId)
	if err != nil {
		return err
	}

	return handleDatabaseCommit(tx)
}
