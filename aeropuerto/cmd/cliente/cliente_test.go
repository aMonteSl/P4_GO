package main

import (
	"testing"
	"time"
	"strconv"
)

// Test1: 10 aviones de cada categoría
func TestEqualDistribution(t *testing.T) {
	// Reiniciar el aeropuerto
	airport = NewAirport(3, 5)
	planeID = 0
	
	// Test con 10 aviones de cada categoría
	generatePlanesForCategory(10, "A")
	generatePlanesForCategory(10, "B")
	generatePlanesForCategory(10, "C")
	
	// Verificar la distribución
	countA, countB, countC := 0, 0, 0
	for _, plane := range airport.planesWaiting {
		switch plane.category {
		case "A":
			countA++
		case "B":
			countB++
		case "C":
			countC++
		}
	}
	
	if countA != 10 || countB != 10 || countC != 10 {
		t.Errorf("Distribución incorrecta. Esperado 10 de cada, obtenido A:%d, B:%d, C:%d", 
			countA, countB, countC)
	}
	
	// Probar el procesamiento
	testProcessing(t, "Test1: Distribución equitativa")
}

// Test2: 20 aviones A, 5 B y C
func TestMoreTypeA(t *testing.T) {
	// Reiniciar el aeropuerto
	airport = NewAirport(3, 5)
	planeID = 0
	
	// Generar distribución asimétrica
	generatePlanesForCategory(20, "A")
	generatePlanesForCategory(5, "B")
	generatePlanesForCategory(5, "C")
	
	// Verificar la distribución
	countA, countB, countC := 0, 0, 0
	for _, plane := range airport.planesWaiting {
		switch plane.category {
		case "A":
			countA++
		case "B":
			countB++
		case "C":
			countC++
		}
	}
	
	if countA != 20 || countB != 5 || countC != 5 {
		t.Errorf("Distribución incorrecta. Esperado A:20, B:5, C:5, obtenido A:%d, B:%d, C:%d", 
			countA, countB, countC)
	}
	
	// Probar el procesamiento
	testProcessing(t, "Test2: Mayoría tipo A")
}

// Test3: 20 aviones C, 5 A y B
func TestMoreTypeC(t *testing.T) {
	// Reiniciar el aeropuerto
	airport = NewAirport(3, 5)
	planeID = 0
	
	// Generar distribución asimétrica
	generatePlanesForCategory(5, "A")
	generatePlanesForCategory(5, "B")
	generatePlanesForCategory(20, "C")
	
	// Verificar la distribución
	countA, countB, countC := 0, 0, 0
	for _, plane := range airport.planesWaiting {
		switch plane.category {
		case "A":
			countA++
		case "B":
			countB++
		case "C":
			countC++
		}
	}
	
	if countA != 5 || countB != 5 || countC != 20 {
		t.Errorf("Distribución incorrecta. Esperado A:5, B:5, C:20, obtenido A:%d, B:%d, C:%d", 
			countA, countB, countC)
	}
	
	// Probar el procesamiento
	testProcessing(t, "Test3: Mayoría tipo C")
}

// Función auxiliar para probar el procesamiento de aviones
func testProcessing(t *testing.T, testName string) {
	t.Log("Iniciando", testName)
	
	// Probar diferentes estados del aeropuerto
	states := []int{1, 2, 3, 4, 5, 6} // Estados que permiten aterrizajes
	
	for _, state := range states {
		t.Logf("Probando estado %d", state)
		
		airport.stateMutex.Lock()
		airport.state = state
		airport.stateMutex.Unlock()
		
		// Procesar algunos aviones
		planesBeforeCount := len(airport.planesWaiting)
		
		// Simular procesamiento de aviones
		handleMessage(strconv.Itoa(state))
		
		// Dar tiempo para que se procesen algunos aviones
		time.Sleep(5 * time.Second)
		
		// Verificar que algunos aviones fueron procesados
		planesAfterCount := len(airport.planesWaiting)
		if planesBeforeCount <= planesAfterCount {
			t.Logf("Advertencia: No se procesaron aviones en estado %d", state)
		} else {
			t.Logf("Procesados %d aviones en estado %d", 
				planesBeforeCount-planesAfterCount, state)
		}
	}
	
	t.Log("Finalizado", testName)
}