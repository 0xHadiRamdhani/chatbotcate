package services

import (
	"whatsapp-bot/internal/config"
	"whatsapp-bot/internal/models"
	"whatsapp-bot/pkg/whatsapp"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
)

type ServiceManager struct {
	DB                *gorm.DB
	Redis             *redis.Client
	WhatsApp          *whatsapp.Client
	Config            *config.Config
	UserService       *UserService
	ContactService    *ContactService
	MessageService    *MessageService
	AutoReplyService  *AutoReplyService
	BroadcastService  *BroadcastService
	GameService       *GameService
	BusinessService   *BusinessService
	ReminderService   *ReminderService
	ModerationService *ModerationService
	UtilityService    *UtilityService
	AnalyticsService  *AnalyticsService
	CleanupService    *CleanupService
}

func NewServiceManager(db *gorm.DB, redis *redis.Client, waClient *whatsapp.Client, cfg *config.Config) *ServiceManager {
	sm := &ServiceManager{
		DB:       db,
		Redis:    redis,
		WhatsApp: waClient,
		Config:   cfg,
	}

	// Initialize all services
	sm.UserService = NewUserService(sm)
	sm.ContactService = NewContactService(sm)
	sm.MessageService = NewMessageService(sm)
	sm.AutoReplyService = NewAutoReplyService(sm)
	sm.BroadcastService = NewBroadcastService(sm)
	sm.GameService = NewGameService(sm)
	sm.BusinessService = NewBusinessService(sm)
	sm.ReminderService = NewReminderService(sm)
	sm.ModerationService = NewModerationService(sm)
	sm.UtilityService = NewUtilityService(sm)
	sm.AnalyticsService = NewAnalyticsService(sm)
	sm.CleanupService = NewCleanupService(sm)

	return sm
}

type UserService struct {
	sm *ServiceManager
}

func NewUserService(sm *ServiceManager) *UserService {
	return &UserService{sm: sm}
}

type ContactService struct {
	sm *ServiceManager
}

func NewContactService(sm *ServiceManager) *ContactService {
	return &ContactService{sm: sm}
}

type MessageService struct {
	sm *ServiceManager
}

func NewMessageService(sm *ServiceManager) *MessageService {
	return &MessageService{sm: sm}
}

type AutoReplyService struct {
	sm *ServiceManager
}

func NewAutoReplyService(sm *ServiceManager) *AutoReplyService {
	return &AutoReplyService{sm: sm}
}

type BroadcastService struct {
	sm *ServiceManager
}

func NewBroadcastService(sm *ServiceManager) *BroadcastService {
	return &BroadcastService{sm: sm}
}

type GameService struct {
	sm *ServiceManager
}

func NewGameService(sm *ServiceManager) *GameService {
	return &GameService{sm: sm}
}

type BusinessService struct {
	sm *ServiceManager
}

func NewBusinessService(sm *ServiceManager) *BusinessService {
	return &BusinessService{sm: sm}
}

type ReminderService struct {
	sm *ServiceManager
}

func NewReminderService(sm *ServiceManager) *ReminderService {
	return &ReminderService{sm: sm}
}

type ModerationService struct {
	sm *ServiceManager
}

func NewModerationService(sm *ServiceManager) *ModerationService {
	return &ModerationService{sm: sm}
}

type UtilityService struct {
	sm *ServiceManager
}

func NewUtilityService(sm *ServiceManager) *UtilityService {
	return &UtilityService{sm: sm}
}

type AnalyticsService struct {
	sm *ServiceManager
}

func NewAnalyticsService(sm *ServiceManager) *AnalyticsService {
	return &AnalyticsService{sm: sm}
}

type CleanupService struct {
	sm *ServiceManager
}

func NewCleanupService(sm *ServiceManager) *CleanupService {
	return &CleanupService{sm: sm}
}