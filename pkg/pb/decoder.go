// Package pb package pb
package pb

import (
	"trade_agent/pkg/taerror"

	"google.golang.org/protobuf/proto"
)

// UnmarshalProto UnmarshalProto
func (c *EventResponse) UnmarshalProto(message []byte) error {
	if err := proto.Unmarshal(message, c); err != nil {
		return taerror.ErrProtoFormatWrong(message, err)
	}
	return nil
}

// UnmarshalProto UnmarshalProto
func (c *HistoryCloseResponse) UnmarshalProto(message []byte) error {
	if err := proto.Unmarshal(message, c); err != nil {
		return taerror.ErrProtoFormatWrong(message, err)
	}
	return nil
}

// UnmarshalProto UnmarshalProto
func (c *OrderStatusHistoryResponse) UnmarshalProto(message []byte) error {
	if err := proto.Unmarshal(message, c); err != nil {
		return taerror.ErrProtoFormatWrong(message, err)
	}
	return nil
}

// UnmarshalProto UnmarshalProto
func (c *StockDetailResponse) UnmarshalProto(message []byte) error {
	if err := proto.Unmarshal(message, c); err != nil {
		return taerror.ErrProtoFormatWrong(message, err)
	}
	return nil
}

// UnmarshalProto UnmarshalProto
func (c *SnapshotResponse) UnmarshalProto(message []byte) error {
	if err := proto.Unmarshal(message, c); err != nil {
		return taerror.ErrProtoFormatWrong(message, err)
	}
	return nil
}

// UnmarshalProto UnmarshalProto
func (c *HistoryTickResponse) UnmarshalProto(message []byte) error {
	if err := proto.Unmarshal(message, c); err != nil {
		return taerror.ErrProtoFormatWrong(message, err)
	}
	return nil
}

// UnmarshalProto UnmarshalProto
func (c *RealTimeTickResponse) UnmarshalProto(message []byte) error {
	if err := proto.Unmarshal(message, c); err != nil {
		return taerror.ErrProtoFormatWrong(message, err)
	}
	return nil
}

// UnmarshalProto UnmarshalProto
func (c *HistoryKbarResponse) UnmarshalProto(message []byte) error {
	if err := proto.Unmarshal(message, c); err != nil {
		return taerror.ErrProtoFormatWrong(message, err)
	}
	return nil
}

// UnmarshalProto UnmarshalProto
func (c *RealTimeBidAskResponse) UnmarshalProto(message []byte) error {
	if err := proto.Unmarshal(message, c); err != nil {
		return taerror.ErrProtoFormatWrong(message, err)
	}
	return nil
}

// UnmarshalProto UnmarshalProto
func (c *VolumeRankResponse) UnmarshalProto(message []byte) error {
	if err := proto.Unmarshal(message, c); err != nil {
		return taerror.ErrProtoFormatWrong(message, err)
	}
	return nil
}
