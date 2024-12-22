package main

import (
    "testing"
    "time"
    "fmt"
    "sync"
)

type CategoryMetrics struct {
    category        string
    landingTime    time.Duration
    processedPlanes int
    mutex          sync.Mutex
}

type TestMetrics struct {
    totalTime     time.Duration
    categoryStats map[string]*CategoryMetrics
    mutex         sync.Mutex
}

func NewTestMetrics() *TestMetrics {
    return &TestMetrics{
        categoryStats: make(map[string]*CategoryMetrics),
    }
}

// Test Case 1: 10 aviones de cada categoría
func TestCase1(t *testing.T) {
    fmt.Println("\n=== Ejecutando Test Case 1: 10A-10B-10C ===")
    metrics := runTestWithMetrics(t, map[string]int{
        "A": 10,
        "B": 10,
        "C": 10,
    })
    reportDetailedMetrics(t, "Test Case 1", metrics)
}

// Test Case 2: 20A-5B-5C
func TestCase2(t *testing.T) {
    fmt.Println("\n=== Ejecutando Test Case 2: 20A-5B-5C ===")
    metrics := runTestWithMetrics(t, map[string]int{
        "A": 20,
        "B": 5,
        "C": 5,
    })
    reportDetailedMetrics(t, "Test Case 2", metrics)
}

// Test Case 3: 5A-5B-20C
func TestCase3(t *testing.T) {
    fmt.Println("\n=== Ejecutando Test Case 3: 5A-5B-20C ===")
    metrics := runTestWithMetrics(t, map[string]int{
        "A": 5,
        "B": 5,
        "C": 20,
    })
    reportDetailedMetrics(t, "Test Case 3", metrics)
}

func runTestWithMetrics(t *testing.T, distribution map[string]int) *TestMetrics {
    airport = NewAirport(3, 5) // Reiniciar el aeropuerto
    metrics := NewTestMetrics()
    
    // Inicializar métricas para cada categoría
    for category := range distribution {
        metrics.categoryStats[category] = &CategoryMetrics{
            category: category,
        }
    }

    startTotalTime := time.Now()
    var wg sync.WaitGroup

    // Procesar cada categoría
    for category, count := range distribution {
        wg.Add(1)
        go func(cat string, cnt int) {
            defer wg.Done()
            categoryStart := time.Now()
            
            // Generar y procesar aviones de esta categoría
            planesProcessed := processPlanesForCategory(cat, cnt)
            
            // Actualizar métricas de la categoría
            catMetrics := metrics.categoryStats[cat]
            catMetrics.mutex.Lock()
            catMetrics.landingTime = time.Since(categoryStart)
            catMetrics.processedPlanes = planesProcessed
            catMetrics.mutex.Unlock()
            
        }(category, count)
    }

    wg.Wait()
    metrics.totalTime = time.Since(startTotalTime)
    
    return metrics
}

func processPlanesForCategory(category string, count int) int {
    processed := 0
    var passengers int
    
    switch category {
    case "A":
        passengers = 101
    case "B":
        passengers = 75
    case "C":
        passengers = 25
    }

    for i := 0; i < count; i++ {
        plane := &Plane{
            id:         len(airport.planesWaiting) + 1,
            passengers: passengers,
            category:   category,
        }
        
        // Simular el proceso de aterrizaje
        if processPlane(plane) {
            processed++
        }
        
        // Pequeña pausa para simular tiempo entre aviones
        time.Sleep(100 * time.Millisecond)
    }
    
    return processed
}

func processPlane(plane *Plane) bool {
    // Simular tiempo de aproximación
    time.Sleep(time.Duration(200+rand.Intn(300)) * time.Millisecond)
    
    // Intentar conseguir pista
    runway := getAvailableRunway()
    if runway == nil {
        return false
    }
    
    // Simular aterrizaje
    time.Sleep(time.Duration(500+rand.Intn(500)) * time.Millisecond)
    runway.mutex.Lock()
    runway.inUse = false
    runway.mutex.Unlock()
    
    // Intentar conseguir puerta
    gate := getAvailableGate()
    if gate == nil {
        return false
    }
    
    // Simular desembarque
    time.Sleep(time.Duration(300+rand.Intn(400)) * time.Millisecond)
    gate.mutex.Lock()
    gate.inUse = false
    gate.mutex.Unlock()
    
    return true
}

func reportDetailedMetrics(t *testing.T, testName string, metrics *TestMetrics) {
    t.Logf("\n=== Métricas detalladas para %s ===", testName)
    t.Logf("Tiempo total de ejecución: %v", metrics.totalTime)
    
    // Reportar métricas por categoría
    for category, stats := range metrics.categoryStats {
        t.Logf("\nCategoría %s:", category)
        t.Logf("  - Tiempo de procesamiento: %v", stats.landingTime)
        t.Logf("  - Aviones procesados: %d", stats.processedPlanes)
        t.Logf("  - Tiempo promedio por avión: %v", 
            stats.landingTime/time.Duration(stats.processedPlanes))
    }
    
    // Calcular y mostrar estadísticas adicionales
    var totalPlanes int
    for _, stats := range metrics.categoryStats {
        totalPlanes += stats.processedPlanes
    }
    
    t.Logf("\nEstadísticas globales:")
    t.Logf("  - Total de aviones procesados: %d", totalPlanes)
    t.Logf("  - Tiempo promedio por avión global: %v", 
        metrics.totalTime/time.Duration(totalPlanes))
}