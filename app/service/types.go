package service

import (
	"time"

	"github.com/fatih/structs"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type (
	// Product definisi type untuk menampilkan product
	Product struct {
		ID            int64      `json:"id"`
		ProductName   string     `json:"product_name"`
		ProductImages []string   `json:"product_images"`
		Desc          string     `json:"desc"`
		Condition     int32      `json:"condition" `
		ConditionAvg  float64    `json:"condition_avg"`
		StartPrice    float64    `json:"start_price"`
		BidMultpl     float64    `json:"bid_multpl"`
		ClosedAT      *time.Time `json:"closed_at"`
		CreatedAT     *time.Time `json:"created_at"`
		Labels        []string   `json:"labels"`
	}
)

func (p *Product) toAPI(model interface{}, db *gorm.DB) Product {
	productMap := structs.Map(model)
	mapstructure.Decode(productMap, &p)

	imageList := []string{}
	images := []models.ProductImage{}
	db.Where("product_id = ?", p.ID).Find(&images)
	for _, img := range images {
		imageList = append(imageList, img.ImageURL)
	}

	labelList := []string{}
	labels := []models.ProductLabel{}
	db.Where("product_id = ?", p.ID).Find(&labels)
	for _, label := range labels {
		labelList = append(labelList, label.Name)
	}

	res := Product{
		ID:            p.ID,
		ProductName:   p.ProductName,
		ProductImages: imageList,
		Desc:          p.Desc,
		Condition:     p.Condition,
		ConditionAvg:  p.ConditionAvg,
		StartPrice:    p.StartPrice,
		BidMultpl:     p.BidMultpl,
		ClosedAT:      p.ClosedAT,
		CreatedAT:     p.CreatedAT,
		Labels:        labelList,
	}

	return res
}
