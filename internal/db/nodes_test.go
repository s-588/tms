package db_test

import (
	"context"
	"slices"
	"testing"

	"github.com/s-588/tms/cmd/models"
	"github.com/s-588/tms/internal/db"
	"github.com/s-588/tms/internal/testutil"
	"github.com/s-588/tms/internal/ui"
)

func TestDB_Nodes(t *testing.T) {
	database := testutil.SetupTestDB(t)

	testCreateNode(t, database)
	testGetNodeByID(t, database)
	testGetNodes(t, database)
	testUpdateNode(t, database)
	testSoftDeleteNode(t, database)
	testHardDeleteNode(t, database)
	testRestoreNode(t, database)
	testBulkSoftDeleteNodes(t, database)
	testBulkHardDeleteNodes(t, database)
	testListNodes(t, database)
	testCalculateDistance(t, database)
}

func testCreateNode(t *testing.T, database db.DB) {
	ctx := context.Background()

	t.Run("creates node successfully", func(t *testing.T) {
		args := db.CreateNodeArgs{
			Name:    models.Optional[string]{Value: "Node A", Set: true},
			Address: "Address A",
			Geom:    models.Point{X: 10, Y: 20},
		}
		node, err := database.CreateNode(ctx, args)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if node.NodeID == 0 {
			t.Errorf("expected node ID to be populated")
		}
	})
}

func testGetNodeByID(t *testing.T, database db.DB) {
	ctx := context.Background()
	node := createTestNode(t, database, "Get Node", 10, 20)

	t.Run("get node by id", func(t *testing.T) {
		got, err := database.GetNodeByID(ctx, node.NodeID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.NodeID != node.NodeID {
			t.Fatalf("got node ID %d, want %d", got.NodeID, node.NodeID)
		}
	})
}

func testGetNodes(t *testing.T, database db.DB) {
	ctx := context.Background()
	node := createTestNode(t, database, "Get All Nodes", 10, 20)

	t.Run("get all nodes", func(t *testing.T) {
		list, total, err := database.GetNodes(ctx, 1, models.NodeFilter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if total == 0 {
			t.Fatalf("expected total pages to be populated")
		}
		if !slices.ContainsFunc(list, func(v models.Node) bool {
			return v.NodeID == node.NodeID
		}) {
			t.Fatalf("slice doesn't contain wanted node; got %v, want %v", list, node)
		}
	})
}

func testUpdateNode(t *testing.T, database db.DB) {
	ctx := context.Background()
	node := createTestNode(t, database, "Update Node", 10, 20)

	t.Run("update node", func(t *testing.T) {
		err := database.UpdateNode(ctx, db.UpdateNodeArgs{
			NodeID:  node.NodeID,
			Name:    models.Optional[string]{Value: "Updated Node", Set: true},
			Address: "Updated Address",
			Geom:    models.Point{X: 11, Y: 22},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		got, err := database.GetNodeByID(ctx, node.NodeID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.Name != "Updated Node" {
			t.Fatalf("expected updated name, got %s", got.Name)
		}
	})
}

func testSoftDeleteNode(t *testing.T, database db.DB) {
	ctx := context.Background()
	node := createTestNode(t, database, "Soft Node", 10, 20)

	t.Run("soft delete node", func(t *testing.T) {
		if err := database.SoftDeleteNode(ctx, node.NodeID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
	})
}

func testHardDeleteNode(t *testing.T, database db.DB) {
	ctx := context.Background()
	node := createTestNode(t, database, "Hard Node", 10, 20)

	t.Run("hard delete node", func(t *testing.T) {
		if err := database.HardDeleteNode(ctx, node.NodeID); err != nil {
			t.Fatalf("hard delete: %v", err)
		}
		if _, err := database.GetNodeByID(ctx, node.NodeID); err == nil {
			t.Fatalf("expected error after hard delete")
		}
	})
}

func testRestoreNode(t *testing.T, database db.DB) {
	ctx := context.Background()
	node := createTestNode(t, database, "Restore Node", 10, 20)

	t.Run("restore node", func(t *testing.T) {
		if err := database.SoftDeleteNode(ctx, node.NodeID); err != nil {
			t.Fatalf("soft delete: %v", err)
		}
		if err := database.RestoreNode(ctx, node.NodeID); err != nil {
			t.Fatalf("restore: %v", err)
		}
		if _, err := database.GetNodeByID(ctx, node.NodeID); err != nil {
			t.Fatalf("get after restore: %v", err)
		}
	})
}

func testBulkSoftDeleteNodes(t *testing.T, database db.DB) {
	ctx := context.Background()
	node := createTestNode(t, database, "Bulk Soft Node", 10, 20)

	t.Run("bulk soft delete nodes", func(t *testing.T) {
		if err := database.BulkSoftDeleteNodes(ctx, []int32{node.NodeID}); err != nil {
			t.Fatalf("bulk soft delete: %v", err)
		}
	})
}

func testBulkHardDeleteNodes(t *testing.T, database db.DB) {
	ctx := context.Background()
	node := createTestNode(t, database, "Bulk Hard Node", 10, 20)

	t.Run("bulk hard delete nodes", func(t *testing.T) {
		if err := database.BulkHardDeleteNodes(ctx, []int32{node.NodeID}); err != nil {
			t.Fatalf("bulk hard delete: %v", err)
		}
		if _, err := database.GetNodeByID(ctx, node.NodeID); err == nil {
			t.Fatalf("expected error after bulk hard delete")
		}
	})
}

func testListNodes(t *testing.T, database db.DB) {
	ctx := context.Background()
	node := createTestNode(t, database, "List Node", 10, 20)

	t.Run("list nodes", func(t *testing.T) {
		items, err := database.ListNodes(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !slices.ContainsFunc(items, func(item ui.ListItem) bool {
			return item.ID == node.NodeID
		}) {
			t.Fatalf("list doesn't contain node; got %v", items)
		}
	})
}

func testCalculateDistance(t *testing.T, database db.DB) {
	ctx := context.Background()
	nodeA := createTestNode(t, database, "Dist Node A", 10, 10)
	nodeB := createTestNode(t, database, "Dist Node B", 20, 20)

	t.Run("calculate distance", func(t *testing.T) {
		distance, err := database.CalculateDistance(ctx, nodeA.NodeID, nodeB.NodeID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if distance < 0 {
			t.Fatalf("expected non-negative distance, got %f", distance)
		}
	})
}

func createTestNode(t *testing.T, database db.DB, name string, x, y float64) models.Node {
	t.Helper()
	ctx := context.Background()

	node, err := database.CreateNode(ctx, db.CreateNodeArgs{
		Name:    models.Optional[string]{Value: name, Set: true},
		Address: name,
		Geom:    models.Point{X: x, Y: y},
	})
	if err != nil {
		t.Fatalf("create node %q: %v", name, err)
	}
	return node
}
