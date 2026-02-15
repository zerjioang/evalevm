package datatype

import (
	"fmt"
	"testing"
)

func TestAnalyzeCFG(t *testing.T) {
	// Grafo de ejemplo con:
	// 1. Un flujo principal
	// 2. Ramas (if)
	// 3. Código inalcanzable (Dead Code)
	// 4. Entry Point explícito
	dotString := `digraph G {
		Entry -> CheckAuth;
		
		CheckAuth -> Dashboard [label="auth_ok"];
		CheckAuth -> ErrorPage [label="auth_fail"];
		
		Dashboard -> ProcessData;
		ProcessData -> SaveDB;
		SaveDB -> Exit;
		
		ErrorPage -> Exit;

		// Código Muerto (Hacker_Function nunca es llamada desde Entry)
		Hacker_Function -> StealData;
		StealData -> SendToServer;
	}`

	metrics, err := AnalyzeCFG(dotString)
	if err != nil {
		panic(err)
	}

	fmt.Println("=== Reporte de Análisis de CFG ===")
	fmt.Printf("Nodos: %d | Aristas: %d\n", metrics.NodeCount, metrics.EdgeCount)
	fmt.Printf("Complejidad Ciclomática: %d\n", metrics.CyclomaticComplexity)

	fmt.Println("\n--- Análisis de Cobertura (Coverage) ---")
	fmt.Printf("Entry Point Detectado: [%s]\n", metrics.DetectedEntryPoint)
	fmt.Printf("Nodos Alcanzables: %d de %d\n", metrics.ReachableNodes, metrics.NodeCount)
	fmt.Printf("Coverage: %.2f%%\n", metrics.CodeCoverage)
	fmt.Printf("Profundidad Máxima (Desde Entry): %d saltos\n", metrics.MaxDepth)

	if metrics.CodeCoverage < 100 {
		fmt.Println("⚠️  ADVERTENCIA: Se detectó código inalcanzable (Dead Code)")
	}

	fmt.Println("\n--- Integridad Estructural ---")
	fmt.Printf("Islas Desconectadas: %d\n", metrics.NumDisconnectedSets)
	for i, comp := range metrics.DisconnectedComponents {
		status := "✅ Activo"
		// Si el componente no contiene el EntryPoint, es código muerto
		containsEntry := false
		for _, n := range comp {
			if n == metrics.DetectedEntryPoint {
				containsEntry = true
				break
			}
		}
		if !containsEntry {
			status = "💀 Muerto / Inalcanzable"
		}

		fmt.Printf("  Grupo %d (Size: %d): %v -> %s\n", i+1, len(comp), comp, status)
	}
}
