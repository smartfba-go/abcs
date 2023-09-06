package failureschi

import (
	"net/http"

	"github.com/go-chi/render"

	"go.smartfba.io/abcs/pkg/failures"
	"go.smartfba.io/abcs/pkg/failures/failureshttp"
)

func Render(w http.ResponseWriter, r *http.Request, err error) {
	status, _ := failures.StatusFromError(err)

	statusCode := failureshttp.StatusCode(status.Code)

	render.Status(r, statusCode)

	render.JSON(w, r, status)
}
