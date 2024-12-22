package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"math/rand"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
	msg    string
)

// Estructuras necesarias
type Plane struct {
	id         int
	passengers int
	category   string
}

type Airport struct {
	runways      []Runway
	gates        []Gate
	controlTower ControlTower
	state        int
	stateMutex   sync.RWMutex
	planesWaiting []*Plane
}

type Runway struct {
	id    int
	inUse bool
	mutex sync.Mutex
}

type Gate struct {
	id    int
	inUse bool
	mutex sync.Mutex
}

type ControlTower struct {
	queue []*Plane
	mutex sync.Mutex
}

// Variables globales del aeropuerto
var (
	airport *Airport
	planeID = 0
	planeMutex sync.Mutex
)

func init() {
	// Inicializar el aeropuerto con 3 pistas y 5 puertas
	airport = NewAirport(3, 5)
	// Generar aviones iniciales de prueba
	generateInitialPlanes()
}

func NewAirport(runwayCount, gateCount int) *Airport {
	airport := &Airport{
		runways:       make([]Runway, runwayCount),
		gates:         make([]Gate, gateCount),
		planesWaiting: make([]*Plane, 0),
	}
	
	for i := range airport.runways {
		airport.runways[i] = Runway{id: i}
	}
	
	for i := range airport.gates {
		airport.gates[i] = Gate{id: i}
	}
	
	return airport
}

func generateInitialPlanes() {
	// Generar aviones de cada categoría
	generatePlanesForCategory(10, "A") // Categoría A: >100 pasajeros
	generatePlanesForCategory(10, "B") // Categoría B: 50-100 pasajeros
	generatePlanesForCategory(10, "C") // Categoría C: <50 pasajeros
}

func generatePlanesForCategory(count int, category string) {
	for i := 0; i < count; i++ {
		var passengers int
		switch category {
		case "A":
			passengers = 101 + rand.Intn(100) // 101-200 pasajeros
		case "B":
			passengers = 50 + rand.Intn(51)   // 50-100 pasajeros
		case "C":
			passengers = 1 + rand.Intn(49)    // 1-49 pasajeros
		}
		
		planeMutex.Lock()
		planeID++
		plane := &Plane{
			id:         planeID,
			passengers: passengers,
			category:   category,
		}
		planeMutex.Unlock()
		
		airport.planesWaiting = append(airport.planesWaiting, plane)
	}
}

func handleMessage(msg string) {
    // Lista de mensajes del sistema a ignorar
    systemMessages := []string{
        "Aeropuerto localizado",
        "Se ha conectado",
        "se ha desconectado",
    }

    // Verificar si es un mensaje del sistema
    for _, sysMsg := range systemMessages {
        if strings.Contains(msg, sysMsg) {
            fmt.Println("Mensaje del sistema:", msg)
            return
        }
    }

    // Dividir el mensaje en líneas por si vienen múltiples números
    numbers := strings.Split(strings.TrimSpace(msg), "\n")
    
    for _, numStr := range numbers {
        // Limpiar el string y verificar si está vacío
        numStr = strings.TrimSpace(numStr)
        if numStr == "" {
            continue
        }

        // Convertir a número
        number, err := strconv.Atoi(numStr)
        if err != nil {
            fmt.Printf("Ignorando mensaje no numérico: %s\n", numStr)
            continue
        }

        fmt.Printf("Procesando estado: %d\n", number)
        
        // Actualizar el estado del aeropuerto
        airport.stateMutex.Lock()
        airport.state = number
        airport.stateMutex.Unlock()

        // Manejar el estado
        switch {
        case number == 0:
            fmt.Println("Aeropuerto Inactivo - No hay aterrizajes")
        case number >= 1 && number <= 3:
            category := string('A' + rune(number-1))
            fmt.Printf("Permitiendo aterrizaje solo categoría %s\n", category)
            handleCategoryOnly(category)
        case number >= 4 && number <= 6:
            category := string('A' + rune(number-4))
            fmt.Printf("Prioridad para categoría %s\n", category)
            handleCategoryPriority(category)
        case number == 7 || number == 8:
            fmt.Println("Manteniendo estado anterior")
        case number == 9:
            fmt.Println("Aeropuerto Cerrado Temporalmente")
        default:
            fmt.Printf("Estado no reconocido: %d - Manteniendo estado anterior\n", number)
        }
    }
}

func handleCategoryOnly(category string) {
	fmt.Printf("Procesando solo aviones de categoría %s\n", category)
	for _, plane := range airport.planesWaiting {
		if plane.category == category {
			go processPlane(plane)
		}
	}
}

func handleCategoryPriority(category string) {
	fmt.Printf("Prioridad para aviones de categoría %s\n", category)
	// Primero procesar aviones de la categoría prioritaria
	for _, plane := range airport.planesWaiting {
		if plane.category == category {
			go processPlane(plane)
		}
	}
	// Luego procesar el resto
	time.Sleep(2 * time.Second) // Dar tiempo a los prioritarios
	for _, plane := range airport.planesWaiting {
		if plane.category != category {
			go processPlane(plane)
		}
	}
}

func processPlane(plane *Plane) {
	// Simular tiempo de aproximación
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	
	// Buscar pista disponible
	runway := getAvailableRunway()
	if runway == nil {
		fmt.Printf("Avión %d en espera - No hay pistas disponibles\n", plane.id)
		return
	}
	
	// Aterrizar
	fmt.Printf("Avión %d (Categoría %s) aterrizando en pista %d\n", 
		plane.id, plane.category, runway.id)
	time.Sleep(time.Duration(2+rand.Intn(3)) * time.Second)
	
	// Liberar pista
	runway.mutex.Lock()
	runway.inUse = false
	runway.mutex.Unlock()
	
	// Buscar puerta
	gate := getAvailableGate()
	if gate == nil {
		fmt.Printf("Avión %d en espera de puerta\n", plane.id)
		return
	}
	
	// Proceso de desembarque
	fmt.Printf("Avión %d desembarcando en puerta %d\n", plane.id, gate.id)
	time.Sleep(time.Duration(3+rand.Intn(4)) * time.Second)
	
	// Liberar puerta
	gate.mutex.Lock()
	gate.inUse = false
	gate.mutex.Unlock()
	
	fmt.Printf("Avión %d ha completado sus operaciones\n", plane.id)
}

func getAvailableRunway() *Runway {
	for i := range airport.runways {
		airport.runways[i].mutex.Lock()
		if !airport.runways[i].inUse {
			airport.runways[i].inUse = true
			airport.runways[i].mutex.Unlock()
			return &airport.runways[i]
		}
		airport.runways[i].mutex.Unlock()
	}
	return nil
}

func getAvailableGate() *Gate {
	for i := range airport.gates {
		airport.gates[i].mutex.Lock()
		if !airport.gates[i].inUse {
			airport.gates[i].inUse = true
			airport.gates[i].mutex.Unlock()
			return &airport.gates[i]
		}
		airport.gates[i].mutex.Unlock()
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		logger.Fatal(err)
	}
	defer conn.Close()
	
	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			msg = string(buf[:n])
			fmt.Println("len: " + strconv.Itoa(n) + " msg: " + msg)
			handleMessage(msg)
		}
	}
}