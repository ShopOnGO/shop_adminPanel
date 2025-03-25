package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "admin/pkg/service"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("error due to connection: %v", err)
	}
	defer conn.Close()

	categoryClient := pb.NewCategoryServiceClient(conn)
	brandClient := pb.NewBrandServiceClient(conn)
	linkClient := pb.NewLinkServiceClient(conn)
	productClient := pb.NewProductServiceClient(conn)
	userClient := pb.NewUserServiceClient(conn)
	statClient := pb.NewStatServiceClient(conn)
	homeClient := pb.NewHomeServiceClient(conn)

	DeleteAllCategories(categoryClient)
	testCategoryService(categoryClient)

	DeleteAllBrands(brandClient)
	testBrandService(brandClient)

	DeleteAllLinks(linkClient)
	testLinkService(linkClient)

	DeleteAllProducts(productClient)
	testProductService(productClient)

	DeleteAllUsers(userClient)
	testUserService(userClient)

	testStatService(statClient)

	testHomeService(homeClient)

}

func testCategoryService(client pb.CategoryServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createResp, err := client.CreateCategory(ctx, &pb.CreateCategoryRequest{
		Name:        "boots",
		Description: "Stylish winter boots",
	})
	if err != nil {
		log.Fatalf("error due to category creation: %v", err)
	}
	fmt.Printf("✅ Category created: ID=%d, Name=%s\n", createResp.Category.Model.Id, createResp.Category.Name)

	getResp, err := client.FindCategoryByID(ctx, &pb.FindCategoryByIDRequest{
		Id: createResp.Category.Model.Id,
	})
	if err != nil {
		log.Fatalf("failed to get category: %v", err)
	}
	fmt.Printf("✅ Category found: ID=%d, Name=%s\n", getResp.Category.Model.Id, getResp.Category.Name)

	updateResp, err := client.UpdateCategory(ctx, &pb.UpdateCategoryRequest{
		Id:          getResp.Category.Model.Id,
		Name:        "boots - Updated",
		Description: "Updated description",
	})
	if err != nil {
		log.Fatalf("error during category update: %v", err)
	}
	fmt.Printf("✅ Category updated: ID=%d, Name=%s\n", updateResp.Category.Model.Id, updateResp.Category.Name)

	featuredResp, err := client.GetFeaturedCategories(ctx, &pb.GetFeaturedCategoriesRequest{
		Amount: 5,
	})
	if err != nil {
		log.Fatalf("error during category retrieval: %v", err)
	}
	fmt.Println("✅ Featured categories:")
	for _, category := range featuredResp.Categories {
		fmt.Printf("   - ID=%d, Name=%s\n", category.Model.Id, category.Name)
	}

	_, err = client.DeleteCategory(ctx, &pb.DeleteCategoryByNameRequest{
		Name: createResp.Category.Name,
	})
	if err != nil {
		log.Fatalf("error during category deletion: %v", err)
	}
	fmt.Printf("✅ Category with ID=%d deleted (soft delete)\n", createResp.Category.Model.Id)
}

func DeleteAllCategories(client pb.CategoryServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := client.GetFeaturedCategories(ctx, &pb.GetFeaturedCategoriesRequest{Amount: 0, Unscoped: true})
	if err != nil {
		log.Fatalf("error during category retrieval: %v", err)
	}

	log.Println("📌 Retrieved categories (including deleted ones):")
	for _, category := range resp.Categories {
		log.Printf("➡ ID=%d, Name=%s, DeletedAt=%v", category.Model.Id, category.Name, category.Model.DeletedAt)
	}

	for _, category := range resp.Categories {
		_, err := client.DeleteCategory(ctx, &pb.DeleteCategoryByNameRequest{Name: category.Name, Unscoped: true})
		if err != nil {
			log.Printf("❌ Error deleting category Name=%s: %v", category.Name, err)
		} else {
			log.Printf("✅ Category Name=%s permanently deleted", category.Name)
		}
	}
}

func testBrandService(client pb.BrandServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createResp, err := client.CreateBrand(ctx, &pb.CreateBrandRequest{
		Name:        "Nike",
		Description: "Top sports brand",
	})
	if err != nil {
		log.Fatalf("error due to brand creation: %v", err)
	}
	fmt.Printf("✅ Brand created: ID=%d, Name=%s\n", createResp.Brand.Model.Id, createResp.Brand.Name)

	getResp, err := client.FindBrandByID(ctx, &pb.FindBrandByIDRequest{
		Id: uint64(createResp.Brand.Model.Id),
	})
	if err != nil {
		log.Fatalf("failed to get brand: %v", err)
	}
	fmt.Printf("✅ Brand found: ID=%d, Name=%s\n", getResp.Brand.Model.Id, getResp.Brand.Name)

	updateResp, err := client.UpdateBrand(ctx, &pb.Brand{
		Model:       &pb.Model{Id: getResp.Brand.Model.Id},
		Name:        "Nike - Updated",
		Description: "Updated brand description",
	})
	if err != nil {
		log.Fatalf("error during brand update: %v", err)
	}
	fmt.Printf("✅ Brand updated: ID=%d, Name=%s\n", updateResp.Brand.Model.Id, updateResp.Brand.Name)

	featuredResp, err := client.GetFeaturedBrands(ctx, &pb.GetFeaturedBrandsRequest{
		Amount:   5,
		Unscoped: true,
	})
	if err != nil {
		log.Fatalf("error during brand retrieval: %v", err)
	}
	fmt.Println("✅ Featured brands:")
	for _, brand := range featuredResp.Brands {
		fmt.Printf("   - ID=%d, Name=%s\n", brand.Model.Id, brand.Name)
	}

	_, err = client.DeleteBrand(ctx, &pb.DeleteBrandRequest{
		Name: createResp.Brand.Name,
	})
	if err != nil {
		log.Fatalf("error during brand deletion: %v", err)
	}
	fmt.Printf("✅ Brand with ID=%d deleted (soft delete)\n", createResp.Brand.Model.Id)
}

func DeleteAllBrands(client pb.BrandServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := client.GetFeaturedBrands(ctx, &pb.GetFeaturedBrandsRequest{Amount: 0, Unscoped: true})
	if err != nil {
		log.Fatalf("error during brand retrieval: %v", err)
	}

	log.Println("📌 Retrieved brands (including deleted ones):")
	for _, brand := range resp.Brands {
		log.Printf("➡ ID=%d, Name=%s, DeletedAt=%v", brand.Model.Id, brand.Name, brand.Model.DeletedAt)
	}

	for _, brand := range resp.Brands {
		_, err := client.DeleteBrand(ctx, &pb.DeleteBrandRequest{Name: brand.Name, Unscoped: true})
		if err != nil {
			log.Printf("❌ Error deleting brand Name=%s: %v", brand.Name, err)
		} else {
			log.Printf("✅ Brand Name=%s permanently deleted", brand.Name)
		}
	}
}

func DeleteAllLinks(client pb.LinkServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetAllLinks(ctx, &pb.GetAllLinksRequest{Limit: 100, Offset: 0})
	if err != nil {
		log.Fatalf("❌ Error fetching links: %v", err)
	}

	log.Println("📌 Retrieved links:")
	for _, link := range resp.Links {
		log.Printf("➡ ID=%d, Hash=%s, URL=%s, DeletedAt=%v", link.Model.Id, link.Hash, link.Url, link.Model.DeletedAt)
	}

	for _, link := range resp.Links {
		_, err := client.Delete(ctx, &pb.DeleteLinkRequest{Id: link.Model.Id, Unscoped: true})
		if err != nil {
			log.Printf("❌ Error permanently deleting link ID=%d: %v", link.Model.Id, err)
		} else {
			log.Printf("✅ Link ID=%d permanently deleted", link.Model.Id)
		}
	}
}

func testLinkService(client pb.LinkServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 🔹 Создание ссылки
	createResp, err := client.Create(ctx, &pb.CreateLinkRequest{
		Url: "https://example.com",
	})
	if err != nil {
		log.Fatalf("❌ Error creating link: %v", err)
	}
	fmt.Printf("✅ Link created: ID=%d, Hash=%s, URL=%s\n", createResp.Link.Model.Id, createResp.Link.Hash, createResp.Link.Url)

	// 🔹 Получение ссылки по хэшу
	getResp, err := client.GetLinkByHash(ctx, &pb.GetLinkByHashRequest{
		Hash: createResp.Link.Hash,
	})
	if err != nil {
		log.Fatalf("❌ Error fetching link by hash: %v", err)
	}
	fmt.Printf("✅ Link found by hash: ID=%d, Hash=%s, URL=%s\n", getResp.Link.Model.Id, getResp.Link.Hash, getResp.Link.Url)

	// 🔹 Обновление ссылки
	updateResp, err := client.Update(ctx, &pb.UpdateLinkRequest{
		Id:   getResp.Link.Model.Id,
		Url:  "https://updated-example.com",
		Hash: getResp.Link.Hash,
	})
	if err != nil {
		log.Fatalf("❌ Error updating link: %v", err)
	}
	fmt.Printf("✅ Link updated: ID=%d, Hash=%s, New URL=%s\n", updateResp.Link.Model.Id, updateResp.Link.Hash, updateResp.Link.Url)

	// 🔹 Получение всех ссылок
	allLinksResp, err := client.GetAllLinks(ctx, &pb.GetAllLinksRequest{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		log.Fatalf("❌ Error fetching all links: %v", err)
	}
	fmt.Println("✅ Retrieved links:")
	for _, link := range allLinksResp.Links {
		fmt.Printf("   - ID=%d, Hash=%s, URL=%s\n", link.Model.Id, link.Hash, link.Url)
	}

	// 🔹 Удаление ссылки (мягкое)
	_, err = client.Delete(ctx, &pb.DeleteLinkRequest{
		Id: createResp.Link.Model.Id,
	})
	if err != nil {
		log.Fatalf("❌ Error soft deleting link: %v", err)
	}
	fmt.Printf("✅ Link with ID=%d soft deleted\n", createResp.Link.Model.Id)
}

func DeleteAllProducts(client pb.ProductServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetFeaturedProducts(ctx, &pb.FeaturedRequest{
		Amount:         2,
		Random:         true,
		IncludeDeleted: true, // Retrieve all products, including deleted ones
	})
	if err != nil {
		log.Fatalf("❌ Error retrieving the product list: %v", err)
	}

	log.Println("📌 Found products (including deleted ones):")
	for _, product := range resp.Products {
		log.Printf("➡ ID=%d, Name=%s, DeletedAt=%v", product.Model.Id, product.Name, product.Model.DeletedAt)
	}

	for _, product := range resp.Products {
		_, err := client.DeleteProduct(ctx, &pb.DeleteProductRequest{
			Id:       uint64(product.Model.Id),
			Unscoped: true, // Hard delete
		})
		if err != nil {
			log.Printf("❌ Error deleting product ID=%d: %v", product.Model.Id, err)
		} else {
			log.Printf("✅ Product ID=%d permanently deleted", product.Model.Id)
		}
	}
}

func testProductService(client pb.ProductServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 🔹 Create a new product
	createResp, err := client.CreateProduct(ctx, &pb.Product{
		Name:        "Test Boots",
		Description: "Winter Boots",
		CategoryId:  1,
		BrandId:     1,
		Stock:       10,
		IsAvailable: true,
	})
	if err != nil {
		log.Fatalf("❌ Error creating product: %v", err)
	}
	fmt.Printf("✅ Product created: ID=%d, Name=%s\n", createResp.Product.Model.Id, createResp.Product.Name)

	// 🔹 Retrieve products by category ID
	getResp, err := client.GetProductsByCategory(ctx, &pb.CategoryRequest{CategoryId: 1})
	if err != nil {
		log.Fatalf("❌ Error retrieving products by category: %v", err)
	}
	fmt.Println("✅ Products found in category 1:")
	for _, product := range getResp.Products {
		fmt.Printf("   - ID=%d, Name=%s\n", product.Model.Id, product.Name)
	}

	// 🔹 Search product by name
	nameResp, err := client.GetProductsByName(ctx, &pb.NameRequest{Name: "Test Boots"})
	if err != nil {
		log.Fatalf("❌ Error searching for product by name: %v", err)
	}
	fmt.Println("✅ Products found by name 'Test Boots':")
	for _, product := range nameResp.Products {
		fmt.Printf("   - ID=%d, Name=%s\n", product.Model.Id, product.Name)
	}

	// 🔹 Update product
	updateResp, err := client.UpdateProduct(ctx, &pb.Product{
		Model:       &pb.Model{Id: createResp.Product.Model.Id},
		Name:        "Updated Boots",
		Description: "Updated Winter Boots",
		CategoryId:  1,
		BrandId:     2,
		Stock:       20,
		IsAvailable: true,
	})
	if err != nil {
		log.Fatalf("❌ Error updating product: %v", err)
	}
	fmt.Printf("✅ Product updated: ID=%d, Name=%s\n", updateResp.Product.Model.Id, updateResp.Product.Name)

	// 🔹 Retrieve featured products
	featuredResp, err := client.GetFeaturedProducts(ctx, &pb.FeaturedRequest{
		Amount:         5,
		Random:         true,
		IncludeDeleted: false, // Only active products
	})
	if err != nil {
		log.Fatalf("❌ Error retrieving featured products: %v", err)
	}
	fmt.Println("✅ Featured products:")
	for _, product := range featuredResp.Products {
		fmt.Printf("   - ID=%d, Name=%s\n", product.Model.Id, product.Name)
	}

	// 🔹 Soft delete the product
	_, err = client.DeleteProduct(ctx, &pb.DeleteProductRequest{
		Id:       uint64(createResp.Product.Model.Id),
		Unscoped: false, // Soft delete
	})
	if err != nil {
		log.Fatalf("❌ Error performing soft delete on product: %v", err)
	}
	fmt.Printf("✅ Product ID=%d deleted (soft delete)\n", createResp.Product.Model.Id)

	// 🔹 Verify that the product is not shown in the normal product list
	featuredAfterDeleteResp, err := client.GetFeaturedProducts(ctx, &pb.FeaturedRequest{
		Amount:         5,
		Random:         true,
		IncludeDeleted: false, // Only active products
	})
	if err != nil {
		log.Fatalf("❌ Error checking product list after deletion: %v", err)
	}
	for _, product := range featuredAfterDeleteResp.Products {
		if product.Model.Id == createResp.Product.Model.Id {
			log.Fatalf("❌ Error: Deleted product ID=%d is still visible", product.Model.Id)
		}
	}
	fmt.Println("✅ Product is not shown in the active list (soft delete works)")

	// 🔹 Hard delete the product
	_, err = client.DeleteProduct(ctx, &pb.DeleteProductRequest{
		Id:       uint64(createResp.Product.Model.Id),
		Unscoped: true, // Permanent deletion
	})
	if err != nil {
		log.Fatalf("❌ Error performing hard delete on product: %v", err)
	}
	fmt.Printf("✅ Product ID=%d permanently deleted\n", createResp.Product.Model.Id)

	// 🔹 Verify that the product is completely removed, even from deleted products
	deletedProductsResp, err := client.GetFeaturedProducts(ctx, &pb.FeaturedRequest{
		Amount:         5,
		Random:         true,
		IncludeDeleted: true, // Check deleted products as well
	})
	if err != nil {
		log.Fatalf("❌ Error checking the list of deleted products: %v", err)
	}
	for _, product := range deletedProductsResp.Products {
		if product.Model.Id == createResp.Product.Model.Id {
			log.Fatalf("❌ Error: Product ID=%d is still present in the database", product.Model.Id)
		}
	}
	fmt.Println("✅ Product completely deleted (hard delete works)")
}

func DeleteAllUsers(client pb.UserServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users := []string{"test@example.com", "test-updated@example.com"}

	for _, email := range users {
		resp, err := client.FindUserByEmail(ctx, &pb.EmailRequest{Email: email})
		if err != nil {
			log.Printf("🔹 User %s not found, skipping...", email)
			continue
		}

		_, err = client.DeleteUser(ctx, &pb.DeleteUserRequest{
			Id:       uint64(resp.User.Model.Id),
			Unscoped: true,
		})
		if err != nil {
			log.Printf("❌ Error deleting user %s: %v", email, err)
		} else {
			log.Printf("✅ User %s successfully deleted", email)
		}
	}
}

func testUserService(client pb.UserServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 📌 Create a new user
	createResp, err := client.CreateUser(ctx, &pb.User{
		Email: "test@example.com",
		Name:  "Test User",
		Role:  pb.UserRole_USER_ROLE_BUYER,
	})
	if err != nil {
		log.Fatalf("❌ Error creating user: %v", err)
	}
	fmt.Printf("✅ User created: ID=%d, Email=%s, Name=%s\n",
		createResp.User.Model.Id, createResp.User.Email, createResp.User.Name)

	// 📌 Find user by email
	getResp, err := client.FindUserByEmail(ctx, &pb.EmailRequest{
		Email: createResp.User.Email,
	})
	if err != nil {
		log.Fatalf("❌ Error finding user: %v", err)
	}
	fmt.Printf("✅ User found: ID=%d, Email=%s, Name=%s\n",
		getResp.User.Model.Id, getResp.User.Email, getResp.User.Name)

	// 📌 Update user details
	updateResp, err := client.UpdateUser(ctx, &pb.User{
		Model: &pb.Model{Id: getResp.User.Model.Id},
		Email: "test-updated@example.com",
		Name:  "Updated User",
		Role:  pb.UserRole_USER_ROLE_ADMIN,
	})
	if err != nil {
		log.Fatalf("❌ Error updating user: %v", err)
	}
	fmt.Printf("✅ User updated: ID=%d, Email=%s, Name=%s, Role=%s\n",
		updateResp.User.Model.Id, updateResp.User.Email, updateResp.User.Name, updateResp.User.Role)

	// 📌 Soft delete user
	_, err = client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id:       uint64(getResp.User.Model.Id),
		Unscoped: false,
	})
	if err != nil {
		log.Fatalf("❌ Error soft deleting user: %v", err)
	}
	fmt.Printf("✅ User with ID=%d has been soft deleted\n", getResp.User.Model.Id)

	// 📌 Hard delete user
	_, err = client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id:       uint64(getResp.User.Model.Id),
		Unscoped: true,
	})
	if err != nil {
		log.Fatalf("❌ Error permanently deleting user: %v", err)
	}
	fmt.Printf("✅ User with ID=%d has been permanently deleted\n", getResp.User.Model.Id)
}

func testStatService(client pb.StatServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 📌 Add a click to a link
	linkID := uint64(123) // Replace with an actual link ID
	_, err := client.AddClick(ctx, &pb.ClickRequest{LinkId: linkID})
	if err != nil {
		log.Fatalf("❌ Error adding click for link ID=%d: %v", linkID, err)
	}
	fmt.Printf("✅ Click successfully added for link ID=%d\n", linkID)
}

func testHomeService(client pb.HomeServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 📌 Request homepage data
	resp, err := client.GetHomeData(ctx, &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("❌ Error retrieving homepage data: %v", err)
	}

	// Check for errors in the response
	if resp.Error != nil {
		log.Fatalf("❌ Service error: Code=%d, Message=%s", resp.Error.Code, resp.Error.Message)
	}

	// 📌 Print retrieved data
	fmt.Println("✅ Homepage data successfully retrieved:")
	fmt.Printf("🔹 Categories: %d\n", len(resp.Categories))
	for _, category := range resp.Categories {
		fmt.Printf("   - ID=%d, Name=%s\n", category.Model.Id, category.Name)
	}

	fmt.Printf("🔹 Featured products: %d\n", len(resp.FeaturedProducts))
	for _, product := range resp.FeaturedProducts {
		fmt.Printf("   - ID=%d, Name=%s\n", product.Model.Id, product.Name)
	}

	fmt.Printf("🔹 Brands: %d\n", len(resp.Brands))
	for _, brand := range resp.Brands {
		fmt.Printf("   - ID=%d, Name=%s\n", brand.Model.Id, brand.Name)
	}
}
