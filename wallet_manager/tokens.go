package wallet_manager

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/g45t345rt/g45w/sc"
)

type TokenFolder struct {
	ID       int32
	Name     string
	ParentId sql.NullInt32
}

type Token struct {
	SCID              string
	Name              string
	MaxSupply         sql.NullInt64 // 1 is an NFT
	TotalSupply       sql.NullInt64
	Decimals          int32 // native dero decimals is 5
	StandardType      sc.SCType
	Metadata          sql.NullString
	IsFavorite        sql.NullBool
	ListOrderFavorite sql.NullInt32
	Image             sql.NullString
	Symbol            sql.NullString
}

func initDatabaseTokens(db *sql.DB) error {
	dbTx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = dbTx.Exec(`
			CREATE TABLE IF NOT EXISTS token_folders (
				id INT PRIMARY KEY AUTOINCREMENT,
				name VARCHAR NOT NULL,
				parent_id INT,
				FOREIGN KEY (parent_id) REFERENCES token_folders(id) ON DELETE CASCADE
			);

			CREATE TABLE IF NOT EXISTS folder_tokens (
				folder_id INT,
				sc_id VARCHAR,
				PRIMARY KEY (folder_id,sc_id),
				FOREIGN KEY (folder_id) REFERENCES token_folders(id) ON DELETE CASCADE
			);

			CREATE TABLE IF NOT EXISTS tokens (
				sc_id VARCHAR PRIMARY KEY,
				name VARCHAR NOT NULL,
				max_supply BIGINT,
				total_supply BIGINT,
				decimals INT NOT NULL,
				standard_type VARCHAR NOT NULL,
				metadata VARCHAR,
				is_favorite BOOL,
				list_order_favorite INT,
				image VARCHAR,
				symbol VARCHAR
			);
		`)
	if err != nil {
		return err
	}

	return handleDatabaseCommit(dbTx)
}

func (w *Wallet) GetTokenFolder(id int32) (*TokenFolder, error) {
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

func (w *Wallet) GetTokenFolderPath(id sql.NullInt32) (string, error) {
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

func (w *Wallet) GetTokenFolderFolders(id sql.NullInt32) ([]TokenFolder, error) {
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

	folder.ID = int32(id)
	return nil
}

type GetTokensParams struct {
	Descending bool
	OrderBy    string
	IsFavorite sql.NullBool
	FolderId   sql.NullInt32
	IsNFT      sql.NullBool
}

func (w *Wallet) GetTokens(params GetTokensParams) ([]Token, error) {
	query := sq.Select("*").From("tokens")

	if params.IsFavorite.Valid {
		query = query.Where(sq.Eq{"is_favorite": params.IsFavorite.Bool})
	}

	if params.FolderId.Valid {
		query = query.RightJoin("folder_tokens ON tokens.sc_id = folder_tokens.sc_id").
			Where(sq.Eq{"folder_tokens.folder_id": params.FolderId.Int32})
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
		INSERT INTO tokens (sc_id,name,max_supply,total_supply,decimals,standard_type,metadata,is_favorite,list_order_favorite,image,symbol)
		VALUES (?,?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(sc_id) DO UPDATE SET
		name = excluded.name,
		max_supply = excluded.max_supply,
		total_supply = excluded.total_supply,
		decimals = excluded.decimals,
		standard_type = excluded.standard_type,
		metadata = excluded.metadata,
		is_favorite = excluded.is_favorite,
		list_order_favorite = excluded.list_order_favorite,
		image = excluded.image,
		symbol = excluded.symbol;
	`, token.SCID, token.Name, token.MaxSupply, token.TotalSupply, token.Decimals,
		token.StandardType, token.Metadata, token.IsFavorite,
		token.ListOrderFavorite, token.Image, token.Symbol)
	if err != nil {
		return err
	}

	return handleDatabaseCommit(tx)
}

func (w *Wallet) DelTokenFolder(id int32) error {
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
