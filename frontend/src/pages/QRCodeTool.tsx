import React, { useState, useEffect, useRef } from 'react';
import { GenerateQRCode } from '../../wailsjs/go/main/App';
import { useSessionState } from '../hooks/useSessionState';

const SIZE_OPTIONS = [128, 256, 512] as const;
const LEVEL_OPTIONS = [
  { value: 'L', label: 'L（低 7%）' },
  { value: 'M', label: 'M（中 15%）' },
  { value: 'Q', label: 'Q（高 25%）' },
  { value: 'H', label: 'H（最高 30%）' },
] as const;

export default function QRCodeTool() {
  const [input, setInput] = useSessionState('qrcode.input', '');
  const [size, setSize] = useSessionState<number>('qrcode.size', 256);
  const [level, setLevel] = useSessionState('qrcode.level', 'M');
  const [qrImage, setQrImage] = useState('');
  const [error, setError] = useState('');
  const requestIdRef = useRef(0);

  useEffect(() => {
    if (!input.trim()) {
      setQrImage('');
      setError('');
      return;
    }

    const currentId = ++requestIdRef.current;
    const timer = setTimeout(async () => {
      try {
        const result = await GenerateQRCode(input, size, level);
        if (currentId === requestIdRef.current) {
          setQrImage(result);
          setError('');
        }
      } catch (e) {
        if (currentId === requestIdRef.current) {
          setQrImage('');
          setError(String(e));
        }
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [input, size, level]);

  const handleDownload = () => {
    if (!qrImage) return;
    const link = document.createElement('a');
    link.href = `data:image/png;base64,${qrImage}`;
    link.download = 'qrcode.png';
    link.click();
  };

  const handleClear = () => {
    setInput('');
    setQrImage('');
    setError('');
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-900 text-white p-4">
      <div className="bg-gray-800 p-10 rounded-xl shadow-2xl w-full max-w-2xl">
        <h1 className="text-3xl font-bold mb-6 text-center">QR コード生成</h1>

        <div className="mb-4">
          <label htmlFor="qr-input" className="block text-sm font-semibold mb-2">
            テキスト / URL
          </label>
          <textarea
            id="qr-input"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            className="w-full h-28 p-3 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500 font-mono text-sm"
            placeholder="QR コードにするテキストまたは URL を入力"
          />
        </div>

        <div className="flex flex-wrap gap-6 mb-6">
          <div>
            <label className="block text-sm font-semibold mb-2">サイズ (px)</label>
            <select
              value={size}
              onChange={(e) => setSize(Number(e.target.value))}
              className="p-2 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
            >
              {SIZE_OPTIONS.map((s) => (
                <option key={s} value={s}>{s} × {s}</option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-semibold mb-2">誤り訂正レベル</label>
            <select
              value={level}
              onChange={(e) => setLevel(e.target.value)}
              className="p-2 rounded bg-gray-700 border border-gray-600 focus:outline-none focus:border-blue-500"
            >
              {LEVEL_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value}>{opt.label}</option>
              ))}
            </select>
          </div>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-900 border border-red-600 rounded text-red-100 text-sm">
            {error}
          </div>
        )}

        {qrImage && (
          <div className="flex flex-col items-center gap-4">
            <img
              src={`data:image/png;base64,${qrImage}`}
              alt="Generated QR Code"
              className="rounded border border-gray-600 bg-white p-2"
              style={{ width: size, height: size }}
            />
            <div className="flex gap-3">
              <button
                onClick={handleDownload}
                className="px-5 py-2 rounded bg-blue-600 hover:bg-blue-500 transition font-semibold"
                type="button"
              >
                PNG ダウンロード
              </button>
              <button
                onClick={handleClear}
                className="px-5 py-2 rounded bg-gray-600 hover:bg-gray-500 transition font-semibold"
                type="button"
              >
                クリア
              </button>
            </div>
          </div>
        )}

        {!qrImage && !error && (
          <div className="flex justify-center mt-4">
            <button
              onClick={handleClear}
              className="px-5 py-2 rounded bg-gray-600 hover:bg-gray-500 transition font-semibold"
              type="button"
            >
              クリア
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
