package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/db"
	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/pb"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertPBToDBRangeStatus converts pb.RangeStatus to db.RangeStatus.
func ConvertPBToDBRangeStatus(status *pb.RangeStatus) (db.RangeStatus, error) {
	if status == nil {
		return db.RangeStatusUNSPECIFIED, errors.New("pb.RangeStatus is nil")
	}

	switch *status {
	case pb.RangeStatus_RANGE_STATUS_UNSPECIFIED:
		return db.RangeStatusUNSPECIFIED, nil
	case pb.RangeStatus_RANGE_STATUS_ACTIVE:
		return db.RangeStatusACTIVE, nil
	case pb.RangeStatus_RANGE_STATUS_EXHAUSTED:
		return db.RangeStatusEXHAUSTED, nil
	case pb.RangeStatus_RANGE_STATUS_RELEASED:
		return db.RangeStatusRELEASED, nil
	default:
		return db.RangeStatusRELEASED, errors.New("unknown pb.RangeStatus value")
	}
}

// ConvertDBToPBRangeStatus converts db.RangeStatus to pb.RangeStatus.
func ConvertDBToPBRangeStatus(status *db.RangeStatus) (pb.RangeStatus, error) {
	if status == nil {
		return pb.RangeStatus_RANGE_STATUS_UNSPECIFIED, errors.New("db.RangeStatus is nil")
	}

	switch *status {
	case db.RangeStatusUNSPECIFIED:
		return pb.RangeStatus_RANGE_STATUS_UNSPECIFIED, nil
	case db.RangeStatusACTIVE:
		return pb.RangeStatus_RANGE_STATUS_ACTIVE, nil
	case db.RangeStatusEXHAUSTED:
		return pb.RangeStatus_RANGE_STATUS_EXHAUSTED, nil
	case db.RangeStatusRELEASED:
		return pb.RangeStatus_RANGE_STATUS_RELEASED, nil
	default:
		return pb.RangeStatus_RANGE_STATUS_UNSPECIFIED, errors.New("unknown db.RangeStatus value")
	}
}

func ConvertPgTextToString(text *pgtype.Text) (string, error) {
	if text.Valid {
		return text.String, nil
	}
	return "", fmt.Errorf("text is null or invalid")
}

func ConvertPgTimeToTime(tz *pgtype.Timestamptz) (time.Time, error) {
	if tz.Valid {
		return tz.Time, nil
	}
	return time.Time{}, fmt.Errorf("timestamp is null")
}

func ConvertRangeToProto(r *db.Range) (*pb.Range, error) {
	rangeStatus, err := ConvertDBToPBRangeStatus(&r.Status)
	if err != nil {
		return nil, fmt.Errorf("invalid RangeStatus value: %v", err)
	}

	region, err := ConvertPgTextToString(&r.Region)
	if err != nil {
		return nil, fmt.Errorf("invalid Region value: %v", err)
	}

	allocatedAt, err := ConvertPgTimeToTime(&r.AllocatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid AllocatedAt value: %v", err)
	}

	updatedAt, err := ConvertPgTimeToTime(&r.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid UpdatedAt value: %v", err)
	}

	return &pb.Range{
		RangeId:     r.RangeID.String(),
		StartId:     r.StartID,
		EndId:       r.EndID,
		ServiceId:   r.ServiceID,
		Region:      region,
		Status:      rangeStatus,
		AllocatedAt: timestamppb.New(allocatedAt),
		UpdatedAt:   timestamppb.New(updatedAt),
	}, nil
}
