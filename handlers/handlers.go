package handlers

import (
	"fmt"
	HandlerContext "go-webapp-starter/context"
	"go-webapp-starter/utils"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func VersionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, hctx *HandlerContext.Context) {
	data, err := ioutil.ReadFile("version")
	utils.CheckErr(err, "versionHandler", "Unable to read version file")
	fmt.Fprintf(w, "Version: %s", data)
}
