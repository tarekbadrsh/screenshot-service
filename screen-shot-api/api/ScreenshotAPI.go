package api

import (
	"fmt"
	"net/http"
	"time"

	"screen-shot-api/bll"
	"screen-shot-api/dto"
	"screen-shot-api/logger"

	"github.com/julienschmidt/httprouter"
)

func configScreenshotsRouter(routes *[]route) {
	*routes = append(*routes, route{method: "GET", path: "/screenshots", handle: getAllScreenshots})
	*routes = append(*routes, route{method: "POST", path: "/screenshots", handle: postScreenshots})
	*routes = append(*routes, route{method: "PUT", path: "/screenshots", handle: putScreenshots})
	*routes = append(*routes, route{method: "GET", path: "/screenshots/:id", handle: getScreenshots})
	*routes = append(*routes, route{method: "DELETE", path: "/screenshots/:id", handle: deleteScreenshots})
}

func getAllScreenshots(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	screenshots, err := bll.GetAllScreenshots()
	if err != nil {
		logger.Error(err)
		writeResponseError(w, err, http.StatusInternalServerError)
		return
	}
	writeResponseJSON(w, screenshots, http.StatusOK)
}

func getScreenshots(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	requestID := ps.ByName("id")
	id, err := bll.ConvertID(requestID)
	if err != nil {
		msg := fmt.Errorf("Error: parameter (id) should be int32; Id=%v; err (%v)", requestID, err)
		logger.Error(msg)
		writeResponseError(w, msg, http.StatusBadRequest)
		return
	}

	screenshot, err := bll.GetScreenshot(id)
	if err != nil {
		msg := fmt.Errorf("Can’t find screenshot (%v); err (%v)", id, err)
		logger.Error(msg)
		writeResponseError(w, msg, http.StatusNotFound)
		return
	}
	writeResponseJSON(w, screenshot, http.StatusOK)
}

func postScreenshots(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	screenshot := &dto.ScreenshotDTO{}
	if err := readJSON(r, screenshot); err != nil {
		logger.Error(err)
		writeResponseError(w, err, http.StatusBadRequest)
		return
	}
	screenshot.CreatedAt = time.Now().Unix()
	result, err := bll.CreateScreenshot(screenshot)
	if err != nil {
		logger.Error(err)
		writeResponseError(w, err, http.StatusInternalServerError)
		return
	}
	writeResponseJSON(w, result, http.StatusCreated)
}

func putScreenshots(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	screenshot := &dto.ScreenshotDTO{}
	if err := readJSON(r, screenshot); err != nil {
		logger.Error(err)
		writeResponseError(w, err, http.StatusBadRequest)
		return
	}

	result, err := bll.UpdateScreenshot(screenshot)
	if err != nil {
		logger.Error(err)
		writeResponseError(w, err, http.StatusBadRequest)
		return
	}
	writeResponseJSON(w, result, http.StatusOK)
}

func deleteScreenshots(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	requestID := ps.ByName("id")
	id, err := bll.ConvertID(requestID)
	if err != nil {
		msg := fmt.Errorf("Error: parameter (id) should be int32; Id=%v; err (%v)", requestID, err)
		logger.Error(msg)
		writeResponseError(w, msg, http.StatusBadRequest)
		return
	}

	err = bll.DeleteScreenshot(id)
	if err != nil {
		msg := fmt.Errorf("Screenshot with id (%v) does not exist; err (%v)", id, err)
		logger.Error(msg)
		writeResponseError(w, msg, http.StatusNotFound)
		return

	}
	writeResponseJSON(w, true, http.StatusOK)
}
