package integration

import (
    "os"
    "os/exec"
    "testing"
    "time"
    "fmt"
    "net"
    "path/filepath"
    "log"
    "syscall"
)

// Variables globales para los procesos
var (
    serverCmd *exec.Cmd
    clientCmd *exec.Cmd
    enaireCmd *exec.Cmd
)

func getProjectRoot() string {
    currentDir, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }
    return filepath.Join(currentDir, "..", "..")
}

func killProcess(cmd *exec.Cmd) {
    if cmd != nil && cmd.Process != nil {
        // Enviar SIGTERM primero
        cmd.Process.Signal(syscall.SIGTERM)
        
        // Esperar un poco y luego forzar el cierre si es necesario
        time.Sleep(time.Second)
        cmd.Process.Kill()
    }
}

func cleanupProcesses() {
    killProcess(clientCmd)
    killProcess(enaireCmd)
    killProcess(serverCmd)
}

func TestIntegration(t *testing.T) {
    defer cleanupProcesses()
    
    projectRoot := getProjectRoot()
    t.Logf("Directorio del proyecto: %s", projectRoot)
    
    // Iniciar servidor
    serverPath := filepath.Join(projectRoot, "cmd", "servidor", "servidor.go")
    serverCmd = exec.Command("go", "run", serverPath)
    serverCmd.Stdout = os.Stdout
    serverCmd.Stderr = os.Stderr
    
    if err := serverCmd.Start(); err != nil {
        t.Fatalf("Error iniciando servidor: %v", err)
    }
    
    // Esperar a que el servidor esté listo
    time.Sleep(3 * time.Second)

    // Ejecutar las configuraciones de prueba
    configs := []struct {
        name string
        numA int
        numB int
        numC int
    }{
        {"10_each", 10, 10, 10},
        {"20A_5BC", 20, 5, 5},
        {"5AB_20C", 5, 5, 20},
    }

    for _, cfg := range configs {
        t.Run(cfg.name, func(t *testing.T) {
            runTestConfiguration(t, cfg.numA, cfg.numB, cfg.numC, projectRoot)
            time.Sleep(5 * time.Second) // Espera entre pruebas
        })
    }
}

func runTestConfiguration(t *testing.T, numA, numB, numC int, projectRoot string) {
    // Limpiar procesos anteriores si existen
    killProcess(clientCmd)
    killProcess(enaireCmd)
    
    // Iniciar cliente
    clientPath := filepath.Join(projectRoot, "cmd", "cliente", "cliente.go")
    clientCmd = exec.Command("go", "run", clientPath)
    clientCmd.Env = append(os.Environ(),
        fmt.Sprintf("PLANES_A=%d", numA),
        fmt.Sprintf("PLANES_B=%d", numB),
        fmt.Sprintf("PLANES_C=%d", numC))
    clientCmd.Stdout = os.Stdout
    clientCmd.Stderr = os.Stderr
    
    if err := clientCmd.Start(); err != nil {
        t.Fatalf("Error iniciando cliente: %v", err)
    }

    time.Sleep(2 * time.Second)

    // Iniciar ENAIRE
    enairePath := filepath.Join(projectRoot, "cmd", "enaire", "enaire.go")
    enaireCmd = exec.Command("go", "run", enairePath)
    enaireCmd.Stdout = os.Stdout
    enaireCmd.Stderr = os.Stderr
    
    if err := enaireCmd.Start(); err != nil {
        t.Fatalf("Error iniciando ENAIRE: %v", err)
    }

    // Tiempo de espera basado en el número total de aviones
    totalPlanes := numA + numB + numC
    waitTime := time.Duration(totalPlanes) * 3 * time.Second
    t.Logf("Esperando %v para procesar %d aviones", waitTime, totalPlanes)
    time.Sleep(waitTime)

    // Verificar conexiones
    verifyConnections(t)
}

func verifyConnections(t *testing.T) {
    conn, err := net.Dial("tcp", "localhost:8000")
    if err != nil {
        t.Errorf("Error verificando conexión: %v", err)
        return
    }
    conn.Close()
    t.Log("Conexión verificada correctamente")
}

func TestMain(m *testing.M) {
    // Ejecutar las pruebas
    code := m.Run()
    
    // Asegurar que todos los procesos se cierren
    cleanupProcesses()
    
    // Esperar un poco más para asegurar que los procesos se han cerrado
    time.Sleep(3 * time.Second)
    
    os.Exit(code)
}