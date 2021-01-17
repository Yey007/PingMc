package data

import (
	"context"
	"sync"

	"gorm.io/driver/sqlite"
	"yey007.github.io/software/pingmc/networking"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

type PingRepo struct {
	db            *gorm.DB
	tempStore     map[uint]TempPingData
	tempStoreLock sync.RWMutex
}

type RecurringPing struct {
	gorm.Model
	Name      string
	GuildID   string
	ChannelID string
	Server    networking.McServer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TempPingData struct {
	RecurringPingID uint
	PreviousData    *networking.PingData
	HasPrevious     bool
}

func NewPingRepo() (*PingRepo, error) {
	var repo PingRepo
	var err error

	repo.db, err = gorm.Open(sqlite.Open("data/database.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = repo.db.AutoMigrate(&RecurringPing{})
	if err != nil {
		return nil, err
	}

	err = repo.db.AutoMigrate(&networking.McServer{})
	if err != nil {
		return nil, err
	}

	repo.tempStore = make(map[uint]TempPingData)

	return &repo, nil
}

func (p *PingRepo) Create(ctx context.Context, ping RecurringPing) error {
	return p.db.WithContext(ctx).Create(&ping).Error
}

func (p *PingRepo) CreateTemp(tempData TempPingData) {
	p.tempStoreLock.Lock()
	p.tempStore[tempData.RecurringPingID] = tempData
	p.tempStoreLock.Unlock()
}

func (p *PingRepo) GetAllInBatches(
	ctx context.Context,
	callback func(ping RecurringPing, tempData TempPingData),
) error {
	var pings []RecurringPing
	return p.db.WithContext(ctx).Preload("Server").FindInBatches(&pings, 100, func(tx *gorm.DB, batch int) error {
		for _, ping := range pings {
			p.tempStoreLock.RLock()
			tempData := p.tempStore[ping.ID]
			p.tempStoreLock.RUnlock()
			callback(ping, tempData)
		}
		return nil
	}).Error
}

func (p *PingRepo) GetByNameGuildID(ctx context.Context, name string, guildID string) (*RecurringPing, error) {
	var ping RecurringPing
	err := p.db.WithContext(ctx).Preload("Server").Where("name = ? AND guild_id = ?", name, guildID).Take(&ping).Error
	return &ping, err
}

func (p *PingRepo) GetByGuildID(ctx context.Context, guildID string) ([]RecurringPing, error) {
	var pings []RecurringPing
	err := p.db.WithContext(ctx).Preload("Server").Where("guild_id = ?", guildID).Find(&pings).Error
	return pings, err
}

func (p *PingRepo) Remove(ctx context.Context, id uint) error {
	err := p.db.WithContext(ctx).Delete(&RecurringPing{}, id).Error
	if err != nil {
		return err
	}
	err = p.db.WithContext(ctx).Where("recurring_ping_id = ?", id).Delete(&networking.McServer{}).Error
	if err != nil {
		return err
	}
	p.tempStoreLock.Lock()
	delete(p.tempStore, id)
	p.tempStoreLock.Unlock()
	return nil
}
