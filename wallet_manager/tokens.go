package wallet_manager

import (
	"bytes"
	"database/sql"
	"fmt"
	"image"

	//"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gioui.org/op/paint"
	sq "github.com/Masterminds/squirrel"

	"github.com/deroproject/derohe/cryptography/crypto"
	"github.com/deroproject/derohe/rpc"
	"github.com/g45t345rt/g45w/app_db/order_column"
	"github.com/g45t345rt/g45w/app_db/schema_version"
	"github.com/g45t345rt/g45w/assets"
	"github.com/g45t345rt/g45w/caching"
	"github.com/g45t345rt/g45w/multi_fetch"
	"github.com/g45t345rt/g45w/sc"
	"github.com/g45t345rt/g45w/sc/dex_sc"
	"github.com/g45t345rt/g45w/sc/g45_sc"
	"github.com/g45t345rt/g45w/sc/unknown_sc"
	"github.com/g45t345rt/g45w/settings"
	"github.com/g45t345rt/g45w/theme"
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
	IsFavorite        bool
	ListOrderFavorite int
	ImageUrl          sql.NullString
	Symbol            sql.NullString
	FolderId          sql.NullInt64
	CreatedTimestamp  sql.NullInt64 // date created on the blockchain
	AddedTimestamp    sql.NullInt64 // date added to the sql table

	imgLoaded bool
	imageOp   *paint.ImageOp
	hash      *crypto.Hash
}

var tokenFavOrderer = order_column.Orderer{
	TableName:   "tokens",
	ColumnName:  "list_order_favorite",
	FilterQuery: "is_favorite = true",
}

func (token *Token) GetHash() crypto.Hash {
	if token.hash != nil {
		return *token.hash
	}

	hash := crypto.HashHexToHash(token.SCID)
	token.hash = &hash
	return hash
}

func (token *Token) DataDirPath() (string, error) {
	cacheDir := settings.CacheDir
	tokenDataDirPath := filepath.Join(cacheDir, "tokens", token.SCID)
	err := os.MkdirAll(tokenDataDirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return tokenDataDirPath, nil
}

func (token *Token) RefreshImageOp() {
	token.imgLoaded = false
}

func (token *Token) LoadImageOp() paint.ImageOp {
	if !token.imgLoaded {
		token.imgLoaded = true
		go func() {
			imgOp, err := token.GetImageOp()
			if err == nil {
				token.imageOp = &imgOp
			}
		}()
	}

	if token.imageOp != nil {
		return *token.imageOp
	}

	return theme.Current.TokenImage
}

var imageMemCache map[string]paint.ImageOp
var imageMemCacheMutex sync.Mutex

func (token *Token) GetImageOp() (imgOp paint.ImageOp, err error) {
	if imageMemCache == nil {
		imageMemCacheMutex.Lock()
		imageMemCache = make(map[string]paint.ImageOp)

		// load default native token image
		img, _ := assets.GetImage("dero.jpg")
		imageMemCache[crypto.ZEROHASH.String()] = paint.NewImageOp(img)
		imageMemCacheMutex.Unlock()
	}

	imageMemCacheMutex.Lock()
	imgOp, ok := imageMemCache[token.SCID]
	imageMemCacheMutex.Unlock()

	if ok {
		return
	}

	if token.ImageUrl.Valid {
		relCachePath := filepath.Join("tokens", token.SCID)
		cacheFileName := "image"

		var imgData []byte
		var exists bool
		exists, err = caching.Get(relCachePath, cacheFileName, &imgData)
		if err != nil {
			return
		}

		if !exists {
			// download from ipfs/http
			var res *http.Response
			res, err = multi_fetch.Fetch(token.ImageUrl.String)
			if err != nil {
				return
			}
			defer res.Body.Close()

			imgData, err = io.ReadAll(res.Body)
			if err != nil {
				return
			}

			err = caching.Store(relCachePath, cacheFileName, imgData)
			if err != nil {
				return
			}
		}

		var img image.Image // jpg, png, gif and webp by importing golang.org/x/image/webp
		buffer := bytes.NewBuffer(imgData)
		img, _, err = image.Decode(buffer)
		if err != nil {
			return
		}

		imgOp = paint.NewImageOp(img)
		imageMemCacheMutex.Lock()
		imageMemCache[token.SCID] = imgOp
		imageMemCacheMutex.Unlock()

		return
	}

	err = fmt.Errorf("no image")
	return
}

func (token *Token) Parse(scId string, scResult rpc.GetSC_Result) error {
	scType := sc.CheckType(scResult.Code)
	token.SCID = scId
	token.StandardType = scType
	token.AddedTimestamp = sql.NullInt64{Int64: time.Now().Unix(), Valid: true}

	switch scType {
	case sc.G45_FAT_TYPE:
		fat := g45_sc.G45_FAT{}
		err := fat.Parse(scId, scResult.VariableStringKeys)
		if err != nil {
			return err
		}

		metadata := g45_sc.TokenMetadata{}
		err = metadata.Parse(fat.Metadata)
		if err != nil {
			return err
		}

		token.Name = metadata.Name
		token.Decimals = int64(fat.Decimals)
		token.MaxSupply = sql.NullInt64{Int64: int64(fat.MaxSupply), Valid: true}
		token.ImageUrl = sql.NullString{String: metadata.Image, Valid: true}
		token.Symbol = sql.NullString{String: metadata.Symbol, Valid: true}
		token.Metadata = sql.NullString{String: fat.Metadata, Valid: true}
		token.CreatedTimestamp = sql.NullInt64{Int64: int64(fat.Timestamp), Valid: true}
	case sc.G45_AT_TYPE:
		at := g45_sc.G45_AT{}
		err := at.Parse(scId, scResult.VariableStringKeys)
		if err != nil {
			return err
		}

		metadata := g45_sc.TokenMetadata{}
		err = metadata.Parse(at.Metadata)
		if err != nil {
			return err
		}

		token.Name = metadata.Name
		token.Decimals = int64(at.Decimals)
		token.MaxSupply = sql.NullInt64{Int64: int64(at.MaxSupply), Valid: true}
		token.ImageUrl = sql.NullString{String: metadata.Image, Valid: true}
		token.Symbol = sql.NullString{String: metadata.Symbol, Valid: true}
		token.Metadata = sql.NullString{String: at.Metadata, Valid: true}
		token.CreatedTimestamp = sql.NullInt64{Int64: int64(at.Timestamp), Valid: true}
	case sc.G45_NFT_TYPE:
		nft := g45_sc.G45_NFT{}
		err := nft.Parse(scId, scResult.VariableStringKeys)
		if err != nil {
			return err
		}

		metadata := g45_sc.NFTMetadata{}
		err = metadata.Parse(nft.Metadata)
		if err != nil {
			return err
		}

		token.Name = metadata.Name
		token.Decimals = 0
		token.MaxSupply = sql.NullInt64{Int64: 1, Valid: true}
		token.ImageUrl = sql.NullString{String: metadata.Image, Valid: true}
		token.Metadata = sql.NullString{String: nft.Metadata, Valid: true}
		token.CreatedTimestamp = sql.NullInt64{Int64: int64(nft.Timestamp), Valid: true}
	case sc.DEX_SC_TYPE:
		dex := dex_sc.Token{}
		err := dex.Parse(scId, scResult.VariableStringKeys)
		if err != nil {
			return err
		}

		token.Name = dex.Name
		token.Decimals = int64(dex.Decimals)
		token.ImageUrl = sql.NullString{String: dex.ImageUrl, Valid: true}
		token.Symbol = sql.NullString{String: dex.Symbol, Valid: true}
	case sc.UNKNOWN_TYPE:
		unknown := unknown_sc.Token{}
		unknown.Parse(scId, scResult.VariableStringKeys)

		token.Name = unknown.Name
		token.Decimals = int64(unknown.Decimals)
		token.ImageUrl = sql.NullString{String: unknown.ImageUrl, Valid: true}
		token.Symbol = sql.NullString{String: unknown.Symbol, Valid: true}
	}

	return nil
}

func GetSC(scId string) (result rpc.GetSC_Result, cached bool, err error) {
	cacheFileName := "get_sc"
	relCachePath := filepath.Join("tokens", scId)
	cached, err = caching.Get(relCachePath, cacheFileName, &result)
	if err != nil {
		return
	}

	if cached {
		return
	}

	err = RPCCall("DERO.GetSC", rpc.GetSC_Params{
		SCID:      scId,
		Variables: true,
		Code:      true,
	}, &result)
	if err != nil {
		return
	}

	err = caching.Store(relCachePath, cacheFileName, result)
	return
}

func GetTokenBySCID(scId string) (token *Token, err error) {
	if scId == crypto.ZEROHASH.String() {
		token = DeroToken()
	} else {
		var result rpc.GetSC_Result
		result, _, err = GetSC(scId)
		if err != nil {
			return
		}

		token = &Token{}
		err = token.Parse(scId, result)
		if err != nil {
			return
		}
	}

	return
}

func DeroToken() *Token {
	scId := crypto.ZEROHASH.String()

	return &Token{
		ID:        -1,
		SCID:      scId,
		Decimals:  5,
		Name:      "Dero",
		Symbol:    sql.NullString{String: "DERO", Valid: true},
		MaxSupply: sql.NullInt64{Int64: 2100000000000, Valid: true}, // max supply is 21,000,000 but don't forget 5 decimals
	}
}

func initTableTokens(db *sql.DB) error {
	version, err := schema_version.GetVersion(db, "tokens")
	if err != nil {
		return err
	}

	if version == 0 {
		_, err = db.Exec(`
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
		if err != nil {
			return err
		}

		version = 1
		err = schema_version.StoreVersion(db, "tokens", version)
		if err != nil {
			return err
		}
	}

	if version == 1 {
		_, err = db.Exec(`
			ALTER TABLE tokens ADD COLUMN created_timestamp BIGTINT;
			ALTER TABLE tokens ADD COLUMN added_timestamp BIGINT;
		`)
		if err != nil {
			return err
		}

		version = 2
		err = schema_version.StoreVersion(db, "tokens", version)
		if err != nil {
			return err
		}
	}

	if version == 2 {
		// There is no ALTER COLUMN in SQLITE -___-
		// Alter is_favorite and list_order_favorite to not null; This was a mistake way back.
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		_, err = tx.Exec(`
			UPDATE tokens SET is_favorite = 0 WHERE is_favorite IS NULL;
			ALTER TABLE tokens RENAME COLUMN is_favorite TO temp_is_favorite;
			ALTER TABLE tokens ADD COLUMN is_favorite BOOL NOT NULL DEFAULT 0; 
			UPDATE tokens SET is_favorite = COALESCE(temp_is_favorite, 0);
			ALTER TABLE tokens DROP COLUMN temp_is_favorite;

			UPDATE tokens SET list_order_favorite = 0 WHERE list_order_favorite IS NULL;
			ALTER TABLE tokens RENAME COLUMN list_order_favorite TO temp_list_order_favorite;
			ALTER TABLE tokens ADD COLUMN list_order_favorite INT NOT NULL DEFAULT 0;
			UPDATE tokens SET list_order_favorite = COALESCE(temp_list_order_favorite, 0);
			ALTER TABLE tokens DROP COLUMN temp_list_order_favorite;
		`)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}

		version = 3
		err = schema_version.StoreVersion(db, "tokens", version)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Wallet) GetTokenFolder(id int64) (*TokenFolder, error) {
	query := sq.Select("*").From("token_folders").Where(sq.Eq{"id": id})

	row := query.RunWith(w.DB).QueryRow()

	var folder TokenFolder
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

	return &folder, nil
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
	defer rows.Close()

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

func (w *Wallet) InsertFolderToken(folder TokenFolder) (int64, error) {
	// can't use UNIQUE() constraint because null does not count as towards uniqueness
	// https://stackoverflow.com/questions/22699409/sqlite-null-and-unique
	// we check manually instead
	exists, err := w.FolderTokenExists(folder)
	if err != nil {
		return -1, err
	}

	if exists {
		return -1, fmt.Errorf("folder already exists")
	}

	result, err := w.DB.Exec(`
		INSERT INTO token_folders (name,parent_id)
		VALUES (?,?);
	`, folder.Name, folder.ParentId)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
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
		&token.ImageUrl,
		&token.Symbol,
		&token.FolderId,
		&token.CreatedTimestamp,
		&token.AddedTimestamp,
		&token.IsFavorite,
		&token.ListOrderFavorite,
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
	query := sq.Select("*").From("tokens").OrderBy("list_order_favorite ASC")

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

		query = query.OrderBy(fmt.Sprintf("%s %s", params.OrderBy, direction))
	}

	rows, err := query.RunWith(w.DB).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			&token.ImageUrl,
			&token.Symbol,
			&token.FolderId,
			&token.CreatedTimestamp,
			&token.AddedTimestamp,
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

func (w *Wallet) InsertToken(token Token) error {
	tx, err := w.DB.Begin()
	if err != nil {
		return err
	}

	row := tx.QueryRow(`
		SELECT COUNT(*) FROM tokens
		WHERE sc_id = ? AND folder_id = ?
	`, token.SCID, token.FolderId)
	var count int
	err = row.Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	/*
		// Update fav list order if not specified
		if token.IsFavorite.Bool && !token.ListOrderFavorite.Valid {
			favLastOrder, err := w.GetFavTokenLastOrder()
			if err != nil {
				return err
			}
			token.ListOrderFavorite = sql.NullInt64{Int64: favLastOrder + 1, Valid: true}
		}*/

	if token.IsFavorite {
		token.ListOrderFavorite, err = tokenFavOrderer.GetNewOrderNumber(tx)
		if err != nil {
			return err
		}

		/*
			err = tokenFavOrderer.Insert(tx, orderNumber)
			if err != nil {
				tx.Rollback()
				return err
			}
		*/
	}

	_, err = tx.Exec(`
		INSERT INTO tokens (sc_id,name,max_supply,total_supply,decimals,standard_type,metadata,is_favorite,list_order_favorite,image,symbol,folder_id,created_timestamp,added_timestamp)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?);
	`, token.SCID, token.Name, token.MaxSupply, token.TotalSupply, token.Decimals,
		token.StandardType, token.Metadata, token.IsFavorite,
		token.ListOrderFavorite, token.ImageUrl, token.Symbol, token.FolderId, token.CreatedTimestamp, token.AddedTimestamp)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	// update token balance right away
	w.Memory.TokenAdd(token.GetHash())
	return nil
}

func (w *Wallet) GetFavTokenLastOrder() (int64, error) {
	row := w.DB.QueryRow(`
		SELECT list_order_favorite FROM tokens
		WHERE is_favorite = true
		ORDER BY list_order_favorite DESC
		LIMIT 1
	`)

	var lastOrder int64
	err := row.Scan(&lastOrder)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}

		return 0, err
	}

	return lastOrder, nil
}

func (w *Wallet) UpdateToken(token Token) error {
	// Update fav list order if not specified
	/*if token.IsFavorite.Bool && !token.ListOrderFavorite.Valid {
		favLastOrder, err := w.GetFavTokenLastOrder()
		if err != nil {
			return err
		}
		token.ListOrderFavorite = sql.NullInt64{Int64: favLastOrder + 1, Valid: true}
	}*/

	tx, err := w.DB.Begin()
	if err != nil {
		return err
	}

	currentToken, err := w.GetToken(token.ID)
	if err != nil {
		return err
	}

	if token.IsFavorite {
		err = tokenFavOrderer.Update(tx, currentToken.ListOrderFavorite, token.ListOrderFavorite)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		err = tokenFavOrderer.Delete(tx, currentToken.ListOrderFavorite)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec(`
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
				folder_id = ?,
				created_timestamp = ?,
				added_timestamp = ?
		WHERE id = ?;
	`, token.SCID, token.Name, token.MaxSupply, token.TotalSupply, token.Decimals,
		token.StandardType, token.Metadata, token.IsFavorite, token.ListOrderFavorite,
		token.ImageUrl, token.Symbol, token.FolderId, token.CreatedTimestamp, token.AddedTimestamp, token.ID)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
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
	tx, err := w.DB.Begin()
	if err != nil {
		return err
	}

	currentToken, err := w.GetToken(id)
	if err != nil {
		return err
	}

	err = tokenFavOrderer.Delete(tx, currentToken.ListOrderFavorite)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM tokens
		WHERE id = ?;
	`, id)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
