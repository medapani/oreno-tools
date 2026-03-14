import { NavLink, Outlet } from 'react-router-dom';

export default function Layout() {
  const baseMenuClass = 'block p-3 rounded transition';

  return (
    <div className="flex h-screen">
      {/* サイドバー */}
      <aside className="w-64 bg-gray-800 text-white flex flex-col">
        <div className="p-6 text-2xl font-bold border-b border-gray-700">oreno tools</div>
        <nav className="flex-1 p-4 space-y-2 overflow-auto">
          <NavLink
            to="/"
            end
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            バイト変換
          </NavLink>
          <NavLink
            to="/network-unit"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            通信速度変換
          </NavLink>
          <NavLink
            to="/data-transfer-rate"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            データ転送速度変換
          </NavLink>
          <NavLink
            to="/transfer-time"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            データ転送時間計算
          </NavLink>
          <NavLink
            to="/cidr-calculator"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            CIDR 計算
          </NavLink>
          <NavLink
            to="/base64"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            Base64 変換
          </NavLink>
          <NavLink
            to="/jwt"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            JWT 検証
          </NavLink>
          <NavLink
            to="/json-formatter"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            JSON/YAML 変換
          </NavLink>
          <NavLink
            to="/url-tool"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            URL エンコード/デコード
          </NavLink>
          <NavLink
            to="/base-converter"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            基数変換/計算
          </NavLink>
          <NavLink
            to="/key-pair"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            鍵ペア作成/検証
          </NavLink>
          <NavLink
            to="/self-signed-cert"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            自己署名証明書作成
          </NavLink>
          <NavLink
            to="/mtls-cert"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            mTLS 証明書作成
          </NavLink>
          <NavLink
            to="/crl-tool"
            className={({ isActive }) => `${baseMenuClass} ${isActive ? 'bg-blue-600 text-white' : 'hover:bg-gray-700 text-gray-100'}`}
          >
            CRL 更新
          </NavLink>

        </nav>
      </aside>

      {/* メインコンテンツ表示エリア */}
      <main className="flex-1 overflow-auto">
        <Outlet />
      </main>
    </div>
  );
}