package wildlifenl

import (
	"log"
	"runtime"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/lib/pq"
)

func generateDescription(text string, scopes []string) string {
	if len(scopes) == 0 {
		return text
	}
	result := make([]string, 0)
	for _, value := range scopes {
		result = append(result, "`"+value+"`")
	}
	return text + "<br/><br/>**Scopes**<br/>" + strings.Join(result, ", ")
}

func handleError(err error) error {
	switch typedError := err.(type) {
	case *pq.Error:
		message := typedError.Message
		detail := typedError.Detail
		if strings.Contains(message, "violates foreign key constraint") {
			detail = detail[4:strings.LastIndex(detail, ")")+1] + " does not exist."
			return huma.Error400BadRequest(detail)
		}
	}
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Println("ERROR", frame.File+" "+strconv.Itoa(frame.Line)+" "+frame.Function+": ", err)
	return huma.Error500InternalServerError("")
}

func generateNotFoundByIDError(objectType string, id string) error {
	return huma.Error404NotFound("No " + objectType + " with ID " + id + " was found")
}

func generateNotFoundForThisUserError(objectType string, id string) error {
	return huma.Error404NotFound("No " + objectType + " with ID " + id + " was found for the current user")
}
