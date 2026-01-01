import Link from 'next/link'

export default function HeroSection() {
  return (
    <div className="bg-gradient-to-r from-primary-600 to-primary-800 text-white">
      <div className="container mx-auto px-4 py-16 md:py-24">
        <div className="max-w-3xl">
          <h1 className="text-4xl md:text-6xl font-bold mb-6">
            Điện thoại chính hãng
            <br />
            Giá tốt nhất thị trường
          </h1>
          <p className="text-lg md:text-xl mb-8 text-gray-100">
            Mua sắm điện thoại thông minh, tablet, laptop với ưu đãi hấp dẫn.
            Bảo hành chính hãng, giao hàng toàn quốc.
          </p>
          <div className="flex flex-col sm:flex-row space-y-4 sm:space-y-0 sm:space-x-4">
            <Link
              href="/products"
              className="px-8 py-3 bg-white text-primary-600 rounded-lg font-semibold hover:bg-gray-100 text-center"
            >
              Mua sắm ngay
            </Link>
            <Link
              href="/about"
              className="px-8 py-3 bg-transparent border-2 border-white text-white rounded-lg font-semibold hover:bg-white hover:text-primary-600 text-center"
            >
              Tìm hiểu thêm
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}
