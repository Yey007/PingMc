package data

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"sync"

	"yey007.github.io/software/pingmc/networking"

	"gorm.io/gorm"
)

type PingRepo struct {
	db            *gorm.DB
	tempStore     map[uint]TempPingData
	tempStoreLock sync.RWMutex
}

type PingPair struct {
	Ping RecurringPing
	Temp TempPingData
}

func NewPingRepo(config Config) (*PingRepo, error) {
	var repo PingRepo
	var err error

	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
	)
	repo.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

func (p *PingRepo) GetAll(ctx context.Context, ch chan PingPair) error {
	defer close(ch)
	rows, err := p.db.WithContext(ctx).Joins("Server").Model(&RecurringPing{}).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var ping RecurringPing
		err = p.db.ScanRows(rows, &ping)
		if err != nil {
			return err
		}

		p.tempStoreLock.RLock()
		temp := p.tempStore[ping.ID]
		p.tempStoreLock.RUnlock()

		ch <- PingPair{
			Ping: ping,
			Temp: temp,
		}
	}
	return nil
}

func (p *PingRepo) GetByNameGuildID(ctx context.Context, name string, guildID string) (RecurringPing, error) {
	var ping RecurringPing
	err := p.db.WithContext(ctx).Preload("Server").Where("name = ? AND guild_id = ?", name, guildID).Take(&ping).Error
	return ping, err
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
