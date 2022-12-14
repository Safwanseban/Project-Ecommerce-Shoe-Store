package controllers

import (
	"fmt"
	"math"
	"path/filepath"
	"strconv"

	"github.com/Safwanseban/Project-Ecommerce/initializers"
	i "github.com/Safwanseban/Project-Ecommerce/initializers"
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

func ListingAllCat(c *gin.Context) {

	var brandss []models.Brand
	var catogorry []models.Catogory
	var shoesizess []models.ShoeSize

	brands := i.DB.Find(&brandss)

	if brandsearch := c.Query("brandsearch"); brandsearch != "" {
		brands := i.DB.Where("brands LIKE ?", "%"+brandsearch+"%").Find(&brandss)
		if brands.Error != nil {
			c.JSON(404, gin.H{
				"err": brands.Error.Error(),
			})
			c.Abort()
			return
		}
	}

	catogory := i.DB.Find(&catogorry)

	if catogorysearch := c.Query("catogorysearch"); catogorysearch != "" {
		catogory = i.DB.Where("catogory LIKE ?", "%"+catogorysearch+"%").Find(&catogorry)
		if catogory.Error != nil {
			c.JSON(404, gin.H{
				"err": brands.Error.Error(),
			})
			c.Abort()
			return
		}
	}
	size := initializers.DB.Find(&shoesizess)
	if sizesearch := c.Query("sizesearch"); sizesearch != "" {
		sizes, _ := strconv.Atoi(sizesearch)
		size = i.DB.Where("size = ?", sizes).Find(&shoesizess)
		if size.Error != nil {
			c.JSON(404, gin.H{
				"err": brands.Error.Error(),
			})
			c.Abort()
			return
		}
	}

	c.JSON(200, gin.H{
		"available brands":     brandss,
		"available carogories": catogorry,
		"available sizes":      shoesizess,
	})
}
func ApplyBrandDiscount(c *gin.Context) {

	var brand struct {
		Brand_id uint
		Discount uint
	}
	if err := c.BindJSON(&brand); err != nil {
		c.JSON(404, gin.H{
			"err": err.Error(),
		})
	}

	record := initializers.DB.Model(&models.Brand{}).Where("id = ?", brand.Brand_id).Update("discount", brand.Discount)
	if record.Error == nil {
		c.JSON(200, gin.H{
			"discount": brand.Discount,
			"msg":      "brand discount applied"})
	}

}
func ProductAdding(c *gin.Context) { //Admin

	prodname := c.PostForm("productname")
	price := c.PostForm("price")
	Price, _ := strconv.Atoi(price)
	description := c.PostForm("description")
	color := c.PostForm("color")
	brand := c.PostForm("brandID")
	brands, _ := strconv.Atoi(brand)
	stock := c.PostForm("stock")
	Stock, _ := strconv.Atoi(stock)
	catogory := c.PostForm("catogoryID")
	catogoryy, _ := strconv.Atoi(catogory)
	size := c.PostForm("sizeID")

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
	BrandDiscount := c.PostForm("BrandDiscount")
	brandDiscount, _ := strconv.Atoi(BrandDiscount)
	var Discount int
	//inserting brand discount on to the products
	insert := initializers.DB.Raw("update brands set discount=? where id=?", brandDiscount, brands).Scan(&models.Brand{})
	if insert.Error != nil {
		c.JSON(404, gin.H{
			"err": insert.Error.Error(),
		})
		c.Abort()
		return
	}
	//comparing whcih type of discount is greater
	if brandDiscount > discount {
		Discount = (Price * brandDiscount) / 100

	} else {
		Discount = (Price * discount) / 100
	}

	// Discount = (Price * discount) / 100
	var count uint
	initializers.DB.Raw("select count(*) from products where product_name=?", prodname).Scan(&count)
	fmt.Println(count)
	if count > 0 {
		c.JSON(404, gin.H{
			"msg": "A product with same name already exists",
		})
		c.Abort()
		return
	}
	product := models.Product{

		Product_name: prodname,

		Price:        uint(Price) - uint(Discount),
		Color:        color,
		Description:  description,
		Actual_Price: uint(Price),
		Discount:     uint(discount),

		Brand_id:   uint(brands),
		CatogoryID: uint(catogoryy),
		ShoeSizeID: uint(Size),
		Image:      image,

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
	if count <= 0 {
		c.JSON(404, gin.H{
			"msg": "product doesnot exist",
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

func TestProductsView(c *gin.Context) {

	// i.DB.Joins("join brands on products.brand_id = brands.id ").Joins("join catogories on products.catogory_id=catogories.id").Joins("join shoe_sizes on products.shoe_size_id=shoe_sizes.id").Scan(&Products)
	record := i.DB.Raw("SELECT product_id,product_name,actual_price,price,image,color,description,sub_pic1,sub_pic2,stock,brands.brands,catogories.catogory,shoe_sizes.size FROM products join brands on products.brand_id = brands.id join catogories on products.catogory_id=catogories.id join shoe_sizes on products.shoe_size_id=shoe_sizes.id").Scan(&Products)
	fmt.Println(record)

	if s := c.Query("search"); s != "" { //search
		i.DB.Raw("SELECT product_id,product_name,actual_price,price,image,color,description,sub_pic1,sub_pic2,stock,brands.brands,catogories.catogory,shoe_sizes.size FROM products join brands on products.brand_id = brands.id join catogories on products.catogory_id=catogories.id join shoe_sizes on products.shoe_size_id=shoe_sizes.id where product_name like ?", "%"+s+"%").Scan(&Products)

	}
	if sort := c.Query("sort"); sort != "" { //sort
		i.DB.Raw("SELECT product_id,product_name,actual_price,price,image,color,description,sub_pic1,sub_pic2,stock,brands.brands,catogories.catogory,shoe_sizes.size FROM products join brands on products.brand_id = brands.id join catogories on products.catogory_id=catogories.id join shoe_sizes on products.shoe_size_id=shoe_sizes.id  order by price ?", sort).Scan(&Products)
	}
	c.JSON(200, gin.H{
		"products": Products,
	})
}

func ProductsView(c *gin.Context) { //user
	brandFilter := c.Query("brandFilter")
	catogoryFilter := c.Query("catogoryFilter")
	sizefilter := c.Query("sizeFilter")
	sizeFilter, _ := strconv.Atoi(sizefilter)

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
	fmt.Println(sizeFilter)

	if brandFilter == "" && catogoryFilter == "" && sizefilter == "" {
		sql = fmt.Sprintf("%s ", sql)
	} else if brandFilter != "" && catogoryFilter != "" && sizefilter == "" {
		sql = fmt.Sprintf("%s where (brands='%s' and catogory='%s')", sql, brandFilter, catogoryFilter)
	} else if brandFilter != "" && catogoryFilter != "" && sizefilter != "" {
		sql = fmt.Sprintf("%s where (brands='%s' and catogory='%s' and size=%d)", sql, brandFilter, catogoryFilter, sizeFilter)
	} else if brandFilter != "" && catogoryFilter == "" && sizefilter != "" {
		sql = fmt.Sprintf("%s where (brands='%s'  and size=%d)", sql, brandFilter, sizeFilter)
	} else if brandFilter == "" && catogoryFilter != "" && sizefilter != "" {
		sql = fmt.Sprintf("%s where (catogory='%s'  and size=%d)", sql, catogoryFilter, sizeFilter)
	} else if brandFilter != "" && catogoryFilter == "" && sizefilter == "" {
		sql = fmt.Sprintf("%s where brands='%s'", sql, brandFilter)
	} else if brandFilter == "" && catogoryFilter != "" && sizefilter == "" {
		sql = fmt.Sprintf("%s where catogory='%s'", sql, catogoryFilter)
	} else if brandFilter == "" && catogoryFilter == "" && sizefilter != "" {
		sql = fmt.Sprintf("%s where size=%d", sql, sizeFilter)
	}

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
