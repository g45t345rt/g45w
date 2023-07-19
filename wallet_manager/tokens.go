package wallet_manager

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/g45t345rt/g45w/sc"
)

type TokenFolder struct {
	ID       int
	Name     string
	ParentId int
}

type Token struct {
	SCID              string
	Name              sql.NullString
	MaxSupply         sql.NullInt64 // 1 is an NFT
	Decimals          sql.NullInt32 // native dero decimals is 5
	StandardType      sc.SCType
	Metadata          sql.NullString
	IsFavorite        sql.NullBool
	ListOrderFavorite sql.NullInt32
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
				parent_id INTEGER,
				FOREIGN KEY (parent_id) REFERENCES token_folders(id) ON DELETE CASCADE
			);

			CREATE TABLE IF NOT EXISTS folder_tokens (
				folder_id INTEGER,
				sc_id VARCHAR,
				PRIMARY KEY (folder_id,sc_id),
				FOREIGN KEY (folder_id) REFERENCES token_folders(id) ON DELETE CASCADE
			);

			CREATE TABLE IF NOT EXISTS tokens (
				sc_id VARCHAR PRIMARY KEY,
				name VARCHAR,
				max_supply BIGINT,
				decimals INTEGER,
				standard_type VARCHAR,
				metadata VARCHAR,
				is_favorite BOOL,
				list_order_favorite INTEGER
			);
		`)
	if err != nil {
		return err
	}

	return handleDatabaseCommit(dbTx)
}

func (w *Wallet) GetTokenFolders(id *int) ([]TokenFolder, error) {
	query := sq.Select("*").From("token_folders").Where(sq.Eq{"parent_id": id})

	rows, err := w.DB.Query(query.ToSql())
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

type GetTokensParams struct {
	Descending bool
	OrderBy    string
	IsFavorite *bool
	FolderId   *int
	IsNFT      *bool
}

func (w *Wallet) GetTokens(params GetTokensParams) ([]Token, error) {
	query := sq.Select("*").From("tokens")

	if params.IsFavorite != nil {
		query = query.Where(sq.Eq{"is_favorite": params.IsFavorite})
	}

	if params.FolderId != nil {
		query = query.RightJoin("folder_tokens ON id = folder_id")
	}

	if params.IsNFT != nil {
		if *params.IsNFT {
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

	rows, err := w.DB.Query(query.ToSql())
	if err != nil {
		return nil, err
	}

	return rowsScanTokens(rows)
}

func rowsScanTokens(rows *sql.Rows) ([]Token, error) {
	var tokens []Token
	for rows.Next() {
		var token Token
		err := rows.Scan(
			&token.SCID,
			&token.Name,
			&token.MaxSupply,
			&token.Decimals,
			&token.StandardType,
			&token.Metadata,
			&token.IsFavorite,
			&token.ListOrderFavorite,
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
		INSERT INTO tokens (sc_id,name,max_supply,decimals,standard_type,metadata,is_favorite,list_order_favorite)
		VALUES (?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(sc_id) DO UPDATE SET
		name = excluded.name,
		max_supply = excluded.max_supply,
		decimals = excluded.decimals,
		standard_type = excluded.standard_type,
		metadata = excluded.metadata,
		is_favorite = excluded.is_favorite,
		list_order_favorite = excluded.list_order_favorite;
	`, token.SCID, token.Name, token.MaxSupply, token.Decimals, token.StandardType,
		token.Metadata, token.IsFavorite,
		token.ListOrderFavorite)
	if err != nil {
		return err
	}

	return handleDatabaseCommit(tx)
}

func (w *Wallet) DelTokenFolder(id int) error {
	tx, err := w.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM token_folders
		WHERE id = ?;
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
