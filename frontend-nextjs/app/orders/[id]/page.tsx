'use client'

import { useState, useEffect } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/store/authStore'
import { orderService } from '@/services/orderService'
import { Order } from '@/types'
import toast from 'react-hot-toast'

export default function OrderDetailPage({ params }: { params: { id: string } }) {
  const router = useRouter()
  const { isAuthenticated } = useAuthStore()
  const [order, setOrder] = useState<Order | null>(null)
  const [loading, setLoading] = useState(true)

  const statusIcons: { [key: string]: string } = {
    pending: '⏳',
    confirmed: '✅',
    processing: '📦',
    shipping: '🚚',
    delivered: '✅',
    cancelled: '❌',
  }

  const statusLabels: { [key: string]: string } = {
    pending: 'Chờ xác nhận',
    confirmed: 'Đã xác nhận',
    processing: 'Đang chuẩn bị hàng',
    shipping: 'Đang giao hàng',
    delivered: 'Đã giao hàng',
    cancelled: 'Đã hủy',
  }

  const statusColors: { [key: string]: string } = {
    pending: 'text-yellow-600',
    confirmed: 'text-blue-600',
    processing: 'text-blue-600',
    shipping: 'text-blue-600',
    delivered: 'text-green-600',
    cancelled: 'text-red-600',
  }

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
      return
    }
    fetchOrderDetail()
  }, [isAuthenticated, params.id, router])

  const fetchOrderDetail = async () => {
    try {
      const response = await orderService.getOrderById(params.id)
      setOrder(response)
    } catch (error: any) {
      if (error.message?.includes('đăng nhập')) {
        router.push('/login')
      } else {
        toast.error('Không thể tải thông tin đơn hàng')
      }
    } finally {
      setLoading(false)
    }
  }

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('vi-VN', {
      style: 'currency',
      currency: 'VND'
    }).format(price)
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString('vi-VN', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const handleBuyAgain = () => {
    toast('Tính năng mua lại đang được phát triển', { icon: 'ℹ️' })
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="max-w-5xl mx-auto px-4">
          <div className="animate-pulse">
            <div className="h-8 bg-gray-200 rounded w-64 mb-6"></div>
            <div className="bg-white rounded-lg shadow p-6">
              <div className="h-4 bg-gray-200 rounded w-full mb-4"></div>
              <div className="h-4 bg-gray-200 rounded w-3/4"></div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!order) {
    return (
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="max-w-5xl mx-auto px-4 text-center">
          <h1 className="text-2xl font-bold mb-4">Không tìm thấy đơn hàng</h1>
          <Link href="/orders" className="text-blue-600 hover:underline">
            ← Quay lại danh sách đơn hàng
          </Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-5xl mx-auto px-4">
        {/* Breadcrumb */}
        <div className="mb-4 text-sm text-gray-600">
          <Link href="/" className="hover:underline">Trang chủ</Link>
          <span className="mx-2">/</span>
          <Link href="/orders" className="hover:underline">Quản lý đơn hàng</Link>
          <span className="mx-2">/</span>
          <span>Đơn hàng #{order.order_number}</span>
        </div>

        {/* Page Header */}
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold">Chi Tiết Đơn Hàng #{order.order_number}</h1>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Left Column - Order Timeline & Items */}
          <div className="lg:col-span-2 space-y-6">
            {/* Order Status Timeline */}
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold mb-4">Lộ trình vận chuyển</h2>
              
              <div className="space-y-4">
                {order.status_history && order.status_history.length > 0 ? (
                  order.status_history.map((history, index) => (
                    <div key={index} className="flex items-start space-x-4">
                      <div
                        className={`flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center ${
                          index === 0 ? 'bg-green-100' : 'bg-gray-100'
                        }`}
                      >
                        <span className="text-xl">
                          {statusIcons[history.status] || '📋'}
                        </span>
                      </div>
                      <div className="flex-1">
                        <h3
                          className={`font-semibold ${
                            statusColors[history.status] || 'text-gray-800'
                          }`}
                        >
                          {statusLabels[history.status] || history.status}
                        </h3>
                        <p className="text-sm text-gray-600 mt-1">{history.note}</p>
                        <p className="text-xs text-gray-400 mt-1">
                          {formatDate(history.created_at || history.createdAt)}
                        </p>
                      </div>
                    </div>
                  ))
                ) : (
                  <div className="flex items-start space-x-4">
                    <div className="flex-shrink-0 w-10 h-10 rounded-full bg-yellow-100 flex items-center justify-center">
                      <span className="text-xl">⏳</span>
                    </div>
                    <div className="flex-1">
                      <h3 className="font-semibold text-yellow-600">
                        {statusLabels[order.status] || 'Đang xử lý'}
                      </h3>
                      <p className="text-sm text-gray-600 mt-1">
                        Đơn hàng đã được tạo thành công.
                      </p>
                      <p className="text-xs text-gray-400 mt-1">
                        {formatDate(order.created_at || order.createdAt)}
                      </p>
                    </div>
                  </div>
                )}
              </div>
            </div>

            {/* Order Items */}
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold mb-4">Sản phẩm đã đặt</h2>
              
              <div className="space-y-4">
                {order.items && order.items.length > 0 ? (
                  order.items.map((item, index) => (
                    <div key={index} className="flex items-center space-x-4 pb-4 border-b last:border-b-0">
                      <div className="relative w-24 h-24 flex-shrink-0">
                        <Image
                          src={
                            item.image && item.image.startsWith('http')
                              ? item.image
                              : item.image
                              ? `/${item.image.replace(/^\//, '')}`
                              : '/placeholder.png'
                          }
                          alt={item.productName || item.name}
                          fill
                          className="object-cover rounded-lg"
                        />
                      </div>
                      <div className="flex-1">
                        <h3 className="font-medium mb-1">{item.productName || item.name}</h3>
                        <p className="text-sm text-gray-500 mb-1">{item.variantName || `${item.storage} - ${item.color}`}</p>
                        <p className="text-sm text-gray-500">Số lượng: x{item.quantity}</p>
                      </div>
                      <div className="text-right">
                        <p className="font-semibold text-lg">
                          {formatPrice(item.price)}
                        </p>
                      </div>
                    </div>
                  ))
                ) : (
                  <p className="text-gray-500 text-center py-4">Không có sản phẩm</p>
                )}
              </div>

              <div className="mt-6 pt-4 border-t space-y-2">
                <div className="flex justify-between text-gray-600">
                  <span>Tạm tính:</span>
                  <span>{formatPrice(order.total_amount || order.totalAmount || order.total)}</span>
                </div>
                <div className="flex justify-between text-gray-600">
                  <span>Phí vận chuyển:</span>
                  <span>Miễn phí</span>
                </div>
                <div className="flex justify-between text-xl font-bold text-red-600">
                  <span>Tổng cộng:</span>
                  <span>{formatPrice(order.total_amount || order.totalAmount || order.total)}</span>
                </div>
              </div>
            </div>
          </div>

          {/* Right Column - Shipping & Payment Info */}
          <div className="space-y-6">
            {/* Shipping Information */}
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold mb-4">Thông tin vận chuyển</h2>
              <div className="space-y-3">
                <div className="flex items-start space-x-2">
                  <span className="text-xl">📦</span>
                  <div>
                    <h3 className="font-semibold">Giao Hàng Nhanh</h3>
                    <p className="text-sm text-gray-600">
                      Mã vận đơn: GHN{order.order_number}
                    </p>
                  </div>
                </div>
                
                <div className="pt-3 border-t">
                  <h4 className="font-medium mb-2">Địa chỉ nhận hàng:</h4>
                  {order.shipping_address || order.shippingAddress ? (
                    <>
                      <p className="text-sm font-medium">{(order.shipping_address || order.shippingAddress)?.fullName}</p>
                      <p className="text-sm text-gray-600">{(order.shipping_address || order.shippingAddress)?.phone}</p>
                      <p className="text-sm text-gray-600 mt-1">
                        {(order.shipping_address || order.shippingAddress)?.address}, {(order.shipping_address || order.shippingAddress)?.ward}, {(order.shipping_address || order.shippingAddress)?.district}, {(order.shipping_address || order.shippingAddress)?.city}
                      </p>
                    </>
                  ) : (
                    <p className="text-sm text-gray-500">Chưa có thông tin</p>
                  )}
                </div>

                <button className="w-full mt-4 px-4 py-2 border border-blue-600 text-blue-600 rounded-lg hover:bg-blue-50">
                  📍 Sao chép địa chỉ
                </button>
              </div>
            </div>

            {/* Payment Information */}
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-lg font-semibold mb-4">Thanh toán</h2>
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Phương thức:</span>
                  <span className="font-medium">
                    {order.payment_method === 'cod' ? 'COD - Thanh toán khi nhận hàng' : 
                     order.payment_method === 'card' ? 'Thẻ tín dụng/ghi nợ' :
                     order.payment_method === 'momo' ? 'Ví MoMo' :
                     order.payment_method === 'vnpay' ? 'VNPay' :
                     order.payment_method}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-gray-600">Trạng thái:</span>
                  <span className={`font-medium ${
                    order.payment_status === 'paid' ? 'text-green-600' : 
                    order.payment_status === 'pending' ? 'text-yellow-600' : 
                    'text-gray-600'
                  }`}>
                    {order.payment_status === 'paid' ? 'Đã thanh toán' : 
                     order.payment_status === 'pending' ? 'Chờ thanh toán' : 
                     'Chưa thanh toán'}
                  </span>
                </div>
              </div>
            </div>

            {/* Actions */}
            <div className="space-y-3">
              <button
                onClick={handleBuyAgain}
                className="w-full px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium"
              >
                Mua lại sản phẩm
              </button>
              <button className="w-full px-6 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 font-medium">
                Gửi yêu cầu hỗ trợ
              </button>
              <Link
                href="/orders"
                className="block w-full px-6 py-3 text-center border border-gray-300 rounded-lg hover:bg-gray-50 font-medium"
              >
                ← Quay lại đơn hàng
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
