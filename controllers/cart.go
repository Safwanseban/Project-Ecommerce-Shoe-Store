package controllers

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	i "github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

var UserEmail string

func AddToCart(c *gin.Context) {
	userEmail := c.GetString("user")
	var user models.User
	var products models.Product
	i.DB.Raw("select id from users where email=?", userEmail).Scan(&user)
	var ProdtDetails struct {
		Product_id uint
		Quantity   uint
	}

	c.BindJSON(&ProdtDetails)
	//geting price for setting totalamount
	i.DB.Raw("select price ,stock from products where product_id=?", ProdtDetails.Product_id).Scan(&products)
	total := products.Price * ProdtDetails.Quantity
	prodid := ProdtDetails.Product_id
	prodqua := ProdtDetails.Quantity
	cart := models.Cart{
		ProductID:   ProdtDetails.Product_id,
		Quantity:    ProdtDetails.Quantity,
		UserId:      user.ID,
		Total_Price: total,
	}
	var Cart []models.Cart
	i.DB.Raw("select cart_id,product_id from carts where user_id=?", user.ID).Scan(&Cart) //geting all the cart details associated to user
	//ranging through the cart inorder to find if the product already exists
	for _, l := range Cart {
		fmt.Println("enterd")
		if l.ProductID == prodid {
			fmt.Println("in")
			i.DB.Raw("select quantity from carts where product_id=? and user_id=?", ProdtDetails.Product_id, user.ID).Scan(&Cart)
			totl := (prodqua + cart.Quantity) * products.Price
			i.DB.Raw("update carts set quantity=?,total_price=? where product_id=? and user_id=? ", prodqua+cart.Quantity, totl, prodid, user.ID).Scan(&Cart)

			c.JSON(400, gin.H{
				"msg": "quantity updated",
			})
			c.Abort()
			return
		}
	}

	record := i.DB.Create(&cart)
	// record:=i.DB.Raw("select quantity from carts where product_id=? and userid=?",ProdtDetails.Product_id,user.ID).Scan(&cart)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"err": record.Error.Error(),
		})
		c.Abort()
		return
	}
	i.DB.Raw("select product_name,brands.brands from products join carts on carts.product_id=products.product_id join brands on brands.id=products.brand_id  where products.product_id=? ", ProdtDetails.Product_id).Scan(&cart)
	c.JSON(200, gin.H{
		"userId": user.ID,
		"msg":    "added to cart",
	})
}

type Cartsinfo []struct {
	User_id      string
	Product_id   string
	Product_Name string
	Price        string
	Email        string
	Quantity     string
	Total_Amount uint
	Total_Price  string
}

func Viewcart(c *gin.Context) {
	var products models.Product
	var userid models.User
	var cart Cartsinfo

	useremail := c.GetString("user")
	fmt.Println(useremail)
	product := c.Query("product")
	newproduct, _ := strconv.Atoi(product)
	quanity := c.Query("quantity")
	newquantity, _ := strconv.Atoi(quanity)
	i.DB.Raw("select price from products where product_id=?", product).Scan(&products)
	i.DB.Raw("select id from users where email=?", useremail).Scan(&userid)
	total := products.Price * uint(newquantity) //specifying total quantity
	if newquantity >= 1 {
		i.DB.Raw("update carts set quantity=? , total_price=? where product_id=? and user_id=?", newquantity, total, newproduct, userid.ID).Scan(&cart)
	} else if newquantity <= 0 {
		i.DB.Raw("delete from carts where product_id =? and user_id=?", newproduct, userid.ID).Scan(&cart)
	}

	//select products.product_name,brands.brands,quantity,users.id,users.email from carts join products on products.product_id=carts.product_id join brands on brands.id=products.brand_id join users on users.id=carts.user_id;
	fmt.Println(UserEmail)
	record := i.DB.Raw("select  products.product_id, products.product_name,products.price,carts.user_id,users.email ,carts.quantity,total_price from carts join products on products.product_id=carts.product_id join users on carts.user_id=users.id where users.email=? ", useremail).Scan(&cart)
	if record.Error != nil {
		c.JSON(404, gin.H{"err": record.Error.Error()})
		c.Abort()
		return
	}
	var totalcartvalue uint

	i.DB.Raw("select sum(total_price) as total from carts where user_id=?", userid.ID).Scan(&totalcartvalue)
	c.JSON(200, gin.H{
		"cart":  cart,
		"total": totalcartvalue,
	})
}

func CheckOutAddAdress(c *gin.Context) {

	useremail := c.GetString("user")
	var user models.User
	i.DB.Raw("select id from users where email=?", useremail).Scan(&user)

	Name := c.PostForm("name")
	Phonenum := c.PostForm("phone_number")
	phonenum, _ := strconv.Atoi(Phonenum)
	pincod := c.PostForm("pincode")
	pincode, _ := strconv.Atoi(pincod)
	area := c.PostForm("area")
	houseadd := c.PostForm("house")
	landmark := c.PostForm("landmark")
	city := c.PostForm("city")
	address := models.Address{
		UserId:       user.ID,
		Name:         Name,
		Phone_number: phonenum,
		Pincode:      pincode,
		Area:         area,
		House:        houseadd,
		Landmark:     landmark,
		City:         city,
	}
	record := i.DB.Create(&address)
	if record.Error != nil {
		c.JSON(404, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"msg": "address added"})

}

var Address []struct {
	UserId       uint
	Address_id   uint
	Name         string
	Phone_number uint
	Pincode      uint
	Area         string
	House        string
	Landmark     string
	City         string
}

func Checkout(c *gin.Context) {
	var user models.User
	var cart models.Cart
	var carts Cartsinfo
	useremail := c.GetString("user")
	i.DB.Raw("select id from users where email=?", useremail).Scan(&user)
	precord := i.DB.Raw("select  products.product_id, products.product_name,products.price,carts.user_id,users.email ,carts.quantity,total_price from carts join products on products.product_id=carts.product_id join users on carts.user_id=users.id where users.email=? ", useremail).Scan(&carts)
	if precord.Error != nil {
		c.JSON(404, gin.H{"err": precord.Error.Error()})
		c.Abort()
		return
	}
	var totalcartvalue uint
	var address models.Address
	addres := c.Query("addressID")
	addressID, _ := strconv.Atoi(addres)
	PaymentMethod := c.Query("PaymentMethod")
	Coupons := c.Query("coupon")

	cod := "COD"
	razorpay := "Razorpay"
	notcompRazorpay := "Needs to complete razorpay payment"
	//getting total cart value
	i.DB.Raw("select sum(total_price) as total from carts where user_id=?", user.ID).Scan(&totalcartvalue)
	if PaymentMethod == cod || PaymentMethod == razorpay {
		for _, l := range carts {
			fmt.Println("entered into carts")
			pud := l.User_id
			Puid, _ := strconv.Atoi(pud)
			pid := l.Product_id
			Piid, _ := strconv.Atoi(pid)
			pname := l.Product_Name
			pprice := l.Price
			pPrice,_:=strconv.Atoi(pprice)
			pquantity:=l.Quantity
			pQuantity,_:=strconv.Atoi(pquantity)

			ordereditems := models.Orderd_Items{UserId: uint(Puid), Product_id: uint(Piid),
				Product_Name: pname, Price: pprice, OrdersID: CreateOrderId(),
				Order_Status: "confirmed", Payment_Status: "pending", Total_amount: uint(pQuantity)*uint(pPrice),
			}
			i.DB.Create(&ordereditems)

		}
	}

	//getting details from address
	record := i.DB.Raw("select address_id, user_id,name,phone_number,pincode,house,area,landmark,city from addresses where user_id=?", user.ID).Scan(&Address)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"err": record.Error.Error(),
		})
		c.Abort()
		return
	}

	i.DB.Raw("select address_id,user_id,name from addresses where address_id=?", addressID).Scan(&address)

	// coupon section

	var Coupondisc struct {
		Discount uint
		Count    uint
		Validity int64
	}
	var Appliedcoup struct {
		User_Id     uint
		Coupon_Code string
		Count       uint
	}
	var flag = 0
	// checking all possible conditions a coupon wont work
	if Coupons == "" {
		c.JSON(300, gin.H{
			"msg": "enter coupon if you have any coupon",
		})
		
	} else if Coupons != "" {
		flag = 1
		i.DB.Raw("select discount,validity,count(*) as count from coupons  where coupon_code=? group by discount,validity", Coupons).Scan(&Coupondisc)
		i.DB.Raw("select user_id,coupon_code,count(*) from applied_coupons where coupon_code=? and user_id=? group by user_id,coupon_code", Coupons, user.ID).Scan(&Appliedcoup)
		if Appliedcoup.Count > 0 {
			c.JSON(300, gin.H{
				"msg": "already applied",
			})
			flag = 2
		}
		fmt.Println(Coupondisc.Validity)
		if Coupondisc.Count <= 0 {
			fmt.Println("hai")
			c.JSON(300, gin.H{

				"msg": "not a valid coupon",
			})
			flag = 2

		}

		if Coupondisc.Validity < time.Now().Local().Unix() && Coupondisc.Validity > 1 {

			c.JSON(300, gin.H{
				"msg": "coupon expired",
			})
			flag = 2

		}

	}
	if flag == 1 {
		fmt.Println("hai")
		Discount := (totalcartvalue * Coupondisc.Discount) / 100
		totalcartvalue = totalcartvalue - Discount

	}

	c.JSON(300, gin.H{
		"address":          Address,
		"total cart value": totalcartvalue,
	})

	//creatimg a orderID
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(9999999999-1000000000) + 1000000000
	id := strconv.Itoa(value)
	orderID := "OID" + id
	fmt.Println(addressID)
	fmt.Println(address.Address_id)
	if address.UserId != user.ID {
		c.JSON(200, gin.H{
			"msg": "enter valid address id",
		})
	}

	if PaymentMethod == razorpay && addressID == int(address.Address_id) && address.UserId == user.ID {
		
		orders := models.Orders{
			UserId:          user.ID,
			Address_id:      uint(addressID),
			Order_id:        orderID,
			Order_Status:    "pending",
			PaymentMethod:   razorpay,
			Applied_Coupons: Coupons,
			Payment_Status:  notcompRazorpay,
			Total_Amount:    totalcartvalue,
		}

		result := i.DB.Create(&orders)
		if result.Error != nil {
			c.JSON(404, gin.H{"err": result.Error.Error()})
			c.Abort()
			return

		}
		var ordereditems models.Orderd_Items

		i.DB.Raw("update orderd_items set  order_status=?,payment_status=?,payment_method=? where user_id=?",  "orderplaced", notcompRazorpay, razorpay, user.ID).Scan(&ordereditems)
		if result.Error == nil {
			c.JSON(300, gin.H{
				"msg": "Go to the Razorpay Page for Order completion",
			})
			c.Abort()
			return
		}
	} else if PaymentMethod == cod && addressID == int(address.Address_id) && address.UserId == user.ID {
		fmt.Println("hai cod")
	
		orders := models.Orders{
			UserId:         user.ID,
			Address_id:     uint(addressID),
			PaymentMethod:  cod,
			Total_Amount:   totalcartvalue,
			Order_id:       orderID,
			Order_Status:   "order Placed",
			Payment_Status: "cash on deleivery",
		}

		result := i.DB.Create(&orders)
		coupns := models.Applied_Coupons{
			UserID:      user.ID,
			Coupon_Code: Coupons,
		}
		i.DB.Create(&coupns)
		if result.Error != nil {
			c.JSON(404, gin.H{"err": result.Error.Error()})
			c.Abort()
			return

		}

		var ordereditems models.Orderd_Items
		i.DB.Raw("update orderd_items set order_status=?,payment_method=? where user_id=?", "orderplaced", cod, user.ID).Scan(&ordereditems)

	} else {
		c.JSON(300, gin.H{
			"msg": "select payment method and address",
		})
		c.Abort()
		return
	}

	i.DB.Raw("delete from carts where user_id=?", user.ID).Scan(&cart)
	c.JSON(200, gin.H{"order": "order placed"})
}
