package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type LabCreateRequest struct {
	Image                string    `json:"image"`
	ExpiryTime           time.Time `json:"expiry_time"`
	WebTerminalEnabled   bool      `json:"web_terminal_enabled"`
	CodeServerEnabled    bool      `json:"code_server_enabled"`
	VNCEnabled           bool      `json:"vnc_enabled"`
	PortProxyEnabled     bool      `json:"port_proxy_enabled"`
	EnvironmentVariables string    `json:"environment_variables"`
}

type LabInfo struct {
	ID         string    `json:"id"`
	Status     LabStatus `json:"status"`
	ExpiryTime time.Time `json:"expiry_time"`
}

func StartAPIServer() {
	e := echo.New()
	e.GET("/status/:lab_id", func(c echo.Context) error {
		lab_id := c.Param("lab_id")
		lab, err := GetLabByID(lab_id)
		if err != nil {
			return c.String(http.StatusNotFound, "Lab not found")
		}
		return c.JSON(200, LabInfo{
			ID:         lab.ID,
			Status:     lab.Status,
			ExpiryTime: lab.ExpiryTime,
		})
	})
	e.POST("/start", func(c echo.Context) error {
		var labCreateRequest LabCreateRequest
		if err := c.Bind(&labCreateRequest); err != nil {
			return c.String(http.StatusBadRequest, "Invalid request")
		}
		lab := Lab{
			Image:                labCreateRequest.Image,
			ExpiryTime:           labCreateRequest.ExpiryTime,
			WebTerminalEnabled:   labCreateRequest.WebTerminalEnabled,
			CodeServerEnabled:    labCreateRequest.CodeServerEnabled,
			VNCEnabled:           labCreateRequest.VNCEnabled,
			PortProxyEnabled:     labCreateRequest.PortProxyEnabled,
			EnvironmentVariables: labCreateRequest.EnvironmentVariables,
		}
		// check if ExpiryTime is in the future
		if lab.ExpiryTime.Before(time.Now()) {
			return c.String(http.StatusBadRequest, "Expiry time is in the past")
		}
		err := lab.Create()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create lab")
		}
		return c.JSON(200, LabInfo{
			ID:         lab.ID,
			Status:     lab.Status,
			ExpiryTime: lab.ExpiryTime,
		})
	})
	e.POST("/stop", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":8888"))
}
