package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
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
	r := gin.Default()

	// Definir la ruta /productos
	r.GET("/productos", productosHandler)

	// Definir la ruta /productparams
	r.GET("/productparams", productParamsHandler)

	// Definir la ruta /products/:id
	r.GET("/products/:id", productoPorIDHandler)

	// Definir la ruta /searchbyquantity
	r.GET("/searchbyquantity", productosPorCantidadHandler)

	// Definir la ruta /buy
	r.GET("/buy", compraHandler)

	// Ejecutar la aplicación en el puerto 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func productosHandler(c *gin.Context) {
	productos, err := cargarProductos()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productos)
}

func productParamsHandler(c *gin.Context) {
	// Obtener los valores de los parámetros
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id debe ser un número entero"})
		return
	}
	nombre := c.Query("nombre")
	precioStr := c.Query("precio")
	precio, err := strconv.ParseFloat(precioStr, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "precio debe ser un número"})
		return
	}
	stockStr := c.Query("stock")
	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "stock debe ser un número entero"})
		return
	}
	codigo := c.Query("codigo")
	publicadoStr := c.Query("publicado")
	publicado, err := strconv.ParseBool(publicadoStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "publicado debe ser true o false"})
		return
	}
	fechaDeCreacionStr := c.Query("fechaDeCreacion")
	fechaDeCreacion, err := time.Parse(time.RFC3339, fechaDeCreacionStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "fechaDeCreacion debe estar en formato RFC3339"})
		return
	}

	// Crear un nuevo producto con los valores de los parámetros
	producto := Producto{
		Id:              id,
		Nombre:          nombre,
		Precio:          precio,
		Stock:           stock,
		Codigo:          codigo,
		Publicado:       publicado,
		FechaDeCreacion: fechaDeCreacion,
	}

	// Devolver el producto en formato JSON
	c.JSON(http.StatusOK, producto)
}

func cargarProductos() ([]Producto, error) {
	// Leer el contenido del archivo productos.json
	contenido, err := os.ReadFile("productos.json")
	if err != nil {
		return nil, err
	}

	// Parsear el contenido en un slice de Producto
	productos := []Producto{}
	if err := json.Unmarshal(contenido, &productos); err != nil {
		return nil, err
	}

	// Agregar el último producto a la lista
	ultimoProducto := Producto{
		Id:              7,
		Nombre:          "Producto 7",
		Precio:          99.99,
		Stock:           10,
		Codigo:          "P007",
		Publicado:       true,
		FechaDeCreacion: time.Now(),
	}
	productos = append(productos, ultimoProducto)

	return productos, nil
}

func productoPorIDHandler(c *gin.Context) {
	// Obtener el ID del producto desde el parámetro de la ruta
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id debe ser un número entero"})
		return
	}

	// Cargar la lista de productos
	productos, err := cargarProductos()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Buscar el producto correspondiente
	var productoEncontrado *Producto
	for i := range productos {
		if productos[i].Id == id {
			productoEncontrado = &productos[i]
			break
		}
	}
	if productoEncontrado == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "producto no encontrado"})
		return
	}

	// Devolver el producto en formato JSON
	c.JSON(http.StatusOK, productoEncontrado)
}

func productosPorCantidadHandler(c *gin.Context) {
	// Obtener los límites de la cantidad de stock desde los parámetros de la ruta
	min, err := strconv.Atoi(c.Query("min"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "min debe ser un número entero"})
		return
	}
	max, err := strconv.Atoi(c.Query("max"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "max debe ser un número entero"})
		return
	}

	// Cargar la lista de productos
	productos, err := cargarProductos()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Buscar los productos que cumplan con la condición
	productosEncontrados := []Producto{}
	for _, p := range productos {
		if p.Stock >= min && p.Stock <= max {
			productosEncontrados = append(productosEncontrados, p)
		}
	}

	// Devolver la lista de productos en formato JSON
	c.JSON(http.StatusOK, productosEncontrados)
}

func compraHandler(c *gin.Context) {
	// Obtener los parámetros de la ruta
	code := c.Query("code_value")
	cantidad, err := strconv.Atoi(c.Query("cantidad"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "cantidad debe ser un número entero"})
		return
	}

	// Cargar la lista de productos
	productos, err := cargarProductos()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Buscar el producto por code_value
	var productoEncontrado *Producto
	for _, p := range productos {
		if p.Codigo == code {
			productoEncontrado = &p
			break
		}
	}

	// Verificar si se encontró el producto
	if productoEncontrado == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	// Calcular el precio total de la compra
	precioTotal := productoEncontrado.Precio * float64(cantidad)

	// Devolver el detalle de la compra
	detalle := gin.H{
		"nombre":      productoEncontrado.Nombre,
		"cantidad":    cantidad,
		"precioTotal": precioTotal,
	}
	c.JSON(http.StatusOK, detalle)
}
