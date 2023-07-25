package wallet_manager

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/g45t345rt/g45w/sc"
)

type TokenFolder struct {
	ID       int64
	Name     string
	ParentId sql.NullInt64
}

type Token struct {
	ID                int64
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

func DeroToken() Token {
	scId := crypto.ZEROHASH.String()
	return Token{
		ID:        -1,
		SCID:      scId,
		Decimals:  5,
		Name:      "Dero",
		Symbol:    sql.NullString{String: "DERO", Valid: true},
		MaxSupply: sql.NullInt64{Int64: 2100000000000, Valid: true}, // max supply is 21,000,000 but don't forget 5 decimals
	}
}

func initDatabaseTokens(db *sql.DB) error {
	_, err := db.Exec(`
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
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				sc_id VARCHAR NOT NULL,
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
				folder_id INTEGER
			);
		`)
	return err
}

func (w *Wallet) GetTokenFolder(id int64) (*TokenFolder, error) {
	query := sq.Select("*").From("token_folders").Where(sq.Eq{"id": id})

	row := query.RunWith(w.DB).QueryRow()

	var folder *TokenFolder
	err := row.Scan(
		&folder.ID,
		&folder.Name,
		&folder.ParentId,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return folder, nil
}

func (w *Wallet) GetTokenFolderPath(id sql.NullInt64) (string, error) {
	if !id.Valid {
		return "root", nil
	}

	query := sq.Select("*").From("token_folders").Where(sq.Eq{"id": id})

	row := query.RunWith(w.DB).QueryRow()

	var folder TokenFolder
	err := row.Scan(
		&folder.ID,
		&folder.Name,
		&folder.ParentId,
	)
	if err != nil {
		return "", err
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

func (w *Wallet) UpdateFolderToken(folder TokenFolder) error {
	exists, err := w.FolderTokenExists(folder)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("folder already exists")
	}

	_, err = w.DB.Exec(`
		UPDATE token_folders
		SET name = ?,
				parent_id = ?
		WHERE id = ?;
	`, folder.Name, folder.ParentId, folder.ID)
	return err
}

func (w *Wallet) FolderTokenExists(folder TokenFolder) (bool, error) {
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
		return false, err
	}

	return count >= 1, nil
}

func (w *Wallet) InsertFolderToken(folder TokenFolder) error {
	// can't use UNIQUE() constraint because null does not count as towards uniqueness
	// https://stackoverflow.com/questions/22699409/sqlite-null-and-unique
	// we check manually instead
	exists, err := w.FolderTokenExists(folder)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("folder already exists")
	}

	result, err := w.DB.Exec(`
		INSERT INTO token_folders (name,parent_id)
		VALUES (?,?);
	`, folder.Name, folder.ParentId)
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

func (w *Wallet) GetToken(id int64) (*Token, error) {
	query := sq.Select("*").From("tokens").Where(sq.Eq{"id": id})
	row := query.RunWith(w.DB).QueryRow()

	var token Token
	err := row.Scan(
		&token.ID,
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
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &token, nil
}

func (w *Wallet) GetTokenCount(folderId sql.NullInt64) (int, error) {
	query := sq.Select("COUNT(*)").From("tokens")

	if folderId.Valid {
		query = query.Where(sq.Eq{"folder_id": folderId.Int64})
	} else {
		query = query.Where(sq.Eq{"folder_id": nil})
	}

	row := query.RunWith(w.DB).QueryRow()

	var count int
	err := row.Scan(&count)
	return count, err
}

type GetTokensParams struct {
	Descending bool
	OrderBy    string
	IsFavorite sql.NullBool
	FolderId   *sql.NullInt64
	IsNFT      sql.NullBool
}

func (w *Wallet) GetTokens(params GetTokensParams) ([]Token, error) {
	query := sq.Select("*").From("tokens")

	if params.IsFavorite.Valid {
		query = query.Where(sq.Eq{"is_favorite": params.IsFavorite.Bool})
	}

	if params.FolderId != nil {
		if params.FolderId.Valid {
			query = query.Where(sq.Eq{"folder_id": params.FolderId})
		} else {
			query = query.Where(sq.Eq{"folder_id": nil})
		}
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
			&token.ID,
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

func (w *Wallet) InsertToken(token Token) error {
	_, err := w.DB.Exec(`
		INSERT INTO tokens (sc_id,name,max_supply,total_supply,decimals,standard_type,metadata,is_favorite,list_order_favorite,image,symbol,folder_id)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?);
	`, token.SCID, token.Name, token.MaxSupply, token.TotalSupply, token.Decimals,
		token.StandardType, token.Metadata, token.IsFavorite,
		token.ListOrderFavorite, token.Image, token.Symbol, token.FolderId)
	return err
}

func (w *Wallet) UpdateToken(token Token) error {
	_, err := w.DB.Exec(`
		UPDATE tokens
		SET sc_id = ?,
				name = ?,
				max_supply = ?,
				total_supply = ?,
				decimals = ?,
				standard_type = ?,
				metadata = ?,
				is_Favorite = ?,
				list_order_favorite = ?,
				image = ?,
				symbol = ?,
				folder_id = ?
		WHERE id = ?;
	`, token.SCID, token.Name, token.MaxSupply, token.TotalSupply, token.Decimals,
		token.StandardType, token.Metadata, token.IsFavorite, token.ListOrderFavorite,
		token.Image, token.Symbol, token.FolderId, token.ID)
	return err
}

func (w *Wallet) DelTokenFolder(id int64) error {
	_, err := w.DB.Exec(`
		PRAGMA recursive_triggers = ON;
		DELETE FROM token_folders
		WHERE id = ?;
		PRAGMA recursive_triggers = OFF;
	`, id)
	return err
}

func (w *Wallet) DelToken(id int64) error {
	_, err := w.DB.Exec(`
		DELETE FROM tokens
		WHERE id = ?;
	`, id)
	return err
}
