package controllers

import (
	"fmt"
	"math"
	"path/filepath"
	"strconv"

	"github.com/Safwanseban/Project-Ecommerce/initializers"
	"github.com/Safwanseban/Project-Ecommerce/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var Products []struct {
	Product_ID   uint
	Product_Name string
	Actual_price uint
	Price        string
	Image        string
	SubPic1      string
	SubPic2      string
	Description  string
	Color        string
	Brands       string
	Stock        uint
	Catogory     string
	Size         uint
}

func ProductAdding(c *gin.Context) { //Admin
	// {
	// 	"product_name":"Nike Running Shoes",
	// 	"price":2500,
	// 	"description":"men running shoes",
	// 	"image":"https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/195110/01/sv01/fnd/IND/fmt/png/Blaze-Unisex-Shoes",
	// 	"subpic1":"https://images.puma.com/image/upload/f_auto,q_auto,b_rgb:fafafa,w_2000,h_2000/global/195110/01/sv01/fnd/IND/fmt/png/Blaze-Unisex-Shoes",
	// 	"color":"black&white",
	// 	"brand_id":1  ,
	// 	"stock":2,
	// 	"

	// }
	// productid:=c.PostForm("productid")
	// productID,_:=strconv.Atoi(productid)
	prodname := c.PostForm("productname")
	price := c.PostForm("price")
	Price, _ := strconv.Atoi(price)
	description := c.PostForm("description")
	color := c.PostForm("color")
	brands := c.PostForm("brand")
	stock := c.PostForm("stock")
	Stock, _ := strconv.Atoi(stock)
	catogory := c.PostForm("catogory")
	size := c.PostForm("size")
	Size, _ := strconv.Atoi(size)
	// images adding
	imagepath, _ := c.FormFile("image")
	extension := filepath.Ext(imagepath.Filename)
	image := uuid.New().String() + extension
	c.SaveUploadedFile(imagepath, "./public/images"+image)

	subpic1, _ := c.FormFile("subpic1")
	extension = filepath.Ext(subpic1.Filename)
	subpic11 := uuid.New().String() + extension
	c.SaveUploadedFile(subpic1, "./public/images"+image)

	subpic2path, _ := c.FormFile("subpic2")
	extension = filepath.Ext(subpic2path.Filename)
	subpic2 := uuid.New().String() + extension
	c.SaveUploadedFile(subpic2path, "./public/images"+image)
	discont := c.PostForm("discount")
	discount, _ := strconv.Atoi(discont)

	Discount := (Price * discount) / 100
	product := models.Product{

		Product_name: prodname,
		Price:        uint(Price)-uint(Discount),
		Color:        color,
		Description:  description,
		Actual_Price: uint(Price),
		Discount:     uint(discount),
		
		Brand: models.Brand{
			Brands: brands,
		},
		Catogory: models.Catogory{
			Catogory: catogory,
		},
		ShoeSize: models.ShoeSize{
			Size: uint(Size),
		},
		Image: image,

		SubPic1: subpic11,
		SubPic2: subpic2,
		Stock:   uint(Stock),
	}

	record := initializers.DB.Create(&product)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"msg": "product already exists",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"msg": "added succesfully",
	})

	// var products models.Product
	// if err := c.ShouldBindJSON(&products); err != nil {
	// 	c.JSON(404, gin.H{"err": err.Error()})
	// 	c.Abort()
	// 	return
	// }
	// record := initializers.DB.Create(&products)
	// if record.Error != nil {
	// 	c.JSON(404, gin.H{"err": record.Error.Error()})
	// 	c.Abort()
	// 	return
	// }
	// c.JSON(http.StatusOK, gin.H{"productname": products.Product_name,
	// 	"color": products.Color, "price": products.Price})

}

type Editproduct struct {
	Product_name string `json:"product_name"`
	Price        uint   `json:"price"`
	Image        string `json:"image"`
	Color        string `json:"color"`
}

func Editproducts(c *gin.Context) { //admin

	params := c.Param("id")

	var editproducts Editproduct
	if err := c.ShouldBindJSON(&editproducts); err != nil {
		c.JSON(404, gin.H{"err": err.Error()})
	}
	var products models.Product
	record := initializers.DB.Model(products).Where("product_id=?", params).Updates(models.Product{Product_name: editproducts.Product_name,
		Price: editproducts.Price, Color: editproducts.Color})
	if record.Error != nil {
		c.JSON(404, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"msg": "updated successfully"})

}

func DeleteProductById(c *gin.Context) { //admin
	params := c.Param("id")
	var products models.Product
	var count uint
initializers.DB.Raw("select count(product_id) from products where product_id=?", params).Scan(&count)
if count<=0{
	c.JSON(404,gin.H{
		"msg":"product doesnot exist",
	})
	c.Abort()
	return
}

	record := initializers.DB.Raw("delete from products where product_id=?", params).Scan(&products)
	if record.Error != nil {
		c.JSON(404, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"msg": "deleted successfully"})
}

func ProductsView(c *gin.Context) { //user
	sql := "SELECT product_id,product_name,actual_price,price,image,color,description,sub_pic1,sub_pic2,stock,brands.brands,catogories.catogory,shoe_sizes.size FROM products join brands on products.brand_id = brands.id join catogories on products.catogory_id=catogories.id join shoe_sizes on products.shoe_size_id=shoe_sizes.id"
	if s := c.Query("s"); s != "" { //search
		sql = fmt.Sprintf("%s WHERE product_name like'%%%s%%'", sql, s)
	}
	if sort := c.Query("sort"); sort != "" { //sort
		sql = fmt.Sprintf("%s order by price %s", sql, sort)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage := 5
	var total int64

	initializers.DB.Raw(sql).Count(&total)

	sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, perPage, (page-1)*perPage)

	record := initializers.DB.Raw(sql).Scan(&Products)
	if record.Error != nil {
		c.JSON(404, gin.H{"err": record.Error.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"products": Products,
		"total": total, "page": page,
		"last_page": math.Ceil(float64(total) / float64(perPage)),
	})

}
func GetProductByID(c *gin.Context) { //user
	params := c.Param("id")
	// var product models.Product
	record := initializers.DB.Raw("SELECT product_id,product_name,price,image,color,stock,brands.brands FROM products join brands on products.brand_id = brands.id where product_id=?", params).Scan(&Products)
	if record.Error != nil {
		c.JSON(404, gin.H{"err": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"product": Products})

}
