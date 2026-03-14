import React, { useState } from 'react';
import { Base64Decode, Base64Encode } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

export default function Base64Tool() {
  const [input, setInput] = useSessionState('base64.input', '');
  const [output, setOutput] = useSessionState('base64.output', '');
  const [error, setError] = useState('');
  const [urlSafe, setUrlSafe] = useSessionState('base64.urlSafe', true);

  const handleEncode = async () => {
    try {
      const encoded = await Base64Encode(input, urlSafe);
      setOutput(encoded);
      setError('');
    } catch (e) {
      setOutput('');
      setError(`エンコードに失敗しました: ${String(e)}`);
    }
  };

  const handleDecode = async () => {
    try {
      const decoded = await Base64Decode(input, urlSafe);
      setOutput(decoded);
      setError('');
    } catch (e) {
      setOutput('');
      setError(`デコードに失敗しました: ${String(e)}`);
    }
  };

  const handleClear = () => {
    setInput('');
    setOutput('');
    setError('');
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-full max-w-2xl">
        <h1 className="text-3xl font-bold mb-6 text-center">Base64 変換</h1>

        <div className="mb-4">
          <label htmlFor="base64-input" className="block text-sm font-semibold mb-2">
            入力テキスト
          </label>
          <textarea
            id="base64-input"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            className="w-full h-40 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
            placeholder="ここにテキストまたはBase64文字列を入力"
          />
        </div>

        <div className="mb-4">
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={urlSafe}
              onChange={(e) => setUrlSafe(e.target.checked)}
              className="w-4 h-4 rounded bg-gray-700 border-gray-600 text-blue-600 focus:ring-blue-500 focus:ring-2"
            />
            <span className="text-sm font-semibold">URL-safe Base64 (RFC 4648)</span>
          </label>
        </div>

        <div className="flex flex-wrap gap-3 mb-4">
          <button
            onClick={handleEncode}
            className="px-4 py-2 rounded bg-blue-600 hover:bg-blue-500 transition font-semibold"
            type="button"
          >
            Base64 エンコード
          </button>
          <button
            onClick={handleDecode}
            className="px-4 py-2 rounded bg-green-600 hover:bg-green-500 transition font-semibold"
            type="button"
          >
            Base64 デコード
          </button>
          <button
            onClick={handleClear}
            className="px-4 py-2 rounded bg-gray-600 hover:bg-gray-500 transition font-semibold"
            type="button"
          >
            クリア
          </button>
        </div>

        <div>
          <label htmlFor="base64-output" className="block text-sm font-semibold mb-2">
            結果
          </label>
          <textarea
            id="base64-output"
            value={error || output}
            readOnly
            className={`w-full h-40 p-3 rounded border focus:outline-none ${error ? 'bg-red-900 border-red-600 text-red-100' : 'bg-gray-700 border-gray-600'}`}
            placeholder="変換結果がここに表示されます"
          />
        </div>
      </div>
    </div>
  );
}
