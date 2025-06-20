// Copyright 2025 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package service

import (
	"bytes"
	"context"
	"fmt"

	"github.com/google/syzkaller/syz-cluster/pkg/api"
	"github.com/google/syzkaller/syz-cluster/pkg/app"
	"github.com/google/syzkaller/syz-cluster/pkg/blob"
	"github.com/google/syzkaller/syz-cluster/pkg/db"
	"github.com/google/uuid"
)

type FindingService struct {
	findingRepo *db.FindingRepository
	blobStorage blob.Storage
}

func NewFindingService(env *app.AppEnvironment) *FindingService {
	return &FindingService{
		findingRepo: db.NewFindingRepository(env.Spanner),
		blobStorage: env.BlobStorage,
	}
}

func (s *FindingService) Save(ctx context.Context, req *api.NewFinding) error {
	finding := &db.Finding{
		ID:        uuid.NewString(),
		SessionID: req.SessionID,
		TestName:  req.TestName,
		Title:     req.Title,
	}
	var err error
	if len(req.Log) > 0 {
		finding.LogURI, err = s.blobStorage.Write(bytes.NewReader(req.Log), "Finding", finding.ID, "log")
		if err != nil {
			return fmt.Errorf("failed to save the log: %w", err)
		}
	}
	if len(req.Report) > 0 {
		finding.ReportURI, err = s.blobStorage.Write(bytes.NewReader(req.Report), "Finding", finding.ID, "report")
		if err != nil {
			return fmt.Errorf("failed to save the report: %w", err)
		}
	}
	// TODO: if it's not actually addded, the blob records will be orphaned.
	err = s.findingRepo.Save(ctx, finding)
	if err == db.ErrFindingExists {
		// It's ok, just ignore.
		return nil
	}
	return err
}

func (s *FindingService) List(ctx context.Context, sessionID string) ([]*api.Finding, error) {
	list, err := s.findingRepo.ListForSession(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query the list: %w", err)
	}
	var ret []*api.Finding
	for _, item := range list {
		finding := &api.Finding{
			Title:  item.Title,
			LogURL: "TODO", // TODO: where to take it from?
		}
		bytes, err := blob.ReadAllBytes(s.blobStorage, item.ReportURI)
		if err != nil {
			return nil, fmt.Errorf("failed to read the report: %w", err)
		}
		finding.Report = string(bytes)
		ret = append(ret, finding)
	}
	return ret, nil
}
