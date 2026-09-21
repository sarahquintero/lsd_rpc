package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"lsd_rpc/ServidorDeAudios/controllers"
	"lsd_rpc/ServidorDeAudios/services"
	"lsd_rpc/ServidorDeAudios/storage"
)

func main() {
	almacen := storage.NewAudioStorage("audios")
	servicio := services.NewAudioService(almacen, "catalogo.json")
	ctrl := controllers.NewAudioController(servicio)

	router := mux.NewRouter()
	ctrl.RegistrarRutas(router)

	fmt.Println("Servidor de audios (REST) escuchando en :8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
