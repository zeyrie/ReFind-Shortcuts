package repo

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zeyrie/ReFind-Shortcuts/internal/domain"
)

type RepoManager struct {
	db     *sql.DB
	dbPath string
}

func (repo *RepoManager) initRepo() {

	dbPath := "./internal/repo/dev_test.db"

	// Open
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open DB: ", err)
	}

	// Check Connection
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}

	repo.db = db
	repo.dbPath = dbPath

	log.Println("Initialized DB")
}

func (repo *RepoManager) InitializeTable() {

	// Create the DB
	repo.initRepo()

	mkCategoryTableQuery := `
	CREATE TABLE IF NOT EXISTS Category(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT
	);`

	mkShortcutTableQuery := `
	CREATE TABLE IF NOT EXISTS Shortcut(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		value TEXT NOT NULL,
		description TEXT NOT NULL,
		note TEXT,
		iconURL TEXT,
		categoryId INTEGER NOT NULL,
		FOREIGN KEY (categoryId) REFERENCES Category(id)
	);`

	if _, err := repo.db.Exec(mkCategoryTableQuery); err != nil {
		log.Fatal("Category Table creation failed", err)
	}

	if _, err := repo.db.Exec(mkShortcutTableQuery); err != nil {
		log.Fatal("Shortcuts Table creation failed", err)
	}

	log.Println("Created tables for Shortcuts")

}

func (repo *RepoManager) InsertCategory(cg domain.Category) (*domain.Category, error) {
	insertQuery := `INSERT INTO Category (title, description) VALUES (?, ?);`

	res, err := repo.db.Exec(insertQuery, cg.Title, getNullString(cg.Description))

	if err != nil {
		log.Printf("Failed to Insert %s: %v\n", cg.Title, err)
		return nil, err
	}

	cg.ID, err = res.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last Insert Id for %s: %v\n", cg.Title, err)
		return nil, err
	}

	return &cg, nil
}

func (repo *RepoManager) InsertShortcut(sh domain.Shortcut) (*domain.Shortcut, error) {
	insertQuery := `INSERT INTO Shortcut (value, description, note, iconURL, categoryId) VALUES ( ?, ?, ?, ?, ? );`

	res, err := repo.db.Exec(
		insertQuery,
		sh.Value,
		sh.Description,
		getNullString(sh.Note),
		getNullString(sh.IconUrl),
		sh.Category,
	)

	if err != nil {
		log.Printf("Failed to Insert %s: %v\n", sh.Description, err)
		return nil, err
	}

	sh.ID, err = res.LastInsertId()
	if err != nil {
		log.Printf("Failed to get last Insert Id for %s: %v", sh.Description, err)
		return nil, err
	}

	return &sh, nil
}

func (repo *RepoManager) FetchAllCategories() ([]domain.Category, error) {
	var categories []domain.Category

	query := `SELECT * FROM Category;`

	rows, err := repo.db.Query(query)
	if err != nil {
		log.Printf("Failed to fetch categories: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var category domain.Category
		var desc sql.NullString

		if err := rows.Scan(&category.ID, &category.Title, &desc); err != nil {
			log.Println("Failed to scan category: ", err)
			return nil, err
		}

		if desc.Valid {
			category.Description = &desc.String
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error encountered while iteration over categories: ", err)
		return nil, err
	}

	return categories, nil
}

func (repo *RepoManager) FetchAllShortcuts() ([]domain.Shortcut, error) {
	var shortcuts []domain.Shortcut

	query := `SELECT * FROM Shortcut;`

	rows, err := repo.db.Query(query)
	if err != nil {
		log.Println("Failed to fetch shortcuts: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var shortcut domain.Shortcut
		var note sql.NullString
		var iconUrl sql.NullString

		if err := rows.Scan(&shortcut.ID, &shortcut.Value, &shortcut.Description, &note, &iconUrl, &shortcut.Category); err != nil {
			log.Println("Failed to scan shortcut: ", err)
			return nil, err
		}

		if note.Valid {
			shortcut.Note = &note.String
		}

		if iconUrl.Valid {
			shortcut.IconUrl = &iconUrl.String
		}

		shortcuts = append(shortcuts, shortcut)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error encountred while iteration over shortcuts: ", err)
		return nil, err
	}

	return shortcuts, nil
}

func (repo *RepoManager) CloseDB() {
	err := repo.db.Close()
	if err != nil {
		log.Println("Failed to cloes DB", err)
	}
}

func getNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{Valid: false}
}
