// 简单的toast实现
export const toast = {
  success: (message: string, options?: { icon?: string }) => {
    console.log('✅', message)
    // 创建简单的toast通知
    const toastEl = document.createElement('div')
    toastEl.className = 'fixed top-4 right-4 bg-green-500 text-white px-4 py-2 rounded-lg shadow-lg z-50 transition-all duration-300'
    toastEl.textContent = `${options?.icon || '✅'} ${message}`
    document.body.appendChild(toastEl)
    
    setTimeout(() => {
      toastEl.style.opacity = '0'
      setTimeout(() => {
        document.body.removeChild(toastEl)
      }, 300)
    }, 3000)
  },
  
  error: (message: string) => {
    console.log('❌', message)
    // 创建简单的toast通知
    const toastEl = document.createElement('div')
    toastEl.className = 'fixed top-4 right-4 bg-red-500 text-white px-4 py-2 rounded-lg shadow-lg z-50 transition-all duration-300'
    toastEl.textContent = `❌ ${message}`
    document.body.appendChild(toastEl)
    
    setTimeout(() => {
      toastEl.style.opacity = '0'
      setTimeout(() => {
        document.body.removeChild(toastEl)
      }, 300)
    }, 3000)
  },

  info: (message: string) => {
    console.log('ℹ️', message)
    // 创建简单的toast通知
    const toastEl = document.createElement('div')
    toastEl.className = 'fixed top-4 right-4 bg-blue-500 text-white px-4 py-2 rounded-lg shadow-lg z-50 transition-all duration-300'
    toastEl.textContent = `ℹ️ ${message}`
    document.body.appendChild(toastEl)
    
    setTimeout(() => {
      toastEl.style.opacity = '0'
      setTimeout(() => {
        document.body.removeChild(toastEl)
      }, 300)
    }, 3000)
  }
}