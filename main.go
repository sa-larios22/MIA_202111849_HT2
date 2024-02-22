package main

import (
	"HT2_MIA_202111849/Comandos"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// exec -path=/home/daniel/Escritorio/MIA_1S2024/Ejemplos_Proyecto/prueba.adjs

func main() {
	for {
		fmt.Println("======================== HOJA DE TRABAJO 2 ========================")
		fmt.Println("======================= INGRESE UN COMANDO ========================")
		fmt.Println("============ Puede finalizar la aplicación con \"exit\" ===========")
		fmt.Println("\t")

		reader := bufio.NewReader(os.Stdin)
		entrada, _ := reader.ReadString('\n')
		eleccion := strings.TrimRight(entrada, "\r\n")
		if eleccion == "exit" {
			break
		}
		comando := Comando(eleccion)
		fmt.Println("Comando: " + comando)

		eleccion = strings.TrimSpace(eleccion)
		fmt.Println("TrimSpace: " + eleccion)

		eleccion = strings.TrimLeft(eleccion, comando)
		fmt.Println("TrimLeft: " + eleccion)

		tokens := SepararTokens(eleccion)
		fmt.Println(tokens)

		funciones(comando, tokens)

		fmt.Println("\tPresione Enter para continuar....")
		fmt.Scanln()
	}
}

func Comando(text string) string {
	var tkn string
	terminar := false
	for i := 0; i < len(text); i++ {
		if terminar {
			if string(text[i]) == " " || string(text[i]) == "-" {
				break
			}
			tkn += string(text[i])
		} else if string(text[i]) != " " && !terminar {
			if string(text[i]) == "#" {
				tkn = text
			} else {
				tkn += string(text[i])
				terminar = true
			}
		}
	}
	return tkn
}

func SepararTokens(texto string) []string {
	var tokens []string
	if texto == "" {
		return tokens
	}
	texto += " "
	var token string
	estado := 0
	for i := 0; i < len(texto); i++ {
		c := string(texto[i])
		if estado == 0 && c == "-" {
			estado = 1
		} else if estado == 0 && c == "#" {
			continue
		} else if estado != 0 {
			if estado == 1 {
				if c == "=" {
					estado = 2
				} else if c == " " {
					continue
				} else if (c == "P" || c == "p") && string(texto[i+1]) == " " && string(texto[i-1]) == "-" {
					estado = 0
					tokens = append(tokens, c)
					token = ""
					continue
				} else if (c == "R" || c == "r") && string(texto[i+1]) == " " && string(texto[i-1]) == "-" {
					estado = 0
					tokens = append(tokens, c)
					token = ""
					continue
				}
			} else if estado == 2 {
				if c == " " {
					continue
				}
				if c == "\"" {
					estado = 3
					continue
				} else {
					estado = 4
				}
			} else if estado == 3 {
				if c == "\"" {
					estado = 4
					continue
				}
			} else if estado == 4 && c == "\"" {
				tokens = []string{}
				continue
			} else if estado == 4 && c == " " {
				estado = 0
				tokens = append(tokens, token)
				token = ""
				continue
			}
			token += c
		}
	}
	return tokens
}

func funciones(token string, tks []string) {

	fmt.Println("\t")
	fmt.Println("===== FUNC FUNCIONES =====")
	fmt.Println("STRING TOKEN: ")
	fmt.Println(token)
	fmt.Println("STRING TKS: ")
	fmt.Println(tks)
	fmt.Println("\t")

	if token != "" {
		if Comandos.Comparar(token, "EXECUTE") {
			fmt.Println("======================== FUNCIÓN EXECUTE ========================")
			FuncionExec(tks)
		} else if Comandos.Comparar(token, "MKDISK") {
			fmt.Println("======================== FUNCIÓN MKDISK ========================")
			Comandos.ValidarDatosMKDISK(tks)
		} else if Comandos.Comparar(token, "FMDISK") {
			fmt.Println("======================== FUNCIÓN FDISK ========================")
			Comandos.ValidarDatosFDISK(tks)
		} else {
			Comandos.Error("ANALIZADOR", "No se reconoce el comando \""+token+"\"")
		}
	}
}

func FuncionExec(tokens []string) {

	path := ""
	for i := 0; i < len(tokens); i++ {
		datos := strings.Split(tokens[i], "=")
		if Comandos.Comparar(datos[0], "path") {
			path = datos[1]
		}
	}
	if path == "" {
		Comandos.Error("EXEC", "Se requiere el parámetro \"path\" para este comando")
		return
	}
	Exec(path)
}

func Exec(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %s", err)
	}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		texto := fileScanner.Text()
		texto = strings.TrimSpace(texto)
		tk := Comando(texto)
		if texto != "" {
			if Comandos.Comparar(tk, "pause") {
				fmt.Println("========================= FUNCIÓN PAUSE =========================")
				var pause string
				Comandos.Mensaje("Pausa", "Presione \"enter\" para continuar...")
				fmt.Scanln(&pause)
				continue
			} else if string(texto[0]) == "#" {
				fmt.Println("========================= COMENTARIO =========================")
				Comandos.Mensaje("Comentario de Texto", texto)
				continue
			}
			texto = strings.TrimLeft(texto, tk)
			tokens := SepararTokens(texto)
			funciones(tk, tokens)
		}
	}
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error al leer el archivo: %s", err)
	}
	file.Close()
}
