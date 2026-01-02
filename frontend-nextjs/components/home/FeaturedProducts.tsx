'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { productService } from '@/services/productService'
import { Product } from '@/types'
import { formatPrice } from '@/lib/utils'
import toast from 'react-hot-toast'

export default function FeaturedProducts() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)

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
          <Link
            key={product.id}
            href={`/products/${product.slug}`}
            className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-4"
          >
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
            <div className="flex items-baseline space-x-2">
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
        ))}
      </div>
    </section>
  )
}
