export default function Footer() {
  return (
    <footer className="mt-auto border-t border-slate-800/80 bg-darker/60">
      <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <p className="text-center text-sm leading-6 text-slate-500">
          © {new Date().getFullYear()}{' '}
          <span className="font-medium text-slate-300">玄机占卜</span>
          <span className="mx-2 text-slate-700">·</span>
          <span className="text-slate-400">心有所问，自有玄机</span>
        </p>
      </div>
    </footer>
  )
}
