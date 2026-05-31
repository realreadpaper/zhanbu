import { useAuthStore } from '../stores/authStore'

export default function Profile() {
  const { user } = useAuthStore()

  return (
    <div className="flex-1 flex items-center justify-center px-4 py-20">
      <div className="text-center max-w-md">
        <div className="w-20 h-20 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white text-3xl font-bold mx-auto mb-4">
          {user?.username?.[0]?.toUpperCase() || '?'}
        </div>
        <h1 className="text-3xl font-bold text-white mb-2">{user?.username}</h1>
        <p className="text-slate-400">{user?.email}</p>
        <p className="text-slate-500 text-sm mt-8">🚧 更多个人设置即将上线</p>
      </div>
    </div>
  )
}
