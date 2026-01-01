'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { productService } from '@/services/productService'
import { Category } from '@/types'
import toast from 'react-hot-toast'

export default function CategoryList() {
  const [categories, setCategories] = useState<Category[]>([])

  useEffect(() => {
    fetchCategories()
  }, [])

  const fetchCategories = async () => {
    try {
      const data = await productService.getCategories()
      setCategories(data)
    } catch (error) {
      toast.error('Không thể tải danh mục')
    }
  }

  return (
    <section>
      <h2 className="text-2xl font-bold mb-6">Danh mục sản phẩm</h2>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {categories.map((category) => (
          <Link
            key={category.id}
            href={`/products?category=${category.slug}`}
            className="relative h-32 md:h-48 rounded-lg overflow-hidden group"
          >
            <Image
              src={category.image || '/placeholder.png'}
              alt={category.name}
              fill
              className="object-cover group-hover:scale-110 transition-transform"
            />
            <div className="absolute inset-0 bg-black bg-opacity-40 flex items-center justify-center">
              <h3 className="text-white text-lg md:text-xl font-bold">
                {category.name}
              </h3>
            </div>
          </Link>
        ))}
      </div>
    </section>
  )
}
