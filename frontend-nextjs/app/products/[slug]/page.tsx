'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { useParams, useRouter } from 'next/navigation'
import { productService } from '@/services/productService'
import { cartService } from '@/services/cartService'
import { ProductDetail, ProductVariant } from '@/types'
import toast from 'react-hot-toast'

export default function ProductDetailPage() {
  const params = useParams()
  const router = useRouter()
  const slug = params.slug as string

  const [product, setProduct] = useState<ProductDetail | null>(null)
  const [selectedVariant, setSelectedVariant] = useState<ProductVariant | null>(null)
  const [selectedImage, setSelectedImage] = useState(0)
  const [quantity, setQuantity] = useState(1)
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState('description')

  useEffect(() => {
    fetchProduct()
  }, [slug])

  const fetchProduct = async () => {
    try {
      const data = await productService.getProductBySlug(slug)
      setProduct(data)
      if (data.variants && data.variants.length > 0) {
        setSelectedVariant(data.variants[0])
      }
    } catch (error) {
      toast.error('Không thể tải thông tin sản phẩm')
      router.push('/products')
    } finally {
      setLoading(false)
    }
  }

  const handleAddToCart = async () => {
    if (!selectedVariant) {
      toast.error('Vui lòng chọn phiên bản sản phẩm')
      return
    }

    try {
      await cartService.addItem({
        variantId: selectedVariant.id,
        quantity
      })
      toast.success('Đã thêm vào giỏ hàng')
    } catch (error: any) {
      toast.error(error.message || 'Không thể thêm vào giỏ hàng')
    }
  }

  const handleBuyNow = async () => {
    if (!selectedVariant) {
      toast.error('Vui lòng chọn phiên bản sản phẩm')
      return
    }

    try {
      await cartService.addItem({
        variantId: selectedVariant.id,
        quantity
      })
      router.push('/checkout')
    } catch (error: any) {
      toast.error(error.message || 'Có lỗi xảy ra')
    }
  }

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND'
    }).format(price)
  }

  const getUniqueStorages = () => {
    if (!product?.variants) return []
    const storages = Array.from(new Set(product.variants.map(v => v.storage)))
    return storages
  }

  const getUniqueColors = () => {
    if (!product?.variants) return []
    const colors = Array.from(new Set(product.variants.map(v => v.color)))
    return colors
  }

  const getVariantByStorageAndColor = (storage: string, color: string) => {
    return product?.variants.find(v => v.storage === storage && v.color === color)
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (!product) return null

  const currentStorage = selectedVariant?.storage || ''
  const currentColor = selectedVariant?.color || ''

  return (
    <div className="min-h-screen bg-gray-50 py-6">
      <div className="max-w-7xl mx-auto px-4">
        {/* Breadcrumb */}
        <div className="text-sm text-gray-600 mb-4">
          <Link href="/" className="hover:text-blue-600">Trang chủ</Link>
          <span className="mx-2">/</span>
          <Link href="/products" className="hover:text-blue-600">Điện thoại</Link>
          <span className="mx-2">/</span>
          <span>{product.brand.name}</span>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
          {/* Left - Images */}
          <div>
            <div className="bg-white rounded-lg p-4 mb-4">
              <div className="relative aspect-square mb-4">
                <Image
                  src={product.images[selectedImage] ? (product.images[selectedImage].startsWith('http') ? product.images[selectedImage] : `/${product.images[selectedImage].replace(/^\//, '')}`) : '/placeholder.png'}
                  alt={product.name}
                  fill
                  className="object-contain"
                />
              </div>
            </div>

            {/* Thumbnails */}
            <div className="flex gap-2 overflow-x-auto">
              {product.images.map((image, index) => (
                <button
                  key={index}
                  onClick={() => setSelectedImage(index)}
                  className={`relative w-20 h-20 flex-shrink-0 rounded border-2 ${
                    selectedImage === index ? 'border-blue-600' : 'border-gray-200'
                  }`}
                >
                  <Image
                    src={image ? (image.startsWith('http') ? image : `/${image.replace(/^\//, '')}`) : '/placeholder.png'}
                    alt={`${product.name} ${index + 1}`}
                    fill
                    className="object-cover rounded"
                  />
                </button>
              ))}
            </div>
          </div>

          {/* Right - Product Info */}
          <div>
            <div className="bg-white rounded-lg p-6">
              <h1 className="text-2xl font-bold mb-2">{product.name}</h1>
              
              {/* Rating */}
              <div className="flex items-center gap-2 mb-4">
                <div className="flex text-yellow-400">
                  {[...Array(5)].map((_, i) => (
                    <svg key={i} className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                  ))}
                </div>
                <span className="text-sm text-gray-600">(101 đánh giá)</span>
              </div>

              {/* Price */}
              {selectedVariant && (
                <div className="mb-6">
                  <div className="flex items-baseline gap-3">
                    <span className="text-3xl font-bold text-red-600">
                      {formatPrice(selectedVariant.price)}
                    </span>
                    {product.minPrice !== product.maxPrice && (
                      <span className="text-lg text-gray-400 line-through">
                        {formatPrice(selectedVariant.price * 1.1)}
                      </span>
                    )}
                    <span className="bg-red-100 text-red-600 text-sm px-2 py-1 rounded">
                      -10%
                    </span>
                  </div>
                  <p className="text-sm text-gray-500 mt-1">Giá chưa bao gồm VAT</p>
                </div>
              )}

              {/* Storage Options */}
              <div className="mb-6">
                <h3 className="font-semibold mb-3">Lựa chọn phiên bản</h3>
                <div className="grid grid-cols-3 gap-2">
                  {getUniqueStorages().map((storage) => {
                    const variant = getVariantByStorageAndColor(storage, currentColor)
                    return (
                      <button
                        key={storage}
                        onClick={() => variant && setSelectedVariant(variant)}
                        disabled={!variant}
                        className={`py-3 px-4 rounded-lg border-2 text-center ${
                          currentStorage === storage
                            ? 'border-blue-600 bg-blue-50 text-blue-600'
                            : 'border-gray-200 hover:border-gray-300'
                        } ${!variant ? 'opacity-50 cursor-not-allowed' : ''}`}
                      >
                        <div className="font-semibold">{storage}</div>
                      </button>
                    )
                  })}
                </div>
              </div>

              {/* Color Options */}
              <div className="mb-6">
                <h3 className="font-semibold mb-3">Lựa chọn màu sắc</h3>
                <div className="grid grid-cols-3 gap-2">
                  {getUniqueColors().map((color) => {
                    const variant = getVariantByStorageAndColor(currentStorage, color)
                    return (
                      <button
                        key={color}
                        onClick={() => variant && setSelectedVariant(variant)}
                        disabled={!variant}
                        className={`py-3 px-4 rounded-lg border-2 text-center ${
                          currentColor === color
                            ? 'border-blue-600 bg-blue-50 text-blue-600'
                            : 'border-gray-200 hover:border-gray-300'
                        } ${!variant ? 'opacity-50 cursor-not-allowed' : ''}`}
                      >
                        <div className="font-semibold">{color}</div>
                      </button>
                    )
                  })}
                </div>
              </div>

              {/* Stock Info */}
              {selectedVariant && (
                <div className="mb-6">
                  <p className="text-sm text-gray-600">
                    Còn {selectedVariant.stock} sản phẩm
                  </p>
                </div>
              )}

              {/* Action Buttons */}
              <div className="flex gap-3 mb-6">
                <button
                  onClick={handleBuyNow}
                  className="flex-1 bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 transition"
                >
                  MUA NGAY
                </button>
                <button
                  onClick={handleAddToCart}
                  className="flex-1 border-2 border-blue-600 text-blue-600 py-3 rounded-lg font-semibold hover:bg-blue-50 transition flex items-center justify-center gap-2"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                  </svg>
                  Thêm vào giỏ
                </button>
              </div>

              {/* Promotions */}
              <div className="border-t pt-4">
                <h3 className="font-semibold mb-3">Khuyến mãi & Ưu đãi</h3>
                <div className="space-y-2 text-sm">
                  <div className="flex items-start gap-2">
                    <svg className="w-5 h-5 text-green-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span>Miễn phí giao hàng toàn quốc cho đơn hàng từ 500.000đ</span>
                  </div>
                  <div className="flex items-start gap-2">
                    <svg className="w-5 h-5 text-green-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span>Trả góp 0% - Duyệt hồ sơ chỉ 5 phút</span>
                  </div>
                  <div className="flex items-start gap-2">
                    <svg className="w-5 h-5 text-green-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span>Bảo hành 12 tháng chính hãng</span>
                  </div>
                  <div className="flex items-start gap-2">
                    <svg className="w-5 h-5 text-green-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                    <span>Đổi trả trong 7 ngày nếu có lỗi từ nhà sản xuất</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Tabs - Description & Reviews */}
        <div className="bg-white rounded-lg p-6 mb-8">
          <div className="border-b mb-6">
            <div className="flex gap-8">
              <button
                onClick={() => setActiveTab('description')}
                className={`pb-4 px-2 font-semibold border-b-2 transition ${
                  activeTab === 'description'
                    ? 'border-blue-600 text-blue-600'
                    : 'border-transparent text-gray-600 hover:text-gray-900'
                }`}
              >
                Mô tả sản phẩm
              </button>
              <button
                onClick={() => setActiveTab('specs')}
                className={`pb-4 px-2 font-semibold border-b-2 transition ${
                  activeTab === 'specs'
                    ? 'border-blue-600 text-blue-600'
                    : 'border-transparent text-gray-600 hover:text-gray-900'
                }`}
              >
                Thông số kỹ thuật
              </button>
              <button
                onClick={() => setActiveTab('reviews')}
                className={`pb-4 px-2 font-semibold border-b-2 transition ${
                  activeTab === 'reviews'
                    ? 'border-blue-600 text-blue-600'
                    : 'border-transparent text-gray-600 hover:text-gray-900'
                }`}
              >
                Đánh giá (101)
              </button>
            </div>
          </div>

          {activeTab === 'description' && (
            <div className="prose max-w-none">
              <h2 className="text-2xl font-bold mb-4">Thiết kế, Tiện điểm biến, siêu nhỏ</h2>
              <p className="text-gray-700 leading-relaxed mb-4">{product.description}</p>
              <p className="text-gray-700 leading-relaxed">
                {product.name} là một trong những sản phẩm flagship mới nhất của {product.brand.name}, 
                mang đến trải nghiệm tuyệt vời với thiết kế cao cấp, hiệu năng mạnh mẽ và nhiều tính năng hiện đại.
              </p>
            </div>
          )}

          {activeTab === 'specs' && (
            <div>
              <h3 className="font-semibold text-lg mb-4">Thông số kỹ thuật</h3>
              <div className="grid grid-cols-2 gap-4">
                <div className="border-b py-2">
                  <span className="text-gray-600">Thương hiệu:</span>
                  <span className="font-semibold ml-2">{product.brand.name}</span>
                </div>
                <div className="border-b py-2">
                  <span className="text-gray-600">Danh mục:</span>
                  <span className="font-semibold ml-2">{product.category.name}</span>
                </div>
                {selectedVariant && (
                  <>
                    <div className="border-b py-2">
                      <span className="text-gray-600">Bộ nhớ:</span>
                      <span className="font-semibold ml-2">{selectedVariant.storage}</span>
                    </div>
                    <div className="border-b py-2">
                      <span className="text-gray-600">Màu sắc:</span>
                      <span className="font-semibold ml-2">{selectedVariant.color}</span>
                    </div>
                  </>
                )}
              </div>
            </div>
          )}

          {activeTab === 'reviews' && (
            <div>
              <p className="text-gray-600">Chưa có đánh giá nào cho sản phẩm này.</p>
            </div>
          )}
        </div>

        {/* Related Products */}
        <div>
          <h2 className="text-2xl font-bold mb-6">Sản phẩm liên quan</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 gap-4">
            {/* Placeholder for related products */}
            <div className="bg-white rounded-lg p-4 text-center text-gray-500">
              Sản phẩm 1
            </div>
            <div className="bg-white rounded-lg p-4 text-center text-gray-500">
              Sản phẩm 2
            </div>
            <div className="bg-white rounded-lg p-4 text-center text-gray-500">
              Sản phẩm 3
            </div>
            <div className="bg-white rounded-lg p-4 text-center text-gray-500">
              Sản phẩm 4
            </div>
            <div className="bg-white rounded-lg p-4 text-center text-gray-500">
              Sản phẩm 5
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
