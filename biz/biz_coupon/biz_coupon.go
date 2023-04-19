package biz_coupon

import (
	"math"
	"time"
)

type bizCoupon struct{}

var bizCouponVar = bizCoupon{}

func BizCoupon() *bizCoupon {
	return &bizCouponVar
}

// GetCondition 领取条件
type GetCondition struct {
	Method       int `json:"method"`    // 发放方式 0.系统发放 1.主动领取
	Number       int `json:"number"`    // 单用户可领取数量上限 min:1
	PeriodNumber int `json:"condition"` // 周期内可领取数量上限 0.不限数量
	Period       int `json:"period"`    // 领取周期 0.不限周期 other:n天
}

// BizCouponPolicy 优惠券策略
type BizCouponPolicy struct {
	CropMode          int   `json:"cropMode"`          // 组织策略 0:全部子组织通用,1.白名单,2.黑名单
	CropList          []int `json:"cropList"`          // 组织策略 组织黑白名单库
	OverlieMode       int   `json:"overlieMode"`       // 叠加策略 0:不允许与任何优惠券叠加,1:任意叠加（rank值较大且不支持叠加的除外）;2:白名单,3.黑名单
	OverlieList       []int `json:"overlieList"`       // 叠加策略 黑白名单库
	GoodsMode         int   `json:"goodsMode"`         // 商品策略 0:全部可用,1:白名单,2.黑名单
	GoodsList         []int `json:"goodsList"`         // 商品策略 黑白名单库
	GoodsCategoryMode int   `json:"goodsCategoryMode"` // 品类策略 0:全部可用,1:白名单,2.黑名单
	GoodsCategoryList []int `json:"goodsCategoryList"` // 品类策略 黑白名单库
	UserMode          int   `json:"userMode"`          // 用户策略 0:全部可领;1.新用户可领;2.老用户可领,3.白名单;4.黑名单
	UserList          []int `json:"userList"`          // 用户策略 黑白名单库
}
type BizCouponCreate struct {
	Owner             int       `json:"owner"`             // 优惠券归属平台方编码（如非saas模式运行可始终定义为0或按需使用）
	CropId            int       `json:"cropId"`            // 发行方内部子组织编码,0.平台方发行,other:指定子组织发行
	OpId              int       `json:"opId"`              // 发行（操作）者用户凭证
	Type              int       `json:"type"`              // 优惠券类型 0.满减 1.折扣 2.赠品
	Name              string    `json:"name"`              // 优惠券名称 etc. 新人优惠券
	Desc              string    `json:"desc"`              // 优惠券描述 etc. 新人专属优惠券
	Amount            float64   `json:"amount"`            // 优惠券基础面值,对应优惠券类型 ect. ¥100 / 9.9折 / 1件
	StepAmount        float64   `json:"stepAmount"`        // 基于基础面值的优惠券阶梯面值 0.不启用
	StepMaximumAmount float64   `json:"stepMaximumAmount"` // 阶梯优惠最高面值
	StepUnit          float64   `json:"stepUnit"`          // 阶梯优惠步长 etc. 50
	MinimumAmount     float64   `json:"minimumAmount"`     // 最低使用金额门槛
	Number            int64     `json:"number"`            // 发行数量 ect. 10000
	GetNumber         int       `json:"getNumber"`         // 单用户可领数量上限
	Period            int       `json:"period"`            // 领取后有效期（天） etc. 30天
	Disabled          int       `json:"disabled"`          // 是否禁用
	Nbf               time.Time `json:"nbf"`               // 可领取时间 et
	Exp               time.Time `json:"exp"`               // 领取失效时间
	Rank              int       `json:"rank"`              // 权重，用于计算叠加优先级策略
	*BizCouponPolicy            // 优惠券策略
}
type BizCouponRes struct {
	Id int `json:"id"`
}

// Query 查询优惠券
// note:cropId  = -1 返回所有子组织优惠券列表
func (rec *bizCoupon) Query(owner, cropId, categoryId int, policy *BizCouponPolicy) (err error) {
	// code here
	return
}

func (rec *bizCoupon) Create(owner int) (err error) {
	// code here
	rec.Query(1, 1, 1, nil)
	return
}

// CalcStepAmount 计算阶梯优惠金额
func (rec *bizCoupon) CalcStepAmount(inputAmount, minimumAmount, stepAmount, stepMaximumAmount, stepUnit float64) (total float64) {
	if inputAmount < minimumAmount+stepUnit {
		return
	}
	n := math.Floor((inputAmount - minimumAmount) / stepUnit)
	if n < 1 {
		return
	}
	total = n * stepAmount
	if total > stepMaximumAmount {
		total = stepMaximumAmount
	}
	return
}

// inspectOverlie 叠加优惠券检查
func (rec *bizCoupon) inspectOverlie() (total float64) {

	return
}
