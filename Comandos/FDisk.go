package Comandos

import (
	"fmt"
	"strconv"
	"strings"
)

type Transition struct {
	partition int
	start     int
	end       int
	before    int
	after     int
}

var startValue int

func ValidarDatosFDISK(tokens []string) {
	size := ""
	unit := "k"
	driveletter := "disco"
	name := ""

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Comparar(tk[0], "size") {
			size = tk[1]
		} else if Comparar(token, "unit") {
			unit = tk[1]
		} else if Comparar(token, "driveletter") {
			driveletter = tk[1]
		} else if Comparar(token, "name") {
			name = tk[1]
		}
	}
	if size == "" || driveletter == "" || name == "" {
		Error("FDISK", "El comando FDISK necesita parámetros obligatorios: Size, Driveletter y Name")
	} else {
		GenerarParticion(size, unit, driveletter, name)
	}
}

func GenerarParticion(s string, u string, dl string, n string) {
	startValue = 0
	i, error_ := strconv.Atoi(s)
	if error_ != nil {
		Error("FDISK", "SIZE debe ser un número entero")
		return
	}

	if i <= 0 {
		Error("FDISK", "SIZE debe ser un número entero")
		return
	}

	if Comparar(u, "b") || Comparar(u, "k") || Comparar(u, "m") {
		if Comparar(u, "k") {
			i = i * 1024
			fmt.Println("Partición en KB")
		} else if Comparar(u, "m") {
			i = i * 1024 * 1024
			fmt.Println("Partición en MB")
		}
	} else {
		Error("FDISK", "Unit no contiene los valores esperados")
	}
}
