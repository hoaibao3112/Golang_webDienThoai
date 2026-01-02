# Images Directory

Thư mục này chứa tất cả các ảnh static cho ứng dụng.

## Cấu trúc thư mục

- `brands/` - Logo thương hiệu (samsung.png, apple.png, xiaomi.png, v.v.)
- `categories/` - Ảnh danh mục sản phẩm
- `products/` - Ảnh sản phẩm

## Cách sử dụng

### 1. Thêm ảnh vào thư mục
Đặt ảnh vào thư mục tương ứng:
```
public/images/brands/samsung.png
public/images/brands/apple.png
public/images/categories/smartphone.jpg
public/images/products/iphone-15-pro.jpg
```

### 2. Truy cập ảnh qua API
Server sẽ serve các ảnh tại endpoint `/images`:

- Brand logo: `http://localhost:8080/images/brands/samsung.png`
- Category image: `http://localhost:8080/images/categories/smartphone.jpg`
- Product image: `http://localhost:8080/images/products/iphone-15-pro.jpg`

### 3. Sử dụng trong Database
Khi tạo brand, category hoặc product, chỉ cần lưu đường dẫn tương đối:

```json
{
  "name": "Samsung",
  "logo": "images/brands/samsung.png"
}
```

hoặc

```json
{
  "name": "Samsung",
  "logo": "/images/brands/samsung.png"
}
```

### 4. Sử dụng trong Frontend
Frontend sẽ tự động xử lý và hiển thị ảnh đúng cách.

## Lưu ý

- Đặt tên file không dấu, không khoảng trắng (dùng `-` hoặc `_`)
- Format được khuyến nghị: PNG cho logo, JPG cho ảnh sản phẩm
- Nên resize ảnh trước khi upload để tối ưu hiệu suất
