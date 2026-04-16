package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/activityGuide/model"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/activityGuide/model/request"
)

var Guide = new(guide)

type guide struct{}

func (s *guide) CreateGuide(guide *model.ActivityGuide) (err error) {
	if guide.IsPreview == nil {
		f := false
		guide.IsPreview = &f
	}
	if guide.ShowStatus == nil {
		t := true
		guide.ShowStatus = &t
	}
	err = global.GVA_DB.Create(guide).Error
	return err
}

func (s *guide) DeleteGuide(ID string) (err error) {
	err = global.GVA_DB.Delete(&model.ActivityGuide{}, "id = ?", ID).Error
	return err
}

func (s *guide) DeleteGuideByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]model.ActivityGuide{}, "id in ?", IDs).Error
	return err
}

func (s *guide) UpdateGuide(guide model.ActivityGuide) (err error) {
	// Select 指定列后，零值也会写入（否则 coverImage 传 "" 无法清空封面）
	err = global.GVA_DB.Model(&model.ActivityGuide{}).Where("id = ?", guide.ID).
		Select("name", "summary", "cover_image", "media", "is_preview", "show_status").
		Updates(&guide).Error
	return err
}

func (s *guide) GetGuide(ID string) (guide model.ActivityGuide, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&guide).Error
	return
}

func (s *guide) GetGuideList(search request.GuideSearch) (list []model.ActivityGuide, total int64, err error) {
	limit := search.PageSize
	offset := search.PageSize * (search.Page - 1)
	db := global.GVA_DB.Model(&model.ActivityGuide{})
	var guides []model.ActivityGuide

	if search.Name != "" {
		db = db.Where("name LIKE ?", "%"+search.Name+"%")
	}
	if search.ID != nil {
		db = db.Where("id = ?", *search.ID)
	}
	if search.IsPreview != nil {
		db = db.Where("is_preview = ?", *search.IsPreview)
	}
	if search.ShowStatus != nil {
		db = db.Where("show_status = ?", *search.ShowStatus)
	}
	if search.StartCreatedAt != nil && search.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", search.StartCreatedAt, search.EndCreatedAt)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("created_at DESC").Find(&guides).Error
	return guides, total, err
}
