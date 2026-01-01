'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { productService } from '@/services/productService'
import { Brand } from '@/types'
import toast from 'react-hot-toast'

export default function BrandList() {
  const [brands, setBrands] = useState<Brand[]>([])

  useEffect(() => {
    fetchBrands()
  }, [])

  const fetchBrands = async () => {
    try {
      const data = await productService.getBrands()
      setBrands(data)
    } catch (error) {
      toast.error('Không thể tải thương hiệu')
    }
  }

  return (
    <section>
      <h2 className="text-2xl font-bold mb-6">Thương hiệu nổi bật</h2>
      <div className="grid grid-cols-3 md:grid-cols-6 gap-4">
        {brands.map((brand) => (
          <Link
            key={brand.id}
            href={`/products?brand=${brand.slug}`}
            className="bg-white rounded-lg shadow hover:shadow-lg transition-shadow p-4 flex items-center justify-center aspect-square"
          >
            <Image
              src={brand.logo || '/placeholder.png'}
              alt={brand.name}
              width={80}
              height={80}
              className="object-contain"
            />
          </Link>
        ))}
      </div>
    </section>
  )
}
