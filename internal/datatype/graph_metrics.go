package datatype

import (
	"fmt"
	"sort"
	"strings"

	"github.com/awalterschulze/gographviz"
)

// CFGMetrics contiene todas las métricas calculadas
type CFGMetrics struct {
	// --- Métricas Básicas ---
	NodeCount int
	EdgeCount int

	// --- Métricas de Complejidad ---
	CyclomaticComplexity int
	GraphDensity         float64

	// --- Métricas de Flujo ---
	MaxOutDegree   int
	MaxInDegree    int
	AvgOutDegree   float64
	AvgInDegree    float64
	BranchingNodes int
	MergeNodes     int
	EntryPoints    int
	ExitPoints     int

	// --- Análisis de Integridad ---
	OrphanNodes            []string
	DisconnectedComponents [][]string
	NumDisconnectedSets    int
	HasCycles              bool

	// --- NUEVO: Coverage y Caminos ---
	DetectedEntryPoint string  // Nombre del nodo considerado "Inicio"
	ReachableNodes     int     // Cantidad de nodos visitables desde el inicio
	CodeCoverage       float64 // % de nodos alcanzables (0.0 a 100.0)
	MaxDepth           int     // Distancia al nodo más lejano (Profundidad del árbol)
}

func (m *CFGMetrics) Headers() []string {
	return []string{
		"Nodes",
		"Edges",
		"McCabe",
		"Density",
		"MaxOut",
		"MaxIn",
		"AvgOut",
		"AvgIn",
		"Branching",
		"Merge",
		"EntryPts",
		"ExitPts",
		"Islands",
		"Cycles",
		"DetectedEntry",
		"Reachable",
		"Coverage %",
		"MaxDepth",
		"Orphans",
	}
}

func (m *CFGMetrics) Rows() []string {
	// Formateo de booleanos
	cycles := "No"
	if m.HasCycles {
		cycles = "YES"
	}

	orphansStr := fmt.Sprintf("%d", len(m.OrphanNodes))

	return []string{
		fmt.Sprintf("%d", m.NodeCount),
		fmt.Sprintf("%d", m.EdgeCount),
		fmt.Sprintf("%d", m.CyclomaticComplexity),
		fmt.Sprintf("%.4f", m.GraphDensity),
		fmt.Sprintf("%d", m.MaxOutDegree),
		fmt.Sprintf("%d", m.MaxInDegree),
		fmt.Sprintf("%.2f", m.AvgOutDegree),
		fmt.Sprintf("%.2f", m.AvgInDegree),
		fmt.Sprintf("%d", m.BranchingNodes),
		fmt.Sprintf("%d", m.MergeNodes),
		fmt.Sprintf("%d", m.EntryPoints),
		fmt.Sprintf("%d", m.ExitPoints),
		fmt.Sprintf("%d", m.NumDisconnectedSets),
		cycles,
		m.DetectedEntryPoint,
		fmt.Sprintf("%d/%d", m.ReachableNodes, m.NodeCount),
		fmt.Sprintf("%.2f%%", m.CodeCoverage),
		fmt.Sprintf("%d", m.MaxDepth),
		orphansStr,
	}
}

// AnalyzeCFG procesa el string DOT y retorna métricas
func AnalyzeCFG(dotData string, toolName string) (*CFGMetrics, error) {
	// 1. Parsear DOT
	parsedGraph, err := gographviz.ParseString(dotData)
	if err != nil {
		return nil, fmt.Errorf("error al parsear DOT: %v", err)
	}
	graph := gographviz.NewGraph()
	if err := gographviz.Analyse(parsedGraph, graph); err != nil {
		return nil, fmt.Errorf("error al analizar grafo: %v", err)
	}

	sanitizeGraph(graph, toolName)

	adj := make(map[string][]string)
	revAdj := make(map[string][]string)
	nodes := make(map[string]bool)

	for _, node := range graph.Nodes.Nodes {
		nodes[node.Name] = true
	}

	edgeCount := 0
	for src, dstMap := range graph.Edges.SrcToDsts {
		for dst, edges := range dstMap {
			count := len(edges)
			edgeCount += count
			nodes[src] = true
			nodes[dst] = true
			for i := 0; i < count; i++ {
				adj[src] = append(adj[src], dst)
				revAdj[dst] = append(revAdj[dst], src)
			}
		}
	}

	nodeCount := len(nodes)
	metrics := &CFGMetrics{
		NodeCount: nodeCount,
		EdgeCount: edgeCount,
	}

	// 3. Métricas estándar
	metrics.CyclomaticComplexity = edgeCount - nodeCount + 2
	if nodeCount > 1 {
		metrics.GraphDensity = float64(edgeCount) / float64(nodeCount*(nodeCount-1))
	}

	var totalOut, totalIn int
	var potentialEntries []string

	for node := range nodes {
		outD := len(adj[node])
		inD := len(revAdj[node])
		totalOut += outD
		totalIn += inD

		if outD > metrics.MaxOutDegree {
			metrics.MaxOutDegree = outD
		}
		if inD > metrics.MaxInDegree {
			metrics.MaxInDegree = inD
		}
		if outD > 1 {
			metrics.BranchingNodes++
		}
		if inD > 1 {
			metrics.MergeNodes++
		}

		if inD == 0 {
			metrics.EntryPoints++
			potentialEntries = append(potentialEntries, node)
		}
		if outD == 0 {
			metrics.ExitPoints++
		}

		if outD == 0 && inD == 0 {
			metrics.OrphanNodes = append(metrics.OrphanNodes, node)
		}
	}

	if nodeCount > 0 {
		metrics.AvgOutDegree = float64(totalOut) / float64(nodeCount)
		metrics.AvgInDegree = float64(totalIn) / float64(nodeCount)
	}

	metrics.HasCycles = detectCycles(nodes, adj)
	metrics.DisconnectedComponents = findWeaklyConnectedComponents(nodes, adj, revAdj)
	// Subtract 1 so that a fully connected graph shows 0 islands
	if len(metrics.DisconnectedComponents) > 1 {
		metrics.NumDisconnectedSets = len(metrics.DisconnectedComponents) - 1
	} else {
		metrics.NumDisconnectedSets = 0
	}
	sort.Strings(metrics.OrphanNodes)

	// --- 4. NUEVO: Cálculo de Coverage y Profundidad ---

	// A. Determinar el Entry Point principal
	metrics.DetectedEntryPoint = determineMainEntry(potentialEntries, nodes)

	// B. Calcular métricas de caminos desde ese Entry Point
	if metrics.DetectedEntryPoint != "" {
		reachable, maxDepth := calculateReachabilityAndDepth(metrics.DetectedEntryPoint, adj)
		metrics.ReachableNodes = reachable
		metrics.MaxDepth = maxDepth
		if nodeCount > 0 {
			metrics.CodeCoverage = (float64(reachable) / float64(nodeCount)) * 100.0
		}
	}

	return metrics, nil
}

// determineMainEntry usa heurística para elegir el nodo inicial
func determineMainEntry(candidates []string, allNodes map[string]bool) string {
	// Estrategia 1: Buscar nombres comunes
	commonNames := []string{"Entry", "entry", "Start", "start", "Main", "main", "Root"}
	for node := range allNodes {
		cleanName := strings.ReplaceAll(node, "\"", "") // Limpiar comillas si las hay
		for _, common := range commonNames {
			if cleanName == common {
				return node
			}
		}
	}

	// Estrategia 2: Usar candidatos con In-Degree 0
	if len(candidates) > 0 {
		sort.Strings(candidates)
		return candidates[0] // Retornar el primero alfabéticamente para determinismo
	}

	// Estrategia 3: Si todo falla (ej. grafo circular puro), retornar cualquiera
	for node := range allNodes {
		return node
	}
	return ""
}

// calculateReachabilityAndDepth realiza un BFS para medir cobertura y profundidad
func calculateReachabilityAndDepth(startNode string, adj map[string][]string) (int, int) {
	visited := make(map[string]bool)
	queue := []string{startNode}

	// Map para rastrear distancia desde el origen
	distance := make(map[string]int)

	visited[startNode] = true
	distance[startNode] = 0

	maxDist := 0
	count := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		count++

		currDist := distance[current]
		if currDist > maxDist {
			maxDist = currDist
		}

		for _, neighbor := range adj[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				distance[neighbor] = currDist + 1
				queue = append(queue, neighbor)
			}
		}
	}

	return count, maxDist
}

// findWeaklyConnectedComponents (Igual que versión anterior)
func findWeaklyConnectedComponents(allNodes map[string]bool, adj, revAdj map[string][]string) [][]string {
	visited := make(map[string]bool)
	var components [][]string

	for node := range allNodes {
		if !visited[node] {
			var component []string
			queue := []string{node}
			visited[node] = true

			for len(queue) > 0 {
				current := queue[0]
				queue = queue[1:]
				component = append(component, current)

				for _, neighbor := range adj[current] {
					if !visited[neighbor] {
						visited[neighbor] = true
						queue = append(queue, neighbor)
					}
				}
				for _, parent := range revAdj[current] {
					if !visited[parent] {
						visited[parent] = true
						queue = append(queue, parent)
					}
				}
			}
			sort.Strings(component)
			components = append(components, component)
		}
	}
	sort.Slice(components, func(i, j int) bool { return len(components[i]) > len(components[j]) })
	return components
}

// detectCycles (Igual que versión anterior)
func detectCycles(allNodes map[string]bool, adj map[string][]string) bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var dfs func(string) bool
	dfs = func(n string) bool {
		visited[n] = true
		recStack[n] = true
		for _, neighbor := range adj[n] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true
			}
		}
		recStack[n] = false
		return false
	}
	for n := range allNodes {
		if !visited[n] {
			if dfs(n) {
				return true
			}
		}
	}
	return false
}

// sanitizeGraph removes nodes that are not part of the CFG logic (legends, metadata, etc.)
func sanitizeGraph(graph *gographviz.Graph, toolName string) {
	nodeName := ""
	switch toolName {
	case "rattle":
		//
		_ = graph.RemoveNode(graph.Name, nodeName)
	case "ethersolve":
		// Ethersolve: "EXIT BLOCK"
		_ = graph.RemoveNode(graph.Name, nodeName)
	case "evm-lisa":
		_ = graph.RemoveSubGraph(graph.Name, "cluster_legend")
	}
}
