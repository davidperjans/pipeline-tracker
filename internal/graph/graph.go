package graph

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var driver neo4j.DriverWithContext

// ConnectToNeo4j initializes the Neo4j driver and verifies connectivity
func ConnectToNeo4j() error {
	var err error
	uri := "bolt://localhost:7687"
	username := "neo4j"
	password := "password" // TODO: Load securely via env var

	driver, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return fmt.Errorf("❌ Failed to create Neo4j driver: %w", err)
	}

	ctx := context.Background()
	if err := driver.VerifyConnectivity(ctx); err != nil {
		return fmt.Errorf("❌ Failed to connect to Neo4j: %w", err)
	}

	fmt.Println("✅ Connected to Neo4j")
	return nil
}

// CreatePipelineNode creates a node for a pipeline run
func CreatePipelineNode(id string, branch string, status string) error {
	ctx := context.Background()
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		// 1. Skapa ny nod
		createQuery := `
			CREATE (n:PipelineRun {id: $id, branch: $branch, status: $status})
		`
		params := map[string]interface{}{
			"id":     id,
			"branch": branch,
			"status": status,
		}
		if _, err := tx.Run(ctx, createQuery, params); err != nil {
			return nil, err
		}

		// 2. Hitta senaste noden (med högsta ID mindre än den nya)
		findPrevQuery := `
			MATCH (prev:PipelineRun)
			WHERE toInteger(prev.id) < toInteger($id)
			RETURN prev
			ORDER BY toInteger(prev.id) DESC
			LIMIT 1
		`
		result, err := tx.Run(ctx, findPrevQuery, map[string]interface{}{"id": id})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			// 3. Skapa relation om föregående nod finns
			createRelQuery := `
				MATCH (prev:PipelineRun {id: $prevId}), (curr:PipelineRun {id: $currId})
				CREATE (prev)-[:PRECEDES]->(curr)
			`
			relParams := map[string]interface{}{
				"prevId": result.Record().Values[0].(neo4j.Node).Props["id"],
				"currId": id,
			}
			if _, err := tx.Run(ctx, createRelQuery, relParams); err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return err
}
