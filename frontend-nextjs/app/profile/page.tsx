'use client'

import { useState, useEffect } from 'react'
import Image from 'next/image'
import { useAuthStore } from '@/store/authStore'
import { useRouter } from 'next/navigation'
import toast from 'react-hot-toast'

export default function ProfilePage() {
  const router = useRouter()
  const { user, isAuthenticated } = useAuthStore()
  const [isEditing, setIsEditing] = useState(false)
  
  // Form states
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [phone, setPhone] = useState('')
  
  // Password change states
  const [currentPassword, setCurrentPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')

  // Activity history mock data
  const activityHistory = [
    { action: 'Đăng nhập tài khoản', description: 'Từ địa chỉ IP 192.168.1.1 trên Chrome/Windows', time: 'Vừa xong', icon: '🔐' },
    { action: 'Cập nhật thông tin tài khoản', description: 'Thay đổi số điện thoại', time: '1 ngày trước', icon: '✏️' },
    { action: 'Thay đổi mật khẩu', description: 'Mật khẩu đã được cập nhật thành công', time: '3 ngày trước', icon: '🔒' },
    { action: 'Đăng nhập trên thiết bị mới', description: 'Đăng nhập từ Safari trên iPhone 15 Pro Max', time: '7 ngày trước', icon: '📱' },
  ]

  // Favorite products mock data
  const favoriteProducts = [
    { id: 1, name: 'Đầm Doux XYZ Pro Max 256GB', price: '29.990.000₫', image: '/images/products/product1.jpg' },
    { id: 2, name: 'Laptop ABC Ultimate™ M3 13.6"', price: '32.990.000₫', image: '/images/products/product2.jpg' },
    { id: 3, name: 'Tai nghe Sony/Wave Pro AirPods', price: '4.490.000₫', image: '/images/products/product3.jpg' },
  ]

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
      return
    }
    
    if (user) {
      setName(user.full_name || user.username || '')
      setEmail(user.email || '')
      setPhone(user.phone || '')
    }
  }, [isAuthenticated, user, router])

  const handleSaveProfile = () => {
    // TODO: Implement API call to update profile
    toast.success('Đã lưu thay đổi thành công')
    setIsEditing(false)
  }

  const handleCancel = () => {
    if (user) {
      setName(user.full_name || user.username || '')
      setEmail(user.email || '')
      setPhone(user.phone || '')
    }
    setIsEditing(false)
  }

  const handleChangePassword = () => {
    if (newPassword !== confirmPassword) {
      toast.error('Mật khẩu xác nhận không khớp')
      return
    }
    
    if (newPassword.length < 6) {
      toast.error('Mật khẩu phải có ít nhất 6 ký tự')
      return
    }

    // TODO: Implement API call to change password
    toast.success('Đổi mật khẩu thành công')
    setCurrentPassword('')
    setNewPassword('')
    setConfirmPassword('')
  }

  const handleActivate2FA = () => {
    toast.info('Tính năng xác thực hai yếu tố sẽ sớm được cập nhật')
  }

  if (!isAuthenticated || !user) {
    return null
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-5xl mx-auto px-4">
        {/* Page Header */}
        <h1 className="text-2xl font-bold mb-6">Tài Khoản Của Tôi</h1>

        {/* Sidebar Navigation */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
          {/* Left Sidebar */}
          <div className="md:col-span-1">
            <div className="bg-white rounded-lg shadow p-4">
              <div className="flex items-center space-x-3 mb-6">
                <div className="w-12 h-12 bg-blue-600 rounded-full flex items-center justify-center text-white font-bold text-xl">
                  {(user.full_name || user.username || 'U').charAt(0).toUpperCase()}
                </div>
                <div>
                  <p className="font-semibold">{user.full_name || user.username}</p>
                  <p className="text-sm text-gray-500">{user.email}</p>
                </div>
              </div>

              <nav className="space-y-2">
                <button className="w-full text-left px-4 py-2 bg-blue-50 text-blue-600 rounded-lg font-medium">
                  👤 Thông tin tài khoản
                </button>
                <button className="w-full text-left px-4 py-2 hover:bg-gray-50 rounded-lg text-gray-700">
                  📦 Đơn hàng của tôi
                </button>
                <button className="w-full text-left px-4 py-2 hover:bg-gray-50 rounded-lg text-gray-700">
                  📍 Sổ địa chỉ
                </button>
                <button className="w-full text-left px-4 py-2 hover:bg-gray-50 rounded-lg text-gray-700">
                  ❤️ Sản phẩm yêu thích
                </button>
                <button className="w-full text-left px-4 py-2 hover:bg-gray-50 rounded-lg text-gray-700">
                  🔔 Thông báo của tôi
                </button>
                <button className="w-full text-left px-4 py-2 hover:bg-gray-50 rounded-lg text-gray-700">
                  🎫 Mã giảm giá của tôi
                </button>
                <button className="w-full text-left px-4 py-2 text-red-600 hover:bg-red-50 rounded-lg">
                  🚪 Đăng xuất
                </button>
              </nav>
            </div>
          </div>

          {/* Main Content */}
          <div className="md:col-span-3 space-y-6">
            {/* Profile Picture Section */}
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-6">
                  <div className="w-24 h-24 bg-blue-600 rounded-full flex items-center justify-center text-white font-bold text-4xl">
                    {(user.full_name || user.username || 'U').charAt(0).toUpperCase()}
                  </div>
                  <div>
                    <h3 className="font-semibold text-lg">Ảnh đại diện</h3>
                    <p className="text-sm text-gray-500">Vừa và nhỏ, không quá 500kb</p>
                  </div>
                </div>
                <button className="px-4 py-2 text-blue-600 border border-blue-600 rounded-lg hover:bg-blue-50">
                  Thay đổi ảnh
                </button>
              </div>
            </div>

            {/* Account Information */}
            <div className="bg-white rounded-lg shadow p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-xl font-semibold">Thông tin tài khoản</h2>
                {!isEditing && (
                  <button
                    onClick={() => setIsEditing(true)}
                    className="text-blue-600 hover:underline"
                  >
                    Chỉnh sửa
                  </button>
                )}
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Họ và Tên
                  </label>
                  <input
                    type="text"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    disabled={!isEditing}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 disabled:bg-gray-50"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Email
                  </label>
                  <input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    disabled={!isEditing}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 disabled:bg-gray-50"
                  />
                </div>
              </div>

              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Số điện thoại
                </label>
                <input
                  type="tel"
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  disabled={!isEditing}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 disabled:bg-gray-50"
                />
              </div>

              {isEditing && (
                <div className="flex justify-end space-x-3">
                  <button
                    onClick={handleCancel}
                    className="px-6 py-2 border border-gray-300 rounded-lg hover:bg-gray-50"
                  >
                    Hủy
                  </button>
                  <button
                    onClick={handleSaveProfile}
                    className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                  >
                    Lưu thay đổi
                  </button>
                </div>
              )}
            </div>

            {/* Activity History */}
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-xl font-semibold mb-4">Lịch sử hoạt động tài khoản</h2>
              <p className="text-sm text-gray-500 mb-4">
                Theo dõi hoạt động của tài khoản để đảm bảo an toàn tài khoản của bạn
              </p>
              
              <div className="space-y-4">
                {activityHistory.map((activity, index) => (
                  <div key={index} className="flex items-start space-x-4 pb-4 border-b last:border-b-0">
                    <div className="text-2xl">{activity.icon}</div>
                    <div className="flex-1">
                      <h3 className="font-medium">{activity.action}</h3>
                      <p className="text-sm text-gray-500">{activity.description}</p>
                    </div>
                    <span className="text-sm text-gray-400 whitespace-nowrap">{activity.time}</span>
                  </div>
                ))}
              </div>
            </div>

            {/* Favorite Products */}
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-xl font-semibold mb-4">Sản phẩm yêu thích</h2>
              <p className="text-sm text-gray-500 mb-4">
                Những sản phẩm yêu thích của bạn được lưu tại đây để dễ dàng tìm kiếm
              </p>
              
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {favoriteProducts.map((product) => (
                  <div key={product.id} className="border rounded-lg overflow-hidden">
                    <div className="relative h-48 bg-gray-100">
                      <Image
                        src={product.image}
                        alt={product.name}
                        fill
                        className="object-cover"
                      />
                    </div>
                    <div className="p-4">
                      <h3 className="font-medium text-sm mb-2 line-clamp-2">{product.name}</h3>
                      <p className="text-red-600 font-bold mb-3">{product.price}</p>
                      <div className="flex space-x-2">
                        <button className="flex-1 bg-blue-600 text-white py-2 rounded-lg text-sm hover:bg-blue-700">
                          Xem ngay
                        </button>
                        <button className="p-2 border border-red-300 text-red-500 rounded-lg hover:bg-red-50">
                          ❤️
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Password & Security */}
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-xl font-semibold mb-4">Mật khẩu và bảo mật</h2>
              
              <div className="mb-6">
                <h3 className="font-medium mb-3">Thay đổi mật khẩu</h3>
                <div className="space-y-3">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Mật khẩu hiện tại
                    </label>
                    <input
                      type="password"
                      value={currentPassword}
                      onChange={(e) => setCurrentPassword(e.target.value)}
                      placeholder="••••••••"
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Mật khẩu mới
                    </label>
                    <input
                      type="password"
                      value={newPassword}
                      onChange={(e) => setNewPassword(e.target.value)}
                      placeholder="••••••••"
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Xác nhận mật khẩu mới
                    </label>
                    <input
                      type="password"
                      value={confirmPassword}
                      onChange={(e) => setConfirmPassword(e.target.value)}
                      placeholder="••••••••"
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                    />
                  </div>
                  <button
                    onClick={handleChangePassword}
                    className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                  >
                    Lưu mật khẩu
                  </button>
                </div>
              </div>

              <div className="border-t pt-6">
                <h3 className="font-medium mb-2">Xác thực hai yếu tố</h3>
                <p className="text-sm text-gray-500 mb-4">
                  Bảo vệ tài khoản bằng mã xác thực 2 lớp. Kích hoạt để tăng cường bảo mật
                </p>
                <button
                  onClick={handleActivate2FA}
                  className="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
                >
                  Kích hoạt
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
