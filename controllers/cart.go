package controllers

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Safwanseban/Project-Ecommerce/initializers"
	i "github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
)

var (
	ErrcantFindProduct    = errors.New("cant find the product")
	ErrcantDecodeProduct  = errors.New("cant find the product")
	ErrUserIdnotvalid     = errors.New("this user not valid")
	ErrcantUpdateUser     = errors.New("cant this product to this cart")
	ErrCantRemoveItemCart = errors.New("cant remove this item from the cart")
	ErrCantGetItem        = errors.New("was unable to get the item for the cart")
	ErrCantBuyCartitem    = errors.New("cannot update the purchase")
)

var UserEmail string

func AddToCart(c *gin.Context) {
	var user models.User
	var products models.Product

	UserEmail = c.GetString("user")

	record := i.DB.Raw("select first_name,id from users where email=?", UserEmail).Scan(&user)

	product := c.Query("product")
	// Quantity:=c.Query("quantity")
	i.DB.Raw("select price from products where product_id=?", product).Scan(&products)

	record2 := i.DB.Raw("insert into carts(product_id,user_id,quantity,total_price) values(?,?,1,?)", product, user.ID, products.Price).Scan(&models.Cart{})

	if record.Error != nil {
		c.JSON(404, gin.H{
			"error": record.Error.Error(),
		})
	}
	fmt.Println(record2)

	// var userNum string
	// i.DB.Raw("select product_name,brands.brands from products join carts on carts.product_id=products.product_id join brands on brands.id=products.brand_id  where products.product_id=? ", product).Scan(&Cart{})
	c.JSON(200, gin.H{
		"userId":   user.ID,
		"Username": user.First_Name,
	})
}

type Cartsinfo []struct {
	User_id      string
	Product_id   string
	Product_Name string
	Price        string
	Email        string
	Quantity     string
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
	cod := "COD"

	//getting total cart value
	i.DB.Raw("select sum(total_price) as total from carts where user_id=?", user.ID).Scan(&totalcartvalue)
	if PaymentMethod!=""{
	for _, i := range carts{
		fmt.Println("entered into carts")
		pud:=i.User_id
		Puid, _ := strconv.Atoi(pud)
		pid := i.Product_id
		Piid, _ := strconv.Atoi(pid)
		pname := i.Product_Name
		pprice := i.Price
		ordereditems := models.Orderd_Items{UserId: uint(Puid),Product_id: uint(Piid),
			Product_Name: pname, Price: pprice, OrdersID: CreateOrderId(),
			Order_Status: "confirmed",Payment_Status: "pending",Total_amount: totalcartvalue,
		}
		initializers.DB.Create(&ordereditems)

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

	c.JSON(404, gin.H{
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
	if address.UserId != user.ID{
		c.JSON(404,gin.H{
			"msg":"enter valid address id",
		})
	}
	if PaymentMethod == cod && addressID == int(address.Address_id) && address.UserId == user.ID {

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
		if result.Error != nil {
			c.JSON(404, gin.H{"err": result.Error.Error()})
			c.Abort()
			return

		}
		var ordereditems models.Orderd_Items
		i.DB.Raw("update orderd_items set order_status=?,payment_method=? where user_id=?","orderplaced",cod,user.ID).Scan(&ordereditems)
	} else {
		c.JSON(404, gin.H{
			"msg": "select payment method and address",
		})
		c.Abort()
		return
	}

	i.DB.Raw("delete from carts where user_id=?", user.ID).Scan(&cart)
	c.JSON(200, gin.H{"order": "order placed"})
}
