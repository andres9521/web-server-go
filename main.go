package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

// Definimos la estructura Productos y sus atributos con la letra inicial en mayuscula para que sean de tipo publico.
type Producto struct {
	Id              int       `json:"id"`
	Nombre          string    `json:"nombre"`
	Precio          float64   `json:"precio"`
	Stock           int       `json:"stock"`
	Codigo          string    `json:"codigo"`
	Publicado       bool      `json:"publicado"`
	FechaDeCreacion time.Time `json:"fecha_de_creacion"`
}

func main() {
	//Crear un router con Gin
	r := gin.Default()

	// Definir la ruta /productos
	r.GET("/productos", productosHandler)

	// Ejecutar la aplicación en el puerto 8081
	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}

}

func cargarProductos() ([]Producto, error) {
	var productos []Producto

	// Leer el archivo productos.json
	data, err := os.ReadFile("productos.json")
	if err != nil {
		return nil, err
	}

	// Parsear el contenido del archivo JSON y almacenar los productos en la variable correspondiente
	if err := json.Unmarshal(data, &productos); err != nil {
		return nil, err
	}
	return productos, nil
}

func productosHandler(c *gin.Context) {
	productos, err := cargarProductos()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productos)
}
