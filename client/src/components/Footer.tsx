export default function Footer() {
  return (
    <footer className="bg-darker/50 border-t border-slate-800 mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2">
            <span className="text-2xl">🔮</span>
            <span className="text-lg font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
              玄机占卜
            </span>
          </div>
          <p className="text-slate-500 text-sm">
            © {new Date().getFullYear()} 玄机占卜 · 探索命运的奥秘 · 仅供娱乐参考
          </p>
        </div>
      </div>
    </footer>
  )
}
