import FeaturedProducts from '@/components/home/FeaturedProducts'
import HeroSection from '@/components/home/HeroSection'
import BrandList from '@/components/home/BrandList'
import CategoryList from '@/components/home/CategoryList'

export default function Home() {
  return (
    <div>
      <HeroSection />
      <div className="container mx-auto px-4 py-8 space-y-12">
        <CategoryList />
        <FeaturedProducts />
        <BrandList />
      </div>
    </div>
  )
}
