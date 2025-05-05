package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/lib/pq"
	
	"github.com/rni0719/test_workflow/internal/controller"
	"github.com/rni0719/test_workflow/internal/repository"
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
}

func NewApp(dbURL string) *App {
	app := &App{}
	app.Initialize(dbURL)
	app.initializeDB()
	return app
}

func (app *App) Initialize(dbURL string) {
	var err error
	app.DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// 接続テスト
	err = app.DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

func (app *App) Run(addr string) {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(app.Router)

	log.Printf("Server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

func (app *App) initializeRoutes() {
	// リポジトリの初期化
	workflowRepo := repository.NewWorkflowRepository(app.DB)
	taskRepo := repository.NewTaskRepository(app.DB)
	
	// コントローラーの初期化
	workflowController := controller.NewWorkflowController(workflowRepo)
	taskController := controller.NewTaskController(taskRepo)
	
	// ルートの設定
	app.Router.HandleFunc("/api/health", controller.HealthCheck).Methods("GET")
	
	// ワークフロー関連
	app.Router.HandleFunc("/api/workflows", workflowController.GetWorkflows).Methods("GET")
	app.Router.HandleFunc("/api/workflows", workflowController.CreateWorkflow).Methods("POST")
	app.Router.HandleFunc("/api/workflows/{id:[0-9]+}", workflowController.GetWorkflow).Methods("GET")
	app.Router.HandleFunc("/api/workflows/{id:[0-9]+}", workflowController.UpdateWorkflow).Methods("PUT")
	app.Router.HandleFunc("/api/workflows/{id:[0-9]+}", workflowController.DeleteWorkflow).Methods("DELETE")
	app.Router.HandleFunc("/api/workflows/{id:[0-9]+}/tasks", workflowController.GetWorkflowTasks).Methods("GET")
	
	// タスク関連
	app.Router.HandleFunc("/api/tasks", taskController.CreateTask).Methods("POST")
	app.Router.HandleFunc("/api/tasks/{id:[0-9]+}", taskController.GetTask).Methods("GET")
	app.Router.HandleFunc("/api/tasks/{id:[0-9]+}", taskController.UpdateTask).Methods("PUT")
	app.Router.HandleFunc("/api/tasks/{id:[0-9]+}", taskController.DeleteTask).Methods("DELETE")
}

func (app *App) initializeDB() {
	// ワークフローテーブル作成
	_, err := app.DB.Exec(`
		CREATE TABLE IF NOT EXISTS workflows (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// タスクテーブル作成
	_, err = app.DB.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			workflow_id INTEGER NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}
