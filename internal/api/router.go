package api

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	ENDPOINT_PREFIX := "/api"

	public := r.PathPrefix(ENDPOINT_PREFIX).Subrouter()
	public.HandleFunc("/status", statusHandler).Methods("GET")

	public.HandleFunc("/register", registerHandler).Methods("POST")
	public.HandleFunc("/login", loginHandler).Methods("POST")
	public.HandleFunc("/refresh", refreshHandler).Methods("POST")
	public.HandleFunc("/logout", logoutHandler).Methods("POST")

	protected := r.PathPrefix(ENDPOINT_PREFIX).Subrouter()
	protected.Use(AuthMiddleware)

	protected.HandleFunc("/vaults", vaultsHandler).Methods("GET")
	protected.HandleFunc("/vault/{vaultId}", vaultHandler).Methods("GET")

	protected.HandleFunc("/collections/{vaultId}", collectionsHandler).Methods("GET")

	protected.HandleFunc("/videos/{collectionId}", videosHandler).Methods("GET")
	protected.HandleFunc("/video/{videoId}", videoHandler).Methods("GET")

	protected.HandleFunc("/galleries/{vaultId}", galleriesHandler).Methods("GET")
	protected.HandleFunc("/gallery/{galleryId}", galleryHandler).Methods("GET")

	protected.HandleFunc("/actors/{vaultId}", actorsHandler).Methods("GET")

	return r
}
