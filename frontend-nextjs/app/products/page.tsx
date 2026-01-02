'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { useSearchParams } from 'next/navigation'
import { productService } from '@/services/productService'
import { Product, Brand, Category } from '@/types'
import toast from 'react-hot-toast'

export default function ProductsPage() {
  const searchParams = useSearchParams()
  const [products, setProducts] = useState<Product[]>([])
  const [brands, setBrands] = useState<Brand[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [loading, setLoading] = useState(true)
  
  // Filters
  const [selectedBrands, setSelectedBrands] = useState<string[]>([])
  const [selectedCategories, setSelectedCategories] = useState<string[]>([])
  const [priceRange, setPriceRange] = useState<string>('all')
  const [ramFilter, setRamFilter] = useState<string[]>([])
  const [sortBy, setSortBy] = useState('newest')
  
  // Pagination
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)

  useEffect(() => {
    fetchData()
  }, [selectedBrands, selectedCategories, priceRange, ramFilter, sortBy, page])

  const fetchData = async () => {
    setLoading(true)
    try {
      const [productsData, brandsData, categoriesData] = await Promise.all([
        productService.getProducts({
          page,
          limit: 12,
          brand: selectedBrands.join(','),
          category: selectedCategories.join(','),
          sort: sortBy
        }),
        productService.getBrands(),
        productService.getCategories()
      ])
      
      setProducts(productsData.data || [])
      setTotalPages(productsData.totalPages || 1)
      setBrands(brandsData || [])
      setCategories(categoriesData || [])
    } catch (error) {
      toast.error('Không thể tải dữ liệu')
    } finally {
      setLoading(false)
    }
  }

  const toggleBrand = (brandSlug: string) => {
    setSelectedBrands(prev =>
      prev.includes(brandSlug)
        ? prev.filter(b => b !== brandSlug)
        : [...prev, brandSlug]
    )
    setPage(1)
  }

  const toggleCategory = (categorySlug: string) => {
    setSelectedCategories(prev =>
      prev.includes(categorySlug)
        ? prev.filter(c => c !== categorySlug)
        : [...prev, categorySlug]
    )
    setPage(1)
  }

  const toggleRam = (ram: string) => {
    setRamFilter(prev =>
      prev.includes(ram)
        ? prev.filter(r => r !== ram)
        : [...prev, ram]
    )
    setPage(1)
  }

  const clearFilters = () => {
    setSelectedBrands([])
    setSelectedCategories([])
    setPriceRange('all')
    setRamFilter([])
    setPage(1)
  }

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND'
    }).format(price)
  }

  const activeFiltersCount = selectedBrands.length + selectedCategories.length + ramFilter.length + (priceRange !== 'all' ? 1 : 0)

  return (
    <div className="min-h-screen bg-gray-50 py-6">
      <div className="max-w-7xl mx-auto px-4">
        <h1 className="text-3xl font-bold mb-6">Điện thoại di động</h1>

        {/* Active Filters */}
        {activeFiltersCount > 0 && (
          <div className="mb-4 flex flex-wrap items-center gap-2">
            <span className="text-sm text-gray-600">Đang lọc:</span>
            {selectedBrands.map(brand => {
              const brandObj = brands.find(b => b.slug === brand)
              return (
                <span key={brand} className="inline-flex items-center gap-1 px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm">
                  Thương hiệu: {brandObj?.name}
                  <button onClick={() => toggleBrand(brand)} className="hover:text-blue-900">×</button>
                </span>
              )
            })}
            {ramFilter.map(ram => (
              <span key={ram} className="inline-flex items-center gap-1 px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm">
                RAM: {ram}
                <button onClick={() => toggleRam(ram)} className="hover:text-blue-900">×</button>
              </span>
            ))}
            <button onClick={clearFilters} className="text-sm text-blue-600 hover:underline">
              Xóa tất cả
            </button>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Sidebar Filters */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg shadow p-4 sticky top-4">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-bold">Bộ lọc</h2>
                {activeFiltersCount > 0 && (
                  <button onClick={clearFilters} className="text-sm text-blue-600 hover:underline">
                    Xóa tất cả
                  </button>
                )}
              </div>

              {/* Price Range Filter */}
              <div className="mb-6">
                <h3 className="font-semibold mb-3">Khoảng giá</h3>
                <div className="space-y-2">
                  {[
                    { value: 'all', label: 'Tất cả' },
                    { value: 'under-2m', label: 'Dưới 2 triệu' },
                    { value: '2m-4m', label: 'Từ 2 - 4 triệu' },
                    { value: '4m-7m', label: 'Từ 4 - 7 triệu' },
                    { value: '7m-13m', label: 'Từ 7 - 13 triệu' },
                    { value: 'over-13m', label: 'Trên 13 triệu' }
                  ].map(range => (
                    <label key={range.value} className="flex items-center cursor-pointer">
                      <input
                        type="radio"
                        name="price"
                        value={range.value}
                        checked={priceRange === range.value}
                        onChange={(e) => { setPriceRange(e.target.value); setPage(1) }}
                        className="mr-2"
                      />
                      <span className="text-sm">{range.label}</span>
                    </label>
                  ))}
                </div>
              </div>

              {/* Brand Filter */}
              <div className="mb-6">
                <h3 className="font-semibold mb-3">Thương hiệu</h3>
                <div className="space-y-2">
                  {brands.map(brand => (
                    <label key={brand.id} className="flex items-center cursor-pointer">
                      <input
                        type="checkbox"
                        checked={selectedBrands.includes(brand.slug)}
                        onChange={() => toggleBrand(brand.slug)}
                        className="mr-2"
                      />
                      <span className="text-sm">{brand.name}</span>
                    </label>
                  ))}
                </div>
              </div>

              {/* RAM Filter */}
              <div className="mb-6">
                <h3 className="font-semibold mb-3">Dung lượng RAM</h3>
                <div className="space-y-2">
                  {['4 GB', '6 GB', '8 GB', '12 GB'].map(ram => (
                    <label key={ram} className="flex items-center cursor-pointer">
                      <input
                        type="checkbox"
                        checked={ramFilter.includes(ram)}
                        onChange={() => toggleRam(ram)}
                        className="mr-2"
                      />
                      <span className="text-sm">{ram}</span>
                    </label>
                  ))}
                </div>
              </div>

              <button
                onClick={clearFilters}
                className="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 font-semibold"
              >
                Áp dụng bộ lọc
              </button>
            </div>
          </div>

          {/* Products Grid */}
          <div className="lg:col-span-3">
            {/* Sort Options */}
            <div className="flex items-center justify-between mb-4 bg-white rounded-lg shadow p-4">
              <span className="text-sm text-gray-600">
                Hiển thị {products.length} sản phẩm
              </span>
              <div className="flex items-center gap-2">
                <span className="text-sm text-gray-600">Sắp xếp:</span>
                <select
                  value={sortBy}
                  onChange={(e) => { setSortBy(e.target.value); setPage(1) }}
                  className="border border-gray-300 rounded px-3 py-1 text-sm focus:outline-none focus:border-blue-500"
                >
                  <option value="newest">Mới nhất</option>
                  <option value="price-asc">Giá thấp đến cao</option>
                  <option value="price-desc">Giá cao đến thấp</option>
                  <option value="name">Tên A-Z</option>
                </select>
              </div>
            </div>

            {loading ? (
              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                {[...Array(9)].map((_, i) => (
                  <div key={i} className="bg-white rounded-lg shadow p-4 animate-pulse">
                    <div className="bg-gray-200 aspect-square rounded mb-3"></div>
                    <div className="bg-gray-200 h-4 rounded mb-2"></div>
                    <div className="bg-gray-200 h-4 rounded w-2/3"></div>
                  </div>
                ))}
              </div>
            ) : products.length === 0 ? (
              <div className="bg-white rounded-lg shadow p-8 text-center">
                <p className="text-gray-600">Không tìm thấy sản phẩm nào</p>
                <button onClick={clearFilters} className="mt-4 text-blue-600 hover:underline">
                  Xóa bộ lọc
                </button>
              </div>
            ) : (
              <>
                <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                  {products.map((product) => (
                    <Link
                      key={product.id}
                      href={`/products/${product.slug}`}
                      className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-4 relative"
                    >
                      {/* Badges */}
                      {product.discount && product.discount > 0 && (
                        <span className="absolute top-2 left-2 bg-red-600 text-white text-xs px-2 py-1 rounded">
                          Trả góp 0%
                        </span>
                      )}
                      {product.isNew && (
                        <span className="absolute top-2 right-2 bg-green-500 text-white text-xs px-2 py-1 rounded">
                          Mới
                        </span>
                      )}

                      <div className="relative aspect-square mb-3">
                        <Image
                          src={product.images[0] ? (product.images[0].startsWith('http') ? product.images[0] : `/${product.images[0].replace(/^\//, '')}`) : '/placeholder.png'}
                          alt={product.name}
                          fill
                          className="object-cover rounded"
                        />
                      </div>

                      <h3 className="font-semibold text-sm mb-2 line-clamp-2 min-h-[40px]">
                        {product.name}
                      </h3>

                      <div className="flex items-baseline gap-2 mb-2">
                        <span className="text-red-600 font-bold text-lg">
                          {formatPrice(product.minPrice)}
                        </span>
                        {product.originalPrice && product.originalPrice > product.minPrice && (
                          <span className="text-gray-400 line-through text-sm">
                            {formatPrice(product.originalPrice)}
                          </span>
                        )}
                      </div>

                      {product.minPrice !== product.maxPrice && (
                        <p className="text-xs text-gray-500">
                          - {formatPrice(product.maxPrice)}
                        </p>
                      )}
                    </Link>
                  ))}
                </div>

                {/* Pagination */}
                {totalPages > 1 && (
                  <div className="flex justify-center items-center gap-2 mt-8">
                    <button
                      onClick={() => setPage(p => Math.max(1, p - 1))}
                      disabled={page === 1}
                      className="px-3 py-2 border rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      ‹
                    </button>

                    {[...Array(totalPages)].map((_, i) => {
                      const pageNum = i + 1
                      if (
                        pageNum === 1 ||
                        pageNum === totalPages ||
                        (pageNum >= page - 1 && pageNum <= page + 1)
                      ) {
                        return (
                          <button
                            key={pageNum}
                            onClick={() => setPage(pageNum)}
                            className={`px-4 py-2 rounded ${
                              page === pageNum
                                ? 'bg-blue-600 text-white'
                                : 'border hover:bg-gray-100'
                            }`}
                          >
                            {pageNum}
                          </button>
                        )
                      } else if (pageNum === page - 2 || pageNum === page + 2) {
                        return <span key={pageNum}>...</span>
                      }
                      return null
                    })}

                    <button
                      onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                      disabled={page === totalPages}
                      className="px-3 py-2 border rounded hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      ›
                    </button>
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
