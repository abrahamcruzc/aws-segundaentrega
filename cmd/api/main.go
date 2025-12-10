package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abrahamcruzc/aws-segundaentrega/internal/adapter/aws"
	apphttp "github.com/abrahamcruzc/aws-segundaentrega/internal/adapter/http"
	"github.com/abrahamcruzc/aws-segundaentrega/internal/adapter/http/handler"
	"github.com/abrahamcruzc/aws-segundaentrega/internal/adapter/storage/dynamodb"
	"github.com/abrahamcruzc/aws-segundaentrega/internal/adapter/storage/postgres"
	"github.com/abrahamcruzc/aws-segundaentrega/internal/adapter/storage/s3"
	"github.com/abrahamcruzc/aws-segundaentrega/internal/config"
	"github.com/abrahamcruzc/aws-segundaentrega/internal/usecase"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// Conectar a PostgreSQL
	db, err := config.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Error al conectar a PostgreSQL: %v", err)
	}
	log.Println("✅ Conectado a PostgreSQL")

	// Ejecutar migraciones
	if err := postgres.RunMigrations(db); err != nil {
		log.Fatalf("Error al ejecutar migraciones: %v", err)
	}
	log.Println("✅ Migraciones ejecutadas")

	// Crear cliente S3 (MinIO)
	s3Client, err := config.NewS3Client(cfg.S3)
	if err != nil {
		log.Fatalf("Error al crear cliente S3: %v", err)
	}
	log.Println("✅ Cliente S3 creado")

	// Crear cliente DynamoDB
	dynamoClient, err := config.NewDynamoDBCient(cfg.DynamoDB)
	if err != nil {
		log.Fatalf("Error al crear cliente DynamoDB: %v", err)
	}
	log.Println("✅ Cliente DynamoDB creado")

	// Inicializar repositorios
	alumnoRepo := postgres.NewAlumnoRepository(db)
	profesorRepo := postgres.NewProfesorRepository(db)
	sesionRepo := dynamodb.NewSesionRepository(dynamoClient, cfg.DynamoDB.TableName)
	fileStorage := s3.NewFileStorage(s3Client, cfg.S3.BucketName, cfg.S3.Endpoint)

	// Crear tabla DynamoDB y bucket S3 si no existen
	ctx := context.Background()
	if err := sesionRepo.CreateTable(ctx); err != nil {
		log.Printf("⚠️  Tabla DynamoDB ya existe o error: %v", err)
	}
	if err := fileStorage.CreateBucket(ctx); err != nil {
		log.Printf("⚠️  Bucket S3 ya existe o error: %v", err)
	}

	// Inicializar notificador (mock para desarrollo local)
	var notifier *aws.SNSMock
	if cfg.SNS.Mock {
		notifier = aws.NewSNSMock()
		log.Println("✅ SNS Mock habilitado")
	}

	// Inicializar casos de uso
	alumnoUseCase := usecase.NewAlumnoUseCase(alumnoRepo, fileStorage, notifier)
	profesorUseCase := usecase.NewProfesorUseCase(profesorRepo)
	sesionUseCase := usecase.NewSesionUseCase(sesionRepo, alumnoRepo)

	// Inicializar handlers
	alumnoHandler := handler.NewAlumnoHandler(alumnoUseCase)
	profesorHandler := handler.NewProfesorHandler(profesorUseCase)
	sesionHandler := handler.NewSesionHandler(sesionUseCase)

	// Configurar router
	router := apphttp.NewRouter(alumnoHandler, profesorHandler, sesionHandler)
	r := router.Setup()

	// Configurar servidor
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("🚀 Servidor iniciado en http://localhost:%s", cfg.Server.Port)
		log.Println("📚 Endpoints disponibles:")
		log.Println("   GET/POST       /alumnos")
		log.Println("   GET/PUT/DELETE /alumnos/{id}")
		log.Println("   POST           /alumnos/{id}/fotoPerfil")
		log.Println("   POST           /alumnos/{id}/email")
		log.Println("   POST           /alumnos/{id}/session/login")
		log.Println("   POST           /alumnos/{id}/session/verify")
		log.Println("   POST           /alumnos/{id}/session/logout")
		log.Println("   GET/POST       /profesores")
		log.Println("   GET/PUT/DELETE /profesores/{id}")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar servidor: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Apagando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error al apagar servidor: %v", err)
	}

	log.Println("👋 Servidor apagado correctamente")
}
