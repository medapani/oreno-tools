import React, { useState } from 'react';
import { useSessionState } from '../hooks/useSessionState';

export default function URLTool() {
  const [input, setInput] = useSessionState('url.input', '');
  const [output, setOutput] = useSessionState('ur.output', '');
  const [error, setError] = useState('');
  const [encodeComponent, setEncodeComponent] = useSessionState('url.encodeComponent', true);

  const handleEncode = () => {
    try {
      const result = encodeComponent ? encodeURIComponent(input) : encodeURI(input);
      setOutput(result);
      setError('');
    } catch (e) {
      setOutput('');
      if (e instanceof Error) {
        setError(`エンコードに失敗しました: ${e.message}`);
      } else {
        setError(`エンコードに失敗しました: ${String(e)}`);
      }
    }
  };

  const handleDecode = () => {
    try {
      const result = decodeURIComponent(input);
      setOutput(result);
      setError('');
    } catch (e) {
      setOutput('');
      if (e instanceof Error) {
        setError(`デコードに失敗しました: ${e.message}`);
      } else {
        setError(`デコードに失敗しました: ${String(e)}`);
      }
    }
  };

  const handleClear = () => {
    setInput('');
    setOutput('');
    setError('');
  };

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(output);
    } catch (e) {
      setError('クリップボードへのコピーに失敗しました');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-full max-w-4xl">
        <h1 className="text-3xl font-bold mb-6 text-center">URL エンコード/デコード</h1>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* 入力セクション */}
          <div>
            <label htmlFor="url-input" className="block text-sm font-semibold mb-2">
              入力
            </label>
            <textarea
              id="url-input"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              className="w-full h-64 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono text-sm"
              placeholder="ここにテキストまたはURLエンコード文字列を入力"
            />
          </div>

          {/* 出力セクション */}
          <div>
            <div className="flex justify-between items-center mb-2">
              <label htmlFor="url-output" className="block text-sm font-semibold">
                出力
              </label>
              {output && (
                <button
                  onClick={handleCopy}
                  className="px-3 py-1 rounded bg-gray-600 hover:bg-gray-500 transition text-sm font-semibold"
                  type="button"
                  title="クリップボードにコピー"
                >
                  コピー
                </button>
              )}
            </div>
            <textarea
              id="url-output"
              value={output}
              readOnly
              className="w-full h-64 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none font-mono text-sm"
              placeholder="エンコード/デコード結果がここに表示されます"
            />
          </div>
        </div>

        {/* オプション */}
        <div className="mt-6 mb-6">
          <label className="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={encodeComponent}
              onChange={(e) => setEncodeComponent(e.target.checked)}
              className="w-4 h-4 rounded bg-gray-700 border-gray-600 text-blue-600 focus:ring-blue-500 focus:ring-2"
            />
            <span className="text-sm font-semibold">
              encodeURIComponent を使用（クエリパラメータ推奨）
            </span>
          </label>
          <p className="text-xs text-gray-400 mt-2 ml-6">
            チェック: パラメータ値用 | チェック解除: URL全体用
          </p>
        </div>

        {/* エラー表示 */}
        {error && (
          <div className="mb-6 p-3 bg-red-900 border border-red-600 rounded text-red-100">
            {error}
          </div>
        )}

        {/* アクションボタン */}
        <div className="flex flex-wrap gap-3 justify-center">
          <button
            onClick={handleEncode}
            className="px-6 py-2 rounded bg-blue-600 hover:bg-blue-500 transition font-semibold"
            type="button"
          >
            URL エンコード
          </button>
          <button
            onClick={handleDecode}
            className="px-6 py-2 rounded bg-green-600 hover:bg-green-500 transition font-semibold"
            type="button"
          >
            URL デコード
          </button>
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
