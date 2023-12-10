package main

import (
	"backend-user/api/config"
	"backend-user/api/routes"
	"backend-user/api/services"
	"fmt"
	"log"

	"github.com/gin-contrib/cors" // Importa el paquete de CORS
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Carga las variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar las variables de entorno: ", err)
		return
	}

	// Inicializa Firebase y obtén el cliente de autenticación
	authClient, err := config.InitFirebase()
	if err != nil {
		fmt.Printf("Error initializing Firebase client: %v\n", err)
		return
	}

	// Initalize the Storage client
	if err := config.InitStorage(); err != nil {
		fmt.Printf("Error initializing Storage client: %v\n", err)
		return
	}

	// Inicializa la conexión a la base de datos PostgreSQL
	db := config.InitDatabase()

	// Configura los servicios y controladores
	authService := services.NewAuthService(authClient, db)

	// Crea una instancia de Gin
	r := gin.Default()

	// Configura CORS para permitir solicitudes desde el puerto 5173 y permitir el encabezado "Authorization"
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173",
		"http://localhost:5174",
		"http://localhost:5175",
		"http://localhost:5176",
		"https://usuario.tisw.cl",
		"https://admin.tisw.cl",
	}
	config.AllowHeaders = []string{"Authorization", "Content-Type"} // Permite los encabezados "Authorization" y "Content-Type"
	r.Use(cors.New(config))                                         // Aplica la configuración de CORS a tu router

	// Configura las rutas con Gin
	routes.ConfigureRegisterRoute(r, authService)      // Configura la ruta de registro
	routes.ConfigureLoginRoute(r, authService)         // Configura la ruta de inicio de sesión
	routes.ConfigureUserInfoRoute(r, authService)      // Configura la ruta de información del usuario
	routes.ConfigureAgendamientoRouter(r, authService) // Configura la ruta de agendamiento
	routes.ConfigureImagenesRouter(r, authService)     // Configura la ruta de imágenes
	// Configura y ejecuta tu servidor Gin aquí
	port := ":8080"
	fmt.Printf("Servidor en ejecución en http://localhost%s\n", port)
	log.Fatal(r.Run(port))
}
