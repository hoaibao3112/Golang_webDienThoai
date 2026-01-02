'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { useRouter } from 'next/navigation'
import { productService } from '@/services/productService'
import { cartService } from '@/services/cartService'
import { useAuthStore } from '@/store/authStore'
import { Product } from '@/types'
import { formatPrice } from '@/lib/utils'
import toast from 'react-hot-toast'

export default function FeaturedProducts() {
  const router = useRouter()
  const { isAuthenticated } = useAuthStore()
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [addingToCart, setAddingToCart] = useState<string | null>(null)

  useEffect(() => {
    fetchProducts()
  }, [])

  const fetchProducts = async () => {
    try {
      const data = await productService.getProducts({ limit: 8, sort: 'newest' })
      setProducts(data.data)
    } catch (error) {
      toast.error('Không thể tải sản phẩm')
    } finally {
      setLoading(false)
    }
  }

  const handleAddToCart = async (e: React.MouseEvent, product: Product) => {
    e.preventDefault()
    e.stopPropagation()

    if (!isAuthenticated) {
      toast.error('Vui lòng đăng nhập để thêm sản phẩm vào giỏ hàng')
      router.push('/login')
      return
    }

    // Check if product has variants
    if (!product.variants || product.variants.length === 0) {
      toast.error('Sản phẩm chưa có phiên bản để bán')
      return
    }

    // Find first available variant with stock > 0
    const availableVariant = product.variants.find(v => v.stock > 0)
    
    if (!availableVariant) {
      toast.error('Sản phẩm đã hết hàng')
      return
    }

    setAddingToCart(product.id)
    try {
      await cartService.addItem({ variantId: availableVariant.id, quantity: 1 })
      toast.success('Đã thêm vào giỏ hàng')
    } catch (error: any) {
      if (error.message?.includes('đăng nhập')) {
        toast.error('Phiên đăng nhập hết hạn, vui lòng đăng nhập lại')
        router.push('/login')
      } else {
        toast.error(error.message || 'Không thể thêm vào giỏ hàng')
      }
    } finally {
      setAddingToCart(null)
    }
  }

  if (loading) {
    return (
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {[...Array(8)].map((_, i) => (
          <div key={i} className="animate-pulse">
            <div className="bg-gray-200 h-48 rounded-lg mb-2"></div>
            <div className="bg-gray-200 h-4 rounded mb-2"></div>
            <div className="bg-gray-200 h-4 rounded w-2/3"></div>
          </div>
        ))}
      </div>
    )
  }

  return (
    <section>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">Sản phẩm nổi bật</h2>
        <Link href="/products" className="text-primary-600 hover:underline">
          Xem tất cả →
        </Link>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {products.map((product) => (
          <div
            key={product.id}
            className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-4 group"
          >
            <Link href={`/products/${product.slug}`} className="block">
              <div className="relative aspect-square mb-3">
                <Image
                  src={(product.images[0] && product.images[0].startsWith('http')) ? product.images[0] : (product.images[0] ? `/${product.images[0].replace(/^\//, '')}` : '/placeholder.png')}
                  alt={product.name}
                  fill
                  className="object-cover rounded-lg"
                />
              </div>
              <h3 className="font-semibold text-sm md:text-base mb-2 line-clamp-2">
                {product.name}
              </h3>
              <div className="flex items-baseline space-x-2 mb-3">
                <span className="text-red-600 font-bold">
                  {formatPrice(product.minPrice)}
                </span>
                {product.minPrice !== product.maxPrice && (
                  <span className="text-xs text-gray-500">
                    - {formatPrice(product.maxPrice)}
                  </span>
                )}
              </div>
            </Link>
            <button
              onClick={(e) => handleAddToCart(e, product)}
              disabled={addingToCart === product.id}
              className="w-full bg-blue-600 text-white py-2 rounded-lg text-sm font-semibold hover:bg-blue-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed"
            >
              {addingToCart === product.id ? 'Đang thêm...' : 'Thêm vào giỏ hàng'}
            </button>
          </div>
        ))}
      </div>
    </section>
  )
}
