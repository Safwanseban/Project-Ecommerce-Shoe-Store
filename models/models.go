package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID             uint   `json:"id" gorm:"primaryKey;unique"  `
	First_Name     string `json:"first_name"  gorm:"not null" validate:"required,min=2,max=50"  `
	Last_Name      string `json:"last_name"    gorm:"not null"    validate:"required,min=1,max=50"  `
	Email          string `json:"email"   gorm:"not null;unique"  validate:"email,required"`
	Password       string `json:"password" gorm:"not null"  validate:"required"`
	Phone          string `json:"phone"   gorm:"not null;unique" validate:"required"`
	Block_status   bool   `json:"block_status " gorm:"not null"   `
	Country        string `json:"country "   `
	City           string `json:"city "   `
	Pincode        uint   `json:"pincode "   `
	Cart           Cart
	Cart_id        uint `json:"cart_id" `
	Address        Address
	Address_id     uint `json:"address_id" `
	Orders         Orders
	Orders_ID      uint `json:"orders_id" `
	Wallet_Balance uint `json:"wallet_balance" `

	Applied_Coupons   Applied_Coupons
	Applied_CouponsID uint
	WishList          WishList
	WishListID        uint
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
type Admin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (admin *Admin) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	admin.Password = string(bytes)
	return nil
}
func (admin *Admin) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

// `json:"cart"  gorm:"foreignKey:User_cart_id;references:Cart_id"`
// `json:"address" gorm:"foreignKey;references:Address_id"`
type Cart struct {
	Cart_id     uint `json:"cart_id" gorm:"primaryKey"  `
	UserId      uint `json:"user_id"   `
	ProductID   uint `json:"product_id"  `
	Quantity    uint `json:"quantity" `
	Total_Price uint `json:"total_price"   `
}
type Cartsinfo struct {
	gorm.Model
	User_id      string
	Product_id   string
	Product_Name string
	Price        string
	Email        string
	Quantity     string
	Total_Price  string
}
type Address struct {
	Address_id   uint   `json:"address_id" gorm:"primaryKey"  `
	UserId       uint   `json:"user_id"  gorm:"not null" `
	Name         string `json:"name"  gorm:"not null" `
	Phone_number int    `json:"phone_number"  gorm:"not null" `
	Pincode      int    `json:"pincode"  gorm:"not null" `
	House        string `json:"house"   `
	Area         string `json:"area"   `
	Landmark     string `json:"landmark"  gorm:"not null" `
	City         string `json:"city"  gorm:"not null" `
}
type Product struct {
	Product_id   uint   `json:"product_id" gorm:"primaryKey" `
	Product_name string `json:"product_name" gorm:"not null"  `
	Price        uint   `json:"price" gorm:"not null"  `
	Actual_Price uint   `json:"actual_price" gorm:"not null"`
	Image        string `json:"image" gorm:"not null"  `
	Cover        string `json:"cover"   `
	SubPic1      string `json:"subpic1"  `
	SubPic2      string `json:"subpic2"  `
	Stock        uint   `json:"stock"  `
	Color        string `json:"color" gorm:"not null"  `
	Description  string `json:"description"   `
	Discount     uint   `json:"discount"`

	Brand      Brand
	Brand_id   uint `json:"brand_id" `
	Cart       Cart
	Cart_id    uint `json:"cart_id" `
	Catogory   Catogory
	CatogoryID uint
	ShoeSize   ShoeSize
	ShoeSizeID uint
	WishList   WishList
	WishListID uint
}
type Brand struct {
	ID       uint   `json:"id" gorm:"primaryKey"  `
	Brands   string `json:"brands" gorm:"not null"  `
	Discount uint   `json:"discount"`
}
type Catogory struct {
	ID       uint `json:"id" gorm:"primaryKey"  `
	Catogory string
}
type ShoeSize struct {
	ID   uint `json:"id" gorm:"primaryKey"  `
	Size uint `json:"size"`
}
type Otp struct {
	gorm.Model
	Mobile string
	Otp    string
}
type PaymentMethod struct {
	COD bool
}
type Orders struct {
	gorm.Model
	UserId          uint   `json:"user_id"  gorm:"not null" `
	Order_id        string `json:"order_id"  gorm:"not null" `
	Total_Amount    uint   `json:"total_amount"  gorm:"not null" `
	Applied_Coupons string `json:"applied_coupons"  `
	Discount        uint   `json:"discount"   `
	PaymentMethod   string `json:"paymentmethod"  gorm:"not null" `
	Payment_Status  string `json:"payment_status"   `
	Order_Status    string `json:"order_status"   `
	Address         Address
	Address_id      uint `json:"address_id"  `
}
type Orderd_Items struct {
	gorm.Model
	UserId          uint `json:"user_id"  gorm:"not null" `
	Product_id      uint
	OrdersID        string
	Product_Name    string
	Price           string
	Order_Status    string
	Payment_Status  string
	PaymentMethod   string
	Applied_Coupons string
	Total_amount    uint
}
type RazorPay struct {
	UserID          uint
	RazorPaymentId  string
	RazorPayOrderID string
	Signature       string
	OrderId         string
	AmountPaid      string
}
type Coupon struct {
	gorm.Model

	Coupon_Code string `json:"coupon_code"`
	Discount    uint   `json:"discount"`
	Quantity    uint   `json:"quantity"`
	Validity    int64  `json:"validity"`
}
type Applied_Coupons struct {
	gorm.Model
	UserID      uint
	Coupon_Code string `json:"coupon_code" `
}
type WishList struct {
	gorm.Model
	UserID     uint
	Product_id uint
}
