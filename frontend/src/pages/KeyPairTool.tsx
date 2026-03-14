import React, { useState } from 'react';
import { ExtractPublicKey, VerifyKeyPair, GeneratePrivateKey } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

export default function KeyPairTool() {
  const [privateKey, setPrivateKey] = useSessionState('keypair.privateKey', '-----BEGIN PRIVATE KEY-----\n\n-----END PRIVATE KEY-----');
  const [publicKey, setPublicKey] = useSessionState('keypair.publicKey', '');
  const [inputPublicKey, setInputPublicKey] = useSessionState('keypair.inputPublicKey', '');
  const [algorithm, setAlgorithm] = useSessionState('keypair.algorithm', 'RS256');
  const [error, setError] = useState('');
  const [verificationResult, setVerificationResult] = useState<{ valid: boolean; message: string } | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const handleGeneratePrivateKey = async () => {
    try {
      setIsLoading(true);
      setError('');
      setVerificationResult(null);
      const generatedPrivateKey = await GeneratePrivateKey(algorithm);
      setPrivateKey(generatedPrivateKey);
      setPublicKey(''); // 公開鍵をクリア
    } catch (e) {
      setError(`秘密鍵の生成に失敗しました: ${String(e)}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleExtractPublicKey = async () => {
    try {
      setIsLoading(true);
      setError('');
      const extractedPublicKey = await ExtractPublicKey(privateKey, algorithm);
      setPublicKey(extractedPublicKey);
    } catch (e) {
      setPublicKey('');
      setError(`公開鍵の抽出に失敗しました: ${String(e)}`);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCopyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text).then(() => {
      alert(`${label}をクリップボードにコピーしました`);
    }).catch(() => {
      alert(`${label}のコピーに失敗しました`);
    });
  };

  const handleVerifyKeyPair = async () => {
    try {
      setIsLoading(true);
      setError('');
      setVerificationResult(null);
      const isValid = await VerifyKeyPair(privateKey, inputPublicKey, algorithm);
      if (isValid) {
        setVerificationResult({ valid: true, message: '✓ 秘密鍵と公開鍵のペアは一致しています' });
      } else {
        setVerificationResult({ valid: false, message: '✗ 秘密鍵と公開鍵のペアが一致していません' });
      }
    } catch (e) {
      setError(`検証に失敗しました: ${String(e)}`);
      setVerificationResult(null);
    } finally {
      setIsLoading(false);
    }
  };

  const handleClear = () => {
    setPrivateKey('-----BEGIN PRIVATE KEY-----\n\n-----END PRIVATE KEY-----');
    setPublicKey('');
    setInputPublicKey('');
    setAlgorithm('RS256');
    setError('');
    setVerificationResult(null);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-full max-w-4xl">
        <h1 className="text-3xl font-bold mb-6 text-center">鍵ペア作成/検証</h1>

        {/* 秘密鍵セクション */}
        <div className="mb-8">
          <h2 className="text-xl font-semibold mb-4">秘密鍵</h2>

          <div className="mb-4">
            <label htmlFor="algorithm" className="block text-sm font-semibold mb-2">
              アルゴリズム
            </label>
            <select
              id="algorithm"
              value={algorithm}
              onChange={(e) => setAlgorithm(e.target.value)}
              className="w-full p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
            >
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

          <div className="mb-4">
            <button
              onClick={handleGeneratePrivateKey}
              disabled={isLoading}
              className="w-full px-4 py-2 rounded bg-indigo-600 hover:bg-indigo-500 disabled:bg-gray-500 disabled:cursor-not-allowed transition font-semibold mb-4"
              type="button"
            >
              {isLoading ? '生成中...' : '秘密鍵を生成'}
            </button>
          </div>

          <div className="mb-4">
            <label htmlFor="private-key" className="block text-sm font-semibold mb-2">
              秘密鍵 (PEM形式)
            </label>
            <textarea
              id="private-key"
              value={privateKey}
              onChange={(e) => setPrivateKey(e.target.value)}
              className="w-full h-64 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono text-xs"
              placeholder="-----BEGIN PRIVATE KEY-----&#10;...&#10;-----END PRIVATE KEY-----"
            />
          </div>

          {privateKey && (
            <button
              onClick={() => handleCopyToClipboard(privateKey, '秘密鍵')}
              className="w-full px-4 py-2 rounded bg-gray-600 hover:bg-gray-500 transition font-semibold text-sm"
              type="button"
            >
              秘密鍵をコピー
            </button>
          )}
        </div>

        {/* 公開鍵セクション */}
        <div className="mb-8">
          <h2 className="text-xl font-semibold mb-4">公開鍵</h2>

          <div className="mb-4">
            <button
              onClick={handleExtractPublicKey}
              disabled={isLoading || !privateKey.trim()}
              className="w-full px-4 py-2 rounded bg-blue-600 hover:bg-blue-500 disabled:bg-gray-500 disabled:cursor-not-allowed transition font-semibold"
              type="button"
            >
              {isLoading ? '処理中...' : '公開鍵を生成'}
            </button>
          </div>

          <div className="mb-4">
            <label htmlFor="public-key" className="block text-sm font-semibold mb-2">
              公開鍵 (PEM形式)
            </label>
            <textarea
              id="public-key"
              value={publicKey}
              readOnly
              className="w-full h-64 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none font-mono text-xs"
              placeholder="秘密鍵から生成される公開鍵"
            />
          </div>

          {publicKey && (
            <button
              onClick={() => handleCopyToClipboard(publicKey, '公開鍵')}
              className="w-full px-4 py-2 rounded bg-green-600 hover:bg-green-500 transition font-semibold text-sm"
              type="button"
            >
              公開鍵をコピー
            </button>
          )}
        </div>

        {/* 鍵ペア検証セクション */}
        <div className="pt-8 border-t border-gray-700">
          <h2 className="text-xl font-semibold mb-4">鍵ペア検証</h2>
          <p className="text-sm text-gray-300 mb-4">秘密鍵と公開鍵のペアが一致しているかを検証します</p>

          <div className="mb-4">
            <label htmlFor="input-public-key" className="block text-sm font-semibold mb-2">
              検証する公開鍵 (PEM形式)
            </label>
            <textarea
              id="input-public-key"
              value={inputPublicKey}
              onChange={(e) => setInputPublicKey(e.target.value)}
              className="w-full h-40 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono text-xs"
              placeholder="-----BEGIN PUBLIC KEY-----&#10;...&#10;-----END PUBLIC KEY-----"
            />
          </div>

          <button
            onClick={handleVerifyKeyPair}
            disabled={isLoading || !privateKey.trim() || !inputPublicKey.trim()}
            className="w-full px-4 py-2 rounded bg-purple-600 hover:bg-purple-500 disabled:bg-gray-500 disabled:cursor-not-allowed transition font-semibold mb-4"
            type="button"
          >
            {isLoading ? '検証中...' : 'ペアを検証'}
          </button>

          {verificationResult && (
            <div className={`p-3 rounded border ${verificationResult.valid
              ? 'bg-green-900 border-green-600 text-green-100'
              : 'bg-red-900 border-red-600 text-red-100'
              }`}>
              {verificationResult.message}
            </div>
          )}
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
