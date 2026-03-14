import React, { useState } from 'react';
import { JWTDecode, JWTEncode, JWTVerify } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

export default function JWTTool() {
  const [payload, setPayload] = useSessionState('jwt.payload', '{\n  "sub": "1234567890",\n  "name": "John Doe",\n  "iat": 1516239022\n}');
  const [secret, setSecret] = useSessionState('jwt.secret', 'your-256-bit-secret');
  const [algorithm, setAlgorithm] = useSessionState('jwt.algorithm', 'HS256');
  const [token, setToken] = useSessionState('jwt.token', '');
  const [header, setHeader] = useSessionState('jwt.header', '');
  const [decodedPayload, setDecodedPayload] = useSessionState('jwt.decodedPayload', '');
  const [verifySecret, setVerifySecret] = useSessionState('jwt.verifySecret', '');
  const [verificationStatus, setVerificationStatus] = useState('');
  const [error, setError] = useState('');

  const handleEncode = async () => {
    try {
      const encoded = await JWTEncode(payload, secret, algorithm);
      setToken(encoded);
      setError('');
      setVerificationStatus('');
    } catch (e) {
      setToken('');
      setError(`エンコードに失敗しました: ${String(e)}`);
    }
  };

  const handleDecode = async () => {
    try {
      const result = await JWTDecode(token);
      setHeader(result.header);
      setDecodedPayload(result.payload);
      setError('');
      setVerificationStatus('');
    } catch (e) {
      setHeader('');
      setDecodedPayload('');
      setError(`デコードに失敗しました: ${String(e)}`);
    }
  };

  const handleVerify = async () => {
    try {
      const result = await JWTVerify(token, verifySecret);
      setHeader(result.header);
      setDecodedPayload(result.payload);
      if (result.valid) {
        setVerificationStatus('✓ 署名検証成功');
        setError('');
      } else {
        setVerificationStatus('✗ 署名検証失敗');
        setError(result.error || '署名が無効です');
      }
    } catch (e) {
      setHeader('');
      setDecodedPayload('');
      setVerificationStatus('✗ 検証エラー');
      setError(`検証に失敗しました: ${String(e)}`);
    }
  };

  const handleClear = () => {
    setPayload('{\n  "sub": "1234567890",\n  "name": "John Doe",\n  "iat": 1516239022\n}');
    setSecret('your-256-bit-secret');
    setAlgorithm('HS256');
    setToken('');
    setHeader('');
    setDecodedPayload('');
    setVerifySecret('');
    setVerificationStatus('');
    setError('');
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-full max-w-4xl">
        <h1 className="text-3xl font-bold mb-6 text-center">JWT 検証</h1>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* エンコードセクション */}
          <div>
            <h2 className="text-xl font-semibold mb-4">エンコード</h2>

            <div className="mb-4">
              <label htmlFor="jwt-payload" className="block text-sm font-semibold mb-2">
                ペイロード (JSON)
              </label>
              <textarea
                id="jwt-payload"
                value={payload}
                onChange={(e) => setPayload(e.target.value)}
                className="w-full h-40 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono text-sm"
                placeholder='{"sub": "1234567890", "name": "John Doe"}'
              />
            </div>

            <div className="mb-4">
              <label htmlFor="jwt-secret" className="block text-sm font-semibold mb-2">
                秘密鍵 (HMAC) / 秘密鍵 (RS/ES/PS/EdDSA - PEM形式)
              </label>
              <textarea
                id="jwt-secret"
                value={secret}
                onChange={(e) => setSecret(e.target.value)}
                className="w-full h-32 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono text-xs"
                placeholder="HMAC: your-256-bit-secret&#10;&#10;RS/ES/PS/EdDSA:&#10;-----BEGIN PRIVATE KEY-----&#10;...&#10;-----END PRIVATE KEY-----"
              />
            </div>

            <div className="mb-4">
              <label htmlFor="jwt-algorithm" className="block text-sm font-semibold mb-2">
                アルゴリズム
              </label>
              <select
                id="jwt-algorithm"
                value={algorithm}
                onChange={(e) => setAlgorithm(e.target.value)}
                className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
              >
                <optgroup label="HMAC">
                  <option value="HS256">HS256</option>
                  <option value="HS384">HS384</option>
                  <option value="HS512">HS512</option>
                </optgroup>
                <optgroup label="RSA">
                  <option value="RS256">RS256</option>
                  <option value="RS384">RS384</option>
                  <option value="RS512">RS512</option>
                </optgroup>
                <optgroup label="RSA-PSS">
                  <option value="PS256">PS256</option>
                  <option value="PS384">PS384</option>
                  <option value="PS512">PS512</option>
                </optgroup>
                <optgroup label="ECDSA">
                  <option value="ES256">ES256</option>
                  <option value="ES384">ES384</option>
                  <option value="ES512">ES512</option>
                </optgroup>
                <optgroup label="EdDSA">
                  <option value="EdDSA">EdDSA(Ed25519)</option>
                </optgroup>
              </select>
            </div>

            <button
              onClick={handleEncode}
              className="w-full px-4 py-2 rounded bg-blue-600 hover:bg-blue-500 transition font-semibold mb-4"
              type="button"
            >
              JWT エンコード
            </button>

            <div>
              <label htmlFor="jwt-token-output" className="block text-sm font-semibold mb-2">
                JWT トークン
              </label>
              <textarea
                id="jwt-token-output"
                value={token}
                onChange={(e) => setToken(e.target.value)}
                className="w-full h-32 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono text-xs break-all"
                placeholder="生成されたJWTトークン"
              />
            </div>
          </div>

          {/* デコードセクション */}
          <div>
            <h2 className="text-xl font-semibold mb-4">デコード</h2>

            <div className="mb-4">
              <label htmlFor="jwt-token-input" className="block text-sm font-semibold mb-2">
                JWT トークン
              </label>
              <textarea
                id="jwt-token-input"
                value={token}
                onChange={(e) => setToken(e.target.value)}
                className="w-full h-32 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono text-xs break-all"
                placeholder="デコードするJWTトークンを入力"
              />
            </div>

            <div className="mb-4">
              <label htmlFor="jwt-verify-secret" className="block text-sm font-semibold mb-2">
                秘密鍵 (HMAC) / 公開鍵 (RS/ES/PS/EdDSA - PEM形式)
              </label>
              <textarea
                id="jwt-verify-secret"
                value={verifySecret}
                onChange={(e) => setVerifySecret(e.target.value)}
                className="w-full h-32 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono text-xs"
                placeholder="HMAC: your-256-bit-secret&#10;&#10;RS/ES/PS/EdDSA:&#10;-----BEGIN PUBLIC KEY-----&#10;...&#10;-----END PUBLIC KEY-----"
              />
            </div>

            <div className="flex gap-2 mb-4">
              <button
                onClick={handleDecode}
                className="flex-1 px-4 py-2 rounded bg-green-600 hover:bg-green-500 transition font-semibold"
                type="button"
              >
                デコード
              </button>
              <button
                onClick={handleVerify}
                className="flex-1 px-4 py-2 rounded bg-purple-600 hover:bg-purple-500 transition font-semibold"
                type="button"
              >
                検証
              </button>
            </div>

            {verificationStatus && (
              <div className={`mb-4 p-2 rounded text-center font-semibold ${verificationStatus.includes('成功')
                ? 'bg-green-700 text-green-100'
                : 'bg-red-700 text-red-100'
                }`}>
                {verificationStatus}
              </div>
            )}

            <div className="mb-4">
              <label htmlFor="jwt-header" className="block text-sm font-semibold mb-2">
                ヘッダー
              </label>
              <textarea
                id="jwt-header"
                value={header}
                readOnly
                className="w-full h-24 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none font-mono text-sm"
                placeholder="デコードされたヘッダー"
              />
            </div>

            <div>
              <label htmlFor="jwt-decoded-payload" className="block text-sm font-semibold mb-2">
                ペイロード
              </label>
              <textarea
                id="jwt-decoded-payload"
                value={decodedPayload}
                readOnly
                className="w-full h-40 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none font-mono text-sm"
                placeholder="デコードされたペイロード"
              />
            </div>
          </div>
        </div>

        {/* エラー表示 */}
        {error && (
          <div className="mt-6 p-3 bg-red-900 border border-red-600 rounded text-red-100">
            {error}
          </div>
        )}

        {/* クリアボタン */}
        <div className="mt-6 flex justify-center">
          <button
            onClick={handleClear}
            className="px-6 py-2 rounded bg-gray-600 hover:bg-gray-500 transition font-semibold"
            type="button"
          >
            クリア
          </button>
        </div>
      </div>
    </div>
  );
}
