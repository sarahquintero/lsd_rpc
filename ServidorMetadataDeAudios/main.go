package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"lsd_rpc/ServidorMetadataDeAudios/controllers"
	"lsd_rpc/ServidorMetadataDeAudios/repository"
	"lsd_rpc/ServidorMetadataDeAudios/services"
)

func main() {
	repo, err := repository.NewAudioRepository("catalogo.json")
	if err != nil {
		log.Fatalf("No se pudo cargar catálogo: %v", err)
	}
	servicio := services.NewAudioService(repo)
	controlador := controllers.NewAudioController(servicio)

	router := mux.NewRouter()
	controlador.RegistrarRutas(router)

	fmt.Println("Servidor de metadatos escuchando en :8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
