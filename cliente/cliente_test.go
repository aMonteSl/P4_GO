package main

import (
	"testing"
	"time"
	"fmt"
)

// Test Case 1: 10 aviones de cada categoría
func TestCase1(t *testing.T) {
	fmt.Println("=== Ejecutando Test Case 1: 10A-10B-10C ===")
	airport = NewAirport(3, 5) // Reiniciar el aeropuerto
	
	// Generar aviones
	generatePlanesForTest(10, "A")
	generatePlanesForTest(10, "B")
	generatePlanesForTest(10, "C")
	
	// Simular operaciones
	runTestOperations(t)
}

// Test Case 2: 20A-5B-5C
func TestCase2(t *testing.T) {
	fmt.Println("=== Ejecutando Test Case 2: 20A-5B-5C ===")
	airport = NewAirport(3, 5) // Reiniciar el aeropuerto
	
	// Generar aviones
	generatePlanesForTest(20, "A")
	generatePlanesForTest(5, "B")
	generatePlanesForTest(5, "C")
	
	// Simular operaciones
	runTestOperations(t)
}

// Test Case 3: 5A-5B-20C
func TestCase3(t *testing.T) {
	fmt.Println("=== Ejecutando Test Case 3: 5A-5B-20C ===")
	airport = NewAirport(3, 5) // Reiniciar el aeropuerto
	
	// Generar aviones
	generatePlanesForTest(5, "A")
	generatePlanesForTest(5, "B")
	generatePlanesForTest(20, "C")
	
	// Simular operaciones
	runTestOperations(t)
}

func generatePlanesForTest(count int, category string) {
	for i := 0; i < count; i++ {
		var passengers int
		switch category {
		case "A":
			passengers = 101
		case "B":
			passengers = 75
		case "C":
			passengers = 25
		}
		
		plane := &Plane{
			id:         len(airport.planesWaiting) + 1,
			passengers: passengers,
			category:   category,
		}
		
		airport.planesWaiting = append(airport.planesWaiting, plane)
	}
}

func runTestOperations(t *testing.T) {
	startTime := time.Now()
	
	// Simular diferentes estados del aeropuerto
	states := []int{1, 4, 2, 5, 3, 6, 0, 9}
	for _, state := range states {
		handleMessage(fmt.Sprintf("%d", state))
		time.Sleep(2 * time.Second) // Dar tiempo para procesar
	}
	
	duration := time.Since(startTime)
	t.Logf("Tiempo total de ejecución: %v", duration)
	t.Logf("Total aviones procesados: %d", len(airport.planesWaiting))
}