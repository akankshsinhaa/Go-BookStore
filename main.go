package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"unicode"

	"github.com/gorilla/sessions"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template
var db *sql.DB

var store *sessions.CookieStore

func init() {
	tpl, _ = template.ParseGlob("templates/*.html")

	// Initialize the session store
	store = sessions.NewCookieStore([]byte("your-secret-key"))
}

func main() {
	file, _ := os.Create("file.log")
	log.SetOutput(file)
	var err error
	db, err = sql.Open("postgres", "user=app_user dbname=app_database password=app_password host=localhost sslmode=disable")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", loginAuthHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)
	http.HandleFunc("/allbooks", profileHandler)
	http.HandleFunc("/addreview", addReviewHandler)
	http.HandleFunc("/viewreviews", viewReviewsHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/admin/addbook", addBookHandler)

	http.ListenAndServe("localhost:8088", nil)
	file.Close()

}

// loginHandler serves the login form
func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("*****loginHandler running*****")
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "login.html", nil)
}

// loginAuthHandler authenticates user login and redirects to allbooks.html
func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("*****loginAuthHandler running*****")

	fmt.Println("*****loginAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Println("*****loginAuthHandler running*****")
	fmt.Println("username:", username, "password:", password)

	// Retrieve the hashed password from the database to compare with the user-supplied password
	var hash string
	stmt := "SELECT Hash FROM bcrypt WHERE Username = $1"
	row := db.QueryRow(stmt, username)
	err := row.Scan(&hash)
	log.Println("hash from db:", hash)
	fmt.Println("hash from db:", hash)
	if err != nil {
		log.Println("error selecting Hash in db by Username")
		fmt.Println("error selecting Hash in db by Username")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}
	// Compare the hashed password with the user's password input
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		// Create a session and store user information
		session, _ := store.Get(r, "session-name")
		session.Values["username"] = username
		session.Save(r, w)

		http.Redirect(w, r, "/allbooks", http.StatusSeeOther) // Redirect to allbooks.html
		return
	}
	log.Println("incorrect password")

	fmt.Println("incorrect password")
	tpl.ExecuteTemplate(w, "login.html", "check username and password")
}

// registerHandler serves the registration form
func registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("*****registerHandler running*****")

	fmt.Println("*****registerHandler running*****")
	tpl.ExecuteTemplate(w, "register.html", nil)
}

// registerAuthHandler creates a new user in the database
func registerAuthHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("*****registerAuthHandler running*****")

	fmt.Println("*****registerAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Println("username:", username, "\npassword:", password)

	fmt.Println("username:", username, "\npassword:", password)

	// Check username criteria
	var nameAlphaNumeric = true
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			nameAlphaNumeric = false
		}
	}
	var nameLength bool
	if 5 <= len(username) && len(username) <= 50 {
		nameLength = true
	}

	// Check password criteria
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			pswdLowercase = true
		case unicode.IsUpper(char):
			pswdUppercase = true
		case unicode.IsNumber(char):
			pswdNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	if 11 < len(password) && len(password) < 60 {
		pswdLength = true
	}

	// Check if username and password criteria are met
	if !nameAlphaNumeric || !nameLength || !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces {
		tpl.ExecuteTemplate(w, "register.html", "please check username and password criteria")
		return
	}

	// Check if the username already exists in the database
	stmt := "SELECT UserID FROM bcrypt WHERE Username = $1"
	row := db.QueryRow(stmt, username)
	var uID string
	err := row.Scan(&uID)
	if err != sql.ErrNoRows {
		log.Println("username already exists, err:", err)
		fmt.Println("username already exists, err:", err)
		tpl.ExecuteTemplate(w, "register.html", "username already taken")
		return
	}

	// Create a hash from the password
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypt err:", err)

		fmt.Println("bcrypt err:", err)
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering an account")
		return
	}

	// Prepare and execute the SQL INSERT statement
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO bcrypt (Username, Hash) VALUES ($1, $2);")
	if err != nil {
		log.Println("error preparing statement:", err)
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering an account")
		return
	}
	defer insertStmt.Close()

	var result sql.Result
	result, err = insertStmt.Exec(username, hash)
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	log.Println("rowsAff:", rowsAff)
	log.Println("lastIns:", lastIns)
	log.Println("err:", err)
	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)
	if err != nil {
		log.Println("error inserting new user")
		fmt.Println("error inserting new user")
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering an account")
		return
	}

	// Create a session and store user information
	session, _ := store.Get(r, "session-name")
	session.Values["username"] = username
	session.Save(r, w)

	http.Redirect(w, r, "/allbooks", http.StatusSeeOther)
}

// profileHandler displays the user's profile
// profileHandler displays the allbooks.html page with book data
// profileHandler displays the allbooks.html page with book data
func profileHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the username from the session
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if false {
		fmt.Print(username)
	}

	// Query the "books" table to retrieve book data
	rows, err := db.Query("SELECT id, title, author, publisher FROM books")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []struct {
		ID        int
		Title     string
		Author    string
		Publisher string
	}

	for rows.Next() {
		var book struct {
			ID        int
			Title     string
			Author    string
			Publisher string
		}
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Publisher); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	// Render the allbooks.html template with the book data
	tpl.ExecuteTemplate(w, "allbooks.html", books)
}

// addReviewHandler handles adding reviews to a book
// addReviewHandler handles adding reviews to a book
// addReviewHandler displays the "addreview.html" page and handles review submission
func addReviewHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the book ID from the query string
	bookID := r.URL.Query().Get("id")

	// Retrieve the username from the session
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if false {
		fmt.Printf(username)
	}

	if r.Method == http.MethodPost {
		// Parse the review text from the form
		reviewText := r.FormValue("review")

		// Retrieve the existing reviews for the book
		// Retrieve the existing reviews for the book
		var existingReviews pq.StringArray
		err := db.QueryRow("SELECT COALESCE(reviews, '{}'::TEXT[]) FROM books WHERE id = $1", bookID).Scan(&existingReviews)
		if err != nil {
			log.Println("Error retrieving existing reviews:", err)

			fmt.Println("Error retrieving existing reviews:", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Append the new review to the existing reviews
		existingReviews = append(existingReviews, reviewText)

		// Update the reviews for the book in the database
		_, err = db.Exec("UPDATE books SET reviews = $1 WHERE id = $2", existingReviews, bookID)
		if err != nil {
			log.Println("Error updating reviews:", err)
			fmt.Println("Error updating reviews:", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Redirect back to the book list
		http.Redirect(w, r, "/allbooks", http.StatusSeeOther)
		return
	}

	// Render the "addreview.html" template when the page is initially loaded
	tpl.ExecuteTemplate(w, "addreview.html", nil)
}

// viewReviewsHandler displays the reviews for a book
func viewReviewsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the book ID from the query string
	bookID := r.URL.Query().Get("id")

	// Retrieve the reviews for the book
	var reviews []string
	err := db.QueryRow("SELECT reviews FROM books WHERE id = $1", bookID).Scan(&reviews)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Render a template to display the reviews
	tpl.ExecuteTemplate(w, "reviews.html", reviews)
}
func adminHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "admin.html", nil)
}
func addBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	title := r.FormValue("title")
	author := r.FormValue("author")
	publisher := r.FormValue("publisher")

	// Insert the book into the database
	_, err := db.Exec("INSERT INTO books (title, author, publisher) VALUES ($1, $2, $3)", title, author, publisher)
	if err != nil {
		log.Println("Error adding book:", err)
		fmt.Println("Error adding book:", err)
		http.Error(w, "Error adding book", http.StatusInternalServerError)
		return
	}

	// Redirect back to the admin page
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
