package list

import (
	"Tiktok/global"
	"math/rand"

	"gorm.io/gorm"
)

/* ---------------------- 入参结构 ---------------------- */

type PageInfo struct {
	Page  int64 `form:"page"`  // 第几页
	Limit int64 `form:"limit"` // 每页条数
}

func (p *PageInfo) Sanitize() {
	if p.Page <= 0 || p.Page > 20 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 10
	}
}

type Options struct {
	PageInfo
	Where    func(db *gorm.DB) *gorm.DB // 自定义过滤条件
	Preloads []string                   // 需要预加载的字段
	Order    string                     // latest / popular / random
	Debug    bool
}

/* ---------------------- 通用查询 ---------------------- */

// Query 返回分页后的列表 + 总条数
func Query[T any](opt Options) (list []T, total int64, err error) {
	// 1. 参数合法性校验
	opt.Sanitize()

	// 2. 构造初始 DB
	db := global.DB.Model(new(T))
	if opt.Debug {
		db = db.Debug()
	}
	if opt.Where != nil {
		db = opt.Where(db)
	}

	// 3. 先算总数（where 之后）
	if err = db.Count(&total).Error; err != nil {
		return
	}

	// 4. 排序 & 分页
	switch opt.Order {
	case "latest":
		db = db.Order("created_time DESC")
	case "popular":
		db = db.Order("likes DESC, created_time DESC")
	case "random":
		db = randomOrder(db)
	default:
		db = db.Order("RAND()")
	}

	offset := (opt.Page - 1) * opt.Limit
	db = db.Offset(int(offset)).Limit(int(opt.Limit))

	// 5. 预加载
	for _, p := range opt.Preloads {
		db = db.Preload(p)
	}

	// 6. 取数据
	err = db.Find(&list).Error
	return
}

/* ---------------------- 随机排序辅助函数 ---------------------- */

func randomOrder(db *gorm.DB) *gorm.DB {
	var ids []int64
	// 把当前 db 的 select/where 原样抄给子查询
	if err := db.Session(&gorm.Session{}). // 复制一份 session，避免污染
						Select("id").
						Find(&ids).Error; err != nil {
		return db // fallback：如果出错就还是用 RAND()
	}

	if len(ids) == 0 {
		return db.Where("1=0") // 无数据
	}

	rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	return db.Where("id IN ?", ids)
}
