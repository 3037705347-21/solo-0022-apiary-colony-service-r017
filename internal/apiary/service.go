package apiary

import (
	"time"

	"example.com/apiary-colony-service/internal/domain"
	"example.com/apiary-colony-service/internal/store"
)

// Service 负责蜂场与蜂箱的录入，以及蜂箱状态迁移。
type Service struct {
	store *store.Memory
	clock func() time.Time
}

func New(memory *store.Memory, clock func() time.Time) *Service {
	if clock == nil {
		clock = time.Now
	}
	return &Service{store: memory, clock: clock}
}

type RegisterApiaryInput struct {
	ID       string
	Name     string
	District string
	Keeper   string
}

type InstallHiveInput struct {
	ID    string
	Model string
}

func (s *Service) RegisterApiary(input RegisterApiaryInput) (domain.Apiary, error) {
	// 先完成全部领域校验：校验失败的登记直接返回，不留下任何占位痕迹。
	apiary, err := domain.NewApiary(input.ID, input.Name, input.District, input.Keeper, s.clock())
	if err != nil {
		return domain.Apiary{}, err
	}
	// 唯一性判定与写入在仓储的同一把写锁内完成，重复或并发的同编号提交
	// 至多一次成功，且已经确认的内容不会被后续提交改写。
	return s.store.RegisterApiary(apiary)
}

func (s *Service) Apiary(id string) (domain.Apiary, error) {
	apiary, ok := s.store.GetApiary(id)
	if !ok {
		return domain.Apiary{}, domain.NotFound("蜂场", id)
	}
	return apiary, nil
}

func (s *Service) InstallHive(apiaryID string, input InstallHiveInput) (domain.Hive, error) {
	// 先确认蜂场存在且已经登记成功；被拒绝或从未接受的编号在这里就是未找到。
	if _, err := s.Apiary(apiaryID); err != nil {
		return domain.Hive{}, err
	}
	model, err := ParseHiveModel(input.Model)
	if err != nil {
		return domain.Hive{}, err
	}
	hive, err := domain.NewHive(input.ID, apiaryID, model, s.clock())
	if err != nil {
		return domain.Hive{}, err
	}
	// 蜂场存在性与蜂箱编号全局唯一在同一把写锁内权威复核并写入，
	// 保证挂在蜂场下的蜂箱都能在蜂场汇总中找到稳定归属。
	if err := s.store.CreateHive(hive); err != nil {
		return domain.Hive{}, err
	}
	return hive, nil
}

func (s *Service) Hive(id string) (domain.Hive, error) {
	hive, ok := s.store.GetHive(id)
	if !ok {
		return domain.Hive{}, domain.NotFound("蜂箱", id)
	}
	return hive, nil
}

func (s *Service) AdvanceHiveState(hiveID, target string) (domain.Transition, error) {
	to, err := ParseHiveState(target)
	if err != nil {
		return domain.Transition{}, err
	}
	return s.store.AdvanceHive(hiveID, to, s.clock())
}
