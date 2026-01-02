'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { useRouter } from 'next/navigation'
import { cartService } from '@/services/cartService'
import { orderService } from '@/services/orderService'
import toast from 'react-hot-toast'

interface CartItem {
  variantId: string
  productName: string
  variant: string
  color: string
  storage: string
  price: number
  quantity: number
  image: string
}

export default function CheckoutPage() {
  const router = useRouter()
  const [items, setItems] = useState<CartItem[]>([])
  const [loading, setLoading] = useState(false)
  
  // Form data
  const [fullName, setFullName] = useState('')
  const [phone, setPhone] = useState('')
  const [address, setAddress] = useState('')
  const [city, setCity] = useState('')
  const [district, setDistrict] = useState('')
  const [ward, setWard] = useState('')
  const [note, setNote] = useState('')
  const [paymentMethod, setPaymentMethod] = useState('COD')
  const [voucherCode, setVoucherCode] = useState('')
  const [discount, setDiscount] = useState(0)

  useEffect(() => {
    fetchCart()
  }, [])

  const fetchCart = async () => {
    try {
      const data = await cartService.getCart()
      const mappedItems = (data.items || []).map((item: any) => ({
        ...item,
        variant: item.variant || `${item.storage} ${item.color}`.trim()
      }))
      setItems(mappedItems)
      if (mappedItems.length === 0) {
        toast.error('Giỏ hàng trống')
        router.push('/cart')
      }
    } catch (error: any) {
      if (error.message?.includes('401') || error.message?.includes('Unauthorized')) {
        toast.error('Vui lòng đăng nhập để thanh toán')
        setTimeout(() => router.push('/auth/login'), 1500)
      } else {
        toast.error('Không thể tải giỏ hàng')
        router.push('/cart')
      }
    }
  }

  const applyVoucher = () => {
    if (voucherCode.trim()) {
      // Simulate voucher validation
      setDiscount(500000)
      toast.success('Đã áp dụng mã giảm giá')
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!fullName || !phone || !address || !city || !district || !ward) {
      toast.error('Vui lòng điền đầy đủ thông tin giao hàng')
      return
    }

    setLoading(true)
    try {
      await orderService.createOrder({
        shippingAddress: {
          fullName,
          phone,
          address,
          city,
          district,
          ward
        },
        voucherCode: voucherCode || undefined
      })
      
      toast.success('Đặt hàng thành công!')
      router.push('/orders')
    } catch (error: any) {
      toast.error(error.message || 'Đặt hàng thất bại')
    } finally {
      setLoading(false)
    }
  }

  const subtotal = items.reduce((sum, item) => sum + item.price * item.quantity, 0)
  const shippingFee = 0 // Miễn phí
  const total = subtotal + shippingFee - discount

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND'
    }).format(price)
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-6xl mx-auto px-4">
        {/* Breadcrumb */}
        <div className="text-sm text-gray-600 mb-4">
          <Link href="/" className="hover:text-blue-600">Trang chủ</Link>
          <span className="mx-2">/</span>
          <Link href="/cart" className="hover:text-blue-600">Giỏ hàng</Link>
          <span className="mx-2">/</span>
          <span>Thanh toán</span>
        </div>

        <h1 className="text-3xl font-bold mb-8">Thanh toán</h1>

        <form onSubmit={handleSubmit}>
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Left Column - Shipping Info & Payment */}
            <div className="lg:col-span-2 space-y-6">
              {/* Shipping Information */}
              <div className="bg-white rounded-lg shadow p-6">
                <h2 className="text-xl font-bold mb-4">1. Thông tin giao hàng</h2>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="md:col-span-1">
                    <label className="block text-sm font-medium mb-2">Họ và tên</label>
                    <input
                      type="text"
                      value={fullName}
                      onChange={(e) => setFullName(e.target.value)}
                      placeholder="Nhập họ và tên"
                      className="w-full border border-gray-300 rounded px-4 py-2 focus:outline-none focus:border-blue-500"
                      required
                    />
                  </div>

                  <div className="md:col-span-1">
                    <label className="block text-sm font-medium mb-2">Số điện thoại</label>
                    <input
                      type="tel"
                      value={phone}
                      onChange={(e) => setPhone(e.target.value)}
                      placeholder="Nhập số điện thoại"
                      className="w-full border border-gray-300 rounded px-4 py-2 focus:outline-none focus:border-blue-500"
                      required
                    />
                  </div>

                  <div className="md:col-span-2">
                    <label className="block text-sm font-medium mb-2">Địa chỉ</label>
                    <input
                      type="text"
                      value={address}
                      onChange={(e) => setAddress(e.target.value)}
                      placeholder="Số nhà, tên đường, phường/xã, quận/huyện, tỉnh/thành"
                      className="w-full border border-gray-300 rounded px-4 py-2 focus:outline-none focus:border-blue-500"
                      required
                    />
                  </div>

                  <div className="md:col-span-1">
                    <label className="block text-sm font-medium mb-2">Tỉnh/Thành phố</label>
                    <input
                      type="text"
                      value={city}
                      onChange={(e) => setCity(e.target.value)}
                      placeholder="Nhập tỉnh/thành phố"
                      className="w-full border border-gray-300 rounded px-4 py-2 focus:outline-none focus:border-blue-500"
                      required
                    />
                  </div>

                  <div className="md:col-span-1">
                    <label className="block text-sm font-medium mb-2">Quận/Huyện</label>
                    <input
                      type="text"
                      value={district}
                      onChange={(e) => setDistrict(e.target.value)}
                      placeholder="Nhập quận/huyện"
                      className="w-full border border-gray-300 rounded px-4 py-2 focus:outline-none focus:border-blue-500"
                      required
                    />
                  </div>

                  <div className="md:col-span-2">
                    <label className="block text-sm font-medium mb-2">Phường/Xã</label>
                    <input
                      type="text"
                      value={ward}
                      onChange={(e) => setWard(e.target.value)}
                      placeholder="Nhập phường/xã"
                      className="w-full border border-gray-300 rounded px-4 py-2 focus:outline-none focus:border-blue-500"
                      required
                    />
                  </div>

                  <div className="md:col-span-2">
                    <label className="block text-sm font-medium mb-2">Ghi chú (Tùy chọn)</label>
                    <textarea
                      value={note}
                      onChange={(e) => setNote(e.target.value)}
                      placeholder="Ghi chú cho người giao hàng"
                      rows={3}
                      className="w-full border border-gray-300 rounded px-4 py-2 focus:outline-none focus:border-blue-500"
                    />
                  </div>
                </div>
              </div>

              {/* Payment Method */}
              <div className="bg-white rounded-lg shadow p-6">
                <h2 className="text-xl font-bold mb-4">2. Phương thức thanh toán</h2>
                
                <div className="space-y-3">
                  <label className="flex items-start p-4 border-2 border-blue-500 rounded-lg cursor-pointer bg-blue-50">
                    <input
                      type="radio"
                      name="payment"
                      value="COD"
                      checked={paymentMethod === 'COD'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                      className="mt-1 mr-3"
                    />
                    <div>
                      <div className="font-semibold">Thanh toán khi nhận hàng (COD)</div>
                      <div className="text-sm text-gray-600">Thanh toán bằng tiền mặt khi shipper giao hàng.</div>
                    </div>
                  </label>

                  <label className="flex items-start p-4 border-2 border-gray-200 rounded-lg cursor-pointer hover:border-gray-300">
                    <input
                      type="radio"
                      name="payment"
                      value="CARD"
                      checked={paymentMethod === 'CARD'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                      className="mt-1 mr-3"
                    />
                    <div>
                      <div className="font-semibold">Thẻ ATM / Visa / Mastercard</div>
                      <div className="text-sm text-gray-600">Thanh toán an toàn qua các thẻ thanh toán trực tuyến.</div>
                    </div>
                  </label>

                  <label className="flex items-start p-4 border-2 border-gray-200 rounded-lg cursor-pointer hover:border-gray-300">
                    <input
                      type="radio"
                      name="payment"
                      value="MOMO"
                      checked={paymentMethod === 'MOMO'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                      className="mt-1 mr-3"
                    />
                    <div className="flex items-center">
                      <div className="w-8 h-8 bg-pink-600 rounded mr-3 flex items-center justify-center text-white font-bold">M</div>
                      <div>
                        <div className="font-semibold">Ví điện tử MoMo</div>
                        <div className="text-sm text-gray-600">Quét mã QR để thanh toán qua ứng dụng MoMo.</div>
                      </div>
                    </div>
                  </label>

                  <label className="flex items-start p-4 border-2 border-gray-200 rounded-lg cursor-pointer hover:border-gray-300">
                    <input
                      type="radio"
                      name="payment"
                      value="VNPAY"
                      checked={paymentMethod === 'VNPAY'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                      className="mt-1 mr-3"
                    />
                    <div className="flex items-center">
                      <div className="w-8 h-8 bg-blue-700 rounded mr-3 flex items-center justify-center text-white font-bold text-xs">VP</div>
                      <div>
                        <div className="font-semibold">VNPAY-QR</div>
                        <div className="text-sm text-gray-600">Thanh toán qua ứng dụng hỗ trợ VNPAY QR và ví điện tử hỗ trợ VNPAY.</div>
                      </div>
                    </div>
                  </label>
                </div>
              </div>
            </div>

            {/* Right Column - Order Summary */}
            <div className="lg:col-span-1">
              <div className="bg-white rounded-lg shadow p-6 sticky top-4">
                <h2 className="text-xl font-bold mb-4">Đơn hàng của bạn</h2>

                {/* Items List */}
                <div className="space-y-3 mb-4 max-h-60 overflow-y-auto">
                  {items.map((item) => (
                    <div key={item.variantId} className="flex gap-3">
                      <div className="relative w-16 h-16 flex-shrink-0">
                        <Image
                          src={item.image ? (item.image.startsWith('http') ? item.image : `/${item.image.replace(/^\//, '')}`) : '/placeholder.png'}
                          alt={item.productName}
                          fill
                          className="object-cover rounded"
                        />
                      </div>
                      <div className="flex-1">
                        <h4 className="font-medium text-sm line-clamp-1">{item.productName}</h4>
                        <p className="text-xs text-gray-600">{item.storage} | {item.color}</p>
                        <p className="text-sm">Số lượng: {item.quantity}</p>
                      </div>
                      <div className="text-right">
                        <p className="font-semibold">{formatPrice(item.price)}</p>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Price Summary */}
                <div className="border-t pt-4 space-y-3 mb-4">
                  <div className="flex justify-between">
                    <span className="text-gray-600">Tạm tính</span>
                    <span className="font-semibold">{formatPrice(subtotal)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">Phí vận chuyển</span>
                    <span className="text-green-600 font-semibold">Miễn phí</span>
                  </div>
                  {discount > 0 && (
                    <div className="flex justify-between text-green-600">
                      <span>Giảm giá</span>
                      <span className="font-semibold">- {formatPrice(discount)}</span>
                    </div>
                  )}
                  <div className="border-t pt-3 flex justify-between">
                    <span className="text-lg font-bold">Tổng cộng</span>
                    <span className="text-xl font-bold text-blue-600">{formatPrice(total)}</span>
                  </div>
                </div>

                <button
                  type="submit"
                  disabled={loading}
                  className="w-full bg-blue-600 text-white py-3 rounded-lg font-semibold hover:bg-blue-700 transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed mb-4"
                >
                  {loading ? 'ĐANG XỬ LÝ...' : 'HOÀN TẤT ĐƠN HÀNG'}
                </button>

                {/* Voucher */}
                <div className="border-t pt-4">
                  <h3 className="font-semibold mb-2">Mã giảm giá</h3>
                  <div className="flex gap-2">
                    <input
                      type="text"
                      placeholder="Nhập mã giảm giá"
                      value={voucherCode}
                      onChange={(e) => setVoucherCode(e.target.value)}
                      className="flex-1 border border-gray-300 rounded px-3 py-2 text-sm"
                    />
                    <button
                      type="button"
                      onClick={applyVoucher}
                      className="px-4 py-2 bg-gray-100 text-gray-700 rounded hover:bg-gray-200 font-semibold text-sm"
                    >
                      Áp dụng
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  )
}
