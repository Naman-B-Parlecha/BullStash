package backup

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
)

type WALManager struct {
	connConfig    *pgx.ConnConfig
	walArchiveDir string
	slotName      string
}

func NewWALManager(connConfig *pgx.ConnConfig, walArchiveDir, slotName string) *WALManager {
	return &WALManager{
		connConfig:    connConfig,
		walArchiveDir: walArchiveDir,
		slotName:      slotName,
	}
}
func (wm *WALManager) Start(ctx context.Context) error {
	conn, err := pgconn.Connect(ctx, wm.connConfig.ConnString())
	if err != nil {
		return fmt.Errorf("failed to connect for replication: %w", err)
	}
	defer conn.Close(ctx)

	// Create replication slot if not exists
	_, err = pglogrepl.CreateReplicationSlot(ctx, conn, wm.slotName, "", pglogrepl.CreateReplicationSlotOptions{
		Mode: pglogrepl.PhysicalReplication,
	})
	if err != nil {
		log.Printf("Replication slot may already exist: %v", err)
	}

	// Start replication
	err = pglogrepl.StartReplication(ctx, conn, wm.slotName, 0, pglogrepl.StartReplicationOptions{})
	if err != nil {
		return fmt.Errorf("failed to start replication: %w", err)
	}

	for {
		msg, err := conn.ReceiveMessage(ctx)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("receive message failed: %w", err)
		}

		switch msg := msg.(type) {
		case *pgproto3.CopyData:
			switch msg.Data[0] {
			case pglogrepl.PrimaryKeepaliveMessageByteID:
				// Handle keepalive
			case pglogrepl.XLogDataByteID:
				xld, err := pglogrepl.ParseXLogData(msg.Data[1:])
				if err != nil {
					return fmt.Errorf("parse xlog data failed: %w", err)
				}
				wm.saveWALSegment(uint64(xld.WALStart), xld.WALData)
			}
		default:
			log.Printf("Received unexpected message: %T", msg)
		}
	}
}

func (wm *WALManager) saveWALSegment(startLSN uint64, data []byte) {
	walPath := filepath.Join(wm.walArchiveDir, fmt.Sprintf("%016X", startLSN))
	if err := os.WriteFile(walPath, data, 0644); err != nil {
		log.Printf("Failed to save WAL segment: %v", err)
	}
}
