package mgrModel

import (
	"encoding/json"
	"time"
)

// video clear level config of encoder
type ConfigVideoQuality struct {
	ID                int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`                       // autocreated key
	TID               int64     `gorm:"column:tid;default:0" json:"tid"`                                      // tenant id
	AppID             int64     `gorm:"column:app_id;default:0" json:"app_id"`                                // game id, configs focus the game when the id >0
	ClientType        int       `gorm:"column:client_type;default:0;NOT NULL" json:"client_type"`             //
	BizType           int       `gorm:"column:biz_type;default:0;NOT NULL" json:"biz_type"`                   // Business type 0:cloud game 1:cloud phone
	Fps               int       `gorm:"column:fps;default:0;NOT NULL" json:"fps"`                             // fps
	Name              string    `gorm:"column:name;NOT NULL" json:"name"`                                     // level name
	ResoWidthMax      int       `gorm:"column:reso_width_max;default:0;NOT NULL" json:"reso_width_max"`       // resolution width
	ResoHeightMax     int       `gorm:"column:reso_height_max;default:0;NOT NULL" json:"reso_height_max"`     // video resolution height
	ResoHeightMin     int       `gorm:"column:reso_height_min;NOT NULL" json:"reso_height_min"`               // The minimum resolution at this image quality level is high
	ResoWidthMin      int       `gorm:"column:reso_width_min;NOT NULL" json:"reso_width_min"`                 // The minimum resolution width at this image quality level
	ResoHeightDefault int       `gorm:"column:reso_height_default;NOT NULL" json:"reso_height_default"`       // The default resolution at this image quality level is high
	ResoWidthDefault  int       `gorm:"column:reso_width_default;NOT NULL" json:"reso_width_default"`         // The default resolution width at this image quality level
	ParamByEncType    string    `gorm:"column:param_by_enc_type;NOT NULL" json:"param_by_enc_type"`           // parameter by encode types, json: [{"encType":"h264","qpMin":18,"qpMax":50,"qpDefault":50,"kbpsMin":4000,"kbpsMax":4000,"kbpsDefault":4000,"fpsMax":60,"fpsMin":60,"fpsDefault":60}]
	AsDefault         int       `gorm:"column:as_default;default:0;NOT NULL" json:"as_default"`               // Is it the default image quality level parameter for initial playback
	VgpuCost          int       `gorm:"column:vgpu_cost;default:0;NOT NULL" json:"vgpu_cost"`                 // The resource consumption value of full speed encoding for vgpu at this image quality level
	Note              string    `gorm:"column:note;NOT NULL" json:"note"`                                     // Remark
	CreateAt          time.Time `gorm:"column:create_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_at"` // Data record creation time
	UpdateAt          time.Time `gorm:"column:update_at;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_at"` // Data update time
	Sort              int       `gorm:"column:sort;default:0;NOT NULL" json:"sort"`                           // Sort order, 0 is the minimum, the larger the number, the higher the weight
}

type ParamByEncType struct {
	EncType         EncType `gorm:"column:enc_type" json:"encType"`                 // The encoder type of this image quality level
	QpMin           int32   `gorm:"column:qp_min" json:"qpMin"`                     // The minimum qp of this image quality level
	QpMax           int32   `gorm:"column:qp_max" json:"qpMax"`                     // The maximum qp of this image quality level
	QpDefault       int32   `gorm:"column:qp_default" json:"qpDefault"`             // The default qp of this image quality level
	KbpsMin         int32   `gorm:"column:kbps_min" json:"kbpsMin"`                 // The minimum kbps of this image quality level
	KbpsMax         int32   `gorm:"column:kbps_max" json:"kbpsMax"`                 // The maximum kbps of this image quality level
	KbpsDefault     int32   `gorm:"column:kbps_default" json:"kbpsDefault"`         // The default kbps of this image quality level
	KbpsCoefficient float32 `gorm:"column:kbps_coefficient" json:"kbpsCoefficient"` // The kbps coefficient of this image quality level
}

type EncType string

const (
	EncType264 EncType = "H264"
	EncType265 EncType = "H265"
	EncTypeVP8 EncType = "VP8"
)

// TableName get sql table name.
func (m *ConfigVideoQuality) TableName() string {
	return "config_video_quality"
}

func (m *ConfigVideoQuality) GetParamByEncType() []*ParamByEncType {
	var params []*ParamByEncType
	if m.ParamByEncType != "" {
		err := json.Unmarshal([]byte(m.ParamByEncType), &params)
		if err != nil {
			return params
		}
	}
	return params
}
