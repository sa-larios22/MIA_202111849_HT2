package Comandos

import (
	"HT2_MIA_202111849/Structs"
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

func ValidarDatosMKDISK(tokens []string) {
	fmt.Println("\t")
	fmt.Println("\t===== VALIDAR DATOS MK DISK =====")
	fmt.Println("\tSTRING TOKENS")
	fmt.Println("\t" + strings.Join(tokens, ", "))
	fmt.Println("\t")

	size := ""
	unit := ""
	error_ := false
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		if Comparar(tk[0], "size") {
			if size == "" {
				size = tk[1]
			} else {
				Error("MKDISK", "parametro SIZE repetido en el comando: "+tk[0])
				return
			}
		} else if Comparar(tk[0], "unit") {
			if unit == "" {
				unit = tk[1]
			} else {
				Error("MKDISK", "parametro U repetido en el comando: "+tk[0])
				return
			}
		} else {
			Error("MKDISK", "no se esperaba el parametro "+tk[0])
			error_ = true
			return
		}
	}
	if unit == "" {
		unit = "M"
	}
	if error_ {
		return
	}
	if size == "" {
		Error("MKDISK", "se requiere parametro Size para este comando")
		return
	} else if !Comparar(unit, "k") && !Comparar(unit, "m") {
		Error("MKDISK", "valores en parametro unit no esperados")
		return
	} else {
		makeFile(size, unit)
	}
}

func makeFile(s string, u string) {
	var disco = Structs.NewMBR()
	size, err := strconv.Atoi(s)
	if err != nil {
		Error("MKDISK", "Size debe ser un número entero")
		return
	}
	if size <= 0 {
		Error("MKDISK", "Size debe ser mayor a 0")
		return
	}
	if Comparar(u, "M") {
		size = 1024 * 1024 * size
	} else if Comparar(u, "k") {
		size = 1024 * size
	}
	var fecha = time.Now().Format("2006-01-02 15:04:05")
	copy(disco.Mbr_fecha_creacion[:], fecha)
	aleatorio, _ := rand.Int(rand.Reader, big.NewInt(999999999))
	entero, _ := strconv.Atoi(aleatorio.String())
	disco.Mbr_dsk_signature = int64(entero)
	copy(disco.Dsk_fit[:], "f")
	disco.Mbr_partition_1 = Structs.NewParticion()
	disco.Mbr_partition_2 = Structs.NewParticion()
	disco.Mbr_partition_3 = Structs.NewParticion()
	disco.Mbr_partition_4 = Structs.NewParticion()

	// Asegurarse de que el tamaño de disco se ajuste según la unidad
	num := int64(size)

	if ArchivoExiste("disco") {
		_ = os.Remove("disco")
	}

	file, err := os.Create("disco")
	defer file.Close()
	if err != nil {
		Error("MKDISK", "No se pudo crear el disco.")
		return
	}
	var vacio int8 = 0
	s1 := &vacio
	num = num - 1
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s1)
	EscribirBytes(file, binario.Bytes())

	file.Seek(num, 0)

	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s1)
	EscribirBytes(file, binario2.Bytes())

	file.Seek(0, 0)
	disco.Mbr_tamano = num + 1

	var binario3 bytes.Buffer
	binary.Write(&binario3, binary.BigEndian, disco)
	EscribirBytes(file, binario3.Bytes())
	file.Close()
	Mensaje("MKDISK", "¡Disco creado correctamente!")
}

func ValidarDatosFDISK(tokens []string) {
	fmt.Println("\t")
	fmt.Println("\t===== VALIDAR DATOS FDISK =====")
	fmt.Println("\tSTRING TOKENS")
	fmt.Println("\t" + strings.Join(tokens, ", "))
	fmt.Println("\t")

	size := ""
	unit := ""
	driveLetter := ""
	name := ""
	error_ := false

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		tk := strings.Split(token, "=")
		switch tk[0] {
		case "size":
			if size == "" {
				size = tk[1]
				// Verificar si el tamaño es válido
				if _, err := strconv.Atoi(size); err != nil {
					Error("FDISK", "El tamaño debe ser un número válido")
					return
				}
				sizeInt, _ := strconv.Atoi(size)
				if sizeInt <= 0 {
					Error("FDISK", "El tamaño debe ser mayor que 0")
					return
				}
			} else {
				Error("FDISK", "Parametro SIZE repetido en el comando: "+tk[0])
				return
			}
		case "unit":
			if unit == "" {
				unit = tk[1]
				// Verificar si la unidad es válida
				if unit != "B" && unit != "K" && unit != "M" {
					Error("FDISK", "Unidad no válida. Las opciones válidas son B, K o M")
					return
				}
			} else {
				Error("FDISK", "Parametro UNIT repetido en el comando: "+tk[0])
				return
			}
		case "driveletter":
			if driveLetter == "" {
				driveLetter = tk[1]
				// Verificar si la letra de unidad es válida (por ejemplo, si el disco existe)
				if !ExisteDisco(driveLetter) {
					Error("FDISK", "El disco "+driveLetter+" no existe")
					return
				}
			} else {
				Error("FDISK", "Parametro DRIVELETTER repetido en el comando: "+tk[0])
				return
			}
		case "name":
			if name == "" {
				name = tk[1]
				// Verificar si el nombre de la partición ya existe
				if ExisteParticion(name) {
					Error("FDISK", "La partición "+name+" ya existe")
					return
				}
			} else {
				Error("FDISK", "Parametro NAME repetido en el comando: "+tk[0])
				return
			}
		default:
			Error("FDISK", "No se esperaba el parametro "+tk[0])
			error_ = true
			return
		}
	}

	if error_ {
		return
	}

	// Continuar con la lógica de FDISK
	// Por ejemplo, aquí se realizarían las operaciones de particionado, y se mostrarían los errores correspondientes si no se pueden realizar.

	// Si la operación es exitosa, mostrar un mensaje de éxito
	// Mensaje("FDISK", "La partición se creó/modificó correctamente")
}
