package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/fatkhur1960/goauction/app"
	"github.com/fatkhur1960/goauction/app/middleware"
	"github.com/fatkhur1960/goauction/app/models"
	"github.com/fatkhur1960/goauction/app/repository"
	"github.com/fatkhur1960/goauction/system/core"
)

// ProductMonitor --
type ProductMonitor struct {
	repo models.ProductQuerySet
}

func (p *ProductMonitor) inspectProduct() error {
	products := []models.Product{}

	log.Println("ProductMonitor] inspecting products...")

	now := time.Now().UTC()
	log.Println("ProductMonitor] inspectProduct now:", now)

	if err := p.repo.ClosedEq(false).ClosedATLte(now).All(&products); err != nil {
		return err
	}

	return p.processCloseProduct(products)
}

func (p *ProductMonitor) processCloseProduct(products []models.Product) error {
	for _, product := range products {
		log.Printf("ProductMonitor] Closing product with name: `%s`", product.ProductName)
		p.createNotifs(&product)
		err := p.repo.IDEq(product.ID).GetUpdater().SetClosed(true).Update()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ProductMonitor) createNotifs(product *models.Product) error {
	notifRepo := repository.NewNotifRepository()
	userRepo := repository.NewUserRepository()
	bidStatus := product.GetBidderStatus(middleware.CurrentUser.ID)
	uid := bidStatus.LatestUserID
	price := bidStatus.LatestBidPrice

	if bidStatus.BidCount != 0 {
		// create notif for product creator
		user, _ := userRepo.GetByID(uid)
		title := fmt.Sprintf("%s ditutup", product.ProductName)
		content := fmt.Sprintf("%s memenangkan bid Anda dengan harga %v", user.FullName, price)
		notifRepo.CreateNotif(product.UserID, title, content, core.GotWinner, product.ID)

		// create notif for bidder
		title = fmt.Sprintf("Selamat Anda menangkan bid untuk %s", product.ProductName)
		content = fmt.Sprintf("Anda memenangkan bid dengan harga %v", price)
		notifRepo.CreateNotif(product.UserID, title, content, core.WinBid, product.ID)
	} else {
		// create notif for product creator
		title := fmt.Sprintf("%s ditutup", product.ProductName)
		content := "Belum ada pemenang untuk bid ini"
		notifRepo.CreateNotif(product.UserID, title, content, core.GotWinner, product.ID)
	}

	return nil
}

// Start --
func (p *ProductMonitor) Start() {
	for {
		log.Println("ProductMonitor] monitor checking...")
		if err := p.inspectProduct(); err != nil {
			log.Fatalf("check product got error: %s\n", err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}

// Stop --
func (p *ProductMonitor) Stop() {}

// NewProductMonitor instance
func NewProductMonitor() Monitor {
	return &ProductMonitor{
		repo: models.NewProductQuerySet(app.DB),
	}
}
